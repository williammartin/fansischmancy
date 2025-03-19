package main

import (
	"fmt"
	"io"
	"os"

	"github.com/williammartin/fansischmancy"
)

func main() {
	writer := fansischmancy.NewWriter(os.Stdout)
	_, err := io.Copy(writer, os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
