package node

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/advanced-blockchain/pkg/crypto"
)

// AuthStatus represents the authentication status of a node
type AuthStatus string

const (
	// AuthStatusUnauthenticated indicates an unauthenticated node
	AuthStatusUnauthenticated AuthStatus = "unauthenticated"
	
	// AuthStatusAuthenticated indicates an authenticated node
	AuthStatusAuthenticated AuthStatus = "authenticated"
	
	// AuthStatusRevoked indicates a revoked authentication
	AuthStatusRevoked AuthStatus = "revoked"
)

// AuthToken represents an authentication token
type AuthToken struct {
	Token      string    // Token value
	NodeID     string    // Node ID
	Expiration time.Time // Expiration time
	Signature  []byte    // Token signature
}

// NodeAuthenticator represents a node authenticator
type NodeAuthenticator struct {
	Nodes           map[string]*NodeAuth // Map of node ID to node authentication
	TokenExpiration time.Duration       // Token expiration duration
	mu              sync.RWMutex         // Mutex for thread safety
}

// NodeAuth represents node authentication information
type NodeAuth struct {
	NodeID         string      // Node ID
	PublicKey      []byte      // Node public key
	Status         AuthStatus  // Authentication status
	TrustScore     float64     // Trust score
	LastAuth       time.Time   // Last authentication time
	CurrentToken   *AuthToken  // Current authentication token
	AuthFactors    []AuthFactor // Authentication factors
}

// AuthFactor represents an authentication factor
type AuthFactor struct {
	Type  string // Factor type
	Value []byte // Factor value
}

// NewNodeAuthenticator creates a new node authenticator
func NewNodeAuthenticator(tokenExpiration time.Duration) *NodeAuthenticator {
	return &NodeAuthenticator{
		Nodes:           make(map[string]*NodeAuth),
		TokenExpiration: tokenExpiration,
	}
}

// RegisterNode registers a node with the authenticator
func (a *NodeAuthenticator) RegisterNode(nodeID string, publicKey []byte) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	// Check if the node already exists
	if _, ok := a.Nodes[nodeID]; ok {
		return fmt.Errorf("node already registered: %s", nodeID)
	}
	
	// Create a new node authentication
	nodeAuth := &NodeAuth{
		NodeID:     nodeID,
		PublicKey:  publicKey,
		Status:     AuthStatusUnauthenticated,
		TrustScore: 1.0,
		AuthFactors: make([]AuthFactor, 0),
	}
	
	// Add the node to the map
	a.Nodes[nodeID] = nodeAuth
	
	return nil
}

// UnregisterNode unregisters a node from the authenticator
func (a *NodeAuthenticator) UnregisterNode(nodeID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	// Check if the node exists
	if _, ok := a.Nodes[nodeID]; !ok {
		return fmt.Errorf("node not registered: %s", nodeID)
	}
	
	// Remove the node from the map
	delete(a.Nodes, nodeID)
	
	return nil
}

// AddAuthFactor adds an authentication factor to a node
func (a *NodeAuthenticator) AddAuthFactor(nodeID, factorType string, factorValue []byte) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	// Check if the node exists
	nodeAuth, ok := a.Nodes[nodeID]
	if !ok {
		return fmt.Errorf("node not registered: %s", nodeID)
	}
	
	// Add the authentication factor
	nodeAuth.AuthFactors = append(nodeAuth.AuthFactors, AuthFactor{
		Type:  factorType,
		Value: factorValue,
	})
	
	return nil
}

// RemoveAuthFactor removes an authentication factor from a node
func (a *NodeAuthenticator) RemoveAuthFactor(nodeID, factorType string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	// Check if the node exists
	nodeAuth, ok := a.Nodes[nodeID]
	if !ok {
		return fmt.Errorf("node not registered: %s", nodeID)
	}
	
	// Find and remove the authentication factor
	var newFactors []AuthFactor
	found := false
	for _, factor := range nodeAuth.AuthFactors {
		if factor.Type != factorType {
			newFactors = append(newFactors, factor)
		} else {
			found = true
		}
	}
	
	if !found {
		return fmt.Errorf("authentication factor not found: %s", factorType)
	}
	
	nodeAuth.AuthFactors = newFactors
	
	return nil
}

