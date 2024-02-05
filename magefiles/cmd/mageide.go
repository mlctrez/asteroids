package main

import (
	"fmt"
	"github.com/magefile/mage/mage"
	"os"
	"strings"
	"time"
)

// main allows GoLand ide to run mage targets.
func main() {
	start := time.Now()
	exitCode := mage.Main()
	fmt.Printf("mage %s took %0.2f seconds\n", strings.Join(os.Args[1:], ", "), time.Since(start).Seconds())
	os.Exit(exitCode)
}
