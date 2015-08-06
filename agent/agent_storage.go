package agent

import (
	"github.com/resourced/resourced/storage"
)

// setStorages configures all kind of storages, i.e. metadata storage.
func (a *Agent) setStorages() error {
	a.MetadataStorages = &storage.MetadataStorages{}

	a.MetadataStorages.ResourcedMaster = storage.NewResourcedMasterMetadataStorage(a.GeneralConfig.ResourcedMaster.URL, a.GeneralConfig.ResourcedMaster.AccessToken)

	a.Db = storage.NewStorage()

	return nil
}
