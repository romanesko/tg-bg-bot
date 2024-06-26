package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type CitiesResponse struct {
	Status int    `json:"status"`
	Descr  string `json:"descr"`
	Data   struct {
		LastUpdate string
		Countries  map[string]struct {
			CountryID   string `json:"country_id"`
			NameRussian string `json:"country_name_rus"`
			ISO         string `json:"country_iso"`
			Cities      map[string]struct {
				ID            string `json:"id"`
				DisplayName   string `json:"city_display_name"`
				NameRussian   string `json:"city_name_rus"`
				RegionRussian string `json:"city_region_rus"`
			} `json:"сities"`
		} `json:"country"`
	} `json:"data"`
}

type CityRow struct {
	CountryID   int
	CountryName string
	CityID      int
	CityName    string
	DisplayName string
}

type Api struct {
	host string
	key  string
}

func fetchCities() (CitiesResponse, error) {
	resp, err := http.Get(fmt.Sprintf("https://bodygraph.online/api_v1/city_list.php?dkey=test_public_key"))
	if err != nil {
		return CitiesResponse{}, err
	}
	defer resp.Body.Close()
	var result CitiesResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

func GetCities() ([]CityRow, error) {

	cityRows := []CityRow{}

	cities, err := fetchCities()
	if err != nil {
		return cityRows, fmt.Errorf("fetchCities: %w", err)
	}
	for countryKey := range cities.Data.Countries {
		country := cities.Data.Countries[countryKey]
		for cityKey := range country.Cities {
			city := country.Cities[cityKey]

			//println(city.ID, city.DisplayName, city.NameRussian, city.RegionRussian)

			countryID, err := strconv.Atoi(country.CountryID)
			if err != nil {
				return cityRows, fmt.Errorf("FindCity: %w", err)
			}
			cityID, err := strconv.Atoi(city.ID)
			if err != nil {
				return cityRows, fmt.Errorf("FindCity: %w", err)
			}

			cityRows = append(cityRows, CityRow{
				CountryID:   countryID,
				CountryName: country.NameRussian,
				CityID:      cityID,
				CityName:    city.NameRussian,
				DisplayName: strings.Replace(city.DisplayName, " ()", "", -1),
			})

		}

	}
	return cityRows, nil
}
