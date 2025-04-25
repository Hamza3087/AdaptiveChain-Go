package crypto

import (
	"crypto/sha256"
	"errors"
	"math/big"
)

// Accumulator represents a cryptographic accumulator
// This is a simplified RSA accumulator implementation
type Accumulator struct {
	Value *big.Int
	N     *big.Int // RSA modulus
}

// NewAccumulator creates a new accumulator with the given RSA modulus
func NewAccumulator(n *big.Int) *Accumulator {
	// Start with a value of 1
	return &Accumulator{
		Value: big.NewInt(1),
		N:     n,
	}
}

// Add adds an element to the accumulator
func (a *Accumulator) Add(element []byte) {
	// Hash the element to get a prime representative
	prime := hashToPrime(element)

	// Update the accumulator value: value^prime mod N
	a.Value.Exp(a.Value, prime, a.N)
}

// VerifyMembership verifies that an element is in the accumulator
func (a *Accumulator) VerifyMembership(element []byte, witness *big.Int) bool {
	// Hash the element to get a prime representative
	prime := hashToPrime(element)

	// Compute witness^prime mod N
	result := new(big.Int).Exp(witness, prime, a.N)

	// Check if the result equals the accumulator value
	return result.Cmp(a.Value) == 0
}

// GenerateWitness generates a membership witness for an element
func (a *Accumulator) GenerateWitness(element []byte, allElements [][]byte) (*big.Int, error) {
	// Find the element in the list
	found := false
	for _, e := range allElements {
		if string(e) == string(element) {
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("element not in accumulator")
	}

	// Compute the product of all primes except the one for the element
	product := big.NewInt(1)
	for _, e := range allElements {
		if string(e) != string(element) {
			prime := hashToPrime(e)
			product.Mul(product, prime)
		}
	}

	// Compute the witness: value^(1/product) mod N
	// This is a simplified version and not secure for real use
	// In practice, we would need to compute the modular inverse
	witness := new(big.Int).Exp(a.Value, new(big.Int).ModInverse(product, a.N), a.N)

	return witness, nil
}

// hashToPrime hashes a byte slice to a prime number
// This is a simplified implementation for demonstration purposes
func hashToPrime(data []byte) *big.Int {
	// Hash the data
	hash := sha256.Sum256(data)

	// Convert the hash to a big integer
	num := new(big.Int).SetBytes(hash[:])

	// Ensure the number is odd
	if num.Bit(0) == 0 {
		num.Add(num, big.NewInt(1))
	}

	// Find the next prime
	for !num.ProbablyPrime(20) {
		num.Add(num, big.NewInt(2))
	}

	return num
}

// BatchAdd adds multiple elements to the accumulator
func (a *Accumulator) BatchAdd(elements [][]byte) {
	for _, element := range elements {
		a.Add(element)
	}
}

// BatchVerify verifies multiple elements in the accumulator
func (a *Accumulator) BatchVerify(elements [][]byte, witnesses []*big.Int) bool {
	if len(elements) != len(witnesses) {
		return false
	}

	for i, element := range elements {
		if !a.VerifyMembership(element, witnesses[i]) {
			return false
		}
	}

	return true
}

// GetValue returns the current value of the accumulator
func (a *Accumulator) GetValue() *big.Int {
	return new(big.Int).Set(a.Value)
}
