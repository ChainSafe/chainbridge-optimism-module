# Sygma optimism module
<a href="https://discord.gg/ykXsJKfhgq">
  <img alt="discord" src="https://img.shields.io/discord/593655374469660673?label=Discord&logo=discord&style=flat" />
</a>

Sygma optimism module is the part of Sygma-core framework. This module brings support of optimism compatible client module.

*Project still in deep beta*
- Chat with us on [discord](https://discord.gg/ykXsJKfhgq).

### Table of Contents

1. [Installation](#installation)
2. [Usage](#usage)
3. [Differences Between EVM and Celo](#differences-between-evm-and-celo)

## Installation
Refer to [installation](https://github.com/ChainSafe/chainbridge-docs/blob/develop/docs/installation.md) guide for assistance in installing.

## Usage
Module should be used along with core [framework](https://github.com/ChainSafe/chainbridge-core).

Since sygma-optimism-module is a package it will require writing some extra code to get it running alongside [chainbridge-core](https://github.com/ChainSafe/chainbridge-core). Here you can find some examples
[Example](https://github.com/ChainSafe/chainbridge-core-example)

### Differences Between EVM and Optimism module

Though Optimism is an EVM-compatible chain, it needs additional checks to verify transaction batches before submitting votes, and therefore is deserving of its own separate module.

The differences pertain to `OptimismClient` which returns only verified blocks when querying latest block.

# ChainSafe Security Policy

## Reporting a Security Bug

We take all security issues seriously, if you believe you have found a security issue within a ChainSafe
project please notify us immediately. If an issue is confirmed, we will take all necessary precautions
to ensure a statement and patch release is made in a timely manner.

Please email us a description of the flaw and any related information (e.g. reproduction steps, version) to
[security at chainsafe dot io](mailto:security@chainsafe.io).

## License

_GNU Lesser General Public License v3.0_
