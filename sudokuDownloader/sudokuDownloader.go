package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

func downloadSudokus(baseURL string) {
	minDifficulty := 1
	maxDifficulty := 15
	instances := []string{"a", "b", "c"}

	var wg sync.WaitGroup
	for i := minDifficulty; i <= maxDifficulty; i++ {
		for _, inst := range instances {
			filename := fmt.Sprintf("s%02d%s.txt", i, inst)
			wg.Add(1)
			go downloadSudoku(baseURL+filename, "../boards/"+filename, &wg)
		}
	}
	wg.Wait()
}

func downloadSudoku(url, filepath string, wg *sync.WaitGroup) error {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err

}

func main() {
	baseURL := "http://lipas.uwasa.fi/~timan/sudoku/"
	downloadSudokus(baseURL)
}
