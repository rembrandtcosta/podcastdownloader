package handler

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	
	"os"
	"testing"
)

func TestAddPodcast(t *testing.T) {
  db, err := mockDb()
  if err != nil {
    os.Remove("./test.db")
    t.Errorf("failed to mock db")
    return
  }
  handler := NewActionHandler(db)
  err = handler.AddPodcast("https://podcasts.files.bbci.co.uk/p02nrvz8.rss")
  if err != nil {
    t.Errorf("couldn't add podcast")
  }
  os.Remove("./test.db")
  return
}

func TestListEpisodes(t *testing.T) {
  db, err := mockDb()
  if err != nil {
    os.Remove("./test.db")
    t.Errorf("failed to mock db")
    return
  }
  handler := NewActionHandler(db)
  handler.AddPodcast("https://podcasts.files.bbci.co.uk/p02nrvz8.rss")
  err = handler.ListEpisodes(0)
  if err != nil {
    t.Errorf("couldn't list podcasts")
  }
  os.Remove("./test.db")
  return
}

func mockDb() (*sql.DB, error) {
  db, err := sql.Open("sqlite3", "./test.db")
  if err != nil {
    return nil, err 
  }

  sqlStmt := `
  CREATE TABLE IF NOT EXISTS podcasts (
    url VARCHAR NOT NULL PRIMARY KEY,
    name VARCHAR
  );
  `

  _, err = db.Exec(sqlStmt)
  if err != nil {
    return nil, err 
  }

  return db, nil
}

