# Radius SDK Contract Examples

This directory contains examples of smart contracts that can be deployed to Radius. They are for demonstration
purposes only and are not intended for production use.

## Compile Script

The ABI and bytecode of the smart contracts are included in the examples, and were generated using the
[solcjs](https://www.npmjs.com/package/solc) compiler. For convenience, a shell script is provided to make compilation even easier.

```sh
./compile.sh SimpleStorage.sol
```

This generates the ABI and bytecode files in the current directory, overwriting existing files if present.

If a contract imports other contracts, you can use the `--no-deps` flag to avoid generating ABI and bytecode files for
those imported contracts:

```sh
./compile.sh AccessTokenSystem.sol --no-deps
```

## Dependencies

Some of the contracts in this directory import other contracts from libraries like OpenZeppelin. To compile these
contracts, you first need to install the dependencies defined in `package.json` by running:

```sh
npm install
```
