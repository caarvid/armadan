package assert

import (
	"log"
)

func Equal[T comparable](value T, expected T, msg string) {
	if value != expected {
		log.Fatalln(msg)
	}
}

func OneOf[T comparable](value T, expected []T, msg string) {
	for _, e := range expected {
		if value == e {
			return
		}
	}

	log.Fatalln(msg)
}
