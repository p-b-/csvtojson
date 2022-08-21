# csvtojson

## Description

This utility is a simple utility to convert CSV files into JSON files.  This utility was born from data being passed to engineers as Excel spreadsheets, which need to be converted to JSON.

It is expected that there is a header row at the start of the CSV file containing column headers.

It works in two modes.  The first is a simple mode that produces a flat JSON structure, with each CSV row forming an element in a JSON array, each JSON element consisting of key/value pairs from the columm header/CSV row.

The second is a template mode that uses a template file to form an output JSON file, with
 each CSV row being formed from the template.

## Installation

Download the source code and run `go build` from a terminal inside the directory containing the go.mod file.  Add the resulting executable to a utility directory that is already part of your PATH, or, add this directory to your PATH.

Or just run it from the directory containing the built executable.

## How to use

### Simple Mode

Each JSON object produced is a simple flat structure consisting of each column as an object key/value pair.  For instance, the following CSV excerpt:

`id:number, book name, author, book length, genre`
`1, Catch 22, "Heller, Jospher", not long enough, satire`

produces:

`{`
`	"books": [{`
``
`			"id": 1,`
`			"book name": "Catch 22",`
`			"author": "Heller, Jospher",`
`			"book length": "not long enough",`
`			"genre": "satire"`
`		},`
`... `

The above example, based on books, is created using the following command

`csvtojson.exe -i="TestData\simple.csv" -o="TestData\books.json" -h=1 -a=books`

-i and -o detail the input and output files.
-h details the number or csv header rows
-a details the json array identifier for the json output file


### Template mode

Here the output is based upon a JSON template input file such as:

`{`
`	"Books": [`
`%Row Start%`
`		{"Book Id": %id%,`
`			"Title": %book name%,`
`            "Author": %author%,`
`            "Book Details": {`
`                "Length": %book length%,`
`                "Genre": %genre%`
`            }`
`	    }`
`%Row End%`
`	]`
`}`

From a CSV file containing data such as:

`id:number, book name, author, book length, genre`
`1, Catch 22, "Heller, Jospher", not long enough, satire`


this will produce output such as the following:

`{`
`	"Books": [`
`		{"Book Id": 1,`
`			"Title": "Catch 22",`
`            "Author": "Heller, Jospher",`
`            "Book Details": {`
`                "Length": "not long enough",`
`                "Genre": "satire"`
`            }`
`	    }`
`	]`
`}`

`csvtojson.exe -i="TestData\books.csv" -o="TestData\books.json" -t="TestData\books_template.txt"`

-i and -o detail the input and output files.
-t details the template file


### CSV File

Note, that the CSV file must consist of a header row and one or more rows of data.

The header columns can be delimited by "", which is necessary if the headers have commas in them.
When output to JSON file, each item is typically surrounded by "".  However, if the header comma has :bool or :number or :bool at the end, such as id:number, then they will be output without these.

The row items can be delimited by "" if they contain commas.


## Limitations

Currently, the template mode will output all values that are in the template, even when they are empty.  It relies on the template being able to form valid JSON.

## Future enhancements

A proper templating system that output correct JSON, and allows array elements and key/values to be omitted if the CSV value is empty.  However, the complexity of the resulting templating language may be counterproductive to actual use.


## License

This software is licensed under the GNU GPL V3.  See the LICENSE file for more details.