package data

import (
  "log"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

func getDbConn(connStr string) *sql.DB {
  db, err := sql.Open("mysql", connStr + "?tls=false&autocommit=true")
  if err != nil {
    log.Fatal(err)
  }

  return db
}
