package proxy

import (
	"database/sql"
	"github.com/atomicptr/web-file-proxy/link"
)

func (p *Proxy) initDatabase() error {
	db, err := sql.Open(p.DatabaseDriver, p.DatabaseUrl)
	if err != nil {
		return err
	}

	p.linkRepository = &link.Repository{Db: db}

	return p.linkRepository.InitSchema()
}
