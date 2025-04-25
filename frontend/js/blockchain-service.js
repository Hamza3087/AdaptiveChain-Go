/**
 * Blockchain Service
 * Simulates a blockchain backend with dynamic data generation
 */

// Dynamic data store - simulates a runtime database
const BlockchainService = {
  // Blockchain state
  blockchain: {
    startTime: Date.now(),
    lastBlockTime: Date.now(),
    blockInterval: 5000, // 5 seconds per block (faster for demo purposes)
    transactionPool: [],
    blocks: [],
    transactions: [],
    nodes: [],
    shards: {},
    forestStructure: { shards: [] }
  },

  // Runtime stats
  stats: {
    blocks: 0,
    transactions: 0,
    nodes: 5,
    networkHealth: 98.7
  },

  // Network stats
  networkStats: {
    totalNodes: 5,
    activeNodes: 5,
    networkLatency: 125,
    byzantineNodes: 0
  },

  // AMF stats
  amfStats: {
    totalShards: 0,
    activeShards: 0,
    shardDepth: 0,
    rebalanceOperations: 0
  },

  // Initialize blockchain data
  initialize() {
    console.log('Initializing blockchain service...');

    // Check if we already have shards - if so, don't recreate the AMF structure
    const hasExistingShards = Object.keys(this.blockchain.shards).length > 0;
    console.log(`Has existing shards: ${hasExistingShards}, count: ${Object.keys(this.blockchain.shards).length}`);

    // Create initial nodes
    this.createInitialNodes();

    // Create initial AMF structure only if we don't have any shards yet
    if (!hasExistingShards) {
      console.log('Creating initial AMF structure');
      this.createInitialAMFStructure();
    } else {
      console.log('Using existing AMF structure');
      // Make sure AMF stats are updated based on the current state
      this.updateAMFStats();
    }

    // Create genesis block if we don't have any blocks yet
    if (this.blockchain.blocks.length === 0) {
      this.createGenesisBlock();
    }

    // Generate a few initial transactions to get started
    if (this.blockchain.transactions.length === 0 && this.blockchain.transactionPool.length === 0) {
      for (let i = 0; i < 5; i++) {
        this.generateRandomTransaction();
      }
    }

    // Start block production
    setInterval(() => this.produceBlock(), this.blockchain.blockInterval);

    // Start random transaction generation - more frequent for demo purposes
    setInterval(() => {
      // Generate 1-3 transactions at a time
      const numTransactions = Math.floor(Math.random() * 3) + 1;
      for (let i = 0; i < numTransactions; i++) {
        this.generateRandomTransaction();
      }
    }, 3000); // Every 3 seconds

    // Force process pending transactions if they've been waiting too long
    setInterval(() => {
      this.forceProcessPendingTransactions();
    }, 5000); // Every 5 seconds (reduced from 15 seconds)

    // Update node metrics periodically
    setInterval(() => this.updateNodeMetrics(), 10000);

    // Update AMF stats periodically
    setInterval(() => {
      this.updateAMFStats();
      this.saveState();
    }, 5000); // Every 5 seconds

    console.log('Blockchain service initialized');
  },

  // Force process pending transactions that have been waiting too long
  forceProcessPendingTransactions() {
    // Check if there are pending transactions
    if (this.blockchain.transactionPool.length === 0) {
      return;
    }

    // Find transactions that have been waiting for more than 15 seconds (reduced from 30)
    const now = Date.now();
    const oldTransactions = this.blockchain.transactionPool.filter(tx =>
      now - tx.timestamp > 15000 // 15 seconds
    );

    if (oldTransactions.length > 0) {
      console.log(`Force processing ${oldTransactions.length} pending transactions that have been waiting too long`);

      // Force produce blocks with all pending transactions
      this.produceBlock(true); // Pass true to process ALL pending transactions
    }
  },

  // Create initial nodes
  createInitialNodes() {
    const nodeAddresses = ['127.0.0.1:8000', '127.0.0.1:8001', '127.0.0.1:8002', '127.0.0.1:8003', '127.0.0.1:8004'];

    nodeAddresses.forEach((address, index) => {
      const nodeId = `node-${index + 1}`;
      const publicKey = this.generateRandomHash();

      this.blockchain.nodes.push({
        id: nodeId,
        address: address,
        status: 'Online',
        uptime: Math.floor(Math.random() * 500000),
        peers: Math.floor(Math.random() * 20) + 10,
        consistencyLevel: 'Strong',
        publicKey: publicKey,
        version: '1.0.0',
        networkState: 'Normal',
        metrics: {
          cpu: Math.floor(Math.random() * 50) + 10,
          memory: Math.floor(Math.random() * 500000000) + 500000000,
          disk: Math.floor(Math.random() * 2000000000) + 7000000000,
          networkIn: (Math.random() * 3).toFixed(1),
          networkOut: (Math.random() * 2).toFixed(1),
          blocks: 0,
          transactions: 0
        }
      });
    });
  },

  // Create initial AMF structure
  createInitialAMFStructure() {
    // Create root shard with enhanced data structure
    const rootShard = {
      id: 'Root Shard',
      parentId: null,
      childIds: [],
      creationTime: Date.now(),
      keys: 0,
      hash: this.generateRandomHash(),
      lastUpdated: Date.now(),
      status: 'Active',
      // Add statistics tracking
      stats: {
        totalTransactions: 0,
        totalValue: 0,
        averageValue: 0,
        uniqueAddresses: new Set(),
        blockCount: 0,
        lastBlockNumber: 0
      },
      operations: [
        { type: 'Creation', timestamp: Date.now() }
      ]
    };

    this.blockchain.shards['Root Shard'] = rootShard;
    this.blockchain.forestStructure.shards.push({
      id: 'Root Shard',
      level: 0,
      keys: 0,
      hash: rootShard.hash
    });

    // Create only Shard-1 initially (main shard) with enhanced data structure
    const mainShardId = 'Shard-1';
    const mainShard = {
      id: mainShardId,
      parentId: 'Root Shard',
      childIds: [],
      creationTime: Date.now(),
      keys: 0,
      hash: this.generateRandomHash(),
      lastUpdated: Date.now(),
      status: 'Active',
      // Add statistics tracking
      stats: {
        totalTransactions: 0,
        totalValue: 0,
        averageValue: 0,
        uniqueAddresses: new Set(),
        blockCount: 0,
        lastBlockNumber: 0
      },
      operations: [
        { type: 'Creation', timestamp: Date.now() }
      ]
    };

    this.blockchain.shards[mainShardId] = mainShard;
    this.blockchain.forestStructure.shards.push({
      id: mainShardId,
      level: 1,
      keys: 0,
      hash: mainShard.hash
    });

    // Update root shard to reference the main shard
    rootShard.childIds.push(mainShardId);

    // Update AMF stats
    this.amfStats.totalShards = 2; // Root + Main shard
    this.amfStats.activeShards = 2;
    this.amfStats.shardDepth = 1;

    console.log('Created initial AMF structure with Root Shard and Shard-1');
  },

  // Create genesis block
  createGenesisBlock() {
    const genesisBlock = {
      number: 1,
      hash: this.generateRandomHash(),
      previousHash: '0x0000000000000000000000000000000000000000000000000000000000000000',
      timestamp: this.blockchain.startTime,
      transactions: 0,
      shardId: 1,
      size: 1024,
      merkleRoot: this.generateRandomHash(),
      stateRoot: this.generateRandomHash(),
      difficulty: 1,
      nonce: 0
    };

    this.blockchain.blocks.push(genesisBlock);
    this.blockchain.lastBlockTime = this.blockchain.startTime;
    this.stats.blocks = 1;

    // Update node metrics
    this.blockchain.nodes.forEach(node => {
      node.metrics.blocks = 1;
    });

    // Update the shard data for the genesis block
    this.updateShardForNewBlock(1, genesisBlock, []);
  },

  // Produce a new block
  produceBlock(forceAll = false) {
    try {
      // Validate blockchain state
      if (!this.blockchain) {
        console.error('Blockchain is not initialized');
        return;
      }

      if (!this.blockchain.transactionPool) {
        console.error('Transaction pool is not initialized');
        this.blockchain.transactionPool = [];
        return;
      }

      if (!this.blockchain.blocks || this.blockchain.blocks.length === 0) {
        console.error('Blocks array is not initialized or empty');
        return;
      }

      // Only create blocks when there are pending transactions
      if (this.blockchain.transactionPool.length === 0) {
        // No pending transactions, skip block production
        return;
      }

      console.log(`Producing blocks for ${this.blockchain.transactionPool.length} pending transactions (forceAll: ${forceAll})`);

      // Group pending transactions by shard
      const transactionsByShard = {};
      this.blockchain.transactionPool.forEach(tx => {
        // Ensure transaction has a valid shardId
        const shardId = tx.shardId || 1; // Default to shard 1 if not specified

        if (!transactionsByShard[shardId]) {
          transactionsByShard[shardId] = [];
        }
        transactionsByShard[shardId].push(tx);
      });

      // Process each shard's transactions
      for (const [shardId, shardTransactions] of Object.entries(transactionsByShard)) {
        try {
          // Skip if no transactions for this shard
          if (shardTransactions.length === 0) continue;

          // If forceAll is true, process all transactions for this shard
          // Otherwise, take up to 100 transactions per block (increased from 50)
          const maxTransactionsPerBlock = forceAll ? shardTransactions.length : 100;
          const pendingTransactions = shardTransactions.slice(0, maxTransactionsPerBlock);
          const transactionCount = pendingTransactions.length;

          // Remove these transactions from the pool
          pendingTransactions.forEach(tx => {
            const index = this.blockchain.transactionPool.findIndex(t => t.hash === tx.hash);
            if (index !== -1) {
              this.blockchain.transactionPool.splice(index, 1);
            }
          });

          // Create new block for this shard
          const lastBlock = this.blockchain.blocks[this.blockchain.blocks.length - 1];
          if (!lastBlock) {
            console.error('Last block is undefined');
            continue;
          }

          const newBlockNumber = lastBlock.number + 1;

          const newBlock = {
            number: newBlockNumber,
            hash: this.generateRandomHash(),
            previousHash: lastBlock.hash,
            timestamp: Date.now(),
            transactions: transactionCount,
            shardId: parseInt(shardId) || 1, // Use the shard ID from the transactions, default to 1
            size: 1024 + (transactionCount * 256),
            merkleRoot: this.generateRandomHash(),
            stateRoot: this.generateRandomHash(),
            difficulty: 1,
            nonce: Math.floor(Math.random() * 1000000)
          };

          // Add block to blockchain
          this.blockchain.blocks.push(newBlock);
          this.blockchain.lastBlockTime = newBlock.timestamp;
          this.stats.blocks++;

          // Update transactions
          pendingTransactions.forEach(tx => {
            tx.blockNumber = newBlockNumber;
            tx.status = 'Confirmed';
            this.blockchain.transactions.push(tx);
          });

          console.log(`Produced block #${newBlockNumber} with ${transactionCount} transactions for Shard ${shardId}`);

          // Update AMF for this shard
          try {
            this.updateShardForNewBlock(parseInt(shardId) || 1, newBlock, pendingTransactions);

            // Dynamically create new shards when blocks are added
            // This ensures the AMF grows as the blockchain grows
            if (this.stats.blocks % 3 === 0) { // Every 3 blocks
              // Check if we need more shards
              const shardCount = Object.keys(this.blockchain.shards).filter(id => id !== 'Root Shard').length;

              // Create a new shard if we have fewer than 10 shards
              if (shardCount < 10 && Math.random() < 0.5) {
                const newShardId = this.createNewShard();
                console.log(`Created new shard ${newShardId} after block #${newBlockNumber}`);
              }
            }
          } catch (shardError) {
            console.error(`Error updating shard ${shardId}:`, shardError);
          }
        } catch (shardProcessingError) {
          console.error(`Error processing transactions for shard ${shardId}:`, shardProcessingError);
        }
      }

      // If there are still pending transactions and forceAll is true, process them in another round
      if (forceAll && this.blockchain.transactionPool.length > 0) {
        console.log(`Still have ${this.blockchain.transactionPool.length} pending transactions, processing in another round`);
        this.produceBlock(true);
      }

      // Update node metrics
      if (this.blockchain.nodes) {
        this.blockchain.nodes.forEach(node => {
          if (node && node.metrics) {
            node.metrics.blocks = this.stats.blocks;
            node.metrics.transactions = this.stats.transactions;
          }
        });
      }

      // Update AMF stats directly based on the current state
      this.updateAMFStats();

      // Save state immediately after producing blocks
      this.saveState();
    } catch (error) {
      console.error('Error producing block:', error);
    }
  },

  // Update AMF stats based on current blockchain state
  updateAMFStats() {
    // Count active shards (excluding Root Shard)
    const activeShards = Object.values(this.blockchain.shards)
      .filter(shard => shard.id !== 'Root Shard' && shard.status === 'Active');

    // Update AMF stats
    this.amfStats.totalShards = Object.keys(this.blockchain.shards).length;
    this.amfStats.activeShards = activeShards.length;

    // Calculate shard depth
    const maxLevel = Math.max(
      ...this.blockchain.forestStructure.shards.map(shard => shard.level || 0)
    );
    this.amfStats.shardDepth = maxLevel;

    console.log(`Updated AMF stats: ${this.amfStats.activeShards} active shards, depth ${this.amfStats.shardDepth}`);
  },

  // Update shard data when a new block is added
  updateShardForNewBlock(shardId, block, transactions) {
    // Find the corresponding shard
    const shardKey = `Shard-${shardId}`;
    let shard = this.blockchain.shards[shardKey];

    // If shard doesn't exist yet, create it
    if (!shard) {
      // Find a parent shard to attach to
      let parentId = 'Root Shard';

      // Create the new shard with enhanced data structure
      shard = {
        id: shardKey,
        parentId: parentId,
        childIds: [],
        creationTime: Date.now(),
        keys: 0,
        hash: this.generateRandomHash(),
        lastUpdated: Date.now(),
        status: 'Active',
        // Add statistics tracking
        stats: {
          totalTransactions: 0,
          totalValue: 0,
          averageValue: 0,
          uniqueAddresses: new Set(),
          blockCount: 0,
          lastBlockNumber: 0
        },
        operations: [
          { type: 'Creation', timestamp: Date.now() }
        ]
      };

      // Add to shards collection
      this.blockchain.shards[shardKey] = shard;

      // Add to forest structure
      this.blockchain.forestStructure.shards.push({
        id: shardKey,
        level: 1,
        keys: 0,
        hash: shard.hash
      });

      // Update parent shard
      const parentShard = this.blockchain.shards[parentId];
      if (parentShard && !parentShard.childIds.includes(shardKey)) {
        parentShard.childIds.push(shardKey);
      }

      // Update AMF stats
      this.amfStats.totalShards++;
      this.amfStats.activeShards++;
    }

    // Update shard data
    shard.keys += transactions.length;
    shard.hash = block.hash;
    shard.lastUpdated = block.timestamp;

    // Ensure shard.stats exists
    if (!shard.stats) {
      console.log(`Creating stats object for shard ${shard.id}`);
      shard.stats = {
        totalTransactions: 0,
        totalValue: 0,
        averageValue: 0,
        uniqueAddresses: new Set(),
        blockCount: 0,
        lastBlockNumber: 0
      };
    }

    // Update shard statistics with real transaction data
    shard.stats.blockCount = (shard.stats.blockCount || 0) + 1;
    shard.stats.lastBlockNumber = block.number;

    if (transactions.length > 0) {
      // Update transaction count
      shard.stats.totalTransactions = (shard.stats.totalTransactions || 0) + transactions.length;

      // Update value statistics
      const totalValue = transactions.reduce((sum, tx) => sum + tx.amount, 0);
      shard.stats.totalValue = (shard.stats.totalValue || 0) + totalValue;
      shard.stats.averageValue = shard.stats.totalValue / shard.stats.totalTransactions;

      // Ensure uniqueAddresses is a Set
      if (!shard.stats.uniqueAddresses || !(shard.stats.uniqueAddresses instanceof Set)) {
        shard.stats.uniqueAddresses = new Set();
      }

      // Update unique addresses
      transactions.forEach(tx => {
        shard.stats.uniqueAddresses.add(tx.from);
        shard.stats.uniqueAddresses.add(tx.to);
      });
    }

    // Ensure shard.operations exists
    if (!shard.operations) {
      console.log(`Creating operations array for shard ${shard.id}`);
      shard.operations = [];
    }

    // Add operation with more detailed transaction data
    const operation = {
      type: 'Block Added',
      timestamp: block.timestamp,
      blockNumber: block.number,
      transactionCount: transactions.length,
      blockHash: block.hash
    };

    // Add transaction details if there are any
    if (transactions.length > 0) {
      try {
        // Calculate total transaction value
        const totalValue = transactions.reduce((sum, tx) => sum + (tx.amount || 0), 0);

        // Add transaction summary
        operation.totalValue = totalValue;
        operation.averageValue = totalValue / transactions.length;

        // Add transaction hashes (limited to first 5 for space)
        operation.transactionHashes = transactions.slice(0, 5).map(tx => tx.hash);

        // Track unique addresses
        const addresses = new Set();
        transactions.forEach(tx => {
          if (tx.from) addresses.add(tx.from);
          if (tx.to) addresses.add(tx.to);
        });
        operation.uniqueAddresses = addresses.size;
      } catch (error) {
        console.error(`Error processing transaction details for shard ${shard.id}:`, error);
        // Add basic operation without detailed transaction data
        operation.error = "Failed to process transaction details";
      }
    }

    shard.operations.push(operation);

    // Ensure forest structure exists
    if (!this.blockchain.forestStructure) {
      console.log('Creating forest structure');
      this.blockchain.forestStructure = {
        shards: []
      };
    }

    if (!this.blockchain.forestStructure.shards) {
      console.log('Creating shards array in forest structure');
      this.blockchain.forestStructure.shards = [];
    }

    // Update forest structure
    let forestShard = this.blockchain.forestStructure.shards.find(s => s && s.id === shard.id);
    if (forestShard) {
      forestShard.keys = shard.keys;
      forestShard.hash = shard.hash;
    } else {
      // If the shard doesn't exist in the forest structure, add it
      console.log(`Adding shard ${shard.id} to forest structure`);
      this.blockchain.forestStructure.shards.push({
        id: shard.id,
        level: 1, // Default level
        keys: shard.keys,
        hash: shard.hash
      });
    }

    // Check if shard needs to be split (if it has too many keys)
    this.checkAndSplitShard(shard);
  },

  // Check if a shard needs to be split and split it if necessary
  checkAndSplitShard(shard) {
    // Lower the threshold for shard splitting to make it happen more frequently
    const MAX_KEYS_PER_SHARD = 50; // Reduced from 500 to 50

    // If shard has too many keys and no children yet, split it
    if (shard.keys > MAX_KEYS_PER_SHARD && shard.childIds.length === 0) {
      console.log(`Splitting shard ${shard.id} with ${shard.keys} keys`);

      // Create two child shards
      const shardIdMatch = shard.id.match(/Shard-(\d+)(?:-(\d+))?/);
      let childId1, childId2;

      if (shardIdMatch) {
        const parentNum = parseInt(shardIdMatch[1]);
        if (shardIdMatch[2]) {
          // Already a second-level shard
          const subNum = parseInt(shardIdMatch[2]);
          childId1 = `${shard.id}-1`;
          childId2 = `${shard.id}-2`;
        } else {
          // First-level shard
          childId1 = `${shard.id}-1`;
          childId2 = `${shard.id}-2`;
        }
      } else {
        // Fallback if ID format doesn't match
        childId1 = `${shard.id}-1`;
        childId2 = `${shard.id}-2`;
      }

      // Calculate how to distribute statistics between child shards
      const keysForChild1 = Math.floor(shard.keys / 2);
      const keysForChild2 = shard.keys - keysForChild1;
      const ratio1 = keysForChild1 / shard.keys;
      const ratio2 = keysForChild2 / shard.keys;

      // Create a copy of the unique addresses set for each child
      const uniqueAddresses1 = new Set([...shard.stats.uniqueAddresses]);
      const uniqueAddresses2 = new Set([...shard.stats.uniqueAddresses]);

      // Create first child shard with proportional statistics
      const childShard1 = {
        id: childId1,
        parentId: shard.id,
        childIds: [],
        creationTime: Date.now(),
        keys: keysForChild1,
        hash: this.generateRandomHash(),
        lastUpdated: Date.now(),
        status: 'Active',
        stats: {
          totalTransactions: Math.floor(shard.stats.totalTransactions * ratio1),
          totalValue: shard.stats.totalValue * ratio1,
          averageValue: shard.stats.averageValue,
          uniqueAddresses: uniqueAddresses1,
          blockCount: Math.floor(shard.stats.blockCount * ratio1),
          lastBlockNumber: shard.stats.lastBlockNumber
        },
        operations: [
          {
            type: 'Creation',
            timestamp: Date.now(),
            parentShard: shard.id,
            initialKeys: keysForChild1
          }
        ]
      };

      // Create second child shard with proportional statistics
      const childShard2 = {
        id: childId2,
        parentId: shard.id,
        childIds: [],
        creationTime: Date.now(),
        keys: keysForChild2,
        hash: this.generateRandomHash(),
        lastUpdated: Date.now(),
        status: 'Active',
        stats: {
          totalTransactions: Math.floor(shard.stats.totalTransactions * ratio2),
          totalValue: shard.stats.totalValue * ratio2,
          averageValue: shard.stats.averageValue,
          uniqueAddresses: uniqueAddresses2,
          blockCount: Math.floor(shard.stats.blockCount * ratio2),
          lastBlockNumber: shard.stats.lastBlockNumber
        },
        operations: [
          {
            type: 'Creation',
            timestamp: Date.now(),
            parentShard: shard.id,
            initialKeys: keysForChild2
          }
        ]
      };

      // Add child shards to collection
      this.blockchain.shards[childId1] = childShard1;
      this.blockchain.shards[childId2] = childShard2;

      // Update parent shard
      shard.childIds = [childId1, childId2];
      shard.operations.push({
        type: 'Split',
        timestamp: Date.now(),
        childIds: [childId1, childId2]
      });

      // Add to forest structure
      const parentLevel = this.blockchain.forestStructure.shards.find(s => s.id === shard.id)?.level || 0;

      this.blockchain.forestStructure.shards.push({
        id: childId1,
        level: parentLevel + 1,
        keys: childShard1.keys,
        hash: childShard1.hash
      });

      this.blockchain.forestStructure.shards.push({
        id: childId2,
        level: parentLevel + 1,
        keys: childShard2.keys,
        hash: childShard2.hash
      });

      // Update AMF stats
      this.amfStats.totalShards += 2;
      this.amfStats.activeShards += 2;
      this.amfStats.shardDepth = Math.max(this.amfStats.shardDepth, parentLevel + 1);
      this.amfStats.rebalanceOperations++;

      console.log(`Split shard ${shard.id} into ${childId1} (${childShard1.keys} keys) and ${childId2} (${childShard2.keys} keys)`);
    }
  },

  // Create a new shard dynamically
  createNewShard() {
    // Find the highest shard ID number
    const shardIds = Object.keys(this.blockchain.shards)
      .filter(id => id !== 'Root Shard')
      .map(id => {
        const match = id.match(/Shard-(\d+)(?:-\d+)?/);
        return match ? parseInt(match[1]) : 0;
      });

    const maxShardId = Math.max(...shardIds, 0);
    const newShardId = maxShardId + 1;
    const newShardKey = `Shard-${newShardId}`;

    console.log(`Creating new shard ${newShardKey} (current total: ${this.amfStats.totalShards})`);

    // Create the new shard
    const newShard = {
      id: newShardKey,
      parentId: 'Root Shard',
      childIds: [],
      creationTime: Date.now(),
      keys: 0,
      hash: this.generateRandomHash(),
      lastUpdated: Date.now(),
      status: 'Active',
      stats: {
        totalTransactions: 0,
        totalValue: 0,
        averageValue: 0,
        uniqueAddresses: new Set(),
        blockCount: 0,
        lastBlockNumber: 0
      },
      operations: [
        { type: 'Creation', timestamp: Date.now() }
      ]
    };

    // Add to shards collection
    this.blockchain.shards[newShardKey] = newShard;

    // Add to forest structure
    this.blockchain.forestStructure.shards.push({
      id: newShardKey,
      level: 1,
      keys: 0,
      hash: newShard.hash
    });

    // Update parent shard
    const rootShard = this.blockchain.shards['Root Shard'];
    if (rootShard && !rootShard.childIds.includes(newShardKey)) {
      rootShard.childIds.push(newShardKey);
    }

    // Update AMF stats
    this.amfStats.totalShards++;
    this.amfStats.activeShards++;

    // Save state immediately after creating a new shard
    this.saveState();

    console.log(`Created new shard ${newShardKey}, updated AMF stats: ${this.amfStats.totalShards} total shards, ${this.amfStats.activeShards} active shards`);

    return newShardId;
  },

  // Generate a random transaction
  generateRandomTransaction() {
    const addresses = [
      '0x7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c',
      '0x6a5b4c3d2e1f0a9b8c7d6e5f4a3b2c1d0e9f8a7b',
      '0x5f4e3d2c1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f6e',
      '0x4d3c2b1a0f9e8d7c6b5a4f3e2d1c0b9a8f7e6d5c',
      '0x3b2a1f0e9d8c7b6a5f4e3d2c1b0a9f8e7d6c5b4',
      '0x2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c7d6e5f4a3',
      '0x1a0f9e8d7c6b5a4f3e2d1c0b9a8f7e6d5c4b3a2',
      '0x0f9e8d7c6b5a4f3e2d1c0b9a8f7e6d5c4b3a291'
    ];

    const fromIndex = Math.floor(Math.random() * addresses.length);
    let toIndex = Math.floor(Math.random() * addresses.length);

    // Make sure from and to are different
    while (toIndex === fromIndex) {
      toIndex = Math.floor(Math.random() * addresses.length);
    }

    const from = addresses[fromIndex];
    const to = addresses[toIndex];
    const amount = parseFloat((Math.random() * 5).toFixed(2));

    // Determine shard based on transaction data
    // Use the first character of the sender's address to determine shard
    // This creates a deterministic but distributed sharding mechanism
    const shardId = this.determineTransactionShard(from, to, amount);

    const tx = {
      hash: this.generateRandomHash(),
      blockNumber: null, // Will be set when included in a block
      timestamp: Date.now(),
      from: from,
      to: to,
      amount: amount,
      status: 'Pending',
      shardId: shardId,
      gasUsed: 21000,
      gasPrice: 20,
      data: '0x',
      signature: {
        r: this.generateRandomHash(),
        s: this.generateRandomHash(),
        v: 27 + Math.floor(Math.random() * 2)
      }
    };

    // Add to transaction pool
    this.blockchain.transactionPool.push(tx);
    this.stats.transactions++;

    // Save state after generating transactions (but not too frequently)
    if (this.blockchain.transactionPool.length % 5 === 0) {
      this.saveState();
    }

    console.log(`Generated transaction ${tx.hash} from ${tx.from} to ${tx.to} for ${tx.amount} (Shard: ${shardId})`);
  },

  // Determine which shard a transaction belongs to based on its data
  determineTransactionShard(from, to, amount) {
    // Get available shard IDs (excluding Root Shard)
    const shardIds = Object.keys(this.blockchain.shards)
      .filter(id => id !== 'Root Shard')
      .map(id => {
        // Extract numeric part if it's in format "Shard-X" or "Shard-X-Y"
        const match = id.match(/Shard-(\d+)(?:-\d+)?/);
        return match ? parseInt(match[1]) : 1;
      });

    // If no shards available yet, default to shard 1
    if (shardIds.length === 0) {
      return 1;
    }

    // Dynamically create new shards as transaction volume increases
    // Create a new shard when we have more than 10 transactions per shard on average
    const totalTransactions = this.stats.transactions;
    const transactionsPerShard = totalTransactions / shardIds.length;

    // Create a new shard if we have more than 10 transactions per shard
    // and we have fewer than 10 shards total
    if (transactionsPerShard > 10 && shardIds.length < 10 && Math.random() < 0.3) {
      // 30% chance to create a new shard when conditions are met
      const newShardId = this.createNewShard();
      shardIds.push(newShardId);
      console.log(`Created new shard ${newShardId} due to transaction volume (${transactionsPerShard.toFixed(1)} tx/shard)`);
    }

    // Use consistent hashing based on the sender's address
    // This ensures the same sender's transactions tend to go to the same shard
    const addressNum = parseInt(from.slice(2, 10), 16);

    // For cross-shard transactions (high value transactions), use a different algorithm
    if (amount > 3.0) {
      // High-value transactions go to shard 1 (main shard) for special handling
      return 1;
    } else {
      // Normal transactions are distributed based on sender address
      // This creates a deterministic but distributed sharding mechanism
      return shardIds[addressNum % shardIds.length];
    }
  },

  // Create a transaction from user input
  createTransaction(transactionData) {
    const { to, amount, data } = transactionData;
    const from = '0x7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c'; // Default sender
    const parsedAmount = parseFloat(amount);

    // Validate recipient address format
    // Check if it starts with 0x and contains only valid hex characters
    // Allow both 40 and 42 character addresses (with or without 0x prefix)
    if (!to) {
      throw new Error('Recipient address is required');
    }

    // Normalize the address to ensure it has 0x prefix
    let normalizedTo = to;
    if (!normalizedTo.startsWith('0x')) {
      normalizedTo = '0x' + normalizedTo;
    }

    // Validate the address format
    const addressRegex = /^0x[0-9a-fA-F]{40}$/;
    if (!addressRegex.test(normalizedTo)) {
      throw new Error('Invalid recipient address format. Must be a valid 0x... address with 40 hex characters');
    }

    // Determine shard based on transaction data
    let shardId;
    if (transactionData.shardId) {
      // If user specified a shard, use it
      shardId = parseInt(transactionData.shardId);
    } else {
      // Otherwise determine shard automatically
      shardId = this.determineTransactionShard(from, normalizedTo, parsedAmount);
    }

    const tx = {
      hash: this.generateRandomHash(),
      blockNumber: null, // Will be set when included in a block
      timestamp: Date.now(),
      from: from,
      to: normalizedTo,
      amount: parsedAmount,
      status: 'Pending',
      shardId: shardId,
      gasUsed: 21000,
      gasPrice: 20,
      data: data || '0x',
      signature: {
        r: this.generateRandomHash(),
        s: this.generateRandomHash(),
        v: 27
      }
    };

    // Add to transaction pool
    this.blockchain.transactionPool.push(tx);
    this.stats.transactions++;

    // Create a new shard if needed based on transaction volume
    // This ensures the AMF grows as transaction volume increases
    if (this.stats.transactions % 5 === 0) { // Every 5 transactions
      // Check if we need more shards
      const shardCount = Object.keys(this.blockchain.shards).filter(id => id !== 'Root Shard').length;

      // Create a new shard if we have fewer than 10 shards
      if (shardCount < 10 && Math.random() < 0.3) {
        const newShardId = this.createNewShard();
        console.log(`Created new shard ${newShardId} after transaction ${tx.hash}`);

        // Update AMF stats
        this.updateAMFStats();
      }
    }

    // Save state immediately after user-created transaction
    this.saveState();

    console.log(`Created transaction ${tx.hash} from ${tx.from} to ${tx.to} for ${tx.amount} (Shard: ${shardId})`);

    return { hash: tx.hash };
  },

  // Update node metrics
  updateNodeMetrics() {
    this.blockchain.nodes.forEach(node => {
      // Randomly update node status
      if (Math.random() < 0.05) {
        if (node.status === 'Online') {
          node.status = 'Syncing';
          node.networkState = 'Syncing';
        } else if (node.status === 'Syncing') {
          node.status = 'Online';
          node.networkState = 'Normal';
        }
      }

      // Update uptime for online nodes
      if (node.status !== 'Offline') {
        node.uptime += 10;
      }

      // Update metrics
      node.metrics.cpu = Math.min(100, Math.max(5, node.metrics.cpu + (Math.random() * 10 - 5)));
      node.metrics.memory = Math.min(2000000000, Math.max(500000000, node.metrics.memory + (Math.random() * 100000000 - 50000000)));
      node.metrics.networkIn = parseFloat((Math.min(10, Math.max(0.1, parseFloat(node.metrics.networkIn) + (Math.random() - 0.5)))).toFixed(1));
      node.metrics.networkOut = parseFloat((Math.min(10, Math.max(0.1, parseFloat(node.metrics.networkOut) + (Math.random() - 0.5)))).toFixed(1));
    });

    // Update network stats
    const activeNodes = this.blockchain.nodes.filter(node => node.status !== 'Offline').length;
    this.networkStats.activeNodes = activeNodes;
    this.networkStats.networkLatency = Math.floor(Math.random() * 50) + 100;
    this.networkStats.byzantineNodes = Math.floor(Math.random() * 2);
  },

  // Update all AMF shards (periodic maintenance)
  updateAMF() {
    // Update all shards
    Object.values(this.blockchain.shards).forEach(shard => {
      // Update hash
      shard.hash = this.generateRandomHash();
      shard.lastUpdated = Date.now();

      // Add operation
      shard.operations.push({
        type: 'Maintenance',
        timestamp: Date.now()
      });

      // Update corresponding forest structure entry
      const forestShard = this.blockchain.forestStructure.shards.find(s => s.id === shard.id);
      if (forestShard) {
        forestShard.hash = shard.hash;
      }
    });

    // Check for shards that need rebalancing
    this.checkForShardRebalancing();

    // Periodically create new shards to increase the active shards count
    // Get the number of top-level shards (direct children of Root Shard)
    const topLevelShards = Object.values(this.blockchain.shards)
      .filter(shard => shard.parentId === 'Root Shard' && shard.id !== 'Root Shard');

    // Create a new shard if we have fewer than 5 top-level shards
    // and we have a reasonable number of transactions
    if (topLevelShards.length < 5 && this.stats.transactions > 20 && Math.random() < 0.4) {
      // 40% chance to create a new shard when conditions are met
      this.createNewShard();
    }

    // Always update AMF stats after maintenance
    this.updateAMFStats();

    // Create additional shards based on transaction volume
    if (this.stats.transactions > 0) {
      // Calculate the ideal number of shards based on transaction volume
      // 1 shard per 10 transactions, up to a maximum of 10 shards
      const idealShardCount = Math.min(10, Math.max(2, Math.floor(this.stats.transactions / 10)));
      const currentShardCount = Object.keys(this.blockchain.shards).filter(id => id !== 'Root Shard').length;

      // If we need more shards, create them
      if (currentShardCount < idealShardCount && Math.random() < 0.5) {
        const newShardId = this.createNewShard();
        console.log(`Created new shard ${newShardId} during AMF maintenance (${currentShardCount} -> ${currentShardCount + 1})`);

        // Update AMF stats again after creating a new shard
        this.updateAMFStats();
      }
    }
  },

  // Check if any shards need rebalancing
  checkForShardRebalancing() {
    // Find leaf shards (shards with no children)
    const leafShards = Object.values(this.blockchain.shards).filter(shard => shard.childIds.length === 0);

    // Check for imbalance between leaf shards
    if (leafShards.length >= 2) {
      // Sort by key count
      leafShards.sort((a, b) => a.keys - b.keys);

      // If the difference between min and max is too large, rebalance
      const minShard = leafShards[0];
      const maxShard = leafShards[leafShards.length - 1];

      if (maxShard.keys > minShard.keys * 3) {
        console.log(`Rebalancing shards: ${minShard.id} (${minShard.keys} keys) and ${maxShard.id} (${maxShard.keys} keys)`);

        // Transfer some keys from max to min
        const keysToTransfer = Math.floor(maxShard.keys * 0.3); // Transfer 30% of keys
        const transferRatio = keysToTransfer / maxShard.keys;

        // Update key counts
        minShard.keys += keysToTransfer;
        maxShard.keys -= keysToTransfer;

        // Transfer proportional statistics
        const txToTransfer = Math.floor(maxShard.stats.totalTransactions * transferRatio);
        const valueToTransfer = maxShard.stats.totalValue * transferRatio;
        const blocksToTransfer = Math.floor(maxShard.stats.blockCount * transferRatio);

        // Update transaction statistics
        minShard.stats.totalTransactions += txToTransfer;
        maxShard.stats.totalTransactions -= txToTransfer;

        // Update value statistics
        minShard.stats.totalValue += valueToTransfer;
        maxShard.stats.totalValue -= valueToTransfer;

        // Update average values
        if (minShard.stats.totalTransactions > 0) {
          minShard.stats.averageValue = minShard.stats.totalValue / minShard.stats.totalTransactions;
        }
        if (maxShard.stats.totalTransactions > 0) {
          maxShard.stats.averageValue = maxShard.stats.totalValue / maxShard.stats.totalTransactions;
        }

        // Update block counts
        minShard.stats.blockCount += blocksToTransfer;
        maxShard.stats.blockCount -= blocksToTransfer;

        // Share unique addresses
        maxShard.stats.uniqueAddresses.forEach(addr => {
          minShard.stats.uniqueAddresses.add(addr);
        });

        // Update hashes
        minShard.hash = this.generateRandomHash();
        maxShard.hash = this.generateRandomHash();

        // Update timestamps
        const now = Date.now();
        minShard.lastUpdated = now;
        maxShard.lastUpdated = now;

        // Add operations
        minShard.operations.push({
          type: 'Rebalance (Receive)',
          timestamp: now,
          keysReceived: keysToTransfer,
          fromShard: maxShard.id
        });

        maxShard.operations.push({
          type: 'Rebalance (Send)',
          timestamp: now,
          keysSent: keysToTransfer,
          toShard: minShard.id
        });

        // Update forest structure
        const minForestShard = this.blockchain.forestStructure.shards.find(s => s.id === minShard.id);
        const maxForestShard = this.blockchain.forestStructure.shards.find(s => s.id === maxShard.id);

        if (minForestShard) {
          minForestShard.keys = minShard.keys;
          minForestShard.hash = minShard.hash;
        }

        if (maxForestShard) {
          maxForestShard.keys = maxShard.keys;
          maxForestShard.hash = maxShard.hash;
        }

        // Update stats
        this.amfStats.rebalanceOperations++;

        console.log(`Rebalanced: ${minShard.id} now has ${minShard.keys} keys, ${maxShard.id} now has ${maxShard.keys} keys`);
      }
    }
  },

  // Add a new node
  addNode(nodeData) {
    const { address, type = 'validator', initialState = 'sync' } = nodeData;

    // Generate a new node ID
    const id = 'node-' + (this.blockchain.nodes.length + 1);

    // Determine initial status based on initialState
    const status = initialState === 'sync' ? 'Syncing' : 'Online';
    const networkState = initialState === 'sync' ? 'Syncing' : 'Normal';

    // Determine consistency level based on node type
    let consistencyLevel = 'Strong';
    if (type === 'full') {
      consistencyLevel = 'Causal';
    } else if (type === 'light') {
      consistencyLevel = 'Eventual';
    }

    // Create a new node
    const newNode = {
      id,
      address,
      status,
      uptime: 0,
      peers: 0,
      nodeType: type,
      consistencyLevel,
      publicKey: this.generateRandomHash(),
      version: '1.0.0',
      networkState,
      metrics: {
        cpu: 10,
        memory: 800000000,
        disk: 8000000000,
        networkIn: '0.5',
        networkOut: '0.3',
        blocks: this.stats.blocks,
        transactions: this.stats.transactions
      }
    };

    // Add to nodes
    this.blockchain.nodes.push(newNode);
    this.networkStats.totalNodes++;
    this.networkStats.activeNodes++;

    console.log(`Added new node ${id} at ${address} (Type: ${type}, Initial State: ${initialState})`);

    return { id };
  },

  // Restart a node
  restartNode(nodeId) {
    const node = this.blockchain.nodes.find(n => n.id === nodeId);

    if (!node) {
      throw new Error('Node not found');
    }

    // Update node status
    node.status = 'Syncing';
    node.uptime = 0;
    node.networkState = 'Syncing';

    console.log(`Restarted node ${nodeId}`);

    return { success: true };
  },

  // Set a node offline
  setNodeOffline(nodeId) {
    const node = this.blockchain.nodes.find(n => n.id === nodeId);

    if (!node) {
      throw new Error('Node not found');
    }

    // Only change status if node is not already offline
    if (node.status.toLowerCase() !== 'offline') {
      // Update node status
      node.status = 'Offline';
      node.networkState = 'Disconnected';
      node.peers = 0;

      // Decrease active nodes count
      this.networkStats.activeNodes--;

      console.log(`Set node ${nodeId} offline`);
    } else {
      console.log(`Node ${nodeId} is already offline`);
    }

    return { success: true };
  },

  // Remove a node
  removeNode(nodeId) {
    const nodeIndex = this.blockchain.nodes.findIndex(n => n.id === nodeId);

    if (nodeIndex === -1) {
      throw new Error('Node not found');
    }

    // Remove node
    this.blockchain.nodes.splice(nodeIndex, 1);
    this.networkStats.totalNodes--;
    this.networkStats.activeNodes--;

    console.log(`Removed node ${nodeId}`);

    return { success: true };
  },

  // Generate a random hash
  generateRandomHash() {
    return '0x' + Array.from({length: 40}, () => '0123456789abcdef'[Math.floor(Math.random() * 16)]).join('');
  },

  // Get blockchain status
  getStatus() {
    return { status: 'running' };
  },

  // Get blockchain stats
  getStats() {
    return this.stats;
  },

  // Get network stats
  getNetworkStats() {
    return this.networkStats;
  },

  // Get AMF stats
  getAMFStats() {
    return this.amfStats;
  },

  // Get AMF forest structure
  getForestStructure() {
    return this.blockchain.forestStructure;
  },

  // Get a specific shard
  getShard(shardId) {
    const shard = this.blockchain.shards[shardId];
    if (!shard) {
      throw new Error('Shard not found');
    }
    return shard;
  },

  // Get operations for a specific shard
  getShardOperations(shardId) {
    const shard = this.blockchain.shards[shardId];
    return shard ? shard.operations : [];
  },

  // Get all nodes
  getNodes() {
    return this.blockchain.nodes;
  },

  // Get a specific node
  getNode(nodeId) {
    const node = this.blockchain.nodes.find(n => n.id === nodeId);
    if (!node) {
      throw new Error('Node not found');
    }
    return node;
  },

  // Get all blocks with pagination and filters
  getBlocks(page = 1, limit = 10, filters = {}) {
    // Ensure blocks are created for all transactions
    this.ensureTransactionsBlockSync();

    // Get all blocks in reverse order (newest first)
    let blocks = [...this.blockchain.blocks].reverse();

    // Apply filters
    if (filters.shardId) {
      blocks = blocks.filter(block => block.shardId.toString() === filters.shardId);
    }

    // Paginate
    const startIndex = (page - 1) * limit;
    const endIndex = startIndex + limit;
    const paginatedBlocks = blocks.slice(startIndex, endIndex);

    return {
      blocks: paginatedBlocks,
      totalPages: Math.ceil(blocks.length / limit)
    };
  },



  // Get latest blocks
  getLatestBlocks(limit = 5) {
    // Ensure blocks are created for all transactions
    this.ensureTransactionsBlockSync();

    return [...this.blockchain.blocks].reverse().slice(0, limit);
  },

  // Get a specific block
  getBlock(blockId) {
    // Ensure blocks are created for all transactions
    this.ensureTransactionsBlockSync();

    if (/^0x[0-9a-f]+$/i.test(blockId)) {
      // It's a hash
      return this.blockchain.blocks.find(block => block.hash === blockId);
    } else {
      // It's a number
      return this.blockchain.blocks.find(block => block.number === parseInt(blockId));
    }
  },

  // Get transactions for a block
  getBlockTransactions(blockId) {
    const block = this.getBlock(blockId);
    if (!block) {
      throw new Error('Block not found');
    }

    // Find transactions for this block
    const transactions = this.blockchain.transactions.filter(tx => tx.blockNumber === block.number);

    // If no transactions are found but the block says it has transactions,
    // generate some placeholder transactions for display purposes
    if (transactions.length === 0 && block.transactions > 0) {
      console.log(`No transactions found for block ${block.number}, generating ${block.transactions} placeholder transactions`);

      const addresses = [
        '0x7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c',
        '0x6a5b4c3d2e1f0a9b8c7d6e5f4a3b2c1d0e9f8a7b',
        '0x5f4e3d2c1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f6e',
        '0x4d3c2b1a0f9e8d7c6b5a4f3e2d1c0b9a8f7e6d5c'
      ];

      const placeholderTransactions = [];

      for (let i = 0; i < block.transactions; i++) {
        const fromIndex = Math.floor(Math.random() * addresses.length);
        let toIndex = Math.floor(Math.random() * addresses.length);

        // Make sure from and to are different
        while (toIndex === fromIndex) {
          toIndex = Math.floor(Math.random() * addresses.length);
        }

        const tx = {
          hash: this.generateRandomHash(),
          blockNumber: block.number,
          timestamp: block.timestamp - Math.floor(Math.random() * 30000), // Within 30 seconds before block
          from: addresses[fromIndex],
          to: addresses[toIndex],
          amount: parseFloat((Math.random() * 5).toFixed(2)),
          status: 'Confirmed',
          shardId: block.shardId,
          gasUsed: 21000,
          gasPrice: 20,
          data: '0x',
          signature: {
            r: this.generateRandomHash(),
            s: this.generateRandomHash(),
            v: 27 + Math.floor(Math.random() * 2)
          }
        };

        placeholderTransactions.push(tx);
      }

      return placeholderTransactions;
    }

    return transactions;
  },

  // Get all transactions with pagination and filters
  getTransactions(page = 1, limit = 10, filters = {}) {
    // Ensure we have blocks and transactions are properly linked
    this.ensureTransactionsBlockSync();

    let transactions = [...this.blockchain.transactions, ...this.blockchain.transactionPool]
      .sort((a, b) => b.timestamp - a.timestamp);

    // Apply filters
    if (filters.address) {
      transactions = transactions.filter(tx =>
        tx.from.toLowerCase() === filters.address.toLowerCase() ||
        tx.to.toLowerCase() === filters.address.toLowerCase()
      );
    }

    if (filters.status) {
      transactions = transactions.filter(tx =>
        tx.status.toLowerCase() === filters.status.toLowerCase()
      );
    }

    if (filters.shardId) {
      transactions = transactions.filter(tx =>
        tx.shardId.toString() === filters.shardId
      );
    }

    // Paginate
    const startIndex = (page - 1) * limit;
    const endIndex = startIndex + limit;
    const paginatedTransactions = transactions.slice(startIndex, endIndex);

    return {
      transactions: paginatedTransactions,
      totalPages: Math.ceil(transactions.length / limit)
    };
  },

  // Ensure transactions reference valid blocks
  ensureTransactionsBlockSync() {
    // If we don't have any blocks yet, nothing to sync
    if (this.blockchain.blocks.length === 0) {
      return;
    }

    // Get all valid block numbers
    const validBlockNumbers = new Set(this.blockchain.blocks.map(block => block.number));

    // Get the latest block for reference
    const latestBlock = this.blockchain.blocks[this.blockchain.blocks.length - 1];

    // Find all unique invalid block numbers referenced by transactions
    const invalidBlockNumbers = new Set();
    this.blockchain.transactions.forEach(tx => {
      if (tx.blockNumber !== null && !validBlockNumbers.has(tx.blockNumber)) {
        invalidBlockNumbers.add(tx.blockNumber);
      }
    });

    // Create missing blocks
    if (invalidBlockNumbers.size > 0) {
      console.log(`Creating ${invalidBlockNumbers.size} missing blocks referenced by transactions`);

      // Convert to array and sort
      const missingBlockNumbers = Array.from(invalidBlockNumbers).sort((a, b) => a - b);

      for (const blockNumber of missingBlockNumbers) {
        // Skip if the block number is less than 1 (genesis block)
        if (blockNumber < 1) continue;

        // Skip if we already have this block
        if (validBlockNumbers.has(blockNumber)) continue;

        // Find transactions for this block
        const blockTransactions = this.blockchain.transactions.filter(tx => tx.blockNumber === blockNumber);
        const transactionCount = blockTransactions.length;

        // Determine the shard ID for this block based on its transactions
        let shardId = 1; // Default to shard 1

        if (blockTransactions.length > 0) {
          // Count transactions by shard
          const shardCounts = {};
          blockTransactions.forEach(tx => {
            shardCounts[tx.shardId] = (shardCounts[tx.shardId] || 0) + 1;
          });

          // Find the most common shard
          let maxCount = 0;
          for (const [sId, count] of Object.entries(shardCounts)) {
            if (count > maxCount) {
              maxCount = count;
              shardId = parseInt(sId);
            }
          }
        }

        // Create the missing block
        const newBlock = {
          number: blockNumber,
          hash: this.generateRandomHash(),
          previousHash: blockNumber > 1 ?
            (this.blockchain.blocks.find(b => b.number === blockNumber - 1)?.hash || this.generateRandomHash()) :
            '0x0000000000000000000000000000000000000000000000000000000000000000',
          timestamp: blockTransactions.length > 0 ?
            Math.max(...blockTransactions.map(tx => tx.timestamp)) + 10000 : // 10 seconds after latest transaction
            latestBlock.timestamp - ((latestBlock.number - blockNumber) * 15000), // Estimate based on block interval
          transactions: transactionCount,
          shardId: shardId,
          size: 1024 + (transactionCount * 256),
          merkleRoot: this.generateRandomHash(),
          stateRoot: this.generateRandomHash(),
          difficulty: 1,
          nonce: Math.floor(Math.random() * 1000000)
        };

        console.log(`Created missing block #${blockNumber} with ${transactionCount} transactions`);

        // Insert the block in the correct position to maintain order
        const insertIndex = this.blockchain.blocks.findIndex(b => b.number > blockNumber);
        if (insertIndex === -1) {
          this.blockchain.blocks.push(newBlock);
        } else {
          this.blockchain.blocks.splice(insertIndex, 0, newBlock);
        }

        this.stats.blocks++;
        validBlockNumbers.add(blockNumber);

        // Update the shard data for this block
        this.updateShardForNewBlock(newBlock.shardId, newBlock, blockTransactions);
      }

      // Save state immediately after creating blocks
      this.saveState();
    }

    // Now fix any remaining transactions with invalid block numbers
    this.blockchain.transactions.forEach(tx => {
      if (tx.blockNumber !== null && !validBlockNumbers.has(tx.blockNumber)) {
        // Assign to the latest block by default
        const targetBlock = latestBlock;

        console.log(`Fixing transaction ${tx.hash} with invalid block number ${tx.blockNumber} -> ${targetBlock.number}`);

        tx.blockNumber = targetBlock.number;
        tx.timestamp = targetBlock.timestamp - Math.floor(Math.random() * 30000); // Within 30 seconds before block
        tx.status = 'Confirmed';
      }
    });
  },

  // Get a specific transaction
  getTransaction(txHash) {
    // Find the transaction
    const tx = this.blockchain.transactions.find(tx => tx.hash === txHash) ||
               this.blockchain.transactionPool.find(tx => tx.hash === txHash);

    // If transaction is found and has a block number, verify it's valid
    if (tx && tx.blockNumber !== null && this.blockchain.blocks.length > 0) {
      const validBlockNumbers = new Set(this.blockchain.blocks.map(block => block.number));

      if (!validBlockNumbers.has(tx.blockNumber)) {
        // Get the latest block
        const latestBlock = this.blockchain.blocks[this.blockchain.blocks.length - 1];

        console.log(`Fixing transaction ${tx.hash} with invalid block number ${tx.blockNumber} -> ${latestBlock.number}`);

        tx.blockNumber = latestBlock.number;
        tx.timestamp = latestBlock.timestamp - Math.floor(Math.random() * 30000); // Within 30 seconds before block
        tx.status = 'Confirmed';
      }
    }

    return tx;
  },

  // Get transaction status by hash
  getTransactionStatus(txHash) {
    const tx = this.getTransaction(txHash);
    if (!tx) {
      return { found: false, status: 'Unknown' };
    }

    return {
      found: true,
      status: tx.status,
      blockNumber: tx.blockNumber,
      timestamp: tx.timestamp,
      confirmationTime: tx.blockNumber ?
        (this.blockchain.blocks.find(b => b.number === tx.blockNumber)?.timestamp - tx.timestamp) / 1000 :
        null
    };
  },

  // Get network activity data for charts
  getNetworkActivity(period = 'day', chartType = 'transactions') {
    try {
      // Ensure we have blocks and transactions
      this.ensureTransactionsBlockSync();

      // Get current time
      const now = Date.now();

      // Determine time intervals based on period
      let timeIntervals = [];
      let intervalMs = 0;
      let intervalCount = 6;

      switch(period.toLowerCase()) {
        case 'hour':
          intervalMs = 10 * 60 * 1000; // 10 minutes
          intervalCount = 6;
          break;
        case 'day':
          intervalMs = 4 * 60 * 60 * 1000; // 4 hours
          intervalCount = 6;
          break;
        case 'week':
          intervalMs = 24 * 60 * 60 * 1000; // 1 day
          intervalCount = 7;
          break;
        case 'month':
          intervalMs = 6 * 24 * 60 * 60 * 1000; // 6 days
          intervalCount = 5;
          break;
        default:
          intervalMs = 4 * 60 * 60 * 1000; // 4 hours
          intervalCount = 6;
      }

      // Create time intervals
      for (let i = intervalCount - 1; i >= 0; i--) {
        const intervalStart = now - (i * intervalMs);
        timeIntervals.push(intervalStart);
      }

      // Format labels based on period
      const labels = timeIntervals.map(time => {
        const date = new Date(time);
        switch(period.toLowerCase()) {
          case 'hour':
            return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
          case 'day':
            return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
          case 'week':
            return date.toLocaleDateString([], { weekday: 'short' });
          case 'month':
            return date.toLocaleDateString([], { month: 'short', day: 'numeric' });
          default:
            return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
        }
      });

      // Prepare data based on chart type
      let datasets = [];

      if (chartType === 'transactions' || !chartType) {
        // Count transactions in each time interval
        const transactionCounts = timeIntervals.map((intervalStart, index) => {
          const intervalEnd = index < timeIntervals.length - 1 ? timeIntervals[index + 1] : now;

          // Count confirmed transactions in this interval
          const confirmedCount = this.blockchain.transactions.filter(tx =>
            tx.timestamp >= intervalStart &&
            tx.timestamp < intervalEnd &&
            tx.status === 'Confirmed'
          ).length;

          // Count pending transactions in this interval
          const pendingCount = this.blockchain.transactionPool.filter(tx =>
            tx.timestamp >= intervalStart &&
            tx.timestamp < intervalEnd
          ).length;

          return { confirmed: confirmedCount, pending: pendingCount };
        });

        // Create datasets for transactions
        datasets = [
          {
            label: 'Confirmed Transactions',
            data: transactionCounts.map(count => count.confirmed),
            borderColor: 'rgba(56, 176, 0, 1)',
            backgroundColor: 'rgba(56, 176, 0, 0.2)'
          },
          {
            label: 'Pending Transactions',
            data: transactionCounts.map(count => count.pending),
            borderColor: 'rgba(255, 190, 11, 1)',
            backgroundColor: 'rgba(255, 190, 11, 0.2)'
          }
        ];
      } else if (chartType === 'blockTime') {
        // Calculate average block time in each interval
        const blockTimes = timeIntervals.map((intervalStart, index) => {
          const intervalEnd = index < timeIntervals.length - 1 ? timeIntervals[index + 1] : now;

          // Get blocks created in this interval
          const blocksInInterval = this.blockchain.blocks.filter(block =>
            block.timestamp >= intervalStart &&
            block.timestamp < intervalEnd
          );

          // Calculate average time between blocks
          if (blocksInInterval.length <= 1) {
            return 15; // Default block time if not enough blocks
          }

          let totalTime = 0;
          let count = 0;

          for (let i = 1; i < blocksInInterval.length; i++) {
            const timeDiff = (blocksInInterval[i].timestamp - blocksInInterval[i-1].timestamp) / 1000; // in seconds
            if (timeDiff > 0 && timeDiff < 300) { // Ignore outliers (> 5 minutes)
              totalTime += timeDiff;
              count++;
            }
          }

          return count > 0 ? totalTime / count : 15; // Average in seconds
        });

        // Create dataset for block time
        datasets = [
          {
            label: 'Block Time (seconds)',
            data: blockTimes,
            borderColor: 'rgba(58, 134, 255, 1)',
            backgroundColor: 'rgba(58, 134, 255, 0.2)'
          }
        ];
      } else if (chartType === 'nodeDistribution') {
        // Count nodes by type
        const nodeTypes = {};
        this.blockchain.nodes.forEach(node => {
          const type = node.nodeType || 'unknown';
          nodeTypes[type] = (nodeTypes[type] || 0) + 1;
        });

        // Create dataset for node distribution
        datasets = [
          {
            label: 'Node Types',
            data: Object.values(nodeTypes),
            backgroundColor: [
              'rgba(58, 134, 255, 0.8)',
              'rgba(131, 56, 236, 0.8)',
              'rgba(56, 176, 0, 0.8)',
              'rgba(255, 190, 11, 0.8)'
            ]
          }
        ];

        // Override labels for node distribution
        labels = Object.keys(nodeTypes);
      }

      return {
        labels: labels,
        datasets: datasets
      };
    } catch (error) {
      console.error('Error generating network activity data:', error);

      // Return fallback data
      return {
        labels: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00'],
        datasets: [
          {
            label: 'Transactions',
            data: [35, 42, 67, 89, 76, 54]
          },
          {
            label: 'Blocks',
            data: [12, 15, 22, 28, 24, 18]
          }
        ]
      };
    }
  },

  // Get proof for a key in a shard
  getProof(shardId, key) {
    // Get the shard
    const shard = this.blockchain.shards[shardId];
    if (!shard) {
      throw new Error(`Shard ${shardId} not found`);
    }

    // Initialize the shard's keyMap if it doesn't exist
    if (!shard.keyMap) {
      shard.keyMap = this.generateShardKeyMap(shard.id, shard.keys);
    }

    // Check if the key exists in the shard
    const keyExists = this.checkKeyInShard(key, shard);

    // If the key doesn't exist, we should return an invalid proof
    if (!keyExists) {
      // For demonstration purposes, we'll still generate a proof but mark it as invalid
      const invalidProof = {
        shardId: shardId,
        key: key,
        value: "INVALID_KEY",
        rootHash: shard.hash,
        timestamp: Date.now(),
        valid: false,
        error: `Key '${key}' not found in shard ${shardId}`,
        siblings: []
      };

      return {
        shardId,
        key,
        value: "INVALID_KEY",
        proof: JSON.stringify(invalidProof),
        metadata: {
          shardHash: shard.hash,
          timestamp: new Date().toISOString(),
          proofType: 'Merkle',
          valid: false,
          error: `Key '${key}' not found in shard ${shardId}`
        }
      };
    }

    // Get the value for this key
    const value = this.getValueForKey(key, shard);

    // Generate a leaf hash for this key-value pair
    const leafHash = this.hashString(key + value);

    // Generate a deterministic path in the Merkle tree for this key
    const keyHash = this.hashString(key);
    const keyPosition = parseInt(keyHash.substring(0, 8), 16) % 256; // Position in a tree with 256 leaves

    // Generate a consistent root hash for this key-value pair
    // In a real implementation, this would be the actual root hash of the Merkle tree
    // For our demo, we'll generate a deterministic root hash based on the shard and key
    const rootHash = this.hashString(shard.id + key + "root");

    // Generate sibling hashes for the Merkle path
    // In a real implementation, these would be the actual sibling hashes from the Merkle tree
    const siblings = [];
    let currentPos = keyPosition;
    let level = 0;

    // Generate siblings for each level of the tree (log2(256) = 8 levels)
    while (level < 8) {
      // Determine if we're a left or right child
      const isRightChild = (currentPos % 2) === 1;

      // Generate the sibling hash - make it deterministic based on the key and level
      const siblingPos = isRightChild ? currentPos - 1 : currentPos + 1;
      const siblingHash = this.hashString(`${shard.id}-${key}-level-${level}-position-${siblingPos}`);

      // Add to siblings array
      siblings.push(siblingHash);

      // Move up the tree
      currentPos = Math.floor(currentPos / 2);
      level++;
    }

    // Create a simulated Merkle proof with realistic structure
    const merkleProof = {
      shardId: shardId,
      key: key,
      value: value,
      rootHash: rootHash, // Use our deterministic root hash
      timestamp: Date.now(),
      leafHash: leafHash,
      treePosition: keyPosition,
      treeHeight: 8,
      siblings: siblings,
      // Include a proof hash that combines all elements
      proofHash: this.hashString(leafHash + siblings.join('') + rootHash)
    };

    // Serialize the proof to a string
    const proofString = JSON.stringify(merkleProof);

    return {
      shardId,
      key,
      value,
      // The actual proof would be a binary format in a real implementation
      proof: proofString,
      // Include additional metadata for display purposes
      metadata: {
        shardHash: shard.hash,
        timestamp: new Date().toISOString(),
        proofType: 'Merkle',
        proofSize: proofString.length,
        treePosition: keyPosition,
        treeHeight: 8
      }
    };
  },

  // Verify a proof
  verifyProof(proofData) {
    // In a real implementation, this would call the backend API
    // For demonstration purposes, we'll implement a basic verification

    // Check if all required fields are present
    if (!proofData || !proofData.shardId || !proofData.key || !proofData.value || !proofData.proof) {
      return {
        valid: false,
        error: "Missing required fields in proof data"
      };
    }

    // Get the shard
    const shard = this.blockchain.shards[proofData.shardId];
    if (!shard) {
      return {
        valid: false,
        error: `Shard ${proofData.shardId} not found`
      };
    }

    // Check if the specific key exists in the shard
    // In a real implementation, this would query the actual state database
    // For our simulation, we'll implement a more realistic check

    // First, initialize the shard's keyMap if it doesn't exist
    if (!shard.keyMap) {
      // Create a deterministic key map based on the shard ID and keys count
      shard.keyMap = this.generateShardKeyMap(shard.id, shard.keys);
    }

    // Check if the specific key exists in the shard's key map
    const keyExists = this.checkKeyInShard(proofData.key, shard);

    // If the key doesn't exist, fail verification immediately
    if (!keyExists) {
      return {
        valid: false,
        error: `Key '${proofData.key}' not found in shard ${proofData.shardId}`,
        details: {
          shardId: proofData.shardId,
          key: proofData.key,
          valueVerified: false,
          keyExists: false,
          reason: "Key not found in shard"
        }
      };
    }

    // Verify that the value matches what we expect for this key
    const expectedValue = this.getValueForKey(proofData.key, shard);
    const valueMatches = expectedValue === proofData.value;

    // If the value doesn't match, fail verification
    if (!valueMatches) {
      return {
        valid: false,
        error: `Value mismatch for key '${proofData.key}' in shard ${proofData.shardId}`,
        details: {
          shardId: proofData.shardId,
          key: proofData.key,
          valueVerified: false,
          keyExists: true,
          valueMatches: false,
          expectedValue: expectedValue,
          providedValue: proofData.value,
          reason: "Value mismatch"
        }
      };
    }

    // Generate a deterministic result based on the proof data
    // This simulates cryptographic verification without actually implementing it
    const hash = this.hashString(proofData.shardId + proofData.key + proofData.value + proofData.proof);
    const lastChar = hash.charAt(hash.length - 1);

    // Verify the Merkle path using the sibling hashes
    // Parse the proof to get the Merkle path details
    let merklePathValid = false;
    let merklePathError = null;

    try {
      // Parse the proof JSON
      const parsedProof = JSON.parse(proofData.proof);

      // Check if we have all the required fields for Merkle path verification
      if (parsedProof.leafHash && parsedProof.siblings && parsedProof.treePosition !== undefined) {
        // For demonstration purposes, we'll make the verification more reliable
        // by using a deterministic approach based on the key and value

        // Generate a deterministic hash from the key and value
        const keyValueHash = this.hashString(proofData.key + proofData.value);

        // In a real implementation, we would verify the entire Merkle path
        // For this demo, we'll use a simplified approach that's more consistent

        // Check if the key is one of our generated keys (starts with "key-")
        if (proofData.key.startsWith("key-")) {
          // For our generated keys, we'll make the verification pass
          // This simulates a correct Merkle path for keys we know exist
          merklePathValid = true;
        } else {
          // For other keys, use the hash check to determine validity
          // This ensures some proofs will still fail for demonstration
          const keyHash = this.hashString(proofData.key);
          merklePathValid = keyHash.charAt(keyHash.length - 1) !== '0';

          if (!merklePathValid) {
            merklePathError = "Merkle path verification failed: invalid key format";
          }
        }
      } else {
        merklePathError = "Merkle path verification failed: missing required proof fields";
      }
    } catch (error) {
      console.error("Error parsing proof for Merkle path verification:", error);
      merklePathError = `Merkle path verification error: ${error.message}`;
    }

    // For demonstration purposes, we'll also use the hash check as a backup
    // This ensures that some proofs will fail even if our Merkle path verification has issues
    if (merklePathValid && lastChar === '0') {
      merklePathValid = false;
      merklePathError = "Merkle path verification failed: hash check failed";
    }

    // A proof is valid only if:
    // 1. The key exists in the shard
    // 2. The value matches the expected value for the key
    // 3. The Merkle path verification passes
    const isValid = keyExists && valueMatches && merklePathValid;

    // Determine the error message if validation failed
    let errorMessage = null;
    if (!isValid) {
      if (!keyExists) {
        errorMessage = `Key '${proofData.key}' not found in shard ${proofData.shardId}`;
      } else if (!valueMatches) {
        errorMessage = `Value mismatch for key '${proofData.key}' in shard ${proofData.shardId}`;
      } else if (!merklePathValid) {
        errorMessage = merklePathError || "Merkle path verification failed";
      } else {
        errorMessage = "Unknown verification error";
      }
    }

    return {
      valid: isValid,
      error: isValid ? null : errorMessage,
      details: {
        shardId: proofData.shardId,
        key: proofData.key,
        valueVerified: valueMatches,
        keyExists: keyExists,
        merklePathValid: merklePathValid,
        merklePathError: merklePathError,
        hashCheck: lastChar !== '0',
        verificationSteps: [
          { step: "Key existence check", passed: keyExists },
          { step: "Value verification", passed: valueMatches },
          { step: "Merkle path verification", passed: merklePathValid }
        ]
      }
    };
  },

  // Helper function to generate a hash string
  hashString(str) {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      const char = str.charCodeAt(i);
      hash = ((hash << 5) - hash) + char;
      hash = hash & hash; // Convert to 32bit integer
    }
    return hash.toString(16);
  },

  // Generate a deterministic key map for a shard
  generateShardKeyMap(shardId, keyCount) {
    // Create a map to store keys and their values
    const keyMap = new Map();

    // Generate a deterministic set of keys based on the shard ID
    const baseKey = this.hashString(shardId);

    // Add keys to the map
    for (let i = 0; i < keyCount; i++) {
      // Generate a deterministic key based on the shard ID and index
      const key = `key-${baseKey.substring(0, 4)}-${i}`;

      // Generate a deterministic value for this key
      const value = `value-${this.hashString(shardId + key).substring(0, 6)}`;

      // Add to the key map
      keyMap.set(key, value);
    }

    return keyMap;
  },

  // Check if a key exists in a shard
  checkKeyInShard(key, shard) {
    // If the key starts with "key-", check if it exists in the keyMap
    if (key.startsWith("key-")) {
      return shard.keyMap.has(key);
    }

    // For user-provided keys that don't follow our naming convention,
    // use a deterministic approach to decide if they exist
    const keyHash = this.hashString(key + shard.id);
    const keyIndex = parseInt(keyHash.substring(0, 4), 16) % (shard.keys * 2);

    // The key exists if the index is less than the number of keys
    return keyIndex < shard.keys;
  },

  // Get the expected value for a key in a shard
  getValueForKey(key, shard) {
    // If the key exists in the keyMap, return its value
    if (shard.keyMap.has(key)) {
      return shard.keyMap.get(key);
    }

    // For keys that don't exist in the keyMap but should be considered valid
    // (based on checkKeyInShard), generate a deterministic value
    return `value-${this.hashString(shard.id + key).substring(0, 6)}`;
  },

  // Get all keys for a shard
  getShardKeys(shardId) {
    // Get the shard
    const shard = this.blockchain.shards[shardId];
    if (!shard) {
      throw new Error(`Shard ${shardId} not found`);
    }

    // Initialize the shard's keyMap if it doesn't exist
    if (!shard.keyMap) {
      shard.keyMap = this.generateShardKeyMap(shard.id, shard.keys);
    }

    // Return the key-value pairs
    return Array.from(shard.keyMap.entries()).map(([key, value]) => ({ key, value }));
  },

  // Compute the Merkle root from a leaf hash, siblings, and position
  computeMerkleRoot(leafHash, siblings, position) {
    // Start with the leaf hash
    let currentHash = leafHash;
    let currentPosition = position;

    // Combine with siblings to compute the root
    for (let i = 0; i < siblings.length; i++) {
      const siblingHash = siblings[i];

      // Determine if the current hash is a left or right child
      const isRightChild = (currentPosition % 2) === 1;

      // Combine the hashes in the correct order
      let combinedHash;
      if (isRightChild) {
        // If we're a right child, the sibling is on the left
        combinedHash = this.hashString(siblingHash + currentHash);
      } else {
        // If we're a left child, the sibling is on the right
        combinedHash = this.hashString(currentHash + siblingHash);
      }

      // Update the current hash and move up the tree
      currentHash = combinedHash;
      currentPosition = Math.floor(currentPosition / 2);
    }

    // The final hash is the root
    return currentHash;
  },

  // Get blockchain node status
  getStatus() {
    return {
      status: 'running',
      version: '1.0.0',
      uptime: Math.floor((Date.now() - this.blockchain.startTime) / 1000),
      peers: this.blockchain.nodes.length,
      syncStatus: '100%',
      pendingTransactions: this.blockchain.transactionPool.length
    };
  },

  // Process all pending transactions immediately
  processAllPendingTransactions() {
    try {
      if (!this.blockchain || !this.blockchain.transactionPool) {
        console.error('Transaction pool is not initialized');
        return {
          success: false,
          message: "Transaction pool is not initialized. Please refresh the page and try again."
        };
      }

      if (this.blockchain.transactionPool.length === 0) {
        return { success: true, message: "No pending transactions to process" };
      }

      const count = this.blockchain.transactionPool.length;
      console.log(`Manually processing all ${count} pending transactions`);

      // Force process all pending transactions
      this.produceBlock(true);

      return {
        success: true,
        message: `Processed ${count} pending transactions`,
        remainingPending: this.blockchain.transactionPool.length
      };
    } catch (error) {
      console.error('Error processing pending transactions:', error);
      return {
        success: false,
        message: `Error processing transactions: ${error.message}`,
        error: error.toString()
      };
    }
  }
};

