package state

import (
	"testing"

	"github.com/advanced-blockchain/pkg/crypto"
)

func TestState(t *testing.T) {
	// Create a new state
	state := NewState()
	
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
	tx := NewTransaction(
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
	
	// Add the transaction to the state
	err = state.AddTransaction(tx)
	if err != nil {
		t.Fatalf("Failed to add transaction to state: %v", err)
	}
	
	// Get the transaction from the state
	retrievedTx, err := state.GetTransaction(tx.ID.String())
	if err != nil {
		t.Fatalf("Failed to get transaction from state: %v", err)
	}
	
	// Verify the transaction
	if retrievedTx.ID != tx.ID {
		t.Errorf("Retrieved transaction has incorrect ID: got %s, want %s", retrievedTx.ID, tx.ID)
	}
	
	// Create a block with the transaction
	previousHash := crypto.ZeroHash()
	block := state.NewBlock(previousHash, []*Transaction{tx}, 1)
	
	// Add the block to the state
	err = state.AddBlock(block)
	if err != nil {
		t.Fatalf("Failed to add block to state: %v", err)
	}
	
	// Get the block from the state
	blockHash := block.ComputeHash().String()
	retrievedBlock, err := state.GetBlock(blockHash)
	if err != nil {
		t.Fatalf("Failed to get block from state: %v", err)
	}
	
	// Verify the block
	if retrievedBlock.Header.Height != block.Header.Height {
		t.Errorf("Retrieved block has incorrect height: got %d, want %d", retrievedBlock.Header.Height, block.Header.Height)
	}
	
	// Get the latest block
	latestBlock := state.GetLatestBlock()
	if latestBlock == nil {
		t.Fatalf("Latest block is nil")
	}
	
	// Verify the latest block
	if latestBlock.Header.Height != block.Header.Height {
		t.Errorf("Latest block has incorrect height: got %d, want %d", latestBlock.Header.Height, block.Header.Height)
	}
	
	// Get the state root
	stateRoot := state.GetStateRoot()
	if stateRoot == crypto.ZeroHash() {
		t.Errorf("State root is zero hash")
	}
	
	// Generate a state proof
	proof, err := state.GenerateStateProof(publicKey)
	if err != nil {
		t.Fatalf("Failed to generate state proof: %v", err)
	}
	
	// Verify the state proof
	valid, err := state.VerifyStateProof(publicKey, proof)
	if err != nil {
		t.Fatalf("Failed to verify state proof: %v", err)
	}
	if !valid {
		t.Errorf("State proof verification failed")
	}
}
