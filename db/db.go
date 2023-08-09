package database

import (
  "database/sql"
  "log"
)


func DbConnect() (*sql.DB, error) {
  db, err := sql.Open("sqlite3", "./database.db")
  if err != nil {
    log.Fatal(err)
  }

  sqlStmt := `
  CREATE TABLE IF NOT EXISTS podcasts (
    url VARCHAR NOT NULL PRIMARY KEY,
    name VARCHAR
  );
  `

  _, err = db.Exec(sqlStmt)
  if err != nil {
    log.Fatal(err)
  }

  return db, nil
}
