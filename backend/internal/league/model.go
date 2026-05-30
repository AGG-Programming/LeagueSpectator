package league

type GameData struct {
	Players []Player `json:"allPlayers"`
	Events  struct {
		Events []Event `json:"Events"`
	} `json:"events"`
	GameTime float64 `json:"gameTime"`
}

type Event struct {
	EventID    int     `json:"EventID"`
	EventName  string  `json:"EventName"`
	EventTime  float64 `json:"EventTime"`
	KillerName *string `json:"KillerName,omitempty"`
	DragonType *string `json:"DragonType,omitempty"`
}

type Player struct {
	ChampionName string  `json:"championName"`
	IsDead       bool    `json:"isDead"`
	Items        []Item  `json:"items"`
	Level        int     `json:"level"`
	Position     string  `json:"position"`
	RespawnTimer float64 `json:"respawnTimer"`
	Runes        Runes   `json:"runes"`
	Scores       Scores  `json:"scores"`
	RiotID       string  `json:"riotId"`
	Spells       Spells  `json:"summonerSpells"`
	Team         string  `json:"team"`
}

type Item struct {
	CanUse      bool   `json:"canUse"`
	Consumable  bool   `json:"consumable"`
	Count       int    `json:"count"`
	DisplayName string `json:"displayName"`
	ItemID      int    `json:"itemID"`
	Price       int    `json:"price"`
	Slot        int    `json:"slot"`
}

type Runes struct {
	Keystone  Rune `json:"keystone"`
	Primary   Rune `json:"primaryRuneTree"`
	Secondary Rune `json:"secondaryRuneTree"`
}

type Rune struct {
	DisplayName string `json:"displayName"`
	ID          int    `json:"id"`
}

type Scores struct {
	Assists    int     `json:"assists"`
	CreepScore int     `json:"creepScore"`
	Deaths     int     `json:"deaths"`
	Kills      int     `json:"kills"`
	WardScore  float64 `json:"wardScore"`
}

type Spells struct {
	SpellOne Spell `json:"summonerSpellOne"`
	SpellTwo Spell `json:"summonerSpellTwo"`
}

type Spell struct {
	DisplayName    string `json:"displayName"`
	RawDisplayName string `json:"rawDisplayName"`
}
