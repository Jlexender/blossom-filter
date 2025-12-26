package core

import (
	"testing"
)

// Table-driven tests for Generic BloomFilter

func TestBloomFilter_New(t *testing.T) {
	tests := []struct {
		name      string
		size      uint32
		wantPanic bool
		wantSize  uint32
	}{
		{"size 1", 1, false, 1},
		{"size 16", 16, false, 16},
		{"size 1024", 1024, false, 1024},
		{"size 0 should panic", 0, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewBloomFilter() did not panic, want panic")
					}
				}()
			}
			f := NewBloomFilter[string](tt.size)
			if !tt.wantPanic {
				if f == nil {
					t.Fatal("NewBloomFilter() = nil, want non-nil")
				}
				if f.Size() != tt.wantSize {
					t.Errorf("Size() = %d, want %d", f.Size(), tt.wantSize)
				}
				if f.elements != 0 {
					t.Errorf("elements = %d, want 0", f.elements)
				}
			}
		})
	}
}

func TestBloomFilter_StringType(t *testing.T) {
	tests := []struct {
		name         string
		size         uint32
		insertData   []string
		checkData    string
		wantContains bool
	}{
		{"single string", 16, []string{"hello"}, "hello", true},
		{"string not inserted", 16, []string{"hello"}, "world", false},
		{"multiple strings", 64, []string{"foo", "bar", "baz"}, "bar", true},
		{"empty string", 16, []string{""}, "", true},
		{"unicode strings", 32, []string{"hello", "世界", "🌍"}, "世界", true},
		{"empty filter", 16, []string{}, "anything", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewBloomFilter[string](tt.size)
			for _, data := range tt.insertData {
				f.Insert(data)
			}
			got := f.Contains(tt.checkData)
			if got != tt.wantContains {
				t.Errorf("Contains(%q) = %v, want %v", tt.checkData, got, tt.wantContains)
			}
		})
	}
}

func TestBloomFilter_IntType(t *testing.T) {
	tests := []struct {
		name         string
		size         uint32
		insertData   []int
		checkData    int
		wantContains bool
	}{
		{"single int", 16, []int{42}, 42, true},
		{"int not inserted", 16, []int{42}, 43, false},
		{"multiple ints", 64, []int{1, 2, 3, 4, 5}, 3, true},
		{"zero value", 16, []int{0}, 0, true},
		{"negative numbers", 32, []int{-10, -5, 0, 5, 10}, -5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewBloomFilter[int](tt.size)
			for _, data := range tt.insertData {
				f.Insert(data)
			}
			got := f.Contains(tt.checkData)
			if got != tt.wantContains {
				t.Errorf("Contains(%d) = %v, want %v", tt.checkData, got, tt.wantContains)
			}
		})
	}
}

func TestBloomFilter_StructType(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	tests := []struct {
		name         string
		size         uint32
		insertData   []Person
		checkData    Person
		wantContains bool
	}{
		{"single struct", 32, []Person{{Name: "Alice", Age: 30}}, Person{Name: "Alice", Age: 30}, true},
		{"different struct", 32, []Person{{Name: "Alice", Age: 30}}, Person{Name: "Bob", Age: 30}, false},
		{"same name different age", 32, []Person{{Name: "Alice", Age: 30}}, Person{Name: "Alice", Age: 31}, false},
		{
			"multiple structs", 64,
			[]Person{
				{Name: "Alice", Age: 30},
				{Name: "Bob", Age: 25},
				{Name: "Charlie", Age: 35},
			},
			Person{Name: "Bob", Age: 25}, true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewBloomFilter[Person](tt.size)
			for _, data := range tt.insertData {
				f.Insert(data)
			}
			got := f.Contains(tt.checkData)
			if got != tt.wantContains {
				t.Errorf("Contains(%+v) = %v, want %v", tt.checkData, got, tt.wantContains)
			}
		})
	}
}

