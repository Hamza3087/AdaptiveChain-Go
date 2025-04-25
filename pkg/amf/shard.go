package amf

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/advanced-blockchain/pkg/crypto"
)

// AdvancedShardStatus represents the status of a shard
type AdvancedShardStatus string

const (
	// AdvancedShardStatusActive indicates an active shard
	AdvancedShardStatusActive AdvancedShardStatus = "active"

	// AdvancedShardStatusSplitting indicates a shard that is being split
	AdvancedShardStatusSplitting AdvancedShardStatus = "splitting"

	// AdvancedShardStatusMerging indicates a shard that is being merged
	AdvancedShardStatusMerging AdvancedShardStatus = "merging"

	// AdvancedShardStatusInactive indicates an inactive shard
	AdvancedShardStatusInactive AdvancedShardStatus = "inactive"
)

// AdvancedShard represents a shard in the Adaptive Merkle Forest
type AdvancedShard struct {
	ID           string                 // Shard ID
	Status       AdvancedShardStatus    // Shard status
	Parent       string                 // Parent shard ID
	Children     []string               // Child shard IDs
	StateRoot    crypto.Hash            // Merkle root of the state trie
	StateAccum   *crypto.Accumulator    // Cryptographic accumulator for the state
	StateMap     map[string]interface{} // State map
	LoadFactor   float64                // Load factor for dynamic rebalancing
	LastAccessed time.Time              // Last accessed time
	mu           sync.RWMutex           // Mutex for thread safety
}

// NewAdvancedShard creates a new shard
func NewAdvancedShard(id string, parent string) *AdvancedShard {
	// Create a new RSA modulus for the accumulator
	// In a real implementation, this would be a secure RSA modulus
	n := crypto.NewHash([]byte(id))
	nBig := new(big.Int).SetBytes(n.Bytes())

	return &AdvancedShard{
		ID:           id,
		Status:       AdvancedShardStatusActive,
		Parent:       parent,
		Children:     make([]string, 0),
		StateRoot:    crypto.ZeroHash(),
		StateAccum:   crypto.NewAccumulator(nBig),
		StateMap:     make(map[string]interface{}),
		LoadFactor:   0.0,
		LastAccessed: time.Now(),
	}
}

// UpdateState updates the state of the shard
func (s *AdvancedShard) UpdateState(key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update the state map
	s.StateMap[key] = value

	// Update the state root
	if err := s.updateStateRoot(); err != nil {
		return err
	}

	// Update the load factor
	s.updateLoadFactor()

	// Update the last accessed time
	s.LastAccessed = time.Now()

	return nil
}

// GetState gets a state value from the shard
func (s *AdvancedShard) GetState(key string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Get the value from the state map
	value, ok := s.StateMap[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	// Update the last accessed time
	s.LastAccessed = time.Now()

	return value, nil
}

// GenerateProof generates a proof for a state value
func (s *AdvancedShard) GenerateProof(key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if the key exists
	value, ok := s.StateMap[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	// Serialize the value
	valueData, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal value: %w", err)
	}

	// Create a Merkle proof
	// In a real implementation, we would generate a proper Merkle proof
	// For this implementation, we'll create a more sophisticated proof structure

	// Compute a hash of the key-value pair
	keyValueData := append([]byte(key), valueData...)
	keyValueHash := crypto.NewHash(keyValueData)

	// Create a proof structure
	type MerkleProof struct {
		Key              string        // The key
		Value            []byte        // The serialized value
		KeyValueHash     crypto.Hash   // Hash of the key-value pair
		StateRoot        crypto.Hash   // State root
		AccumulatorValue *big.Int      // Accumulator value
		Siblings         []crypto.Hash // Sibling hashes for the Merkle path
	}

	// Create a proof
	proof := MerkleProof{
		Key:              key,
		Value:            valueData,
		KeyValueHash:     keyValueHash,
		StateRoot:        s.StateRoot,
		AccumulatorValue: s.StateAccum.GetValue(),
		Siblings:         make([]crypto.Hash, 0), // In a real implementation, we would include sibling hashes
	}

	// Serialize the proof
	proofData, err := json.Marshal(proof)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal proof: %w", err)
	}

	// Compress the proof using probabilistic techniques
	// In a real implementation, we would use a more sophisticated compression technique
	// For now, we'll just use a simple compression approach
	compressedProof := compressProof(proofData)

	return compressedProof, nil
}

// compressProof compresses a proof using probabilistic techniques
func compressProof(proofData []byte) []byte {
	// This is a simplified implementation for demonstration purposes
	// In a real implementation, we would use a more sophisticated compression technique

	// For now, we'll just return the original proof data
	// In a real implementation, we would use techniques like Bloom filters or sketches
	return proofData
}

