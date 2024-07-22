package service

import "time"

const (
	interval = 20_000
)

type clientConnectionsStorageService interface {
	RemoveOutdatedConnections() error
}

type periodicClientConnectionsCleanUp struct {
	storage  clientConnectionsStorageService
	interval int
}

func newClientConnectionsCleanUpService(storage clientConnectionsStorageService) *periodicClientConnectionsCleanUp {
	cleanUpService := &periodicClientConnectionsCleanUp{storage, interval}
	go func() {
		cleanUpService.storage.RemoveOutdatedConnections()
		time.Sleep(time.Millisecond * time.Duration(cleanUpService.interval))
	}()
	return cleanUpService
}
