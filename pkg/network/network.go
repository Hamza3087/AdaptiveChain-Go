package network

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/advanced-blockchain/pkg/crypto"
)

// Network represents a network manager
type Network struct {
	NodeID          string                         // Node ID
	Address         string                         // Node address
	KeyPair         *crypto.KeyPair                // Node key pair
	Peers           map[string]*Peer               // Map of peer ID to peer
	MessageHandlers map[MessageType]MessageHandler // Map of message type to handler
	Listener        net.Listener                   // Network listener
	Running         bool                           // Whether the network is running
	mu              sync.RWMutex                   // Mutex for thread safety
}

// MessageHandler represents a message handler function
type MessageHandler func(*Message, *Peer) error

// NewNetwork creates a new network manager
func NewNetwork(nodeID, address string, keyPair *crypto.KeyPair) *Network {
	return &Network{
		NodeID:          nodeID,
		Address:         address,
		KeyPair:         keyPair,
		Peers:           make(map[string]*Peer),
		MessageHandlers: make(map[MessageType]MessageHandler),
		Running:         false,
	}
}

// Start starts the network manager
func (n *Network) Start() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Check if the network is already running
	if n.Running {
		return fmt.Errorf("network already running")
	}

	// Start listening for incoming connections
	listener, err := net.Listen("tcp", n.Address)
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	n.Listener = listener
	n.Running = true

	// Start accepting connections in a separate goroutine
	go n.acceptConnections()

	return nil
}

// Stop stops the network manager
func (n *Network) Stop() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Check if the network is running
	if !n.Running {
		return fmt.Errorf("network not running")
	}

	// Close the listener
	if n.Listener != nil {
		if err := n.Listener.Close(); err != nil {
			return fmt.Errorf("failed to close listener: %w", err)
		}
		n.Listener = nil
	}

	// Disconnect from all peers
	for _, peer := range n.Peers {
		if peer.GetStatus() == PeerStatusConnected {
			if err := peer.Disconnect(); err != nil {
				// Log the error but continue
				fmt.Printf("Error disconnecting from peer %s: %v\n", peer.ID, err)
			}
		}
	}

	n.Running = false

	return nil
}

// AddPeer adds a peer to the network
func (n *Network) AddPeer(id, address string, publicKey []byte) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Check if the peer already exists
	if _, ok := n.Peers[id]; ok {
		return fmt.Errorf("peer already exists: %s", id)
	}

	// Create a new peer
	peer := NewPeer(id, address, publicKey)

	// Add the peer to the map
	n.Peers[id] = peer

	return nil
}

// RemovePeer removes a peer from the network
func (n *Network) RemovePeer(id string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Check if the peer exists
	peer, ok := n.Peers[id]
	if !ok {
		return fmt.Errorf("peer not found: %s", id)
	}

	// Disconnect from the peer if connected
	if peer.GetStatus() == PeerStatusConnected {
		if err := peer.Disconnect(); err != nil {
			return fmt.Errorf("failed to disconnect from peer: %w", err)
		}
	}

	// Remove the peer from the map
	delete(n.Peers, id)

	return nil
}

// ConnectToPeer connects to a peer
func (n *Network) ConnectToPeer(id string) error {
	n.mu.RLock()
	defer n.mu.RUnlock()

	// Check if the peer exists
	peer, ok := n.Peers[id]
	if !ok {
		return fmt.Errorf("peer not found: %s", id)
	}

	// Connect to the peer
	if err := peer.Connect(); err != nil {
		return fmt.Errorf("failed to connect to peer: %w", err)
	}

	// Send a handshake message
	if err := n.sendHandshake(peer); err != nil {
		return fmt.Errorf("failed to send handshake: %w", err)
	}

	// Start receiving messages from the peer in a separate goroutine
	go n.receiveMessages(peer)

	return nil
}

