package playvs

import (
	"encoding/json"
	"time"
)

type client struct{}

func CreateClient() client {
	return client{}
}

// TODO examine other libraries to see if this is how they handle their magic constants

type region string

const (
	EasternRegion region = "17a567ac-c0cb-401f-85de-619f84bcb75b"
)

type metaSeason string

const (
	MetaSeason metaSeason = "95c742a7-8f9c-4417-a459-8c5b930d79c5"
)

type getter struct {
	region     region
	metaseason metaSeason
}

func (c client) GetRegion(region region) getter {
	return getter{
		region:     region,
		metaseason: MetaSeason,
	}
}

type getAllLeagueTeamsResult struct {
	Data struct {
		GetTeams struct {
			Teams []struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				State  string `json:"state"`
				Esport struct {
					ID       string `json:"id"`
					Rating   string `json:"rating"`
					Typename string `json:"__typename"`
				} `json:"esport"`
				School struct {
					ID       string `json:"id"`
					Name     string `json:"name"`
					LogoURL  string `json:"logoUrl"`
					Slug     string `json:"slug"`
					Typename string `json:"__typename"`
				} `json:"school"`
				Typename string `json:"__typename"`
			} `json:"teams"`
			TotalCount int    `json:"totalCount"`
			Typename   string `json:"__typename"`
		} `json:"getTeams"`
	} `json:"data"`
	Extensions struct {
		TraceID string `json:"traceId"`
	} `json:"extensions"`
}

type team struct {
	ID    string
	Name  string
	State string
}

const (
	playVsEndpoint = "https://api.playvs.com/graphql"
)

func (g getter) Teams() ([]*team, error) {
	payload := map[string]interface{}{
		"operationName": "getAllLeagueTeams",
		"variables": map[string]interface{}{
			"filters": map[string]interface{}{
				"metaseasonId": g.metaseason,
				"leagueId":     g.region,
				"esportSlugs":  []string{"league-of-legends"},
				"keyword":      "",
			},
			// TODO approximately 75 teams total
			"limit":  100,
			"offset": 0,
		},
		"extensions": map[string]interface{}{
			"persistedQuery": map[string]interface{}{
				"version":    1,
				"sha256Hash": "fdd6c95ee9f8ea96a45a87ab89f822ebfa41f3eb4348e4e4e595733aa7cbb570",
			},
		},
	}

	result, err := performRequest("POST", playVsEndpoint, payload)
	if err != nil {
		return []*team{}, err
	}

	var leagueTeams getAllLeagueTeamsResult
	if err := json.Unmarshal(result, &leagueTeams); err != nil {
		return []*team{}, err
	}

	var teams []*team

	for _, leagueTeam := range leagueTeams.Data.GetTeams.Teams {
		teams = append(teams, &team{
			ID:    leagueTeam.ID,
			Name:  leagueTeam.Name,
			State: leagueTeam.State,
		})
	}

	return teams, nil
}

