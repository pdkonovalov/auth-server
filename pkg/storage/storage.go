package storage

type Storage interface {
	Shutdown() error
	WriteNewJti(string) (string, error)
	FindJti(string) (string, bool, error)
	DeleteJti(string) error
}
