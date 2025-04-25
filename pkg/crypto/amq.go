package crypto

import (
	"encoding/binary"
	"hash/fnv"
	"math"
	"math/bits"
	"sync"
)

// BloomFilter represents a Bloom filter for approximate membership queries
type BloomFilter struct {
	bits      []uint64     // Bit array
	numBits   uint         // Number of bits in the filter
	numHashes uint         // Number of hash functions
	numItems  uint         // Number of items in the filter
	mu        sync.RWMutex // Mutex for thread safety
}

// NewBloomFilter creates a new Bloom filter
func NewBloomFilter(expectedItems uint, falsePositiveRate float64) *BloomFilter {
	// Calculate the optimal number of bits and hash functions
	numBits := uint(math.Ceil(-float64(expectedItems) * math.Log(falsePositiveRate) / math.Pow(math.Log(2), 2)))
	numHashes := uint(math.Ceil(float64(numBits) / float64(expectedItems) * math.Log(2)))

	// Ensure minimum values
	if numBits < 1 {
		numBits = 1
	}
	if numHashes < 1 {
		numHashes = 1
	}

	// Calculate the number of uint64 words needed
	numWords := (numBits + 63) / 64

	return &BloomFilter{
		bits:      make([]uint64, numWords),
		numBits:   numBits,
		numHashes: numHashes,
		numItems:  0,
	}
}

// Add adds an item to the Bloom filter
func (bf *BloomFilter) Add(item []byte) {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	// Compute the hash values
	h1, h2 := bf.hashValues(item)

	// Set the bits
	for i := uint(0); i < bf.numHashes; i++ {
		// Use the double hashing technique to generate multiple hash values
		hashValue := (h1 + i*h2) % bf.numBits
		wordIndex := hashValue / 64
		bitIndex := hashValue % 64
		bf.bits[wordIndex] |= 1 << bitIndex
	}

	bf.numItems++
}

// Contains checks if an item might be in the Bloom filter
func (bf *BloomFilter) Contains(item []byte) bool {
	bf.mu.RLock()
	defer bf.mu.RUnlock()

	// Compute the hash values
	h1, h2 := bf.hashValues(item)

	// Check the bits
	for i := uint(0); i < bf.numHashes; i++ {
		// Use the double hashing technique to generate multiple hash values
		hashValue := (h1 + i*h2) % bf.numBits
		wordIndex := hashValue / 64
		bitIndex := hashValue % 64
		if (bf.bits[wordIndex] & (1 << bitIndex)) == 0 {
			return false
		}
	}

	return true
}

// hashValues computes two hash values for an item
func (bf *BloomFilter) hashValues(item []byte) (uint, uint) {
	// Use FNV hash for simplicity
	// In a real implementation, we would use more sophisticated hash functions
	h := fnv.New64a()
	h.Write(item)
	h1 := uint(h.Sum64() % uint64(bf.numBits))

	// Use a different seed for the second hash
	h.Reset()
	h.Write([]byte("second-hash-seed"))
	h.Write(item)
	h2 := uint(h.Sum64() % uint64(bf.numBits))

	return h1, h2
}

// Clear clears the Bloom filter
func (bf *BloomFilter) Clear() {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	for i := range bf.bits {
		bf.bits[i] = 0
	}
	bf.numItems = 0
}

// EstimatedFalsePositiveRate returns the estimated false positive rate
func (bf *BloomFilter) EstimatedFalsePositiveRate() float64 {
	bf.mu.RLock()
	defer bf.mu.RUnlock()

	// Calculate the estimated false positive rate
	// p = (1 - e^(-kn/m))^k
	// where k is the number of hash functions, n is the number of items, and m is the number of bits
	exponent := -float64(bf.numHashes) * float64(bf.numItems) / float64(bf.numBits)
	return math.Pow(1.0-math.Exp(exponent), float64(bf.numHashes))
}

