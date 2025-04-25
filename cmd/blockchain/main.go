package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/advanced-blockchain/pkg/crypto"
	"github.com/advanced-blockchain/pkg/node"
)

func main() {
	// Parse command line flags
	nodeID := flag.String("id", "node-1", "Node ID")
	address := flag.String("addr", "127.0.0.1:8000", "Node address")
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
	node := node.NewNode(*nodeID, *address, keyPair, config)
	
	// Start the node
	if err := node.Start(); err != nil {
		fmt.Printf("Error starting node: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Node started with ID %s at address %s\n", *nodeID, *address)
	
	// Export the public key
	publicKey, err := keyPair.ExportPublicKey()
	if err != nil {
		fmt.Printf("Error exporting public key: %v\n", err)
	} else {
		fmt.Printf("Public key: %x\n", publicKey)
	}
	
	// Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Wait for a signal
	sig := <-sigChan
	fmt.Printf("Received signal %s, shutting down...\n", sig)
	
	// Stop the node
	if err := node.Stop(); err != nil {
		fmt.Printf("Error stopping node: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Node stopped")
}
