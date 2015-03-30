package main

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mikkyang/id3-go"
)

type AudioInfo struct {
	Title    string
	Artist   string
	Album    string
	Track    int
	Type     string
	Path     string
	File     string
	Duration time.Duration
}

func GetAudioInfo(path string) AudioInfo {
	_, name := filepath.Split(path)
	spl := strings.Split(path, ".")
	ext := spl[len(spl)-1]
	switch ext {
	case "mp3":
		f, err := id3.Open(path)
		if err != nil {
			log.Println(err)
			return AudioInfo{Title: name, File: name, Path: path}
		}
		defer f.Close()
		tf := f.Frame("TRCK")
		var itrack int
		if tf != nil {
			track := strings.Trim(tf.String(), "\x00")
			itrack, _ = strconv.Atoi(track)
		}
		return AudioInfo{
			Title:  f.Title(),
			Artist: f.Artist(),
			Album:  f.Album(),
			Track:  itrack,
			Type:   "mp3",
			File:   name,
			Path:   path,
		}
	case "ogg":
		return AudioInfo{
			Title: name,
			Type:  "oga",
			File:  name,
			Path:  path,
		}
	case "m4a":
		return AudioInfo{
			Title: name,
			Type:  "m4a",
			File:  name,
			Path:  path,
		}
		// no ogg support yet
	default:
		return AudioInfo{Title: name, File: name, Path: path}
	}
}

func isAudio(path string) bool {
	sections := strings.Split(path, ".")
	ext := strings.ToLower(sections[len(sections)-1])
	switch ext {
	case "mp3", "ogg", "m4a":
		return true
	default:
		return false
	}
}
