package main

import "io/ioutil"
import "fmt"

func main() {
	html, err := ioutil.ReadFile("./spec.html")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(html))
}
