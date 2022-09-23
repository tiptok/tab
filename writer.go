package tab

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"os"
	"path/filepath"
)

type Writer interface {
	WriteTo(w io.Writer) (n int64, err error)
	Save(fileName string) error
}

type XLXSWriterTo struct {
	data  [][]string
	title []string
}

func (wt *XLXSWriterTo) WriteTo(w io.Writer) (n int64, err error) {
	var file *excelize.File
	file, err = wt.newFile()
	if err != nil {
		return 0, nil
	}
	return file.WriteTo(w)
}

func (wt *XLXSWriterTo) Save(fileName string) error {
	var file *excelize.File
	var err error
	if err = checkPath(fileName); err != nil {
		return err
	}
	file, err = wt.newFile()
	if err != nil {
		return nil
	}
	return file.SaveAs(fileName)
}

func checkPath(fileName string) error {
	dir := filepath.Dir(fileName)
	if !Exists(dir) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (wt *XLXSWriterTo) newFile() (*excelize.File, error) {
	sheet := "Sheet1"
	file := excelize.NewFile()
	streamWriter, err := file.NewStreamWriter(sheet)
	if err != nil {
		return nil, err
	}

	if len(wt.title) == 0 {
		return nil, fmt.Errorf("未设置数据表头")
	}
	if err := streamWriter.SetRow("A1", stringsToInterfaces(wt.title)); err != nil {
		return nil, err
	}
	var rowID = 2
	for i := 0; i < len(wt.data); i++ {
		row := stringsToInterfaces(wt.data[i])
		cell, _ := excelize.CoordinatesToCellName(1, rowID)
		if err := streamWriter.SetRow(cell, row); err != nil {
			return nil, err
		}
		rowID += 1
	}

	if err := streamWriter.Flush(); err != nil {
		return nil, err
	}
	return file, nil
}

func stringsToInterfaces(input []string) []interface{} {
	output := make([]interface{}, len(input))
	for i, v := range input {
		output[i] = v
	}
	return output
}

func NewXLXSWriterTo(title []string, data [][]string) *XLXSWriterTo {
	return &XLXSWriterTo{
		data:  data,
		title: title,
	}
}

type CSVWriterTo struct {
	data  [][]string
	title []string
}

func (xw *CSVWriterTo) WriteTo(w io.Writer) (n int64, err error) {
	var file = bytes.NewBuffer(nil)
	_, err = xw.write(file)
	if err != nil {
		return 0, nil
	}
	return file.WriteTo(w)
}

func (xw *CSVWriterTo) Save(fileName string) error {
	if err := checkPath(fileName); err != nil {
		return err
	}
	csvFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer csvFile.Close()
	if _, err := xw.write(csvFile); err != nil {
		return err
	}
	return nil
}

func (xw *CSVWriterTo) write(w io.Writer) (*csv.Writer, error) {
	_, err := w.Write([]byte("\xEF\xBB\xBF")) //写入UTF-8 BOM
	if err != nil {
		return nil, err
	}
	csvWriter := csv.NewWriter(w)
	if err := csvWriter.Write(xw.title); err != nil {
		return nil, err
	}
	if err := csvWriter.WriteAll(xw.data); err != nil {
		return nil, err
	}
	csvWriter.Flush()
	return csvWriter, csvWriter.Error()
}

func NewCSVWriterTo(title []string, data [][]string) *CSVWriterTo {
	return &CSVWriterTo{
		data:  data,
		title: title,
	}
}
