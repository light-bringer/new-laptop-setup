package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "valid yaml",
			fn: func(t *testing.T) {
				cfg, err := Load(filepath.Join("..", "..", "testdata", "awsx.yaml"))
				if err != nil {
					t.Fatalf("Load() error = %v", err)
				}

				if cfg == nil {
					t.Fatal("Load() returned nil config")
				}

				if got := cfg.Defaults.Region; got != "us-east-1" {
					t.Fatalf("Defaults.Region = %q, want %q", got, "us-east-1")
				}

				checks := map[string]ProfileAlias{
					"nonprod": {
						AWSProfile:  "AdministratorAccess-515216334008",
						AccountID:   "515216334008",
						Region:      "us-east-1",
						Description: "AppSec Common NonProd - Admin",
					},
					"nonprod-ro": {
						AWSProfile:  "ReadOnly-515216334008",
						AccountID:   "515216334008",
						Region:      "us-east-1",
						Description: "AppSec Common NonProd - ReadOnly",
					},
					"prod": {
						AWSProfile:  "ReadOnly-657536724158",
						AccountID:   "657536724158",
						Region:      "us-east-1",
						Description: "AppSec Common Prod - ReadOnly",
					},
				}

				for name, want := range checks {
					got, ok := cfg.Profiles[name]
					if !ok {
						t.Fatalf("Profiles[%q] missing", name)
					}
					if got != want {
						t.Fatalf("Profiles[%q] = %#v, want %#v", name, got, want)
					}
				}
			},
		},
		{
			name: "missing file",
			fn: func(t *testing.T) {
				cfg, err := Load("/nonexistent/path.yaml")
				if err != nil {
					t.Fatalf("Load() error = %v", err)
				}
				if cfg == nil {
					t.Fatal("Load() returned nil config")
				}
				if cfg.Profiles == nil {
					t.Fatal("Load() returned nil Profiles map")
				}
				if len(cfg.Profiles) != 0 {
					t.Fatalf("len(Profiles) = %d, want 0", len(cfg.Profiles))
				}
				if cfg.Defaults != (Defaults{}) {
					t.Fatalf("Defaults = %#v, want zero value", cfg.Defaults)
				}
			},
		},
		{
			name: "empty file",
			fn: func(t *testing.T) {
				f, err := os.CreateTemp(t.TempDir(), "empty-*.yaml")
				if err != nil {
					t.Fatalf("CreateTemp() error = %v", err)
				}
				if err := f.Close(); err != nil {
					t.Fatalf("Close() error = %v", err)
				}

				cfg, err := Load(f.Name())
				if err != nil {
					t.Fatalf("Load() error = %v", err)
				}
				if cfg == nil {
					t.Fatal("Load() returned nil config")
				}
				if cfg.Profiles == nil {
					t.Fatal("Load() returned nil Profiles map")
				}
				if len(cfg.Profiles) != 0 {
					t.Fatalf("len(Profiles) = %d, want 0", len(cfg.Profiles))
				}
			},
		},
		{
			name: "malformed yaml",
			fn: func(t *testing.T) {
				f, err := os.CreateTemp(t.TempDir(), "bad-*.yaml")
				if err != nil {
					t.Fatalf("CreateTemp() error = %v", err)
				}
				if _, err := f.WriteString("{invalid: [yaml"); err != nil {
					t.Fatalf("WriteString() error = %v", err)
				}
				if err := f.Close(); err != nil {
					t.Fatalf("Close() error = %v", err)
				}

				if _, err := Load(f.Name()); err == nil {
					t.Fatal("Load() error = nil, want error")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.fn)
	}
}
