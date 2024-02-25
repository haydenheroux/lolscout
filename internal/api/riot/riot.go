package riot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type client struct {
	APIKey string
}

func CreateClient(apiKey string) client {
	return client{
		APIKey: apiKey,
	}
}

type Account struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type getter struct {
	client   client
	gameName string
	tagLine  string
}

func (c client) Get(gameName, tagLine string) getter {
	return getter{
		client:   c,
		gameName: gameName,
		tagLine:  tagLine,
	}
}

func (g getter) Account() (*Account, error) {
	url := fmt.Sprintf("https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s", g.gameName, g.tagLine)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Riot-Token", g.client.APIKey)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("failed to get account by riot id")
	}

	contents, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result Account

	if err := json.Unmarshal(contents, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
