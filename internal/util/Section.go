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
	ActionableAccounts     []string
	ActionableEnvironments []string
	ActionableApplications []string
	Accounts               []sectionio.AccountGraph
}

type UpdateBlocklistInput struct {
	ctx       context.Context
	Log       zerolog.Logger
	AccountId int64
	AppId     int64
	EnvName   string
	Ips       sectionio.IpRestrictions
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
					in := &UpdateBlocklistInput{
						ctx:       s.Auth,
						AccountId: int64(account.Id),
						AppId:     int64(app.Id),
						EnvName:   env.EnvironmentName,
						Ips:       ips,
					}
					in.Log = s.Logger.With().
						Int64("account", in.AccountId).
						Int64("app", in.AppId).
						Str("env", in.EnvName).
						Logger()
					go func(in *UpdateBlocklistInput) {
						in.Log.Info().Msgf("adding %d ips to blocklist", len(in.Ips.IpBlacklist), in.AppId, in.EnvName)
						_, resp, err := s.Client.EnvironmentApi.EnvironmentIpRestrictionsPost(in.ctx, in.AccountId, in.AppId, in.EnvName, in.Ips)
						if err != nil {
							in.Log.Err(err).Str("statusCode", resp.Status).Msg("failed to add Ips to blocklist")
							return
						}
						in.Log.Info().Msg("successfully updated ip restrictions")
					}(in)
				}
			}
		}
	}
}
