package api

import (
	"time"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/riot/lol"
	log "github.com/sirupsen/logrus"
)

type client struct {
	APIKey string
	Client *golio.Client
}

func CreateClient(apiKey string) client {
	return client{
		APIKey: apiKey,
		Client: golio.NewClient(apiKey,
			golio.WithRegion(api.RegionNorthAmerica),
			golio.WithLogger(log.New())),
	}
}

func (c client) SummonerByPUUID(puuid string) (*lol.Summoner, error) {
	summoner, err := c.Client.Riot.LoL.Summoner.GetByPUUID(puuid)
	if err != nil {
		return nil, err
	}

	return summoner, err
}

type getter struct {
	api      client
	summoner *lol.Summoner
	queues   []QueueType
}

func (c client) Get(summoner *lol.Summoner, queues []QueueType) getter {
	return getter{
		c,
		summoner,
		queues,
	}
}

func (g getter) options() []*lol.MatchListOptions {
	var options []*lol.MatchListOptions

	for _, queue := range g.queues {
		q := int(queue)

		options = append(options, &lol.MatchListOptions{
			Queue: &q,
		})
	}

	return options
}

func (g getter) Recent(name string) ([]*lol.Match, error) {
	matchIds, err := g.api.Client.Riot.LoL.Match.List(g.summoner.PUUID, 0, 20)
	if err != nil {
		return []*lol.Match{}, err
	}

	var matches []*lol.Match

	for _, matchId := range matchIds {
		match, err := g.api.Client.Riot.LoL.Match.Get(matchId)
		if err != nil {
			return []*lol.Match{}, err
		}

		matches = append(matches, match)
	}

	return matches, nil
}

func (g getter) Until(summoner *lol.Summoner, predicate func(*lol.Match) bool) ([]*lol.Match, error) {
	var matches []*lol.Match

	for _, options := range g.options() {
		for result := range g.api.Client.Riot.LoL.Match.ListStream(summoner.PUUID, options) {
			if result.Error != nil {
				break
			}

			match, err := g.api.Client.Riot.LoL.Match.Get(result.MatchID)
			if err != nil {
				break
			}

			if !predicate(match) {
				break
			}

			matches = append(matches, match)
		}
	}

	return matches, nil
}

func (g getter) Since(startTime time.Time) ([]*lol.Match, error) {
	return g.Between(startTime, time.Now())
}

func (g getter) Between(startTime time.Time, endTime time.Time) ([]*lol.Match, error) {
	var matches []*lol.Match

	for _, options := range g.options() {
		options.StartTime = startTime
		options.EndTime = endTime

		for result := range g.api.Client.Riot.LoL.Match.ListStream(g.summoner.PUUID, options) {
			if result.Error != nil {
				break
			}

			match, err := g.api.Client.Riot.LoL.Match.Get(result.MatchID)
			if err != nil {
				break
			}

			matches = append(matches, match)
		}
	}

	return matches, nil
}
