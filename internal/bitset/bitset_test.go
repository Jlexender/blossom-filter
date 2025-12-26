package bitset

import (
	"testing"
)

func TestBitsetNew(t *testing.T) {
	tests := []struct {
		name     string
		size     uint32
		wantSize uint32
		wantLen  int
	}{
		{
			name:     "size 1",
			size:     1,
			wantSize: 1,
			wantLen:  1,
		},
		{
			name:     "size 8",
			size:     8,
			wantSize: 8,
			wantLen:  1,
		},
		{
			name:     "size 9",
			size:     9,
			wantSize: 9,
			wantLen:  2,
		},
		{
			name:     "size 16",
			size:     16,
			wantSize: 16,
			wantLen:  2,
		},
		{
			name:     "size 100",
			size:     100,
			wantSize: 100,
			wantLen:  13,
		},
		{
			name:     "size 1024",
			size:     1024,
			wantSize: 1024,
			wantLen:  128,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := NewBitset(tt.size)
			if bs == nil {
				t.Fatal("expected bitset not to be nil")
			}
			if bs.Size() != tt.wantSize {
				t.Errorf("Size() = %d, want %d", bs.Size(), tt.wantSize)
			}
			if len(bs.List()) != tt.wantLen {
				t.Errorf("len(List()) = %d, want %d", len(bs.List()), tt.wantLen)
			}
		})
	}
}

func TestBitsetSet(t *testing.T) {
	tests := []struct {
		name      string
		size      uint32
		setIndex  uint32
		wantError bool
	}{
		{
			name:      "set first bit",
			size:      16,
			setIndex:  0,
			wantError: false,
		},
		{
			name:      "set middle bit",
			size:      16,
			setIndex:  8,
			wantError: false,
		},
		{
			name:      "set last bit",
			size:      16,
			setIndex:  15,
			wantError: false,
		},
		{
			name:      "set out of range",
			size:      16,
			setIndex:  16,
			wantError: true,
		},
		{
			name:      "set way out of range",
			size:      16,
			setIndex:  100,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := NewBitset(tt.size)
			err := bs.Set(tt.setIndex)
			if (err != nil) != tt.wantError {
				t.Errorf("Set() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				isSet, err := bs.IsSet(tt.setIndex)
				if err != nil {
					t.Fatalf("IsSet() error = %v", err)
				}
				if !isSet {
					t.Errorf("IsSet() = false, want true after Set()")
				}
			}
		})
	}
}

func TestBitsetIsSet(t *testing.T) {
	tests := []struct {
		name      string
		size      uint32
		setBits   []uint32
		checkBit  uint32
		wantSet   bool
		wantError bool
	}{
		{
			name:      "check set bit",
			size:      16,
			setBits:   []uint32{5},
			checkBit:  5,
			wantSet:   true,
			wantError: false,
		},
		{
			name:      "check unset bit",
			size:      16,
			setBits:   []uint32{5},
			checkBit:  6,
			wantSet:   false,
			wantError: false,
		},
		{
			name:      "check multiple set bits",
			size:      16,
			setBits:   []uint32{0, 5, 10, 15},
			checkBit:  10,
			wantSet:   true,
			wantError: false,
		},
		{
			name:      "check out of range",
			size:      16,
			setBits:   []uint32{},
			checkBit:  16,
			wantSet:   false,
			wantError: true,
		},
		{
			name:      "check first bit",
			size:      100,
			setBits:   []uint32{0, 50, 99},
			checkBit:  0,
			wantSet:   true,
			wantError: false,
		},
		{
			name:      "check last bit",
			size:      100,
			setBits:   []uint32{0, 50, 99},
			checkBit:  99,
			wantSet:   true,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := NewBitset(tt.size)
			for _, idx := range tt.setBits {
				if err := bs.Set(idx); err != nil {
					t.Fatalf("Set(%d) error = %v", idx, err)
				}
			}
			isSet, err := bs.IsSet(tt.checkBit)
			if (err != nil) != tt.wantError {
				t.Errorf("IsSet() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && isSet != tt.wantSet {
				t.Errorf("IsSet() = %v, want %v", isSet, tt.wantSet)
			}
		})
	}
}

func TestBitsetUnset(t *testing.T) {
	tests := []struct {
		name         string
		size         uint32
		setBits      []uint32
		unsetBit     uint32
		wantError    bool
		wantSetAfter bool
	}{
		{
			name:         "unset set bit",
			size:         16,
			setBits:      []uint32{5},
			unsetBit:     5,
			wantError:    false,
			wantSetAfter: false,
		},
		{
			name:         "unset already unset bit",
			size:         16,
			setBits:      []uint32{5},
			unsetBit:     6,
			wantError:    false,
			wantSetAfter: false,
		},
		{
			name:         "unset out of range",
			size:         16,
			setBits:      []uint32{},
			unsetBit:     16,
			wantError:    true,
			wantSetAfter: false,
		},
		{
			name:         "unset one of many",
			size:         32,
			setBits:      []uint32{0, 5, 10, 15, 20},
			unsetBit:     10,
			wantError:    false,
			wantSetAfter: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := NewBitset(tt.size)
			for _, idx := range tt.setBits {
				if err := bs.Set(idx); err != nil {
					t.Fatalf("Set(%d) error = %v", idx, err)
				}
			}
			err := bs.Unset(tt.unsetBit)
			if (err != nil) != tt.wantError {
				t.Errorf("Unset() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				isSet, err := bs.IsSet(tt.unsetBit)
				if err != nil {
					t.Fatalf("IsSet() error = %v", err)
				}
				if isSet != tt.wantSetAfter {
					t.Errorf("IsSet() after Unset() = %v, want %v", isSet, tt.wantSetAfter)
				}
			}
		})
	}
}

