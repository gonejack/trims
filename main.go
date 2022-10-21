package main

import "log"

func main() {
	err := new(trims).run()
	if err != nil {
		log.Fatal(err)
	}
}
