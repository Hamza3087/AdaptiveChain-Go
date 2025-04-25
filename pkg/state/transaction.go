package state

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/advanced-blockchain/pkg/crypto"
)

// Transaction represents a blockchain transaction
type Transaction struct {
	ID        crypto.Hash // Transaction ID (hash of the transaction)
	From      []byte      // Sender's public key
	To        []byte      // Recipient's public key
	Amount    uint64      // Transaction amount
	Data      []byte      // Additional data
	Timestamp time.Time   // Transaction timestamp
	ShardID   uint64      // Shard ID
	Signature []byte      // Transaction signature
}

// NewTransaction creates a new transaction
func NewTransaction(from, to []byte, amount uint64, data []byte, shardID uint64) *Transaction {
	tx := &Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Data:      data,
		Timestamp: time.Now(),
		ShardID:   shardID,
	}

	// Compute the transaction ID
	tx.ID = tx.ComputeHash()

	return tx
}

// ComputeHash computes the hash of the transaction
func (tx *Transaction) ComputeHash() crypto.Hash {
	// Marshal the transaction without the signature and ID
	txCopy := *tx
	txCopy.Signature = nil
	txCopy.ID = crypto.ZeroHash()

	data, err := json.Marshal(txCopy)
	if err != nil {
		// In case of error, return a zero hash
		return crypto.ZeroHash()
	}

	return crypto.NewHash(data)
}

// Sign signs the transaction with the given key pair
func (tx *Transaction) Sign(keyPair *crypto.KeyPair) error {
	// Compute the hash of the transaction
	hash := tx.ComputeHash()

	// Sign the hash
	signature, err := keyPair.Sign(hash[:])
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Set the signature
	tx.Signature = signature

	return nil
}

// Verify verifies the transaction signature
func (tx *Transaction) Verify() (bool, error) {
	// Check if the transaction has a signature
	if tx.Signature == nil {
		return false, fmt.Errorf("transaction has no signature")
	}

	// Compute the hash of the transaction
	hash := tx.ComputeHash()

	// Verify the signature
	return crypto.VerifySignature(tx.From, hash[:], tx.Signature)
}

// MarshalJSON marshals the transaction to JSON
func (tx *Transaction) MarshalJSON() ([]byte, error) {
	type Alias Transaction
	return json.Marshal(&struct {
		ID        string `json:"id"`
		From      string `json:"from"`
		To        string `json:"to"`
		Signature string `json:"signature,omitempty"`
		Timestamp string `json:"timestamp"`
		*Alias
	}{
		ID:        tx.ID.String(),
		From:      fmt.Sprintf("%x", tx.From),
		To:        fmt.Sprintf("%x", tx.To),
		Signature: fmt.Sprintf("%x", tx.Signature),
		Timestamp: tx.Timestamp.Format(time.RFC3339),
		Alias:     (*Alias)(tx),
	})
}

// UnmarshalJSON unmarshals the transaction from JSON
func (tx *Transaction) UnmarshalJSON(data []byte) error {
	type Alias Transaction
	aux := &struct {
		ID        string `json:"id"`
		From      string `json:"from"`
		To        string `json:"to"`
		Signature string `json:"signature,omitempty"`
		Timestamp string `json:"timestamp"`
		*Alias
	}{
		Alias: (*Alias)(tx),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse the ID
	id, err := crypto.HashFromString(aux.ID)
	if err != nil {
		return err
	}
	tx.ID = id

	// Parse the timestamp
	timestamp, err := time.Parse(time.RFC3339, aux.Timestamp)
	if err != nil {
		return err
	}
	tx.Timestamp = timestamp

	return nil
}

// Marshal marshals the transaction to JSON
func (tx *Transaction) Marshal() ([]byte, error) {
	return json.Marshal(tx)
}

// Unmarshal unmarshals the transaction from JSON
func (tx *Transaction) Unmarshal(data []byte) error {
	return json.Unmarshal(data, tx)
}
