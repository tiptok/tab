package tab

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/extrame/xls"
	"github.com/tiptok/tab/x"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"sync"
	"unicode/utf8"
)

type Reader interface {
	Read(importer *options, r File) ([][]string, error)
	Header() HeaderInfo
}

var _ Reader = (*XLXSReader)(nil)

type XLXSReader struct {
	headerInfo HeaderInfo
}

func (reader *XLXSReader) Read(excelImport *options, f File) ([][]string, error) {
	excelFile, err := excelize.OpenReader(f)
	if err != nil {
		return nil, err
	}
	sheets := excelFile.GetSheetList()
	if len(sheets) > 0 && sheets[0] != excelImport.Sheet {
		excelImport.Sheet = sheets[0]
	}
	rows, err := excelFile.Rows(excelImport.Sheet)
	if err != nil {
		return nil, err
	}
	var (
		rowDataList = make([][]string, 0) //数据列表
		rowIndex    int                   //行计数
		lenColumn   int
	)
	for rows.Next() {
		rowIndex++
		cols, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		if rowIndex < excelImport.RowBegin {
			continue
		}
		if rowIndex == excelImport.RowBegin {
			reader.headerInfo = NewHeaderInfo(cols)
			lenColumn = len(cols)
			continue
		}
		if len(cols) > 0 {
			if len(cols) < lenColumn {
				padding := make([]string, lenColumn-len(cols))
				cols = append(cols, padding...)
			}
			rowDataList = append(rowDataList, cols)
		}
	}
	return rowDataList, nil
}
func (reader *XLXSReader) Header() HeaderInfo {
	return reader.headerInfo
}

var _ Reader = (*XLSReader)(nil)

type XLSReader struct {
	headerInfo HeaderInfo
}

func (reader *XLSReader) Read(excelImport *options, f File) ([][]string, error) {
	r, ok := f.(io.ReadSeeker)
	if !ok {
		return [][]string{}, fmt.Errorf("file must io.ReadSeeker")
	}
	wb, err := xls.OpenReader(r, "utf-8")
	if err != nil {
		return nil, err
	}
	sheet := wb.GetSheet(0)
	var (
		rowDataList = make([][]string, 0) //数据列表
		rowIndex    int                   //行计数
		lenColumn   int
	)
	for rowIndex <= int(sheet.MaxRow) {
		row := sheet.Row(rowIndex)
		cols := make([]string, 0)
		if row.LastCol() > 0 {
			for j := 0; j < row.LastCol(); j++ {
				col := row.Col(j)
				cols = append(cols, col)
			}
		}
		rowIndex++
		if rowIndex < excelImport.RowBegin {
			continue
		}
		if rowIndex == excelImport.RowBegin {
			reader.headerInfo = NewHeaderInfo(cols)
			lenColumn = len(cols)
			continue
		}
		if len(cols) == lenColumn {
			rowDataList = append(rowDataList, cols)
		}
	}
	return rowDataList, nil
}
func (reader *XLSReader) Header() HeaderInfo {
	return reader.headerInfo
}

var _ Reader = (*CSVReader)(nil)

type CSVReader struct {
	headerInfo HeaderInfo
}

func (reader *CSVReader) Read(excelImport *options, f File) ([][]string, error) {
	utf8Reader, err := reader.PrepareCheck(f)
	if err != nil {
		return nil, err
	}
	csvReader := csv.NewReader(utf8Reader)
	csvReader.FieldsPerRecord = -1
	csvReader.LazyQuotes = true

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	var (
		rowDataList = make([][]string, 0) //数据列表
		rowIndex    int                   //行计数
		lenColumn   int
	)
	for i := 0; i < len(records); i++ {
		rowIndex++
		if rowIndex < excelImport.RowBegin {
			continue
		}
		cols := records[i]
		if rowIndex == excelImport.RowBegin {
			reader.headerInfo = NewHeaderInfo(cols)
			lenColumn = len(cols)
			continue
		}
		if len(cols) > 0 {
			if len(cols) != lenColumn {
				padding := make([]string, lenColumn-len(cols))
				cols = append(cols, padding...)
			}
			rowDataList = append(rowDataList, cols)
		}
	}
	return rowDataList, nil
}
func (reader *CSVReader) Header() HeaderInfo {
	return reader.headerInfo
}
func (reader *CSVReader) PrepareCheck(r io.Reader) (io.Reader, error) {
	return GBKToUtf8(r)
}
func GBKToUtf8(readIn io.Reader) (io.Reader, error) {
	var (
		err      error
		fileByte []byte
	)
	fileByte, err = io.ReadAll(readIn)
	if err != nil {
		return nil, err
	}

	if utf8.Valid(fileByte) {
		return bytes.NewReader(fileByte), nil
	} else {
		utf8Reader := transform.NewReader(bytes.NewReader(fileByte), simplifiedchinese.GBK.NewDecoder())
		return utf8Reader, nil
	}
}

type HtmlTableReader struct {
	headerInfo HeaderInfo
	page       *x.Page
	err        error
	doInit     sync.Once
}

func (reader *HtmlTableReader) Read(options *options, f File) ([][]string, error) {
	var (
		rowDataList = make([][]string, 0) //数据列表
	)
	reader.doInit.Do(func() {
		reader.init(f)
	})
	if reader.err != nil {
		return rowDataList, reader.err
	}
	l := reader.page.Len()
	if reader.page != nil && l > 0 && l > options.TableIndex {
		t := reader.page.Tables[options.TableIndex]
		reader.headerInfo = HeaderInfo{
			Columns: t.Header,
		}
		return t.Rows, nil
	}
	return rowDataList, nil
}

func (reader *HtmlTableReader) init(f File) {
	page, err := x.New(context.Background(), f)
	if err != nil {
		reader.err = err
		return
	}
	reader.page = page
}

func (reader *HtmlTableReader) Header() HeaderInfo {
	return reader.headerInfo
}
