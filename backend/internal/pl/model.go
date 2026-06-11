package pl

import (
	"errors"
	"fmt"
	"sort"
)

type PrimeLeagueResponse struct {
	GroupTitle string `json:"group_title"`
	Ranking    struct {
		Pages map[string][]PageItem `json:"_pages"`
	} `json:"_ranking"`
}

type PageItem struct {
	TeamID   int    `json:"team_id"`
	Wins     int    `json:"ssp_stats_wins"`
	Losses   int    `json:"ssp_stats_losses"`
	Points   int    `json:"ssp_stats_points"`
	Position int    `json:"ssp_stats_rank"`
	Tag      string `json:"_short"`
	Img      string `json:"_img"`
}

type TeamStandings struct {
	Target   PageItem
	Leading  PageItem
	Trailing PageItem
	Last     PageItem
}

func (p *PrimeLeagueResponse) GetTeamStandings(targetID int) (*TeamStandings, error) {
	var result TeamStandings

	for _, teams := range p.Ranking.Pages {
		if len(teams) == 0 {
			return nil, errors.New("no teams found in response data")
		}

		sort.Slice(teams, func(i, j int) bool {
			return teams[i].Position < teams[j].Position
		})

		targetIdx := -1
		for i, team := range teams {
			if team.TeamID == targetID {
				targetIdx = i
				break
			}
		}

		if targetIdx == -1 {
			return nil, fmt.Errorf("target team ID %d not found", targetID)
		}

		result = TeamStandings{
			Target: teams[targetIdx],
			Last:   teams[len(teams)-1],
		}

		if targetIdx > 0 {
			result.Leading = teams[targetIdx-1]
		}

		if targetIdx < len(teams)-1 {
			result.Trailing = teams[targetIdx+1]
		}
	}
	return &result, nil
}
