package models

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDynamicGameData_UnmarshalJSON(t *testing.T) {
	riotJson, err := os.ReadFile("../../testdata/allgamedata.json")
	if err != nil {
		t.Fatal(err)
	}

	expected := DynamicGameData{
		BlueTeam: Team{
			Score:      0,
			Gold:       0,
			Objectives: []Objective{},
			Players: []Player{
				{
					ChampionName:    "Annie",
					Icon:            "",
					IsDead:          false,
					Level:           1,
					Position:        "TOP",
					RespawnTimer:    0.0,
					RiotId:          "Riot Tuxedo",
					PlayerTotalGold: 0,
					Runes: Runes{
						Keystone: Rune{
							DisplayName: "Electrocute",
							Icon:        "",
						},
						Secondary: Rune{
							DisplayName: "Sorcery",
							Icon:        "",
						},
					},
					Scores: Scores{
						Assists:    0,
						CreepScore: 0,
						Deaths:     0,
						Kills:      0,
						WardScore:  0.0,
					},
					Spells: []Spell{
						{
							DisplayName: "Flash",
							Icon:        "",
						},
						{
							DisplayName: "Ignite",
							Icon:        "",
						},
					},
				},
			},
		},
		RedTeam: Team{
			Objectives: []Objective{},
			Players:    []Player{},
		},
		Timers: []Timer{
			{
				Type:      "baron",
				SpawnTime: 1200,
				Alive:     false,
			},
		},
		GameTime: 0.000000000,
	}

	var d DynamicGameData
	err = d.UnmarshalJSON(riotJson)

	assert.NoError(t, err)
	assert.Equal(t, expected, d)
}
