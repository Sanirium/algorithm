package bloomfilter

import (
	"github.com/twmb/murmur3"
	"hash/fnv"
	"math"
)

type Server interface {
	LoadData() []string
	Contains(value string) bool
	StoreData(value string) bool
}

type hashFunction func(h1, h2 uint32) uint32

type BloomFilter struct {
	bitsArray        []byte
	hashFunctions    []hashFunction
	numBits          uint32
	maxSize          uint32
	numHashFunctions uint32
	seed             uint32
}

func NewBloomFilter(n uint32, fpRate float64, seed uint32) *BloomFilter {
	if fpRate <= 0 || fpRate >= 1 {
		panic("fpRate must be in (0,1)")
	}
	ln2 := math.Ln2
	numBits := uint32(math.Ceil(float64(n) * math.Abs(math.Log(fpRate)) / (ln2 * ln2)))
	k := uint32(math.Max(1, math.Round((float64(numBits)/float64(n))*ln2)))

	byteLen := (numBits + 7) / 8
	bf := &BloomFilter{
		maxSize:          n,
		seed:             seed,
		numBits:          numBits,
		bitsArray:        make([]byte, byteLen),
		numHashFunctions: k,
	}
	bf.hashFunctions = initHashFunctions(k, numBits)
	return bf
}

func (b *BloomFilter) HashFunctions() []hashFunction {
	cp := make([]hashFunction, len(b.hashFunctions))
	copy(cp, b.hashFunctions)
	return cp
}

func (b *BloomFilter) Insert(value string) {
	for _, p := range b.key2Positions(value) {
		writeBit(b.bitsArray, p)
	}
}

func (b *BloomFilter) Contains(value string) bool {
	for _, p := range b.key2Positions(value) {
		if !readBit(b.bitsArray, p) {
			return false
		}
	}
	return true
}

func (b *BloomFilter) key2Positions(key string) []uint32 {
	h1 := murmur3.SeedSum32(b.seed, []byte(key))

	f := fnv.New32a()
	_, _ = f.Write([]byte(key))
	h2 := f.Sum32()

	pos := make([]uint32, b.numHashFunctions)
	for i, hf := range b.hashFunctions {
		pos[i] = hf(h1, h2)
	}
	return pos
}

func initHashFunctions(k, numBits uint32) []hashFunction {
	hs := make([]hashFunction, k)
	for i := uint32(0); i < k; i++ {
		j := i
		hs[i] = func(h1, h2 uint32) uint32 {
			return (h1 + j*h2 + j*j) % numBits
		}
	}
	return hs
}

func findBitCoordinates(index uint32) (uint32, uint32) {
	byteIndex := index / 8
	bitOffset := index % 8
	return byteIndex, bitOffset
}

func readBit(bitsArray []byte, index uint32) bool {
	element, bit := findBitCoordinates(index)
	return (bitsArray[element] & (1 << bit)) != 0
}

func writeBit(bitsArray []byte, index uint32) {
	element, bit := findBitCoordinates(index)
	bitsArray[element] |= (1 << bit)
}
