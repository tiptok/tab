# table data extractor for Go

## Installation

```bash
go get github.com/tiptok/tab
```

## Usage

Example read table from excel(xlsx、xls)、html-table,and write to a file 

```go
package main

import (
	"fmt"
	"os"
	"github.com/tiptok/tab"
)

func main() {
	filename := "testdata/test.xlsx"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	tabulator := tab.NewTabulator(tab.WithReadFrom(tab.XLSX))
	table, err := tabulator.Open(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(table)
	
	err = tabulator.Write(tab.CSV, filename)
	if err != nil {
		panic(err)
	}
}
```

More example In  [example_test.go](https://github.com/tiptok/tab/blob/main/example_test.go)

## Dependency

- github.com/extrame/xls 
- github.com/xuri/excelize/v2
- github.com/nfx/go-htmltable#page