type teamRosterResult struct {
	Errors []struct {
		Message   string `json:"message"`
		Locations []struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"locations"`
		Path       []any `json:"path"`
		Extensions struct {
			ErrorType     string `json:"errorType"`
			OriginalError struct {
				JseShortmsg string `json:"jse_shortmsg"`
				JseInfo     struct {
				} `json:"jse_info"`
				Message    string `json:"message"`
				Extensions struct {
					Code string `json:"code"`
				} `json:"extensions"`
			} `json:"originalError"`
			Code      string `json:"code"`
			Exception struct {
				JseShortmsg string `json:"jse_shortmsg"`
				JseInfo     struct {
				} `json:"jse_info"`
				Message string `json:"message"`
			} `json:"exception"`
		} `json:"extensions"`
	} `json:"errors"`
	Data struct {
		Team struct {
			ID     string `json:"id"`
			School struct {
				ID               string `json:"id"`
				CompetitionGroup string `json:"competitionGroup"`
				Name             string `json:"name"`
				LogoURL          string `json:"logoUrl"`
				Typename         string `json:"__typename"`
			} `json:"school"`
			IsPlayerLed bool `json:"isPlayerLed"`
			IsHidden    bool `json:"isHidden"`
			Esport      struct {
				ID       string `json:"id"`
				Slug     string `json:"slug"`
				Rating   string `json:"rating"`
				Typename string `json:"__typename"`
			} `json:"esport"`
			EnrolledSeasons []struct {
				ID                          string    `json:"id"`
				RostersLockAt               time.Time `json:"rostersLockAt"`
				Name                        string    `json:"name"`
				StartsAt                    time.Time `json:"startsAt"`
				EndsAt                      time.Time `json:"endsAt"`
				TeamRegistrationEndsAt      time.Time `json:"teamRegistrationEndsAt"`
				SuggestedRegistrationEndsAt time.Time `json:"suggestedRegistrationEndsAt"`
				TeamDeregistrationEndsAt    time.Time `json:"teamDeregistrationEndsAt"`
				Phases                      []struct {
					ID       string    `json:"id"`
					Type     string    `json:"type"`
					StartsAt time.Time `json:"startsAt"`
					EndsAt   time.Time `json:"endsAt"`
					Name     string    `json:"name"`
					Status   string    `json:"status"`
					Typename string    `json:"__typename"`
				} `json:"phases"`
				Metaseason struct {
					ID         string    `json:"id"`
					IsPromoted bool      `json:"isPromoted"`
					StartsAt   time.Time `json:"startsAt"`
					Name       string    `json:"name"`
					EndsAt     time.Time `json:"endsAt"`
					Typename   string    `json:"__typename"`
				} `json:"metaseason"`
				Typename string `json:"__typename"`
			} `json:"enrolledSeasons"`
			Leagues []struct {
				ID               string `json:"id"`
				CompetitionModel string `json:"competitionModel"`
				Name             string `json:"name"`
				DisplayName      string `json:"displayName"`
				EsportID         string `json:"esportId"`
				Typename         string `json:"__typename"`
			} `json:"leagues"`
			Coaches []struct {
				ID                   string    `json:"id"`
				Name                 string    `json:"name"`
				FirstName            string    `json:"firstName"`
				LastName             string    `json:"lastName"`
				LastSeen             time.Time `json:"lastSeen"`
				Email                any       `json:"email"`
				Phone                any       `json:"phone"`
				PhoneExt             any       `json:"phoneExt"`
				IsPhoneNumberVisible bool      `json:"isPhoneNumberVisible"`
				AvatarURL            any       `json:"avatarUrl"`
				Roles                []struct {
					UserID       string `json:"userId"`
					RoleName     string `json:"roleName"`
					ResourceID   string `json:"resourceId"`
					ResourceType string `json:"resourceType"`
					Typename     string `json:"__typename"`
				} `json:"roles"`
				Typename string `json:"__typename"`
			} `json:"coaches"`
			Name   string `json:"name"`
			Roster struct {
				ID            string `json:"id"`
				TeamID        string `json:"teamId"`
				NumPlayers    int    `json:"numPlayers"`
				MaxNumPlayers int    `json:"maxNumPlayers"`
				Players       []struct {
					IsCaptain   bool      `json:"isCaptain"`
					EffectiveAt time.Time `json:"effectiveAt"`
					User        struct {
						ID                  string    `json:"id"`
						LastSeen            time.Time `json:"lastSeen"`
						Name                string    `json:"name"`
						AvatarURL           any       `json:"avatarUrl"`
						UserEsportPlatforms []any     `json:"userEsportPlatforms"`
						Roles               []struct {
							UserID       string `json:"userId"`
							RoleName     string `json:"roleName"`
							ResourceType string `json:"resourceType"`
							ResourceID   string `json:"resourceId"`
							Typename     string `json:"__typename"`
						} `json:"roles"`
						SchoolRoleStatus []struct {
							Status   string `json:"status"`
							Typename string `json:"__typename"`
						} `json:"schoolRoleStatus"`
						Typename string `json:"__typename"`
					} `json:"user"`
					Typename string `json:"__typename"`
				} `json:"players"`
				Formats []struct {
					TeamSize int `json:"teamSize"`
					Starters []struct {
						Position struct {
							ID           string `json:"id"`
							Index        int    `json:"index"`
							Name         string `json:"name"`
							Abbreviation string `json:"abbreviation"`
							Colloquial   string `json:"colloquial"`
							Typename     string `json:"__typename"`
						} `json:"position"`
						Player struct {
							EffectiveAt time.Time `json:"effectiveAt"`
							IsCaptain   bool      `json:"isCaptain"`
							User        struct {
								ID           string    `json:"id"`
								Name         string    `json:"name"`
								LastSeen     time.Time `json:"lastSeen"`
								AvatarURL    string    `json:"avatarUrl"`
								SeasonPasses []struct {
									ID           string    `json:"id"`
									LeagueID     string    `json:"leagueId"`
									MetaseasonID string    `json:"metaseasonId"`
									ConsumedAt   time.Time `json:"consumedAt"`
									Typename     string    `json:"__typename"`
								} `json:"seasonPasses"`
								Roles []struct {
									UserID       string `json:"userId"`
									RoleName     string `json:"roleName"`
									ResourceType string `json:"resourceType"`
									ResourceID   string `json:"resourceId"`
									Typename     string `json:"__typename"`
								} `json:"roles"`
								SchoolRoleStatus []struct {
									Status   string `json:"status"`
									Typename string `json:"__typename"`
								} `json:"schoolRoleStatus"`
								UserProviderAccounts []struct {
									ID                      string `json:"id"`
									UserID                  string `json:"userId"`
									ProviderName            string `json:"providerName"`
									ProviderAccountID       string `json:"providerAccountId"`
									ProviderDisplayName     string `json:"providerDisplayName"`
									ProviderIntegrationType string `json:"providerIntegrationType"`
									Typename                string `json:"__typename"`
								} `json:"userProviderAccounts"`
								Typename string `json:"__typename"`
							} `json:"user"`
							Typename string `json:"__typename"`
						} `json:"player"`
						Typename string `json:"__typename"`
					} `json:"starters"`
					Substitutes []struct {
						IsCaptain bool `json:"isCaptain"`
						User      struct {
							ID           string    `json:"id"`
							Name         string    `json:"name"`
							AvatarURL    any       `json:"avatarUrl"`
							LastSeen     time.Time `json:"lastSeen"`
							SeasonPasses []struct {
								ID           string `json:"id"`
								LeagueID     string `json:"leagueId"`
								MetaseasonID string `json:"metaseasonId"`
								ConsumedAt   any    `json:"consumedAt"`
								Typename     string `json:"__typename"`
							} `json:"seasonPasses"`
							Roles []struct {
								UserID       string `json:"userId"`
								RoleName     string `json:"roleName"`
								ResourceType string `json:"resourceType"`
								ResourceID   string `json:"resourceId"`
								Typename     string `json:"__typename"`
							} `json:"roles"`
							SchoolRoleStatus []struct {
								Status   string `json:"status"`
								Typename string `json:"__typename"`
							} `json:"schoolRoleStatus"`
							UserProviderAccounts []any  `json:"userProviderAccounts"`
							Typename             string `json:"__typename"`
						} `json:"user"`
						Typename string `json:"__typename"`
					} `json:"substitutes"`
					Typename string `json:"__typename"`
				} `json:"formats"`
				Typename string `json:"__typename"`
			} `json:"roster"`
			Typename string `json:"__typename"`
		} `json:"team"`
	} `json:"data"`
	Extensions struct {
		TraceID string `json:"traceId"`
	} `json:"extensions"`
}

type teamGetter struct {
	getter getter
	teamId string
}

func (g getter) Get(team *team) teamGetter {
	return teamGetter{
		getter: g,
		teamId: team.ID,
	}
}

func (g getter) GetTeam(teamId string) teamGetter {
	return teamGetter{
		getter: g,
		teamId: teamId,
	}
}

func (tg teamGetter) DisplayNames() ([]string, error) {
	payload := map[string]interface{}{
		"operationName": "teamRoster",
		"variables": map[string]interface{}{
			"id":                         tg.teamId,
			"metaseasonId":               tg.getter.metaseason,
			"includeSlotExclusionsField": false,
			"isPublic":                   false,
			"isCoach":                    false,
		},
		"extensions": map[string]interface{}{
			"persistedQuery": map[string]interface{}{
				"version":    1,
				"sha256Hash": "3b1fa794463895123f7179a73165c3732fe3a6dc5138b6a7b6276a6f8c0619fa",
			},
		},
	}

	result, err := performRequest("POST", playVsEndpoint, payload)
	if err != nil {
		return []string{}, err
	}

	var roster teamRosterResult
	if err := json.Unmarshal(result, &roster); err != nil {
		return []string{}, err
	}

	var displayNames []string

	for _, format := range roster.Data.Team.Roster.Formats {
		for _, starter := range format.Starters {
			for _, account := range starter.Player.User.UserProviderAccounts {
				// TODO Move to constant
				if account.ProviderName == "Riot" {
					displayNames = append(displayNames, account.ProviderDisplayName)
				}
			}
		}

		// TODO Substitute struct not defined; probably the same as starter's, but not sure
		// for _, substitute := range format.Substitutes {
		// for _, account := range substitute.User.UserProviderAccounts {
		// playerNames = append(playerNames, account.ProviderDisplayName)
		// }
		// }
	}

	return displayNames, nil
}
