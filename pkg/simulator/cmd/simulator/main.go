package main

import (
	"mnms/pkg/simulator/internal/app"

	"github.com/bobbae/q"
)

func main() {
	err := app.Execute()
	if err != nil {
		q.Q(err)
	}
}
