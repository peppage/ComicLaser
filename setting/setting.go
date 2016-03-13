package setting

import (
	"github.com/pelletier/go-toml"
)

var (
	ComicFolder   string
	HttpPort      string
	LogLevel      string
	DatabaseName  string
	ServerLogging bool
	config        *toml.TomlTree
)

const APP_VER = "alpha4"

func init() {
	var err error
	config, err = toml.LoadFile("conf.toml")
	if err != nil {
		panic("Error loading conf.toml " + err.Error())
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
		LogLevel = config.Get("server.log_level").(string)
	}

	if config.Has("database.name") {
		DatabaseName = config.Get("database.name").(string)
	}

	if config.Has("server.server_logging") {
		ServerLogging = config.Get("server.server_logging").(bool)
	}
}
