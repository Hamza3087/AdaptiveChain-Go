package network

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/advanced-blockchain/pkg/crypto"
)

// MessageType represents the type of a message
type MessageType string

const (
	// MessageTypeHandshake indicates a handshake message
	MessageTypeHandshake MessageType = "handshake"

	// MessageTypePing indicates a ping message
	MessageTypePing MessageType = "ping"

	// MessageTypePong indicates a pong message
	MessageTypePong MessageType = "pong"

	// MessageTypeTransaction indicates a transaction message
	MessageTypeTransaction MessageType = "transaction"

	// MessageTypeBlock indicates a block message
	MessageTypeBlock MessageType = "block"

	// MessageTypeConsensus indicates a consensus message
	MessageTypeConsensus MessageType = "consensus"

	// MessageTypePeerDiscovery indicates a peer discovery message
	MessageTypePeerDiscovery MessageType = "peer_discovery"
)

// Message represents a network message
type Message struct {
	Type      MessageType // Message type
	Sender    string      // Sender ID
	Recipient string      // Recipient ID
	Timestamp time.Time   // Message timestamp
	Payload   []byte      // Message payload
	Signature []byte      // Message signature
}

// NewMessage creates a new message
func NewMessage(messageType MessageType, sender, recipient string, payload []byte) *Message {
	return &Message{
		Type:      messageType,
		Sender:    sender,
		Recipient: recipient,
		Timestamp: time.Now(),
		Payload:   payload,
	}
}

// Sign signs the message with the given key pair
func (m *Message) Sign(keyPair *crypto.KeyPair) error {
	// Compute the hash of the message
	hash := m.ComputeHash()

	// Sign the hash
	signature, err := keyPair.Sign(hash[:])
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	// Set the signature
	m.Signature = signature

	return nil
}

// Verify verifies the message signature
func (m *Message) Verify(publicKey []byte) (bool, error) {
	// Check if the message has a signature
	if m.Signature == nil {
		return false, fmt.Errorf("message has no signature")
	}

	// Compute the hash of the message
	hash := m.ComputeHash()

	// Verify the signature
	return crypto.VerifySignature(publicKey, hash[:], m.Signature)
}

// ComputeHash computes the hash of the message
func (m *Message) ComputeHash() crypto.Hash {
	// Marshal the message without the signature
	messageCopy := *m
	messageCopy.Signature = nil

	data, err := json.Marshal(messageCopy)
	if err != nil {
		// In case of error, return a zero hash
		return crypto.ZeroHash()
	}

	return crypto.NewHash(data)
}

// Marshal marshals the message to JSON
func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshal unmarshals the message from JSON
func (m *Message) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}

// Unmarshal unmarshals the handshake payload from JSON
func (p *HandshakePayload) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

// Unmarshal unmarshals the ping payload from JSON
func (p *PingPayload) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

// Unmarshal unmarshals the pong payload from JSON
func (p *PongPayload) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

// Unmarshal unmarshals the transaction payload from JSON
func (p *TransactionPayload) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

// Unmarshal unmarshals the block payload from JSON
func (p *BlockPayload) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

// Unmarshal unmarshals the consensus payload from JSON
func (p *ConsensusPayload) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

// Unmarshal unmarshals the peer discovery payload from JSON
func (p *PeerDiscoveryPayload) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

// HandshakePayload represents the payload of a handshake message
type HandshakePayload struct {
	Version   string    // Protocol version
	NodeID    string    // Node ID
	PublicKey []byte    // Node public key
	Timestamp time.Time // Handshake timestamp
	Nonce     uint64    // Nonce for challenge-response
}

// PingPayload represents the payload of a ping message
type PingPayload struct {
	Timestamp time.Time // Ping timestamp
	Nonce     uint64    // Nonce for challenge-response
}

// PongPayload represents the payload of a pong message
type PongPayload struct {
	Timestamp time.Time // Pong timestamp
	Nonce     uint64    // Nonce from the ping message
}

