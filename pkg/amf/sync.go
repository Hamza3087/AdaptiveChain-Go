package amf

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/advanced-blockchain/pkg/crypto"
)

// SyncStatus represents the status of a synchronization operation
type SyncStatus string

const (
	// SyncStatusPending indicates a pending synchronization
	SyncStatusPending SyncStatus = "pending"

	// SyncStatusInProgress indicates a synchronization in progress
	SyncStatusInProgress SyncStatus = "in_progress"

	// SyncStatusCompleted indicates a completed synchronization
	SyncStatusCompleted SyncStatus = "completed"

	// SyncStatusFailed indicates a failed synchronization
	SyncStatusFailed SyncStatus = "failed"
)

// SyncOperation represents a cross-shard synchronization operation
type SyncOperation struct {
	ID            string     // Operation ID
	SourceShardID string     // Source shard ID
	TargetShardID string     // Target shard ID
	Keys          []string   // Keys to synchronize
	Status        SyncStatus // Operation status
	StartTime     time.Time  // Start time
	EndTime       time.Time  // End time
	Error         string     // Error message, if any
}

// ShardSynchronizer handles cross-shard synchronization
type ShardSynchronizer struct {
	forest     *AdaptiveMerkleForest     // Reference to the forest
	operations map[string]*SyncOperation // Map of operation ID to operation
	mu         sync.RWMutex              // Mutex for thread safety
}

// NewShardSynchronizer creates a new shard synchronizer
func NewShardSynchronizer(forest *AdaptiveMerkleForest) *ShardSynchronizer {
	return &ShardSynchronizer{
		forest:     forest,
		operations: make(map[string]*SyncOperation),
	}
}

// SyncShards synchronizes state between two shards
func (s *ShardSynchronizer) SyncShards(sourceID, targetID string, keys []string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a new operation
	opID := fmt.Sprintf("sync-%d", time.Now().UnixNano())
	op := &SyncOperation{
		ID:            opID,
		SourceShardID: sourceID,
		TargetShardID: targetID,
		Keys:          keys,
		Status:        SyncStatusPending,
		StartTime:     time.Now(),
	}

	// Add the operation to the map
	s.operations[opID] = op

	// Start the synchronization in a separate goroutine
	go s.performSync(opID)

	return opID, nil
}

// GetSyncStatus gets the status of a synchronization operation
func (s *ShardSynchronizer) GetSyncStatus(opID string) (*SyncOperation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Get the operation
	op, ok := s.operations[opID]
	if !ok {
		return nil, fmt.Errorf("operation not found: %s", opID)
	}

	return op, nil
}

// performSync performs the actual synchronization
func (s *ShardSynchronizer) performSync(opID string) {
	s.mu.Lock()
	op := s.operations[opID]
	op.Status = SyncStatusInProgress
	s.mu.Unlock()

	// Get the source and target shards
	sourceShard, err := s.forest.GetShard(op.SourceShardID)
	if err != nil {
		s.setSyncError(opID, fmt.Sprintf("failed to get source shard: %v", err))
		return
	}

	targetShard, err := s.forest.GetShard(op.TargetShardID)
	if err != nil {
		s.setSyncError(opID, fmt.Sprintf("failed to get target shard: %v", err))
		return
	}

	// Synchronize each key
	for _, key := range op.Keys {
		// Get the value from the source shard
		value, err := sourceShard.GetState(key)
		if err != nil {
			s.setSyncError(opID, fmt.Sprintf("failed to get value for key %s: %v", key, err))
			return
		}

		// Generate a proof for the value
		proof, err := sourceShard.GenerateProof(key)
		if err != nil {
			s.setSyncError(opID, fmt.Sprintf("failed to generate proof for key %s: %v", key, err))
			return
		}

		// Serialize the value for verification
		valueStr, err := json.Marshal(value)
		if err != nil {
			s.setSyncError(opID, fmt.Sprintf("failed to marshal value for key %s: %v", key, err))
			return
		}

		// Verify the proof
		valid, err := sourceShard.VerifyProof(key, string(valueStr), proof)
		if err != nil {
			s.setSyncError(opID, fmt.Sprintf("failed to verify proof for key %s: %v", key, err))
			return
		}
		if !valid {
			s.setSyncError(opID, fmt.Sprintf("invalid proof for key %s", key))
			return
		}

		// Update the target shard
		if err := targetShard.UpdateState(key, value); err != nil {
			s.setSyncError(opID, fmt.Sprintf("failed to update target shard for key %s: %v", key, err))
			return
		}
	}

	// Mark the operation as completed
	s.mu.Lock()
	op.Status = SyncStatusCompleted
	op.EndTime = time.Now()
	s.mu.Unlock()
}

