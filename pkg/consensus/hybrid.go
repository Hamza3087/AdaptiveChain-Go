package consensus

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/advanced-blockchain/pkg/crypto"
	"github.com/advanced-blockchain/pkg/state"
)

// ConsensusStatus represents the status of the consensus
type ConsensusStatus string

const (
	// ConsensusStatusIdle indicates an idle consensus
	ConsensusStatusIdle ConsensusStatus = "idle"
	
	// ConsensusStatusPreparing indicates a preparing consensus
	ConsensusStatusPreparing ConsensusStatus = "preparing"
	
	// ConsensusStatusCommitting indicates a committing consensus
	ConsensusStatusCommitting ConsensusStatus = "committing"
	
	// ConsensusStatusFinalized indicates a finalized consensus
	ConsensusStatusFinalized ConsensusStatus = "finalized"
)

// HybridConsensus represents a hybrid consensus mechanism
// combining Proof of Work (PoW) and Delegated Byzantine Fault Tolerance (dBFT)
type HybridConsensus struct {
	State            *state.State        // Blockchain state
	Validators       []*Validator        // List of validators
	ValidatorMap     map[string]*Validator // Map of validator ID to validator
	LeaderID         string              // Current leader ID
	Status           ConsensusStatus     // Consensus status
	Round            uint64              // Current consensus round
	Difficulty       uint64              // Current mining difficulty
	PrepareThreshold uint64              // Threshold for prepare phase
	CommitThreshold  uint64              // Threshold for commit phase
	PrepareVotes     map[string]bool     // Map of validator ID to prepare vote
	CommitVotes      map[string]bool     // Map of validator ID to commit vote
	ProposedBlock    *state.Block        // Proposed block
	FinalizedBlock   *state.Block        // Finalized block
	mu               sync.RWMutex        // Mutex for thread safety
}

// Validator represents a consensus validator
type Validator struct {
	ID            string          // Validator ID
	PublicKey     []byte          // Validator public key
	Stake         uint64          // Validator stake
	ReputationScore float64       // Validator reputation score
	LastActive    time.Time       // Last active time
}

// NewHybridConsensus creates a new hybrid consensus
func NewHybridConsensus(state *state.State) *HybridConsensus {
	return &HybridConsensus{
		State:            state,
		Validators:       make([]*Validator, 0),
		ValidatorMap:     make(map[string]*Validator),
		Status:           ConsensusStatusIdle,
		Round:            0,
		Difficulty:       1000000, // Initial difficulty
		PrepareThreshold: 2,       // Initial prepare threshold (2/3 of validators)
		CommitThreshold:  2,       // Initial commit threshold (2/3 of validators)
		PrepareVotes:     make(map[string]bool),
		CommitVotes:      make(map[string]bool),
	}
}

// AddValidator adds a validator to the consensus
func (c *HybridConsensus) AddValidator(id string, publicKey []byte, stake uint64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Check if the validator already exists
	if _, ok := c.ValidatorMap[id]; ok {
		return fmt.Errorf("validator already exists: %s", id)
	}
	
	// Create a new validator
	validator := &Validator{
		ID:              id,
		PublicKey:       publicKey,
		Stake:           stake,
		ReputationScore: 1.0, // Initial reputation score
		LastActive:      time.Now(),
	}
	
	// Add the validator to the list and map
	c.Validators = append(c.Validators, validator)
	c.ValidatorMap[id] = validator
	
	// Update the thresholds
	c.updateThresholds()
	
	return nil
}

// RemoveValidator removes a validator from the consensus
func (c *HybridConsensus) RemoveValidator(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Check if the validator exists
	if _, ok := c.ValidatorMap[id]; !ok {
		return fmt.Errorf("validator not found: %s", id)
	}
	
	// Remove the validator from the map
	delete(c.ValidatorMap, id)
	
	// Remove the validator from the list
	var newValidators []*Validator
	for _, v := range c.Validators {
		if v.ID != id {
			newValidators = append(newValidators, v)
		}
	}
	c.Validators = newValidators
	
	// Update the thresholds
	c.updateThresholds()
	
	return nil
}

