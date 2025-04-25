package state

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync"
	"time"

	"github.com/advanced-blockchain/pkg/crypto"
)

// CompressionAlgorithm represents a compression algorithm
type CompressionAlgorithm string

const (
	// CompressionAlgorithmGzip indicates Gzip compression
	CompressionAlgorithmGzip CompressionAlgorithm = "gzip"

	// CompressionAlgorithmSnappy indicates Snappy compression
	CompressionAlgorithmSnappy CompressionAlgorithm = "snappy"

	// CompressionAlgorithmLZ4 indicates LZ4 compression
	CompressionAlgorithmLZ4 CompressionAlgorithm = "lz4"
)

// ArchiveStatus represents the status of an archive
type ArchiveStatus string

const (
	// ArchiveStatusCreating indicates a creating archive
	ArchiveStatusCreating ArchiveStatus = "creating"

	// ArchiveStatusReady indicates a ready archive
	ArchiveStatusReady ArchiveStatus = "ready"

	// ArchiveStatusFailed indicates a failed archive
	ArchiveStatusFailed ArchiveStatus = "failed"
)

// StateArchive represents an archived state
type StateArchive struct {
	ID             string               // Archive ID
	Height         uint64               // Block height
	Timestamp      time.Time            // Archive timestamp
	StateRoot      crypto.Hash          // State root
	Algorithm      CompressionAlgorithm // Compression algorithm
	CompressedData []byte               // Compressed state data
	Status         ArchiveStatus        // Archive status
	Error          string               // Error message, if any
}

// StateCompressor represents a state compressor
type StateCompressor struct {
	state         *State                   // Reference to the state
	archives      map[string]*StateArchive // Map of archive ID to archive
	algorithm     CompressionAlgorithm     // Compression algorithm
	pruningHeight uint64                   // Pruning height
	mu            sync.RWMutex             // Mutex for thread safety
}

// NewStateCompressor creates a new state compressor
func NewStateCompressor(state *State, algorithm CompressionAlgorithm) *StateCompressor {
	return &StateCompressor{
		state:         state,
		archives:      make(map[string]*StateArchive),
		algorithm:     algorithm,
		pruningHeight: 100, // Default pruning height
	}
}

// SetPruningHeight sets the pruning height
func (sc *StateCompressor) SetPruningHeight(height uint64) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.pruningHeight = height
	sc.state.SetPruningHeight(height)
}

// CreateArchive creates a new state archive
func (sc *StateCompressor) CreateArchive() (*StateArchive, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Get the latest block
	latestBlock := sc.state.GetLatestBlock()
	if latestBlock == nil {
		return nil, errors.New("no blocks in the state")
	}

	// Create a new archive
	archiveID := fmt.Sprintf("archive-%d", time.Now().UnixNano())
	archive := &StateArchive{
		ID:        archiveID,
		Height:    latestBlock.Header.Height,
		Timestamp: time.Now(),
		StateRoot: sc.state.GetStateRoot(),
		Algorithm: sc.algorithm,
		Status:    ArchiveStatusCreating,
	}

	// Add the archive to the map
	sc.archives[archiveID] = archive

	// Compress the state in a separate goroutine
	go sc.compressState(archiveID)

	return archive, nil
}

// compressState compresses the state for an archive
func (sc *StateCompressor) compressState(archiveID string) {
	sc.mu.Lock()
	archive := sc.archives[archiveID]
	sc.mu.Unlock()

	// Serialize the state
	stateData, err := sc.serializeState()
	if err != nil {
		sc.setArchiveError(archiveID, fmt.Sprintf("failed to serialize state: %v", err))
		return
	}

	// Compress the state data
	compressedData, err := sc.compressData(stateData)
	if err != nil {
		sc.setArchiveError(archiveID, fmt.Sprintf("failed to compress state data: %v", err))
		return
	}

	// Update the archive
	sc.mu.Lock()
	archive.CompressedData = compressedData
	archive.Status = ArchiveStatusReady
	sc.mu.Unlock()
}

// serializeState serializes the state
func (sc *StateCompressor) serializeState() ([]byte, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	// Create a serializable state
	type SerializableState struct {
		Accounts     map[string]*Account
		Transactions map[string]*Transaction
		Blocks       map[string]*Block
		StateRoot    crypto.Hash
	}

	serializableState := SerializableState{
		Accounts:     sc.state.accounts,
		Transactions: sc.state.transactions,
		Blocks:       sc.state.blocks,
		StateRoot:    sc.state.stateRoot,
	}

	// Serialize the state
	return json.Marshal(serializableState)
}

