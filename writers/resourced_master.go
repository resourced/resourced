package writers

// NewResourcedMaster is ResourcedMaster constructor.
func NewResourcedMaster() *ResourcedMaster {
	rm := &ResourcedMaster{}
	return rm
}

// ResourcedMaster is a writer that serialize readers data to ResourcedMaster.
type ResourcedMaster struct {
	Http
}
