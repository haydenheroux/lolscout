package api

import (
	"encoding/json"
	"fmt"
	"io"
	"lolscout/data"
	"net/http"
	"time"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/riot/lol"
	log "github.com/sirupsen/logrus"
)

type API struct {
	APIKey string
	Client *golio.Client
}

func New(apiKey string) API {
	return API{
		APIKey: apiKey,
		Client: golio.NewClient(apiKey,
			golio.WithRegion(api.RegionNorthAmerica),
			golio.WithLogger(log.New())),
	}
}

type summonerByTagResult struct {
	PUUID string `json:"puuid"`
	Name  string `json:"gameName"`
	Tag   string `json:"tagLine"`
}

// TODO change when golio is updated with PR #60
func (api API) TODO_SummonerByTag_TODO(name, tag string) (*lol.Summoner, error) {
	url := fmt.Sprintf("https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s", name, tag)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Riot-Token", api.APIKey)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	contents, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result summonerByTagResult

	if err := json.Unmarshal(contents, &result); err != nil {
		return nil, err
	}

	return api.SummonerByPUUID(result.PUUID)
}

func (api API) SummonerByPUUID(puuid string) (*lol.Summoner, error) {
	summoner, err := api.Client.Riot.LoL.Summoner.GetByPUUID(puuid)
	if err != nil {
		return nil, err
	}

	return summoner, err
}

type getter struct {
	api      API
	summoner *lol.Summoner
	queues   []data.Queue
}

func (api API) Get(summoner *lol.Summoner, queues []data.Queue) getter {
	return getter{
		api,
		summoner,
		queues,
	}
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

	for result := range g.api.Client.Riot.LoL.Match.ListStream(summoner.PUUID) {
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
		// TODO Hack
		queue_ := int(queue)
		options.Queue = &queue_

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
