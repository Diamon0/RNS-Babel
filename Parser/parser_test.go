package parser

import (
	"os"
	"testing"
)

// Set to true if you want to see the full print of results
var FullVerbose bool = true

func TestParseFile(t *testing.T) {
	t.Log("Testing parseFile...")

	file, err := os.OpenFile("Data/testOG.csv", os.O_RDONLY, 0444)
	if err != nil {
        t.Errorf("File testOG.csv could not be read with error:\n%v", err)
	}
    defer file.Close()
	t.Log("File opened...")

    records, err := parseFile(file)
	if err != nil {
		t.Errorf("Failed to parse file with error:\n%v", err)
	}
	t.Log("File parsed...")

    if FullVerbose {
        t.Log("Printing parsed records...")
        t.Logf("%v", records)
    }

	t.Log("parseFile Passed!")
}

func TestParseGameFiles(t *testing.T) {
    t.Log("Testing ParseGameFiles...")

    gameFiles, err := ParseGameFiles(".")
    if err != nil {
        t.Errorf("Failed to parse game files with error:\n%v", err)
    }
    t.Log("Game files parsed...")

    if FullVerbose {
        t.Log("Printing parsed game files...")
        t.Logf("%+v", gameFiles)
    }

    t.Log("ParseGameFiles Passed!")
}
