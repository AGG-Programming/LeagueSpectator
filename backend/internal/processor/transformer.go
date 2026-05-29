package processor

import (
	"slices"

	"github.com/AGG-Programming/LeagueSpectator/internal/league"
	"github.com/AGG-Programming/LeagueSpectator/pkg/models"
)

const BlueTeam = "ORDER"

type Cache interface {
	GetChampion(id string) string
	GetRune(id int) string
	GetItem(id string) string
	GetSpell(id string) string
}

type Processor struct {
	cache Cache
}

func NewProcessor(cache Cache) *Processor {
	return &Processor{
		cache: cache,
	}
}

func (p *Processor) Transformer(data league.GameData) (models.DynamicGameData, error) {
	var blueTeam, redTeam []string
	for _, player := range data.Players {
		if player.Team == BlueTeam {
			blueTeam = append(blueTeam, player.RiotID)
		} else {
			redTeam = append(redTeam, player.RiotID)
		}
	}

	blueScore, redScore := p.getTeamScore(data.Players)
	blueObjectives, redObjectives := p.getTeamObjectives(data.Events.Events, blueTeam, redTeam)
	bluePlayers, redPlayers := p.getPlayers(data.Players)

	output := models.DynamicGameData{
		BlueTeam: models.Team{
			Score:      blueScore,
			Gold:       0, //TODO: calculate team gold
			Objectives: blueObjectives,
			Players:    bluePlayers,
		},
		RedTeam: models.Team{
			Score:      redScore,
			Gold:       0, //TODO: calculate team gold
			Objectives: redObjectives,
			Players:    redPlayers,
		},
		Timers:   []models.Timer{},
		GameTime: data.GameTime,
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

func (p *Processor) getTeamObjectives(events []league.Event, blueTeam []string, redTeam []string) ([]models.Objective, []models.Objective) {
	var blueDrakeOrderKey, redDrakeOrderKey int
	var blueObjectives, redObjectives []models.Objective
	for _, event := range events {
		switch event.EventName {
		case "TurretKilled", "InhibKilled":
			{
				if slices.Contains(blueTeam, *event.KillerName) {
					blueObjectives = append(blueObjectives, models.Objective{
						Key:   event.EventName,
						Icon:  "",
						Kills: nil, //TODO: calculate all teams turret kills
					})
				} else if slices.Contains(redTeam, *event.KillerName) {
					redObjectives = append(redObjectives, models.Objective{
						Key:   event.EventName,
						Icon:  "",
						Kills: nil, //TODO: calculate all teams turret kills
					})
				}
			}
		case "DragonKill":
			{
				if slices.Contains(blueTeam, *event.KillerName) {
					blueDrakeOrderKey++
					blueObjectives = append(blueObjectives, models.Objective{
						Key:      *event.DragonType,
						Icon:     "",
						OrderKey: &blueDrakeOrderKey,
					})
				} else if slices.Contains(redTeam, *event.KillerName) {
					redDrakeOrderKey++
					redObjectives = append(redObjectives, models.Objective{
						Key:      *event.DragonType,
						Icon:     "",
						OrderKey: &redDrakeOrderKey,
					})
				}
			}
		case "HeraldKill":
			{
				if slices.Contains(blueTeam, *event.KillerName) {
					blueObjectives = append(blueObjectives, models.Objective{
						Key:  event.EventName,
						Icon: "",
					})
				} else if slices.Contains(redTeam, *event.KillerName) {
					redObjectives = append(redObjectives, models.Objective{
						Key:  event.EventName,
						Icon: "",
					})
				}
			}
		case "BaronKill":
			{
				isActive := false //TODO: call function
				remaining := 0.0  //TODO: call function
				if slices.Contains(blueTeam, *event.KillerName) {
					blueObjectives = append(blueObjectives, models.Objective{
						Key:           event.EventName,
						Icon:          "",
						IsActive:      &isActive,
						RemainingTime: &remaining,
					})
				} else if slices.Contains(redTeam, *event.KillerName) {
					redObjectives = append(redObjectives, models.Objective{
						Key:           event.EventName,
						Icon:          "",
						IsActive:      &isActive,
						RemainingTime: &remaining,
					})
				}
			}
		case "GrubKill":
			{
				//TODO: calculate grub kills
			}
		}
	}
	return blueObjectives, redObjectives
}

func (p *Processor) getPlayers(players []league.Player) ([]models.Player, []models.Player) {
	var bluePlayers, redPlayers []models.Player
	for _, player := range players {
		pl := models.Player{
			ChampionName:    player.ChampionName,
			Icon:            p.cache.GetChampion(player.ChampionName),
			IsDead:          player.IsDead,
			Level:           player.Level,
			Position:        player.Position,
			RespawnTimer:    player.RespawnTimer,
			RiotId:          player.RiotID,
			PlayerTotalGold: 0, //TODO: Calculate Gold
			Runes: models.Runes{
				Keystone: models.Rune{
					DisplayName: player.Runes.Keystone.DisplayName,
					Icon:        p.cache.GetRune(player.Runes.Keystone.ID),
				},
				Primary: models.Rune{
					DisplayName: player.Runes.Primary.DisplayName,
					Icon:        p.cache.GetRune(player.Runes.Primary.ID),
				},
				Secondary: models.Rune{
					DisplayName: player.Runes.Secondary.DisplayName,
					Icon:        p.cache.GetRune(player.Runes.Secondary.ID),
				},
			},
			Items: nil, //TODO:
			Scores: models.Scores{
				Assists:    player.Scores.Assists,
				CreepScore: player.Scores.CreepScore, //TODO: check if float64
				Deaths:     player.Scores.Deaths,
				Kills:      player.Scores.Kills,
				WardScore:  player.Scores.WardScore,
			},
			Spells: nil, //TODO:
		}

		if player.Team == BlueTeam {
			bluePlayers = append(bluePlayers, pl)
		} else {
			redPlayers = append(redPlayers, pl)
		}
	}
	return bluePlayers, redPlayers
}
