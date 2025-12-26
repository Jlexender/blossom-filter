package main

import (
	"alex/bvs/pkg/core"
	"fmt"
)

func main() {
	// Example 1: String filter
	fmt.Println("=== String Filter Example ===")
	stringFilter := core.NewBloomFilter[string](1024)

	// Insert some words
	words := []string{"hello", "world", "bloom", "filter"}
	for _, word := range words {
		stringFilter.Insert(word)
	}

	// Check membership
	fmt.Printf("Contains 'hello': %v\n", stringFilter.Contains("hello"))
	fmt.Printf("Contains 'goodbye': %v\n", stringFilter.Contains("goodbye"))

	// Example 2: Integer filter
	fmt.Println("\n=== Integer Filter Example ===")
	intFilter := core.NewBloomFilter[int](2048)

	// Insert some numbers
	for i := 0; i < 100; i += 10 {
		intFilter.Insert(i)
	}

	fmt.Printf("Contains 50: %v\n", intFilter.Contains(50))
	fmt.Printf("Contains 51: %v\n", intFilter.Contains(51))

	// Example 3: Custom struct filter
	fmt.Println("\n=== Struct Filter Example ===")
	type Person struct {
		Name string
		Age  int
	}

	personFilter := core.NewBloomFilter[Person](4096)

	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}

	for _, person := range people {
		personFilter.Insert(person)
	}

	fmt.Printf("Contains Alice(30): %v\n", personFilter.Contains(Person{Name: "Alice", Age: 30}))
	fmt.Printf("Contains Alice(31): %v\n", personFilter.Contains(Person{Name: "Alice", Age: 31}))
	fmt.Printf("Contains Dave(40): %v\n", personFilter.Contains(Person{Name: "Dave", Age: 40}))

	// Example 4: Filter size info
	fmt.Println("\n=== Filter Info ===")
	fmt.Printf("String filter bit size: %d\n", stringFilter.Size())
	fmt.Printf("Integer filter bit size: %d\n", intFilter.Size())
	fmt.Printf("Person filter bit size: %d\n", personFilter.Size())
}
