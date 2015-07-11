package writers

func init() {
	Register("ResourcedMaster", NewResourcedMaster)
}

// NewResourcedMaster is ResourcedMaster constructor.
func NewResourcedMaster() IWriter {
	return &ResourcedMaster{}
}

// ResourcedMaster is a writer that serialize readers data to ResourcedMaster.
type ResourcedMaster struct {
	Http
}
