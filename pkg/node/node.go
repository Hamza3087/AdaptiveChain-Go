package node

import (
	"fmt"
	"sync"
	"time"

	"github.com/advanced-blockchain/pkg/amf"
	"github.com/advanced-blockchain/pkg/consensus"
	"github.com/advanced-blockchain/pkg/crypto"
	"github.com/advanced-blockchain/pkg/network"
	"github.com/advanced-blockchain/pkg/state"
)

// NodeStatus represents the status of a node
type NodeStatus string

const (
	// NodeStatusInitializing indicates an initializing node
	NodeStatusInitializing NodeStatus = "initializing"

	// NodeStatusRunning indicates a running node
	NodeStatusRunning NodeStatus = "running"

	// NodeStatusStopping indicates a stopping node
	NodeStatusStopping NodeStatus = "stopping"

	// NodeStatusStopped indicates a stopped node
	NodeStatusStopped NodeStatus = "stopped"
)

// Node represents a blockchain node
type Node struct {
	ID                  string                         // Node ID
	Address             string                         // Node address
	KeyPair             *crypto.KeyPair                // Node key pair
	State               *state.State                   // Blockchain state
	Forest              *amf.AdaptiveMerkleForest      // Adaptive Merkle Forest
	Consensus           *consensus.HybridConsensus     // Consensus mechanism
	Network             *network.Network               // Network manager
	ByzantineDetector   *consensus.ByzantineDetector   // Byzantine detector
	AdaptiveConsistency *consensus.AdaptiveConsistency // Adaptive consistency
	Status              NodeStatus                     // Node status
	StartTime           time.Time                      // Node start time
	Config              *NodeConfig                    // Node configuration
	mu                  sync.RWMutex                   // Mutex for thread safety
}

// NodeConfig represents the configuration of a node
type NodeConfig struct {
	SplitThreshold    float64       // Threshold for splitting a shard
	RebalanceInterval time.Duration // Interval for rebalancing the forest
	PruningHeight     uint64        // Height for pruning old blocks
	ConsensusInterval time.Duration // Interval for consensus rounds
	TelemetryInterval time.Duration // Interval for telemetry collection
}

// NewNode creates a new blockchain node
func NewNode(id, address string, keyPair *crypto.KeyPair, config *NodeConfig) *Node {
	// Create the blockchain state
	state := state.NewState()

	// Create the adaptive Merkle forest
	forest := amf.NewAdaptiveMerkleForest(config.SplitThreshold, config.RebalanceInterval)

	// Create the consensus mechanism
	hybridConsensus := consensus.NewHybridConsensus(state)

	// Create the network manager
	networkManager := network.NewNetwork(id, address, keyPair)

	// Create the Byzantine detector
	byzantineDetector := consensus.NewByzantineDetector()

	// Create the adaptive consistency
	adaptiveConsistency := consensus.NewAdaptiveConsistency()

	// Create the node
	node := &Node{
		ID:                  id,
		Address:             address,
		KeyPair:             keyPair,
		State:               state,
		Forest:              forest,
		Consensus:           hybridConsensus,
		Network:             networkManager,
		ByzantineDetector:   byzantineDetector,
		AdaptiveConsistency: adaptiveConsistency,
		Status:              NodeStatusStopped,
		Config:              config,
	}

	// Set the pruning height
	state.SetPruningHeight(config.PruningHeight)

	return node
}

// Start starts the node
func (n *Node) Start() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Check if the node is already running
	if n.Status == NodeStatusRunning {
		return fmt.Errorf("node already running")
	}

	// Set the status to initializing
	n.Status = NodeStatusInitializing

	// Start the network
	if err := n.Network.Start(); err != nil {
		n.Status = NodeStatusStopped
		return fmt.Errorf("failed to start network: %w", err)
	}

	// Register message handlers
	n.registerMessageHandlers()

	// Set the status to running
	n.Status = NodeStatusRunning
	n.StartTime = time.Now()

	// Start the consensus in a separate goroutine
	go n.runConsensus()

	// Start telemetry collection in a separate goroutine
	go n.collectTelemetry()

	return nil
}

