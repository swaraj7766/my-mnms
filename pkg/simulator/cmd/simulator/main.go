package main

import (
	"mnms/pkg/simulator/internal/app"

	"github.com/qeof/q"
)

func main() {
	err := app.Execute()
	if err != nil {
		q.Q(err)
	}
}
