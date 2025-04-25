package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/advanced-blockchain/pkg/crypto"
	"github.com/advanced-blockchain/pkg/node"
)

var blockchainNode *node.Node

func main() {
	// Parse command line flags
	nodeID := flag.String("id", "api-node", "Node ID")
	address := flag.String("addr", "127.0.0.1:8001", "Node address")
	apiAddr := flag.String("api", "127.0.0.1:8080", "API server address")
	keyFile := flag.String("key", "", "Private key file")
	splitThreshold := flag.Float64("split", 10.0, "Shard split threshold")
	rebalanceInterval := flag.Duration("rebalance", 5*time.Minute, "Shard rebalance interval")
	pruningHeight := flag.Uint64("prune", 1000, "Block pruning height")
	consensusInterval := flag.Duration("consensus", 10*time.Second, "Consensus interval")
	telemetryInterval := flag.Duration("telemetry", 1*time.Minute, "Telemetry interval")
	flag.Parse()

	// Create or load a key pair
	var keyPair *crypto.KeyPair
	var err error
	if *keyFile != "" {
		// Load the key pair from file
		keyBytes, err := os.ReadFile(*keyFile)
		if err != nil {
			fmt.Printf("Error loading key file: %v\n", err)
			os.Exit(1)
		}

		keyPair, err = crypto.ImportPrivateKey(keyBytes)
		if err != nil {
			fmt.Printf("Error importing private key: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Loaded key pair from file")
	} else {
		// Generate a new key pair
		keyPair, err = crypto.GenerateKeyPair()
		if err != nil {
			fmt.Printf("Error generating key pair: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Generated new key pair")
	}

	// Create the node configuration
	config := &node.NodeConfig{
		SplitThreshold:    *splitThreshold,
		RebalanceInterval: *rebalanceInterval,
		PruningHeight:     *pruningHeight,
		ConsensusInterval: *consensusInterval,
		TelemetryInterval: *telemetryInterval,
	}

	// Create the node
	blockchainNode = node.NewNode(*nodeID, *address, keyPair, config)

	// Start the node
	if err := blockchainNode.Start(); err != nil {
		fmt.Printf("Error starting node: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Node started with ID %s at address %s\n", *nodeID, *address)

	// Set up API routes
	http.HandleFunc("/api/transaction", handleTransaction)
	http.HandleFunc("/api/block", handleBlock)
	http.HandleFunc("/api/latest-block", handleLatestBlock)
	http.HandleFunc("/api/shard", handleShard)
	http.HandleFunc("/api/proof", handleProof)
	http.HandleFunc("/api/status", handleStatus)

	// Start the API server
	fmt.Printf("API server listening on %s\n", *apiAddr)
	if err := http.ListenAndServe(*apiAddr, nil); err != nil {
		fmt.Printf("Error starting API server: %v\n", err)
		os.Exit(1)
	}
}

func handleTransaction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Get a transaction by ID
		txID := r.URL.Query().Get("id")
		if txID == "" {
			http.Error(w, "Missing transaction ID", http.StatusBadRequest)
			return
		}

		tx, err := blockchainNode.GetTransaction(txID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting transaction: %v", err), http.StatusNotFound)
			return
		}

		// Marshal the transaction to JSON
		txJSON, err := json.Marshal(tx)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshaling transaction: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(txJSON)

	case http.MethodPost:
		// Create a new transaction
		var request struct {
			To      string `json:"to"`
			Amount  uint64 `json:"amount"`
			Data    string `json:"data"`
			ShardID uint64 `json:"shardId"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, fmt.Sprintf("Error decoding request: %v", err), http.StatusBadRequest)
			return
		}

		// Decode the recipient address
		to, err := hex.DecodeString(request.To)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error decoding recipient address: %v", err), http.StatusBadRequest)
			return
		}

		// Create the transaction
		tx, err := blockchainNode.CreateTransaction(to, request.Amount, []byte(request.Data), request.ShardID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating transaction: %v", err), http.StatusInternalServerError)
			return
		}

		// Return the transaction ID
		response := struct {
			ID string `json:"id"`
		}{
			ID: tx.ID.String(),
		}

		// Marshal the response to JSON
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshaling response: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get a block by hash
	blockHash := r.URL.Query().Get("hash")
	if blockHash == "" {
		http.Error(w, "Missing block hash", http.StatusBadRequest)
		return
	}

	block, err := blockchainNode.GetBlock(blockHash)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting block: %v", err), http.StatusNotFound)
		return
	}

	// Marshal the block to JSON
	blockJSON, err := json.Marshal(block)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshaling block: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(blockJSON)
}

func handleLatestBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the latest block
	block := blockchainNode.GetLatestBlock()
	if block == nil {
		http.Error(w, "No blocks in the blockchain", http.StatusNotFound)
		return
	}

	// Marshal the block to JSON
	blockJSON, err := json.Marshal(block)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshaling block: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(blockJSON)
}

func handleShard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get a shard state value
	shardID := r.URL.Query().Get("id")
	if shardID == "" {
		http.Error(w, "Missing shard ID", http.StatusBadRequest)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}

	value, err := blockchainNode.GetShardState(shardID, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting shard state: %v", err), http.StatusNotFound)
		return
	}

	// Marshal the value to JSON
	valueJSON, err := json.Marshal(value)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshaling value: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(valueJSON)
}

func handleProof(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Generate a proof
		shardID := r.URL.Query().Get("id")
		if shardID == "" {
			http.Error(w, "Missing shard ID", http.StatusBadRequest)
			return
		}

		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "Missing key", http.StatusBadRequest)
			return
		}

		proof, err := blockchainNode.GenerateShardProof(shardID, key)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error generating proof: %v", err), http.StatusInternalServerError)
			return
		}

		// Return the proof
		response := struct {
			Proof string `json:"proof"`
		}{
			Proof: hex.EncodeToString(proof),
		}

		// Marshal the response to JSON
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshaling response: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)

	case http.MethodPost:
		// Verify a proof
		var request struct {
			ShardID string `json:"shardId"`
			Key     string `json:"key"`
			Value   string `json:"value"`
			Proof   string `json:"proof"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, fmt.Sprintf("Error decoding request: %v", err), http.StatusBadRequest)
			return
		}

		// Decode the proof
		proof, err := hex.DecodeString(request.Proof)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error decoding proof: %v", err), http.StatusBadRequest)
			return
		}

		// Verify the proof
		valid, err := blockchainNode.VerifyShardProof(request.ShardID, request.Key, request.Value, proof)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error verifying proof: %v", err), http.StatusInternalServerError)
			return
		}

		// Return the verification result
		response := struct {
			Valid bool `json:"valid"`
		}{
			Valid: valid,
		}

		// Marshal the response to JSON
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshaling response: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the node status
	status := struct {
		Status               string  `json:"status"`
		Uptime               string  `json:"uptime"`
		PeerCount            int     `json:"peerCount"`
		ConnectedPeerCount   int     `json:"connectedPeerCount"`
		ConsistencyLevel     string  `json:"consistencyLevel"`
		NetworkState         string  `json:"networkState"`
		PartitionProbability float64 `json:"partitionProbability"`
	}{
		Status:               string(blockchainNode.GetStatus()),
		Uptime:               blockchainNode.GetUptime().String(),
		PeerCount:            blockchainNode.GetPeerCount(),
		ConnectedPeerCount:   blockchainNode.GetConnectedPeerCount(),
		ConsistencyLevel:     string(blockchainNode.GetConsistencyLevel()),
		NetworkState:         string(blockchainNode.GetNetworkState()),
		PartitionProbability: blockchainNode.GetPartitionProbability(),
	}

	// Marshal the status to JSON
	statusJSON, err := json.Marshal(status)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshaling status: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(statusJSON)
}
