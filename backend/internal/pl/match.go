package pl

import "sort"

func (c *Client) NextMatch(matches []MatchResponse) *MatchResponse {
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].MatchTime == matches[j].MatchTime {
			return matches[i].MatchPlayday < matches[j].MatchPlayday
		}
		return matches[i].MatchTime < matches[j].MatchTime
	})
	if len(matches) == 0 {
		return nil
	}
	return &matches[0]
}
