package utils

import (
	"flag"
	"log"
)

var (
	Verbose = flag.Bool("v", false, "Enable verbose logging")
)

// VPrint prints verbose logs if enabled
func VPrint(format string, v ...interface{}) {
	if *Verbose {
		log.Printf(format, v...)
	}
}
