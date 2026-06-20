package pl

import (
	"fmt"
	"sort"
)

func (c *Client) NextMatch(matches []MatchResponse) (MatchResponse, error) {
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].MatchTime == matches[j].MatchTime {
			return matches[i].MatchPlayday < matches[j].MatchPlayday
		}
		return matches[i].MatchTime < matches[j].MatchTime
	})
	if len(matches) == 0 {
		return MatchResponse{}, fmt.Errorf("no matches available")
	}
	return matches[0], nil
}
