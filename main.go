//
// PACKAGES
//
package main

//
// IMPORTS
//
import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

//
// CONSTANTS
//
const APP_NAME = "CSVhead"
const APP_VERSION = "0.1.0"
const CLI_NAME = "csvhead"

//
// DERIVED CONSTANTS
//
var APP_LABEL string = fmt.Sprintf("%s v%s", APP_NAME, APP_VERSION)

//
// Structs
//

//
// VARIABLES
//
var ARGC int
var ARGV []string
var output_debug bool = false
var output_version bool = false
var no_header_row bool = false
var skip_lines int = 0
var head_limit int = -1

//
// FUNCTIONS
//

func init() {
	skip_lines_msg := "Specify the number of initial lines to skip (e.g. comments, copyright notices, empty rows)."

	flag.BoolVar(&output_debug, "debug", false, "Output debug info.")
	flag.BoolVar(&output_version, "version", false, "Output the version of this app.")
	flag.IntVar(&skip_lines, "K", 0, skip_lines_msg)          // SKIP_LINES
	flag.IntVar(&skip_lines, "skip-lines", 0, skip_lines_msg) // SKIP_LINES
	flag.IntVar(&head_limit, "n", -1, "How many lines to output.")

	flag.Parse()

	if output_version {
		fmt.Printf("%s", APP_LABEL)
		os.Exit(0)
	}
	ARGV = flag.Args()
	ARGC = len(ARGV)
	// log.Printf("init() | ARGC = %d | ARGV = %v\n", ARGC, ARGV)
}

//
// MAIN ENTRYPOINT
//
func main() {
	// log.Printf("%s", APP_LABEL)
	// log.Printf("head_limit = %d", head_limit)
	// log.Printf("skip_lines = %d", skip_lines)

	reader := bufio.NewReader(os.Stdin)

	r := csv.NewReader(reader)
	w := csv.NewWriter(os.Stdout)

	var i int = 0

	for {
		i += 1
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if skip_lines == 0 || i > skip_lines {
			if head_limit == -1 || i <= (head_limit+skip_lines) {
				w.Write(record)
				if err := w.Error(); err != nil {
					log.Fatalln("error writing csv:", err)
				}
			} else {
				break
			}
		}
		if i%2 == 0 {
			w.Flush()
		}
	}
	w.Flush()
}
