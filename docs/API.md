# API Documentation

The blockchain system provides a RESTful API for interacting with the blockchain.

## Endpoints

### Transactions

#### Get a Transaction

```
GET /api/transaction?id=<txid>
```

Returns a transaction by ID.

**Response:**

```json
{
  "id": "1234567890abcdef",
  "from": "abcdef1234567890",
  "to": "0987654321fedcba",
  "amount": 100,
  "data": "dGVzdCBkYXRh",
  "timestamp": "2023-01-01T12:00:00Z",
  "shardId": 1,
  "signature": "abcdef1234567890"
}
```

#### Create a Transaction

```
POST /api/transaction
```

Creates a new transaction.

**Request:**

```json
{
  "to": "0987654321fedcba",
  "amount": 100,
  "data": "test data",
  "shardId": 1
}
```

**Response:**

```json
{
  "id": "1234567890abcdef"
}
```

### Blocks

#### Get a Block

```
GET /api/block?hash=<hash>
```

Returns a block by hash.

**Response:**

```json
{
  "header": {
    "version": 1,
    "previousHash": "0000000000000000000000000000000000000000000000000000000000000000",
    "merkleRoot": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "timestamp": "2023-01-01T12:00:00Z",
    "height": 1,
    "difficulty": 1000000,
    "nonce": 12345
  },
  "transactions": [
    {
      "id": "1234567890abcdef",
      "from": "abcdef1234567890",
      "to": "0987654321fedcba",
      "amount": 100,
      "data": "dGVzdCBkYXRh",
      "timestamp": "2023-01-01T12:00:00Z",
      "shardId": 1,
      "signature": "abcdef1234567890"
    }
  ],
  "stateRoot": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
  "shardId": 1,
  "signature": "abcdef1234567890",
  "metadata": {
    "key1": "value1",
    "key2": "value2"
  }
}
```

#### Get the Latest Block

```
GET /api/latest-block
```

Returns the latest block.

**Response:**

Same as the response for `GET /api/block?hash=<hash>`.

### Shards

#### Get a Shard State Value

```
GET /api/shard?id=<shard-id>&key=<key>
```

Returns a shard state value.

**Response:**

```json
{
  "value": "test value"
}
```

### Proofs

#### Generate a Proof

```
GET /api/proof?id=<shard-id>&key=<key>
```

Generates a proof for a shard state value.

**Response:**

```json
{
  "proof": "1234567890abcdef"
}
```

#### Verify a Proof

```
POST /api/proof
```

Verifies a proof for a shard state value.

**Request:**

```json
{
  "shardId": "shard-1",
  "key": "test-key",
  "value": "test-value",
  "proof": "1234567890abcdef"
}
```

**Response:**

```json
{
  "valid": true
}
```

### Status

#### Get the Node Status

```
GET /api/status
```

Returns the node status.

**Response:**

```json
{
  "status": "running",
  "uptime": "1h2m3s",
  "peerCount": 10,
  "connectedPeerCount": 5,
  "consistencyLevel": "strong",
  "networkState": "normal",
  "partitionProbability": 0.1
}
```
