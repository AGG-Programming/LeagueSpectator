package processor

import (
	"github.com/AGG-Programming/LeagueSpectator/internal/league"
	"github.com/AGG-Programming/LeagueSpectator/internal/pl"
	"github.com/AGG-Programming/LeagueSpectator/pkg/models"
)

const BlueTeam = "ORDER"

type Cache interface {
	GetChampion(id string) string
	GetRune(id int) string
	GetItem(id int) string
	GetSpell(id string) string
	GetUlt(id string) string
}

type Processor struct {
	cache                Cache
	LastProcessedEventID int
	BlueObjectives       []models.Objective
	RedObjectives        []models.Objective
	BlueDrakeOrderKey    int
	RedDrakeOrderKey     int
	BlueGrubKills        int
	RedGrubKills         int
}

func NewProcessor(cache Cache) *Processor {
	return &Processor{
		cache: cache,
	}
}

func (p *Processor) Transformer(data league.GameResponse) (models.DynamicGameData, error) {
	var blueTeam, redTeam []string
	for _, player := range data.Players {
		if player.Team == BlueTeam {
			blueTeam = append(blueTeam, player.RiotID)
		} else {
			redTeam = append(redTeam, player.RiotID)
		}
	}

	blueScore, redScore := p.getTeamScore(data.Players)
	bluePlayers, redPlayers := p.getPlayers(data.Players)

	p.getTeamObjectives(data.Events.Events, blueTeam, redTeam)

	output := models.DynamicGameData{
		BlueTeam: models.Team{
			Score:      blueScore,
			Objectives: p.BlueObjectives,
			Players:    bluePlayers,
		},
		RedTeam: models.Team{
			Score:      redScore,
			Objectives: p.RedObjectives,
			Players:    redPlayers,
		},
		Timers:   []models.Timer{},
		GameTime: data.GameData.GameTime,
	}
	return output, nil
}

func (p *Processor) getTeamScore(players []league.Player) (int, int) {
	var blueScore, redScore int
	for _, player := range players {
		if player.Team == BlueTeam {
			blueScore = blueScore + player.Scores.Kills
		} else {
			redScore = redScore + player.Scores.Kills
		}
	}
	return blueScore, redScore
}

func (p *Processor) TransformPL(data pl.PrimeLeagueResponse, targetID int) (*models.PrimeLeague, error) {
	teams, err := data.GetTeamStandings(targetID)
	if err != nil {
		return nil, err
	}

	output := models.PrimeLeague{
		GroupTitle: data.GroupTitle,
		TargetTeam: models.PrimeTeam{
			Tag:      teams.Target.Tag,
			Wins:     teams.Target.Wins,
			Losses:   teams.Target.Losses,
			Points:   teams.Target.Points,
			Position: teams.Target.Position,
			Img:      teams.Target.Img,
		},
		LeadingTeam: models.PrimeTeam{
			Tag:      teams.Leading.Tag,
			Wins:     teams.Leading.Wins,
			Losses:   teams.Leading.Losses,
			Points:   teams.Leading.Points,
			Position: teams.Leading.Position,
			Img:      teams.Leading.Img,
		},
		TrailingTeam: models.PrimeTeam{
			Tag:      teams.Trailing.Tag,
			Wins:     teams.Trailing.Wins,
			Losses:   teams.Trailing.Losses,
			Points:   teams.Trailing.Points,
			Position: teams.Trailing.Position,
			Img:      teams.Trailing.Img,
		},
		LastTeam: models.PrimeTeam{
			Tag:      teams.Last.Tag,
			Wins:     teams.Last.Wins,
			Losses:   teams.Last.Losses,
			Points:   teams.Last.Points,
			Position: teams.Last.Position,
			Img:      teams.Last.Img,
		},
	}
	return &output, nil
}
