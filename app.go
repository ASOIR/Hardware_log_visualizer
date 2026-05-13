package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type ParseResult struct {
	Headers    []string    `json:"headers"`
	Times      []string    `json:"times"`
	Data       [][]float64 `json:"data"` // Data[column_index][row_index]
	Files      []string    `json:"files"`
	LoadedFile string      `json:"loadedFile"`
}

// ParseCSV reads and parses the given CSV file.
// Uses goroutines to parallelize the string-to-float conversions for multi-core performance.
func (a *App) OpenFile() (*ParseResult, error) {
	filepathValue, err := wailsruntime.OpenFileDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title: "Select Hardware Log CSV",
		Filters: []wailsruntime.FileFilter{
			{DisplayName: "CSV Files", Pattern: "*.csv;*.CSV"},
		},
	})
	if err != nil {
		return nil, err
	}
	if filepathValue == "" {
		return nil, fmt.Errorf("No file selected")
	}
	return a.ParseCSV(filepathValue)
}

func (a *App) ParseCSV(targetPath string) (*ParseResult, error) {
	file, err := os.Open(targetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Some hardware monitors use malformed CSV with different fields per record.
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV has not enough rows")
	}

	// Remove BOM from the first header if present
	records[0][0] = strings.TrimPrefix(records[0][0], "\xef\xbb\xbf")

	headers := records[0]
	numRows := len(records) - 1

	// Scan siblings
	var siblingFiles []string
	dir := filepath.Dir(targetPath)
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".csv") {
			siblingFiles = append(siblingFiles, filepath.Join(dir, e.Name()))
		}
	}

	// Result structure
	result := &ParseResult{
		Headers:    []string{},
		Times:      make([]string, numRows),
		Data:       [][]float64{},
		Files:      siblingFiles,
		LoadedFile: targetPath,
	}

	// Identify timestamp columns (usually Date and Time)
	timeIndices := make([]int, 0)
	valueIndices := make([]int, 0)

	for i, h := range headers {
		if strings.TrimSpace(h) == "" {
			continue
		}
		if i == 0 || i == 1 { // Assuming first two are Date and Time
			timeIndices = append(timeIndices, i)
		} else {
			result.Headers = append(result.Headers, h)
			valueIndices = append(valueIndices, i)
		}
	}

	// Initialize data matrix [col][row] slice
	numValCols := len(valueIndices)
	result.Data = make([][]float64, numValCols)
	for i := range result.Data {
		result.Data[i] = make([]float64, numRows)
	}

	// 1. Process Timestamps (in main thread as it's string concatenation, fast enough)
	for i := 0; i < numRows; i++ {
		row := records[i+1]
		var ts []string
		for _, idx := range timeIndices {
			if idx < len(row) {
				ts = append(ts, row[idx])
			}
		}
		result.Times[i] = strings.Join(ts, " ")
	}

	// 2. Process Values using Goroutines for multi-core scaling
	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup

	rowsPerWorker := numRows / numWorkers
	if rowsPerWorker == 0 {
		rowsPerWorker = 1
		numWorkers = 1
	}
	wg.Add(numWorkers)

	for w := 0; w < numWorkers; w++ {
		startRow := w * rowsPerWorker
		endRow := startRow + rowsPerWorker
		if w == numWorkers-1 {
			endRow = numRows // Give remaining to the last worker
		}

		go func(start, end int) {
			defer wg.Done()
			for r := start; r < end; r++ {
				record := records[r+1]
				for colIdx, vIdx := range valueIndices {
					if vIdx < len(record) {
						// Clean up common hardware monitor issues e.g. "_ " or empty
						strVal := strings.TrimSpace(record[vIdx])
						if strVal == "" || strVal == "No" || strVal == "Yes" {
							if strVal == "Yes" {
								result.Data[colIdx][r] = 1
							} else {
								result.Data[colIdx][r] = 0 // default 0 or NaN, here using 0
							}
						} else {
							val, err := strconv.ParseFloat(strVal, 64)
							if err == nil {
								result.Data[colIdx][r] = val
							}
						}
					}
				}
			}
		}(startRow, endRow)
	}

	wg.Wait()

	return result, nil
}
