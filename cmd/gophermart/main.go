package main

import (
	fun "github.com/AlexGrtsnk/go-prof-diploma-1/internal/functions"
)

func main() {
	if err := fun.Run(); err != nil {
		panic(err)
	}

}
