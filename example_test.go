package tab

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestExamples(t *testing.T) {
	ExampleReadCSVFile()
	ExampleReadXLSFile()
	ExampleReadXLSXFile()
	ExampleReadHtmlTable()
}

func ExampleReadCSVFile() {
	filename := "testdata/test.csv"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
		return
	}
	os.Remove(filename)
	tabulator := NewTabulator(WithReadFrom(CSV))
	table, err := tabulator.Open(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(table)
	tabulator.Write(CSV, filename)
	tabulator.Write(XLSX, "testdata/test.xlsx")
}

func ExampleReadXLSFile() {
	filename := "testdata/test.xls"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
		return
	}
	//os.Remove(filename)
	tabulator := NewTabulator(WithReadFrom(XLS))
	table, err := tabulator.Open(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(table)
	//tabulator.Write(CSV, filename)
}

func ExampleReadXLSXFile() {
	filename := "testdata/test.xlsx"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	os.Remove(filename)
	tabulator := NewTabulator(WithReadFrom(XLSX))
	table, err := tabulator.Open(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(table)
	tabulator.Write(XLSX, filename)
}

func ExampleReadHtmlTable() {
	url := "https://baike.baidu.com/item/%E4%B8%AD%E5%9B%BD%E4%BC%81%E4%B8%9A500%E5%BC%BA/5542706?fr=aladdin"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	tabulator := NewTabulator(WithReadFrom(HtmlTable), WithTableIndex(0))
	t, err := tabulator.Open(resp.Body)
	if err != nil {
		panic(err)
	}
	var file = "testdata/html-table.csv"
	err = tabulator.Write(CSV, file)
	if err != nil {
		panic(err)
	}
	fmt.Println(t)
	os.Remove(file)
}