// Stop stops the node
func (n *Node) Stop() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Check if the node is running
	if n.Status != NodeStatusRunning {
		return fmt.Errorf("node not running")
	}

	// Set the status to stopping
	n.Status = NodeStatusStopping

	// Stop the network
	if err := n.Network.Stop(); err != nil {
		// Log the error but continue
		fmt.Printf("Error stopping network: %v\n", err)
	}

	// Set the status to stopped
	n.Status = NodeStatusStopped

	return nil
}

// registerMessageHandlers registers message handlers
func (n *Node) registerMessageHandlers() {
	// Register handshake handler
	n.Network.RegisterMessageHandler(network.MessageTypeHandshake, n.handleHandshake)

	// Register ping handler
	n.Network.RegisterMessageHandler(network.MessageTypePing, n.handlePing)

	// Register pong handler
	n.Network.RegisterMessageHandler(network.MessageTypePong, n.handlePong)

	// Register transaction handler
	n.Network.RegisterMessageHandler(network.MessageTypeTransaction, n.handleTransaction)

	// Register block handler
	n.Network.RegisterMessageHandler(network.MessageTypeBlock, n.handleBlock)

	// Register consensus handler
	n.Network.RegisterMessageHandler(network.MessageTypeConsensus, n.handleConsensus)

	// Register peer discovery handler
	n.Network.RegisterMessageHandler(network.MessageTypePeerDiscovery, n.handlePeerDiscovery)
}

// handleHandshake handles a handshake message
func (n *Node) handleHandshake(message *network.Message, peer *network.Peer) error {
	// Unmarshal the handshake payload
	var payload network.HandshakePayload
	if err := payload.Unmarshal(message.Payload); err != nil {
		return fmt.Errorf("failed to unmarshal handshake payload: %w", err)
	}

	// Update the peer's public key
	peer.PublicKey = payload.PublicKey

	return nil
}

// handlePing handles a ping message
func (n *Node) handlePing(message *network.Message, peer *network.Peer) error {
	// Unmarshal the ping payload
	var payload network.PingPayload
	if err := payload.Unmarshal(message.Payload); err != nil {
		return fmt.Errorf("failed to unmarshal ping payload: %w", err)
	}

	// Create a pong message
	pongMessage, err := network.NewPongMessage(n.ID, peer.ID, payload.Nonce)
	if err != nil {
		return fmt.Errorf("failed to create pong message: %w", err)
	}

	// Send the pong message
	if err := n.Network.SendMessage(peer.ID, pongMessage); err != nil {
		return fmt.Errorf("failed to send pong message: %w", err)
	}

	return nil
}

// handlePong handles a pong message
func (n *Node) handlePong(message *network.Message, peer *network.Peer) error {
	// Unmarshal the pong payload
	var payload network.PongPayload
	if err := payload.Unmarshal(message.Payload); err != nil {
		return fmt.Errorf("failed to unmarshal pong payload: %w", err)
	}

	// Compute the latency
	latency := time.Since(payload.Timestamp)

	// Update the peer's latency
	peer.UpdateLatency(latency)

	return nil
}

