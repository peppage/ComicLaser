package setting

import (
	"github.com/pelletier/go-toml"
)

var (
	ComicFolder  string
	HttpPort     string
	LogLevel     int
	DatabaseName string
	config       *toml.TomlTree
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

	if config.Has("server.http_port") {
		HttpPort = config.Get("server.http_port").(string)
	}

	if config.Has("server.log_level") {
		LogLevel = int(config.Get("server.log_level").(int64))
	}

	if config.Has("database.name") {
		DatabaseName = config.Get("database.name").(string)
	}
}
