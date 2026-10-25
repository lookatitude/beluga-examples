package main

import (
	"errors"
	"testing"
)

func TestAPIKeyFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    string
		wantErr error
	}{
		{name: "unset", value: "", wantErr: errMissingAPIKey},
		{name: "whitespace only", value: "  \t", wantErr: errMissingAPIKey},
		{name: "set", value: "sk-test", want: "sk-test"},
		{name: "trimmed", value: " sk-test\n", want: "sk-test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getenv := func(key string) string {
				if key != "OPENAI_API_KEY" {
					t.Fatalf("unexpected env lookup %q", key)
				}
				return tt.value
			}
			got, err := apiKeyFromEnv(getenv)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("key = %q, want %q", got, tt.want)
			}
		})
	}
}
