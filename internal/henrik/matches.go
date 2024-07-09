package henrik

import (
	"fmt"
	"github.com/spf13/viper"
)

func GetMatches(region, name, tag string) ([]Match, error) {
	url := fmt.Sprintf("%s/v3/matches/%s/%s/%s", APIURL, region, name, tag)
	headers := map[string]string{
		"Authorization": viper.GetString("henrik.api_key"),
	}

	response, err := fetchWithRetry(url, headers)
	if err != nil {
		return nil, err
	}

	var matches Response[[]Match]
	if err := response.UnmarshalJSON(&matches); err != nil {
		return nil, err
	}

	return matches.Data, nil
}