// Flag to track initialization status
BlockchainService.initialized = false;

// Add methods to save and load blockchain state
BlockchainService.saveState = function() {
  try {
    const state = {
      blockchain: this.blockchain,
      stats: this.stats,
      amfStats: this.amfStats,
      networkStats: this.networkStats
    };
    localStorage.setItem('blockchainState', JSON.stringify(state));
    console.log('Blockchain state saved to localStorage');
  } catch (error) {
    console.error('Error saving blockchain state:', error);
  }
};

BlockchainService.loadState = function() {
  try {
    const savedState = localStorage.getItem('blockchainState');
    if (savedState) {
      const state = JSON.parse(savedState);
      this.blockchain = state.blockchain;
      this.stats = state.stats;

      // Load AMF stats if available
      if (state.amfStats) {
        this.amfStats = state.amfStats;
      }

      // Load network stats if available
      if (state.networkStats) {
        this.networkStats = state.networkStats;
      }

      console.log('Blockchain state loaded from localStorage with AMF stats:', this.amfStats);
      return true;
    }
    return false;
  } catch (error) {
    console.error('Error loading blockchain state:', error);
    return false;
  }
};

// Original initialize method
const originalInitialize = BlockchainService.initialize;

// Override initialize method to set the initialized flag
BlockchainService.initialize = function() {
  if (this.initialized) {
    console.log('Blockchain service already initialized');
    return;
  }

  console.log('Initializing blockchain service...');

  // Try to load state from localStorage first
  const stateLoaded = this.loadState();

  if (!stateLoaded) {
    // If no saved state, initialize from scratch
    originalInitialize.call(this);
  }

  this.initialized = true;
  console.log('Blockchain service initialization complete');

  // Set up auto-save every 5 seconds
  setInterval(() => {
    this.saveState();
  }, 5000);
};

// Initialize immediately
BlockchainService.initialize();

// Also initialize on DOMContentLoaded as a fallback
document.addEventListener('DOMContentLoaded', () => {
  if (!BlockchainService.initialized) {
    BlockchainService.initialize();
  }
});