// CountingBloomFilter represents a counting Bloom filter
type CountingBloomFilter struct {
	counters    []uint8      // Counter array
	numCounters uint         // Number of counters in the filter
	numHashes   uint         // Number of hash functions
	numItems    uint         // Number of items in the filter
	mu          sync.RWMutex // Mutex for thread safety
}

// NewCountingBloomFilter creates a new counting Bloom filter
func NewCountingBloomFilter(expectedItems uint, falsePositiveRate float64) *CountingBloomFilter {
	// Calculate the optimal number of counters and hash functions
	numCounters := uint(math.Ceil(-float64(expectedItems) * math.Log(falsePositiveRate) / math.Pow(math.Log(2), 2)))
	numHashes := uint(math.Ceil(float64(numCounters) / float64(expectedItems) * math.Log(2)))

	// Ensure minimum values
	if numCounters < 1 {
		numCounters = 1
	}
	if numHashes < 1 {
		numHashes = 1
	}

	return &CountingBloomFilter{
		counters:    make([]uint8, numCounters),
		numCounters: numCounters,
		numHashes:   numHashes,
		numItems:    0,
	}
}

// Add adds an item to the counting Bloom filter
func (cbf *CountingBloomFilter) Add(item []byte) {
	cbf.mu.Lock()
	defer cbf.mu.Unlock()

	// Compute the hash values
	h1, h2 := cbf.hashValues(item)

	// Increment the counters
	for i := uint(0); i < cbf.numHashes; i++ {
		// Use the double hashing technique to generate multiple hash values
		hashValue := (h1 + i*h2) % cbf.numCounters
		// Increment the counter, but don't overflow
		if cbf.counters[hashValue] < 255 {
			cbf.counters[hashValue]++
		}
	}

	cbf.numItems++
}

// Remove removes an item from the counting Bloom filter
func (cbf *CountingBloomFilter) Remove(item []byte) {
	cbf.mu.Lock()
	defer cbf.mu.Unlock()

	// Compute the hash values
	h1, h2 := cbf.hashValues(item)

	// Decrement the counters
	for i := uint(0); i < cbf.numHashes; i++ {
		// Use the double hashing technique to generate multiple hash values
		hashValue := (h1 + i*h2) % cbf.numCounters
		// Decrement the counter, but don't underflow
		if cbf.counters[hashValue] > 0 {
			cbf.counters[hashValue]--
		}
	}

	if cbf.numItems > 0 {
		cbf.numItems--
	}
}

// Contains checks if an item might be in the counting Bloom filter
func (cbf *CountingBloomFilter) Contains(item []byte) bool {
	cbf.mu.RLock()
	defer cbf.mu.RUnlock()

	// Compute the hash values
	h1, h2 := cbf.hashValues(item)

	// Check the counters
	for i := uint(0); i < cbf.numHashes; i++ {
		// Use the double hashing technique to generate multiple hash values
		hashValue := (h1 + i*h2) % cbf.numCounters
		if cbf.counters[hashValue] == 0 {
			return false
		}
	}

	return true
}

// hashValues computes two hash values for an item
func (cbf *CountingBloomFilter) hashValues(item []byte) (uint, uint) {
	// Use FNV hash for simplicity
	// In a real implementation, we would use more sophisticated hash functions
	h := fnv.New64a()
	h.Write(item)
	h1 := uint(h.Sum64() % uint64(cbf.numCounters))

	// Use a different seed for the second hash
	h.Reset()
	h.Write([]byte("second-hash-seed"))
	h.Write(item)
	h2 := uint(h.Sum64() % uint64(cbf.numCounters))

	return h1, h2
}

// Clear clears the counting Bloom filter
func (cbf *CountingBloomFilter) Clear() {
	cbf.mu.Lock()
	defer cbf.mu.Unlock()

	for i := range cbf.counters {
		cbf.counters[i] = 0
	}
	cbf.numItems = 0
}

