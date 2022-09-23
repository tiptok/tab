package tab

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

func NewTabulator(optionsFunc ...optionFunc) *Tabulator {
	op := newOptions(optionsFunc...)
	reader := makeReaderByReadFrom(op.ReadFrom)
	if op.FileName != "" {
		reader = makeReaderByFileExt(op.FileName)
	}
	if op.Reader != nil {
		reader = op.Reader
	}

	return &Tabulator{
		reader:  reader,
		options: op,
	}
}

func makeReaderByReadFrom(from source) Reader {
	var reader Reader
	switch from {
	case XLSX:
		reader = &XLXSReader{}
	case XLS:
		reader = &XLSReader{}
	case CSV:
		reader = &CSVReader{}
	case HtmlTable:
		reader = &HtmlTableReader{}
	}
	return reader
}

func makeReaderByFileExt(filename string) Reader {
	var reader Reader
	ext := filepath.Ext(filename)
	switch ext {
	case ".xlsx":
		reader = &XLXSReader{}
	case ".xls":
		reader = &XLSReader{}
	case ".csv":
		reader = &CSVReader{}
	case ".html":
	}
	return reader
}

func newOptions(optionsFunc ...optionFunc) *options {
	op := &options{ReadFrom: XLSX, Sheet: "Sheet1", RowBegin: 1, ColBegin: 1}
	for i := range optionsFunc {
		optionsFunc[i](op)
	}
	return op
}

func makeWriterByWriteTo(to source, table *Table) Writer {
	var w Writer
	switch to {
	case XLSX:
		w = &XLXSWriterTo{data: table.Rows, title: table.HeaderInfo.Columns}
	case CSV:
		w = &CSVWriterTo{data: table.Rows, title: table.HeaderInfo.Columns}
	default:
		w = &XLXSWriterTo{data: table.Rows, title: table.HeaderInfo.Columns}
	}
	return w
}

type Tabulator struct {
	reader  Reader
	options *options
	Table   *Table
}

func (tab *Tabulator) Open(r File) (*Table, error) {
	rows, err := tab.Reader().Read(tab.options, r)
	if err != nil {
		return nil, err
	}
	header := tab.Reader().Header()
	tab.Table = &Table{
		HeaderInfo: header,
		Rows:       rows,
	}
	return tab.Table, nil
}

func (tab *Tabulator) Write(to source, path string) error {
	wt := makeWriterByWriteTo(to, tab.Table)
	return wt.Save(path)
}

func (tab *Tabulator) WriteTo(to source, w io.Writer) (n int64, err error) {
	wt := makeWriterByWriteTo(to, tab.Table)
	return wt.WriteTo(w)
}

func (tab *Tabulator) Reader() Reader {
	return tab.reader
}

type Table struct {
	HeaderInfo HeaderInfo
	Rows       [][]string
}

func (table *Table) Len() int {
	return len(table.Rows)
}

func (table *Table) String() string {
	return fmt.Sprintf("Table[%s] (%d rows)", strings.Join(table.HeaderInfo.Columns, ", "), len(table.Rows))
}
