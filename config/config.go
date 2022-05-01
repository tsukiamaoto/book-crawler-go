package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Databases       map[string]*Database
	ServerAddress   string
	Redis           *Redis
	ProxyServerAddr string
}

type Redis struct {
	Address  string
	Password string
	DB       int
}

type Database struct {
	Name   string
	Source string
}

func LoadConfig() *Config {
	config := viper.New()
	config.AddConfigPath(".")
	config.AddConfigPath("..")
	config.SetConfigName("app")
	config.SetConfigType("yaml")

	config.AutomaticEnv()

	err := config.ReadInConfig() // Find and read the config file
	if err != nil {
		panic("讀取設定檔出現錯誤，錯誤的原因為" + err.Error())
	}

	var dbs = make(map[string]*Database)
	dbs["shopCart"] = getDatabase("shopCart", config)
	dbs["default"] = getDatabase("default", config)

	serverAddress := fmt.Sprintf("%s:%d", config.GetString("application.host"), config.GetInt("application.port"))

	redis := &Redis{
		Address:  config.GetString("redis.host"),
		Password: config.GetString("redis.password"),
		DB:       config.GetInt("redis.db"),
	}

	proxyServerAddress := config.GetString("proxyServer.host")

	return &Config{
		Databases:       dbs,
		ServerAddress:   serverAddress,
		Redis:           redis,
		ProxyServerAddr: proxyServerAddress,
	}
}

func getDatabase(name string, config *viper.Viper) *Database {
	dbName := config.GetString(fmt.Sprintf("databases.%s.dbname", name))
	source := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Taipei",
		config.GetString(fmt.Sprintf("databases.%s.host", name)),
		config.GetInt(fmt.Sprintf("databases.%s.port", name)),
		config.GetString(fmt.Sprintf("databases.%s.user", name)),
		config.GetString(fmt.Sprintf("databases.%s.password", name)),
		config.GetString(fmt.Sprintf("databases.%s.dbname", name)),
	)

	return &Database{
		Name:   dbName,
		Source: source,
	}
}
