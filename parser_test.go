package main

import (
	"fmt"
	"github.com/Vilsol/ue4pak/parser"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseAllAsFiles(t *testing.T) {
	color.NoColor = false

	paks, err := filepath.Glob("paks/*.pak")

	if err != nil {
		panic(err)
	}

	for _, f := range paks {
		fmt.Println("Parsing file:", f)

		file, err := os.OpenFile(f, os.O_RDONLY, 0644)

		if err != nil {
			panic(err)
		}

		pak := parser.Parse(file)

		summaries := make(map[string]*parser.FPackageFileSummary, 0)

		// First pass, parse summaries
		for j, record := range pak.Index.Records {
			trimmed := strings.Trim(record.FileName, "\x00")
			if strings.HasSuffix(trimmed, "uasset") {
				fmt.Printf("Reading Record: %d: %#v\n", j, record)
				summaries[trimmed[0:strings.Index(trimmed, ".uasset")]] = record.ReadUAsset(file)
			}
		}

		// Second pass, parse exports
		for j, record := range pak.Index.Records {
			trimmed := strings.Trim(record.FileName, "\x00")
			if strings.HasSuffix(trimmed, "uexp") {
				summary, ok := summaries[trimmed[0:strings.Index(trimmed, ".uexp")]]

				if !ok {
					fmt.Printf("Unable to read record. Missing uasset: %d: %#v\n", j, record)
					continue
				}

				fmt.Printf("Reading Record: %d: %#v\n", j, record)

				record.ReadUExp(file, summary)
			}
		}
	}
}

func TestParseAllAsBytes(t *testing.T) {
	color.NoColor = false

	paks, err := filepath.Glob("paks/*.pak")

	if err != nil {
		panic(err)
	}

	for _, f := range paks {
		fmt.Println("Parsing file:", f)

		data, err := ioutil.ReadFile(f)

		if err != nil {
			panic(err)
		}

		reader := &parser.PakByteReader{
			Bytes: data,
		}

		pak := parser.Parse(reader)

		summaries := make(map[string]*parser.FPackageFileSummary, 0)

		// First pass, parse summaries
		for j, record := range pak.Index.Records {
			trimmed := strings.Trim(record.FileName, "\x00")
			if strings.HasSuffix(trimmed, "uasset") {
				fmt.Printf("Reading Record: %d: %#v\n", j, record)
				summaries[trimmed[0:strings.Index(trimmed, ".uasset")]] = record.ReadUAsset(reader)
			}
		}

		// Second pass, parse exports
		for j, record := range pak.Index.Records {
			trimmed := strings.Trim(record.FileName, "\x00")
			if strings.HasSuffix(trimmed, "uexp") {
				summary, ok := summaries[trimmed[0:strings.Index(trimmed, ".uexp")]]

				if !ok {
					fmt.Printf("Unable to read record. Missing uasset: %d: %#v\n", j, record)
					continue
				}

				fmt.Printf("Reading Record: %d: %#v\n", j, record)

				record.ReadUExp(reader, summary)
			}
		}
	}
}