package setting

import (
	"github.com/pelletier/go-toml"
)

var (
	ComicFolder string
	config      *toml.TomlTree
)

const APP_VER = "alpha1"

func init() {
	var err error
	config, err = toml.LoadFile("conf.toml")
	if err != nil {
		panic("Missing conf.toml")
	}
}

func Initialize() {
	if config.Has("server.comic_folder") {
		ComicFolder = config.Get("server.comic_folder").(string)
	}
}