// setSyncError sets an error for a synchronization operation
func (s *ShardSynchronizer) setSyncError(opID, errorMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	op := s.operations[opID]
	op.Status = SyncStatusFailed
	op.Error = errorMsg
	op.EndTime = time.Now()
}

// SyncCommitment represents a cryptographic commitment for cross-shard operations
type SyncCommitment struct {
	OperationID           string      // Operation ID
	SourceShardID         string      // Source shard ID
	TargetShardID         string      // Target shard ID
	Keys                  []string    // Keys being synchronized
	SourceRoot            crypto.Hash // Source shard state root
	TargetRoot            crypto.Hash // Target shard state root
	Signature             []byte      // Signature of the commitment
	HomomorphicCommitment []byte      // Homomorphic commitment for the operation
}

// PartialStateTransfer represents a partial state transfer between shards
type PartialStateTransfer struct {
	OperationID   string                 // Operation ID
	SourceShardID string                 // Source shard ID
	TargetShardID string                 // Target shard ID
	Keys          []string               // Keys being transferred
	Values        map[string]interface{} // Values being transferred
	Proofs        map[string][]byte      // Proofs for the values
	Commitment    *SyncCommitment        // Commitment for the transfer
}

// CreateSyncCommitment creates a cryptographic commitment for a cross-shard operation
func (s *ShardSynchronizer) CreateSyncCommitment(opID string, keyPair *crypto.KeyPair) (*SyncCommitment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Get the operation
	op, ok := s.operations[opID]
	if !ok {
		return nil, fmt.Errorf("operation not found: %s", opID)
	}

	// Check if the operation is completed
	if op.Status != SyncStatusCompleted {
		return nil, fmt.Errorf("operation not completed: %s", opID)
	}

	// Get the source and target shards
	sourceShard, err := s.forest.GetShard(op.SourceShardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source shard: %w", err)
	}

	targetShard, err := s.forest.GetShard(op.TargetShardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get target shard: %w", err)
	}

	// Create a homomorphic authenticated data structure
	hads, err := crypto.NewHomomorphicAuthenticatedDataStructure()
	if err != nil {
		return nil, fmt.Errorf("failed to create homomorphic authenticated data structure: %w", err)
	}

	// Add each key-value pair to the data structure
	for _, key := range op.Keys {
		// Get the value from the source shard
		value, err := sourceShard.GetState(key)
		if err != nil {
			return nil, fmt.Errorf("failed to get value for key %s: %w", key, err)
		}

		// Serialize the value
		valueBytes, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}

		// Add the key-value pair to the data structure
		if err := hads.Add(key, valueBytes); err != nil {
			return nil, fmt.Errorf("failed to add key-value pair to data structure: %w", err)
		}
	}

	// Compute a hash of the data structure
	hadsHash := hads.ComputeHash()

	// Create the commitment
	commitment := &SyncCommitment{
		OperationID:           opID,
		SourceShardID:         op.SourceShardID,
		TargetShardID:         op.TargetShardID,
		Keys:                  op.Keys,
		SourceRoot:            sourceShard.GetStateRoot(),
		TargetRoot:            targetShard.GetStateRoot(),
		HomomorphicCommitment: hadsHash.Bytes(),
	}

	// Serialize the commitment
	data, err := json.Marshal(commitment)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal commitment: %w", err)
	}

	// Sign the commitment
	signature, err := keyPair.Sign(data)
	if err != nil {
		return nil, fmt.Errorf("failed to sign commitment: %w", err)
	}

	// Set the signature
	commitment.Signature = signature

	return commitment, nil
}

// VerifySyncCommitment verifies a cryptographic commitment for a cross-shard operation
func (s *ShardSynchronizer) VerifySyncCommitment(commitment *SyncCommitment, publicKey []byte) (bool, error) {
	// Serialize the commitment without the signature
	commitmentCopy := *commitment
	commitmentCopy.Signature = nil

	data, err := json.Marshal(commitmentCopy)
	if err != nil {
		return false, fmt.Errorf("failed to marshal commitment: %w", err)
	}

	// Verify the signature
	return crypto.VerifySignature(publicKey, data, commitment.Signature)
}

