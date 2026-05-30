package processor

import (
	"slices"
	"strings"

	"github.com/AGG-Programming/LeagueSpectator/internal/league"
	"github.com/AGG-Programming/LeagueSpectator/pkg/models"
)

func (p *Processor) getTeamObjectives(events []league.Event, blueTeam []string, redTeam []string) {
	bluePlayers := p.removeHashtagsInPlace(blueTeam)
	redPlayers := p.removeHashtagsInPlace(redTeam)

	for _, event := range events {
		if event.EventID == 0 && p.LastProcessedEventID > 0 {
			p.Reset()
		}

		if event.EventID <= p.LastProcessedEventID {
			continue
		}

		switch event.EventName {
		case "TurretKilled":
			{
				//TODO: Does not work yet
				if slices.Contains(bluePlayers, *event.KillerName) {
					p.BlueObjectives = append(p.BlueObjectives, models.Objective{
						Key:   event.EventName,
						Icon:  "",
						Kills: nil, //TODO: calculate total turret kills
					})
				} else if slices.Contains(redPlayers, *event.KillerName) {
					p.RedObjectives = append(p.RedObjectives, models.Objective{
						Key:   event.EventName,
						Icon:  "",
						Kills: nil, //TODO: calculate total turret kills
					})
				}
			}
		case "DragonKill":
			{
				//TODO: Order key is adjusted for every drake. Key should stay per drake
				if slices.Contains(bluePlayers, *event.KillerName) {
					p.BlueDrakeOrderKey++
					p.BlueObjectives = append(p.BlueObjectives, models.Objective{
						Key:      *event.DragonType,
						Icon:     "",
						OrderKey: &p.BlueDrakeOrderKey,
					})
				} else if slices.Contains(redPlayers, *event.KillerName) {
					p.RedDrakeOrderKey++
					p.RedObjectives = append(p.RedObjectives, models.Objective{
						Key:      *event.DragonType,
						Icon:     "",
						OrderKey: &p.RedDrakeOrderKey,
					})
				}
			}
		case "HeraldKill":
			{
				if slices.Contains(bluePlayers, *event.KillerName) {
					p.BlueObjectives = append(p.BlueObjectives, models.Objective{
						Key:  event.EventName,
						Icon: "",
					})
				} else if slices.Contains(redPlayers, *event.KillerName) {
					p.RedObjectives = append(p.RedObjectives, models.Objective{
						Key:  event.EventName,
						Icon: "",
					})
				}
			}
		case "BaronKill":
			{
				isActive := false //TODO: call function
				remaining := 0.0  //TODO: call function
				if slices.Contains(bluePlayers, *event.KillerName) {
					p.BlueObjectives = append(p.BlueObjectives, models.Objective{
						Key:           event.EventName,
						Icon:          "",
						IsActive:      &isActive,
						RemainingTime: &remaining,
					})
				} else if slices.Contains(redPlayers, *event.KillerName) {
					p.RedObjectives = append(p.RedObjectives, models.Objective{
						Key:           event.EventName,
						Icon:          "",
						IsActive:      &isActive,
						RemainingTime: &remaining,
					})
				}
			}
		case "HordeKill":
			{
				//TODO: Do not create a new objective for every grub. edit one per team
				if slices.Contains(bluePlayers, *event.KillerName) {
					p.BlueGrubKills++
					kills := p.BlueGrubKills
					p.BlueObjectives = append(p.BlueObjectives, models.Objective{
						Key:   event.EventName,
						Icon:  "",
						Kills: &kills,
					})
				} else if slices.Contains(redPlayers, *event.KillerName) {
					p.RedGrubKills++
					kills := p.RedGrubKills
					p.RedObjectives = append(p.RedObjectives, models.Objective{
						Key:   event.EventName,
						Icon:  "",
						Kills: &kills,
					})
				}
			}
		}

		p.LastProcessedEventID = event.EventID
	}
}

func (p *Processor) removeHashtagsInPlace(tags []string) []string {
	var newTags []string
	for i, str := range tags {
		parts := strings.SplitN(str, "#", 2)
		tags[i] = parts[0]
		newTags = append(newTags, tags[i])
	}
	return newTags
}

func (p *Processor) Reset() {
	p.LastProcessedEventID = -1
	p.BlueObjectives = []models.Objective{}
	p.RedObjectives = []models.Objective{}
	p.BlueDrakeOrderKey = 0
	p.RedDrakeOrderKey = 0
	p.BlueGrubKills = 0
	p.RedGrubKills = 0
}
