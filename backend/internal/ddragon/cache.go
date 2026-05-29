package ddragon

type Cache struct {
	champions map[string]string
	runes     map[int]string
	items     map[string]string
	spells    map[string]string
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

	return cache, nil
}

func (c *Cache) GetChampion(id string) string {
	return c.champions[id]
}
func (c *Cache) GetRune(id int) string {
	return c.runes[id]
}
func (c *Cache) GetItem(id string) string {
	return c.items[id]
}
func (c *Cache) GetSpell(id string) string {
	return c.spells[id]
}
