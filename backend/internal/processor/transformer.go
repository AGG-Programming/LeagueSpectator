package processor

import (
	"github.com/AGG-Programming/LeagueSpectator/internal/league"
	"github.com/AGG-Programming/LeagueSpectator/pkg/models"
)

const BlueTeam = "ORDER"

type Cache interface {
	GetChampion(id string) string
	GetRune(id int) string
	GetItem(id int) string
	GetSpell(id string) string
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
