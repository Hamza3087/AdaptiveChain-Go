package main

import (
	"fmt"
	"log"
	"time"

	"github.com/advanced-blockchain/pkg/amf"
	"github.com/advanced-blockchain/pkg/consensus"
	"github.com/advanced-blockchain/pkg/crypto"
	"github.com/advanced-blockchain/pkg/network"
	"github.com/advanced-blockchain/pkg/node"
	"github.com/advanced-blockchain/pkg/state"
)

func main() {
	fmt.Println("Starting Advanced Blockchain System...")

	// Initialize the blockchain state
	blockchainState := state.NewState()
	fmt.Println("Blockchain state initialized")

	// Initialize the consensus mechanism
	hybridConsensus := consensus.NewHybridConsensus(blockchainState)
	fmt.Println("Hybrid consensus mechanism initialized")

	// Initialize the adaptive consistency orchestrator
	// We're not using this variable directly, but it's initialized for demonstration purposes
	_ = consensus.NewAdaptiveConsistency()
	fmt.Println("Adaptive consistency orchestrator initialized")

	// Initialize the conflict resolver
	// We're not using this variable directly, but it's initialized for demonstration purposes
	_ = consensus.NewConflictResolver(0.5)
	fmt.Println("Conflict resolver initialized")

	// Initialize the AMF
	adaptiveMerkleForest := amf.NewAdaptiveMerkleForest(100, 10*time.Second)
	fmt.Println("Adaptive Merkle Forest initialized")

	// Initialize the shard synchronizer
	shardSynchronizer := amf.NewShardSynchronizer(adaptiveMerkleForest)
	fmt.Println("Shard synchronizer initialized")

	// Create a key pair for the node
	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}
	fmt.Println("Key pair generated")

	// Export the public key
	publicKey, err := keyPair.ExportPublicKey()
	if err != nil {
		log.Fatalf("Failed to export public key: %v", err)
	}

	// Initialize the network
	// We're not using this variable directly, but it's initialized for demonstration purposes
	_ = network.NewNetwork("node-1", "localhost:8000", keyPair)
	fmt.Println("Network initialized")

	// Create node configuration
	nodeConfig := &node.NodeConfig{
		SplitThreshold:    100,
		RebalanceInterval: 10 * time.Second,
		PruningHeight:     100,
		ConsensusInterval: 5 * time.Second,
		TelemetryInterval: 1 * time.Second,
	}

	// Initialize the node
	node := node.NewNode("node-1", "localhost:8000", keyPair, nodeConfig)
	fmt.Println("Node initialized")

	// Register the node with the consensus mechanism
	err = hybridConsensus.AddValidator(node.ID, publicKey, 100)
	if err != nil {
		log.Fatalf("Failed to add validator: %v", err)
	}
	fmt.Println("Node registered as validator")

	// Create a state compressor
	stateCompressor := state.NewStateCompressor(blockchainState, state.CompressionAlgorithmGzip)
	fmt.Println("State compressor initialized")

	// Set the pruning height
	stateCompressor.SetPruningHeight(100)
	fmt.Println("Pruning height set to 100")

	// Create a root shard
	rootShard, err := adaptiveMerkleForest.GetShard("Root Shard")
	if err != nil {
		log.Fatalf("Failed to get root shard: %v", err)
	}
	fmt.Println("Root shard created")

	// Update the state of the root shard
	err = rootShard.UpdateState("key1", "value1")
	if err != nil {
		log.Fatalf("Failed to update state: %v", err)
	}
	fmt.Println("Root shard state updated")

	// Generate a proof for the state
	proof, err := rootShard.GenerateProof("key1")
	if err != nil {
		log.Fatalf("Failed to generate proof: %v", err)
	}
	fmt.Println("Proof generated")

	// Verify the proof
	valid, err := rootShard.VerifyProof("key1", "\"value1\"", proof)
	if err != nil {
		log.Fatalf("Failed to verify proof: %v", err)
	}
	if valid {
		fmt.Println("Proof verified successfully")
	} else {
		fmt.Println("Proof verification failed")
	}

	// Generate a compressed proof
	compressedProof, err := rootShard.GenerateCompressedProof("key1")
	if err != nil {
		log.Fatalf("Failed to generate compressed proof: %v", err)
	}
	fmt.Println("Compressed proof generated")

	// Verify the compressed proof
	valid, err = rootShard.VerifyCompressedProof("key1", "\"value1\"", compressedProof)
	if err != nil {
		log.Fatalf("Failed to verify compressed proof: %v", err)
	}
	if valid {
		fmt.Println("Compressed proof verified successfully")
	} else {
		fmt.Println("Compressed proof verification failed")
	}

	// Create a new shard
	shard2, err := adaptiveMerkleForest.GetShard("Shard-2")
	if err != nil {
		log.Fatalf("Failed to get shard: %v", err)
	}
	fmt.Println("Shard-2 created")

	// Update the state of the new shard
	err = shard2.UpdateState("key2", "value2")
	if err != nil {
		log.Fatalf("Failed to update state: %v", err)
	}
	fmt.Println("Shard-2 state updated")

	// Synchronize state between shards
	syncOpID, err := shardSynchronizer.SyncShards("Root Shard", "Shard-2", []string{"key1"})
	if err != nil {
		log.Fatalf("Failed to synchronize shards: %v", err)
	}
	fmt.Printf("Synchronization operation started with ID: %s\n", syncOpID)

	// Wait for the synchronization to complete
	time.Sleep(1 * time.Second)

	// Get the synchronization status
	syncOp, err := shardSynchronizer.GetSyncStatus(syncOpID)
	if err != nil {
		log.Fatalf("Failed to get synchronization status: %v", err)
	}
	fmt.Printf("Synchronization status: %s\n", syncOp.Status)

	// Create a commitment for the synchronization
	commitment, err := shardSynchronizer.CreateSyncCommitment(syncOpID, keyPair)
	if err != nil {
		log.Fatalf("Failed to create commitment: %v", err)
	}
	fmt.Println("Commitment created")

	// Verify the commitment
	valid, err = shardSynchronizer.VerifySyncCommitment(commitment, publicKey)
	if err != nil {
		log.Fatalf("Failed to verify commitment: %v", err)
	}
	if valid {
		fmt.Println("Commitment verified successfully")
	} else {
		fmt.Println("Commitment verification failed")
	}

	// Create a partial state transfer
	transfer, err := shardSynchronizer.CreatePartialStateTransfer(syncOpID, keyPair)
	if err != nil {
		log.Fatalf("Failed to create partial state transfer: %v", err)
	}
	fmt.Println("Partial state transfer created")

	// Apply the partial state transfer
	err = shardSynchronizer.ApplyPartialStateTransfer(transfer, publicKey)
	if err != nil {
		log.Fatalf("Failed to apply partial state transfer: %v", err)
	}
	fmt.Println("Partial state transfer applied")

	// Create an atomic cross-shard operation
	keys := map[string][]string{
		"Root Shard": {"key3"},
		"Shard-2":    {"key4"},
	}
	values := map[string]interface{}{
		"key3": "value3",
		"key4": "value4",
	}
	atomicOp, err := shardSynchronizer.CreateAtomicCrossShardOperation([]string{"Root Shard", "Shard-2"}, keys, values)
	if err != nil {
		log.Fatalf("Failed to create atomic cross-shard operation: %v", err)
	}
	fmt.Printf("Atomic cross-shard operation started with ID: %s\n", atomicOp.ID)

	// Wait for the atomic operation to complete
	time.Sleep(1 * time.Second)

	// Get the atomic operation status
	atomicOpStatus, err := shardSynchronizer.GetAtomicOperationStatus(atomicOp.ID)
	if err != nil {
		log.Fatalf("Failed to get atomic operation status: %v", err)
	}
	fmt.Printf("Atomic operation status: %s\n", atomicOpStatus.Status)

	// Start a consensus round
	err = hybridConsensus.StartConsensus()
	if err != nil {
		log.Fatalf("Failed to start consensus: %v", err)
	}
	fmt.Println("Consensus round started")

	// Submit prepare votes
	err = hybridConsensus.PrepareVote(node.ID, true)
	if err != nil {
		log.Fatalf("Failed to submit prepare vote: %v", err)
	}
	fmt.Println("Prepare vote submitted")

	// Submit commit votes
	err = hybridConsensus.CommitVote(node.ID, true)
	if err != nil {
		log.Fatalf("Failed to submit commit vote: %v", err)
	}
	fmt.Println("Commit vote submitted")

	// Get the consensus status
	consensusStatus := hybridConsensus.GetConsensusStatus()
	fmt.Printf("Consensus status: %s\n", consensusStatus)

	// Create a state archive
	archive, err := stateCompressor.CreateArchive()
	if err != nil {
		log.Fatalf("Failed to create state archive: %v", err)
	}
	fmt.Printf("State archive created with ID: %s\n", archive.ID)

	// Wait for the archive to be created
	time.Sleep(1 * time.Second)

	// Get the archive
	archive, err = stateCompressor.GetArchive(archive.ID)
	if err != nil {
		log.Fatalf("Failed to get state archive: %v", err)
	}
	fmt.Printf("Archive status: %s\n", archive.Status)

	// Create a compact state representation
	compact, err := stateCompressor.CreateCompactStateRepresentation()
	if err != nil {
		log.Fatalf("Failed to create compact state representation: %v", err)
	}
	fmt.Printf("Compact state representation created with height: %d\n", compact.Height)

	// Serialize the compact state representation
	compactData, err := stateCompressor.SerializeCompactStateRepresentation(compact)
	if err != nil {
		log.Fatalf("Failed to serialize compact state representation: %v", err)
	}
	fmt.Printf("Compact state representation serialized to %d bytes\n", len(compactData))

	// Deserialize the compact state representation
	compact2, err := stateCompressor.DeserializeCompactStateRepresentation(compactData)
	if err != nil {
		log.Fatalf("Failed to deserialize compact state representation: %v", err)
	}
	fmt.Printf("Compact state representation deserialized with height: %d\n", compact2.Height)

	// Prune the state
	err = stateCompressor.PruneState()
	if err != nil {
		log.Fatalf("Failed to prune state: %v", err)
	}
	fmt.Println("State pruned")

	// Compact the state
	err = stateCompressor.CompactState()
	if err != nil {
		log.Fatalf("Failed to compact state: %v", err)
	}
	fmt.Println("State compacted")

	// Create a Bloom filter
	bloomFilter := crypto.NewBloomFilter(1000, 0.01)
	fmt.Println("Bloom filter created")

	// Add items to the Bloom filter
	bloomFilter.Add([]byte("item1"))
	bloomFilter.Add([]byte("item2"))
	fmt.Println("Items added to Bloom filter")

	// Check if items are in the Bloom filter
	if bloomFilter.Contains([]byte("item1")) {
		fmt.Println("item1 is in the Bloom filter")
	}
	if bloomFilter.Contains([]byte("item3")) {
		fmt.Println("item3 is in the Bloom filter (false positive)")
	} else {
		fmt.Println("item3 is not in the Bloom filter")
	}

	// Create a Cuckoo filter
	cuckooFilter := crypto.NewCuckooFilter(1000, 0.01)
	fmt.Println("Cuckoo filter created")

	// Add items to the Cuckoo filter
	cuckooFilter.Add([]byte("item1"))
	cuckooFilter.Add([]byte("item2"))
	fmt.Println("Items added to Cuckoo filter")

	// Check if items are in the Cuckoo filter
	if cuckooFilter.Contains([]byte("item1")) {
		fmt.Println("item1 is in the Cuckoo filter")
	}
	if cuckooFilter.Contains([]byte("item3")) {
		fmt.Println("item3 is in the Cuckoo filter (false positive)")
	} else {
		fmt.Println("item3 is not in the Cuckoo filter")
	}

	// Create a homomorphic authenticated data structure
	hads, err := crypto.NewHomomorphicAuthenticatedDataStructure()
	if err != nil {
		log.Fatalf("Failed to create homomorphic authenticated data structure: %v", err)
	}
	fmt.Println("Homomorphic authenticated data structure created")

	// Add items to the homomorphic authenticated data structure
	err = hads.Add("key1", []byte("value1"))
	if err != nil {
		log.Fatalf("Failed to add item to homomorphic authenticated data structure: %v", err)
	}
	err = hads.Add("key2", []byte("value2"))
	if err != nil {
		log.Fatalf("Failed to add item to homomorphic authenticated data structure: %v", err)
	}
	fmt.Println("Items added to homomorphic authenticated data structure")

	// Generate proofs for the items
	proof1, err := hads.GenerateProof("key1")
	if err != nil {
		log.Fatalf("Failed to generate proof: %v", err)
	}
	proof2, err := hads.GenerateProof("key2")
	if err != nil {
		log.Fatalf("Failed to generate proof: %v", err)
	}
	fmt.Println("Proofs generated")

	// Verify the proofs
	valid, err = hads.VerifyProof(proof1)
	if err != nil {
		log.Fatalf("Failed to verify proof: %v", err)
	}
	if valid {
		fmt.Println("Proof 1 verified successfully")
	} else {
		fmt.Println("Proof 1 verification failed")
	}
	valid, err = hads.VerifyProof(proof2)
	if err != nil {
		log.Fatalf("Failed to verify proof: %v", err)
	}
	if valid {
		fmt.Println("Proof 2 verified successfully")
	} else {
		fmt.Println("Proof 2 verification failed")
	}

	// Merge the proofs
	mergedProof, err := hads.MergeProofs([][]byte{proof1, proof2})
	if err != nil {
		log.Fatalf("Failed to merge proofs: %v", err)
	}
	fmt.Println("Proofs merged")

	// Verify the merged proof
	valid, err = hads.VerifyMergedProof(mergedProof)
	if err != nil {
		log.Fatalf("Failed to verify merged proof: %v", err)
	}
	if valid {
		fmt.Println("Merged proof verified successfully")
	} else {
		fmt.Println("Merged proof verification failed")
	}

	fmt.Println("Advanced Blockchain System demonstration completed successfully!")
}