// updateThresholds updates the prepare and commit thresholds
func (c *HybridConsensus) updateThresholds() {
	// Set the thresholds to 2/3 of the number of validators
	numValidators := len(c.Validators)
	c.PrepareThreshold = uint64((numValidators * 2) / 3)
	c.CommitThreshold = uint64((numValidators * 2) / 3)
	
	// Ensure the thresholds are at least 1
	if c.PrepareThreshold < 1 {
		c.PrepareThreshold = 1
	}
	if c.CommitThreshold < 1 {
		c.CommitThreshold = 1
	}
}

// StartConsensus starts a new consensus round
func (c *HybridConsensus) StartConsensus() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Check if a consensus is already in progress
	if c.Status != ConsensusStatusIdle {
		return fmt.Errorf("consensus already in progress: %s", c.Status)
	}
	
	// Increment the round
	c.Round++
	
	// Reset the votes
	c.PrepareVotes = make(map[string]bool)
	c.CommitVotes = make(map[string]bool)
	
	// Set the status to preparing
	c.Status = ConsensusStatusPreparing
	
	// Elect a leader
	if err := c.electLeader(); err != nil {
		return fmt.Errorf("failed to elect leader: %w", err)
	}
	
	// Create a proposed block
	if err := c.createProposedBlock(); err != nil {
		return fmt.Errorf("failed to create proposed block: %w", err)
	}
	
	return nil
}