// compressData compresses data using the configured algorithm
func (sc *StateCompressor) compressData(data []byte) ([]byte, error) {
	// Compress the data using the configured algorithm
	switch sc.algorithm {
	case CompressionAlgorithmGzip:
		return sc.compressGzip(data)
	case CompressionAlgorithmSnappy:
		// In a real implementation, we would use the Snappy algorithm
		// For now, we'll just use Gzip
		return sc.compressGzip(data)
	case CompressionAlgorithmLZ4:
		// In a real implementation, we would use the LZ4 algorithm
		// For now, we'll just use Gzip
		return sc.compressGzip(data)
	default:
		return nil, fmt.Errorf("unsupported compression algorithm: %s", sc.algorithm)
	}
}

// compressGzip compresses data using Gzip
func (sc *StateCompressor) compressGzip(data []byte) ([]byte, error) {
	// Create a buffer for the compressed data
	var buf bytes.Buffer

	// Create a Gzip writer
	gzipWriter := gzip.NewWriter(&buf)

	// Write the data to the Gzip writer
	if _, err := gzipWriter.Write(data); err != nil {
		return nil, err
	}

	// Close the Gzip writer
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// decompressData decompresses data using the configured algorithm
func (sc *StateCompressor) decompressData(data []byte, algorithm CompressionAlgorithm) ([]byte, error) {
	// Decompress the data using the specified algorithm
	switch algorithm {
	case CompressionAlgorithmGzip:
		return sc.decompressGzip(data)
	case CompressionAlgorithmSnappy:
		// In a real implementation, we would use the Snappy algorithm
		// For now, we'll just use Gzip
		return sc.decompressGzip(data)
	case CompressionAlgorithmLZ4:
		// In a real implementation, we would use the LZ4 algorithm
		// For now, we'll just use Gzip
		return sc.decompressGzip(data)
	default:
		return nil, fmt.Errorf("unsupported compression algorithm: %s", algorithm)
	}
}

// decompressGzip decompresses data using Gzip
func (sc *StateCompressor) decompressGzip(data []byte) ([]byte, error) {
	// Create a reader for the compressed data
	reader := bytes.NewReader(data)

	// Create a Gzip reader
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	// Read the decompressed data
	decompressed, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}

	// Close the Gzip reader
	if err := gzipReader.Close(); err != nil {
		return nil, err
	}

	return decompressed, nil
}

// setArchiveError sets an error for an archive
func (sc *StateCompressor) setArchiveError(archiveID, errorMsg string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	archive := sc.archives[archiveID]
	archive.Status = ArchiveStatusFailed
	archive.Error = errorMsg
}

// GetArchive gets an archive by ID
func (sc *StateCompressor) GetArchive(archiveID string) (*StateArchive, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	archive, ok := sc.archives[archiveID]
	if !ok {
		return nil, fmt.Errorf("archive not found: %s", archiveID)
	}

	return archive, nil
}

// RestoreFromArchive restores the state from an archive
func (sc *StateCompressor) RestoreFromArchive(archiveID string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Get the archive
	archive, ok := sc.archives[archiveID]
	if !ok {
		return fmt.Errorf("archive not found: %s", archiveID)
	}

	// Check if the archive is ready
	if archive.Status != ArchiveStatusReady {
		return fmt.Errorf("archive not ready: %s", archive.Status)
	}

	// Decompress the state data
	stateData, err := sc.decompressData(archive.CompressedData, archive.Algorithm)
	if err != nil {
		return fmt.Errorf("failed to decompress state data: %w", err)
	}

	// Deserialize the state
	if err := sc.deserializeState(stateData); err != nil {
		return fmt.Errorf("failed to deserialize state: %w", err)
	}

	return nil
}

