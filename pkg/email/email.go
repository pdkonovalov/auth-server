package email

import (
	"github.com/pdkonovalov/auth-server/pkg/config"
	"github.com/pdkonovalov/auth-server/pkg/storage"
)

type Email struct {
}

func Init(config *config.Config, storage storage.Storage) (*Email, error) {
	return nil, nil
}

func (*Email) SendAllert(guid string) {

}
