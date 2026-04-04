package validator

import "testing"

func TestVaildataURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantUrl bool
	}{
		{"valid https", "https://example.com", true},
		{"valid http", "http://example.com", true},
		{"valid with path", "https://example.com/path/to/page", true},
		{"missing scheme", "example.com", false},
		{"ftp scheme", "ftp://example.com", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err != nil) == tt.wantUrl {
				t.Errorf("ValidateURL(%q) error = %v, want error: %v", tt.url, err, tt.wantUrl)
			}
		})
	}
}
