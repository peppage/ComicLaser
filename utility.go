package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

type folder struct {
	URL  string `json:"url_path"`
	Name string `json:"name"`
}

func getStructure(path string) ([]folder, int) {
	folders := []folder{}
	comicCount := 0

	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if f.IsDir() {

			t := folder{
				URL:  "/folders/" + f.Name(),
				Name: f.Name(),
			}
			folders = append(folders, t)
		} else if strings.Contains(f.Name(), "cbr") || strings.Contains(f.Name(), "cbz") {
			comicCount++
		}
	}

	return folders, comicCount
}

func getSubStructure(fullpath string, name string) ([]folder, int) {
	folders := []folder{}
	comicCount := 0

	files, _ := ioutil.ReadDir(filepath.Join(fullpath, name))
	for _, f := range files {
		if f.IsDir() {

			t := folder{
				URL:  "/folders/" + name + "/" + f.Name(),
				Name: f.Name(),
			}
			folders = append(folders, t)
		} else if strings.Contains(f.Name(), "cbr") || strings.Contains(f.Name(), "cbz") {
			comicCount++
		}
	}

	return folders, comicCount
}
