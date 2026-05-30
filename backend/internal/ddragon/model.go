package ddragon

type ChampionResponse struct {
	Data map[string]struct {
		Name  string `json:"name"`
		Image struct {
			Full string `json:"full"`
		} `json:"image"`
	} `json:"data"`
}

type ItemResponse struct {
	Data map[string]Item `json:"data"`
}

type Item struct {
	Image ImageInfo `json:"image"`
}

type ImageInfo struct {
	Full string `json:"full"`
}

type SpellResponse struct {
	Data map[string]struct {
		ID    string `json:"id"`
		Image struct {
			Full string `json:"full"`
		} `json:"image"`
	} `json:"data"`
}
type RuneResponse []RuneTree

type RuneTree struct {
	ID    int        `json:"id"`
	Icon  string     `json:"icon"`
	Slots []RuneSlot `json:"slots"`
}

type RuneSlot struct {
	Runes []Rune `json:"runes"`
}

type Rune struct {
	ID   int    `json:"id"`
	Icon string `json:"icon"`
}
