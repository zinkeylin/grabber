package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const empty = ""

func main() {
	src := flag.String("src", empty, "path to URL")
	dest := flag.String("dest", ".", "path to output directory")
	flag.Parse()
	if *src == empty {
		fmt.Println("missing src flag")
		return
	}
 	bs, err := ioutil.ReadFile(*src)
	if err != nil {
		fmt.Println(err)
		return
	}
	outputDir, err := os.Open(*dest)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outputDir.Close()
	urls := strings.Split(string(bs), "\n")
	for i, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		out, err := os.Create(outputDir.Name() + "/" + strconv.Itoa(i) + ".txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer out.Close()
		io.Copy(out, resp.Body)
	}
}