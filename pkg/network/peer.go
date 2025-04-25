package network

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/advanced-blockchain/pkg/crypto"
)

// PeerStatus represents the status of a peer
type PeerStatus string

const (
	// PeerStatusConnecting indicates a connecting peer
	PeerStatusConnecting PeerStatus = "connecting"
	
	// PeerStatusConnected indicates a connected peer
	PeerStatusConnected PeerStatus = "connected"
	
	// PeerStatusDisconnected indicates a disconnected peer
	PeerStatusDisconnected PeerStatus = "disconnected"
	
	// PeerStatusBanned indicates a banned peer
	PeerStatusBanned PeerStatus = "banned"
)

// Peer represents a network peer
type Peer struct {
	ID            string          // Peer ID
	Address       string          // Peer address
	PublicKey     []byte          // Peer public key
	Status        PeerStatus      // Peer status
	LastSeen      time.Time       // Last seen time
	Latency       time.Duration   // Network latency
	Conn          net.Conn        // Network connection
	mu            sync.RWMutex    // Mutex for thread safety
}

// NewPeer creates a new peer
func NewPeer(id, address string, publicKey []byte) *Peer {
	return &Peer{
		ID:        id,
		Address:   address,
		PublicKey: publicKey,
		Status:    PeerStatusDisconnected,
		LastSeen:  time.Now(),
	}
}

// Connect connects to the peer
func (p *Peer) Connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// Check if the peer is already connected
	if p.Status == PeerStatusConnected {
		return fmt.Errorf("peer already connected: %s", p.ID)
	}
	
	// Check if the peer is banned
	if p.Status == PeerStatusBanned {
		return fmt.Errorf("peer is banned: %s", p.ID)
	}
	
	// Set the status to connecting
	p.Status = PeerStatusConnecting
	
	// Connect to the peer
	conn, err := net.Dial("tcp", p.Address)
	if err != nil {
		p.Status = PeerStatusDisconnected
		return fmt.Errorf("failed to connect to peer: %w", err)
	}
	
	// Set the connection
	p.Conn = conn
	
	// Set the status to connected
	p.Status = PeerStatusConnected
	
	// Update the last seen time
	p.LastSeen = time.Now()
	
	return nil
}

// Disconnect disconnects from the peer
func (p *Peer) Disconnect() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// Check if the peer is connected
	if p.Status != PeerStatusConnected {
		return fmt.Errorf("peer not connected: %s", p.ID)
	}
	
	// Close the connection
	if p.Conn != nil {
		if err := p.Conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
		p.Conn = nil
	}
	
	// Set the status to disconnected
	p.Status = PeerStatusDisconnected
	
	return nil
}

// Ban bans the peer
func (p *Peer) Ban() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// Disconnect if connected
	if p.Status == PeerStatusConnected {
		if p.Conn != nil {
			if err := p.Conn.Close(); err != nil {
				return fmt.Errorf("failed to close connection: %w", err)
			}
			p.Conn = nil
		}
	}
	
	// Set the status to banned
	p.Status = PeerStatusBanned
	
	return nil
}

// Unban unbans the peer
func (p *Peer) Unban() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// Set the status to disconnected
	if p.Status == PeerStatusBanned {
		p.Status = PeerStatusDisconnected
	}
}

// Send sends data to the peer
func (p *Peer) Send(data []byte) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	// Check if the peer is connected
	if p.Status != PeerStatusConnected {
		return fmt.Errorf("peer not connected: %s", p.ID)
	}
	
	// Check if the connection is valid
	if p.Conn == nil {
		return fmt.Errorf("peer connection is nil: %s", p.ID)
	}
	
	// Send the data
	_, err := p.Conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}
	
	return nil
}

// Receive receives data from the peer
func (p *Peer) Receive(buffer []byte) (int, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	// Check if the peer is connected
	if p.Status != PeerStatusConnected {
		return 0, fmt.Errorf("peer not connected: %s", p.ID)
	}
	
	// Check if the connection is valid
	if p.Conn == nil {
		return 0, fmt.Errorf("peer connection is nil: %s", p.ID)
	}
	
	// Receive the data
	n, err := p.Conn.Read(buffer)
	if err != nil {
		return 0, fmt.Errorf("failed to receive data: %w", err)
	}
	
	// Update the last seen time
	p.LastSeen = time.Now()
	
	return n, nil
}

// UpdateLatency updates the peer's latency
func (p *Peer) UpdateLatency(latency time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.Latency = latency
}

// GetStatus gets the peer's status
func (p *Peer) GetStatus() PeerStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	return p.Status
}

// GetLatency gets the peer's latency
func (p *Peer) GetLatency() time.Duration {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	return p.Latency
}

// GetLastSeen gets the peer's last seen time
func (p *Peer) GetLastSeen() time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	return p.LastSeen
}

// VerifySignature verifies a signature from the peer
func (p *Peer) VerifySignature(message, signature []byte) (bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	// Verify the signature
	return crypto.VerifySignature(p.PublicKey, message, signature)
}
