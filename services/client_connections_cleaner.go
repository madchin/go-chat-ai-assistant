package service

import "time"

const (
	interval = 20_000
)

type connectionsStorageRemover interface {
	RemoveOutdatedConnections()
}

type periodicClientConnectionsCleanUp struct {
	storage  connectionsStorageRemover
	interval int
}

func newClientConnectionsCleanUpService(storage connectionsStorageRemover) *periodicClientConnectionsCleanUp {
	cleanUpService := &periodicClientConnectionsCleanUp{storage, interval}
	go func() {
		cleanUpService.storage.RemoveOutdatedConnections()
		time.Sleep(time.Millisecond * time.Duration(cleanUpService.interval))
	}()
	return cleanUpService
}
