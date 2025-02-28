# Radius SDK Contract Examples

This repository contains examples of how to interact with smart contracts using the Radius SDK.
The ABI and bytecode of the smart contracts are included in the examples, and were generated using the
[solcjs](https://www.npmjs.com/package/solc) compiler.

```sh
solcjs SimpleStorage.sol --bin --abi --optimize
solcjs StableCoin.sol --bin --abi --optimize --base-path . --include-path node_modules/
```
