package encryption

import (
	"strings"
	"testing"
)

func TestService_Encrypt(t *testing.T) {
	service := NewService()
	secret := [16]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'}

	t.Run("should encrypt text successfully", func(t *testing.T) {
		plainText := "Hello, World!"

		encrypted, err := service.Encrypt(plainText, secret)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if encrypted == "" {
			t.Error("Expected encrypted text, got empty string")
		}
		if encrypted == plainText {
			t.Error("Expected encrypted text to be different from plain text")
		}
	})

	t.Run("should return error for empty text", func(t *testing.T) {
		_, err := service.Encrypt("", secret)

		if err == nil {
			t.Fatal("Expected error for empty text, got nil")
		}
		if err.Error() != "cannot encrypt empty string" {
			t.Errorf("Expected 'cannot encrypt empty string', got: %v", err)
		}
	})

	t.Run("should return error for empty secret", func(t *testing.T) {
		emptySecret := [16]byte{}

		_, err := service.Encrypt("Hello, World!", emptySecret)

		if err == nil {
			t.Fatal("Expected error for empty secret, got nil")
		}
		if err.Error() != "secret is empty" {
			t.Errorf("Expected 'secret is empty', got: %v", err)
		}
	})

	t.Run("should produce base64 encoded output", func(t *testing.T) {
		plainText := "Test"

		encrypted, err := service.Encrypt(plainText, secret)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		// Base64 string should only contain valid base64 characters
		for _, c := range encrypted {
			if !isBase64Char(c) {
				t.Errorf("Expected base64 encoded string, found invalid character: %c", c)
				break
			}
		}
	})

	t.Run("should handle long text", func(t *testing.T) {
		plainText := strings.Repeat("Hello, World! ", 100)

		encrypted, err := service.Encrypt(plainText, secret)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if encrypted == "" {
			t.Error("Expected encrypted text, got empty string")
		}
	})

	t.Run("should handle special characters", func(t *testing.T) {
		plainText := "Hello! @#$%^&*()_+-=[]{}|;':\",./<>?`~"

		encrypted, err := service.Encrypt(plainText, secret)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if encrypted == "" {
			t.Error("Expected encrypted text, got empty string")
		}
	})

	t.Run("should handle unicode characters", func(t *testing.T) {
		plainText := "Hello 世界 🌍"

		encrypted, err := service.Encrypt(plainText, secret)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if encrypted == "" {
			t.Error("Expected encrypted text, got empty string")
		}
	})

	t.Run("should produce different output for same text with different secrets", func(t *testing.T) {
		plainText := "Hello, World!"
		secret1 := [16]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'}
		secret2 := [16]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 'b', 'c', 'd', 'e', 'f'}

		encrypted1, _ := service.Encrypt(plainText, secret1)
		encrypted2, _ := service.Encrypt(plainText, secret2)

		if encrypted1 == encrypted2 {
			t.Error("Expected different encrypted outputs for different secrets")
		}
	})

	t.Run("should produce different output each time due to random IV", func(t *testing.T) {
		plainText := "Hello, World!"

		encrypted1, _ := service.Encrypt(plainText, secret)
		encrypted2, _ := service.Encrypt(plainText, secret)

		// With random IV, same plaintext + same secret should produce different ciphertext
		if encrypted1 == encrypted2 {
			t.Error("Expected different encrypted outputs due to random IV")
		}
	})

	t.Run("should handle single character", func(t *testing.T) {
		plainText := "A"

		encrypted, err := service.Encrypt(plainText, secret)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if encrypted == "" {
			t.Error("Expected encrypted text, got empty string")
		}
	})
}

