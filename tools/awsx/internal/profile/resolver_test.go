package profile

import (
	"testing"

	"github.com/light-bringer/new-laptop-setup/tools/awsx/internal/config"
)

func TestResolverResolve(t *testing.T) {
	tests := []struct {
		name       string
		resolver   *Resolver
		input      string
		want       ResolvedProfile
		wantErr    bool
	}{
		{
			name:     "known alias resolves",
			resolver: NewResolver(&config.Config{Profiles: map[string]config.ProfileAlias{"nonprod": {AWSProfile: "AdministratorAccess-515216334008", AccountID: "515216334008", Region: "us-east-1"}}}),
			input:    "nonprod",
			want: ResolvedProfile{
				AWSProfileName: "AdministratorAccess-515216334008",
				AccountID:      "515216334008",
				Region:         "us-east-1",
				IsAlias:        true,
			},
		},
		{
			name:     "unknown name falls through",
			resolver: NewResolver(&config.Config{Profiles: map[string]config.ProfileAlias{"nonprod": {AWSProfile: "AdministratorAccess-515216334008", AccountID: "515216334008", Region: "us-east-1"}}}),
			input:    "AdministratorAccess-515216334008",
			want: ResolvedProfile{
				AWSProfileName: "AdministratorAccess-515216334008",
				IsAlias:        false,
			},
		},
		{
			name: "alias with no region inherits defaults",
			resolver: NewResolver(&config.Config{
				Profiles: map[string]config.ProfileAlias{
					"prod": {AWSProfile: "ReadOnly-657536724158", AccountID: "657536724158", Region: ""},
				},
				Defaults: config.Defaults{Region: "us-east-1"},
			}),
			input: "prod",
			want: ResolvedProfile{
				AWSProfileName: "ReadOnly-657536724158",
				AccountID:      "657536724158",
				Region:         "us-east-1",
				IsAlias:        true,
			},
		},
		{
			name:     "nil config",
			resolver: NewResolver(nil),
			input:    "anything",
			want: ResolvedProfile{
				AWSProfileName: "anything",
				IsAlias:        false,
			},
		},
		{
			name:     "empty config",
			resolver: NewResolver(&config.Config{}),
			input:    "anything",
			want: ResolvedProfile{
				AWSProfileName: "anything",
				IsAlias:        false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.resolver.Resolve(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Resolve() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("Resolve() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
