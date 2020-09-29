package services

import (
	"github.com/Tomer-Alony/IPLocator/src/models"
	"github.com/Tomer-Alony/IPLocator/src/store"
)

type IPService struct {
	context store.DataStoreService
}

func NewIPServiceRepo(store store.DataStoreService) *IPService {
	// Create and return a new instance of the IPService Repository
	return &IPService{context: store}
}

func (repo *IPService) FindCountry(ip string) (*models.IP, error) {
	ipDetails, err := repo.context.GetIPDetails(ip)

	if err != nil {
		return nil, err
	}

	return &ipDetails, nil
}