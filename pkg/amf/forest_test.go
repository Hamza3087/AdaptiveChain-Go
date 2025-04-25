package amf

import (
	"fmt"
	"testing"
	"time"
)

func TestAdaptiveMerkleForest(t *testing.T) {
	// Create a new forest
	forest := NewAdaptiveMerkleForest(10, 5*time.Minute)

	// Create a new shard
	shardID := "shard-1"
	shard, err := forest.GetShard(shardID)
	if err != nil {
		t.Fatalf("Failed to create shard: %v", err)
	}

	// Verify the shard fields
	if shard.ID != shardID {
		t.Errorf("Shard has incorrect ID: got %s, want %s", shard.ID, shardID)
	}
	if shard.Status != AdvancedShardStatusActive {
		t.Errorf("Shard has incorrect status: got %s, want %s", shard.Status, AdvancedShardStatusActive)
	}

	// Update the shard state
	key := "test-key"
	value := "test-value"
	updateErr := shard.UpdateState(key, value)
	if updateErr != nil {
		t.Fatalf("Failed to update shard state: %v", updateErr)
	}

	// Verify the shard has the value
	retrievedValue, err := shard.GetState(key)
	if err != nil {
		t.Errorf("Failed to get value from shard: %v", err)
	}
	if retrievedValue != value {
		t.Errorf("Retrieved value is incorrect: got %v, want %v", retrievedValue, value)
	}

	// Generate a proof for the value
	proof, err := shard.GenerateProof(key)
	if err != nil {
		t.Fatalf("Failed to generate proof: %v", err)
	}

	// Verify the proof
	valid, err := shard.VerifyProof(key, value, proof)
	if err != nil {
		t.Fatalf("Failed to verify proof: %v", err)
	}
	if !valid {
		t.Errorf("Proof verification failed")
	}

	// Test shard splitting
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		if err := shard.UpdateState(key, value); err != nil {
			t.Fatalf("Failed to update shard state: %v", err)
		}
	}

	// Check if the shard should split
	if !shard.ShouldSplit(10) {
		t.Errorf("Shard should split but doesn't")
	}

	// Split the shard
	childShards, err := shard.Split()
	if err != nil {
		t.Fatalf("Failed to split shard: %v", err)
	}

	// Verify the child shards
	if len(childShards) != 2 {
		t.Errorf("Split should create 2 child shards, got %d", len(childShards))
	}

	// Verify the parent shard status
	if shard.Status != AdvancedShardStatusInactive {
		t.Errorf("Parent shard has incorrect status after split: got %s, want %s", shard.Status, AdvancedShardStatusInactive)
	}

	// Verify the child shard IDs
	expectedLeftID := shardID + "-0"
	expectedRightID := shardID + "-1"

	if childShards[0].ID != expectedLeftID && childShards[0].ID != expectedRightID {
		t.Errorf("Child shard has incorrect ID: got %s, want %s or %s", childShards[0].ID, expectedLeftID, expectedRightID)
	}

	if childShards[1].ID != expectedLeftID && childShards[1].ID != expectedRightID {
		t.Errorf("Child shard has incorrect ID: got %s, want %s or %s", childShards[1].ID, expectedLeftID, expectedRightID)
	}

	// Verify the child shard parent IDs
	if childShards[0].Parent != shardID {
		t.Errorf("Child shard has incorrect parent ID: got %s, want %s", childShards[0].Parent, shardID)
	}

	if childShards[1].Parent != shardID {
		t.Errorf("Child shard has incorrect parent ID: got %s, want %s", childShards[1].Parent, shardID)
	}

	// Verify the child shard statuses
	if childShards[0].Status != AdvancedShardStatusActive {
		t.Errorf("Child shard has incorrect status: got %s, want %s", childShards[0].Status, AdvancedShardStatusActive)
	}

	if childShards[1].Status != AdvancedShardStatusActive {
		t.Errorf("Child shard has incorrect status: got %s, want %s", childShards[1].Status, AdvancedShardStatusActive)
	}

	// Test forest rebalancing
	if err := forest.Rebalance(); err != nil {
		t.Fatalf("Failed to rebalance forest: %v", err)
	}

	// Verify the forest has the child shards
	if len(forest.Shards) != 3 { // parent + 2 children
		t.Errorf("Forest has incorrect number of shards: got %d, want 3", len(forest.Shards))
	}

	// Verify the forest has the correct number of active shards
	if forest.GetShardCount() != 2 { // only the 2 children should be active
		t.Errorf("Forest has incorrect number of active shards: got %d, want 2", forest.GetShardCount())
	}
}
