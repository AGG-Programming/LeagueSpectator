package models

import (
	"encoding/json"
	"strconv"
)

type DynamicGameData struct {
	BlueTeam Team    `json:"blueTeam"`
	RedTeam  Team    `json:"redTeam"`
	Timers   []Timer `json:"timers"`
	GameTime float64 `json:"gameTime"`
}

type StaticGameData struct {
	Players []Player `json:"players"`
}

type Timer struct {
	Type      string `json:"type"`
	SpawnTime int    `json:"SpawnTime"`
	Alive     bool   `json:"alive"`
}

type Team struct {
	Score      int         `json:"score"`
	Gold       int         `json:"gold"`
	Objectives []Objective `json:"objectives"`
	Players    []Player    `json:"players"`
}

type Objective struct {
	Key           string `json:"key"`
	Icon          string `json:"icon"`
	Kills         *int   `json:"kills,omitempty"`
	IsActive      *bool  `json:"isActive,omitempty"`
	RemainingTime *int   `json:"remainingTime,omitempty"`
	OrderKey      *int   `json:"orderKey,omitempty"`
}

type Player struct {
	ChampionName    string  `json:"championName"`
	Icon            string  `json:"icon"`
	IsDead          bool    `json:"isDead"`
	Level           int     `json:"level"`
	Position        string  `json:"position"`
	RespawnTimer    float64 `json:"respawnTimer"`
	RiotId          string  `json:"riotId"`
	PlayerTotalGold int     `json:"playerTotalGold"`
	Runes           Runes   `json:"runes"`
	Items           []Item  `json:"items"`
	Scores          Scores  `json:"scores"`
	Spells          []Spell `json:"spells"`
}

type Spell struct {
	DisplayName string `json:"displayName"`
	Icon        string `json:"icon"`
	Extended    string
}

type Scores struct {
	Assists    int     `json:"assists"`
	CreepScore int     `json:"creepScore"`
	Deaths     int     `json:"deaths"`
	Kills      int     `json:"kills"`
	WardScore  float64 `json:"wardScore"`
}

type Item struct {
	Id         string `json:"id"`
	Icon       string `json:"icon"`
	Slot       int    `json:"slot"`
	Consumable bool   `json:"consumable"`
	Count      int    `json:"count"`
}

type Runes struct {
	Keystone  Rune `json:"keystone"`
	Primary   Rune `json:"primary"`
	Secondary Rune `json:"secondary"`
}

type Rune struct {
	DisplayName string `json:"displayName"`
	Icon        string `json:"icon"`
	Extended    string
}

