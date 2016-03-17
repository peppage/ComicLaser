package main

import (
	"io/ioutil"
	"net/url"
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

			u, _ := UrlEncoded("/folders/" + f.Name())
			t := folder{
				URL:  u,
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

			u, _ := UrlEncoded("/folders/" + f.Name())
			t := folder{
				URL:  u,
				Name: f.Name(),
			}
			folders = append(folders, t)
		} else if strings.Contains(f.Name(), "cbr") || strings.Contains(f.Name(), "cbz") {
			comicCount++
		}
	}

	return folders, comicCount
}

// UrlEncoded encodes a string like Javascript's encodeURIComponent()
func UrlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
