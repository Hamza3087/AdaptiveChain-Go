package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/advanced-blockchain/pkg/crypto"
)

// State represents the blockchain state
type State struct {
	mu            sync.RWMutex
	accounts      map[string]*Account
	transactions  map[string]*Transaction
	blocks        map[string]*Block
	latestBlock   *Block
	stateRoot     crypto.Hash
	stateAccum    *crypto.Accumulator
	pruningHeight uint64
}

// Account represents a blockchain account
type Account struct {
	Address   []byte
	Balance   uint64
	Nonce     uint64
	Data      map[string][]byte
	StateRoot crypto.Hash
}

// NewState creates a new blockchain state
func NewState() *State {
	// Create a new RSA modulus for the accumulator
	// In a real implementation, this would be a secure RSA modulus
	n := crypto.NewHash([]byte("state-accumulator"))
	nBig := new(big.Int).SetBytes(n.Bytes())

	return &State{
		accounts:      make(map[string]*Account),
		transactions:  make(map[string]*Transaction),
		blocks:        make(map[string]*Block),
		stateRoot:     crypto.ZeroHash(),
		stateAccum:    crypto.NewAccumulator(nBig),
		pruningHeight: 0,
	}
}

// GetAccount gets an account by address
func (s *State) GetAccount(address []byte) (*Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	addressStr := fmt.Sprintf("%x", address)
	account, ok := s.accounts[addressStr]
	if !ok {
		return nil, fmt.Errorf("account not found: %s", addressStr)
	}

	return account, nil
}

// CreateAccount creates a new account
func (s *State) CreateAccount(address []byte) (*Account, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	addressStr := fmt.Sprintf("%x", address)
	if _, ok := s.accounts[addressStr]; ok {
		return nil, fmt.Errorf("account already exists: %s", addressStr)
	}

	account := &Account{
		Address:   address,
		Balance:   0,
		Nonce:     0,
		Data:      make(map[string][]byte),
		StateRoot: crypto.ZeroHash(),
	}

	s.accounts[addressStr] = account

	// Update the state root
	s.updateStateRoot()

	return account, nil
}

// UpdateAccount updates an account
func (s *State) UpdateAccount(account *Account) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	addressStr := fmt.Sprintf("%x", account.Address)
	if _, ok := s.accounts[addressStr]; !ok {
		return fmt.Errorf("account not found: %s", addressStr)
	}

	s.accounts[addressStr] = account

	// Update the state root
	s.updateStateRoot()

	return nil
}

// AddTransaction adds a transaction to the state
func (s *State) AddTransaction(tx *Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verify the transaction
	valid, err := tx.Verify()
	if err != nil {
		return fmt.Errorf("failed to verify transaction: %w", err)
	}
	if !valid {
		return errors.New("invalid transaction signature")
	}

	// Check if the transaction already exists
	txID := tx.ID.String()
	if _, ok := s.transactions[txID]; ok {
		return fmt.Errorf("transaction already exists: %s", txID)
	}

	// Get the sender account
	senderStr := fmt.Sprintf("%x", tx.From)
	sender, ok := s.accounts[senderStr]
	if !ok {
		// Create a new account for the sender
		var err error
		sender, err = s.CreateAccount(tx.From)
		if err != nil {
			return fmt.Errorf("failed to create sender account: %w", err)
		}
	}

	// Check if the sender has enough balance
	if sender.Balance < tx.Amount {
		return fmt.Errorf("insufficient balance: %d < %d", sender.Balance, tx.Amount)
	}

	// Get the recipient account
	recipientStr := fmt.Sprintf("%x", tx.To)
	recipient, ok := s.accounts[recipientStr]
	if !ok {
		// Create a new account for the recipient
		var err error
		recipient, err = s.CreateAccount(tx.To)
		if err != nil {
			return fmt.Errorf("failed to create recipient account: %w", err)
		}
	}

	// Update the sender's balance and nonce
	sender.Balance -= tx.Amount
	sender.Nonce++

	// Update the recipient's balance
	recipient.Balance += tx.Amount

	// Add the transaction to the state
	s.transactions[txID] = tx

	// Update the state root
	s.updateStateRoot()

	return nil
}

