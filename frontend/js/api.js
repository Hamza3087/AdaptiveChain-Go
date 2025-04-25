/**
 * API Service for Advanced Blockchain Frontend
 * Uses the BlockchainService for data instead of making HTTP requests
 */

// Ensure BlockchainService is initialized
if (!BlockchainService.initialized) {
    console.log('API: Initializing BlockchainService');
    BlockchainService.initialize();
}

// Dashboard API functions
const DashboardAPI = {
    // Get blockchain stats (blocks, transactions, nodes, network health)
    async getStats() {
        return BlockchainService.getStats();
    },

    // Get latest transactions
    async getLatestTransactions(limit = 5) {
        try {
            const result = BlockchainService.getTransactions(1, limit);
            return result.transactions;
        } catch (error) {
            console.error('Error fetching latest transactions:', error);
            throw error;
        }
    },

    // Get latest blocks for blockchain visualization
    async getLatestBlocks(limit = 5) {
        try {
            return BlockchainService.getLatestBlocks(limit);
        } catch (error) {
            console.error('Error fetching latest blocks:', error);
            throw error;
        }
    },

    // Get network activity data for charts
    async getNetworkActivity(period = 'day', chartType = 'transactions') {
        return BlockchainService.getNetworkActivity(period, chartType);
    }
};

// Blocks API functions
const BlocksAPI = {
    // Get blocks with pagination and filters
    async getBlocks(page = 1, limit = 10, filters = {}) {
        return BlockchainService.getBlocks(page, limit, filters);
    },

    // Get a specific block by number or hash
    async getBlock(blockId) {
        const block = BlockchainService.getBlock(blockId);
        if (!block) {
            throw new Error('Block not found');
        }
        return block;
    },

    // Get transactions in a block
    async getBlockTransactions(blockId) {
        return BlockchainService.getBlockTransactions(blockId);
    }
};

// Transactions API functions
const TransactionsAPI = {
    // Get transactions with pagination and filters
    async getTransactions(page = 1, limit = 10, filters = {}) {
        return BlockchainService.getTransactions(page, limit, filters);
    },

    // Get a specific transaction by hash
    async getTransaction(txHash) {
        const transaction = BlockchainService.getTransaction(txHash);
        if (!transaction) {
            throw new Error('Transaction not found');
        }
        return transaction;
    },

    // Get transaction status by hash
    async getTransactionStatus(txHash) {
        return BlockchainService.getTransactionStatus(txHash);
    },

    // Create a new transaction
    async createTransaction(transactionData) {
        return BlockchainService.createTransaction(transactionData);
    },

    // Process all pending transactions immediately
    async processAllPendingTransactions() {
        return BlockchainService.processAllPendingTransactions();
    }
};

// Nodes API functions
const NodesAPI = {
    // Get all nodes with status
    async getNodes() {
        return BlockchainService.getNodes();
    },

    // Get network stats
    async getNetworkStats() {
        return BlockchainService.getNetworkStats();
    },

    // Get a specific node by ID
    async getNode(nodeId) {
        return BlockchainService.getNode(nodeId);
    },

    // Add a new node
    async addNode(nodeData) {
        return BlockchainService.addNode(nodeData);
    },

    // Restart a node
    async restartNode(nodeId) {
        return BlockchainService.restartNode(nodeId);
    },

    // Set a node offline
    async setNodeOffline(nodeId) {
        return BlockchainService.setNodeOffline(nodeId);
    },

    // Remove a node
    async removeNode(nodeId) {
        return BlockchainService.removeNode(nodeId);
    }
};

// AMF API functions
const AMFAPI = {
    // Get AMF stats
    async getStats() {
        // Force update AMF stats before returning them
        if (BlockchainService.updateAMFStats) {
            BlockchainService.updateAMFStats();
        }
        return BlockchainService.getAMFStats();
    },

    // Get AMF forest structure
    async getForestStructure() {
        return BlockchainService.getForestStructure();
    },

    // Get a specific shard by ID
    async getShard(shardId) {
        return BlockchainService.getShard(shardId);
    },

    // Get operations for a specific shard
    async getShardOperations(shardId) {
        return BlockchainService.getShardOperations(shardId);
    },

    // Get all keys for a specific shard
    async getShardKeys(shardId) {
        return BlockchainService.getShardKeys(shardId);
    },

    // Get proof for a key in a shard
    async getProof(shardId, key) {
        return BlockchainService.getProof(shardId, key);
    },

    // Verify a proof
    async verifyProof(proofData) {
        // Validate input
        if (!proofData || !proofData.shardId || !proofData.key || !proofData.proof) {
            return {
                valid: false,
                error: "Invalid proof data: missing required fields"
            };
        }

        try {
            // Call the blockchain service to verify the proof
            const result = await BlockchainService.verifyProof(proofData);

            // Add additional verification metadata
            return {
                ...result,
                verificationTime: new Date().toISOString(),
                verificationMethod: "AMF Merkle Path Verification",
                proofType: "Merkle"
            };
        } catch (error) {
            console.error("Error in proof verification:", error);
            return {
                valid: false,
                error: error.message || "Unknown verification error",
                details: {
                    errorType: error.name || "Error",
                    timestamp: new Date().toISOString()
                }
            };
        }
    },

    // Force update AMF stats and create new shards if needed
    async updateAMF() {
        // Create a new shard if we have fewer than expected
        const shardCount = Object.keys(BlockchainService.blockchain.shards)
            .filter(id => id !== 'Root Shard').length;

        // Calculate ideal shard count based on transaction volume
        const txCount = BlockchainService.stats.transactions;
        const idealShardCount = Math.min(10, Math.max(2, Math.floor(txCount / 10) + 1));

        // Create new shards if needed
        if (shardCount < idealShardCount) {
            for (let i = shardCount; i < idealShardCount; i++) {
                BlockchainService.createNewShard();
            }
        }

        // Update AMF stats
        if (BlockchainService.updateAMFStats) {
            BlockchainService.updateAMFStats();
        }

        return BlockchainService.getAMFStats();
    }
};

// Status API functions
const StatusAPI = {
    // Get blockchain node status
    async getStatus() {
        try {
            // Ensure BlockchainService is initialized
            if (!BlockchainService.initialized) {
                console.log('StatusAPI: Initializing BlockchainService');
                BlockchainService.initialize();
            }

            // Get status from BlockchainService
            return BlockchainService.getStatus();
        } catch (error) {
            console.error('Error getting blockchain status:', error);

            // Return a default status if there's an error
            return {
                status: 'running',
                version: '1.0.0',
                uptime: 0,
                peers: 0,
                syncStatus: 'Unknown',
                pendingTransactions: 0,
                error: error.message
            };
        }
    },

    // Process all pending transactions immediately
    async processAllPendingTransactions() {
        try {
            return BlockchainService.processAllPendingTransactions();
        } catch (error) {
            console.error('Error processing pending transactions:', error);
            return {
                success: false,
                message: error.message
            };
        }
    }
};

// Export all API functions
const API = {
    Dashboard: DashboardAPI,
    Blocks: BlocksAPI,
    Transactions: TransactionsAPI,
    Nodes: NodesAPI,
    AMF: AMFAPI,
    Status: StatusAPI
};
