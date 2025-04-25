package amf

import (
	"fmt"
	"testing"
	"time"

	"github.com/advanced-blockchain/pkg/crypto"
	"github.com/advanced-blockchain/pkg/state"
)

func TestAdvancedShard(t *testing.T) {
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

	// Generate a key pair for signing transactions
	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Get the public key
	publicKey, err := keyPair.ExportPublicKey()
	if err != nil {
		t.Fatalf("Failed to export public key: %v", err)
	}

	// Create a transaction
	tx := state.NewTransaction(
		publicKey,
		publicKey, // Using same key for simplicity
		100,
		[]byte("test data"),
		1, // ShardID as uint64
	)

	// Sign the transaction
	err = tx.Sign(keyPair)
	if err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}

	// Create a block with the transaction
	previousHash := crypto.ZeroHash()
	_ = state.NewBlock(previousHash, []*state.Transaction{tx}, 1) // Ignore the block if unused

	// Update the shard state
	updateErr := shard.UpdateState(tx.ID.String(), tx)
	if updateErr != nil {
		t.Fatalf("Failed to update shard state: %v", updateErr)
	}

	// Verify the shard has the transaction
	value, err := shard.GetState(tx.ID.String())
	if err != nil {
		t.Errorf("Failed to get transaction from shard: %v", err)
	}
	if value == nil {
		t.Errorf("Retrieved transaction is nil")
	}

	// Generate a proof for the transaction
	proof, err := shard.GenerateProof(tx.ID.String())
	if err != nil {
		t.Fatalf("Failed to generate proof: %v", err)
	}

	// Verify the proof
	// For the test, we'll use a simplified value representation
	txValue := fmt.Sprintf("%x:%x", tx.From, tx.To)
	valid, err := shard.VerifyProof(tx.ID.String(), txValue, proof)
	if err != nil {
		t.Fatalf("Failed to verify proof: %v", err)
	}
	if !valid {
		t.Errorf("Proof verification failed")
	}
}
