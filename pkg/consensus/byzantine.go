package consensus

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/advanced-blockchain/pkg/crypto"
)

// ByzantineDetectionStatus represents the status of a Byzantine detection
type ByzantineDetectionStatus string

const (
	// ByzantineDetectionStatusNormal indicates a normal node
	ByzantineDetectionStatusNormal ByzantineDetectionStatus = "normal"
	
	// ByzantineDetectionStatusSuspicious indicates a suspicious node
	ByzantineDetectionStatusSuspicious ByzantineDetectionStatus = "suspicious"
	
	// ByzantineDetectionStatusByzantine indicates a Byzantine node
	ByzantineDetectionStatusByzantine ByzantineDetectionStatus = "byzantine"
)

// ByzantineDetector represents a Byzantine fault detector
type ByzantineDetector struct {
	Nodes               map[string]*ByzantineNode // Map of node ID to node
	SuspicionThreshold  float64                   // Threshold for suspicion
	ByzantineThreshold  float64                   // Threshold for Byzantine detection
	DetectionWindow     time.Duration             // Time window for detection
	DetectionHistory    map[string][]ByzantineEvent // Map of node ID to detection history
	mu                  sync.RWMutex              // Mutex for thread safety
}

// ByzantineNode represents a node in the Byzantine detection system
type ByzantineNode struct {
	ID              string                    // Node ID
	PublicKey       []byte                    // Node public key
	Status          ByzantineDetectionStatus  // Node status
	SuspicionScore  float64                   // Suspicion score
	LastUpdated     time.Time                 // Last updated time
}

// ByzantineEvent represents a Byzantine detection event
type ByzantineEvent struct {
	NodeID      string                    // Node ID
	EventType   string                    // Event type
	Timestamp   time.Time                 // Event timestamp
	Description string                    // Event description
	Score       float64                   // Event score
}

// NewByzantineDetector creates a new Byzantine detector
func NewByzantineDetector() *ByzantineDetector {
	return &ByzantineDetector{
		Nodes:              make(map[string]*ByzantineNode),
		SuspicionThreshold: 0.7,
		ByzantineThreshold: 0.9,
		DetectionWindow:    24 * time.Hour,
		DetectionHistory:   make(map[string][]ByzantineEvent),
	}
}

// AddNode adds a node to the detector
func (d *ByzantineDetector) AddNode(id string, publicKey []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	// Check if the node already exists
	if _, ok := d.Nodes[id]; ok {
		return fmt.Errorf("node already exists: %s", id)
	}
	
	// Create a new node
	node := &ByzantineNode{
		ID:             id,
		PublicKey:      publicKey,
		Status:         ByzantineDetectionStatusNormal,
		SuspicionScore: 0.0,
		LastUpdated:    time.Now(),
	}
	
	// Add the node to the map
	d.Nodes[id] = node
	
	// Initialize the detection history
	d.DetectionHistory[id] = make([]ByzantineEvent, 0)
	
	return nil
}

// RemoveNode removes a node from the detector
func (d *ByzantineDetector) RemoveNode(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	// Check if the node exists
	if _, ok := d.Nodes[id]; !ok {
		return fmt.Errorf("node not found: %s", id)
	}
	
	// Remove the node from the map
	delete(d.Nodes, id)
	
	// Remove the detection history
	delete(d.DetectionHistory, id)
	
	return nil
}

// ReportEvent reports a Byzantine detection event
func (d *ByzantineDetector) ReportEvent(nodeID, eventType, description string, score float64) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	// Check if the node exists
	node, ok := d.Nodes[nodeID]
	if !ok {
		return fmt.Errorf("node not found: %s", nodeID)
	}
	
	// Create a new event
	event := ByzantineEvent{
		NodeID:      nodeID,
		EventType:   eventType,
		Timestamp:   time.Now(),
		Description: description,
		Score:       score,
	}
	
	// Add the event to the detection history
	d.DetectionHistory[nodeID] = append(d.DetectionHistory[nodeID], event)
	
	// Update the node's suspicion score
	d.updateSuspicionScore(nodeID)
	
	// Update the node's status
	d.updateNodeStatus(node)
	
	return nil
}

// updateSuspicionScore updates the suspicion score of a node
func (d *ByzantineDetector) updateSuspicionScore(nodeID string) {
	// Get the node
	node := d.Nodes[nodeID]
	
	// Get the detection history
	history := d.DetectionHistory[nodeID]
	
	// Filter the history to only include events within the detection window
	var recentEvents []ByzantineEvent
	now := time.Now()
	for _, event := range history {
		if now.Sub(event.Timestamp) <= d.DetectionWindow {
			recentEvents = append(recentEvents, event)
		}
	}
	
	// Compute the suspicion score as the average of the event scores
	var totalScore float64
	for _, event := range recentEvents {
		totalScore += event.Score
	}
	
	if len(recentEvents) > 0 {
		node.SuspicionScore = totalScore / float64(len(recentEvents))
	} else {
		node.SuspicionScore = 0.0
	}
	
	// Update the last updated time
	node.LastUpdated = now
}

