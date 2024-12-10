package util

import (
	"context"
	"strconv"

	sectionio "github.com/dpc-sdp/go-section-io"
	"github.com/rs/zerolog"
)

type Section struct {
	Auth                   context.Context
	Client                 *sectionio.APIClient
	BlockedIps             sectionio.IpRestrictions
	Logger                 zerolog.Logger
	IPTracker              *IPTracker
	ActionableAccounts     []string
	ActionableEnvironments []string
	ActionableApplications []string
	Accounts               []sectionio.AccountGraph
}

func (s *Section) Init() (bool, error) {
	apps, _, err := s.Client.AccountApi.AccountGraph(s.Auth)
	s.Logger.Debug().Msg("fetching account graph")

	if err != nil {
		s.Logger.Debug().Err(err)
		return false, err
	}

	s.Accounts = apps
	s.Logger.Info().Str("status", "success").Msg("initialised application list")

	return true, nil
}

func (s *Section) IsActionableAccount(id int32) bool {
	if len(s.ActionableAccounts) == 0 || (len(s.ActionableAccounts) == 1 && s.ActionableAccounts[0] == "") {
		return true
	}

	for _, a := range s.ActionableAccounts {
		actionableId, _ := strconv.ParseInt(a, 10, 32)
		if id == int32(actionableId) {
			return true
		}
	}

	return false
}

func (s *Section) IsActionableApplication(name string) bool {
	if len(s.ActionableApplications) == 0 || (len(s.ActionableApplications) == 1 && s.ActionableApplications[0] == "") {
		return true
	}

	for _, a := range s.ActionableApplications {
		if name == a {
			return true
		}
	}

	return false
}

func (s *Section) IsActionableEnvironment(name string) bool {
	if len(s.ActionableEnvironments) == 0 || (len(s.ActionableEnvironments) == 1 && s.ActionableEnvironments[0] == "") {
		return true
	}

	for _, e := range s.ActionableEnvironments {
		if name == e {
			return true
		}
	}

	return false
}

func (s *Section) AddIpRestrictionsToAllApplications(ips sectionio.IpRestrictions) {
	for _, account := range s.Accounts {
		if !s.IsActionableAccount(account.Id) {
			s.Logger.Debug().Int32("account", account.Id).Msg("account exclusion matched")
			continue
		}

		for _, app := range account.Applications {
			if !s.IsActionableApplication(app.ApplicationName) {
				s.Logger.Debug().Str("name", app.ApplicationName).Msg("application exclusion matched")
				continue
			}
			for _, env := range app.Environments {
				if s.IsActionableEnvironment(env.EnvironmentName) {
					go func(ctx context.Context, accountId int64, app sectionio.AccountGraphApplications, envName string, ipRestrction sectionio.IpRestrictions, l zerolog.Logger) {
						for _, ip := range ips.IpBlacklist {
							s.IPTracker.TrackIP(ip)
						}
						ips.IpBlacklist = append(ips.IpBlacklist, s.IPTracker.BackoffIPs()...)
						_, _, err := s.Client.EnvironmentApi.EnvironmentIpRestrictionsPost(ctx, accountId, int64(app.Id), envName, ips)
						if err != nil {
							l.Error().Err(err)
							return
						}
						l.Info().Int64("account", accountId).Strs("ips", ips.IpBlacklist).Str("env", envName).Str("app", app.ApplicationName).Msg("successfully updated ip restrictions")
					}(s.Auth, int64(account.Id), app, env.EnvironmentName, ips, s.Logger)
				}
			}
		}
	}
}
