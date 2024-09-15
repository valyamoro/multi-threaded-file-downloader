package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
	urls := []struct {
		url      string
		filename string
	}{
		{"http://example.com/file1.txt", "file1.xt"},
		{"http://example.com/file2.txt", "file2.txt"},
	}

	var wg sync.WaitGroup

	for _, file := range urls {
		wg.Add(1)

		go func(url, filepath string) {
			defer wg.Done()
			err := downloadFile(url, filepath)
			if err != nil {
				fmt.Printf("Ошибка при загрузке файла: %s: %v\n", filepath, url)
				return
			}

			fmt.Printf("Файл %s успешно загружен\n", filepath)
		}(file.url, file.filename)
	}

	wg.Wait()
}