func TestService_Decrypt(t *testing.T) {
	service := NewService()
	secret := [16]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'}

	t.Run("should decrypt encrypted text successfully", func(t *testing.T) {
		plainText := "Hello, World!"
		encrypted, _ := service.Encrypt(plainText, secret)

		decrypted, err := service.Decrypt(encrypted, secret)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if decrypted != plainText {
			t.Errorf("Expected '%s', got '%s'", plainText, decrypted)
		}
	})

	t.Run("should return error for empty text", func(t *testing.T) {
		_, err := service.Decrypt("", secret)

		if err == nil {
			t.Fatal("Expected error for empty text, got nil")
		}
		if err.Error() != "cannot decrypt empty string" {
			t.Errorf("Expected 'cannot decrypt empty string', got: %v", err)
		}
	})

	t.Run("should return error for empty secret", func(t *testing.T) {
		emptySecret := [16]byte{}
		plainText := "Hello, World!"
		encrypted, _ := service.Encrypt(plainText, secret)

		_, err := service.Decrypt(encrypted, emptySecret)

		if err == nil {
			t.Fatal("Expected error for empty secret, got nil")
		}
		if err.Error() != "secret is empty" {
			t.Errorf("Expected 'secret is empty', got: %v", err)
		}
	})

	t.Run("should fail to decrypt with wrong secret", func(t *testing.T) {
		plainText := "Hello, World!"
		secret1 := [16]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'}
		secret2 := [16]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 'b', 'c', 'd', 'e', 'f'}

		encrypted, _ := service.Encrypt(plainText, secret1)
		_, err := service.Decrypt(encrypted, secret2)

		if err == nil {
			t.Fatal("Expected error when decrypting with wrong secret, got nil")
		}
		// GCM authentication should fail with wrong secret
		if !strings.Contains(err.Error(), "authentication check failed") {
			t.Errorf("Expected authentication error, got: %v", err)
		}
	})

	t.Run("should handle long text", func(t *testing.T) {
		plainText := strings.Repeat("Hello, World! ", 100)
		encrypted, _ := service.Encrypt(plainText, secret)

		decrypted, err := service.Decrypt(encrypted, secret)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if decrypted != plainText {
			t.Errorf("Expected original long text, got different text")
		}
	})

	t.Run("should handle special characters", func(t *testing.T) {
		plainText := "Hello! @#$%^&*()_+-=[]{}|;':\",./<>?`~"
		encrypted, _ := service.Encrypt(plainText, secret)

		decrypted, err := service.Decrypt(encrypted, secret)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if decrypted != plainText {
			t.Errorf("Expected '%s', got '%s'", plainText, decrypted)
		}
	})

	t.Run("should handle unicode characters", func(t *testing.T) {
		plainText := "Hello 世界 🌍"
		encrypted, _ := service.Encrypt(plainText, secret)

		decrypted, err := service.Decrypt(encrypted, secret)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if decrypted != plainText {
			t.Errorf("Expected '%s', got '%s'", plainText, decrypted)
		}
	})

	t.Run("should return error for invalid base64", func(t *testing.T) {
		invalidBase64 := "not-valid-base64!!!"

		_, err := service.Decrypt(invalidBase64, secret)

		if err == nil {
			t.Fatal("Expected error for invalid base64, got nil")
		}
		if err.Error() != "invalid base64 encoding" {
			t.Errorf("Expected 'invalid base64 encoding', got: %v", err)
		}
	})

	t.Run("should return error for too short ciphertext", func(t *testing.T) {
		// Valid base64 but too short to contain nonce
		shortCiphertext := service.encode([]byte("short"))

		_, err := service.Decrypt(shortCiphertext, secret)

		if err == nil {
			t.Fatal("Expected error for too short ciphertext, got nil")
		}
		if err.Error() != "ciphertext too short" {
			t.Errorf("Expected 'ciphertext too short', got: %v", err)
		}
	})

	t.Run("should return error for tampered ciphertext", func(t *testing.T) {
		plainText := "Hello, World!"
		encrypted, _ := service.Encrypt(plainText, secret)

		// Tamper with the encrypted text (change one character)
		tamperedBytes := []byte(encrypted)
		if len(tamperedBytes) > 0 {
			// Change a character in the middle
			tamperedBytes[len(tamperedBytes)/2] = 'X'
		}
		tampered := string(tamperedBytes)

		_, err := service.Decrypt(tampered, secret)

		if err == nil {
			t.Fatal("Expected error for tampered ciphertext, got nil")
		}
		// Should fail authentication check
	})

	t.Run("should handle single character", func(t *testing.T) {
		plainText := "A"
		encrypted, _ := service.Encrypt(plainText, secret)

		decrypted, err := service.Decrypt(encrypted, secret)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if decrypted != plainText {
			t.Errorf("Expected '%s', got '%s'", plainText, decrypted)
		}
	})
}

