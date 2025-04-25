package consensus

import (
	"fmt"
	"sync"
	"time"
)

// ConsistencyLevel represents a consistency level
type ConsistencyLevel string

const (
	// ConsistencyLevelStrong indicates strong consistency
	ConsistencyLevelStrong ConsistencyLevel = "strong"

	// ConsistencyLevelEventual indicates eventual consistency
	ConsistencyLevelEventual ConsistencyLevel = "eventual"

	// ConsistencyLevelCausal indicates causal consistency
	ConsistencyLevelCausal ConsistencyLevel = "causal"
)

// NetworkState represents the state of the network
type NetworkState string

const (
	// NetworkStateNormal indicates a normal network state
	NetworkStateNormal NetworkState = "normal"

	// NetworkStateDegraded indicates a degraded network state
	NetworkStateDegraded NetworkState = "degraded"

	// NetworkStatePartitioned indicates a partitioned network state
	NetworkStatePartitioned NetworkState = "partitioned"
)

// AdaptiveConsistency represents an adaptive consistency orchestrator
type AdaptiveConsistency struct {
	CurrentLevel         ConsistencyLevel         // Current consistency level
	NetworkState         NetworkState             // Current network state
	PartitionProbability float64                  // Probability of network partition
	Timeouts             map[string]time.Duration // Map of operation type to timeout
	RetryCount           map[string]int           // Map of operation type to retry count
	NetworkTelemetry     []NetworkEvent           // Network telemetry events
	TelemetryWindow      time.Duration            // Time window for telemetry
	mu                   sync.RWMutex             // Mutex for thread safety
}

// NetworkEvent represents a network telemetry event
type NetworkEvent struct {
	Timestamp  time.Time     // Event timestamp
	Latency    time.Duration // Network latency
	PacketLoss float64       // Packet loss rate
	Throughput float64       // Network throughput
	ErrorRate  float64       // Error rate
}

// NewAdaptiveConsistency creates a new adaptive consistency orchestrator
func NewAdaptiveConsistency() *AdaptiveConsistency {
	return &AdaptiveConsistency{
		CurrentLevel:         ConsistencyLevelStrong,
		NetworkState:         NetworkStateNormal,
		PartitionProbability: 0.0,
		Timeouts:             make(map[string]time.Duration),
		RetryCount:           make(map[string]int),
		NetworkTelemetry:     make([]NetworkEvent, 0),
		TelemetryWindow:      1 * time.Hour,
	}
}

// SetConsistencyLevel sets the consistency level
func (a *AdaptiveConsistency) SetConsistencyLevel(level ConsistencyLevel) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.CurrentLevel = level
}

// GetConsistencyLevel gets the current consistency level
func (a *AdaptiveConsistency) GetConsistencyLevel() ConsistencyLevel {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.CurrentLevel
}

// SetNetworkState sets the network state
func (a *AdaptiveConsistency) SetNetworkState(state NetworkState) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.NetworkState = state

	// Update the consistency level based on the network state
	a.updateConsistencyLevel()
}

// GetNetworkState gets the current network state
func (a *AdaptiveConsistency) GetNetworkState() NetworkState {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.NetworkState
}

// SetTimeout sets the timeout for an operation type
func (a *AdaptiveConsistency) SetTimeout(operationType string, timeout time.Duration) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.Timeouts[operationType] = timeout
}

// GetTimeout gets the timeout for an operation type
func (a *AdaptiveConsistency) GetTimeout(operationType string) (time.Duration, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	timeout, ok := a.Timeouts[operationType]
	if !ok {
		return 0, fmt.Errorf("timeout not found for operation type: %s", operationType)
	}

	return timeout, nil
}

// SetRetryCount sets the retry count for an operation type
func (a *AdaptiveConsistency) SetRetryCount(operationType string, count int) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.RetryCount[operationType] = count
}

// GetRetryCount gets the retry count for an operation type
func (a *AdaptiveConsistency) GetRetryCount(operationType string) (int, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	count, ok := a.RetryCount[operationType]
	if !ok {
		return 0, fmt.Errorf("retry count not found for operation type: %s", operationType)
	}

	return count, nil
}

// AddNetworkEvent adds a network telemetry event
func (a *AdaptiveConsistency) AddNetworkEvent(event NetworkEvent) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Add the event to the telemetry
	a.NetworkTelemetry = append(a.NetworkTelemetry, event)

	// Prune old events
	a.pruneNetworkTelemetry()

	// Update the network state and partition probability
	a.updateNetworkState()

	// Update the consistency level
	a.updateConsistencyLevel()

	// Update the timeouts and retry counts
	a.updateTimeoutsAndRetries()
}