func (d *DynamicGameData) UnmarshalJSON(data []byte) error {
	var response struct {
		AllPlayers []struct {
			ChampionName string  `json:"championName"`
			IsDead       bool    `json:"isDead"`
			Level        int     `json:"level"`
			Position     string  `json:"position"`
			RespawnTimer float64 `json:"respawnTimer"`
			SummonerName string  `json:"summonerName"`
			Team         string  `json:"team"`
			Scores       struct {
				Assists    int     `json:"assists"`
				CreepScore int     `json:"creepScore"`
				Deaths     int     `json:"deaths"`
				Kills      int     `json:"kills"`
				WardScore  float64 `json:"wardScore"`
			} `json:"scores"`
			SummonerSpells struct {
				SummonerSpellOne struct {
					DisplayName string `json:"displayName"`
					Extended    string `json:"rawDisplayName"`
				} `json:"SummonerSpellOne"`
				SummonerSpellTwo struct {
					DisplayName string `json:"displayName"`
					Extended    string `json:"rawDisplayName"`
				} `json:"SummonerSpellTwo"`
			} `json:"summonerSpells"`
			Items []struct {
				DisplayName string `json:"displayName"`
				ItemID      int    `json:"itemID"`
				Slot        int    `json:"slot"`
				Consumable  bool   `json:"consumable"`
				Count       int    `json:"count"`
			} `json:"items"`
			Runes struct {
				Keystone struct {
					DisplayName string `json:"displayName"`
				} `json:"keystone"`
				PrimaryRuneTree struct {
					DisplayName string `json:"displayName"`
					Extended    string `json:"rawDisplayName"`
				} `json:"primaryRuneTree"`
				SecondaryRuneTree struct {
					DisplayName string `json:"displayName"`
					Extended    string `json:"rawDisplayName"`
				} `json:"secondaryRuneTree"`
			} `json:"runes"`
		} `json:"allPlayers"`
		Events struct {
			Events []struct {
				EventName  string  `json:"EventName"`
				EventTime  float64 `json:"EventTime"`
				KillerName string  `json:"KillerName"`
				DragonType string  `json:"DragonType"`
				AcingTeam  string  `json:"AcingTeam"`
			} `json:"Events"`
		} `json:"events"`
		GameData struct {
			GameTime float64 `json:"gameTime"`
		} `json:"gameData"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	d.GameTime = response.GameData.GameTime
	d.BlueTeam = Team{Players: make([]Player, 0), Objectives: make([]Objective, 0)}
	d.RedTeam = Team{Players: make([]Player, 0), Objectives: make([]Objective, 0)}

	for _, rp := range response.AllPlayers {
		var myItems []Item
		playerGoldFromItems := 0

		for _, ri := range rp.Items {
			myItems = append(myItems, Item{
				Id:         ri.DisplayName,
				Icon:       strconv.Itoa(ri.ItemID),
				Slot:       ri.Slot,
				Consumable: ri.Consumable,
				Count:      ri.Count,
			})
		}

		p := Player{
			ChampionName:    rp.ChampionName,
			IsDead:          rp.IsDead,
			Level:           rp.Level,
			Position:        rp.Position,
			RespawnTimer:    rp.RespawnTimer,
			RiotId:          rp.SummonerName,
			PlayerTotalGold: playerGoldFromItems,
			Items:           myItems,
			Scores: Scores{
				Assists:    rp.Scores.Assists,
				CreepScore: rp.Scores.CreepScore,
				Deaths:     rp.Scores.Deaths,
				Kills:      rp.Scores.Kills,
				WardScore:  rp.Scores.WardScore,
			},
			Spells: []Spell{
				{
					DisplayName: rp.SummonerSpells.SummonerSpellOne.DisplayName,
					Extended:    rp.SummonerSpells.SummonerSpellOne.Extended,
				},
				{
					DisplayName: rp.SummonerSpells.SummonerSpellTwo.DisplayName,
					Extended:    rp.SummonerSpells.SummonerSpellTwo.Extended,
				},
			},
			Runes: Runes{
				Keystone: Rune{
					DisplayName: rp.Runes.Keystone.DisplayName,
				},
				Primary: Rune{
					DisplayName: rp.Runes.PrimaryRuneTree.DisplayName,
					Extended:    rp.Runes.PrimaryRuneTree.Extended,
				},
				Secondary: Rune{
					DisplayName: rp.Runes.SecondaryRuneTree.DisplayName,
					Extended:    rp.Runes.SecondaryRuneTree.Extended,
				},
			},
		}

		if rp.Team == "ORDER" {
			d.BlueTeam.Players = append(d.BlueTeam.Players, p)
		} else {
			d.RedTeam.Players = append(d.RedTeam.Players, p)
		}
	}

	d.calculateTeamStatsAndObjectives(response.Events.Events)

	return nil
}

func (d *DynamicGameData) calculateTeamStatsAndObjectives(events []struct {
	EventName  string  `json:"EventName"`
	EventTime  float64 `json:"EventTime"`
	KillerName string  `json:"KillerName"`
	DragonType string  `json:"DragonType"`
	AcingTeam  string  `json:"AcingTeam"`
}) {
	blueKills := 0
	redKills := 0

	for _, event := range events {
		if event.EventName == "ChampionKill" {
			if d.isBluePlayer(event.KillerName) {
				blueKills++
			} else {
				redKills++
			}
		}
	}

	d.BlueTeam.Score = blueKills
	d.RedTeam.Score = redKills

	d.Timers = []Timer{
		{Type: "baron", SpawnTime: 1200, Alive: d.GameTime >= 1200},
	}
}

func (d *DynamicGameData) isBluePlayer(name string) bool {
	for _, p := range d.BlueTeam.Players {
		if p.RiotId == name {
			return true
		}
	}
	return false
}
