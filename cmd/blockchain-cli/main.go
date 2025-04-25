package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/advanced-blockchain/pkg/crypto"
	"github.com/advanced-blockchain/pkg/state"
)

func main() {
	// Parse command line flags
	nodeAddr := flag.String("node", "127.0.0.1:8000", "Node address")
	keyFile := flag.String("key", "", "Private key file")
	flag.Parse()

	// Check the command
	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	// Load the key pair if provided
	var keyPair *crypto.KeyPair
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
	}

	// Process the command
	command := args[0]
	switch command {
	case "generate-key":
		generateKey(args[1:])
	case "get-transaction":
		getTransaction(*nodeAddr, args[1:])
	case "get-block":
		getBlock(*nodeAddr, args[1:])
	case "get-latest-block":
		getLatestBlock(*nodeAddr)
	case "create-transaction":
		createTransaction(*nodeAddr, keyPair, args[1:])
	case "get-shard-state":
		getShardState(*nodeAddr, args[1:])
	case "verify-proof":
		verifyProof(*nodeAddr, args[1:])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: blockchain-cli [options] <command> [args]")
	fmt.Println("Options:")
	fmt.Println("  -node <address>    Node address (default: 127.0.0.1:8000)")
	fmt.Println("  -key <file>        Private key file")
	fmt.Println("Commands:")
	fmt.Println("  generate-key <file>                                  Generate a new key pair and save it to a file")
	fmt.Println("  get-transaction <txid>                               Get a transaction by ID")
	fmt.Println("  get-block <hash>                                     Get a block by hash")
	fmt.Println("  get-latest-block                                     Get the latest block")
	fmt.Println("  create-transaction <to> <amount> <data> <shard-id>   Create a new transaction")
	fmt.Println("  get-shard-state <shard-id> <key>                     Get a shard state value")
	fmt.Println("  verify-proof <shard-id> <key> <value> <proof>        Verify a shard state proof")
}

func generateKey(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: blockchain-cli generate-key <file>")
		os.Exit(1)
	}

	// Generate a new key pair
	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		fmt.Printf("Error generating key pair: %v\n", err)
		os.Exit(1)
	}

	// Export the private key
	privateKey, err := keyPair.ExportPrivateKey()
	if err != nil {
		fmt.Printf("Error exporting private key: %v\n", err)
		os.Exit(1)
	}

	// Export the public key
	publicKey, err := keyPair.ExportPublicKey()
	if err != nil {
		fmt.Printf("Error exporting public key: %v\n", err)
		os.Exit(1)
	}

	// Save the private key to file
	if err := os.WriteFile(args[0], privateKey, 0600); err != nil {
		fmt.Printf("Error saving private key: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated new key pair\n")
	fmt.Printf("Private key saved to %s\n", args[0])
	fmt.Printf("Public key: %x\n", publicKey)
}

func getTransaction(nodeAddr string, args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: blockchain-cli get-transaction <txid>")
		os.Exit(1)
	}

	// Connect to the node
	conn, err := net.Dial("tcp", nodeAddr)
	if err != nil {
		fmt.Printf("Error connecting to node: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Create a request message
	message := fmt.Sprintf("GET_TRANSACTION %s\n", args[0])

	// Send the request
	if _, err := conn.Write([]byte(message)); err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		os.Exit(1)
	}

	// Receive the response
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error receiving response: %v\n", err)
		os.Exit(1)
	}

	// Parse the response
	response := string(buffer[:n])
	if strings.HasPrefix(response, "ERROR") {
		fmt.Printf("Error: %s\n", strings.TrimPrefix(response, "ERROR "))
		os.Exit(1)
	}

	// Print the transaction
	fmt.Println(response)
}

func getBlock(nodeAddr string, args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: blockchain-cli get-block <hash>")
		os.Exit(1)
	}

	// Connect to the node
	conn, err := net.Dial("tcp", nodeAddr)
	if err != nil {
		fmt.Printf("Error connecting to node: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Create a request message
	message := fmt.Sprintf("GET_BLOCK %s\n", args[0])

	// Send the request
	if _, err := conn.Write([]byte(message)); err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		os.Exit(1)
	}

	// Receive the response
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error receiving response: %v\n", err)
		os.Exit(1)
	}

	// Parse the response
	response := string(buffer[:n])
	if strings.HasPrefix(response, "ERROR") {
		fmt.Printf("Error: %s\n", strings.TrimPrefix(response, "ERROR "))
		os.Exit(1)
	}

	// Print the block
	fmt.Println(response)
}