// updateNodeStatus updates the status of a node
func (d *ByzantineDetector) updateNodeStatus(node *ByzantineNode) {
	// Update the node status based on the suspicion score
	if node.SuspicionScore >= d.ByzantineThreshold {
		node.Status = ByzantineDetectionStatusByzantine
	} else if node.SuspicionScore >= d.SuspicionThreshold {
		node.Status = ByzantineDetectionStatusSuspicious
	} else {
		node.Status = ByzantineDetectionStatusNormal
	}
}

// GetNodeStatus gets the status of a node
func (d *ByzantineDetector) GetNodeStatus(nodeID string) (ByzantineDetectionStatus, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	// Check if the node exists
	node, ok := d.Nodes[nodeID]
	if !ok {
		return "", fmt.Errorf("node not found: %s", nodeID)
	}
	
	return node.Status, nil
}

// GetNodeSuspicionScore gets the suspicion score of a node
func (d *ByzantineDetector) GetNodeSuspicionScore(nodeID string) (float64, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	// Check if the node exists
	node, ok := d.Nodes[nodeID]
	if !ok {
		return 0.0, fmt.Errorf("node not found: %s", nodeID)
	}
	
	return node.SuspicionScore, nil
}

// GetByzantineNodes gets all Byzantine nodes
func (d *ByzantineDetector) GetByzantineNodes() []*ByzantineNode {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	var byzantineNodes []*ByzantineNode
	for _, node := range d.Nodes {
		if node.Status == ByzantineDetectionStatusByzantine {
			byzantineNodes = append(byzantineNodes, node)
		}
	}
	
	return byzantineNodes
}

// GetSuspiciousNodes gets all suspicious nodes
func (d *ByzantineDetector) GetSuspiciousNodes() []*ByzantineNode {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	var suspiciousNodes []*ByzantineNode
	for _, node := range d.Nodes {
		if node.Status == ByzantineDetectionStatusSuspicious {
			suspiciousNodes = append(suspiciousNodes, node)
		}
	}
	
	return suspiciousNodes
}

// ComputeEntropyBasedScore computes an entropy-based score for a node
func (d *ByzantineDetector) ComputeEntropyBasedScore(nodeID string) (float64, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	// Check if the node exists
	_, ok := d.Nodes[nodeID]
	if !ok {
		return 0.0, fmt.Errorf("node not found: %s", nodeID)
	}
	
	// Get the detection history
	history := d.DetectionHistory[nodeID]
	
	// Filter the history to only include events within the detection window
	var recentEvents []ByzantineEvent
	now := time.Now()
	for _, event := range history {
		if now.Sub(event.Timestamp) <= d.DetectionWindow {
			recentEvents = append(recentEvents, event)
		}
	}
	
	// Compute the entropy of the event scores
	var entropy float64
	var totalScore float64
	
	// Count the occurrences of each event type
	eventTypeCounts := make(map[string]int)
	for _, event := range recentEvents {
		eventTypeCounts[event.EventType]++
		totalScore += event.Score
	}
	
	// Compute the entropy
	for _, count := range eventTypeCounts {
		p := float64(count) / float64(len(recentEvents))
		entropy -= p * math.Log2(p)
	}
	
	// Normalize the entropy to [0, 1]
	maxEntropy := math.Log2(float64(len(eventTypeCounts)))
	if maxEntropy > 0 {
		entropy /= maxEntropy
	}
	
	// Combine the entropy with the average score
	var avgScore float64
	if len(recentEvents) > 0 {
		avgScore = totalScore / float64(len(recentEvents))
	}
	
	// The final score is a weighted combination of the entropy and the average score
	// Higher entropy means more diverse events, which could indicate Byzantine behavior
	// Higher average score means more severe events
	return 0.3*entropy + 0.7*avgScore, nil
}

// VerifyNodeSignature verifies a signature from a node
func (d *ByzantineDetector) VerifyNodeSignature(nodeID string, message, signature []byte) (bool, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	// Check if the node exists
	node, ok := d.Nodes[nodeID]
	if !ok {
		return false, fmt.Errorf("node not found: %s", nodeID)
	}
	
	// Verify the signature
	return crypto.VerifySignature(node.PublicKey, message, signature)
}

// SetSuspicionThreshold sets the suspicion threshold
func (d *ByzantineDetector) SetSuspicionThreshold(threshold float64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	d.SuspicionThreshold = threshold
	
	// Update the status of all nodes
	for _, node := range d.Nodes {
		d.updateNodeStatus(node)
	}
}

// SetByzantineThreshold sets the Byzantine threshold
func (d *ByzantineDetector) SetByzantineThreshold(threshold float64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	d.ByzantineThreshold = threshold
	
	// Update the status of all nodes
	for _, node := range d.Nodes {
		d.updateNodeStatus(node)
	}
}

// SetDetectionWindow sets the detection window
func (d *ByzantineDetector) SetDetectionWindow(window time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	d.DetectionWindow = window
	
	// Update the suspicion scores of all nodes
	for nodeID := range d.Nodes {
		d.updateSuspicionScore(nodeID)
	}
	
	// Update the status of all nodes
	for _, node := range d.Nodes {
		d.updateNodeStatus(node)
	}
}
