package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const empty = ""

// Парсинг флагов
func parseFlags() (string, string, error) {
	flagSrc := flag.String("src", empty, "path to URL")
	flagDest := flag.String("dest", empty, "path to output directory")
	flag.Parse()
	src, dest := *flagSrc, *flagDest
	if src == empty {
		return empty, empty, errors.New("missing src flag")
	}
	if dest == empty {
		return empty, empty, errors.New("missing dest flag")
	}
	return src, dest, nil
}

// Чтение файла
func readFile(src string) (string, error) {
	bs, err := ioutil.ReadFile(src)
	if err != nil {
		return empty, err
	}
	return string(bs), nil
}

// Парсинг ссылок (срезом)
func parseURLs(rawURLs string) []string {
	urls := strings.Split(string(rawURLs), "\n")
	for i, url := range urls {
		if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
			urls[i] = "http://" + url
		}
	}
	return urls
}

// проверка директории (возвращает её путь, если она существует)
func checkDir(dest string) (string, error) {
	outputDir, err := os.Open(dest)
	defer outputDir.Close()
	if err != nil {
		return empty, err
	}
	return outputDir.Name(), nil
}

// логирование результата
func log(url string, status bool) {
	if status {
		fmt.Printf("%s - ok\n", url)
	} else {
		fmt.Printf("%s - not ok\n", url)
	}
}

// HTTP-запрос
func getBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log(url, false)
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

// Создание файла
func makeHTMLFile(name, path string) (*os.File, error) {
	file, err := os.Create(path + "/" + name + ".html")
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Запись body в файл
func writeFile(file *os.File, body []byte) error {
	_, err := file.Write(body)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	src, dest, err := parseFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	path, err := checkDir(dest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rawURLs, err := readFile(src)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	urls := parseURLs(rawURLs)

	var wg sync.WaitGroup
	// Цикл по срезу
	for i, url := range urls {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
		
			body, err := getBody(url)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			out, err := makeHTMLFile(strconv.Itoa(i), path)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = writeFile(out, body)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			log(url, true)
		
		}(i, url)

	}

	wg.Wait()

}
