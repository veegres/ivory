package auth

type Provider[C any, S any] interface {
	DeleteConfig()
	SetConfig(config C) error
	Verify(subject S) (string, error)
	Connect(config C) error
}
