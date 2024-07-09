package henrik

import (
	"fmt"
	"github.com/go-zoox/fetch"
	"github.com/spf13/viper"
)

func GetMatches(region, name, tag string) ([]Match, error) {
	response, err := fetch.Get(fmt.Sprintf("%s/v3/matches/%s/%s/%s", APIURL, region, name, tag), &fetch.Config{
		Headers: map[string]string{
			"Authorization": viper.GetString("henrik.api_key"),
		},
	})
	if err != nil {
		return nil, err
	}

	var rank Response[[]Match]
	if err := response.UnmarshalJSON(&rank); err != nil {
		return nil, err
	}

	return rank.Data, nil
}