// deserializeState deserializes the state
func (sc *StateCompressor) deserializeState(data []byte) error {
	// Create a serializable state
	type SerializableState struct {
		Accounts     map[string]*Account
		Transactions map[string]*Transaction
		Blocks       map[string]*Block
		StateRoot    crypto.Hash
	}

	// Deserialize the state
	var serializableState SerializableState
	if err := json.Unmarshal(data, &serializableState); err != nil {
		return err
	}

	// Update the state
	sc.state.accounts = serializableState.Accounts
	sc.state.transactions = serializableState.Transactions
	sc.state.blocks = serializableState.Blocks
	sc.state.stateRoot = serializableState.StateRoot

	// Find the latest block
	var latestBlock *Block
	for _, block := range sc.state.blocks {
		if latestBlock == nil || block.Header.Height > latestBlock.Header.Height {
			latestBlock = block
		}
	}
	sc.state.latestBlock = latestBlock

	return nil
}

// PruneState prunes the state
func (sc *StateCompressor) PruneState() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Get the latest block
	latestBlock := sc.state.latestBlock
	if latestBlock == nil {
		return errors.New("no blocks in the state")
	}

	// Compute the pruning height
	pruningHeight := latestBlock.Header.Height - sc.pruningHeight
	if pruningHeight < 0 {
		pruningHeight = 0
	}

	// Prune blocks older than the pruning height
	for hash, block := range sc.state.blocks {
		if block.Header.Height < pruningHeight {
			delete(sc.state.blocks, hash)
		}
	}

	// Prune transactions that are not in any remaining block
	transactionsInBlocks := make(map[string]bool)
	for _, block := range sc.state.blocks {
		for _, tx := range block.Transactions {
			transactionsInBlocks[tx.ID.String()] = true
		}
	}
	for txID := range sc.state.transactions {
		if !transactionsInBlocks[txID] {
			delete(sc.state.transactions, txID)
		}
	}

	return nil
}

// CompactState compacts the state
func (sc *StateCompressor) CompactState() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Create a new state
	newState := NewState()

	// Copy the accounts
	for addressStr, account := range sc.state.accounts {
		newState.accounts[addressStr] = account
	}

	// Copy only the transactions that are in blocks
	transactionsInBlocks := make(map[string]bool)
	for _, block := range sc.state.blocks {
		for _, tx := range block.Transactions {
			txID := tx.ID.String()
			transactionsInBlocks[txID] = true
			newState.transactions[txID] = tx
		}
	}

	// Copy the blocks
	for hash, block := range sc.state.blocks {
		newState.blocks[hash] = block
	}

	// Set the latest block
	newState.latestBlock = sc.state.latestBlock

	// Set the state root
	newState.stateRoot = sc.state.stateRoot

	// Update the state
	sc.state.accounts = newState.accounts
	sc.state.transactions = newState.transactions
	sc.state.blocks = newState.blocks
	sc.state.latestBlock = newState.latestBlock
	sc.state.stateRoot = newState.stateRoot

	return nil
}

// StateAccumulator represents a state accumulator
type StateAccumulator struct {
	accumulator *crypto.Accumulator // Cryptographic accumulator
	values      map[string][]byte   // Map of key to value
	mu          sync.RWMutex        // Mutex for thread safety
}

// NewStateAccumulator creates a new state accumulator
func NewStateAccumulator() *StateAccumulator {
	// Create a new RSA modulus for the accumulator
	// In a real implementation, this would be a secure RSA modulus
	n := crypto.NewHash([]byte("state-accumulator"))
	nBig := new(big.Int).SetBytes(n.Bytes())

	return &StateAccumulator{
		accumulator: crypto.NewAccumulator(nBig),
		values:      make(map[string][]byte),
	}
}

// Add adds a key-value pair to the accumulator
func (sa *StateAccumulator) Add(key string, value []byte) {
	sa.mu.Lock()
	defer sa.mu.Unlock()

	// Add the key-value pair to the map
	sa.values[key] = value

	// Add the key-value pair to the accumulator
	data := make([]byte, 0, len(key)+len(value))
	data = append(data, []byte(key)...)
	data = append(data, value...)
	sa.accumulator.Add(data)
}

// Remove removes a key-value pair from the accumulator
func (sa *StateAccumulator) Remove(key string) {
	sa.mu.Lock()
	defer sa.mu.Unlock()

	// Remove the key-value pair from the map
	delete(sa.values, key)

	// Rebuild the accumulator
	sa.rebuildAccumulator()
}

