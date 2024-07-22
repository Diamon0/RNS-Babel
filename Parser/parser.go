package parser

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"sync"
	"unicode"
)

// TODO: Replace for dynamic lookup inside the Data folder,
// categorize based on the name.
// Alternatively, and probably the way we should go about it, use the SheetList.csv file!!!
var KnownNameFiles []string = []string{"HBS", "Item", "Move", "Potion", "Trinkets"}
var KnownDescriptionFiles []string = []string{"HBS", "Item", "Move", "Potion", "Trinkets"}
var KnownTitleFiles []string = []string{"Character", "Enemy", "NPC"}

// Don't forget the one without any extra '_x'
var KnownStringFiles []string = []string{"Art_Description", "Art_Title", "Effect", "Intro", "Item", "Menu", "Music", "Unlock"}
var KnownStringEnumFiles []string = []string{"Dialog"}

var KnownDialogueFiles []string = []string{"birds", "dragon", "frog", "hearts", "mouse", "other", "shopkeeper", "test", "wolf"}

// For later use
type FileType int8

const (
	TypeName FileType = iota
	TypeDescription
	TypeTitle
	TypeString
	TypeStringEnum
	TypeDialogue
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
	FontSize int

	// Awaiting description.
	// The default is 3.
	// Maps to offsetAmount
	OffsetAmount int

	// Awaiting description.
	// The default is 0.
	// Maps to offsetAmountFancy
	OffsetAmountFancy int

	// Awaiting description.
	// The default is 0.
	// Maps to offsetAmountDialog
	OffsetAmountDialog int

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
		return errors.New("LanguageFile has no File referenced")
	}

	records, err := parseFile(lf.File)
	if err != nil {
		return err
	}

	lf.Languages = make([]Language, len(records[0])-1)

	// Yes I know this can all be done in a more compact way
	// But I wanted to make it clearer and less 'arcane'
	// This will however make it 10 times more painful to edit if anything changes
	// Unless you are smart with search and replace (Diamon: I use neovim btw)
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

		lf.Languages[i-1].FontSize = val
	}

	for i, offsetAmount := range records[7] {
		if i == 0 {
			continue
		}

		val, err := strconv.Atoi(offsetAmount)
		if err != nil {
			return err
		}

		lf.Languages[i-1].OffsetAmount = val
	}

	for i, offsetAmountFancy := range records[8] {
		if i == 0 {
			continue
		}

		val, err := strconv.Atoi(offsetAmountFancy)
		if err != nil {
			return err
		}

		lf.Languages[i-1].OffsetAmountFancy = val
	}

	for i, offsetAmountDialog := range records[9] {
		if i == 0 {
			continue
		}

		val, err := strconv.Atoi(offsetAmountDialog)
		if err != nil {
			return err
		}

		lf.Languages[i-1].OffsetAmountDialog = val
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

// Keeping just in case
//
// This function was welded together with duck tape and a hammer!
// Sorry, my sanity is decreasing (Diamon)
// Also, this is also built on hopes and dreams,
// as it assumes that levels will be placed in order, instead of "level 3, level 1, level 0, level 2"
//
// NEVERMIND I took a second look through the files, some aren't ordered as above,
// so for the sake of keeping consistency, this will be scrapped and I'll just go on a per-line basis;
// I would be getting in the way of other tools working after all
/*
func ParseKeyLevelString(sheet *[]*KeyLevelStrings, record [][]string) error {
    langs := make([]string, 0)
    translations := KeyLevelStrings{}

    for i := 2; i < len(record[0]); i++ {
        langs = append(langs, record[0][i])
    }

    // Loop rows, ignoring first one
    for row := 1; row < len(record); row++ {
        // Check first value for non-empty
        if record[row][0] == "" {
            *sheet = append(*sheet, nil)
            continue
        }

        for i := 0; i+2 < len(record[row]); i++ {

        }
    }
}
*/

// The individual translation for each key
type Translation struct {
	Language string
	String   string
}

type KeyLevelStrings struct {
	Key     string
	Level   int
	Strings []Translation
}

func ParseKeyLevelStrings(sheet *[]KeyLevelStrings, records [][]string) error {
	langs := make([]string, 0)

	// Populate languages
	for i := 2; i < len(records[0]); i++ {
		langs = append(langs, records[0][i])
	}

	// Check every row past the first one
	for row := 1; row < len(records); row++ {
		// Add empty row if empty (For the purposes of not messing with the file structure)
		if records[row][0] == "" {
			emptyKLS := &KeyLevelStrings{}
			*sheet = append(*sheet, *emptyKLS)
			continue
		}

		newKLS := &KeyLevelStrings{}

		newKLS.Key = records[row][0]
		val, err := strconv.Atoi(records[row][1])
		if err != nil {
			return err
		}
		newKLS.Level = val

		// Add the translations
		for lang := 0; lang+2 < len(records[row]); lang++ {
			newTranslation := &Translation{
				Language: langs[lang],
				String:   records[row][lang+2],
			}

			newKLS.Strings = append(newKLS.Strings, *newTranslation)
		}

		*sheet = append(*sheet, *newKLS)
	}

	return nil
}

// A translation with the format key,language
type KeyStrings struct {
	Key     string
	Strings []Translation
}

func ParseKeyStrings(sheet *[]KeyStrings, records [][]string) error {
	langs := make([]string, 0)

	// Populate languages
	for i := 1; i < len(records[0]); i++ {
		langs = append(langs, records[0][i])
	}

	// Check every row past the first one
	for row := 1; row < len(records); row++ {
		// Add empty row if empty (For the purposes of not messing with the file structure)
		if records[row][0] == "" {
			emptyKS := &KeyStrings{}
			*sheet = append(*sheet, *emptyKS)
			continue
		}
		newKS := &KeyStrings{}

		newKS.Key = records[row][0]

		// Add the translations
		for lang := 0; lang+1 < len(records[row]); lang++ {
			newTranslation := &Translation{
				Language: langs[lang],
				String:   records[row][lang+1],
			}

			newKS.Strings = append(newKS.Strings, *newTranslation)
		}

		*sheet = append(*sheet, *newKS)
	}

	return nil
}

// Struct for NameSheets.
// No I won't make the Strings variable into a pointer to an array,
// not for my own sanity, but for anyone who wishes to use it later.
// (Who cares about a few bytes of duplicate data)
type NameSheet struct {
	File    *os.File
	Strings []KeyLevelStrings
}

func (ns *NameSheet) Parse() error {
	if ns.File == nil {
		return errors.New("NameSheet has no File referenced")
	}

	records, err := parseFile(ns.File)
	if err != nil {
		return err
	}

	if err = ParseKeyLevelStrings(&ns.Strings, records); err != nil {
		return err
	}

	return nil
}

func (ns *NameSheet) Update() error {
	return errors.New("Not yet implemented")
}

// Struct for DescriptionSheets
type DescriptionSheet struct {
	File    *os.File
	Strings []KeyLevelStrings
}

func (ds *DescriptionSheet) Parse() error {
	if ds.File == nil {
		return errors.New("DescriptionSheet has no File referenced")
	}

	records, err := parseFile(ds.File)
	if err != nil {
		return err
	}

	if err = ParseKeyLevelStrings(&ds.Strings, records); err != nil {
		return err
	}

	return nil
}

func (ns *DescriptionSheet) Update() error {
	return errors.New("Not yet implemented")
}

// Struct for TitleSheets
type TitleSheet struct {
	File    *os.File
	Strings []KeyLevelStrings
}

func (ts *TitleSheet) Parse() error {
	if ts.File == nil {
		return errors.New("TitleSheet has no File referenced")
	}

	records, err := parseFile(ts.File)
	if err != nil {
		return err
	}

	if err = ParseKeyLevelStrings(&ts.Strings, records); err != nil {
		return err
	}

	return nil
}

func (ns *TitleSheet) Update() error {
	return errors.New("Not yet implemented")
}

// Struct for StringSheets
type StringSheet struct {
	File    *os.File
	Strings []KeyStrings
}

func (ss *StringSheet) Parse() error {
	if ss.File == nil {
		return errors.New("StringSheet has no File referenced")
	}

	records, err := parseFile(ss.File)

	if err = ParseKeyStrings(&ss.Strings, records); err != nil {
		return err
	}

	return nil
}

func (ss *StringSheet) Update() error {
	return errors.New("Not yet implemented")
}

// Struct for StringEnumSheet.
// For now, only Strings_Dialog.csv uses it,
// and it is handled the same as other string sheets
type StringEnumSheet struct {
	File    *os.File
	Strings []KeyStrings
}

func (sse *StringEnumSheet) Parse() error {
	if sse.File == nil {
		return errors.New("StringSheet has no File referenced")
	}

	records, err := parseFile(sse.File)

	if err = ParseKeyStrings(&sse.Strings, records); err != nil {
		return err
	}

	return nil
}

func (ss *StringEnumSheet) Update() error {
	return errors.New("Not yet implemented")
}

type DialogueStrings struct {
	Type           int
	FlagScript     any
	ExpressionVar0 any
	Translations   []Translation
}

func ParseDialogueStrings(sheet *[]DialogueStrings, records [][]string) error {
	langs := make([]string, 0)

	// Populate languages
	for i := 3; i < len(records[0]); i++ {
		langs = append(langs, records[0][i])
	}

	// Check every row past the first one
	for row := 1; row < len(records); row++ {
		// Add empty row if empty
		if records[row][0] == "" {
			emptyDS := &DialogueStrings{}
			*sheet = append(*sheet, *emptyDS)
			continue
		}

		newDS := &DialogueStrings{}

		val, err := strconv.Atoi(records[row][0])
		if err != nil {
			return err
		}

		newDS.Type = val
		if unicode.IsDigit(rune(records[row][1][0])) {
			val, err = strconv.Atoi(records[row][1])
			if err != nil {
				return err
			}

			newDS.FlagScript = val
		} else {
			newDS.FlagScript = records[row][1]
		}

		if unicode.IsDigit(rune(records[row][2][0])) {
			val, err = strconv.Atoi(records[row][2])
			if err != nil {
				return err
			}

			newDS.ExpressionVar0 = val
		} else {
			newDS.ExpressionVar0 = records[row][2]
		}

		for lang := 0; lang+3 < len(records[row]); lang++ {
			newTranslation := &Translation{
				Language: langs[lang],
				String:   records[row][lang+3],
			}

			newDS.Translations = append(newDS.Translations, *newTranslation)
		}

		*sheet = append(*sheet, *newDS)
	}

	return nil
}

type DialogueFile struct {
	File    *os.File
	Strings []DialogueStrings
}

func (df *DialogueFile) Parse() error {
	if df.File == nil {
		return errors.New("DialogueFile has no File referenced")
	}

	records, err := parseFile(df.File)
	if err != nil {
		return err
	}

	if err = ParseDialogueStrings(&df.Strings, records); err != nil {
		return err
	}

	return nil
}

func (df *DialogueFile) Update() error {
	return errors.New("Not yet implemented")
}

// Files used to manage language data
type LanguageFiles struct {
	Languages LanguageFile
	Sheets    []TranslationFile
	Dialogues []TranslationFile
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

// Concurrently opens and parses the language files.
func parseLanguageFileConcurrent(filePath string, fileType FileType, fileCollection *chan *TranslationFile, wg *sync.WaitGroup) {
	defer wg.Done()
    defer func() {
        if r := recover(); r != nil {}
    }()

	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var newLangFile TranslationFile

	switch fileType {
	case TypeName:
		newLangFile = &NameSheet{
			File: file,
		}
	case TypeDescription:
		newLangFile = &DescriptionSheet{
			File: file,
		}
	case TypeTitle:
		newLangFile = &TitleSheet{
			File: file,
		}
	case TypeString:
		newLangFile = &StringSheet{
			File: file,
		}
	case TypeStringEnum:
		newLangFile = &StringEnumSheet{
			File: file,
		}
	case TypeDialogue:
		newLangFile = &DialogueFile{
			File: file,
		}
	default:
		panic(errors.New("A goroutine failed to parse file: " + file.Name()))
	}

	if err = newLangFile.Parse(); err != nil {
		panic(err)
	}

	// Add the file to the channel
	*fileCollection <- &newLangFile
}

func parseLanguageFile(filePath string, fileType FileType) (*TranslationFile, error) {
	var newTranslationFile TranslationFile
    file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return &newTranslationFile, err
	}
	defer file.Close()

	switch fileType {
	case TypeName:
		newTranslationFile = &NameSheet{
			File: file,
		}
	case TypeDescription:
		newTranslationFile = &DescriptionSheet{
			File: file,
		}
	case TypeTitle:
		newTranslationFile = &TitleSheet{
			File: file,
		}
	case TypeString:
		newTranslationFile = &StringSheet{
			File: file,
		}
	case TypeStringEnum:
		newTranslationFile = &StringEnumSheet{
			File: file,
		}
	case TypeDialogue:
		newTranslationFile = &DialogueFile{
			File: file,
		}
	default:
        return &newTranslationFile, errors.New("Failed to parse file: "+file.Name())
	}

	if err = newTranslationFile.Parse(); err != nil {
		return &newTranslationFile, err
	}

    return &newTranslationFile, nil
}

func ParseGameFiles(gamePath string) (LanguageFiles, error) {
	var languageFiles LanguageFiles

	// Open and parse the Languages file
	file, err := os.OpenFile(gamePath+"/Data/LanguageEnable.csv", os.O_RDWR, 0644)
	if err != nil {
		return languageFiles, err
	}
	defer file.Close()

	languageFile := &LanguageFile{
		File: file,
	}

	if err = languageFile.Parse(); err != nil {
		return languageFiles, err
	}

	// Add the Languages file
	languageFiles.Languages = *languageFile

	// Open and parse the Name files
	for _, name := range KnownNameFiles {
        nameSheet, err := parseLanguageFile(gamePath+"/Data/Names_"+name+".csv", TypeName)
		if err != nil {
			return languageFiles, err
		}

		// Add the Names file
		languageFiles.Sheets = append(languageFiles.Sheets, *nameSheet)
	}

	// Open and parse the Description files
	for _, name := range KnownDescriptionFiles {
        descriptionSheet, err := parseLanguageFile(gamePath+"/Data/Descriptions_"+name+".csv", TypeDescription)
		if err != nil {
			return languageFiles, err
		}

		// Add the Descriptions file
		languageFiles.Sheets = append(languageFiles.Sheets, *descriptionSheet)
	}

	// Open and parse the Title files
	for _, name := range KnownTitleFiles {
        titleSheet, err := parseLanguageFile(gamePath+"/Data/Titles_"+name+".csv", TypeTitle)
		if err != nil {
			return languageFiles, err
		}

		// Add the Titles file
		languageFiles.Sheets = append(languageFiles.Sheets, *titleSheet)
	}

	// Open and parse the String files
	// First do the odd one
    stringSheet, err := parseLanguageFile(gamePath+"/Data/Strings.csv", TypeString)
	if err != nil {
		return languageFiles, err
	}

	languageFiles.Sheets = append(languageFiles.Sheets, *stringSheet)

	// Now do the other ones known
	for _, name := range KnownStringFiles {
		stringSheet, err = parseLanguageFile(gamePath+"/Data/Strings_"+name+".csv", TypeString)
		if err != nil {
			return languageFiles, err
		}

		languageFiles.Sheets = append(languageFiles.Sheets, *stringSheet)
	}

	// Same for the StringEnums
	for _, name := range KnownStringEnumFiles {
        stringEnumSheet, err := parseLanguageFile(gamePath+"/Data/Strings_"+name+".csv", TypeStringEnum)
		if err != nil {
			return languageFiles, err
		}

		languageFiles.Sheets = append(languageFiles.Sheets, *stringEnumSheet)
	}

	return languageFiles, nil
}

func ParseGameFilesConcurrent(gamePath string) (LanguageFiles, error) {
	var languageFiles LanguageFiles

	// Open and parse the Languages file
	file, err := os.OpenFile(gamePath+"/Data/LanguageEnable.csv", os.O_RDWR, 0644)
	if err != nil {
		return languageFiles, err
	}
	defer file.Close()

	languageFile := &LanguageFile{
		File: file,
	}

	if err = languageFile.Parse(); err != nil {
		return languageFiles, err
	}

	// Add the Languages file
	languageFiles.Languages = *languageFile

	// Create the Wait Group
	wg := sync.WaitGroup{}

	// Create the Sheets channel
    totalSheetCount := len(KnownNameFiles) + len(KnownDescriptionFiles) + len(KnownTitleFiles) + len(KnownStringFiles)+1 + len(KnownStringEnumFiles)
	sheetChan := make(chan *TranslationFile, totalSheetCount)
    defer close(sheetChan)

	// Open and parse the Name files
	for _, name := range KnownNameFiles {
		wg.Add(1)
		go parseLanguageFileConcurrent(gamePath+"/Data/Names_"+name+".csv", TypeName, &sheetChan, &wg)
	}

	// Open and parse the Description files
	for _, name := range KnownDescriptionFiles {
		wg.Add(1)
		go parseLanguageFileConcurrent(gamePath+"/Data/Descriptions_"+name+".csv", TypeDescription, &sheetChan, &wg)
	}

	// Open and parse the Title files
	for _, name := range KnownTitleFiles {
		wg.Add(1)
		go parseLanguageFileConcurrent(gamePath+"/Data/Titles_"+name+".csv", TypeTitle, &sheetChan, &wg)
	}

	// Open and parse the String files
	// First do the odd one
	wg.Add(1)
	go parseLanguageFileConcurrent(gamePath+"/Data/Strings.csv", TypeString, &sheetChan, &wg)

	// Now do the other ones known
	for _, name := range KnownStringFiles {
		wg.Add(1)
		go parseLanguageFileConcurrent(gamePath+"/Data/Strings_"+name+".csv", TypeString, &sheetChan, &wg)
	}

	// Same for the StringEnums
	for _, name := range KnownStringEnumFiles {
		wg.Add(1)
		go parseLanguageFileConcurrent(gamePath+"/Data/Strings_"+name+".csv", TypeStringEnum, &sheetChan, &wg)
	}

    // Wait for the group
    wg.Wait()

	// Add the sheet channel contents
    for range totalSheetCount {
        sheet, ok := <-sheetChan
        if !ok {
            break
        }
		languageFiles.Sheets = append(languageFiles.Sheets, *sheet)
	}

	return languageFiles, nil
}