func TestBloomFilter_Duplicates(t *testing.T) {
	tests := []struct {
		name         string
		size         uint32
		insertData   []string
		wantElements uint32
	}{
		{"insert same string twice", 16, []string{"hello", "hello"}, 1},
		{"insert different strings", 32, []string{"hello", "world"}, 2},
		{"insert same value multiple times", 32, []string{"x", "x", "x", "x"}, 1},
		{"insert multiple with some duplicates", 64, []string{"a", "b", "a", "c", "b"}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewBloomFilter[string](tt.size)
			for _, data := range tt.insertData {
				f.Insert(data)
			}
			if f.elements != tt.wantElements {
				t.Errorf("elements = %d, want %d", f.elements, tt.wantElements)
			}
			// Verify all unique inserted items are contained
			seen := make(map[string]bool)
			for _, data := range tt.insertData {
				if !seen[data] {
					if !f.Contains(data) {
						t.Errorf("Contains(%q) = false, want true", data)
					}
					seen[data] = true
				}
			}
		})
	}
}

func TestBloomFilter_HashFunctionUpdates(t *testing.T) {
	tests := []struct {
		name           string
		size           uint32
		insertSequence []string
		wantHashCounts []int
	}{
		{
			"hash count decreases with elements",
			16,
			[]string{"first", "second"},
			[]int{11, 5},
		},
		{
			"single insertion",
			16,
			[]string{"single"},
			[]int{11},
		},
		{
			"multiple insertions",
			32,
			[]string{"a", "b", "c", "d"},
			[]int{22, 11, 7, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewBloomFilter[string](tt.size)
			for i, data := range tt.insertSequence {
				f.Insert(data)
				if len(f.hashes) != tt.wantHashCounts[i] {
					t.Errorf("after insert %d: hash count = %d, want %d", i, len(f.hashes), tt.wantHashCounts[i])
				}
			}
		})
	}
}

func TestBloomFilter_EmptyFilter(t *testing.T) {
	t.Run("empty string filter", func(t *testing.T) {
		f := NewBloomFilter[string](16)
		if f.Contains("anything") {
			t.Error("empty filter should not contain any elements")
		}
	})

	t.Run("empty int filter", func(t *testing.T) {
		f := NewBloomFilter[int](16)
		if f.Contains(42) {
			t.Error("empty filter should not contain any elements")
		}
	})
}

func TestBloomFilter_LargeDataset(t *testing.T) {
	tests := []struct {
		name     string
		size     uint32
		numItems int
	}{
		{"100 items in size 1024", 1024, 100},
		{"1000 items in size 16384", 16384, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewBloomFilter[string](tt.size)
			
			// Insert items using unique strings
			for i := 0; i < tt.numItems; i++ {
				f.Insert("item_" + string(rune(i)))
			}
			
			// Verify at least some items are contained
			missingCount := 0
			for i := 0; i < tt.numItems; i++ {
				if !f.Contains("item_" + string(rune(i))) {
					missingCount++
					if missingCount > tt.numItems/10 {
						t.Fatalf("Too many items missing: %d/%d", missingCount, tt.numItems)
					}
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkBloomFilter_InsertString(b *testing.B) {
	f := NewBloomFilter[string](8192)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Insert("test string")
	}
}

func BenchmarkBloomFilter_ContainsString(b *testing.B) {
	f := NewBloomFilter[string](8192)
	for i := 0; i < 1000; i++ {
		f.Insert("test")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Contains("test")
	}
}

func BenchmarkBloomFilter_InsertInt(b *testing.B) {
	f := NewBloomFilter[int](8192)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Insert(i)
	}
}

func BenchmarkBloomFilter_ContainsInt(b *testing.B) {
	f := NewBloomFilter[int](8192)
	for i := 0; i < 1000; i++ {
		f.Insert(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Contains(i % 1000)
	}
}

func BenchmarkBloomFilter_Struct(b *testing.B) {
	type Person struct {
		Name string
		Age  int
	}
	f := NewBloomFilter[Person](8192)
	p := Person{Name: "Alice", Age: 30}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Insert(p)
		f.Contains(p)
	}
}