func getLatestBlock(nodeAddr string) {
	// Connect to the node
	conn, err := net.Dial("tcp", nodeAddr)
	if err != nil {
		fmt.Printf("Error connecting to node: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Create a request message
	message := "GET_LATEST_BLOCK\n"

	// Send the request
	if _, err := conn.Write([]byte(message)); err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		os.Exit(1)
	}

	// Receive the response
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error receiving response: %v\n", err)
		os.Exit(1)
	}

	// Parse the response
	response := string(buffer[:n])
	if strings.HasPrefix(response, "ERROR") {
		fmt.Printf("Error: %s\n", strings.TrimPrefix(response, "ERROR "))
		os.Exit(1)
	}

	// Print the block
	fmt.Println(response)
}

func createTransaction(nodeAddr string, keyPair *crypto.KeyPair, args []string) {
	if len(args) != 4 {
		fmt.Println("Usage: blockchain-cli create-transaction <to> <amount> <data> <shard-id>")
		os.Exit(1)
	}

	// Check if we have a key pair
	if keyPair == nil {
		fmt.Println("Error: No private key provided")
		fmt.Println("Use the -key option to provide a private key file")
		os.Exit(1)
	}

	// Parse the arguments
	to, err := hex.DecodeString(args[0])
	if err != nil {
		fmt.Printf("Error parsing to address: %v\n", err)
		os.Exit(1)
	}

	amount, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		fmt.Printf("Error parsing amount: %v\n", err)
		os.Exit(1)
	}

	data := []byte(args[2])

	shardID, err := strconv.ParseUint(args[3], 10, 64)
	if err != nil {
		fmt.Printf("Error parsing shard ID: %v\n", err)
		os.Exit(1)
	}

	// Export the public key
	publicKey, err := keyPair.ExportPublicKey()
	if err != nil {
		fmt.Printf("Error exporting public key: %v\n", err)
		os.Exit(1)
	}

	// Create a new transaction
	tx := state.NewTransaction(publicKey, to, amount, data, shardID)

	// Sign the transaction
	if err := tx.Sign(keyPair); err != nil {
		fmt.Printf("Error signing transaction: %v\n", err)
		os.Exit(1)
	}

	// Connect to the node
	conn, err := net.Dial("tcp", nodeAddr)
	if err != nil {
		fmt.Printf("Error connecting to node: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Marshal the transaction
	txBytes, err := tx.Marshal()
	if err != nil {
		fmt.Printf("Error marshaling transaction: %v\n", err)
		os.Exit(1)
	}

	// Create a request message
	message := fmt.Sprintf("CREATE_TRANSACTION %x\n", txBytes)

	// Send the request
	if _, err := conn.Write([]byte(message)); err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		os.Exit(1)
	}

	// Receive the response
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error receiving response: %v\n", err)
		os.Exit(1)
	}

	// Parse the response
	response := string(buffer[:n])
	if strings.HasPrefix(response, "ERROR") {
		fmt.Printf("Error: %s\n", strings.TrimPrefix(response, "ERROR "))
		os.Exit(1)
	}

	// Print the transaction ID
	fmt.Printf("Transaction created: %s\n", response)
}

func getShardState(nodeAddr string, args []string) {
	if len(args) != 2 {
		fmt.Println("Usage: blockchain-cli get-shard-state <shard-id> <key>")
		os.Exit(1)
	}

	// Connect to the node
	conn, err := net.Dial("tcp", nodeAddr)
	if err != nil {
		fmt.Printf("Error connecting to node: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Create a request message
	message := fmt.Sprintf("GET_SHARD_STATE %s %s\n", args[0], args[1])

	// Send the request
	if _, err := conn.Write([]byte(message)); err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		os.Exit(1)
	}

	// Receive the response
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error receiving response: %v\n", err)
		os.Exit(1)
	}

	// Parse the response
	response := string(buffer[:n])
	if strings.HasPrefix(response, "ERROR") {
		fmt.Printf("Error: %s\n", strings.TrimPrefix(response, "ERROR "))
		os.Exit(1)
	}

	// Print the state value
	fmt.Println(response)
}

func verifyProof(nodeAddr string, args []string) {
	if len(args) != 4 {
		fmt.Println("Usage: blockchain-cli verify-proof <shard-id> <key> <value> <proof>")
		os.Exit(1)
	}

	// Connect to the node
	conn, err := net.Dial("tcp", nodeAddr)
	if err != nil {
		fmt.Printf("Error connecting to node: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Create a request message
	message := fmt.Sprintf("VERIFY_PROOF %s %s %s %s\n", args[0], args[1], args[2], args[3])

	// Send the request
	if _, err := conn.Write([]byte(message)); err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		os.Exit(1)
	}

	// Receive the response
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error receiving response: %v\n", err)
		os.Exit(1)
	}

	// Parse the response
	response := string(buffer[:n])
	if strings.HasPrefix(response, "ERROR") {
		fmt.Printf("Error: %s\n", strings.TrimPrefix(response, "ERROR "))
		os.Exit(1)
	}

	// Print the verification result
	fmt.Println(response)
}