// electLeader elects a leader for the current round
func (c *HybridConsensus) electLeader() error {
	// Check if there are any validators
	if len(c.Validators) == 0 {
		return fmt.Errorf("no validators available")
	}
	
	// Use a verifiable random function (VRF) to elect a leader
	// For simplicity, we'll just use a random number
	// In a real implementation, we would use a proper VRF
	
	// Generate a random number
	max := big.NewInt(int64(len(c.Validators)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return fmt.Errorf("failed to generate random number: %w", err)
	}
	
	// Select the leader
	leaderIndex := n.Int64()
	c.LeaderID = c.Validators[leaderIndex].ID
	
	return nil
}

// createProposedBlock creates a proposed block for the current round
func (c *HybridConsensus) createProposedBlock() error {
	// Get the latest block
	latestBlock := c.State.GetLatestBlock()
	var previousHash crypto.Hash
	if latestBlock != nil {
		previousHash = latestBlock.ComputeHash()
	} else {
		previousHash = crypto.ZeroHash()
	}
	
	// Create a new block
	c.ProposedBlock = state.NewBlock(previousHash, nil, 0)
	
	// Set the block height
	if latestBlock != nil {
		c.ProposedBlock.Header.Height = latestBlock.Header.Height + 1
	} else {
		c.ProposedBlock.Header.Height = 0
	}
	
	// Set the difficulty
	c.ProposedBlock.Header.Difficulty = c.Difficulty
	
	return nil
}

// PrepareVote submits a prepare vote for the proposed block
func (c *HybridConsensus) PrepareVote(validatorID string, vote bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Check if the consensus is in the preparing phase
	if c.Status != ConsensusStatusPreparing {
		return fmt.Errorf("consensus not in preparing phase: %s", c.Status)
	}
	
	// Check if the validator exists
	validator, ok := c.ValidatorMap[validatorID]
	if !ok {
		return fmt.Errorf("validator not found: %s", validatorID)
	}
	
	// Record the vote
	c.PrepareVotes[validatorID] = vote
	
	// Update the validator's last active time
	validator.LastActive = time.Now()
	
	// Check if we have enough prepare votes
	if c.countPrepareVotes() >= c.PrepareThreshold {
		// Move to the committing phase
		c.Status = ConsensusStatusCommitting
	}
	
	return nil
}

// CommitVote submits a commit vote for the proposed block
func (c *HybridConsensus) CommitVote(validatorID string, vote bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Check if the consensus is in the committing phase
	if c.Status != ConsensusStatusCommitting {
		return fmt.Errorf("consensus not in committing phase: %s", c.Status)
	}
	
	// Check if the validator exists
	validator, ok := c.ValidatorMap[validatorID]
	if !ok {
		return fmt.Errorf("validator not found: %s", validatorID)
	}
	
	// Record the vote
	c.CommitVotes[validatorID] = vote
	
	// Update the validator's last active time
	validator.LastActive = time.Now()
	
	// Check if we have enough commit votes
	if c.countCommitVotes() >= c.CommitThreshold {
		// Finalize the block
		c.FinalizedBlock = c.ProposedBlock
		c.Status = ConsensusStatusFinalized
		
		// Add the block to the state
		if err := c.State.AddBlock(c.FinalizedBlock); err != nil {
			// Log the error
			fmt.Printf("Error adding block to state: %v\n", err)
		}
		
		// Update the difficulty
		c.updateDifficulty()
		
		// Update the validator reputation scores
		c.updateReputationScores()
	}
	
	return nil
}

// countPrepareVotes counts the number of positive prepare votes
func (c *HybridConsensus) countPrepareVotes() uint64 {
	count := uint64(0)
	for _, vote := range c.PrepareVotes {
		if vote {
			count++
		}
	}
	return count
}

// countCommitVotes counts the number of positive commit votes
func (c *HybridConsensus) countCommitVotes() uint64 {
	count := uint64(0)
	for _, vote := range c.CommitVotes {
		if vote {
			count++
		}
	}
	return count
}

// updateDifficulty updates the mining difficulty
func (c *HybridConsensus) updateDifficulty() {
	// In a real implementation, we would adjust the difficulty based on the block time
	// For simplicity, we'll just keep it constant
}

// updateReputationScores updates the reputation scores of the validators
func (c *HybridConsensus) updateReputationScores() {
	// Update the reputation scores based on the votes
	for _, validator := range c.Validators {
		// Check if the validator participated in the prepare phase
		if prepareVote, ok := c.PrepareVotes[validator.ID]; ok {
			if prepareVote {
				// Increase the reputation score for a positive vote
				validator.ReputationScore *= 1.01
			} else {
				// Decrease the reputation score for a negative vote
				validator.ReputationScore *= 0.99
			}
		} else {
			// Decrease the reputation score for not participating
			validator.ReputationScore *= 0.95
		}
		
		// Check if the validator participated in the commit phase
		if commitVote, ok := c.CommitVotes[validator.ID]; ok {
			if commitVote {
				// Increase the reputation score for a positive vote
				validator.ReputationScore *= 1.01
			} else {
				// Decrease the reputation score for a negative vote
				validator.ReputationScore *= 0.99
			}
		} else {
			// Decrease the reputation score for not participating
			validator.ReputationScore *= 0.95
		}
		
		// Ensure the reputation score is within bounds
		if validator.ReputationScore < 0.1 {
			validator.ReputationScore = 0.1
		}
		if validator.ReputationScore > 10.0 {
			validator.ReputationScore = 10.0
		}
	}
}

// GetConsensusStatus gets the current consensus status
func (c *HybridConsensus) GetConsensusStatus() ConsensusStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return c.Status
}

// GetLeaderID gets the current leader ID
func (c *HybridConsensus) GetLeaderID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return c.LeaderID
}

// GetRound gets the current consensus round
func (c *HybridConsensus) GetRound() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return c.Round
}

// GetProposedBlock gets the proposed block
func (c *HybridConsensus) GetProposedBlock() *state.Block {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return c.ProposedBlock
}

// GetFinalizedBlock gets the finalized block
func (c *HybridConsensus) GetFinalizedBlock() *state.Block {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return c.FinalizedBlock
}
