package ddragon

type Cache struct {
	champions map[string]string
	runes     map[int]string
	items     map[int]string
	spells    map[string]string
	ults      map[string]string
}

func NewCache(client *Client) (*Cache, error) {
	var err error
	version, err := client.GetLatestPatchVersion()
	if err != nil {
		return nil, err
	}
	cache := &Cache{}

	cache.champions, err = client.GetChampions(version)
	if err != nil {
		return nil, err
	}

	cache.items, err = client.GetItems(version)
	if err != nil {
		return nil, err
	}

	cache.runes, err = client.GetRunes(version)
	if err != nil {
		return nil, err
	}

	cache.spells, err = client.GetSpells(version)
	if err != nil {
		return nil, err
	}
	champs := make([]string, 0, len(cache.champions))
	for k := range cache.champions {
		champs = append(champs, k)
	}
	cache.ults, err = client.GetUlts(version, champs)

	return cache, nil
}

func (c *Cache) GetChampion(id string) string {
	return c.champions[id]
}
func (c *Cache) GetRune(id int) string {
	return c.runes[id]
}
func (c *Cache) GetItem(id int) string {
	return c.items[id]
}
func (c *Cache) GetSpell(id string) string {
	return c.spells[id]
}
func (c *Cache) GetUlt(id string) string {
	return c.ults[id]
}
