package models

type DynamicGameData struct {
	BlueTeam Team    `json:"blueTeam"`
	RedTeam  Team    `json:"redTeam"`
	Timers   []Timer `json:"timers"`
	GameTime float64 `json:"gameTime"`
}

type Timer struct {
	Type      string `json:"type"`
	SpawnTime int    `json:"SpawnTime"`
	Alive     bool   `json:"alive"`
}

type Team struct {
	Score      int         `json:"score"`
	Objectives []Objective `json:"objectives"`
	Players    []Player    `json:"players"`
}

type Objective struct {
	Key           string   `json:"key"`
	Icon          string   `json:"icon"`
	Kills         *int     `json:"kills,omitempty"`
	IsActive      *bool    `json:"isActive,omitempty"`
	RemainingTime *float64 `json:"remainingTime,omitempty"`
	OrderKey      *int     `json:"orderKey,omitempty"`
}

type Player struct {
	ChampionName string  `json:"championName"`
	Icon         string  `json:"icon"`
	IsDead       bool    `json:"isDead"`
	Level        int     `json:"level"`
	Position     string  `json:"position"`
	RespawnTimer float64 `json:"respawnTimer"`
	RiotId       string  `json:"riotId"`
	Runes        Runes   `json:"runes"`
	Items        []Item  `json:"items"`
	Scores       Scores  `json:"scores"`
	Spells       []Spell `json:"spells"`
	UltIcon      string  `json:"ultIcon"`
	QuestIcon    string  `json:"questIcon"`
}

type Spell struct {
	DisplayName string `json:"displayName"`
	Icon        string `json:"icon"`
}

type Scores struct {
	Assists    int     `json:"assists"`
	CreepScore int     `json:"creepScore"`
	Deaths     int     `json:"deaths"`
	Kills      int     `json:"kills"`
	WardScore  float64 `json:"wardScore"`
}

type Item struct {
	Id         int    `json:"id"`
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
}

type PrimeLeague struct {
	GroupTitle   string        `json:"groupTitle"`
	TargetTeam   PrimeTeam     `json:"targetTeam"`
	LeadingTeam  *PrimeTeam    `json:"leadingTeam,omitempty"`
	TrailingTeam *PrimeTeam    `json:"trailingTeam,omitempty"`
	LastTeam     *PrimeTeam    `json:"lastTeam,omitempty"`
	NextMatch    NextMatch     `json:"nextMatch"`
	CurrentMatch *CurrentMatch `json:"currentMatch,omitempty"`
}

type PrimeTeam struct {
	Tag        string `json:"tag"`
	Wins       int    `json:"wins"`
	Losses     int    `json:"losses"`
	Points     *int   `json:"points,omitempty"`
	Position   *int   `json:"position,omitempty"`
	Img        string `json:"img"`
	MatchScore *int   `json:"matchScore,omitempty"`
}

type NextMatch struct {
	Tag       string `json:"tag"`
	Img       string `json:"img"`
	Status    string `json:"status"`
	MatchTime int    `json:"matchTime"`
}

type CurrentMatch struct {
	Opponent1 PrimeTeam `json:"opponent1"`
	Opponent2 PrimeTeam `json:"opponent2"`
}
