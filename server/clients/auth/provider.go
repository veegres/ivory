package auth

type Provider[C any, S any] interface {
	Connect(config C) error
	Verify(subject S) (string, error)
	Configured() bool
	SetConfig(config C) error
	DeleteConfig()
}
