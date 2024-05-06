/*
inspector generates K8s cluster diagnostics reports.

Usage:

	inspector [flags]

The flags are:

	-v
	    Print out to terminal all operations.
	-n
	    Kubernetes namespace. If not provided `default` is used.
*/
package main

import (
	"os"

	"github.com/qba73/inspector"
)

func main() {
	os.Exit(inspector.Main())
}