// Authenticate authenticates a node
func (a *NodeAuthenticator) Authenticate(nodeID string, keyPair *crypto.KeyPair, factors map[string][]byte) (*AuthToken, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	// Check if the node exists
	nodeAuth, ok := a.Nodes[nodeID]
	if !ok {
		return nil, fmt.Errorf("node not registered: %s", nodeID)
	}
	
	// Check if the node is revoked
	if nodeAuth.Status == AuthStatusRevoked {
		return nil, fmt.Errorf("node authentication revoked: %s", nodeID)
	}
	
	// Verify the authentication factors
	for _, factor := range nodeAuth.AuthFactors {
		value, ok := factors[factor.Type]
		if !ok {
			return nil, fmt.Errorf("missing authentication factor: %s", factor.Type)
		}
		
		// In a real implementation, we would verify the factor properly
		// For now, we'll just check if the values match
		if string(value) != string(factor.Value) {
			return nil, fmt.Errorf("invalid authentication factor: %s", factor.Type)
		}
	}
	
	// Generate a random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	
	// Encode the token as base64
	tokenStr := base64.StdEncoding.EncodeToString(tokenBytes)
	
	// Create a new token
	token := &AuthToken{
		Token:      tokenStr,
		NodeID:     nodeID,
		Expiration: time.Now().Add(a.TokenExpiration),
	}
	
	// Sign the token
	tokenData := fmt.Sprintf("%s:%s:%d", token.Token, token.NodeID, token.Expiration.Unix())
	signature, err := keyPair.Sign([]byte(tokenData))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}
	
	token.Signature = signature
	
	// Update the node authentication
	nodeAuth.Status = AuthStatusAuthenticated
	nodeAuth.LastAuth = time.Now()
	nodeAuth.CurrentToken = token
	
	return token, nil
}

// VerifyToken verifies an authentication token
func (a *NodeAuthenticator) VerifyToken(token *AuthToken) (bool, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	// Check if the node exists
	nodeAuth, ok := a.Nodes[token.NodeID]
	if !ok {
		return false, fmt.Errorf("node not registered: %s", token.NodeID)
	}
	
	// Check if the node is authenticated
	if nodeAuth.Status != AuthStatusAuthenticated {
		return false, fmt.Errorf("node not authenticated: %s", token.NodeID)
	}
	
	// Check if the token has expired
	if token.Expiration.Before(time.Now()) {
		return false, fmt.Errorf("token expired")
	}
	
	// Verify the token signature
	tokenData := fmt.Sprintf("%s:%s:%d", token.Token, token.NodeID, token.Expiration.Unix())
	valid, err := crypto.VerifySignature(nodeAuth.PublicKey, []byte(tokenData), token.Signature)
	if err != nil {
		return false, fmt.Errorf("failed to verify token signature: %w", err)
	}
	
	return valid, nil
}

// RevokeAuthentication revokes a node's authentication
func (a *NodeAuthenticator) RevokeAuthentication(nodeID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	// Check if the node exists
	nodeAuth, ok := a.Nodes[nodeID]
	if !ok {
		return fmt.Errorf("node not registered: %s", nodeID)
	}
	
	// Revoke the authentication
	nodeAuth.Status = AuthStatusRevoked
	nodeAuth.CurrentToken = nil
	
	return nil
}

// UpdateTrustScore updates a node's trust score
func (a *NodeAuthenticator) UpdateTrustScore(nodeID string, score float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	// Check if the node exists
	nodeAuth, ok := a.Nodes[nodeID]
	if !ok {
		return fmt.Errorf("node not registered: %s", nodeID)
	}
	
	// Update the trust score
	nodeAuth.TrustScore = score
	
	return nil
}

// GetNodeAuthStatus gets a node's authentication status
func (a *NodeAuthenticator) GetNodeAuthStatus(nodeID string) (AuthStatus, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	// Check if the node exists
	nodeAuth, ok := a.Nodes[nodeID]
	if !ok {
		return "", fmt.Errorf("node not registered: %s", nodeID)
	}
	
	return nodeAuth.Status, nil
}

// GetNodeTrustScore gets a node's trust score
func (a *NodeAuthenticator) GetNodeTrustScore(nodeID string) (float64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	// Check if the node exists
	nodeAuth, ok := a.Nodes[nodeID]
	if !ok {
		return 0.0, fmt.Errorf("node not registered: %s", nodeID)
	}
	
	return nodeAuth.TrustScore, nil
}

// GetAuthenticatedNodes gets all authenticated nodes
func (a *NodeAuthenticator) GetAuthenticatedNodes() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	var authenticatedNodes []string
	for nodeID, nodeAuth := range a.Nodes {
		if nodeAuth.Status == AuthStatusAuthenticated {
			authenticatedNodes = append(authenticatedNodes, nodeID)
		}
	}
	
	return authenticatedNodes
}

// CleanupExpiredTokens cleans up expired authentication tokens
func (a *NodeAuthenticator) CleanupExpiredTokens() {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	now := time.Now()
	
	for _, nodeAuth := range a.Nodes {
		if nodeAuth.Status == AuthStatusAuthenticated && nodeAuth.CurrentToken != nil {
			if nodeAuth.CurrentToken.Expiration.Before(now) {
				// The token has expired
				nodeAuth.Status = AuthStatusUnauthenticated
				nodeAuth.CurrentToken = nil
			}
		}
	}
}
