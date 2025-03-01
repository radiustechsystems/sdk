# Changelog
All notable changes to the Radius TypeScript SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Type exports for `BigNumberish`, `BytesLike`, and `HttpClient`, required for custom `Signer` implementations

### Security
- Added pnpm override for `esbuild` requiring v0.25.0 or above, to address a known security vulnerability

## 1.0.0
### Added
- Initial SDK implementation
- Radius client creation
- Account management
- Contract deployment and interaction
