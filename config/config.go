package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	DB     *DBConfig     `json:"db"`
	Image  *ImageConfig  `json:"image"`
	Email  *EmailConfig  `json:"email"`
	Assets *AssetsConfig `json:"assets"`
}

type DBConfig struct {
	Dialect  string `json:"dialect"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Charset  string `json:"charset"`
}

type ImageConfig struct {
	Url          string `json:"url"`
	AwsS3Region  string `json:"awss3region"`
	AwsS3Bucket  string `json:"awss3bucket"`
	AwsId        string `json:"awsid"`
	AwsSecretKey string `json:"awssecretkey"`
}

type EmailConfig struct {
	Server   string `json:"server"`
	Port     string `json:"port"`
	Username string `jason:"username"`
	Password string `json:"password"`
	Sender   string `json:"sender"`
	To       string `json:"to"`
}

type AssetsConfig struct {
	Static   string `json:"static"`
	Template string `json:"template"`
}

func GetConfig() *Config {
	config := Config{}
	file, err := os.Open("config.development.json")
	if err != nil {
		log.Println(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Println(err)
	}
	return &config
}
