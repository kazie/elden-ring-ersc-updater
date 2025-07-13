package main

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestUpdateVersionInFile tests the updateVersionInFile function
func TestUpdateVersionInFile(t *testing.T) {
	// Create a temporary config file
	content := `current_version = "v1.0.0"
elden_ring_game_path = "./test_output/"
github_read_token = "test_token"
ignore_ini_file = true
`
	tmpfile, err := os.CreateTemp("", "config*.toml")
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp config file: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp config file: %v", err)
	}

	// Test the function with a recovery for os.Exit
	defer func() {
		if r := recover(); r != nil {
			// Expected panic from os.Exit
		}
	}()

	// Call the function
	updateVersionInFile("v2.0.0", tmpfile.Name())

	// Read the file and check if the version was updated
	updatedContent, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read updated config file: %v", err)
	}

	expected := `current_version = "v2.0.0"
elden_ring_game_path = "./test_output/"
github_read_token = "test_token"
ignore_ini_file = true
`
	if string(updatedContent) != expected {
		t.Errorf("Expected content to be:\n%s\nGot:\n%s", expected, string(updatedContent))
	}
}

// TestUnzipDataIntoFolder tests the unzipDataIntoFolder function
func TestUnzipDataIntoFolder(t *testing.T) {
	// Create a test zip file in memory
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// Add a file to the zip
	f, err := w.Create("test.txt")
	if err != nil {
		t.Fatalf("Failed to create file in zip: %v", err)
	}
	_, err = f.Write([]byte("test content"))
	if err != nil {
		t.Fatalf("Failed to write to file in zip: %v", err)
	}

	// Add a settings file to test ignoring
	f, err = w.Create("SeamlessCoop/ersc_settings.ini")
	if err != nil {
		t.Fatalf("Failed to create settings file in zip: %v", err)
	}
	_, err = f.Write([]byte("settings content"))
	if err != nil {
		t.Fatalf("Failed to write to settings file in zip: %v", err)
	}

	// Close the zip writer
	err = w.Close()
	if err != nil {
		t.Fatalf("Failed to close zip writer: %v", err)
	}

	// Create a temporary directory for output
	tempDir, err := os.MkdirTemp("", "test_output")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test the function with a recovery for os.Exit
	defer func() {
		if r := recover(); r != nil {
			// Expected panic from os.Exit
		}
	}()

	// Call the function
	unzipDataIntoFolder(buf.Bytes(), tempDir, true)

	// Check if the test.txt file exists
	testFilePath := filepath.Join(tempDir, "test.txt")
	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist", testFilePath)
	}

	// Check if the settings file was ignored
	settingsFilePath := filepath.Join(tempDir, "SeamlessCoop/ersc_settings.ini")
	if _, err := os.Stat(settingsFilePath); !os.IsNotExist(err) {
		t.Errorf("Expected settings file %s to be ignored", settingsFilePath)
	}
}
