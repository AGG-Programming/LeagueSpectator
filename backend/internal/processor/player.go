package processor

import (
	"strings"

	"github.com/AGG-Programming/LeagueSpectator/internal/league"
	"github.com/AGG-Programming/LeagueSpectator/pkg/models"
)

func (p *Processor) getPlayers(players []league.Player) ([]models.Player, []models.Player) {
	var bluePlayers, redPlayers []models.Player
	for _, player := range players {
		pl := models.Player{
			ChampionName: player.ChampionName,
			Icon:         p.cache.GetChampion(player.ChampionName),
			IsDead:       player.IsDead,
			Level:        player.Level,
			Position:     player.Position,
			RespawnTimer: player.RespawnTimer,
			RiotId:       player.RiotID,
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
			Items: p.getItems(player.Items),
			Scores: models.Scores{
				Assists:    player.Scores.Assists,
				CreepScore: player.Scores.CreepScore,
				Deaths:     player.Scores.Deaths,
				Kills:      player.Scores.Kills,
				WardScore:  player.Scores.WardScore,
			},
			Spells:  p.getSpells(player.Spells),
			UltIcon: p.cache.GetUlt(player.ChampionName),
		}

		if player.Team == BlueTeam {
			bluePlayers = append(bluePlayers, pl)
		} else {
			redPlayers = append(redPlayers, pl)
		}
	}
	return bluePlayers, redPlayers
}

func (p *Processor) getItems(items []league.Item) []models.Item {
	var playerItems []models.Item
	for _, item := range items {
		pItem := models.Item{
			Id:         item.ItemID,
			Icon:       p.cache.GetItem(item.ItemID),
			Slot:       item.Slot,
			Consumable: item.Consumable,
			Count:      item.Count,
		}
		playerItems = append(playerItems, pItem)
	}
	return playerItems
}

func (p *Processor) getSpells(spells league.Spells) []models.Spell {
	spellOneId := p.getSummonerSpellName(spells.SpellOne.RawDisplayName)
	spellTwoId := p.getSummonerSpellName(spells.SpellTwo.RawDisplayName)

	playerSpells := []models.Spell{
		{
			DisplayName: spells.SpellOne.DisplayName,
			Icon:        p.cache.GetSpell(spellOneId),
		},
		{
			DisplayName: spells.SpellTwo.DisplayName,
			Icon:        p.cache.GetSpell(spellTwoId),
		},
	}

	return playerSpells
}

func (p *Processor) getSummonerSpellName(input string) string {
	parts := strings.Split(input, "_")
	if len(parts) > 2 {
		return parts[2]
	}

	return ""
}
