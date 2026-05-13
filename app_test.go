package main

import (
	"context"
	"path/filepath"
	"testing"
)

func TestParseCSV_Normal(t *testing.T) {
	app := NewApp()
	app.startup(context.Background())

	csvPath := filepath.Join("testdata", "dummy.csv")
	result, err := app.ParseCSV(csvPath)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatalf("Expected result, got nil")
	}

	// 1. Check Headers (Excluding Date and Time)
	expectedHeaders := []string{"CPU Core 0 [°C]", "GPU Load [%]", "Fan Speed [RPM]", "Active [Yes/No]"}
	if len(result.Headers) != len(expectedHeaders) {
		t.Errorf("Expected %d headers, got %d", len(expectedHeaders), len(result.Headers))
	}
	for i, h := range result.Headers {
		if h != expectedHeaders[i] {
			t.Errorf("Expected header %s at %d, got %s", expectedHeaders[i], i, h)
		}
	}

	// 2. Check Rows length
	expectedRows := 3
	if len(result.Times) != expectedRows {
		t.Errorf("Expected %d rows in Times, got %d", expectedRows, len(result.Times))
	}

	// 3. Check Data dimensions
	if len(result.Data) != len(expectedHeaders) {
		t.Errorf("Expected %d cols in Data, got %d", len(expectedHeaders), len(result.Data))
	}
	for i := range result.Data {
		if len(result.Data[i]) != expectedRows {
			t.Errorf("Expected Data col %d to have %d rows, got %d", i, expectedRows, len(result.Data[i]))
		}
	}

	// 4. Check specific data points
	// CPU Core 0 at row 0: 45.5
	if result.Data[0][0] != 45.5 {
		t.Errorf("Expected result.Data[0][0] = 45.5, got %v", result.Data[0][0])
	}

	// Active [Yes/No] at row 0 (Yes should be 1), row 1 (No should be 0)
	if result.Data[3][0] != 1.0 {
		t.Errorf("Expected result.Data[3][0] (Yes) = 1.0, got %v", result.Data[3][0])
	}
	if result.Data[3][1] != 0.0 {
		t.Errorf("Expected result.Data[3][1] (No) = 0.0, got %v", result.Data[3][1])
	}
}

func TestParseCSV_FileNotFound(t *testing.T) {
	app := NewApp()
	app.startup(context.Background())

	csvPath := filepath.Join("testdata", "non_existent.csv")
	_, err := app.ParseCSV(csvPath)

	if err == nil {
		t.Fatalf("Expected error for non-existent file, got nil")
	}
}

func TestParseCSV_MissingContent(t *testing.T) {
	app := NewApp()
	app.startup(context.Background())

	csvPath := filepath.Join("testdata", "dummy_missing.csv")
	result, err := app.ParseCSV(csvPath)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatalf("Expected result, got nil")
	}

	// row 1 (0-indexed) CPU is empty -> parsed as 0.0
	if result.Data[0][1] != 0.0 {
		t.Errorf("Expected empty value to parse as 0.0, got %v", result.Data[0][1])
	}

	// row 2 CPU is "No" -> parsed as 0.0
	if result.Data[0][2] != 0.0 {
		t.Errorf("Expected 'No' value to parse as 0.0, got %v", result.Data[0][2])
	}
}
