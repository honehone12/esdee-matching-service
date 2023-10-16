package context

type Metadata interface {
	Name() string
	Version() string
}
