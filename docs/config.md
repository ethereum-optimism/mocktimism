# Configuration

Mocktimism can be configured via cli flags or a `mocktimism.toml` file.

## Table of Contents
- [Global Configuration](#global-configuration)
- [Chain Configuration](#chain-configuration)
- [Anvil Options](#anvil-options) 
---

## Example TOML
Below is an example of a `mocktimism.toml` configuration file:

```toml
[profile.default]
state = "/path/to/state"
silent = false

# l1 chain
[[profile.default.chains]]
id = "mainnet"
base_chain_id = "mainnet"

# Fork options
fork_chain_id = 1
fork_url = "https://mainnet.alchemy.infura.io"
block_base_fee_per_gas = 420

# Chain options
chain_id = 10
gas_limit = 420

# EVM options
accounts = 10
balance = 1000
steps-tracing = true

# Server options
allow-origin = "*"
port = 8545
host = "127.0.0.1"
block_time = 12
prune_history = false

# l2 chain
[[profile.default.chains]]
id = "optimism"
base_chain_id = "mainnet"

# Fork options
fork_chain_id = 10
fork_url = "https://op.alchemy.infura.io"
block_base_fee_per_gas = 420

# Chain options
chain_id = 10
gas_limit = 420

# EVM options
accounts = 10
balance = 1000
steps-tracing = true

# Server options
allow-origin = "*"
port = 8546
host = "127.0.0.1"
block_time = 2
prune_history = false
```

## Global Configuration
The global configuration options are:

- `state`: Path to the directory where Mocktimism will store its state.
- `silent`: A boolean indicating whether Mocktimism should run in silent mode.

## Chain Configuration
Chains are defined under `profile.default.chains`. Each chain has its own configuration options:

- `id`: A unique identifier for the chain.
- `base_chain_id`: The ID of the chain that this chain is based on.

### Fork options
Options related to the fork of the chain:

- `fork_chain_id`: The ID of the chain to fork.
- `fork_url`: URL of the chain to fork.
- `block_base_fee_per_gas`: The base fee per gas for the block.

### Chain options
Options related to the operation of the chain:

- `chain_id`: A unique identifier for the chain.
- `gas_limit`: The gas limit for the chain.

### EVM options
Options related to the Ethereum Virtual Machine (EVM):

- `accounts`: Number of accounts in the EVM.
- `balance`: The balance for each account in the EVM.
- `steps-tracing`: A boolean indicating whether tracing of steps in the EVM is enabled.

### Server options
Options related to the Mocktimism server:

- `allow-origin`: Allowed origin for cross-origin requests.
- `port`: Port on which the server will listen.
- `host`: Host on which the server will run.
- `block_time`: Time in seconds between blocks.
- `prune_history`: A boolean indicating whether the history should be pruned.

