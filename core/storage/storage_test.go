package storage

import "testing"

func TestValidateContent(t *testing.T) {
	png := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52}
	if err := ValidateContent(png); err != nil {
		t.Fatalf("valid PNG rejected: %v", err)
	}
	if err := ValidateContent([]byte("not an image")); err == nil {
		t.Fatal("non-image content should be rejected")
	}
}
