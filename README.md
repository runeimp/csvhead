CSVhead v0.3.0
==============

POSIX head for tablular data inspired by [csvkit][].

Written in Go. I love Python. But I prefere distributing binaries with no dependencies.

```bash
$ csvhead -h
CSVhead v0.3.0

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

```

### Example usage

```bash
$ cat example.csv | csvhead -3
LAST,FIRST,EMAIL,STATUS,NOTES
"Addams, esq.",Gomez,gomez.addams@deathmail.com,ready,"line 1
line 2
line 3
line 4
line 5"
Addams,Morticia,Morticia.Addams@deathmail.com,waiting,"No one knows how long we've kept her waiting.
We can only beg forginess at this point."
```

Note that even though the first record has a multi-line field only the first three records are output.

```bash
$ cat example.csv | csvhead -3 -K 1
"Addams, esq.",Gomez,gomez.addams@deathmail.com,ready,"line 1
line 2
line 3
line 4
line 5"
Addams,Morticia,Morticia.Addams@deathmail.com,waiting,"No one knows how long we've kept her waiting.
We can only beg forginess at this point."
Vinterslaus,Milo,milov@aol.com,activated,Milo seems unstoppable. Go Milo!
```

And here with `-K` and `--skip-lines` implimenting part of the common arguments shared across all [csvkit][] tools.

### ToDo

* [ ] Allow reading of CSV files specified as an argument to csvhead. Not just stdin.
* [ ] Implement the rest of the [csvkit][] common arguments
	* [ ] `-d DELIMITER, --delimiter DELIMITER`
	* [ ] `-t, --tabs`
	* [ ] `-q QUOTECHAR, --quotechar QUOTECHAR`
	* [ ] `-u {0,1,2,3}, --quoting {0,1,2,3}`
	* [ ] `-b, --no-doublequote`
	* [ ] `-p ESCAPECHAR, --escapechar ESCAPECHAR`
	* [ ] `-z FIELD_SIZE_LIMIT, --maxfieldsize FIELD_SIZE_LIMIT`
	* [ ] `-e ENCODING, --encoding ENCODING`
	* [ ] `-L LOCALE, --locale LOCALE`
	* [ ] `-S, --skipinitialspace`
	* [ ] `--blanks`
	* [ ] `--date-format DATE_FORMAT`
	* [ ] `--datetime-format DATETIME_FORMAT`
	* [ ] `-H, --no-header-row`
	* [x] `-K SKIP_LINES, --skip-lines SKIP_LINES`
	* [ ] `-v, --verbose`
	* [ ] `-l, --linenumbers`
	* [ ] `--zero`
	* [x] `-V, --version`






[csvkit]: https://csvkit.readthedocs.io/

