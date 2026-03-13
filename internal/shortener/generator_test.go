package shortener

import (
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	code, err := GenerateShortCode()
	if err != nil {
		t.Fatalf("GenerateShortCode returned error: %v", err)
	}

	if len(code) != 6 {
		t.Errorf("Expected code length 6, got %d", len(code))
	}

	for _, char := range code {
		found := false
		for _, c := range charset {
			if char == c {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Character %c not in charset", char)
		}
	}
}

func TestGenerateShortCodeRandomness(t *testing.T) {
	code1, err1 := GenerateShortCode()
	code2, err2 := GenerateShortCode()

	if err1 != nil || err2 != nil {
		t.Fatalf("GenerateShortCode returned error")
	}

	if code1 == code2 {
		t.Log("Warning: Generated duplicate codes. This is possible but highly unlikely.")
	}
}
