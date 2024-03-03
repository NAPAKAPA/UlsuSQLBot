package database

import (
	"testing"
)

var (
	scripts = []string{
		"DROP DATABASE IF EXISTS ulsu;",
		"CREATE DATABASE ulsu;",
		"CREATE TABLE ulsu.dbo.Table_1 (Field INT);",
		"CREATE TABLE ulsu.dbo.Table_2 (Field INT);",
		"ALTER TABLE ulsu.dbo.Table_1 ADD column_2 INT;",
		"ALTER TABLE ulsu.dbo.Table_1 DROP COLUMN column_2;",
	}
)

func TestMsSql_Open(t *testing.T) {
	drv := MsSql{}
	_, err := drv.Open()
	if err != nil {
		t.Error(err)
		return
	}
	for _, script := range scripts {
		err = drv.Exec(script)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
