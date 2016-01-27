package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mattn/nicome"
)

type config struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

var num = flag.Int("n", 500, "number of comments")

func main() {
	flag.Parse()

	var file string
	if runtime.GOOS == "windows" {
		file = filepath.Join(os.Getenv("APPDATA"), "nicome", "config.json")
	} else {
		file = filepath.Join(os.Getenv("HOME"), ".config", "nicome", "config.json")
	}
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var cfg config
	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	client := nicome.NewClient(cfg.Mail, cfg.Password)
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Login(); err != nil {
		log.Fatal(err)
	}

	for _, arg := range flag.Args() {
		comments, err := client.Comments(arg, *num)
		if err != nil {
			log.Fatal(err)
		}

		for _, comment := range comments {
			fmt.Println(comment.No, comment.Text)
		}
	}
}
