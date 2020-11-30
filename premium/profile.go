package premium

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func fetchPremiumLink(username string) string {
	return fmt.Sprint("https://api.mojang.com/users/profiles/minecraft/", username)
}

type jsonResponse struct {
	id   string
	name string
	demo bool
}

func FetchProfile(username string) (bool, error) {
	response, err := http.Get(fetchPremiumLink(username))
	if err != nil {
		return false, err
	}

	if response.StatusCode == http.StatusNoContent {
		return false, nil
	}

	if response.StatusCode == http.StatusOK {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return false, err
		}

		var response jsonResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return false, err
		}

		if response.demo {
			return false, nil
		}

		return true, nil
	}

	return false, errors.New("unexpected status code")
}
