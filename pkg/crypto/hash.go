package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
)

// Hash represents a 32-byte SHA-256 hash
type Hash [32]byte

// NewHash creates a new hash from a byte slice
func NewHash(data []byte) Hash {
	return sha256.Sum256(data)
}

// ZeroHash returns a hash with all zeros
func ZeroHash() Hash {
	return Hash{}
}

// String returns the hash as a hexadecimal string
func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

// Bytes returns the hash as a byte slice
func (h Hash) Bytes() []byte {
	return h[:]
}

// FromString creates a hash from a hexadecimal string
func HashFromString(s string) (Hash, error) {
	var hash Hash
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return hash, err
	}

	if len(bytes) != 32 {
		return hash, errors.New("invalid hash length")
	}

	copy(hash[:], bytes)
	return hash, nil
}

// Equal checks if two hashes are equal
func (h Hash) Equal(other Hash) bool {
	return h == other
}

// MerkleHash computes the Merkle hash of two child hashes
func MerkleHash(left, right Hash) Hash {
	// Concatenate the two hashes and compute a new hash
	var data [64]byte
	copy(data[:32], left[:])
	copy(data[32:], right[:])
	return sha256.Sum256(data[:])
}

// ComputeMerkleRoot computes the Merkle root of a list of hashes
func ComputeMerkleRoot(hashes []Hash) Hash {
	if len(hashes) == 0 {
		return ZeroHash()
	}

	if len(hashes) == 1 {
		return hashes[0]
	}

	// If the number of hashes is odd, duplicate the last hash
	if len(hashes)%2 != 0 {
		hashes = append(hashes, hashes[len(hashes)-1])
	}

	// Compute the next level of the Merkle tree
	var nextLevel []Hash
	for i := 0; i < len(hashes); i += 2 {
		nextLevel = append(nextLevel, MerkleHash(hashes[i], hashes[i+1]))
	}

	// Recursively compute the Merkle root of the next level
	return ComputeMerkleRoot(nextLevel)
}

// RandomFloat generates a random float between 0 and 1
func RandomFloat() float64 {
	// Generate a random 64-bit integer
	nBig, err := rand.Int(rand.Reader, big.NewInt(1<<63-1))
	if err != nil {
		// In case of error, return a default value
		return 0.5
	}

	// Convert to float64 and normalize to [0, 1]
	return float64(nBig.Int64()) / float64(1<<63-1)
}