// EstimatedFalsePositiveRate returns the estimated false positive rate
func (cbf *CountingBloomFilter) EstimatedFalsePositiveRate() float64 {
	cbf.mu.RLock()
	defer cbf.mu.RUnlock()

	// Calculate the estimated false positive rate
	// p = (1 - e^(-kn/m))^k
	// where k is the number of hash functions, n is the number of items, and m is the number of counters
	exponent := -float64(cbf.numHashes) * float64(cbf.numItems) / float64(cbf.numCounters)
	return math.Pow(1.0-math.Exp(exponent), float64(cbf.numHashes))
}

// CuckooFilter represents a Cuckoo filter for approximate membership queries
type CuckooFilter struct {
	buckets         []uint64     // Buckets array
	numBuckets      uint         // Number of buckets in the filter
	bucketSize      uint         // Number of entries per bucket
	fingerprintSize uint         // Size of fingerprint in bits
	maxKicks        uint         // Maximum number of kicks in insertion
	numItems        uint         // Number of items in the filter
	mu              sync.RWMutex // Mutex for thread safety
}

// NewCuckooFilter creates a new Cuckoo filter
func NewCuckooFilter(expectedItems uint, falsePositiveRate float64) *CuckooFilter {
	// Calculate the optimal number of buckets and fingerprint size
	fingerprintSize := uint(math.Ceil(-math.Log2(falsePositiveRate)))
	if fingerprintSize < 4 {
		fingerprintSize = 4 // Minimum fingerprint size
	}
	if fingerprintSize > 32 {
		fingerprintSize = 32 // Maximum fingerprint size
	}

	bucketSize := uint(4)                                        // Fixed bucket size
	numBuckets := nextPowerOfTwo(expectedItems / bucketSize * 2) // Ensure power of 2 for fast modulo

	// Calculate the number of uint64 words needed
	entriesPerWord := 64 / fingerprintSize
	wordsPerBucket := (bucketSize + entriesPerWord - 1) / entriesPerWord
	numWords := numBuckets * wordsPerBucket

	return &CuckooFilter{
		buckets:         make([]uint64, numWords),
		numBuckets:      numBuckets,
		bucketSize:      bucketSize,
		fingerprintSize: fingerprintSize,
		maxKicks:        500, // Default max kicks
		numItems:        0,
	}
}

// Add adds an item to the Cuckoo filter
func (cf *CuckooFilter) Add(item []byte) bool {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	// Compute the fingerprint and hash values
	fingerprint := cf.computeFingerprint(item)
	i1 := cf.hashValue(item)
	i2 := cf.altIndex(i1, fingerprint)

	// Try to insert into bucket i1
	if cf.insertIntoBucket(i1, fingerprint) {
		cf.numItems++
		return true
	}

	// Try to insert into bucket i2
	if cf.insertIntoBucket(i2, fingerprint) {
		cf.numItems++
		return true
	}

	// Both buckets are full, start cuckoo kicking
	i := i1
	for kicks := uint(0); kicks < cf.maxKicks; kicks++ {
		// Pick a random entry from the bucket
		entryIndex := uint(RandomFloat() * float64(cf.bucketSize))

		// Swap fingerprints
		oldFingerprint := cf.getFingerprint(i, entryIndex)
		cf.setFingerprint(i, entryIndex, fingerprint)
		fingerprint = oldFingerprint

		// Calculate the alternate bucket
		i = cf.altIndex(i, fingerprint)

		// Try to insert into the alternate bucket
		if cf.insertIntoBucket(i, fingerprint) {
			cf.numItems++
			return true
		}
	}

	// Insertion failed after max kicks
	return false
}

// Contains checks if an item might be in the Cuckoo filter
func (cf *CuckooFilter) Contains(item []byte) bool {
	cf.mu.RLock()
	defer cf.mu.RUnlock()

	// Compute the fingerprint and hash values
	fingerprint := cf.computeFingerprint(item)
	i1 := cf.hashValue(item)
	i2 := cf.altIndex(i1, fingerprint)

	// Check both buckets
	return cf.containsInBucket(i1, fingerprint) || cf.containsInBucket(i2, fingerprint)
}

