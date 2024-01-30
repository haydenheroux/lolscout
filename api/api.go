package api

import (
	"time"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/riot/lol"
	log "github.com/sirupsen/logrus"
)

type API struct {
	client *golio.Client
}

func New(key string) API {
	return API{
		client: golio.NewClient(key,
			golio.WithRegion(api.RegionNorthAmerica),
			golio.WithLogger(log.New())),
	}
}

type getter struct {
	api      API
	summoner *lol.Summoner
	// TODO Make array of enum
	queues []int
}

func (api API) Get(summoner *lol.Summoner, queues []int) getter {
	return getter{
		api,
		summoner,
		queues,
	}
}

func (api API) Summoner(name string) (*lol.Summoner, error) {
	summoner, err := api.client.Riot.LoL.Summoner.GetByName(name)
	if err != nil {
		return nil, err
	}

	return summoner, err
}

func (g getter) Recent(name string) ([]*lol.Match, error) {
	matchIds, err := g.api.client.Riot.LoL.Match.List(g.summoner.PUUID, 0, 20)
	if err != nil {
		return []*lol.Match{}, err
	}

	var matches []*lol.Match

	for _, matchId := range matchIds {
		match, err := g.api.client.Riot.LoL.Match.Get(matchId)
		if err != nil {
			return []*lol.Match{}, err
		}

		matches = append(matches, match)
	}

	return matches, nil
}

func (g getter) Until(summoner *lol.Summoner, predicate func(*lol.Match) bool) ([]*lol.Match, error) {
	var matches []*lol.Match

	for result := range g.api.client.Riot.LoL.Match.ListStream(summoner.PUUID) {
		if result.Error != nil {
			break
		}

		match, err := g.api.client.Riot.LoL.Match.Get(result.MatchID)
		if err != nil {
			break
		}

		if !predicate(match) {
			break
		}

		matches = append(matches, match)
	}

	return matches, nil
}

func (g getter) From(startTime time.Time) ([]*lol.Match, error) {
	return g.Between(startTime, time.Now())
}

func (g getter) Between(startTime time.Time, endTime time.Time) ([]*lol.Match, error) {
	var matches []*lol.Match

	options := &lol.MatchListOptions{
		StartTime: startTime,
		EndTime:   endTime,
	}

	for _, queue := range g.queues {
		options.Queue = &queue

		for result := range g.api.client.Riot.LoL.Match.ListStream(g.summoner.PUUID, options) {
			if result.Error != nil {
				break
			}

			match, err := g.api.client.Riot.LoL.Match.Get(result.MatchID)
			if err != nil {
				break
			}

			matches = append(matches, match)
		}
	}

	return matches, nil
}
