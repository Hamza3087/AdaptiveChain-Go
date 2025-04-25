package consensus

import (
	"testing"

	"github.com/advanced-blockchain/pkg/crypto"
	"github.com/advanced-blockchain/pkg/state"
)

func TestHybridConsensus(t *testing.T) {
	// Create a new state
	state := state.NewState()
	
	// Create a new consensus
	consensus := NewHybridConsensus(state)
	
	// Generate key pairs for validators
	keyPair1, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair 1: %v", err)
	}
	
	keyPair2, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair 2: %v", err)
	}
	
	keyPair3, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair 3: %v", err)
	}
	
	// Get the public keys
	publicKey1, err := keyPair1.ExportPublicKey()
	if err != nil {
		t.Fatalf("Failed to export public key 1: %v", err)
	}
	
	publicKey2, err := keyPair2.ExportPublicKey()
	if err != nil {
		t.Fatalf("Failed to export public key 2: %v", err)
	}
	
	publicKey3, err := keyPair3.ExportPublicKey()
	if err != nil {
		t.Fatalf("Failed to export public key 3: %v", err)
	}
	
	// Add validators
	err = consensus.AddValidator("validator-1", publicKey1, 100)
	if err != nil {
		t.Fatalf("Failed to add validator 1: %v", err)
	}
	
	err = consensus.AddValidator("validator-2", publicKey2, 100)
	if err != nil {
		t.Fatalf("Failed to add validator 2: %v", err)
	}
	
	err = consensus.AddValidator("validator-3", publicKey3, 100)
	if err != nil {
		t.Fatalf("Failed to add validator 3: %v", err)
	}
	
	// Start a consensus round
	err = consensus.StartConsensus()
	if err != nil {
		t.Fatalf("Failed to start consensus: %v", err)
	}
	
	// Verify the consensus status
	if consensus.Status != ConsensusStatusPreparing {
		t.Errorf("Consensus has incorrect status: got %s, want %s", consensus.Status, ConsensusStatusPreparing)
	}
	
	// Submit prepare votes
	err = consensus.PrepareVote("validator-1", true)
	if err != nil {
		t.Fatalf("Failed to submit prepare vote 1: %v", err)
	}
	
	err = consensus.PrepareVote("validator-2", true)
	if err != nil {
		t.Fatalf("Failed to submit prepare vote 2: %v", err)
	}
	
	// Verify the consensus status
	if consensus.Status != ConsensusStatusCommitting {
		t.Errorf("Consensus has incorrect status: got %s, want %s", consensus.Status, ConsensusStatusCommitting)
	}
	
	// Submit commit votes
	err = consensus.CommitVote("validator-1", true)
	if err != nil {
		t.Fatalf("Failed to submit commit vote 1: %v", err)
	}
	
	err = consensus.CommitVote("validator-2", true)
	if err != nil {
		t.Fatalf("Failed to submit commit vote 2: %v", err)
	}
	
	// Verify the consensus status
	if consensus.Status != ConsensusStatusFinalized {
		t.Errorf("Consensus has incorrect status: got %s, want %s", consensus.Status, ConsensusStatusFinalized)
	}
	
	// Verify the finalized block
	if consensus.FinalizedBlock == nil {
		t.Fatalf("Finalized block is nil")
	}
	
	// Verify the block was added to the state
	latestBlock := state.GetLatestBlock()
	if latestBlock == nil {
		t.Fatalf("Latest block is nil")
	}
	
	// Verify the block height
	if latestBlock.Header.Height != 0 {
		t.Errorf("Latest block has incorrect height: got %d, want 0", latestBlock.Header.Height)
	}
}
