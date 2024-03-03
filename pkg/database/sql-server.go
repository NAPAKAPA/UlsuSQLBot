package database

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"net/url"
)

type MsSql struct {
	database *sql.DB
}

func (drv *MsSql) Open() (sqlDb *sql.DB, err error) {
	query := url.Values{}
	query.Add("app name", "MyAppName")
	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword("sa", "admini12!"),
		Host:     fmt.Sprintf("%s:%d", "localhost", 1433),
		RawQuery: query.Encode(),
	}
	drv.database, err = sql.Open("sqlserver", u.String())
	if err != nil {
		return nil, err
	}
	return
}

func (drv *MsSql) Exec(sql string, args ...any) (err error) {
	_, err = drv.Open()
	if err != nil {
		return err
	}
	db := drv.database
	if db == nil {
		return
	}
	if len(args) > 0 {
		_, err = db.Exec(sql, args)
		return db.Close()
	}
	_, err = db.Exec(sql)
	if err != nil {
		return
	}
	return db.Close()
}
