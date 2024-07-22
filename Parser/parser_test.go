package parser

import (
	"os"
	"testing"
)

// Set to true if you want to see the full print of results
var FullVerbose bool = false

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
        t.Log("Printing Language file:")
        t.Log(gameFiles.Languages.Languages)
        t.Log("Printing Sheet files:")
    }

    t.Log("Printing sheets:")
    for _, sheet := range gameFiles.Sheets {
        switch v := sheet.(type) {
            case *NameSheet:
                t.Logf("Found Name Sheet: %v", v.File.Name())
            case *DescriptionSheet:
                t.Logf("Found Description Sheet: %v", v.File.Name())
            case *TitleSheet:
                t.Logf("Found Title Sheet: %v", v.File.Name())
            default:
                t.Errorf("Found Sheet of Unknown Type (HOW?): %+v", v)
        }

        if FullVerbose {
            t.Logf("%+v", sheet)
        }
    }

    t.Log("ParseGameFiles Passed!")
}
