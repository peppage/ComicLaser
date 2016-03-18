package model

import (
	"net/url"
	"path/filepath"
	"strings"
)

type dbFolder struct {
	Name       string
	ComicCount int
}

type ComicFolder struct {
	Current string   `json:"current"`
	Folders []Folder `json:"folders"`
	Comics  `json:"comics"`
}

type Folder struct {
	URL  string `json:"url_path"`
	Name string `json:"name"`
}

type Comics struct {
	Count int    `json:"count"`
	URL   string `json:"url_path"`
}

func GetFolderData(base string, path string) (ComicFolder, error) {
	f := []dbFolder{}
	err := db.Select(&f, `SELECT folder AS name, COUNT(folder) as comiccount FROM COMICS
        WHERE folder LIKE $1 GROUP BY folder ORDER BY folder`, base+path+"%", base+path)

	var cfolder ComicFolder
	subFolderNames := make(map[string]int)
	for _, v := range f {

		if v.Name == base {
			uv := url.Values{}
			uv.Add("folder", filepath.Join(base, path))

			cfolder.Current = path
			cfolder.Count = v.ComicCount
			cfolder.URL = "/comiclist?" + uv.Encode()
		} else {
			baseRemoved := strings.Replace(v.Name, base+"\\"+path, "", -1)
			splits := strings.Split(baseRemoved, "\\")
			subFolderNames[splits[0]] = 0

		}
	}
	for k := range subFolderNames {
		u, _ := urlEncoded(k)
		cfolder.Folders = append(cfolder.Folders, Folder{
			URL:  "/folders/" + u,
			Name: k,
		})
	}

	return cfolder, err
}

// UrlEncoded encodes a string like Javascript's encodeURIComponent()
func urlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
