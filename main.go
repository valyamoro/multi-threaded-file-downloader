package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("не удалось загрузить файл: %v", err)
	}
	defer resp.Body.Close()

	outFile, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("не удалось записать в файл: %v", err)
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Ошибка загрузки .env файла: %v\n", err)
		return
	}

	urlsEnv := os.Getenv("FILE_URLS")
	if urlsEnv == "" {
		fmt.Println("Переменная окружения FILE_URLS не установлена")
		return
	}

	urls := strings.Split(urlsEnv, ",")
	var downloadItems []struct {
		url      string
		filename string
	}
	for _, item := range urls {
		parts := strings.Split(item, "|")
		if len(parts) != 2 {
			fmt.Printf("Некорректный формат в строке: %s\n", item)
			continue
		}
		downloadItems = append(downloadItems, struct {
			url      string
			filename string
		}{parts[0], parts[1]})
	}

	var wg sync.WaitGroup
	for _, file := range downloadItems {
		wg.Add(1)
		go func(url, filepath string) {
			defer wg.Done()
			err := downloadFile(url, filepath)
			if err != nil {
				fmt.Printf("Ошибка при загрузке файла %s: %v\n", filepath, url)
				return
			}
			fmt.Printf("Файл %s успешно загружен\n", filepath)
		}(file.url, file.filename)
	}

	wg.Wait()
}