// handleTransaction handles a transaction message
func (n *Node) handleTransaction(message *network.Message, peer *network.Peer) error {
	// Unmarshal the transaction payload
	var payload network.TransactionPayload
	if err := payload.Unmarshal(message.Payload); err != nil {
		return fmt.Errorf("failed to unmarshal transaction payload: %w", err)
	}

	// Unmarshal the transaction
	var tx state.Transaction
	if err := tx.Unmarshal(payload.Transaction); err != nil {
		return fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	// Verify the transaction
	valid, err := tx.Verify()
	if err != nil {
		return fmt.Errorf("failed to verify transaction: %w", err)
	}
	if !valid {
		return fmt.Errorf("invalid transaction signature")
	}

	// Add the transaction to the state
	if err := n.State.AddTransaction(&tx); err != nil {
		return fmt.Errorf("failed to add transaction to state: %w", err)
	}

	// Update the shard state
	shardID := fmt.Sprintf("shard-%d", tx.ShardID)
	shard, err := n.Forest.GetShard(shardID)
	if err != nil {
		return fmt.Errorf("failed to get shard: %w", err)
	}

	if err := shard.UpdateState(tx.ID.String(), &tx); err != nil {
		return fmt.Errorf("failed to update shard state: %w", err)
	}

	return nil
}

// handleBlock handles a block message
func (n *Node) handleBlock(message *network.Message, peer *network.Peer) error {
	// Unmarshal the block payload
	var payload network.BlockPayload
	if err := payload.Unmarshal(message.Payload); err != nil {
		return fmt.Errorf("failed to unmarshal block payload: %w", err)
	}

	// Unmarshal the block
	var block state.Block
	if err := block.Unmarshal(payload.Block); err != nil {
		return fmt.Errorf("failed to unmarshal block: %w", err)
	}

	// Verify the block
	publicKey, err := n.KeyPair.ExportPublicKey()
	if err != nil {
		return fmt.Errorf("failed to export public key: %w", err)
	}

	valid, err := block.Verify(publicKey)
	if err != nil {
		return fmt.Errorf("failed to verify block: %w", err)
	}
	if !valid {
		return fmt.Errorf("invalid block signature")
	}

	// Add the block to the state
	if err := n.State.AddBlock(&block); err != nil {
		return fmt.Errorf("failed to add block to state: %w", err)
	}

	return nil
}

// handleConsensus handles a consensus message
func (n *Node) handleConsensus(message *network.Message, peer *network.Peer) error {
	// Unmarshal the consensus payload
	var payload network.ConsensusPayload
	if err := payload.Unmarshal(message.Payload); err != nil {
		return fmt.Errorf("failed to unmarshal consensus payload: %w", err)
	}

	// Handle the consensus message based on its type
	switch payload.Type {
	case "prepare":
		// Handle prepare message
		// ...
	case "commit":
		// Handle commit message
		// ...
	default:
		return fmt.Errorf("unknown consensus message type: %s", payload.Type)
	}

	return nil
}

// handlePeerDiscovery handles a peer discovery message
func (n *Node) handlePeerDiscovery(message *network.Message, peer *network.Peer) error {
	// Unmarshal the peer discovery payload
	var payload network.PeerDiscoveryPayload
	if err := payload.Unmarshal(message.Payload); err != nil {
		return fmt.Errorf("failed to unmarshal peer discovery payload: %w", err)
	}

	// Add the peers to the network
	for _, peerInfo := range payload.Peers {
		// Skip if the peer is already known
		if _, err := n.Network.GetPeer(peerInfo.ID); err == nil {
			continue
		}

		// Add the peer to the network
		if err := n.Network.AddPeer(peerInfo.ID, peerInfo.Address, peerInfo.PublicKey); err != nil {
			// Log the error but continue
			fmt.Printf("Error adding peer %s: %v\n", peerInfo.ID, err)
			continue
		}
	}

	return nil
}

// runConsensus runs the consensus mechanism
func (n *Node) runConsensus() {
	ticker := time.NewTicker(n.Config.ConsensusInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check if the node is still running
			n.mu.RLock()
			status := n.Status
			n.mu.RUnlock()

			if status != NodeStatusRunning {
				// The node has been stopped, so exit the loop
				return
			}

			// Start a new consensus round
			if err := n.Consensus.StartConsensus(); err != nil {
				// Log the error but continue
				fmt.Printf("Error starting consensus: %v\n", err)
				continue
			}
		}
	}
}

// collectTelemetry collects network telemetry
func (n *Node) collectTelemetry() {
	ticker := time.NewTicker(n.Config.TelemetryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check if the node is still running
			n.mu.RLock()
			status := n.Status
			n.mu.RUnlock()

			if status != NodeStatusRunning {
				// The node has been stopped, so exit the loop
				return
			}

			// Collect telemetry from connected peers
			connectedPeers := n.Network.GetConnectedPeers()

			// Compute average latency
			var totalLatency time.Duration
			for _, peer := range connectedPeers {
				totalLatency += peer.GetLatency()
			}

			var avgLatency time.Duration
			if len(connectedPeers) > 0 {
				avgLatency = totalLatency / time.Duration(len(connectedPeers))
			}

			// Create a network event
			event := consensus.NetworkEvent{
				Timestamp:  time.Now(),
				Latency:    avgLatency,
				PacketLoss: 0.0, // Placeholder
				Throughput: 0.0, // Placeholder
				ErrorRate:  0.0, // Placeholder
			}

			// Add the event to the adaptive consistency
			n.AdaptiveConsistency.AddNetworkEvent(event)
		}
	}
}

