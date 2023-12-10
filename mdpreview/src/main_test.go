package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestParseContent(t *testing.T) {
	input := []byte("# Hello, World!\n\n## Another heading\n\n  - List item 1\n  - List item 2\n")

	var expectedBuffer bytes.Buffer
	expectedBuffer.WriteString(header)
	expectedBuffer.WriteString("<h1>Hello, World!</h1>\n\n<h2>Another heading</h2>\n\n<ul>\n<li>List item 1</li>\n<li>List item 2</li>\n</ul>\n")
	expectedBuffer.WriteString(footer)

	var outputBuffer bytes.Buffer
	outputBuffer.WriteString(string(parseContent(input)))

	if !bytes.Equal(expectedBuffer.Bytes(), outputBuffer.Bytes()) {
		t.Errorf("Expected\n\n %s, but got\n\n %s", expectedBuffer.String(), outputBuffer.String())
	}
}

func TestRun(t *testing.T) {
	// Create a temporary file with markdown content
	tmpFile, err := os.CreateTemp("", "mdpreview_*.md")
	outName := fmt.Sprintf("%s.html", filepath.Base(tmpFile.Name()))
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}
	defer os.Remove(tmpFile.Name())

	markdownContent := []byte("# Hello, World!\n\n## Another heading\n\n  - List item 1\n  - List item 2\n")
	if _, err := tmpFile.Write(markdownContent); err != nil {
		t.Fatalf("Failed to write markdown content to temporary file: %s", err)
	}

	// Call the run function with the temporary file
	err = run(tmpFile.Name())
	if err != nil {
		t.Errorf("Failed to run: %s", err)
	}

	// TODO: Add assertions to check the output or side effects of the run function
	// Check if file exists
	_, err = os.Stat(outName)
	if os.IsNotExist(err) {
		t.Errorf("File %s was not created", outName)
	}

	//compare the the expected output to content of the outName file
	expectedBuffer := parseContent(markdownContent)
	content, err := os.ReadFile(outName)
	if err != nil {
		t.Errorf("Failed to read file: %s", err)
	}

	if !bytes.Equal(content, expectedBuffer) {
		t.Errorf("Expected %s, but got %s", expectedBuffer, string(content))
	}

	// Clean up
	os.Remove(outName)

}

func TestSaveHTML(t *testing.T) {
	testFileName := "test.html"

	var testDataBuffer bytes.Buffer
	testDataBuffer.WriteString(header)
	testDataBuffer.WriteString("<h1>Hello, World!</h1>\n")
	testDataBuffer.WriteString(footer)

	err := saveHTML(testFileName, testDataBuffer.Bytes())
	if err != nil {
		t.Errorf("Failed to save HTML: %s", err)
	}

	// Check if file exists
	_, err = os.Stat(testFileName)
	if os.IsNotExist(err) {
		t.Errorf("File %s was not created", testFileName)
	}

	// Check if file content is correct
	content, err := os.ReadFile(testFileName)
	if err != nil {
		t.Errorf("Failed to read file: %s", err)
	}

	if !bytes.Equal(content, testDataBuffer.Bytes()) {
		t.Errorf("Expected %s, but got %s", testDataBuffer.String(), string(content))
	}

	// Clean up
	os.Remove(testFileName)
}
