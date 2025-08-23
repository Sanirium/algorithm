package bloomfilter

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestNewBloomFilter_InvalidFpPanics(t *testing.T) {
	t.Parallel()

	cases := []float64{-0.1, 0, 1, 1.1}
	for _, fp := range cases {
		func(fp float64) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("expected panic for fpRate=%v", fp)
				}
			}()
			_ = NewBloomFilter(100, fp, 123)
		}(fp)
	}
}

func TestInsertAndContains_NoFalseNegatives(t *testing.T) {
	t.Parallel()

	n := uint32(1000)
	bf := NewBloomFilter(n, 0.01, 42)

	inserted := make([]string, n)
	for i := uint32(0); i < n; i++ {
		k := fmt.Sprintf("in_%d", i)
		inserted[i] = k
		bf.Insert(k)
	}

	for _, k := range inserted {
		if !bf.Contains(k) {
			t.Fatalf("false negative for key %q", k)
		}
	}
}

func TestFalsePositiveRateApprox(t *testing.T) {
	t.Parallel()

	n := uint32(5000)
	wantFPR := 0.01
	bf := NewBloomFilter(n, wantFPR, 777)

	for i := uint32(0); i < n; i++ {
		bf.Insert(fmt.Sprintf("in_%d", i))
	}

	const trials = 20000
	fp := 0
	for i := 0; i < trials; i++ {
		if bf.Contains(fmt.Sprintf("out_%d", i+1_000_000)) {
			fp++
		}
	}
	got := float64(fp) / float64(trials)

	if got > wantFPR*3.0 {
		t.Fatalf("FPR too high: got=%.4f, want<=%.4f", got, wantFPR*3.0)
	}
}

func TestPositionsDeterministicAndInRange(t *testing.T) {
	t.Parallel()

	bf := NewBloomFilter(1000, 0.02, 10101)

	key := "hello-world"
	p1 := bf.key2Positions(key)
	p2 := bf.key2Positions(key)

	if !reflect.DeepEqual(p1, p2) {
		t.Fatalf("positions not deterministic: %v vs %v", p1, p2)
	}

	if uint32(len(p1)) != bf.numHashFunctions {
		t.Fatalf("positions count=%d, want=%d", len(p1), bf.numHashFunctions)
	}

	for i, pos := range p1 {
		if pos >= bf.numBits {
			t.Fatalf("position[%d]=%d out of range [0,%d)", i, pos, bf.numBits)
		}
	}
}

func TestBitOps(t *testing.T) {
	t.Parallel()

	m := uint32(128)
	bits := make([]byte, (m+7)/8)

	for i := uint32(0); i < m; i++ {
		if readBit(bits, i) {
			t.Fatalf("bit %d must be 0 initially", i)
		}
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for k := 0; k < 200; k++ {
		i := uint32(rng.Intn(int(m)))
		writeBit(bits, i)
		if !readBit(bits, i) {
			t.Fatalf("bit %d should be 1 after writeBit", i)
		}
	}
}

func TestHashFunctionsGetterReturnsCopy(t *testing.T) {
	t.Parallel()

	bf := NewBloomFilter(1000, 0.01, 1)
	got := bf.HashFunctions()

	if len(got) != int(bf.numHashFunctions) {
		t.Fatalf("getter len=%d, want=%d", len(got), bf.numHashFunctions)
	}

	before := len(bf.hashFunctions)
	got = append(got, got[0])
	after := len(bf.hashFunctions)
	if before != after {
		t.Fatalf("internal hashFunctions length changed: before=%d after=%d", before, after)
	}
}

func TestKMatchesFormula(t *testing.T) {
	t.Parallel()

	n := uint32(1234)
	fpRate := 0.015
	bf := NewBloomFilter(n, fpRate, 9)

	ln2 := math.Ln2
	m := uint32(math.Ceil(float64(n) * math.Abs(math.Log(fpRate)) / (ln2 * ln2)))
	wantK := uint32(math.Round((float64(m) / float64(n)) * ln2))

	if diff := int(bf.numHashFunctions) - int(wantK); diff < -1 || diff > 1 {
		t.Fatalf("k=%d, want≈%d (±1 allowed)", bf.numHashFunctions, wantK)
	}
}
