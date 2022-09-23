package main

import (
	"flag"
	"github.com/tiptok/tab"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"
)

// go run main.go -f ../../testdata/test.xlsx -w test.csv -p=false

var fileName *string = flag.String("f", "testdata/test.csv", "input file")
var debug *bool = flag.Bool("d", false, "enable debug model")
var write *string = flag.String("w", "", "write to file")
var profile *bool = flag.Bool("p", false, "enable profile")

func main() {
	flag.Parse()
	begin := time.Now()
	if *profile {
		f, _ := os.OpenFile("cpu.profile", os.O_CREATE|os.O_RDWR, 0644)
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	var tabulator *tab.Tabulator = tab.NewTabulator(tab.WithFileName(file.Name()))
	table, err := tabulator.Open(file)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("read cost:", time.Since(begin))
	log.Println("file:", *fileName)
	log.Println(table)
	log.Println("last row:", table.Rows[table.Len()-1])
	if *debug {
		log.Println("rows:")
		for i := range table.Rows {
			log.Println("len", len(table.Rows[i]), table.Rows[i])
		}
	}

	if len(*write) > 0 {
		begin = time.Now()
		ext := filepath.Ext(*write)
		path := *write
		if ext == ".csv" {
			tabulator.Write(tab.CSV, path)
		} else {
			tabulator.Write(tab.XLSX, path)
		}
		log.Println("write cost:", time.Since(begin), " path:", path)
	}
}
