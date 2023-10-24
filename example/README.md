# Mocktimism example [WIP]

This is a simple example fixture of using mocktimism. This example fixture will test minting an nft on l2 as a deposit tx from l1 against a mocktimism devnet

## Dependencies

To use this fixture you will need bun and foundry installed

- [bun](https://bun.sh/)
- [foundry](https://github.com/foundry-rs/foundry)

## Getting started

1. Install node modules

```bash
bun install
```

2. Run test

```bash
bun test
```

Bun will spin up a mocktimism and use [op-viem](https://github.com/base-org/op-viem) to execute a contract mint on l2 from l1

