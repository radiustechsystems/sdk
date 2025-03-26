// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "@openzeppelin/contracts/token/ERC1155/ERC1155.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

contract AccessTokenSystem is ERC1155, Ownable {
    using Strings for uint256;
    using ECDSA for bytes32;

    // Access tier data
    struct AccessTier {
        uint256 price;
        uint256 ttl;
        bool active;
    }

    // Token expiration tracking
    mapping(address => mapping(uint256 => uint256)) public expiresAt;

    // Bit-mapping for revocations (save gas by using a single uint256)
    // Each bit represents revocation status for a specific tier
    mapping(address => uint256) public revocations;

    // Mapping for tiers
    mapping(uint256 => AccessTier) public tiers;

    // Domain separator for EIP-712 signatures
    bytes32 public immutable DOMAIN_SEPARATOR;

    // Events
    event AccessPurchased(address indexed consumer, uint256 indexed tierId, uint256 expiryTime);
    event BatchAccessPurchased(address indexed consumer, uint256[] tierIds, uint256[] expiryTimes);
    event AccessRevoked(address indexed consumer, uint256 indexed tierId);
    event TierCreated(uint256 indexed tierId, uint256 price, uint256 ttl);
    event TierStatusChanged(uint256 indexed tierId, bool active);

    constructor(string memory _uri) ERC1155(_uri) Ownable(msg.sender) {
        // Create domain separator for signature verification
        DOMAIN_SEPARATOR = keccak256(
            abi.encode(
                keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
                keccak256(bytes("AccessTokenSystem")),
                keccak256(bytes("1")),
                block.chainid,
                address(this)
            )
        );
    }

    /**
     * @dev Create a new access tier
     * @param tierId The ID for this tier
     * @param price Price in wei
     * @param ttl Time to live in seconds
     * @param active Whether this tier is available for purchase
     */
    function createTier(
        uint256 tierId,
        uint256 price,
        uint256 ttl,
        bool active
	) external onlyOwner {
        tiers[tierId] = AccessTier(price, ttl, active);
        emit TierCreated(tierId, price, ttl);
    }

    /**
     * @dev Set the active status of a tier
     */
    function setTierStatus(uint256 tierId, bool active) external onlyOwner {
        tiers[tierId].active = active;
        emit TierStatusChanged(tierId, active);
    }

    /**
     * @dev Purchase access to a single tier
     */
    function purchaseAccess(uint256 tierId) external payable {
        AccessTier memory tier = tiers[tierId];
        require(tier.active, "Tier not available");
        require(msg.value >= tier.price, "Insufficient payment");

        // Set expiration
        uint256 expiry = block.timestamp + tier.ttl;
        expiresAt[msg.sender][tierId] = expiry;

        // Clear any previous revocation for this tier
        if (_isRevoked(msg.sender, tierId)) {
            _toggleRevocation(msg.sender, tierId);
        }

        // Mint token
        _mint(msg.sender, tierId, 1, "");

        // Process payment
        uint256 refund = msg.value - tier.price;
        if (refund > 0) {
            payable(msg.sender).transfer(refund);
        }
        payable(owner()).transfer(tier.price);

        emit AccessPurchased(msg.sender, tierId, expiry);
    }

    /**
     * @dev Purchase access to multiple tiers in one transaction
     */
    function batchPurchaseAccess(uint256[] calldata tierIds) external payable {
        uint256 totalCost = 0;
        uint256[] memory expiryTimes = new uint256[](tierIds.length);
        uint256[] memory amounts = new uint256[](tierIds.length);

        // Calculate total cost and validate tiers
        for (uint256 i = 0; i < tierIds.length; i++) {
            uint256 tierId = tierIds[i];
            AccessTier memory tier = tiers[tierId];
            require(tier.active, "Tier not available");

            totalCost += tier.price;
            expiryTimes[i] = block.timestamp + tier.ttl;
            amounts[i] = 1;

            // Set expiration
            expiresAt[msg.sender][tierId] = expiryTimes[i];

            // Clear any previous revocation
            if (_isRevoked(msg.sender, tierId)) {
                _toggleRevocation(msg.sender, tierId);
            }
        }

        require(msg.value >= totalCost, "Insufficient payment");

        // Mint tokens
        _mintBatch(msg.sender, tierIds, amounts, "");

        // Process payment
        uint256 refund = msg.value - totalCost;
        if (refund > 0) {
            payable(msg.sender).transfer(refund);
        }
        payable(owner()).transfer(totalCost);

        emit BatchAccessPurchased(msg.sender, tierIds, expiryTimes);
    }

    /**
     * @dev Check if access is valid for a user and tier
     */
    function isValid(address user, uint256 tierId) public view returns (bool) {
        return
            balanceOf(user, tierId) > 0 &&
            expiresAt[user][tierId] > block.timestamp &&
            !_isRevoked(user, tierId);
    }

    /**
     * @dev Check if access has been revoked
     */
    function _isRevoked(address user, uint256 tierId) internal view returns (bool) {
        // Use bit operations to check if the bit at position tierId is set
        return (revocations[user] & (1 << (tierId % 256))) != 0;
    }

    /**
     * @dev Toggle revocation status for a tier
     */
    function _toggleRevocation(address user, uint256 tierId) internal {
        // Use XOR to toggle the bit at position tierId
        revocations[user] ^= (1 << (tierId % 256));
    }

    /**
     * @dev Revoke access for a user and tier
     */
    function revokeAccess(address user, uint256 tierId) external onlyOwner {
        require(balanceOf(user, tierId) > 0, "No access token");
        if (!_isRevoked(user, tierId)) {
            _toggleRevocation(user, tierId);
        }
        emit AccessRevoked(user, tierId);
    }

	/**
	 * @dev Verify a signed message from a user to validate access off-chain
	 * @param user Address of the user
	 * @param tierId ID of the access tier
	 * @param challenge A string challenge to prevent replay attacks
	 * @param signature Signature from the user
	 */
	function verifyAccess(
		address user,
		uint256 tierId,
		string calldata challenge,
		bytes calldata signature
	) external view returns (bool) {
		// First check if access is valid
		if (!isValid(user, tierId)) {
			return false;
		}

		// Create the message hash
		bytes32 messageHash = keccak256(abi.encodePacked(
			"\x19Ethereum Signed Message:\n",
			Strings.toString(bytes(challenge).length),
			challenge
		));

		// Recover signer from signature
		bytes32 r;
		bytes32 s;
		uint8 v;

		if (signature.length != 65) {
			return false;
		}

		assembly {
			r := calldataload(add(signature.offset, 0))
			s := calldataload(add(signature.offset, 32))
			v := byte(0, calldataload(add(signature.offset, 64)))
		}

		if (v < 27) {
			v += 27;
		}

		address recoveredSigner = ecrecover(messageHash, v, r, s);

		// Verify the signer is the user
		return recoveredSigner == user;
	}

    /**
     * @dev Override batch transfer function to transfer expiration data
     */
    function safeBatchTransferFrom(
        address from,
        address to,
        uint256[] memory ids,
        uint256[] memory amounts,
        bytes memory data
    ) public override {
        super.safeBatchTransferFrom(from, to, ids, amounts, data);

        // Transfer expiration data for each token
        for (uint256 i = 0; i < ids.length; i++) {
            uint256 id = ids[i];
            uint256 amount = amounts[i];

            if (amount > 0 && expiresAt[from][id] > 0) {
                expiresAt[to][id] = expiresAt[from][id];
                expiresAt[from][id] = 0;

                // Transfer revocation status
                if (_isRevoked(from, id)) {
                    if (!_isRevoked(to, id)) {
                        _toggleRevocation(to, id);
                    }
                    if (_isRevoked(from, id)) {
                        _toggleRevocation(from, id);
                    }
                }
            }
        }
    }
}
