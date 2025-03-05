# Radius SDKs

Official software development kits for building applications on [Radius](https://radiustech.xyz/), a high-performance
smart contract platform that enables near-instant settlement and can process millions of transactions per second.

## Overview

Next-generation payments require processing power that is orders of magnitude more efficient than what is currently
available. Built by the team that brought USDC to market, Radius is the result of many years rethinking smart contract
scalability from first principles.

Unlike blockchains that sequentially process a limited batch of transactions at a time, our distributed execution layer
handles multiple transactions simultaneously. Our platform has demonstrated that it can process over 2.8 million
transactions per second with near zero latency and cost, far exceeding any other system that exists today.

Radius is fully EVM-compatible and provides SDKs in [multiple programming languages](#available-sdks) with a clean,
consistent interface, enabling developers to easily add instant payments to their apps with just a few lines of code.

No block times. No bidding wars. Just instant settlement.

## Available SDKs

- [Go SDK](go/README.md)
- [Python SDK](python/README.md) (coming soon)
- [Rust SDK](rust/README.md) (coming soon)
- [TypeScript SDK](typescript/README.md)

## Use Cases

Radius is capable of handling millions of micro-payments per second at a cost that makes doing so economically viable.
This is particularly well-suited for AI agent use cases, and equally so for any application that requires massive scale,
instant settlement, and cryptographic guarantees.

### AI Payments
- AI agents buying products and services in real-time
- Pay-per-API-call data access at fractions of a cent
- Pay-per-compute, storage, bandwidth request
- High-frequency trading settlement

### Traditional Payments
- High-frequency trading and settlement
- Pay-per-use services and subscriptions
- Real-time revenue sharing and splits

### Beyond Payments
- Decentralized social networks and content systems
- Gaming and virtual world state management
- IoT sensor networks and data marketplaces
- Identity and attestation systems

## Available Today

The Radius invite-only testnet launched in January 2025 with major stablecoin issuers and AI labs already onboard.
Radius supports both simple payments and other EVM-compatible smart contracts, so developers can experience the
efficiency of its parallel execution design.

The next trillion transactions won't come from humans typing on keyboards. They'll come from AI agents making
split-second decisions. We're building the infrastructure necessary to make that future possible.

Ready to build the future? Start [here](https://docs.radiustech.xyz/radius-testnet-access).

## Contributing

We welcome contributions to all Radius SDKs! Please see:

- [General Contributing Guide](CONTRIBUTING.md) - Repository-wide guidelines and principles
- [Go SDK Contributing Guide](go/CONTRIBUTING.md) - Go-specific guidelines
- [TypeScript SDK Contributing Guide](typescript/CONTRIBUTING.md) - TypeScript-specific guidelines

## Support

- [Website](https://radiustech.xyz/)
- [Testnet Access](https://docs.radiustech.xyz/radius-testnet-access)
- [GitHub Issues](https://github.com/radiustechsystems/sdk/issues)

## License

All Radius SDKs are released under the [MIT License](LICENSE).
