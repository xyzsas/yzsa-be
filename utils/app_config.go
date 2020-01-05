package utils

import "github.com/BurntSushi/toml"

type config struct {
	MongoURI     string
	DatabaseName string

	Salt string

	RedisAddr string
	RedisPWD  string

	RedisTokenId int
	RedisCacheId int
}

var Config = new(config)

func init() {
	path := "./config.toml"
	if _, err := toml.DecodeFile(path, &Config); err != nil {
		panic("Failed to load config!")
	}
}