// CreateTransaction creates a new transaction
func (n *Node) CreateTransaction(to []byte, amount uint64, data []byte, shardID uint64) (*state.Transaction, error) {
	// Export the public key
	publicKey, err := n.KeyPair.ExportPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to export public key: %w", err)
	}

	// Create a new transaction
	tx := state.NewTransaction(publicKey, to, amount, data, shardID)

	// Sign the transaction
	if err := tx.Sign(n.KeyPair); err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Add the transaction to the state
	if err := n.State.AddTransaction(tx); err != nil {
		return nil, fmt.Errorf("failed to add transaction to state: %w", err)
	}

	// Update the shard state
	shardIDStr := fmt.Sprintf("shard-%d", shardID)
	shard, err := n.Forest.GetShard(shardIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get shard: %w", err)
	}

	if err := shard.UpdateState(tx.ID.String(), tx); err != nil {
		return nil, fmt.Errorf("failed to update shard state: %w", err)
	}

	// Broadcast the transaction
	txBytes, err := tx.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}

	// Create a transaction message
	message, err := network.NewTransactionMessage(n.ID, "", tx.ID.String(), txBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction message: %w", err)
	}

	// Broadcast the message
	if err := n.Network.BroadcastMessage(message); err != nil {
		return nil, fmt.Errorf("failed to broadcast transaction message: %w", err)
	}

	return tx, nil
}

// GetTransaction gets a transaction by ID
func (n *Node) GetTransaction(txID string) (*state.Transaction, error) {
	return n.State.GetTransaction(txID)
}

// GetBlock gets a block by hash
func (n *Node) GetBlock(blockHash string) (*state.Block, error) {
	return n.State.GetBlock(blockHash)
}

// GetLatestBlock gets the latest block
func (n *Node) GetLatestBlock() *state.Block {
	return n.State.GetLatestBlock()
}

// GetStateRoot gets the state root
func (n *Node) GetStateRoot() crypto.Hash {
	return n.State.GetStateRoot()
}

// GetShardState gets the state of a shard
func (n *Node) GetShardState(shardID string, key string) (interface{}, error) {
	return n.Forest.GetState(shardID, key)
}

// GenerateShardProof generates a proof for a shard state
func (n *Node) GenerateShardProof(shardID string, key string) ([]byte, error) {
	return n.Forest.GenerateProof(shardID, key)
}

// VerifyShardProof verifies a proof for a shard state
func (n *Node) VerifyShardProof(shardID string, key string, value string, proof []byte) (bool, error) {
	return n.Forest.VerifyProof(shardID, key, value, proof)
}

// GetStatus gets the node status
func (n *Node) GetStatus() NodeStatus {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.Status
}

// GetUptime gets the node uptime
func (n *Node) GetUptime() time.Duration {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if n.Status != NodeStatusRunning {
		return 0
	}

	return time.Since(n.StartTime)
}

// GetPeerCount gets the number of peers
func (n *Node) GetPeerCount() int {
	return n.Network.GetPeerCount()
}

// GetConnectedPeerCount gets the number of connected peers
func (n *Node) GetConnectedPeerCount() int {
	return n.Network.GetConnectedPeerCount()
}

// GetConsistencyLevel gets the current consistency level
func (n *Node) GetConsistencyLevel() consensus.ConsistencyLevel {
	return n.AdaptiveConsistency.GetConsistencyLevel()
}

// GetNetworkState gets the current network state
func (n *Node) GetNetworkState() consensus.NetworkState {
	return n.AdaptiveConsistency.GetNetworkState()
}

// GetPartitionProbability gets the probability of a network partition
func (n *Node) GetPartitionProbability() float64 {
	return n.AdaptiveConsistency.PredictPartitionProbability()
}
