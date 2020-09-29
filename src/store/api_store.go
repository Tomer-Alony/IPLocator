package store

import (
	"encoding/json"
	"fmt"
	"github.com/Tomer-Alony/IPLocator/src/errors"
	"github.com/Tomer-Alony/IPLocator/src/models"
	"io/ioutil"
	"net/http"
)

type IPReader struct {
	Country string `json:"country_name"`
	City string `json:"city"`
}

func NewAPIDataStore() DataStoreService {
	return &DataStore{}
}

func (store DataStore) Access(path, key string) DataStoreService {
	return &DataStore{
		Path: path,
		Key:  key,
	}
}

func (store *DataStore) GetIPDetails(ip string) (models.IP, error) {
	requestURL := fmt.Sprintf(store.Path, ip, store.Key)
	response, err := http.Get(requestURL)

	if err != nil {
		return models.IP{}, errors.HttpError{Message: err.Error(), Code: 500}
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return models.IP{}, errors.HttpError{Message: err.Error(), Code: 500}
	}

	var ipObject IPReader
	json.Unmarshal([]byte(body), &ipObject)

	if ipObject.City == "" || ipObject.Country == "" {
		return models.IP{}, errors.HttpError{Message: "IP not found", Code: 500}
	}

	return models.IP{Country: ipObject.Country, City: ipObject.City}, nil
}