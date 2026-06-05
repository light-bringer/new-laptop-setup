package profile

import "github.com/light-bringer/new-laptop-setup/tools/awsx/internal/config"

type ResolvedProfile struct {
	AWSProfileName string
	AccountID      string
	Region         string
	IsAlias        bool
}

type Resolver struct {
	cfg *config.Config
}

func NewResolver(cfg *config.Config) *Resolver {
	return &Resolver{cfg: cfg}
}

func (r *Resolver) Resolve(name string) (ResolvedProfile, error) {
	if r == nil || r.cfg == nil || len(r.cfg.Profiles) == 0 {
		return ResolvedProfile{AWSProfileName: name}, nil
	}

	alias, ok := r.cfg.Profiles[name]
	if !ok {
		return ResolvedProfile{AWSProfileName: name}, nil
	}

	return ResolvedProfile{
		AWSProfileName: alias.AWSProfile,
		AccountID:      alias.AccountID,
		Region:         alias.Region,
		IsAlias:        true,
	}, nil
}
