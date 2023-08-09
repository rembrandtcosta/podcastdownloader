package main

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

func (h *ActionHandler) insert(podcast Podcast) (int, error) {
  res, err := h.Database.Exec("INSERT INTO podcasts VALUES(?, ?);",
			      podcast.URL, podcast.Name)
  if err != nil {
    return 0, err 
  }

  var id int64
  if id, err = res.LastInsertId(); err != nil {
    return 0, err 
  }

  return int(id), nil
}

func (h *ActionHandler) delete(podcast Podcast) (int, error) {
  res, err := h.Database.Exec("DELETE FROM podcasts WHERE url = ?;",
			      podcast.URL)
  if err != nil {
    return 0, err 
  }

  var id int64
  if id, err = res.LastInsertId(); err != nil {
    return 0, err 
  }

  return int(id), nil
}
