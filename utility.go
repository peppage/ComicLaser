package main

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
)

type folder struct {
	URL  string `json:"url_path"`
	Name string `json:"name"`
}

func getFolders(path string) []folder {

	folders := []folder{}

	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if f.IsDir() {

			t := folder{
				URL:  "/folders/" + strconv.Itoa(len(folders)),
				Name: filepath.Join(path, f.Name()),
			}
			folders = append(folders, t)
		}
	}

	return folders
}
