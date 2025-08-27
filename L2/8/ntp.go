package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

const (
	DefaultNTPServer = "pool.ntp.org"
)

func main() {
	ntpTime, err := ntp.Time(DefaultNTPServer)
	if err != nil {
		// Exit with error and write error to stderr
		fmt.Fprintf(os.Stderr, "Error getting time from %s: %v\n", DefaultNTPServer, err)
		os.Exit(1)
	}

	// Write time to stdout (explicitly)
	fmt.Fprintf(os.Stdout, "%v\n", (ntpTime))
}
