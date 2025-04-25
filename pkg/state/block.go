package state

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/advanced-blockchain/pkg/crypto"
)

// Block represents a blockchain block
type Block struct {
	Header       BlockHeader       // Block header
	Transactions []*Transaction    // Block transactions
	StateRoot    crypto.Hash       // Merkle root of the state trie
	ShardID      uint64            // Shard ID
	Signature    []byte            // Block signature
	Metadata     map[string][]byte // Additional metadata
}

// BlockHeader represents a blockchain block header
type BlockHeader struct {
	Version      uint32      // Block version
	PreviousHash crypto.Hash // Hash of the previous block
	MerkleRoot   crypto.Hash // Merkle root of the transactions
	Timestamp    time.Time   // Block timestamp
	Height       uint64      // Block height
	Difficulty   uint64      // Mining difficulty
	Nonce        uint64      // Nonce for mining
}

// NewBlock creates a new block
func NewBlock(previousHash crypto.Hash, transactions []*Transaction, shardID uint64) *Block {
	// Create a new block
	block := &Block{
		Header: BlockHeader{
			Version:      1,
			PreviousHash: previousHash,
			Timestamp:    time.Now(),
			Height:       0, // Will be set by the blockchain
			Difficulty:   0, // Will be set by the blockchain
			Nonce:        0, // Will be set by mining
		},
		Transactions: transactions,
		ShardID:      shardID,
		Metadata:     make(map[string][]byte),
	}

	// Compute the Merkle root of the transactions
	block.Header.MerkleRoot = block.ComputeMerkleRoot()

	return block
}

// ComputeHash computes the hash of the block
func (b *Block) ComputeHash() crypto.Hash {
	// Marshal the block header
	data, err := json.Marshal(b.Header)
	if err != nil {
		// In case of error, return a zero hash
		return crypto.ZeroHash()
	}

	return crypto.NewHash(data)
}

// ComputeMerkleRoot computes the Merkle root of the transactions
func (b *Block) ComputeMerkleRoot() crypto.Hash {
	if len(b.Transactions) == 0 {
		return crypto.ZeroHash()
	}

	// Compute the hashes of the transactions
	var hashes []crypto.Hash
	for _, tx := range b.Transactions {
		hashes = append(hashes, tx.ID)
	}

	// Compute the Merkle root
	return crypto.ComputeMerkleRoot(hashes)
}

// Sign signs the block with the given key pair
func (b *Block) Sign(keyPair *crypto.KeyPair) error {
	// Compute the hash of the block
	hash := b.ComputeHash()

	// Sign the hash
	signature, err := keyPair.Sign(hash[:])
	if err != nil {
		return fmt.Errorf("failed to sign block: %w", err)
	}

	// Set the signature
	b.Signature = signature

	return nil
}

// Verify verifies the block signature
func (b *Block) Verify(publicKey []byte) (bool, error) {
	// Check if the block has a signature
	if b.Signature == nil {
		return false, fmt.Errorf("block has no signature")
	}

	// Compute the hash of the block
	hash := b.ComputeHash()

	// Verify the signature
	return crypto.VerifySignature(publicKey, hash[:], b.Signature)
}

// AddTransaction adds a transaction to the block
func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)

	// Update the Merkle root
	b.Header.MerkleRoot = b.ComputeMerkleRoot()
}

// SetMetadata sets a metadata value
func (b *Block) SetMetadata(key string, value []byte) {
	b.Metadata[key] = value
}

// GetMetadata gets a metadata value
func (b *Block) GetMetadata(key string) []byte {
	return b.Metadata[key]
}

// MarshalJSON marshals the block to JSON
func (b *Block) MarshalJSON() ([]byte, error) {
	type Alias Block
	return json.Marshal(&struct {
		Header       BlockHeader       `json:"header"`
		Transactions []*Transaction    `json:"transactions"`
		StateRoot    string            `json:"stateRoot"`
		ShardID      uint64            `json:"shardId"`
		Signature    string            `json:"signature,omitempty"`
		Metadata     map[string]string `json:"metadata,omitempty"`
		*Alias
	}{
		Header:       b.Header,
		Transactions: b.Transactions,
		StateRoot:    b.StateRoot.String(),
		ShardID:      b.ShardID,
		Signature:    fmt.Sprintf("%x", b.Signature),
		Metadata:     convertMetadataToString(b.Metadata),
		Alias:        (*Alias)(b),
	})
}

// UnmarshalJSON unmarshals the block from JSON
func (b *Block) UnmarshalJSON(data []byte) error {
	type Alias Block
	aux := &struct {
		Header       BlockHeader       `json:"header"`
		Transactions []*Transaction    `json:"transactions"`
		StateRoot    string            `json:"stateRoot"`
		ShardID      uint64            `json:"shardId"`
		Signature    string            `json:"signature,omitempty"`
		Metadata     map[string]string `json:"metadata,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(b),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse the state root
	stateRoot, err := crypto.HashFromString(aux.StateRoot)
	if err != nil {
		return err
	}
	b.StateRoot = stateRoot

	// Convert metadata back to bytes
	b.Metadata = convertMetadataFromString(aux.Metadata)

	return nil
}

// Marshal marshals the block to JSON
func (b *Block) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

// Unmarshal unmarshals the block from JSON
func (b *Block) Unmarshal(data []byte) error {
	return json.Unmarshal(data, b)
}

// Helper function to convert metadata map to string values for JSON
func convertMetadataToString(metadata map[string][]byte) map[string]string {
	result := make(map[string]string)
	for k, v := range metadata {
		result[k] = fmt.Sprintf("%x", v)
	}
	return result
}

// Helper function to convert metadata map from string values from JSON
func convertMetadataFromString(metadata map[string]string) map[string][]byte {
	result := make(map[string][]byte)
	for k, v := range metadata {
		// Parse hex string to bytes
		var bytes []byte
		fmt.Sscanf(v, "%x", &bytes)
		result[k] = bytes
	}
	return result
}
