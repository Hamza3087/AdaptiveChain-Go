package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// KeyPair represents a cryptographic key pair for digital signatures
type KeyPair struct {
	privateKey *ecdsa.PrivateKey
}

// GenerateKeyPair creates a new ECDSA key pair
func GenerateKeyPair() (*KeyPair, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}
	return &KeyPair{privateKey: privateKey}, nil
}

// ExportPublicKey exports the public key as a byte array
func (kp *KeyPair) ExportPublicKey() ([]byte, error) {
	if kp.privateKey == nil {
		return nil, errors.New("private key is nil")
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&kp.privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	return pem.EncodeToMemory(pemBlock), nil
}

// ExportPrivateKey exports the private key as a byte array
func (kp *KeyPair) ExportPrivateKey() ([]byte, error) {
	if kp.privateKey == nil {
		return nil, errors.New("private key is nil")
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(kp.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	pemBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	return pem.EncodeToMemory(pemBlock), nil
}

// ImportPrivateKey imports a private key from a byte array
func ImportPrivateKey(privateKeyBytes []byte) (*KeyPair, error) {
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &KeyPair{privateKey: privateKey}, nil
}

// ImportPublicKey imports a public key from a byte array
func ImportPublicKey(publicKeyBytes []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	genericPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey, ok := genericPublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not an ECDSA public key")
	}

	return publicKey, nil
}

// Sign signs a message with the private key
func (kp *KeyPair) Sign(message []byte) ([]byte, error) {
	if kp.privateKey == nil {
		return nil, errors.New("private key is nil")
	}

	hash := sha256.Sum256(message)

	r, s, err := ecdsa.Sign(rand.Reader, kp.privateKey, hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}

	// Combine r and s into a single signature
	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

// VerifySignature verifies a signature with a public key
func VerifySignature(publicKeyBytes []byte, message []byte, signature []byte) (bool, error) {
	publicKey, err := ImportPublicKey(publicKeyBytes)
	if err != nil {
		return false, fmt.Errorf("failed to import public key: %w", err)
	}

	hash := sha256.Sum256(message)

	// Split signature into r and s
	sigLen := len(signature)
	if sigLen%2 != 0 {
		return false, errors.New("invalid signature length")
	}

	rBytes := signature[:sigLen/2]
	sBytes := signature[sigLen/2:]

	var r, s ecdsa.PublicKey
	r.X.SetBytes(rBytes)
	s.Y.SetBytes(sBytes)

	return ecdsa.Verify(publicKey, hash[:], r.X, s.Y), nil
}
