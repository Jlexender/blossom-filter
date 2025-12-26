package core

import (
	"alex/bvs/internal/bitset"
	"alex/bvs/internal/hash"
	"fmt"
)

// BloomFilter is a type-safe probabilistic data structure for testing set membership.
// The generic type T is constrained to comparable to ensure proper equality checks.
type BloomFilter[T comparable] struct {
	bs       *bitset.Bitset
	hashes   []hash.Hash
	elements uint32
}

// mapToBytes converts a comparable value to bytes for hashing.
func mapToBytes[T comparable](obj T) []byte {
	return []byte(fmt.Sprintf("%T.%v", obj, obj))
}

// NewBloomFilter creates a new type-safe bloom filter with the specified bit size.
// The size must be greater than 0 or it will panic.
func NewBloomFilter[T comparable](size uint32) *BloomFilter[T] {
	if size == 0 {
		panic("size must be greater than 0")
	}

	return &BloomFilter[T]{
		bs:       bitset.NewBitset(size),
		hashes:   hash.NewHashList(size),
		elements: 0,
	}
}

// Insert adds an element to the bloom filter.
// If the element is already present (or appears to be due to hash collisions),
// it will not be added again.
func (bf *BloomFilter[T]) Insert(data T) {
	if bf.Contains(data) {
		return
	}

	bitset := bf.bs
	bf.elements++
	for _, h := range bf.hashes {
		hashsum := h.Compute(mapToBytes(data))
		bitset.Set(hashsum % bitset.Size())
		bf.hashes = hash.UpdateList(bf.hashes, bf.Size(), bf.elements)
	}
}

// Contains checks if an element might be in the bloom filter.
// Returns true if the element might be present (with possible false positives).
// Returns false if the element is definitely not present.
func (bf *BloomFilter[T]) Contains(data T) bool {
	bitset := bf.bs

	for _, h := range bf.hashes {
		hashsum := h.Compute(mapToBytes(data))

		set, _ := bitset.IsSet(hashsum % bitset.Size())
		if !set {
			return false
		}
	}

	return true
}

// Size returns the total bit size of the bloom filter.
func (bf *BloomFilter[T]) Size() uint32 {
	return bf.bs.Size()
}
