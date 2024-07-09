package henrik

import (
	"fmt"
	"github.com/go-zoox/fetch"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GetRank(region, name, tag string) (Rank, error) {
	response, err := fetch.Get(fmt.Sprintf("%s/v2/mmr/%s/%s/%s", APIURL, region, name, tag), &fetch.Config{
		Headers: map[string]string{
			"Authorization": viper.GetString("henrik.api_key"),
		},
	})
	if err != nil {
		return Rank{}, err
	}

	log.Infof("Response: %s", response.Body)

	var rank Response[Rank]
	if err := response.UnmarshalJSON(&rank); err != nil {
		return Rank{}, err
	}

	return rank.Data, nil
}
