package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strconv"
)

func printUsage() {
   usage := `
    Podcast Downloader
    usage: podcastdownloader command [options]

    Options:

    Commands:
    add URL       adds the podcast from the rss feed url
    list          lists all added podcasts
    remove N      removes the n-th podcast from the list
    episodes N    lists all episodes from the n-th podcast - at most 10 a time
    update N      updates episodes from n-th podcast
    download N M  downloads m-th episode from the n-th podcast `  
   fmt.Println(usage)
   os.Exit(0)
}

func printNumberError() {
  error := `error: argument should be a number`
  fmt.Println(error)
  os.Exit(0)
}

func main() {
  cmd := os.Args[1:]

  if len(cmd) < 1 {
    printUsage()
  }

  db, err := DbConnect()
  if err != nil {
    log.Println(err)
    os.Exit(0)
  }

  handler := NewActionHandler(db)
  err = handler.InitHandler()
  if err != nil {
    log.Fatal(err)
  }

  action := cmd[0]

  if action == "add" {
    if len(cmd) < 2 {
      printUsage()
    }
    urlFeed := cmd[1]
    handler.AddPodcast(urlFeed)
  } else if action == "list" {
    handler.ListPodcast()
  } else if action == "remove" {
    if len(cmd) < 2 {
      printUsage()
    }
    arg := cmd[1]
    n, err := strconv.ParseInt(arg, 10, 64)
    if err != nil {
      printNumberError()
    }
    handler.RemovePodcast(n)
  } else if action == "episodes" {
    if len(cmd) < 2 {
      printUsage()
    }
    arg := cmd[1]
    n, err := strconv.ParseInt(arg, 10, 64)
    if err != nil {
      printNumberError()
    }
    handler.ListEpisodes(n)
  } else if action == "download" {
    if len(cmd) < 3 {
      printUsage()
    }

    arg := cmd[1:3]
    n, err := strconv.ParseInt(arg[0], 10, 64)
    if err != nil {
      printNumberError()
    }
    m, err := strconv.ParseInt(arg[1], 10, 64)
    if err != nil {
      printNumberError()
    }
    handler.DownloadEpisode(n, m)
  } else if action == "update" {
    if len(cmd) < 2 {
      printUsage()
    }

    arg := cmd[1]
    n, err := strconv.ParseInt(arg, 10, 64)
    if err != nil {
      printNumberError()
    }

    handler.UpdateEpisodes(n)
  }

}