// pruneNetworkTelemetry prunes old network telemetry events
func (a *AdaptiveConsistency) pruneNetworkTelemetry() {
	// Filter the telemetry to only include events within the telemetry window
	var recentEvents []NetworkEvent
	now := time.Now()
	for _, event := range a.NetworkTelemetry {
		if now.Sub(event.Timestamp) <= a.TelemetryWindow {
			recentEvents = append(recentEvents, event)
		}
	}

	a.NetworkTelemetry = recentEvents
}

// updateNetworkState updates the network state and partition probability
func (a *AdaptiveConsistency) updateNetworkState() {
	// Compute the average latency, packet loss, throughput, and error rate
	var totalLatency time.Duration
	var totalPacketLoss float64
	var totalThroughput float64
	var totalErrorRate float64

	for _, event := range a.NetworkTelemetry {
		totalLatency += event.Latency
		totalPacketLoss += event.PacketLoss
		totalThroughput += event.Throughput
		totalErrorRate += event.ErrorRate
	}

	var avgLatency time.Duration
	var avgPacketLoss float64
	var avgErrorRate float64

	if len(a.NetworkTelemetry) > 0 {
		avgLatency = totalLatency / time.Duration(len(a.NetworkTelemetry))
		avgPacketLoss = totalPacketLoss / float64(len(a.NetworkTelemetry))
		// We're not using avgThroughput for now
		avgErrorRate = totalErrorRate / float64(len(a.NetworkTelemetry))
	}

	// Update the network state based on the average metrics
	if avgPacketLoss > 0.5 || avgErrorRate > 0.5 {
		a.NetworkState = NetworkStatePartitioned
		a.PartitionProbability = 0.9
	} else if avgPacketLoss > 0.2 || avgErrorRate > 0.2 || avgLatency > 1*time.Second {
		a.NetworkState = NetworkStateDegraded
		a.PartitionProbability = 0.5
	} else {
		a.NetworkState = NetworkStateNormal
		a.PartitionProbability = 0.1
	}
}

// updateConsistencyLevel updates the consistency level based on the network state
func (a *AdaptiveConsistency) updateConsistencyLevel() {
	// Update the consistency level based on the network state
	switch a.NetworkState {
	case NetworkStateNormal:
		a.CurrentLevel = ConsistencyLevelStrong
	case NetworkStateDegraded:
		a.CurrentLevel = ConsistencyLevelCausal
	case NetworkStatePartitioned:
		a.CurrentLevel = ConsistencyLevelEventual
	}
}

// updateTimeoutsAndRetries updates the timeouts and retry counts based on the network state
func (a *AdaptiveConsistency) updateTimeoutsAndRetries() {
	// Update the timeouts and retry counts based on the network state
	switch a.NetworkState {
	case NetworkStateNormal:
		a.Timeouts["read"] = 1 * time.Second
		a.Timeouts["write"] = 2 * time.Second
		a.RetryCount["read"] = 3
		a.RetryCount["write"] = 3
	case NetworkStateDegraded:
		a.Timeouts["read"] = 2 * time.Second
		a.Timeouts["write"] = 4 * time.Second
		a.RetryCount["read"] = 5
		a.RetryCount["write"] = 5
	case NetworkStatePartitioned:
		a.Timeouts["read"] = 5 * time.Second
		a.Timeouts["write"] = 10 * time.Second
		a.RetryCount["read"] = 10
		a.RetryCount["write"] = 10
	}
}

// PredictPartitionProbability predicts the probability of a network partition
func (a *AdaptiveConsistency) PredictPartitionProbability() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.PartitionProbability
}

// ShouldRetry determines if an operation should be retried
func (a *AdaptiveConsistency) ShouldRetry(operationType string, attempt int) (bool, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Get the retry count for the operation type
	maxRetries, ok := a.RetryCount[operationType]
	if !ok {
		return false, fmt.Errorf("retry count not found for operation type: %s", operationType)
	}

	// Check if we should retry
	return attempt < maxRetries, nil
}

// GetBackoffDuration gets the backoff duration for a retry attempt
func (a *AdaptiveConsistency) GetBackoffDuration(operationType string, attempt int) (time.Duration, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Get the timeout for the operation type
	baseTimeout, ok := a.Timeouts[operationType]
	if !ok {
		return 0, fmt.Errorf("timeout not found for operation type: %s", operationType)
	}

	// Compute the backoff duration using exponential backoff
	// backoff = baseTimeout * 2^attempt
	backoff := baseTimeout
	for i := 0; i < attempt; i++ {
		backoff *= 2
	}

	return backoff, nil
}
