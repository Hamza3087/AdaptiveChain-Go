# CLI Documentation

The blockchain system provides a command-line interface (CLI) for interacting with the blockchain.

## Commands

### Generate a Key Pair

```bash
blockchain-cli generate-key <file>
```

Generates a new key pair and saves the private key to the specified file.

**Example:**

```bash
blockchain-cli generate-key key.pem
```

### Get a Transaction

```bash
blockchain-cli get-transaction <txid>
```

Gets a transaction by ID.

**Example:**

```bash
blockchain-cli get-transaction 1234567890abcdef
```

### Get a Block

```bash
blockchain-cli get-block <hash>
```

Gets a block by hash.

**Example:**

```bash
blockchain-cli get-block 1234567890abcdef
```

### Get the Latest Block

```bash
blockchain-cli get-latest-block
```

Gets the latest block.

**Example:**

```bash
blockchain-cli get-latest-block
```

### Create a Transaction

```bash
blockchain-cli -key <key-file> create-transaction <to> <amount> <data> <shard-id>
```

Creates a new transaction.

**Example:**

```bash
blockchain-cli -key key.pem create-transaction 0987654321fedcba 100 "test data" 1
```

### Get a Shard State Value

```bash
blockchain-cli get-shard-state <shard-id> <key>
```

Gets a shard state value.

**Example:**

```bash
blockchain-cli get-shard-state shard-1 test-key
```

### Verify a Proof

```bash
blockchain-cli verify-proof <shard-id> <key> <value> <proof>
```

Verifies a proof for a shard state value.

**Example:**

```bash
blockchain-cli verify-proof shard-1 test-key test-value 1234567890abcdef
```

## Options

### Node Address

```bash
-node <address>
```

Specifies the address of the node to connect to. Default is `127.0.0.1:8000`.

**Example:**

```bash
blockchain-cli -node 192.168.1.100:8000 get-latest-block
```

### Private Key File

```bash
-key <file>
```

Specifies the file containing the private key to use for signing transactions.

**Example:**

```bash
blockchain-cli -key key.pem create-transaction 0987654321fedcba 100 "test data" 1
```