// VerifyProof verifies a proof for a state value
func (s *AdvancedShard) VerifyProof(key string, value string, proof []byte) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Decompress the proof
	// In a real implementation, we would use a more sophisticated decompression technique
	decompressedProof := decompressProof(proof)

	// Deserialize the proof
	var p struct {
		Key              string        // The key
		Value            []byte        // The serialized value
		KeyValueHash     crypto.Hash   // Hash of the key-value pair
		StateRoot        crypto.Hash   // State root
		AccumulatorValue *big.Int      // Accumulator value
		Siblings         []crypto.Hash // Sibling hashes for the Merkle path
	}
	if err := json.Unmarshal(decompressedProof, &p); err != nil {
		return false, fmt.Errorf("failed to unmarshal proof: %w", err)
	}

	// Verify the key matches
	if p.Key != key {
		return false, fmt.Errorf("proof key mismatch: %s != %s", p.Key, key)
	}

	// Verify the value matches
	if string(p.Value) != value {
		return false, fmt.Errorf("proof value mismatch")
	}

	// Verify the key-value hash
	keyValueData := append([]byte(key), p.Value...)
	expectedKeyValueHash := crypto.NewHash(keyValueData)
	if !p.KeyValueHash.Equal(expectedKeyValueHash) {
		return false, fmt.Errorf("proof key-value hash mismatch")
	}

	// Verify the state root
	if !p.StateRoot.Equal(s.StateRoot) {
		return false, fmt.Errorf("proof state root mismatch")
	}

	// Verify the accumulator value
	if p.AccumulatorValue.Cmp(s.StateAccum.GetValue()) != 0 {
		return false, fmt.Errorf("proof accumulator value mismatch")
	}

	// Verify the Merkle path using the sibling hashes
	if len(p.Siblings) > 0 {
		// Compute the Merkle root from the key-value hash and siblings
		computedRoot, err := computeMerkleRootFromPath(p.KeyValueHash, p.Siblings, key)
		if err != nil {
			return false, fmt.Errorf("failed to compute Merkle root: %w", err)
		}

		// Verify the computed root matches the state root
		if !computedRoot.Equal(s.StateRoot) {
			return false, fmt.Errorf("Merkle path verification failed: root mismatch")
		}
	}

	return true, nil
}

// computeMerkleRootFromPath computes the Merkle root from a leaf hash and its siblings
func computeMerkleRootFromPath(leafHash crypto.Hash, siblings []crypto.Hash, key string) (crypto.Hash, error) {
	// Start with the leaf hash
	currentHash := leafHash

	// Compute the position of the leaf in the tree based on the key
	position := computeLeafPosition(key, len(siblings))

	// Combine with siblings to compute the root
	for i, sibling := range siblings {
		// Determine if the current hash should be the left or right child
		// based on the position bit at the current level
		isRightChild := (position >> i) & 1

		var left, right crypto.Hash
		if isRightChild == 1 {
			left = sibling
			right = currentHash
		} else {
			left = currentHash
			right = sibling
		}

		// Combine the hashes
		combined := append(left.Bytes(), right.Bytes()...)
		currentHash = crypto.NewHash(combined)
	}

	return currentHash, nil
}

// computeLeafPosition computes the position of a leaf in the Merkle tree based on its key
func computeLeafPosition(key string, height int) uint64 {
	// Hash the key to get a deterministic position
	h := crypto.NewHash([]byte(key))

	// Use the first 8 bytes of the hash as a uint64
	var position uint64
	for i := 0; i < 8 && i < len(h); i++ {
		position = (position << 8) | uint64(h[i])
	}

	// Mask to the appropriate number of bits based on the height
	mask := uint64((1 << height) - 1)
	return position & mask
}

// decompressProof decompresses a proof
func decompressProof(compressedProof []byte) []byte {
	// This is a simplified implementation for demonstration purposes
	// In a real implementation, we would use a more sophisticated decompression technique

	// For now, we'll just return the original compressed proof
	// In a real implementation, we would use techniques like Bloom filters or sketches
	return compressedProof
}

// GenerateCompressedProof generates a compressed proof for a state value
func (s *AdvancedShard) GenerateCompressedProof(key string) ([]byte, error) {
	// Generate a regular proof
	proof, err := s.GenerateProof(key)
	if err != nil {
		return nil, err
	}

	// Create a probabilistic proof compressor
	compressor := crypto.NewProbabilisticProofCompressor(100, 0.01)

	// Compress the proof
	compressedProof := compressor.CompressProof(proof)

	return compressedProof, nil
}

