package main

import (
	"flag"
	"fmt"
	"mnms"
	"os"
	"strings"

	"github.com/qeof/q"
)

var Usage = func() {
	fmt.Printf("\n")
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func CheckArgs(args []string) {
	if len(args) == 0 {
		q.Q("no args exist")
		Usage()
		mnms.DoExit(1)
	}
	cmd := strings.Join(args, " ")
	q.Q(cmd)
	if strings.HasPrefix(cmd, "help") {
		fmt.Fprintf(os.Stderr, "%s\n", mnms.HelpCmd(cmd))
		mnms.DoExit(1)
	}
}