// DisconnectFromPeer disconnects from a peer
func (n *Network) DisconnectFromPeer(id string) error {
	n.mu.RLock()
	defer n.mu.RUnlock()

	// Check if the peer exists
	peer, ok := n.Peers[id]
	if !ok {
		return fmt.Errorf("peer not found: %s", id)
	}

	// Disconnect from the peer
	if err := peer.Disconnect(); err != nil {
		return fmt.Errorf("failed to disconnect from peer: %w", err)
	}

	return nil
}

// BroadcastMessage broadcasts a message to all connected peers
func (n *Network) BroadcastMessage(message *Message) error {
	n.mu.RLock()
	defer n.mu.RUnlock()

	// Sign the message
	if err := message.Sign(n.KeyPair); err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	// Marshal the message
	data, err := message.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Send the message to all connected peers
	for _, peer := range n.Peers {
		if peer.GetStatus() == PeerStatusConnected {
			if err := peer.Send(data); err != nil {
				// Log the error but continue
				fmt.Printf("Error sending message to peer %s: %v\n", peer.ID, err)
			}
		}
	}

	return nil
}

// SendMessage sends a message to a specific peer
func (n *Network) SendMessage(peerID string, message *Message) error {
	n.mu.RLock()
	defer n.mu.RUnlock()

	// Check if the peer exists
	peer, ok := n.Peers[peerID]
	if !ok {
		return fmt.Errorf("peer not found: %s", peerID)
	}

	// Check if the peer is connected
	if peer.GetStatus() != PeerStatusConnected {
		return fmt.Errorf("peer not connected: %s", peerID)
	}

	// Sign the message
	if err := message.Sign(n.KeyPair); err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	// Marshal the message
	data, err := message.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Send the message to the peer
	if err := peer.Send(data); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// RegisterMessageHandler registers a message handler
func (n *Network) RegisterMessageHandler(messageType MessageType, handler MessageHandler) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.MessageHandlers[messageType] = handler
}

// acceptConnections accepts incoming connections
func (n *Network) acceptConnections() {
	for {
		// Accept a connection
		conn, err := n.Listener.Accept()
		if err != nil {
			// Check if the network is still running
			n.mu.RLock()
			running := n.Running
			n.mu.RUnlock()

			if !running {
				// The network has been stopped, so exit the loop
				break
			}

			// Log the error and continue
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go n.handleConnection(conn)
	}
}

// handleConnection handles an incoming connection
func (n *Network) handleConnection(conn net.Conn) {
	// Receive the handshake message
	buffer := make([]byte, 4096)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error receiving handshake: %v\n", err)
		conn.Close()
		return
	}

	// Unmarshal the message
	var message Message
	if err := message.Unmarshal(buffer[:bytesRead]); err != nil {
		fmt.Printf("Error unmarshaling handshake: %v\n", err)
		conn.Close()
		return
	}

	// Check if the message is a handshake
	if message.Type != MessageTypeHandshake {
		fmt.Printf("Expected handshake message, got %s\n", message.Type)
		conn.Close()
		return
	}

	// Unmarshal the handshake payload
	var payload HandshakePayload
	if err := payload.Unmarshal(message.Payload); err != nil {
		fmt.Printf("Error unmarshaling handshake payload: %v\n", err)
		conn.Close()
		return
	}

	// Check if the peer exists
	n.mu.Lock()
	peer, ok := n.Peers[payload.NodeID]
	if !ok {
		// Create a new peer
		peer = NewPeer(payload.NodeID, conn.RemoteAddr().String(), payload.PublicKey)
		n.Peers[payload.NodeID] = peer
	}

	// Set the connection
	peer.Conn = conn
	peer.Status = PeerStatusConnected
	peer.LastSeen = time.Now()
	n.mu.Unlock()

	// Send a handshake response
	if err := n.sendHandshake(peer); err != nil {
		fmt.Printf("Error sending handshake response: %v\n", err)
		conn.Close()
		return
	}

	// Start receiving messages from the peer
	go n.receiveMessages(peer)
}

