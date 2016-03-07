package agent

import (
	"github.com/resourced/resourced/storage"
)

// setStorages configures all kind of storages, i.e. metadata storage.
func (a *Agent) setStorages() error {
	a.Db = storage.NewStorage()
	return nil
}