// rebuildAccumulator rebuilds the accumulator
func (sa *StateAccumulator) rebuildAccumulator() {
	// Create a new RSA modulus for the accumulator
	// In a real implementation, this would be a secure RSA modulus
	n := crypto.NewHash([]byte("state-accumulator"))
	nBig := new(big.Int).SetBytes(n.Bytes())

	// Create a new accumulator
	sa.accumulator = crypto.NewAccumulator(nBig)

	// Add all key-value pairs to the accumulator
	for key, value := range sa.values {
		data := make([]byte, 0, len(key)+len(value))
		data = append(data, []byte(key)...)
		data = append(data, value...)
		sa.accumulator.Add(data)
	}
}

// GenerateProof generates a proof for a key-value pair
func (sa *StateAccumulator) GenerateProof(key string) ([]byte, error) {
	sa.mu.RLock()
	defer sa.mu.RUnlock()

	// Check if the key exists
	value, ok := sa.values[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	// Create a proof structure
	type Proof struct {
		Key   string
		Value []byte
	}

	// Create a proof
	proof := Proof{
		Key:   key,
		Value: value,
	}

	// Serialize the proof
	proofBytes, err := json.Marshal(proof)
	if err != nil {
		return nil, err
	}

	return proofBytes, nil
}

// VerifyProof verifies a proof for a key-value pair
func (sa *StateAccumulator) VerifyProof(proof []byte) (bool, error) {
	sa.mu.RLock()
	defer sa.mu.RUnlock()

	// Deserialize the proof
	var p struct {
		Key   string
		Value []byte
	}
	if err := json.Unmarshal(proof, &p); err != nil {
		return false, err
	}

	// Check if the key exists
	value, ok := sa.values[p.Key]
	if !ok {
		return false, nil
	}

	// Check if the value matches
	return bytes.Equal(value, p.Value), nil
}

// GetValue gets a value from the accumulator
func (sa *StateAccumulator) GetValue(key string) ([]byte, error) {
	sa.mu.RLock()
	defer sa.mu.RUnlock()

	// Check if the key exists
	value, ok := sa.values[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	return value, nil
}

// GetAccumulatorValue gets the accumulator value
func (sa *StateAccumulator) GetAccumulatorValue() *big.Int {
	sa.mu.RLock()
	defer sa.mu.RUnlock()

	return sa.accumulator.GetValue()
}

// CompactStateRepresentation represents a compact state representation
type CompactStateRepresentation struct {
	Height       uint64      // Block height
	StateRoot    crypto.Hash // State root
	Accumulator  *big.Int    // Accumulator value
	AccountCount uint64      // Number of accounts
	BlockCount   uint64      // Number of blocks
	TxCount      uint64      // Number of transactions
	Timestamp    time.Time   // Timestamp
}

// CreateCompactStateRepresentation creates a compact state representation
func (sc *StateCompressor) CreateCompactStateRepresentation() (*CompactStateRepresentation, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	// Get the latest block
	latestBlock := sc.state.latestBlock
	if latestBlock == nil {
		return nil, errors.New("no blocks in the state")
	}

	// Create a state accumulator
	accumulator := NewStateAccumulator()

	// Add all accounts to the accumulator
	for addressStr, account := range sc.state.accounts {
		// Serialize the account
		accountBytes, err := json.Marshal(account)
		if err != nil {
			return nil, err
		}

		// Add the account to the accumulator
		accumulator.Add(addressStr, accountBytes)
	}

	// Create a compact state representation
	compact := &CompactStateRepresentation{
		Height:       latestBlock.Header.Height,
		StateRoot:    sc.state.stateRoot,
		Accumulator:  accumulator.GetAccumulatorValue(),
		AccountCount: uint64(len(sc.state.accounts)),
		BlockCount:   uint64(len(sc.state.blocks)),
		TxCount:      uint64(len(sc.state.transactions)),
		Timestamp:    time.Now(),
	}

	return compact, nil
}

// SerializeCompactStateRepresentation serializes a compact state representation
func (sc *StateCompressor) SerializeCompactStateRepresentation(compact *CompactStateRepresentation) ([]byte, error) {
	// Serialize the compact state representation
	return json.Marshal(compact)
}

// DeserializeCompactStateRepresentation deserializes a compact state representation
func (sc *StateCompressor) DeserializeCompactStateRepresentation(data []byte) (*CompactStateRepresentation, error) {
	// Deserialize the compact state representation
	var compact CompactStateRepresentation
	if err := json.Unmarshal(data, &compact); err != nil {
		return nil, err
	}

	return &compact, nil
}