// TransactionPayload represents the payload of a transaction message
type TransactionPayload struct {
	TransactionID string // Transaction ID
	Transaction   []byte // Serialized transaction
}

// BlockPayload represents the payload of a block message
type BlockPayload struct {
	BlockHash string // Block hash
	Block     []byte // Serialized block
}

// ConsensusPayload represents the payload of a consensus message
type ConsensusPayload struct {
	Round uint64 // Consensus round
	Type  string // Consensus message type
	Data  []byte // Consensus message data
}

// PeerDiscoveryPayload represents the payload of a peer discovery message
type PeerDiscoveryPayload struct {
	Peers []PeerInfo // List of peers
}

// PeerInfo represents information about a peer
type PeerInfo struct {
	ID        string // Peer ID
	Address   string // Peer address
	PublicKey []byte // Peer public key
}

// NewHandshakeMessage creates a new handshake message
func NewHandshakeMessage(sender, recipient, nodeID string, publicKey []byte, version string) (*Message, error) {
	// Create the handshake payload
	payload := HandshakePayload{
		Version:   version,
		NodeID:    nodeID,
		PublicKey: publicKey,
		Timestamp: time.Now(),
		Nonce:     uint64(time.Now().UnixNano()),
	}

	// Marshal the payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal handshake payload: %w", err)
	}

	// Create the message
	return NewMessage(MessageTypeHandshake, sender, recipient, payloadBytes), nil
}

// NewPingMessage creates a new ping message
func NewPingMessage(sender, recipient string) (*Message, error) {
	// Create the ping payload
	payload := PingPayload{
		Timestamp: time.Now(),
		Nonce:     uint64(time.Now().UnixNano()),
	}

	// Marshal the payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ping payload: %w", err)
	}

	// Create the message
	return NewMessage(MessageTypePing, sender, recipient, payloadBytes), nil
}

// NewPongMessage creates a new pong message
func NewPongMessage(sender, recipient string, pingNonce uint64) (*Message, error) {
	// Create the pong payload
	payload := PongPayload{
		Timestamp: time.Now(),
		Nonce:     pingNonce,
	}

	// Marshal the payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pong payload: %w", err)
	}

	// Create the message
	return NewMessage(MessageTypePong, sender, recipient, payloadBytes), nil
}

// NewTransactionMessage creates a new transaction message
func NewTransactionMessage(sender, recipient, transactionID string, transaction []byte) (*Message, error) {
	// Create the transaction payload
	payload := TransactionPayload{
		TransactionID: transactionID,
		Transaction:   transaction,
	}

	// Marshal the payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction payload: %w", err)
	}

	// Create the message
	return NewMessage(MessageTypeTransaction, sender, recipient, payloadBytes), nil
}

// NewBlockMessage creates a new block message
func NewBlockMessage(sender, recipient, blockHash string, block []byte) (*Message, error) {
	// Create the block payload
	payload := BlockPayload{
		BlockHash: blockHash,
		Block:     block,
	}

	// Marshal the payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal block payload: %w", err)
	}

	// Create the message
	return NewMessage(MessageTypeBlock, sender, recipient, payloadBytes), nil
}

// NewConsensusMessage creates a new consensus message
func NewConsensusMessage(sender, recipient string, round uint64, consensusType string, data []byte) (*Message, error) {
	// Create the consensus payload
	payload := ConsensusPayload{
		Round: round,
		Type:  consensusType,
		Data:  data,
	}

	// Marshal the payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal consensus payload: %w", err)
	}

	// Create the message
	return NewMessage(MessageTypeConsensus, sender, recipient, payloadBytes), nil
}

// NewPeerDiscoveryMessage creates a new peer discovery message
func NewPeerDiscoveryMessage(sender, recipient string, peers []PeerInfo) (*Message, error) {
	// Create the peer discovery payload
	payload := PeerDiscoveryPayload{
		Peers: peers,
	}

	// Marshal the payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal peer discovery payload: %w", err)
	}

	// Create the message
	return NewMessage(MessageTypePeerDiscovery, sender, recipient, payloadBytes), nil
}
