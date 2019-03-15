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
	"regexp"
	"strconv"
)

//
// CONSTANTS
//
const APP_NAME = "CSVhead"
const APP_VERSION = "0.3.0"
const CLI_NAME = "csvhead"
const help_msg = APP_NAME + " v" + APP_VERSION + `

POSIX head utility for tabular data

OPTIONS:
 -n COUNT        Number of lines to output
 -COUNT          Shortcut for -n
 -c COUNT        Number of characters to output
 -h | --help     Output this help info
 -K COUNT | --skip-lines COUNT
                 Specify the number of lines to skip (e.g. comments, copyright notices, empty rows).
 -v | --version  Output the version number of this app

This tool was inspired by and is designed to work along with csvkit and similar tools.
`

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
var output_help bool = false
var output_version bool = false
var no_header_row bool = false
var character_limit int = 0
var skip_lines int = 0
var head_limit int = 10

//
// FUNCTIONS
//

func init() {
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-c":
			i++
			char_limit, _ := strconv.Atoi(os.Args[i])
			character_limit = char_limit
		case "-h", "--help":
			output_help = true
		case "-K", "--skip-lines":
			i++
			skipper, _ := strconv.Atoi(os.Args[i])
			skip_lines = skipper
		case "-n":
			i++
			limiter, _ := strconv.Atoi(os.Args[i])
			head_limit = limiter
		case "-V", "--version":
			output_version = true
		default:
			match, _ := regexp.MatchString(`-\d+`, os.Args[i])
			if match {
				limiter, _ := strconv.Atoi(os.Args[i][1:])
				head_limit = limiter
			}
		}
	}

	if output_help {
		fmt.Println(help_msg)
		os.Exit(0)
	}

	if output_version {
		fmt.Printf("%s\n", APP_LABEL)
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

	var c int = 0
	var i int = 0
	end_here := false

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
			if i <= (head_limit + skip_lines) {
				if character_limit > 0 {
					for index, field := range record {
						log.Printf("%3d | field = %v", len(field), field)
						if end_here == false && c+len(field) <= character_limit {
							c += len(field)
							if c == character_limit {
								end_here = true
								break
							}
						} else {
							end_here = true
							if c < character_limit {
								max := character_limit - c
								// log.Printf("%3d | max = %v", len(field), max)
								field = field[:max]
								log.Printf("%d | c = %d | %3d | max = %v | field = '%s'", character_limit, c, len(field), max, field)
								record[index] = field
								c += max
							} else {
								record = record[:index]
								break
							}
						}

						c += 1
					}
					c -= 1
				}
				// log.Printf("record = %v", len(record))
				w.Write(record)
				if err := w.Error(); err != nil {
					log.Fatalln("error writing csv:", err)
				}
				if end_here {
					break
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
