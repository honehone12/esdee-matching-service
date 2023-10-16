package context

type BasicMetadata struct {
	name    string
	version string
}

func NewBasicMetadata(name string, version string) *BasicMetadata {
	return &BasicMetadata{
		name:    name,
		version: version,
	}
}

func (b *BasicMetadata) Name() string {
	return b.name
}

func (b *BasicMetadata) Version() string {
	return b.version
}
