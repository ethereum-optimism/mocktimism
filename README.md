# What is Mocktimism?

Mocktimism is an all-inclusive development tool for OP Stack, similar to hardhat and anvil. It serves as a valuable resource for OP Stack developers in their local and end-to-end development efforts.

Local development and testing are essential components of software development, and Mocktimism plays a key role in the broader **just works** initiative, aimed at ensuring the compatibility of existing Ethereum tools with Optimism.

Developing on the OP Stack comes with its own set of challenges, including:

- L1 Gas fees, which can be surprising for OP stack chain developers.
- Anvil and hardhat lack support for Optimism native bridging.
- Tedious setup when spinning up multiple devnets for applications like Evo-online.
- The Optimism devnet is resource-intensive, relies on Docker, and lacks support for hardhat/anvil features such as impersonation.

To draw an analogy, think of Mocktimism as the **docker-compose** for hardhat and anvil. Just as Docker is used to configure a single container, and docker-compose configures multiple containers, Mocktimism is the tool that OP-chain developers use to configure multiple chains when working with anvil and hardhat.