// VerifyCompressedProof verifies a compressed proof for a state value
func (s *AdvancedShard) VerifyCompressedProof(key string, value string, compressedProof []byte) (bool, error) {
	// Create a probabilistic proof compressor
	compressor := crypto.NewProbabilisticProofCompressor(100, 0.01)

	// Generate the expected proof
	expectedProof, err := s.GenerateProof(key)
	if err != nil {
		return false, err
	}

	// Verify the compressed proof
	return compressor.VerifyCompressedProof(compressedProof, expectedProof), nil
}

// Split splits the shard into two child shards
func (s *AdvancedShard) Split() ([]*AdvancedShard, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the shard is active
	if s.Status != AdvancedShardStatusActive {
		return nil, fmt.Errorf("shard is not active: %s", s.ID)
	}

	// Set the shard status to splitting
	s.Status = AdvancedShardStatusSplitting

	// Create two child shards
	leftID := s.ID + "-0"
	rightID := s.ID + "-1"

	leftShard := NewAdvancedShard(leftID, s.ID)
	rightShard := NewAdvancedShard(rightID, s.ID)

	// Distribute the state between the two child shards
	// For simplicity, we'll just split the state map in half
	i := 0
	for key, value := range s.StateMap {
		if i%2 == 0 {
			leftShard.StateMap[key] = value
		} else {
			rightShard.StateMap[key] = value
		}
		i++
	}

	// Update the state roots of the child shards
	if err := leftShard.updateStateRoot(); err != nil {
		return nil, err
	}
	if err := rightShard.updateStateRoot(); err != nil {
		return nil, err
	}

	// Update the load factors of the child shards
	leftShard.updateLoadFactor()
	rightShard.updateLoadFactor()

	// Update the children of the parent shard
	s.Children = []string{leftID, rightID}

	// Set the parent shard status to inactive
	s.Status = AdvancedShardStatusInactive

	return []*AdvancedShard{leftShard, rightShard}, nil
}

// Merge merges the shard with another shard
func (s *AdvancedShard) Merge(other *AdvancedShard) (*AdvancedShard, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the shards are active
	if s.Status != AdvancedShardStatusActive || other.Status != AdvancedShardStatusActive {
		return nil, errors.New("both shards must be active")
	}

	// Check if the shards have the same parent
	if s.Parent != other.Parent {
		return nil, errors.New("shards must have the same parent")
	}

	// Set the shard statuses to merging
	s.Status = AdvancedShardStatusMerging
	other.Status = AdvancedShardStatusMerging

	// Create a new merged shard
	mergedID := s.Parent
	mergedShard := NewAdvancedShard(mergedID, "")

	// Combine the state maps
	for key, value := range s.StateMap {
		mergedShard.StateMap[key] = value
	}
	for key, value := range other.StateMap {
		mergedShard.StateMap[key] = value
	}

	// Update the state root of the merged shard
	if err := mergedShard.updateStateRoot(); err != nil {
		return nil, err
	}

	// Update the load factor of the merged shard
	mergedShard.updateLoadFactor()

	// Set the original shards to inactive
	s.Status = AdvancedShardStatusInactive
	other.Status = AdvancedShardStatusInactive

	return mergedShard, nil
}

// updateStateRoot updates the state root of the shard
func (s *AdvancedShard) updateStateRoot() error {
	// Compute the state root as the Merkle root of all state values
	var stateHashes []crypto.Hash
	for key, value := range s.StateMap {
		// Serialize the key-value pair
		data, err := json.Marshal(map[string]interface{}{
			"key":   key,
			"value": value,
		})
		if err != nil {
			return fmt.Errorf("failed to marshal key-value pair: %w", err)
		}

		// Compute the hash of the key-value pair
		hash := crypto.NewHash(data)
		stateHashes = append(stateHashes, hash)

		// Add the key-value pair to the accumulator
		s.StateAccum.Add(data)
	}

	// Compute the Merkle root
	s.StateRoot = crypto.ComputeMerkleRoot(stateHashes)

	return nil
}

// updateLoadFactor updates the load factor of the shard
func (s *AdvancedShard) updateLoadFactor() {
	// Compute the load factor as the number of state entries
	s.LoadFactor = float64(len(s.StateMap))
}

// ShouldSplit checks if the shard should be split
func (s *AdvancedShard) ShouldSplit(threshold float64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Status == AdvancedShardStatusActive && s.LoadFactor > threshold
}

// ShouldMerge checks if the shard should be merged
func (s *AdvancedShard) ShouldMerge(threshold float64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Status == AdvancedShardStatusActive && s.LoadFactor < threshold
}

// IsActive checks if the shard is active
func (s *AdvancedShard) IsActive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Status == AdvancedShardStatusActive
}

// GetLoadFactor gets the load factor of the shard
func (s *AdvancedShard) GetLoadFactor() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.LoadFactor
}

// GetStateRoot gets the state root of the shard
func (s *AdvancedShard) GetStateRoot() crypto.Hash {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.StateRoot
}
