# Node Configuration

This document describes the configuration options for the blockchain node.

## Command-Line Options

### Node Identity

```bash
-id <node-id>
```

Specifies the ID of the node. Default is `node-1`.

**Example:**

```bash
blockchain -id my-node
```

### Node Address

```bash
-addr <address>
```

Specifies the address of the node. Default is `127.0.0.1:8000`.

**Example:**

```bash
blockchain -addr 192.168.1.100:8000
```

### Private Key File

```bash
-key <file>
```

Specifies the file containing the private key to use for the node. If not provided, a new key pair will be generated.

**Example:**

```bash
blockchain -key key.pem
```

### Shard Split Threshold

```bash
-split <threshold>
```

Specifies the threshold for splitting a shard. Default is `10.0`.

**Example:**

```bash
blockchain -split 20.0
```

### Shard Rebalance Interval

```bash
-rebalance <interval>
```

Specifies the interval for rebalancing the forest. Default is `5m` (5 minutes).

**Example:**

```bash
blockchain -rebalance 10m
```

### Block Pruning Height

```bash
-prune <height>
```

Specifies the height for pruning old blocks. Default is `1000`.

**Example:**

```bash
blockchain -prune 2000
```

### Consensus Interval

```bash
-consensus <interval>
```

Specifies the interval for consensus rounds. Default is `10s` (10 seconds).

**Example:**

```bash
blockchain -consensus 5s
```

### Telemetry Interval

```bash
-telemetry <interval>
```

Specifies the interval for telemetry collection. Default is `1m` (1 minute).

**Example:**

```bash
blockchain -telemetry 30s
```

## API Server Options

### API Address

```bash
-api <address>
```

Specifies the address of the API server. Default is `127.0.0.1:8080`.

**Example:**

```bash
blockchain-api -api 192.168.1.100:8080
```

## Configuration File

In addition to command-line options, the node can be configured using a configuration file. The configuration file is a JSON file with the following structure:

```json
{
  "id": "node-1",
  "address": "127.0.0.1:8000",
  "keyFile": "key.pem",
  "splitThreshold": 10.0,
  "rebalanceInterval": "5m",
  "pruningHeight": 1000,
  "consensusInterval": "10s",
  "telemetryInterval": "1m",
  "api": {
    "address": "127.0.0.1:8080"
  }
}
```

To use a configuration file, specify the `-config` option:

```bash
blockchain -config config.json
```

Command-line options take precedence over configuration file options.
