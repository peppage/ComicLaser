package main

import (
	"io/ioutil"
	"path/filepath"
)

type folder struct {
	URL  string `json:"url_path"`
	Name string `json:"name"`
}

func getStructure(path string) []folder {
	folders := []folder{}

	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if f.IsDir() {

			t := folder{
				URL:  "/folders/" + f.Name(),
				Name: f.Name(),
			}
			folders = append(folders, t)
		}
	}

	return folders
}

func getSubStructure(fullpath string, name string) []folder {
	folders := []folder{}

	files, _ := ioutil.ReadDir(filepath.Join(fullpath, name))
	for _, f := range files {
		if f.IsDir() {

			t := folder{
				URL:  "/folders/" + name + "/" + f.Name(),
				Name: f.Name(),
			}
			folders = append(folders, t)
		}
	}

	return folders
}
