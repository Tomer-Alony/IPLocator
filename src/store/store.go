package store

import "github.com/Tomer-Alony/IPLocator/src/models"

type DataStore struct {
	Path string
	Key string
}

type DataStoreService interface {
	Access(path, key string) DataStoreService
	GetIPDetails(ip string) (models.IP, error)
}