// AddBlock adds a block to the state
func (s *State) AddBlock(block *Block) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the block already exists
	blockHash := block.ComputeHash().String()
	if _, ok := s.blocks[blockHash]; ok {
		return fmt.Errorf("block already exists: %s", blockHash)
	}

	// Check if the previous block exists
	if block.Header.PreviousHash != crypto.ZeroHash() {
		prevHashStr := block.Header.PreviousHash.String()
		if _, ok := s.blocks[prevHashStr]; !ok {
			return fmt.Errorf("previous block not found: %s", prevHashStr)
		}
	}

	// Process the transactions in the block
	for _, tx := range block.Transactions {
		// Skip if the transaction already exists
		txID := tx.ID.String()
		if _, ok := s.transactions[txID]; ok {
			continue
		}

		// Add the transaction to the state
		if err := s.AddTransaction(tx); err != nil {
			return fmt.Errorf("failed to add transaction: %w", err)
		}
	}

	// Add the block to the state
	s.blocks[blockHash] = block

	// Update the latest block
	if s.latestBlock == nil || block.Header.Height > s.latestBlock.Header.Height {
		s.latestBlock = block
	}

	// Update the state root
	s.updateStateRoot()

	// Prune old blocks if necessary
	s.pruneBlocks()

	return nil
}

// GetTransaction gets a transaction by ID
func (s *State) GetTransaction(txID string) (*Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, ok := s.transactions[txID]
	if !ok {
		return nil, fmt.Errorf("transaction not found: %s", txID)
	}

	return tx, nil
}

// GetBlock gets a block by hash
func (s *State) GetBlock(blockHash string) (*Block, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	block, ok := s.blocks[blockHash]
	if !ok {
		return nil, fmt.Errorf("block not found: %s", blockHash)
	}

	return block, nil
}

// GetLatestBlock gets the latest block
func (s *State) GetLatestBlock() *Block {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.latestBlock
}

// GetStateRoot gets the state root
func (s *State) GetStateRoot() crypto.Hash {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.stateRoot
}

// updateStateRoot updates the state root
func (s *State) updateStateRoot() {
	// Compute the state root as the Merkle root of all account state roots
	var accountRoots []crypto.Hash
	for _, account := range s.accounts {
		accountRoots = append(accountRoots, account.StateRoot)
	}

	s.stateRoot = crypto.ComputeMerkleRoot(accountRoots)

	// Update the state accumulator
	s.updateStateAccumulator()
}

// updateStateAccumulator updates the state accumulator
func (s *State) updateStateAccumulator() {
	// Add all accounts to the accumulator
	for _, account := range s.accounts {
		// Serialize the account
		data, err := json.Marshal(account)
		if err != nil {
			continue
		}

		// Add the account to the accumulator
		s.stateAccum.Add(data)
	}
}

// GenerateStateProof generates a proof for an account
func (s *State) GenerateStateProof(address []byte) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	addressStr := fmt.Sprintf("%x", address)
	account, ok := s.accounts[addressStr]
	if !ok {
		return nil, fmt.Errorf("account not found: %s", addressStr)
	}

	// Serialize the account
	data, err := json.Marshal(account)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal account: %w", err)
	}

	// Generate a witness for the account
	// In a real implementation, we would generate a proper witness
	// For now, we'll just return the serialized account
	return data, nil
}

// VerifyStateProof verifies a proof for an account
func (s *State) VerifyStateProof(address []byte, proof []byte) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// In a real implementation, we would verify the proof properly
	// For now, we'll just check if the account exists
	addressStr := fmt.Sprintf("%x", address)
	_, ok := s.accounts[addressStr]

	return ok, nil
}

// pruneBlocks prunes old blocks
func (s *State) pruneBlocks() {
	if s.latestBlock == nil || s.pruningHeight == 0 {
		return
	}

	// Prune blocks older than the pruning height
	for hash, block := range s.blocks {
		if block.Header.Height < s.latestBlock.Header.Height-s.pruningHeight {
			delete(s.blocks, hash)
		}
	}
}

// SetPruningHeight sets the pruning height
func (s *State) SetPruningHeight(height uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pruningHeight = height
}

// NewBlock creates a new block with the given transactions
func (s *State) NewBlock(previousHash crypto.Hash, transactions []*Transaction, shardID uint64) *Block {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a new block
	block := NewBlock(previousHash, transactions, shardID)

	// Set the block height
	if s.latestBlock == nil {
		block.Header.Height = 0
	} else {
		block.Header.Height = s.latestBlock.Header.Height + 1
	}

	// Set the state root
	block.StateRoot = s.stateRoot

	return block
}
