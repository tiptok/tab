package tab

import "io"

type HeaderInfo struct {
	Columns []string
}

func NewHeaderInfo(cols []string) HeaderInfo {
	return HeaderInfo{
		Columns: cols,
	}
}

type File interface {
	io.Reader
	//io.Seeker
}

type source string

const (
	XLSX      source = "xlsx"
	XLS       source = "xls"
	CSV       source = "csv"
	HtmlTable source = "html-table"
)

type options struct {
	ReadFrom source

	FileName string
	// Excel sheet name
	Sheet string
	// Excel sheet begin of row
	RowBegin int
	// Excel sheet begin of column
	ColBegin int

	// TableIndex
	TableIndex int
	// html url address
	// HtmlUrl string
	Reader Reader
}

type optionFunc func(o *options)

func WithReadFrom(f source) optionFunc {
	return func(o *options) {
		o.ReadFrom = f
	}
}

func WithFileName(s string) optionFunc {
	return func(o *options) {
		o.FileName = s
	}
}

func WithSheet(s string) optionFunc {
	return func(o *options) {
		o.Sheet = s
	}
}

func WithRowBegin(i int) optionFunc {
	return func(o *options) {
		o.RowBegin = i
	}
}

func WithColBegin(i int) optionFunc {
	return func(o *options) {
		o.ColBegin = i
	}
}

func WithTableIndex(i int) optionFunc {
	return func(o *options) {
		o.TableIndex = i
	}
}

func WithReader(r Reader) optionFunc {
	return func(o *options) {
		o.Reader = r
	}
}
