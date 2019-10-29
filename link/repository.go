package link

import (
	"database/sql"
)

type Repository struct {
	Db *sql.DB
}

func (r *Repository) InitSchema() error {
	// create table
	stmt, err := r.Db.Prepare("create table if not exists links (" +
		"uid integer primary key autoincrement," +
		"identifier varchar(255)," +
		"url text," +
		"content_type varchar(255))")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	// create index
	stmt, err = r.Db.Prepare("create unique index if not exists ident_index on links (identifier)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()

	return err
}

func (r *Repository) FindAll() ([]*Link, error) {
	rows, err := r.Db.Query("select uid, identifier, url, content_type from links")
	if err != nil {
		return nil, err
	}

	var links []*Link

	for rows.Next() {
		link := Link{}

		err = rows.Scan(&link.Uid, &link.Identifier, &link.Url, &link.ContentType)
		if err != nil {
			return nil, err
		}

		links = append(links, &link)
	}

	return links, nil
}

func (r *Repository) FindByIdentifier(ident string) (*Link, error) {
	rows, err := r.Db.Query("select uid, identifier, url, content_type from links where identifier = ?", ident)
	if err != nil {
		return nil, err
	}

	link := Link{}

	rows.Next()
	err = rows.Scan(&link.Uid, &link.Identifier, &link.Url, &link.ContentType)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (r *Repository) InsertNew(identifier, url, contentType string) error {
	stmt, err := r.Db.Prepare("insert into links(identifier, url, content_type) values(?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(identifier, url, contentType)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteByUid(uid int) error {
	stmt, err := r.Db.Prepare("delete from links where uid = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(uid)
	if err != nil {
		return err
	}

	return nil
}
