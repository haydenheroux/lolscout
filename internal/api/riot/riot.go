package riot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
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
		return &Account{}, err
	}

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusForbidden {
			return &Account{}, errors.New("possibly bad riot api key")
		}

		if res.StatusCode == http.StatusTooManyRequests {
			retry := res.Header.Get("Retry-After")
			seconds, err := strconv.Atoi(retry)

			if err != nil {
				logrus.Debug(err)
				return nil, err
			}

			logrus.Infof("rate limited, waiting %d seconds", seconds)

			time.Sleep(time.Duration(seconds) * time.Second)

			return g.Account()
		}

		if res.StatusCode == http.StatusNotFound {
			return &Account{}, errors.New("user not found")
		}

		logrus.Infof("status code %d", res.StatusCode)

		return &Account{}, errors.New("unknown error")
	}

	contents, err := io.ReadAll(res.Body)
	if err != nil {
		return &Account{}, err
	}

	var result Account

	if err := json.Unmarshal(contents, &result); err != nil {
		return &Account{}, err
	}

	return &result, nil
}