// sendHandshake sends a handshake message to a peer
func (n *Network) sendHandshake(peer *Peer) error {
	// Export the public key
	publicKey, err := n.KeyPair.ExportPublicKey()
	if err != nil {
		return fmt.Errorf("failed to export public key: %w", err)
	}

	// Create a handshake message
	message, err := NewHandshakeMessage(n.NodeID, peer.ID, n.NodeID, publicKey, "1.0")
	if err != nil {
		return fmt.Errorf("failed to create handshake message: %w", err)
	}

	// Sign the message
	if err := message.Sign(n.KeyPair); err != nil {
		return fmt.Errorf("failed to sign handshake message: %w", err)
	}

	// Marshal the message
	data, err := message.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal handshake message: %w", err)
	}

	// Send the message to the peer
	if err := peer.Send(data); err != nil {
		return fmt.Errorf("failed to send handshake message: %w", err)
	}

	return nil
}

// receiveMessages receives messages from a peer
func (n *Network) receiveMessages(peer *Peer) {
	buffer := make([]byte, 4096)

	for {
		// Receive a message
		bytesRead, err := peer.Receive(buffer)
		if err != nil {
			// Check if the peer is still connected
			if peer.GetStatus() != PeerStatusConnected {
				// The peer has been disconnected, so exit the loop
				break
			}

			// Log the error and disconnect from the peer
			fmt.Printf("Error receiving message from peer %s: %v\n", peer.ID, err)
			peer.Disconnect()
			break
		}

		// Unmarshal the message
		var message Message
		if err := message.Unmarshal(buffer[:bytesRead]); err != nil {
			fmt.Printf("Error unmarshaling message from peer %s: %v\n", peer.ID, err)
			continue
		}

		// Verify the message signature
		valid, err := message.Verify(peer.PublicKey)
		if err != nil {
			fmt.Printf("Error verifying message signature from peer %s: %v\n", peer.ID, err)
			continue
		}
		if !valid {
			fmt.Printf("Invalid message signature from peer %s\n", peer.ID)
			continue
		}

		// Handle the message
		n.handleMessage(&message, peer)
	}
}

// handleMessage handles a received message
func (n *Network) handleMessage(message *Message, peer *Peer) {
	// Get the message handler
	n.mu.RLock()
	handler, ok := n.MessageHandlers[message.Type]
	n.mu.RUnlock()

	if !ok {
		// No handler for this message type
		fmt.Printf("No handler for message type %s\n", message.Type)
		return
	}

	// Handle the message
	if err := handler(message, peer); err != nil {
		fmt.Printf("Error handling message from peer %s: %v\n", peer.ID, err)
	}
}

// GetPeer gets a peer by ID
func (n *Network) GetPeer(id string) (*Peer, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	// Check if the peer exists
	peer, ok := n.Peers[id]
	if !ok {
		return nil, fmt.Errorf("peer not found: %s", id)
	}

	return peer, nil
}

// GetConnectedPeers gets all connected peers
func (n *Network) GetConnectedPeers() []*Peer {
	n.mu.RLock()
	defer n.mu.RUnlock()

	var connectedPeers []*Peer
	for _, peer := range n.Peers {
		if peer.GetStatus() == PeerStatusConnected {
			connectedPeers = append(connectedPeers, peer)
		}
	}

	return connectedPeers
}

// GetPeerCount gets the number of peers
func (n *Network) GetPeerCount() int {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return len(n.Peers)
}

// GetConnectedPeerCount gets the number of connected peers
func (n *Network) GetConnectedPeerCount() int {
	n.mu.RLock()
	defer n.mu.RUnlock()

	count := 0
	for _, peer := range n.Peers {
		if peer.GetStatus() == PeerStatusConnected {
			count++
		}
	}

	return count
}

// IsRunning checks if the network is running
func (n *Network) IsRunning() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.Running
}
