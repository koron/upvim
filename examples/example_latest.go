package main

import (
	"fmt"
	"log"

	"github.com/koron/upvim/internal/github"
)

func main() {
	r, err := github.Latest("vim", "vim-win32-installer")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v", r)
}
