package amf

import (
	"fmt"
	"sync"
	"time"
)

// AdaptiveMerkleForest represents a forest of Merkle trees with adaptive sharding
type AdaptiveMerkleForest struct {
	Shards            map[string]*AdvancedShard // Map of shard ID to shard
	SplitThreshold    float64                   // Threshold for splitting a shard
	MergeThreshold    float64                   // Threshold for merging shards
	RebalanceInterval time.Duration             // Interval for rebalancing the forest
	lastRebalance     time.Time                 // Last rebalance time
	mu                sync.RWMutex              // Mutex for thread safety
}

// NewAdaptiveMerkleForest creates a new adaptive Merkle forest
func NewAdaptiveMerkleForest(splitThreshold float64, rebalanceInterval time.Duration) *AdaptiveMerkleForest {
	return &AdaptiveMerkleForest{
		Shards:            make(map[string]*AdvancedShard),
		SplitThreshold:    splitThreshold,
		MergeThreshold:    splitThreshold / 4, // Set merge threshold to 1/4 of split threshold
		RebalanceInterval: rebalanceInterval,
		lastRebalance:     time.Now(),
	}
}

// GetShard gets a shard by ID, creating it if it doesn't exist
func (f *AdaptiveMerkleForest) GetShard(id string) (*AdvancedShard, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	// Check if the shard exists
	shard, ok := f.Shards[id]
	if ok {
		return shard, nil
	}
	
	// Create a new shard
	shard = NewAdvancedShard(id, "")
	f.Shards[id] = shard
	
	return shard, nil
}

// UpdateState updates the state of a shard
func (f *AdaptiveMerkleForest) UpdateState(shardID string, key string, value interface{}) error {
	// Get the shard
	shard, err := f.GetShard(shardID)
	if err != nil {
		return err
	}
	
	// Update the state
	if err := shard.UpdateState(key, value); err != nil {
		return err
	}
	
	// Check if rebalancing is needed
	f.checkRebalance()
	
	return nil
}

// GetState gets a state value from a shard
func (f *AdaptiveMerkleForest) GetState(shardID string, key string) (interface{}, error) {
	// Get the shard
	shard, err := f.GetShard(shardID)
	if err != nil {
		return nil, err
	}
	
	// Get the state
	return shard.GetState(key)
}

// GenerateProof generates a proof for a state value
func (f *AdaptiveMerkleForest) GenerateProof(shardID string, key string) ([]byte, error) {
	// Get the shard
	shard, err := f.GetShard(shardID)
	if err != nil {
		return nil, err
	}
	
	// Generate the proof
	return shard.GenerateProof(key)
}

// VerifyProof verifies a proof for a state value
func (f *AdaptiveMerkleForest) VerifyProof(shardID string, key string, value string, proof []byte) (bool, error) {
	// Get the shard
	shard, err := f.GetShard(shardID)
	if err != nil {
		return false, err
	}
	
	// Verify the proof
	return shard.VerifyProof(key, value, proof)
}

// Rebalance rebalances the forest
func (f *AdaptiveMerkleForest) Rebalance() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	// Update the last rebalance time
	f.lastRebalance = time.Now()
	
	// Check each shard for splitting
	for id, shard := range f.Shards {
		if shard.ShouldSplit(f.SplitThreshold) {
			// Split the shard
			childShards, err := shard.Split()
			if err != nil {
				return fmt.Errorf("failed to split shard %s: %w", id, err)
			}
			
			// Add the child shards to the forest
			for _, childShard := range childShards {
				f.Shards[childShard.ID] = childShard
			}
		}
	}
	
	// Check for merging sibling shards
	// This is a simplified implementation for demonstration purposes
	// In a real implementation, we would need to identify sibling shards properly
	for id1, shard1 := range f.Shards {
		if !shard1.IsActive() || !shard1.ShouldMerge(f.MergeThreshold) {
			continue
		}
		
		// Find a sibling shard
		for id2, shard2 := range f.Shards {
			if id1 == id2 || !shard2.IsActive() || !shard2.ShouldMerge(f.MergeThreshold) {
				continue
			}
			
			// Check if they are siblings (have the same parent)
			if shard1.Parent == shard2.Parent && shard1.Parent != "" {
				// Merge the shards
				mergedShard, err := shard1.Merge(shard2)
				if err != nil {
					return fmt.Errorf("failed to merge shards %s and %s: %w", id1, id2, err)
				}
				
				// Add the merged shard to the forest
				f.Shards[mergedShard.ID] = mergedShard
				
				// We've merged a pair of shards, so break out of the inner loop
				break
			}
		}
	}
	
	return nil
}

// checkRebalance checks if rebalancing is needed
func (f *AdaptiveMerkleForest) checkRebalance() {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	// Check if it's time to rebalance
	if time.Since(f.lastRebalance) >= f.RebalanceInterval {
		// Rebalance in a separate goroutine to avoid blocking
		go func() {
			if err := f.Rebalance(); err != nil {
				// Log the error
				fmt.Printf("Error rebalancing forest: %v\n", err)
			}
		}()
	}
}

// GetShardCount gets the number of active shards
func (f *AdaptiveMerkleForest) GetShardCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	count := 0
	for _, shard := range f.Shards {
		if shard.IsActive() {
			count++
		}
	}
	
	return count
}

// FindShardForKey finds the appropriate shard for a key
// This is a simplified implementation for demonstration purposes
func (f *AdaptiveMerkleForest) FindShardForKey(key string) (string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	// In a real implementation, we would use a more sophisticated algorithm
	// For now, we'll just use a simple hash-based approach
	
	// Compute a hash of the key
	hash := 0
	for _, c := range key {
		hash = (hash*31 + int(c)) % 1000000
	}
	
	// Find an active shard with the closest ID
	var bestShard *AdvancedShard
	var bestDiff int = 1000000
	
	for _, shard := range f.Shards {
		if !shard.IsActive() {
			continue
		}
		
		// Compute a hash of the shard ID
		shardHash := 0
		for _, c := range shard.ID {
			shardHash = (shardHash*31 + int(c)) % 1000000
		}
		
		// Compute the difference
		diff := hash - shardHash
		if diff < 0 {
			diff = -diff
		}
		
		// Update the best shard if this one is closer
		if diff < bestDiff {
			bestShard = shard
			bestDiff = diff
		}
	}
	
	if bestShard == nil {
		return "", fmt.Errorf("no active shards found")
	}
	
	return bestShard.ID, nil
}

// GetAllShards gets all active shards
func (f *AdaptiveMerkleForest) GetAllShards() []*AdvancedShard {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	var shards []*AdvancedShard
	for _, shard := range f.Shards {
		if shard.IsActive() {
			shards = append(shards, shard)
		}
	}
	
	return shards
}