// Remove removes an item from the Cuckoo filter
func (cf *CuckooFilter) Remove(item []byte) bool {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	// Compute the fingerprint and hash values
	fingerprint := cf.computeFingerprint(item)
	i1 := cf.hashValue(item)
	i2 := cf.altIndex(i1, fingerprint)

	// Try to remove from bucket i1
	if cf.removeFromBucket(i1, fingerprint) {
		if cf.numItems > 0 {
			cf.numItems--
		}
		return true
	}

	// Try to remove from bucket i2
	if cf.removeFromBucket(i2, fingerprint) {
		if cf.numItems > 0 {
			cf.numItems--
		}
		return true
	}

	// Item not found
	return false
}

// computeFingerprint computes a fingerprint for an item
func (cf *CuckooFilter) computeFingerprint(item []byte) uint64 {
	// Use FNV hash for simplicity
	// In a real implementation, we would use more sophisticated hash functions
	h := fnv.New64a()
	h.Write(item)
	hash := h.Sum64()

	// Ensure the fingerprint is not zero (zero is reserved for empty entries)
	fingerprint := hash & ((1 << cf.fingerprintSize) - 1)
	if fingerprint == 0 {
		fingerprint = 1
	}

	return fingerprint
}

// hashValue computes the primary hash value for an item
func (cf *CuckooFilter) hashValue(item []byte) uint {
	// Use FNV hash for simplicity
	h := fnv.New64a()
	h.Write(item)
	hash := h.Sum64()

	// Modulo by number of buckets (which is a power of 2)
	return uint(hash & (uint64(cf.numBuckets) - 1))
}

// altIndex computes the alternate index for a given index and fingerprint
func (cf *CuckooFilter) altIndex(index uint, fingerprint uint64) uint {
	// Use the fingerprint to compute the alternate index
	// i2 = i1 ^ hash(fingerprint)
	h := fnv.New64a()
	binary.Write(h, binary.LittleEndian, fingerprint)
	hash := h.Sum64()

	// Modulo by number of buckets (which is a power of 2)
	return uint((uint64(index) ^ hash) & (uint64(cf.numBuckets) - 1))
}

// insertIntoBucket inserts a fingerprint into a bucket
func (cf *CuckooFilter) insertIntoBucket(bucketIndex uint, fingerprint uint64) bool {
	// Check if there's an empty slot in the bucket
	for i := uint(0); i < cf.bucketSize; i++ {
		if cf.getFingerprint(bucketIndex, i) == 0 {
			cf.setFingerprint(bucketIndex, i, fingerprint)
			return true
		}
	}
	return false
}

// containsInBucket checks if a fingerprint is in a bucket
func (cf *CuckooFilter) containsInBucket(bucketIndex uint, fingerprint uint64) bool {
	for i := uint(0); i < cf.bucketSize; i++ {
		if cf.getFingerprint(bucketIndex, i) == fingerprint {
			return true
		}
	}
	return false
}

// removeFromBucket removes a fingerprint from a bucket
func (cf *CuckooFilter) removeFromBucket(bucketIndex uint, fingerprint uint64) bool {
	for i := uint(0); i < cf.bucketSize; i++ {
		if cf.getFingerprint(bucketIndex, i) == fingerprint {
			cf.setFingerprint(bucketIndex, i, 0)
			return true
		}
	}
	return false
}

// getFingerprint gets a fingerprint from a bucket
func (cf *CuckooFilter) getFingerprint(bucketIndex uint, entryIndex uint) uint64 {
	// Calculate the word and bit position
	entriesPerWord := 64 / cf.fingerprintSize
	wordsPerBucket := (cf.bucketSize + entriesPerWord - 1) / entriesPerWord
	wordIndex := bucketIndex*wordsPerBucket + entryIndex/entriesPerWord
	bitPosition := (entryIndex % entriesPerWord) * cf.fingerprintSize

	// Extract the fingerprint
	mask := uint64((1 << cf.fingerprintSize) - 1)
	return (cf.buckets[wordIndex] >> bitPosition) & mask
}

