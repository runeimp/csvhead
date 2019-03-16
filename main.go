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
const APP_VERSION = "0.4.0"
const CLI_NAME = "csvhead"
const ERROR_INVALID_ARGUMENT = 1

//
// DERIVED CONSTANTS
//
var APP_LABEL string = fmt.Sprintf("%s v%s", APP_NAME, APP_VERSION)

const HELP_MSG = APP_NAME + " v" + APP_VERSION + `

POSIX head utility for tabular data

USAGE: csvhead [OPTIONS] file1 file2...

OPTIONS:
 -n COUNT | --lines COUNT
                 Number of lines to output
 -COUNT          Shortcut for -n
 -CountOptions   Count is a decimal number optionally followed by a size letter ('b', 'k', 'm' for blocks, Kilobytes or Megabytes), or 'l' to mean count by lines, or other option letters ('cqv').
 -c COUNT | --bytes COUNT
                 Number of characters to output
 -h | --help     Output this help info
 -K COUNT | --skip-lines COUNT
                 Specify the number of lines to skip (e.g. comments, copyright notices, empty rows).
 -q | --quiet | --silent
                 Never print file name headers.
 -v | --verbose  Always print file name headers.
 -V | --version  Output the version number of this app

This tool was inspired by and is designed to work along with csvkit and similar tools.
`

//
// Structs
//

//
// VARIABLES
//
var ARGC int
var ARGV []string
var characters_maxed bool = false
var character_limit int = 0
var csv_files []string = []string{}
var head_limit int = 10
var no_header_row bool = false
var output_file_headers int = 1
var output_help bool = false
var output_version bool = false
var skip_lines int = 0

//
// FUNCTIONS
//
func csv_parser(reader io.Reader, line_count int, char_count int) (int, int) {
	r := csv.NewReader(reader)
	w := csv.NewWriter(os.Stdout)

	var c int = 0
	var i int = 0

	if char_count > 0 {
		c = char_count
	}
	// log.Printf("csv_parser() | i = %d | char_count = %d | character_limit = %d", i, char_count, character_limit)

	for {
		i += 1
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// log.Printf("csv_parser() | i = %d | skip_lines = %d | head_limit = %d | characters_maxed = %v", i, skip_lines, head_limit, characters_maxed)

		if skip_lines == 0 || i > skip_lines {
			if i <= (head_limit + skip_lines) {
				if character_limit > 0 {
					for index, field := range record {
						// log.Printf("csv_parser() | %3d | field = %v", len(field), field)
						if characters_maxed == false && c+len(field) <= character_limit {
							c += len(field)
							if c == character_limit {
								characters_maxed = true
								break
							}
						} else {
							characters_maxed = true
							if c < character_limit {
								max := character_limit - c
								// log.Printf("csv_parser() | %3d | max = %v", len(field), max)
								field = field[:max]
								// log.Printf("csv_parser() | %d | c = %d | %3d | max = %v | field = '%s'", character_limit, c, len(field), max, field)
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
				if characters_maxed {
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

	return i, c
}

func init() {
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-c", "--bytes":
			i++
			parse_posix_group(os.Args[i])
			char_limit, _ := strconv.Atoi(os.Args[i])
			character_limit = char_limit
		case "-h", "--help":
			output_help = true
		case "-K", "--skip-lines":
			i++
			skipper, _ := strconv.Atoi(os.Args[i])
			skip_lines = skipper
		case "-n", "--lines":
			i++
			limiter, _ := strconv.Atoi(os.Args[i])
			head_limit = limiter
		case "-V", "--version":
			output_version = true
		case "-q", "--quiet", "--silent":
			output_file_headers = 0
		case "-v", "--verbose":
			output_file_headers = 2
		default:
			match, _ := regexp.MatchString(`-[^-][^-]+`, os.Args[i])
			if match {
				parse_posix_group(os.Args[i])
			} else {
				match, _ := regexp.MatchString(`-\d+`, os.Args[i])
				if match {
					limiter, _ := strconv.Atoi(os.Args[i][1:])
					head_limit = limiter
				} else {
					csv_files = append(csv_files, os.Args[i])
				}
			}
		}
	}

	if output_help {
		fmt.Println(HELP_MSG)
		os.Exit(0)
	}

	if output_version {
		fmt.Printf("%s\n", APP_LABEL)
		os.Exit(0)
	} else if output_file_headers == 2 && len(csv_files) == 0 {
		fmt.Printf("%s\n", APP_LABEL)
		os.Exit(ERROR_INVALID_ARGUMENT)
	}

}

func parse_posix_group(group string) {
	r, _ := regexp.Compile(`(-?(\d+)([bklm]?))?([cqv]*)`)
	match := r.FindStringSubmatch(group)

	// log.Printf("parse_posix_group() |    group = %v\n", group)
	// log.Printf("parse_posix_group() |    match = %v\n", match)
	count := match[2]
	modifier := match[3]
	options := match[4]
	char_limit, _ := strconv.Atoi(count)
	// log.Printf("parse_posix_group() |    count = %v\n", count)
	// log.Printf("parse_posix_group() | modifier = %v\n", modifier)
	// log.Printf("parse_posix_group() |  options = %v\n", options)
	switch modifier {
	case "b": // Blocks
		character_limit = char_limit * 512
	case "k": // KiloBytes
		character_limit = char_limit * 1024
	case "l": // Lines
		head_limit = char_limit
	case "m": // MegaBytes
		character_limit = char_limit * 1048576
	default:
		character_limit = char_limit
	}
	// log.Printf("parse_posix_group() | character_limit = %d\n", character_limit)
	// log.Printf("parse_posix_group() | head_limit = %d\n", head_limit)
	for i, v := range options {
		log.Printf("parse_posix_group() | i = %d | v = %c (%T)\n", i, v, v)
		switch v {
		case 'c':
			if character_limit == 0 && head_limit > 0 {
				character_limit = head_limit
			}
		case 'q':
			output_file_headers = 0
		case 'v':
			output_file_headers = 2
		}
	}
}

//
// MAIN ENTRYPOINT
//
func main() {
	// log.Printf("%s", APP_LABEL)
	// log.Printf("head_limit = %d", head_limit)
	// log.Printf("skip_lines = %d", skip_lines)
	// log.Printf("csv_files = %q", csv_files)
	csv_count := len(csv_files)
	char_count := 0
	line_count := 0

	if csv_count > 0 {
		for _, file := range csv_files {
			if output_file_headers == 2 || output_file_headers == 1 && csv_count > 1 {
				header := fmt.Sprintf("==> %s <==\n", file)
				char_count = len(header) + char_count
				fmt.Printf(header)
				line_count += 1
			}
			csv_file, _ := os.Open(file)
			reader := bufio.NewReader(csv_file)
			line_count, char_count = csv_parser(reader, line_count, char_count)
			if characters_maxed {
				os.Exit(0)
			}
		}
		log.Printf("line_count = %d | char_count = %d", line_count, char_count)
	} else {
		reader := bufio.NewReader(os.Stdin)
		csv_parser(reader, 0, 0)
	}
}
