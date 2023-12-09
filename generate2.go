package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type CharacterEntry struct {
	Cedict *CedictEntry `json:"cedict"`
	Ids    *IdsEntry    `json:"ids"`
}

type CedictEntry struct {
	Simplified    string     `json:"simplified"`
	Traditional   string     `json:"traditional"`
	Pronunciation []string   `json:"pronunciation"`
	Definitions   [][]string `json:"definitions"`
}

type IdsEntry struct {
	Decomposition string   `json:"decomposition"`
	Tags          []string `json:"tags"`
}

type CharacterMap map[string]CharacterEntry

func main() {
	args := os.Args[1:]

	var all_data CharacterMap
	var wg sync.WaitGroup
	var mu sync.Mutex
	progress := 0

	file, err := os.Open("combined.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	err = json.Unmarshal(bytes, &all_data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	buildFolder := "docs"
	if _, err := os.Stat(buildFolder); os.IsNotExist(err) {
		os.Mkdir(buildFolder, os.ModePerm)
	}

	if len(args) > 0 {
		// If there is an argument, generate a test HTML file for the given character
		testCharacter := args[0]
		generateTestFile(testCharacter, all_data[testCharacter], all_data)
	} else {
		total := len(all_data)               // Total number of entries
		semaphore := make(chan struct{}, 64) // Limit to 24 goroutines

		for character, entry := range all_data {
			wg.Add(1)
			semaphore <- struct{}{} // Acquire a token

			go func(character string, entry CharacterEntry) {
				defer wg.Done()
				defer func() { <-semaphore }() // Release the token

				filename := filepath.Join(buildFolder, character+".html")
				file, err := os.Create(filename)
				if err != nil {
					fmt.Println("Error creating file:", err)
					return
				}
				defer file.Close()

				htmlContent := buildHTML(character, entry, all_data)
				_, err = file.WriteString(htmlContent)
				if err != nil {
					fmt.Println("Error writing to file:", err)
					return
				}

				// Update progress
				mu.Lock()
				progress++
				fmt.Printf("\rProgress: %d/%d", progress, total)
				mu.Unlock()
			}(character, entry)
		}

		wg.Wait() // Wait for all goroutines to finish
		fmt.Println("\nAll entries processed.")
	}
}

func buildHTML(character string, entry CharacterEntry, all_data CharacterMap) string {
	var result strings.Builder
	result.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<title>" + character + "</title>\n<link rel='stylesheet' href='style.css'>\n</head>\n<body>\n")

	entryHTML := buildEntry(character, entry, all_data, 0)
	result.WriteString(entryHTML)

	result.WriteString("</body>\n</html>\n")
	return result.String()
}

const MaxDepth = 10

func findMaximalParts(character string, all_data CharacterMap, depth int) string {
	var result strings.Builder

	for i := 0; i < len(character); {
		// Find the longest matching substring in all_data starting at index i
		longestMatch := ""
		for j := i + 1; j <= len(character); j++ {
			if i == 0 && j == len(character) {
				continue // Skip the whole string
			}
			substr := character[i:j]
			if _, ok := all_data[substr]; ok {
				longestMatch = substr
			}
		}

		if longestMatch != "" {
			// Process the longest matching substring
			subEntry := all_data[longestMatch]
			subEntryHTML := buildEntry(longestMatch, subEntry, all_data, depth+1)
			result.WriteString(subEntryHTML)
			// Move the index past the processed substring
			i += len(longestMatch)
		} else {
			// No match found, move to the next character
			i++
		}
	}

	return result.String()
}

func buildEntry(character string, entry CharacterEntry, all_data CharacterMap, depth int) string {
	var result strings.Builder

	if depth == 0 {
		result.WriteString("<div class='component'>")
	} else {
		// if subcomponent, indent it
		result.WriteString("<div class='component'>")
	}

	_, components := "", []string{}

	if entry.Ids != nil {
		_, components = processDecomposition(entry.Ids.Decomposition)
	}

	runes := []rune(character)
	// Building the Cedict information
	cedict := entry.Cedict
	if entry.Cedict != nil {
		oppositeCharacter := cedict.Simplified
		if character == cedict.Simplified {
			oppositeCharacter = cedict.Traditional
		}

		if oppositeCharacter != character && oppositeCharacter != "" {
			result.WriteString(fmt.Sprintf("<h1>%s<br>%s</h1>\n", character, oppositeCharacter))
		} else {
			result.WriteString(fmt.Sprintf("<h1>%s</h1>\n", character))
		}
		// optionally display decomposition string but it might not be necessary
		// if decompString != "" {
		// 	result.WriteString(fmt.Sprintf("<p>%s</p>\n", decompString))
		// }

		if (len(cedict.Pronunciation) > 0) && (len(cedict.Definitions) > 0) {
			result.WriteString(fmt.Sprintf("<div class='entry' style='padding-left: %dpx'>", len(runes)*40))
			for i, pronunciation := range cedict.Pronunciation {
				definitions := strings.Join(cedict.Definitions[i], "; ")
				result.WriteString(fmt.Sprintf("<h2>%s</h2>\n<p>Definitions: %s</p>\n", pronunciation, definitions))
			}
			result.WriteString("</div>")
		}
	}

	// if this entry is multiple characters long, we need to run buildEntry on each of the substrings of this entry that are in the combined.json dictionary
	if depth < MaxDepth && len(runes) > 1 {
		maximalPartsHTML := findMaximalParts(character, all_data, depth)
		result.WriteString(maximalPartsHTML)
	}

	// if this entry exists, it means this is a single character, and we can recurse into its components if it has them
	if entry.Ids != nil {
		if len(components) > 1 {
			for _, comp := range components {
				if subEntry, ok := all_data[comp]; ok {
					subEntryHTML := buildEntry(comp, subEntry, all_data, depth+1)
					result.WriteString(subEntryHTML)
				}
			}
		}
	}

	// close opening div
	result.WriteString("</div>")

	// // Building the IDS sub-entries
	// if entry.Ids != nil {
	// 	components := strings.Split(entry.Ids.Decomposition, " ")
	// 	for _, comp := range components {
	// 		if subEntry, ok := all_data[comp]; ok {
	// 			subHTML := buildEntry(comp, subEntry, all_data, depth+1)
	// 			result.WriteString(subHTML)
	// 		}
	// 	}
	// }

	return result.String()
}

func processDecomposition(decomp string) (string, []string) {
	var result strings.Builder
	var components []string

	runes := []rune(decomp)
	for _, rune := range runes {
		char := string(rune)
		if isCJKCompositionIndicator(rune) {
			// continue // Skip composition indicators
			result.WriteString(char)

		} else {
			result.WriteString(fmt.Sprintf("<a href=\"%s.html\">%s</a>", char, char))
			components = append(components, char)
		}
	}

	return result.String(), components
}

func isCJKCompositionIndicator(char rune) bool {
	switch char {
	case '\u2FF0', '\u2FF1', '\u2FF2', '\u2FF3', '\u2FF4',
		'\u2FF5', '\u2FF6', '\u2FF7', '\u2FF8', '\u2FF9',
		'\u2FFA', '\u2FFB':
		return true
	}
	return false
}

func generateTestFile(character string, entry CharacterEntry, all_data CharacterMap) {
	filename := character + ".html"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating test file:", err)
		return
	}
	defer file.Close()

	htmlContent := buildHTML(character, entry, all_data) // Pass the 'all_data' parameter
	_, err = file.WriteString(htmlContent)
	if err != nil {
		fmt.Println("Error writing to test file:", err)
	}
	fmt.Printf("Test file generated for character: %s\n", character)
}
