package metrics

import "github.com/KnutZuidema/golio/riot/lol"

func FilterBySummoner(matches []*lol.Match, summoner *lol.Summoner) []*lol.Match {
	return FilterBySummoners(matches, []*lol.Summoner{summoner})
}

func FilterBySummoners(matches []*lol.Match, summoners []*lol.Summoner) []*lol.Match {
	var result []*lol.Match

	for _, match := range matches {
		if hasSummoners(match, summoners) {
			result = append(result, match)
		}
	}

	return result
}

func hasSummoners(match *lol.Match, summoners []*lol.Summoner) bool {
	for _, summoner := range summoners {
		participated := false

		for _, participantPUUID := range match.Metadata.Participants {
			if summoner.PUUID == participantPUUID {
				participated = true
			}
		}

		if !participated {
			return false
		}
	}

	return true
}
