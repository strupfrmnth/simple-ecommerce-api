package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/strupfrmnth/simple-ecommerce-api/pkg/route"
)

type Configuration struct {
	Port string `json:"port"`
	Rps  int    `json:"rps"`
}

func init() {
	log.SetPrefix("[GIN] ")
}

func main() {
	fmt.Println("Start Application")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// load config file
	file, err := os.Open("configs/config.json")
	defer file.Close()
	if err != nil {
		log.Fatalln(err, "file open error")
	}

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatalln(err, "decode error")
	}
	fmt.Println("port:", configuration.Port)
	fmt.Println("rps:", configuration.Rps)

	log.Fatal(route.RunAPI(configuration.Port, configuration.Rps))
}
