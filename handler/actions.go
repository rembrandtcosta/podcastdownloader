package handler

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/rembrandtcosta/podcastdownloader/models"
)

type ActionHandler struct { 
  Podcasts []models.Podcast
  Database *sql.DB
}

func NewActionHandler(database *sql.DB) *ActionHandler {
  return &ActionHandler{
    Podcasts: make([]models.Podcast, 0),
    Database: database,
  }
}

func (h *ActionHandler) InitHandler() error {
  rows, err := h.Database.Query("SELECT * FROM podcasts")
  if err != nil {
    return err
  }
  defer rows.Close()

  data := []models.Podcast{}
  for rows.Next() {
    p := models.Podcast{}
    err = rows.Scan(&p.URL, &p.Name)
    if err != nil {
      return nil
    }
    p.Feed, err = updateFeed(p.URL)
    if err != nil {
      return nil
    }
    data = append(data, p)
  }

  h.Podcasts = data
  return nil
}

func (h *ActionHandler) AddPodcast(rssFeedURL string) {
  rss, err := updateFeed(rssFeedURL)
  if err != nil {
    log.Printf("Error decode: %v\n", err)
    return 
  }

  fmt.Printf("Channel title: %v added to podcast list\n" , rss.Channel.Title)
 
  podcast := models.Podcast{
    Name: rss.Channel.Title,
    URL: rssFeedURL,
    Feed: rss,
  }

  h.Podcasts = append(h.Podcasts, podcast)
  h.insert(podcast)
}

func updateFeed(rssFeedURL string) (models.Rss, error) { 
  rss := models.Rss{}

  resp, err := http.Get(rssFeedURL) 
  if err != nil {
    log.Printf("Error GET: %v\n", err)
    return rss, err
  }
  defer resp.Body.Close()

  decoder := xml.NewDecoder(resp.Body)
  err = decoder.Decode(&rss)
  if err != nil {
    return rss, err
  } 

  return rss, nil
}

func (h *ActionHandler) ListPodcast() {
  i := 0
  for _, podcast := range h.Podcasts  {
    fmt.Printf("%d. %s\n", i, podcast.Name)
    i++
  }
}

func (h *ActionHandler) ListEpisodes(n int64) {
  if n > int64(len(h.Podcasts))-1 {
    fmt.Println("no such podcast")
    return
  }
  podcast := h.Podcasts[n]
  for i, item := range podcast.Feed.Channel.Items {
    fmt.Printf("%v. %v\n", i, item.Title)
    if (i+1)%10 == 0 {
      fmt.Print("...")
      fmt.Scanln()
      fmt.Printf("\033[1A\033[K")
    }
  }
}
 
func (h *ActionHandler) RemovePodcast(n int64) {
  if n > int64(len(h.Podcasts))-1 {
    fmt.Println("no such podcast")
    return
  }
  
  h.delete(h.Podcasts[n])
  h.Podcasts = append(h.Podcasts[:n], h.Podcasts[n+1:]...) 
}

func (h *ActionHandler) UpdateEpisodes(n int64) {
  if n > int64(len(h.Podcasts))-1 {
    fmt.Println("no such podcast")
    return
  }
  rssFeedURL := h.Podcasts[n].URL
  rss, err := updateFeed(rssFeedURL)
  if err != nil {
    log.Printf("Error decode: %v\n", err)
    return 
  }
    
  podcast := models.Podcast{
    Name: rss.Channel.Title,
    URL: rssFeedURL,
    Feed: rss,
  }

  h.Podcasts[n] = podcast
}

func (h *ActionHandler) DownloadEpisode(n int64, m int64) {
  if n > int64(len(h.Podcasts))-1 {
    fmt.Println("no such podcast")
    return
  }

  podcast := h.Podcasts[n]

  episodes := podcast.Feed.Channel.Items

  if m > int64(len(episodes))-1 {
    fmt.Println("no such episode")
    return
  }

  for i, item := range podcast.Feed.Channel.Items {
    fmt.Printf("%v. %v\n", i, item.Enclosure.Url)
  }

  episode := episodes[m]

  episodeURL := episode.Enclosure.Url
  episodeName := episode.Title

  downloadFile(episodeURL, episodeName)
}

func downloadFile(URL, fileName string) error {
  resp, err := http.Get(URL)
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  if resp.StatusCode != 200 {
    fmt.Println("error: received non 200 response code")
  }

  file, err := os.Create(fileName + ".mp3")
  if err != nil {
    return err 
  }
  defer file.Close()

  _, err = io.Copy(file, resp.Body)
  if err != nil {
    return err 
  }

  return nil
}

