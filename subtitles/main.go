package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ce(err error) {
	if err != nil {
		panic(err)
	}
}

type FileList []string

func getAllFiles(dir string) map[string]FileList {
	videoFiles := make(map[string]FileList)
	pat := regexp.MustCompile("(?i)(s[0-9]{2}e[0-9]{2})")
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		results := pat.FindStringSubmatch(path)
		if len(results) == 0 {
			return nil
		}
		file := filepath.Base(path)
		if strings.HasPrefix(file, ".") {
			return nil
		}
		episode := strings.ToLower(results[1])
		videoFiles[episode] = append(videoFiles[episode], file)
		return nil
	})

	return videoFiles
}

func renameGroups(episode string, files []string) {
	var videoFile string
	var subtitleFiles FileList
	for _, file := range files {
		ext := filepath.Ext(file)
		if ext == ".mkv" {
			videoFile = file
		} else if ext == ".ass" || ext == ".srt" {
			subtitleFiles = append(subtitleFiles, file)
		}
	}

	if len(videoFile) == 0 {
		// no video file found, touch nothing
		return
	}

	_, videoFilename := filepath.Split(videoFile)
	videoFileBasename := strings.TrimSuffix(videoFilename,
		filepath.Ext(videoFilename))

	for _, subtitle := range subtitleFiles {
		index := strings.LastIndex(subtitle, ".")
		if index == -1 {
			panic("Invalid index")
		}
		headPart := subtitle[:index]
		parts := strings.Split(headPart, ".")

		sz := len(parts)
		lastPart := parts[sz-1]
		if !isAllASCIICharactor(lastPart) {
			parts = parts[0 : sz-1]
		}

		newName := strings.Replace(subtitle, headPart, videoFileBasename, -1)

		if subtitle == newName {
			return
		}

		log.Printf("%s -> %s", subtitle, newName)
		os.Rename(subtitle, newName)
	}
}

func isAllASCIICharactor(s string) bool {
	runes := []rune(s)
	for _, r := range runes {
		if r >= 128 {
			return false
		}
	}
	return true
}

func main() {
	wd, err := os.Getwd()
	ce(err)
	episodeFiles := getAllFiles(wd)
	for episode, files := range episodeFiles {
		renameGroups(episode, files)
	}
}
