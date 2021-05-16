package main

import (
	"fmt"
	"os"

	// "fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func init() {
	fileName := "config.toml"
	// cwd, _ := os.Getwd()
	splits := strings.Split(filepath.Base(fileName), ".")
	viper.SetConfigName(filepath.Base(splits[0]))
	// viper.AddConfigPath(cwd)
	viper.AddConfigPath(filepath.Dir(fileName))
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

}

func checkKey(key string) {
	if !viper.IsSet(key) {
		fmt.Printf("Configuration key %s not found; aborting \n", key)
		os.Exit(1)
	}
}

func MustGetString(key string) string {
	checkKey(key)
	return viper.GetString(key)
}

func MustGetInt(key string) int {
	checkKey(key)
	return viper.GetInt(key)
}

func MustGetBool(key string) bool {
	checkKey(key)
	return viper.GetBool(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func Get(key string) interface{} {
	return viper.Get(key)
}

func Unmarshal(key string, rawVal interface{}) {
	viper.UnmarshalKey(key, rawVal)
}