// CreatePartialStateTransfer creates a partial state transfer between shards
func (s *ShardSynchronizer) CreatePartialStateTransfer(opID string, keyPair *crypto.KeyPair) (*PartialStateTransfer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Get the operation
	op, ok := s.operations[opID]
	if !ok {
		return nil, fmt.Errorf("operation not found: %s", opID)
	}

	// Check if the operation is completed
	if op.Status != SyncStatusCompleted {
		return nil, fmt.Errorf("operation not completed: %s", opID)
	}

	// Get the source shard
	sourceShard, err := s.forest.GetShard(op.SourceShardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source shard: %w", err)
	}

	// Create a commitment for the operation
	commitment, err := s.CreateSyncCommitment(opID, keyPair)
	if err != nil {
		return nil, fmt.Errorf("failed to create commitment: %w", err)
	}

	// Create a partial state transfer
	transfer := &PartialStateTransfer{
		OperationID:   opID,
		SourceShardID: op.SourceShardID,
		TargetShardID: op.TargetShardID,
		Keys:          op.Keys,
		Values:        make(map[string]interface{}),
		Proofs:        make(map[string][]byte),
		Commitment:    commitment,
	}

	// Add each key-value pair to the transfer
	for _, key := range op.Keys {
		// Get the value from the source shard
		value, err := sourceShard.GetState(key)
		if err != nil {
			return nil, fmt.Errorf("failed to get value for key %s: %w", key, err)
		}

		// Generate a proof for the value
		proof, err := sourceShard.GenerateProof(key)
		if err != nil {
			return nil, fmt.Errorf("failed to generate proof for key %s: %w", key, err)
		}

		// Add the key-value pair and proof to the transfer
		transfer.Values[key] = value
		transfer.Proofs[key] = proof
	}

	return transfer, nil
}

// ApplyPartialStateTransfer applies a partial state transfer to a shard
func (s *ShardSynchronizer) ApplyPartialStateTransfer(transfer *PartialStateTransfer, publicKey []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verify the commitment
	valid, err := s.VerifySyncCommitment(transfer.Commitment, publicKey)
	if err != nil {
		return fmt.Errorf("failed to verify commitment: %w", err)
	}
	if !valid {
		return fmt.Errorf("invalid commitment")
	}

	// Get the target shard
	targetShard, err := s.forest.GetShard(transfer.TargetShardID)
	if err != nil {
		return fmt.Errorf("failed to get target shard: %w", err)
	}

	// Apply each key-value pair to the target shard
	for key, value := range transfer.Values {
		// Verify the proof
		proof := transfer.Proofs[key]
		valueStr, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}

		// Get the source shard
		sourceShard, err := s.forest.GetShard(transfer.SourceShardID)
		if err != nil {
			return fmt.Errorf("failed to get source shard: %w", err)
		}

		// Verify the proof
		valid, err := sourceShard.VerifyProof(key, string(valueStr), proof)
		if err != nil {
			return fmt.Errorf("failed to verify proof for key %s: %w", key, err)
		}
		if !valid {
			return fmt.Errorf("invalid proof for key %s", key)
		}

		// Update the target shard
		if err := targetShard.UpdateState(key, value); err != nil {
			return fmt.Errorf("failed to update target shard for key %s: %w", key, err)
		}
	}

	// Create a new operation for the transfer
	opID := fmt.Sprintf("transfer-%d", time.Now().UnixNano())
	op := &SyncOperation{
		ID:            opID,
		SourceShardID: transfer.SourceShardID,
		TargetShardID: transfer.TargetShardID,
		Keys:          transfer.Keys,
		Status:        SyncStatusCompleted,
		StartTime:     time.Now(),
		EndTime:       time.Now(),
	}

	// Add the operation to the map
	s.operations[opID] = op

	return nil
}

// AtomicCrossShardOperation represents an atomic operation across multiple shards
type AtomicCrossShardOperation struct {
	ID          string                     // Operation ID
	ShardIDs    []string                   // Shard IDs involved in the operation
	Keys        map[string][]string        // Keys to operate on for each shard
	Values      map[string]interface{}     // Values to set for each key
	Commitments map[string]*SyncCommitment // Commitments for each shard
	Status      SyncStatus                 // Operation status
	StartTime   time.Time                  // Start time
	EndTime     time.Time                  // End time
	Error       string                     // Error message, if any
}

