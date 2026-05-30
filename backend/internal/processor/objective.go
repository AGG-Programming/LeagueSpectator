package processor

import (
	"slices"

	"github.com/AGG-Programming/LeagueSpectator/internal/league"
	"github.com/AGG-Programming/LeagueSpectator/pkg/models"
)

func (p *Processor) getTeamObjectives(events []league.Event, blueTeam []string, redTeam []string) {
	for _, event := range events {
		if event.EventID == 0 || event.EventName == "GameStart" {
			p.Reset()
		}

		if event.EventID <= p.LastProcessedEventID {
			continue
		}
		switch event.EventName {
		case "TurretKilled":
			{
				if slices.Contains(blueTeam, *event.KillerName) {
					p.BlueObjectives = append(p.BlueObjectives, models.Objective{
						Key:   event.EventName,
						Icon:  "",
						Kills: nil, //TODO: calculate total turret kills
					})
				} else if slices.Contains(redTeam, *event.KillerName) {
					p.RedObjectives = append(p.RedObjectives, models.Objective{
						Key:   event.EventName,
						Icon:  "",
						Kills: nil, //TODO: calculate total turret kills
					})
				}
			}
		case "DragonKill":
			{
				if slices.Contains(blueTeam, *event.KillerName) {
					p.BlueDrakeOrderKey++
					p.BlueObjectives = append(p.BlueObjectives, models.Objective{
						Key:      *event.DragonType,
						Icon:     "",
						OrderKey: &p.BlueDrakeOrderKey,
					})
				} else if slices.Contains(redTeam, *event.KillerName) {
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
				if slices.Contains(blueTeam, *event.KillerName) {
					p.BlueObjectives = append(p.BlueObjectives, models.Objective{
						Key:  event.EventName,
						Icon: "",
					})
				} else if slices.Contains(redTeam, *event.KillerName) {
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
				if slices.Contains(blueTeam, *event.KillerName) {
					p.BlueObjectives = append(p.BlueObjectives, models.Objective{
						Key:           event.EventName,
						Icon:          "",
						IsActive:      &isActive,
						RemainingTime: &remaining,
					})
				} else if slices.Contains(redTeam, *event.KillerName) {
					p.RedObjectives = append(p.RedObjectives, models.Objective{
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
		p.LastProcessedEventID = event.EventID
	}
}

func (p *Processor) Reset() {
	p.LastProcessedEventID = -1
	p.BlueObjectives = []models.Objective{}
	p.RedObjectives = []models.Objective{}
	p.BlueDrakeOrderKey = 0
	p.RedDrakeOrderKey = 0
}
