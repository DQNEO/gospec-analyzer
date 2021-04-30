package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func showUsage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("tsv2json file.tsv > file.json")
}
func main() {
	if len(os.Args) == 1 {
		showUsage()
		return
	}
	fname := os.Args[1]
	r, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	tsv, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(tsv), "\r\n")
	var mp map[string]string = make(map[string]string, len(lines))
	for _, line := range lines {
		if len(line) == 0 { // EOF
			break
		}
		fields := strings.Split(line, "\t")
		mp[fields[0]] = fields[1]
	}
	bytes, err := json.MarshalIndent(mp, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(bytes))
}
