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

	// Парсинг флагов
	src := flag.String("src", empty, "path to URL")
	dest := flag.String("dest", empty, "path to output directory")
	flag.Parse()
	if *src == empty {
		fmt.Println("missing src flag")
		os.Exit(1)
	}
	if *dest == empty {
		fmt.Println("missing dest flag")
		os.Exit(1)
	}
	// Чтение файла
	bs, err := ioutil.ReadFile(*src)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Открытие директории
	outputDir, err := os.Open(*dest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outputDir.Close()
	// Парсинг файла (срезом)
	urls := strings.Split(string(bs), "\n")
	// Цикл по срезу
	for i, url := range urls {
		if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
			url = "http://" + url
		}
		// HTTP-запрос
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			fmt.Println(url, "- not ok")
			continue
		}
		defer resp.Body.Close()
		// Создание файла
		out, err := os.Create(outputDir.Name() + "/" + strconv.Itoa(i) + ".html")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer out.Close()
		// Запись в файл
		io.Copy(out, resp.Body)
		fmt.Println(url, "- ok")
	}
}
