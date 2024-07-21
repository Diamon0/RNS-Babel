package parser

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
)

// Interface for any file where text is changed according to the language.
type TranslationFile interface {
	Parse() error
	Update() error
}

// A struct describing a Language
type Language struct {
	// The name of the language.
	// Maps to Lang
	Name string

	// Name of the language in its own language.
	// Maps to Desc
	NativeName string

	// Whether to enable the language.
	// The default is 1.
	// Maps to enabled
	Enabled bool

	// Whether it uses an external font.
	// The default is 0.
	// Maps to externalFont
	ExternalFont bool

	// Name of the external font; set ExternalFont to true if you are using this.
	// Maps to font
	FontName string

	// Whether the language uses full width characters, such as Japanese.
	// The default is 0.
	// Maps to fullWidth
	FullWidth bool

	// The size of the font.
	// The default is 55.
	// Maps to fontSize
	FontSize uint

	// Awaiting description.
	// The default is 3.
	// Maps to offsetAmount
	OffsetAmount uint

	// Awaiting description.
	// The default is 0.
	// Maps to offsetAmountFancy
	OffsetAmountFancy uint

	// Awaiting description.
	// The default is 0.
	// Maps to offsetAmountDialog
	OffsetAmountDialog uint

	// The width to assign each character.
	// The default is 40.
	// Maps to characterWidth
	CharacterWidth int

	// The width to assign each fancy character.
	// The default is 56.
	// Maps to characterWidthFancy
	CharacterWidthFancy int

	// The width to assign each character when used in dialogue.
	// The default is 40.
	// Maps to characterWidthDialog
	CharacterWidthDialogue int
}

// Struct for the file where you set language data
type LanguageFile struct {
	File      *os.File
	Languages []Language
}

func (lf *LanguageFile) Parse() error {
    if lf.File == nil {
        return errors.New("Struct has no File referenced")
    }

	records, err := parseFile(lf.File)
	if err != nil {
		return err
	}

	lf.Languages = make([]Language, len(records[0])-1)

	// Yes I know this can all be done in a more compact way
	// But I wanted to make it clearer and less 'arcane'
    // This will however make it 10 times more painful to edit if anything changes
    // Unless you are smart with search and replace (I use neovim btw)
	for i, language := range records[0] {
		if i == 0 {
			continue
		}

		lf.Languages[i-1].Name = language
	}

	for i, desc := range records[1] {
		if i == 0 {
			continue
		}

		lf.Languages[i-1].NativeName = desc
	}

	for i, enabled := range records[2] {
		if i == 0 {
			continue
		}

        val, err := strconv.ParseBool(enabled)
        if err != nil {
            return err
        }

		lf.Languages[i-1].Enabled = val
	}

	for i, externalFont := range records[3] {
		if i == 0 {
			continue
		}

        val, err := strconv.ParseBool(externalFont)
        if err != nil {
            return err
        }

		lf.Languages[i-1].ExternalFont = val
	}

	for i, font := range records[4] {
		if i == 0 {
			continue
		}

		lf.Languages[i-1].FontName = font
	}

	for i, fullWidth := range records[5] {
		if i == 0 {
			continue
		}

        val, err := strconv.ParseBool(fullWidth)
        if err != nil {
            return err
        }

		lf.Languages[i-1].FullWidth = val
	}

	for i, fontSize := range records[6] {
		if i == 0 {
			continue
		}

        val, err := strconv.Atoi(fontSize)
        if err != nil {
            return err
        }

		lf.Languages[i-1].FontSize = uint(val)
	}

	for i, offsetAmount := range records[7] {
		if i == 0 {
			continue
		}

        val, err := strconv.Atoi(offsetAmount)
        if err != nil {
            return err
        }

		lf.Languages[i-1].OffsetAmount = uint(val)
	}

	for i, offsetAmountFancy := range records[8] {
		if i == 0 {
			continue
		}

        val, err := strconv.Atoi(offsetAmountFancy)
        if err != nil {
            return err
        }

		lf.Languages[i-1].OffsetAmountFancy = uint(val)
	}

	for i, offsetAmountDialog := range records[9] {
		if i == 0 {
			continue
		}

        val, err := strconv.Atoi(offsetAmountDialog)
        if err != nil {
            return err
        }

		lf.Languages[i-1].OffsetAmountDialog = uint(val)
	}

	for i, characterWidth := range records[10] {
		if i == 0 {
			continue
		}

        val, err := strconv.Atoi(characterWidth)
        if err != nil {
            return err
        }

		lf.Languages[i-1].CharacterWidth = val
	}

	for i, characterWidthFancy := range records[11] {
		if i == 0 {
			continue
		}

        val, err := strconv.Atoi(characterWidthFancy)
        if err != nil {
            return err
        }

		lf.Languages[i-1].CharacterWidthFancy = val
	}

	for i, characterWidthDialog := range records[12] {
		if i == 0 {
			continue
		}

        val, err := strconv.Atoi(characterWidthDialog)
        if err != nil {
            return err
        }

		lf.Languages[i-1].CharacterWidthDialogue = val
	}

    return nil
}

// The individual translation for each key
type Translation struct {
	Language string
	String   string
}

// A translation with the format key,level,language
type KeyLevelStrings struct {
	// Ok, so, the key is the key,
	// the array index is the level,
	// and the translation is just that
	Translation map[string][]Translation
}

// A translation with the format key,language
type KeyStrings struct {
	Translation map[string]Translation
}

// TODO:
// Pending implementation
// (Diamon)
type DialogueStrings struct{}

// Struct for NameSheets
type NameSheet struct {
	File    *os.File
	Strings []KeyLevelStrings
}

// Struct for DescriptionSheets
type DescriptionSheet struct {
	File    *os.File
	Strings []KeyLevelStrings
}

// Struct for TitleSheets
type TitleSheet struct {
	File    *os.File
	Strings []KeyLevelStrings
}

// Struct for StringSheets
type StringSheet struct {
	File    *os.File
	Strings []KeyStrings
}

// Struct for StringSheetEnum
// For now, only Strings_Dialog.csv uses it
type StringSheetEnum struct {
	File    *os.File
	Strings []KeyStrings
}

// TODO:
// Pending Implementation
// (Diamon)
type DialogueFile struct {
	File *os.File
}

// Files used to manage language data
type LanguageFiles struct {
	Languages LanguageFile
	Sheets    TranslationFile

	// TODO:
	// Pending Implementation
	// (Diamon)
	Dialogues TranslationFile
}

// Used to initially parse CSV files
//
// I may or may not consider changing reimplementing this later.
// Probably not,
// but who knows
func parseFile(file *os.File) ([][]string, error) {
	r := csv.NewReader(file)

	records, err := r.ReadAll()
	if err != nil {
		return records, err
	}

	return records, nil
}

func ParseGameFiles(gamePath string) (LanguageFiles, error) {
    var languageFiles LanguageFiles

    file, err := os.OpenFile(gamePath+"/Data/LanguageEnable.csv", os.O_RDWR, 0644)
    if err != nil {
        return languageFiles, err
    }
    defer file.Close()

    languageFile := LanguageFile{
        File: file,
    }

    err = languageFile.Parse()
    if err != nil {
        return languageFiles, err
    }

    languageFiles.Languages = languageFile

    return languageFiles, nil
}