// CreateAtomicCrossShardOperation creates an atomic operation across multiple shards
func (s *ShardSynchronizer) CreateAtomicCrossShardOperation(shardIDs []string, keys map[string][]string, values map[string]interface{}) (*AtomicCrossShardOperation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a new operation
	opID := fmt.Sprintf("atomic-%d", time.Now().UnixNano())
	op := &AtomicCrossShardOperation{
		ID:          opID,
		ShardIDs:    shardIDs,
		Keys:        keys,
		Values:      values,
		Commitments: make(map[string]*SyncCommitment),
		Status:      SyncStatusPending,
		StartTime:   time.Now(),
	}

	// Start the operation in a separate goroutine
	go s.performAtomicOperation(op)

	return op, nil
}

// performAtomicOperation performs an atomic operation across multiple shards
func (s *ShardSynchronizer) performAtomicOperation(op *AtomicCrossShardOperation) {
	s.mu.Lock()
	op.Status = SyncStatusInProgress
	s.mu.Unlock()

	// Create a homomorphic authenticated data structure
	hads, err := crypto.NewHomomorphicAuthenticatedDataStructure()
	if err != nil {
		s.setAtomicOperationError(op, fmt.Sprintf("failed to create homomorphic authenticated data structure: %v", err))
		return
	}

	// Phase 1: Prepare - add all key-value pairs to the data structure
	for shardID, shardKeys := range op.Keys {
		// Get the shard
		_, err := s.forest.GetShard(shardID)
		if err != nil {
			s.setAtomicOperationError(op, fmt.Sprintf("failed to get shard %s: %v", shardID, err))
			return
		}

		// Add each key-value pair to the data structure
		for _, key := range shardKeys {
			// Get the value
			value, ok := op.Values[key]
			if !ok {
				s.setAtomicOperationError(op, fmt.Sprintf("value not found for key %s", key))
				return
			}

			// Serialize the value
			valueBytes, err := json.Marshal(value)
			if err != nil {
				s.setAtomicOperationError(op, fmt.Sprintf("failed to marshal value for key %s: %v", key, err))
				return
			}

			// Add the key-value pair to the data structure
			if err := hads.Add(key, valueBytes); err != nil {
				s.setAtomicOperationError(op, fmt.Sprintf("failed to add key-value pair to data structure: %v", err))
				return
			}
		}
	}

	// Phase 2: Commit - update each shard and create commitments
	for shardID, shardKeys := range op.Keys {
		// Get the shard
		shard, err := s.forest.GetShard(shardID)
		if err != nil {
			s.setAtomicOperationError(op, fmt.Sprintf("failed to get shard %s: %v", shardID, err))
			return
		}

		// Update each key-value pair in the shard
		for _, key := range shardKeys {
			// Get the value
			value, ok := op.Values[key]
			if !ok {
				s.setAtomicOperationError(op, fmt.Sprintf("value not found for key %s", key))
				return
			}

			// Update the shard
			if err := shard.UpdateState(key, value); err != nil {
				s.setAtomicOperationError(op, fmt.Sprintf("failed to update shard %s for key %s: %v", shardID, key, err))
				return
			}
		}

		// Create a commitment for the shard
		commitment := &SyncCommitment{
			OperationID:           op.ID,
			SourceShardID:         shardID,
			TargetShardID:         "",
			Keys:                  shardKeys,
			SourceRoot:            shard.GetStateRoot(),
			TargetRoot:            crypto.ZeroHash(),
			HomomorphicCommitment: hads.ComputeHash().Bytes(),
		}

		// Add the commitment to the operation
		op.Commitments[shardID] = commitment
	}

	// Mark the operation as completed
	s.mu.Lock()
	op.Status = SyncStatusCompleted
	op.EndTime = time.Now()
	s.mu.Unlock()
}

// setAtomicOperationError sets an error for an atomic operation
func (s *ShardSynchronizer) setAtomicOperationError(op *AtomicCrossShardOperation, errorMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	op.Status = SyncStatusFailed
	op.Error = errorMsg
	op.EndTime = time.Now()
}

// GetAtomicOperationStatus gets the status of an atomic operation
func (s *ShardSynchronizer) GetAtomicOperationStatus(opID string) (*AtomicCrossShardOperation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Find the operation
	for _, op := range s.operations {
		if op.ID == opID {
			// Convert the operation to an atomic operation
			atomicOp := &AtomicCrossShardOperation{
				ID:        opID,
				Status:    op.Status,
				StartTime: op.StartTime,
				EndTime:   op.EndTime,
				Error:     op.Error,
			}
			return atomicOp, nil
		}
	}

	return nil, fmt.Errorf("operation not found: %s", opID)
}