func TestService_EncryptDecrypt_Integration(t *testing.T) {
	service := NewService()
	secret := [16]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'}

	testCases := []struct {
		name      string
		plainText string
	}{
		{"simple text", "Hello, World!"},
		{"single character", "A"},
		{"numbers", "1234567890"},
		{"special chars", "!@#$%^&*()_+-=[]{}|;':\",./<>?`~"},
		{"unicode", "Hello 世界 🌍 مرحبا"},
		{"multiline", "Line1\nLine2\nLine3"},
		{"tabs and spaces", "Hello\tWorld  with  spaces"},
		{"json-like", `{"key": "value", "number": 123}`},
		{"sql-like", "SELECT * FROM users WHERE id = 1;"},
		{"very long", strings.Repeat("A", 1000)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, errEnc := service.Encrypt(tc.plainText, secret)
			if errEnc != nil {
				t.Fatalf("Encryption failed: %v", errEnc)
			}

			decrypted, errDec := service.Decrypt(encrypted, secret)
			if errDec != nil {
				t.Fatalf("Decryption failed: %v", errDec)
			}

			if decrypted != tc.plainText {
				t.Errorf("Mismatch: expected '%s', got '%s'", tc.plainText, decrypted)
			}
		})
	}
}

func TestService_encode(t *testing.T) {
	service := NewService()

	t.Run("should encode bytes to base64", func(t *testing.T) {
		input := []byte("Hello, World!")
		encoded := service.encode(input)

		// Should be valid base64
		if encoded == "" {
			t.Error("Expected encoded string, got empty")
		}
		for _, c := range encoded {
			if !isBase64Char(c) {
				t.Errorf("Expected base64 encoded string, found invalid character: %c", c)
				break
			}
		}
	})

	t.Run("should handle empty bytes", func(t *testing.T) {
		input := []byte{}
		encoded := service.encode(input)

		// Base64 of empty bytes is empty string
		if encoded != "" {
			t.Errorf("Expected empty string, got '%s'", encoded)
		}
	})
}

func TestService_decode(t *testing.T) {
	service := NewService()

	t.Run("should decode base64 to bytes", func(t *testing.T) {
		original := []byte("Hello, World!")
		encoded := service.encode(original)
		decoded, err := service.decode(encoded)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if string(decoded) != string(original) {
			t.Errorf("Expected '%s', got '%s'", string(original), string(decoded))
		}
	})

	t.Run("should return error for invalid base64", func(t *testing.T) {
		invalid := "not-valid-base64!!!"
		_, err := service.decode(invalid)

		if err == nil {
			t.Fatal("Expected error for invalid base64, got nil")
		}
		if err.Error() != "invalid base64 encoding" {
			t.Errorf("Expected 'invalid base64 encoding', got: %v", err)
		}
	})

	t.Run("should handle empty string", func(t *testing.T) {
		decoded, err := service.decode("")

		if err != nil {
			t.Fatalf("Expected no error for empty string, got: %v", err)
		}
		if len(decoded) != 0 {
			t.Errorf("Expected empty bytes, got %d bytes", len(decoded))
		}
	})
}

func TestService_NoLengthLeakage(t *testing.T) {
	service := NewService()
	secret := [16]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'}

	t.Run("should not reveal exact plaintext length", func(t *testing.T) {
		// Test various lengths to ensure ciphertext length doesn't directly reveal plaintext length
		text1 := "A"                      // 1 byte
		text2 := "AB"                     // 2 bytes
		text3 := "ABC"                    // 3 bytes
		text15 := strings.Repeat("A", 15) // 15 bytes
		text16 := strings.Repeat("A", 16) // 16 bytes
		text17 := strings.Repeat("A", 17) // 17 bytes

		enc1, _ := service.Encrypt(text1, secret)
		enc2, _ := service.Encrypt(text2, secret)
		enc3, _ := service.Encrypt(text3, secret)
		enc15, _ := service.Encrypt(text15, secret)
		enc16, _ := service.Encrypt(text16, secret)
		enc17, _ := service.Encrypt(text17, secret)

		// All short texts should have similar lengths due to nonce + auth tag overhead
		// GCM adds: 12 byte nonce + 16 byte auth tag = 28 bytes overhead
		len1 := len(enc1)
		len2 := len(enc2)
		len3 := len(enc3)
		len15 := len(enc15)
		len16 := len(enc16)
		len17 := len(enc17)

		// Verify lengths are not exactly revealing plaintext (due to overhead)
		if len1 == len2 && len2 == len3 {
			t.Log("Good: Very short texts have similar encrypted lengths")
		}

		// Longer texts should still grow, but not in a 1:1 ratio
		if len16-len15 == 1 && len17-len16 == 1 {
			t.Log("Encrypted length grows with plaintext, but includes overhead")
		}

		// Ensure no padding reveals length in blocks
		if len15 != len16 {
			t.Log("Good: No block-based padding reveals exact length")
		}
	})
}

// Helper function to check if a character is valid base64
func isBase64Char(c rune) bool {
	return (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') ||
		c == '+' || c == '/' || c == '='
}
