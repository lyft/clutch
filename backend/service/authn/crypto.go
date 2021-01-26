package authn

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// AES-256 encryption.
type cryptographer struct {
	gcm cipher.AEAD
}

func newCryptographer(passphrase string) (*cryptographer, error) {
	if len(passphrase) == 0 {
		return nil, fmt.Errorf("passphrase was empty")
	}

	// Hash the passphrase with SHA-256. Alternatively we could have required a 64-byte hex string.
	// The other option would be a 32-byte key, but entropy with a 32-byte user-provided key is likely to be lower since
	// use of non-printable characters is unlikely.
	key := sha256.Sum256([]byte(passphrase))

	cb, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(cb)
	if err != nil {
		return nil, err
	}

	return &cryptographer{gcm: gcm}, nil
}

// Return decrypted bytes using input with nonce as the prefix and ciphertext as remaining bytes.
func (c *cryptographer) Decrypt(b []byte) ([]byte, error) {
	// Ensure we don't panic when splitting the nonce and ciphertext.
	if len(b) < c.gcm.NonceSize() {
		return nil, fmt.Errorf("invalid nonce+cipher, bytes for decryption are smaller than algorithm's nonce size")
	}

	nonce, ciphertext := b[:c.gcm.NonceSize()], b[c.gcm.NonceSize():]

	out, err := c.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// Return an encrypted input bytes with nonce as the prefix.
func (c *cryptographer) Encrypt(b []byte) ([]byte, error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("input bytes were empty, could not encrypt")
	}

	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("error generating nonce: %w", err)
	}

	// Prefix output with nonce, allocate capacity using provided overhead.
	out := make([]byte, c.gcm.NonceSize(), c.gcm.NonceSize()+c.gcm.Overhead()+len(b))
	copy(out, nonce)

	return c.gcm.Seal(out, nonce, b, nil), nil
}