// setFingerprint sets a fingerprint in a bucket
func (cf *CuckooFilter) setFingerprint(bucketIndex uint, entryIndex uint, fingerprint uint64) {
	// Calculate the word and bit position
	entriesPerWord := 64 / cf.fingerprintSize
	wordsPerBucket := (cf.bucketSize + entriesPerWord - 1) / entriesPerWord
	wordIndex := bucketIndex*wordsPerBucket + entryIndex/entriesPerWord
	bitPosition := (entryIndex % entriesPerWord) * cf.fingerprintSize

	// Clear the old fingerprint
	mask := uint64((1 << cf.fingerprintSize) - 1)
	cf.buckets[wordIndex] &= ^(mask << bitPosition)

	// Set the new fingerprint
	cf.buckets[wordIndex] |= (fingerprint & mask) << bitPosition
}

// Clear clears the Cuckoo filter
func (cf *CuckooFilter) Clear() {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	for i := range cf.buckets {
		cf.buckets[i] = 0
	}
	cf.numItems = 0
}

// EstimatedFalsePositiveRate returns the estimated false positive rate
func (cf *CuckooFilter) EstimatedFalsePositiveRate() float64 {
	cf.mu.RLock()
	defer cf.mu.RUnlock()

	// Calculate the estimated false positive rate
	// p = 1 - (1 - 1/2^f)^(2*n)
	// where f is the fingerprint size and n is the number of items
	return 1.0 - math.Pow(1.0-1.0/math.Pow(2, float64(cf.fingerprintSize)), 2.0*float64(cf.numItems))
}

// nextPowerOfTwo returns the next power of two greater than or equal to n
func nextPowerOfTwo(n uint) uint {
	if n == 0 {
		return 1
	}
	return 1 << (64 - bits.LeadingZeros64(uint64(n-1)))
}

// ProbabilisticProofCompressor represents a probabilistic proof compressor
type ProbabilisticProofCompressor struct {
	filter *BloomFilter // Bloom filter for compression
}

// NewProbabilisticProofCompressor creates a new probabilistic proof compressor
func NewProbabilisticProofCompressor(expectedItems uint, falsePositiveRate float64) *ProbabilisticProofCompressor {
	return &ProbabilisticProofCompressor{
		filter: NewBloomFilter(expectedItems, falsePositiveRate),
	}
}

// CompressProof compresses a proof
func (ppc *ProbabilisticProofCompressor) CompressProof(proof []byte) []byte {
	// This is a simplified implementation for demonstration purposes
	// In a real implementation, we would use a more sophisticated compression technique

	// Add the proof to the Bloom filter
	ppc.filter.Add(proof)

	// Compute a hash of the proof to use as a fingerprint
	h := fnv.New64a()
	h.Write(proof)
	fingerprint := h.Sum64()

	// Create a compressed representation of the proof
	// For demonstration purposes, we'll just use the fingerprint
	// In a real implementation, we would use a more sophisticated compression technique
	compressedProof := make([]byte, 8)
	binary.BigEndian.PutUint64(compressedProof, fingerprint)

	return compressedProof
}

// VerifyCompressedProof verifies a compressed proof
func (ppc *ProbabilisticProofCompressor) VerifyCompressedProof(compressedProof []byte, expectedProof []byte) bool {
	// This is a simplified implementation for demonstration purposes
	// In a real implementation, we would use a more sophisticated verification technique

	// Check if the compressed proof is valid
	if len(compressedProof) != 8 {
		return false
	}

	// Extract the fingerprint from the compressed proof
	fingerprint := binary.BigEndian.Uint64(compressedProof)

	// Compute the fingerprint of the expected proof
	h := fnv.New64a()
	h.Write(expectedProof)
	expectedFingerprint := h.Sum64()

	// Check if the fingerprints match
	if fingerprint != expectedFingerprint {
		return false
	}

	// Check if the expected proof is in the Bloom filter
	return ppc.filter.Contains(expectedProof)
}
