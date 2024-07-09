package henrik

import (
	"fmt"
	"github.com/spf13/viper"
)

func GetRank(region, name, tag string) (Rank, error) {
	url := fmt.Sprintf("%s/v2/mmr/%s/%s/%s", APIURL, region, name, tag)
	headers := map[string]string{
		"Authorization": viper.GetString("henrik.api_key"),
	}

	response, err := fetchWithRetry(url, headers)
	if err != nil {
		return Rank{}, err
	}

	var rank Response[Rank]
	if err := response.UnmarshalJSON(&rank); err != nil {
		return Rank{}, err
	}

	return rank.Data, nil
}
