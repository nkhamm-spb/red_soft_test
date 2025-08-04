package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func GetJson(ctx context.Context, url string) (*map[string]interface{}, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Error in request %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Wrong http status code: %d in url: %s", resp.StatusCode, url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func GetGender(ctx context.Context, name string, surname string) (string, error) {
	url := fmt.Sprintf("https://api.genderize.io?name=%s", url.QueryEscape(fmt.Sprintf("%s %s", name, surname)))
	jsonMap, err := GetJson(ctx, url)

	if err != nil {
		return "", err
	}

	if gender, ok := (*jsonMap)["gender"].(string); ok {
		log.Printf("Getted gender: %s for user: %s %s\n", gender, name, surname)
		return gender, nil
	} else {
		return "", fmt.Errorf("Gender wrong result format")
	}
}

func GetAge(ctx context.Context, name string, surname string) (int, error) {
	url := fmt.Sprintf("https://api.agify.io?name=%s", url.QueryEscape(fmt.Sprintf("%s %s", name, surname)))
	jsonMap, err := GetJson(ctx, url)

	if err != nil {
		return 0, err
	}

	if age, ok := (*jsonMap)["age"].(float64); ok {
		log.Printf("Getted age: %d for user: %s %s\n", int(age), name, surname)
		return int(age), nil
	} else {
		return 0, fmt.Errorf("Age wrong result format")
	}
}

func GetNationalize(ctx context.Context, name string, surname string) (string, error) {
	url := fmt.Sprintf("https://api.nationalize.io?name=%s", url.QueryEscape(fmt.Sprintf("%s %s", name, surname)))
	jsonMap, err := GetJson(ctx, url)

	if err != nil {
		return "", err
	}

	countryList, ok := (*jsonMap)["country"].([]interface{})
	if !ok {
		return "", fmt.Errorf("Nationalize wrong result format")
	}

	itemMap, ok := countryList[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Nationalize wrong result format")
	}

	if nationalize, ok := itemMap["country_id"].(string); ok {
		log.Printf("Getted nationalize: %s for user: %s %s\n", nationalize, name, surname)
		return nationalize, nil
	} else {
		return "", fmt.Errorf("Nationalize wrong result format")
	}
}
