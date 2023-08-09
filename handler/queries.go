package handler 

import (
  "github.com/rembrandtcosta/podcastdownloader/models"
)

func (h *ActionHandler) insert(podcast models.Podcast) (int, error) {
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

func (h *ActionHandler) delete(podcast models.Podcast) (int, error) {
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

