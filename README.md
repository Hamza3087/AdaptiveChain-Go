# AdaptiveChain-Go

A sophisticated blockchain system in Go that demonstrates advanced state verification, dynamic sharding, and Merkle proof verification through an interactive web interface.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Run Instructions](#run-instructions)
- [Usage Guide](#usage-guide)
- [Merkle Proof Verification](#merkle-proof-verification)
- [Technical Details](#technical-details)
- [Development](#development)

## Overview

The Adaptive Merkle Forest (AMF) is an advanced data structure that combines the benefits of Merkle trees with sharding to create a scalable and efficient blockchain system. This simulation demonstrates how AMF works, how shards are created and managed, and how Merkle proofs are generated and verified.

## Features

### Adaptive Merkle Forest (AMF) with Advanced Complexity

- **Interactive Blockchain Dashboard**: Real-time statistics and visualization of the blockchain
- **Dynamic Block Creation**: Blocks are created organically as transactions grow
- **Adaptive Sharding**: Dynamic shard generation based on transaction volume
- **Merkle Proof Generation & Verification**: Create and verify cryptographic proofs for data in shards
- **Cross-Shard Operations**: Simulate operations that span multiple shards
- **Network Simulation**: Visualize nodes and network status
- **Transaction Pool**: See pending transactions and how they're processed into blocks

### Blockchain Data Structure Innovations

- **Advanced Block Composition**: Extended blockchain block structure with cryptographic accumulators for compact state representation
- **State Compression and Archival**: Advanced state management with state pruning algorithms and efficient archival mechanisms

## Architecture

The project consists of two main components:

1. **Backend API (Go)**:

   - RESTful API for blockchain operations
   - Implements the core blockchain logic
   - Manages shards, blocks, and transactions
   - Handles proof generation and verification

2. **Frontend (HTML/CSS/JavaScript)**:
   - Interactive dashboard for visualizing the blockchain
   - AMF visualization with shard management
   - Proof verification interface
   - Block and transaction explorer

The system is organized into the following packages:

- `pkg/amf`: Adaptive Merkle Forest implementation
- `pkg/blockchain`: Core blockchain logic
- `pkg/crypto`: Cryptographic utilities
- `pkg/network`: Network simulation
- `cmd/blockchain-api`: API server

## Run Instructions

### Prerequisites

- Go 1.16 or higher (for the backend API)
- A modern web browser (Chrome, Firefox, Edge, etc.)

### Running the Backend API

1. Navigate to the blockchain-api directory:

   ```
   cd cmd/blockchain-api
   ```

2. Run the Go application:

   ```
   go run main.go
   ```

   This will start the API server on `http://localhost:8080`.

### Running the Frontend

There are two ways to run the frontend:

#### Option 1: Using the built-in API server

If the backend API is running, you can access the frontend through it:

1. Open your browser and navigate to:
   ```
   http://localhost:8080
   ```

#### Option 2: Opening the HTML files directly

You can also open the HTML files directly in your browser:

1. Navigate to the frontend directory
2. Open any of the HTML files in your browser:
   - `index.html` - Dashboard
   - `blocks.html` - Block explorer
   - `transactions.html` - Transaction explorer
   - `amf.html` - AMF visualization and proof verification
   - `nodes.html` - Node management

## Usage Guide

### Dashboard

The dashboard provides an overview of the blockchain system:

- **Blockchain Statistics**: Total blocks, transactions, and nodes
- **Recent Blocks**: Latest blocks added to the blockchain
- **Transaction Pool**: Pending transactions waiting to be included in blocks
- **Network Status**: Status of the blockchain network and nodes

### Blocks & Transactions

The blocks and transactions pages allow you to explore the blockchain data:

- **Block Explorer**: View details of each block including hash, timestamp, and transactions
- **Transaction Explorer**: View details of each transaction including sender, receiver, and amount
- **Search**: Search for specific blocks or transactions by hash or ID

### AMF Visualization

The AMF page provides a visualization of the Adaptive Merkle Forest:

- **Forest Structure**: Visual representation of shards and their relationships
- **Shard Management**: Create new shards and view shard details
- **Proof Verification**: Generate and verify Merkle proofs for data in shards
- **Shard Keys**: View the keys stored in each shard

### Nodes

The nodes page allows you to manage the blockchain nodes:

- **Node List**: View all nodes in the network
- **Node Management**: Add, restart, or remove nodes
- **Network Statistics**: View network performance metrics

## Merkle Proof Verification

### How Proofs Work

Merkle proofs are cryptographic proofs that allow you to verify that a piece of data is included in a Merkle tree without having to store the entire tree. In the AMF system:

1. Each shard maintains a Merkle tree of its key-value pairs
2. A proof consists of:

   - The key and value being proven
   - The leaf hash (hash of the key-value pair)
   - A set of sibling hashes (the Merkle path)
   - The position of the key in the tree
   - The root hash of the Merkle tree

3. To verify a proof, the system:
   - Checks if the key exists in the shard
   - Verifies that the value matches the expected value
   - Recomputes the Merkle root using the leaf hash and siblings
   - Compares the computed root with the expected root hash

### Verifying Proofs

To verify a proof:

1. Go to the AMF page
2. Click on a shard to view its details
3. Click the "Verify Proof" button in the visualization controls
4. Enter a shard ID and key
5. Click "Verify Proof" to see the verification result

The verification result will show:

- Whether the proof is valid
- Details about the key and value
- Information about each verification step
- Any errors that occurred during verification

### Finding Valid Keys

To find valid keys for proof verification:

1. Go to the AMF page
2. Click on a shard to view its details
3. Click the "View Keys" button in the Shard Keys section
4. Browse the list of keys in the shard
5. Click on a key to automatically fill it in the proof verification form
6. Click "Verify Proof" to verify the proof

## Technical Details

### Sharding

The AMF system uses sharding to distribute data across multiple shards:

- Each shard is responsible for a subset of the key-value pairs
- Shards are organized in a hierarchical structure (the forest)
- New shards are created dynamically as the transaction volume grows
- Cross-shard operations are handled through a coordination protocol

### Merkle Trees

Each shard maintains a Merkle tree of its key-value pairs:

- Leaf nodes contain the hash of a key-value pair
- Internal nodes contain the hash of their children
- The root hash represents the entire state of the shard
- Changes to the shard state result in a new root hash

### Proof Generation

Proofs are generated as follows:

1. The system locates the key-value pair in the appropriate shard
2. It computes the leaf hash of the key-value pair
3. It identifies the position of the leaf in the Merkle tree
4. It collects the sibling hashes along the path from the leaf to the root
5. It packages the key, value, leaf hash, siblings, and root hash into a proof

### Proof Verification

Proofs are verified as follows:

1. The system checks if the key exists in the specified shard
2. It verifies that the value matches the expected value for the key
3. It recomputes the Merkle root using the leaf hash and siblings
4. It compares the computed root with the expected root hash
5. If all checks pass, the proof is considered valid

## Development

### Project Structure

```
/
├── cmd/
│   └── blockchain-api/     # Backend API
│       └── main.go         # API entry point
├── pkg/
│   ├── amf/                # AMF implementation
│   ├── blockchain/         # Core blockchain logic
│   ├── crypto/             # Cryptographic utilities
│   └── network/            # Network simulation
├── frontend/
│   ├── css/                # Stylesheets
│   ├── js/                 # JavaScript files
│   │   ├── api.js          # API client
│   │   ├── blockchain-service.js # Frontend blockchain service
│   │   └── utils.js        # Utility functions
│   ├── index.html          # Dashboard
│   ├── blocks.html         # Block explorer
│   ├── transactions.html   # Transaction explorer
│   ├── amf.html            # AMF visualization
│   └── nodes.html          # Node management
└── README.md               # This file
```

### Backend API

The backend API provides the following endpoints:

- `/api/status` - Get blockchain status
- `/api/blocks` - Get blocks
- `/api/transactions` - Get transactions
- `/api/amf/stats` - Get AMF statistics
- `/api/amf/forest` - Get AMF forest structure
- `/api/amf/shard/:id` - Get shard details
- `/api/amf/proof` - Generate and verify proofs
- `/api/nodes` - Manage nodes

### Frontend

The frontend is built with vanilla HTML, CSS, and JavaScript:

- **HTML**: Structure of the pages
- **CSS**: Styling and animations
- **JavaScript**: Interactivity and API communication

The main JavaScript files are:

- **api.js**: Client for communicating with the backend API
- **blockchain-service.js**: Frontend service for blockchain operations
- **utils.js**: Utility functions for formatting, validation, etc.

## Troubleshooting

### Proof Verification Issues

If you encounter issues with proof verification:

1. **Use the "View Keys" Feature**: Always use the "View Keys" button to find valid keys for a shard
2. **Check Key Format**: Make sure the key format is correct (keys should start with "key-")
3. **Verify Shard ID**: Ensure you're using the correct shard ID
4. **Check Error Messages**: The verification result will show detailed error messages

### Common Errors

- **"Key not found in shard"**: The key doesn't exist in the specified shard
- **"Value mismatch"**: The value doesn't match the expected value for the key
- **"Merkle path verification failed"**: The Merkle path doesn't compute to the expected root hash

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgements

This project is a simulation for educational purposes and is not intended for production use. It demonstrates blockchain concepts and the Adaptive Merkle Forest data structure in an interactive way.