func TestBitsetToggle(t *testing.T) {
	tests := []struct {
		name         string
		size         uint32
		initialSet   []uint32
		toggleBit    uint32
		wantError    bool
		wantSetAfter bool
	}{
		{
			name:         "toggle unset to set",
			size:         16,
			initialSet:   []uint32{},
			toggleBit:    5,
			wantError:    false,
			wantSetAfter: true,
		},
		{
			name:         "toggle set to unset",
			size:         16,
			initialSet:   []uint32{5},
			toggleBit:    5,
			wantError:    false,
			wantSetAfter: false,
		},
		{
			name:         "toggle out of range",
			size:         16,
			initialSet:   []uint32{},
			toggleBit:    16,
			wantError:    true,
			wantSetAfter: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := NewBitset(tt.size)
			for _, idx := range tt.initialSet {
				if err := bs.Set(idx); err != nil {
					t.Fatalf("Set(%d) error = %v", idx, err)
				}
			}
			err := bs.Toggle(tt.toggleBit)
			if (err != nil) != tt.wantError {
				t.Errorf("Toggle() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				isSet, err := bs.IsSet(tt.toggleBit)
				if err != nil {
					t.Fatalf("IsSet() error = %v", err)
				}
				if isSet != tt.wantSetAfter {
					t.Errorf("IsSet() after Toggle() = %v, want %v", isSet, tt.wantSetAfter)
				}
			}
		})
	}
}

func TestBitsetToggleTwice(t *testing.T) {
	bs := NewBitset(16)
	idx := uint32(5)

	// Initially unset
	isSet, _ := bs.IsSet(idx)
	if isSet {
		t.Fatal("bit should be initially unset")
	}

	// Toggle once - should be set
	if err := bs.Toggle(idx); err != nil {
		t.Fatalf("first Toggle() error = %v", err)
	}
	isSet, _ = bs.IsSet(idx)
	if !isSet {
		t.Error("bit should be set after first toggle")
	}

	// Toggle twice - should be unset again
	if err := bs.Toggle(idx); err != nil {
		t.Fatalf("second Toggle() error = %v", err)
	}
	isSet, _ = bs.IsSet(idx)
	if isSet {
		t.Error("bit should be unset after second toggle")
	}
}

func TestBitsetOperationsSequence(t *testing.T) {
	tests := []struct {
		name       string
		size       uint32
		operations []struct {
			op      string // "set", "unset", "toggle"
			index   uint32
			wantSet bool
		}
	}{
		{
			name: "set then unset",
			size: 16,
			operations: []struct {
				op      string
				index   uint32
				wantSet bool
			}{
				{op: "set", index: 5, wantSet: true},
				{op: "unset", index: 5, wantSet: false},
			},
		},
		{
			name: "toggle sequence",
			size: 16,
			operations: []struct {
				op      string
				index   uint32
				wantSet bool
			}{
				{op: "toggle", index: 5, wantSet: true},
				{op: "toggle", index: 5, wantSet: false},
				{op: "toggle", index: 5, wantSet: true},
			},
		},
		{
			name: "mixed operations",
			size: 32,
			operations: []struct {
				op      string
				index   uint32
				wantSet bool
			}{
				{op: "set", index: 10, wantSet: true},
				{op: "toggle", index: 10, wantSet: false},
				{op: "set", index: 10, wantSet: true},
				{op: "unset", index: 10, wantSet: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := NewBitset(tt.size)
			for i, op := range tt.operations {
				var err error
				switch op.op {
				case "set":
					err = bs.Set(op.index)
				case "unset":
					err = bs.Unset(op.index)
				case "toggle":
					err = bs.Toggle(op.index)
				}
				if err != nil {
					t.Fatalf("operation %d (%s) error = %v", i, op.op, err)
				}
				isSet, err := bs.IsSet(op.index)
				if err != nil {
					t.Fatalf("IsSet() after operation %d error = %v", i, err)
				}
				if isSet != op.wantSet {
					t.Errorf("after operation %d (%s): IsSet() = %v, want %v", i, op.op, isSet, op.wantSet)
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkBitsetSet(b *testing.B) {
	bs := NewBitset(1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bs.Set(uint32(i % 1024))
	}
}

func BenchmarkBitsetIsSet(b *testing.B) {
	bs := NewBitset(1024)
	for i := uint32(0); i < 1024; i++ {
		bs.Set(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bs.IsSet(uint32(i % 1024))
	}
}

func BenchmarkBitsetToggle(b *testing.B) {
	bs := NewBitset(1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bs.Toggle(uint32(i % 1024))
	}
}
