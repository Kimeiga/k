// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"sync"
// )

// type CharacterEntry struct {
// 	Decomposition string   `json:"decomposition"`
// 	Tags          []string `json:"tags"`
// }

// type DecompositionMap map[string][]CharacterEntry

// func main() {
// 	args := os.Args[1:]

// 	var decompositions DecompositionMap
// 	var wg sync.WaitGroup
// 	var mu sync.Mutex
// 	progress := 0

// 	file, err := os.Open("ids.json")
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		return
// 	}
// 	defer file.Close()

// 	bytes, err := io.ReadAll(file)
// 	if err != nil {
// 		fmt.Println("Error reading file:", err)
// 		return
// 	}

// 	err = json.Unmarshal(bytes, &decompositions)
// 	if err != nil {
// 		fmt.Println("Error unmarshalling JSON:", err)
// 		return
// 	}

// 	if len(args) > 0 {
// 		// If there is an argument, generate a test HTML file for the given character
// 		testCharacter := args[0]
// 		generateTestFile(testCharacter, decompositions[testCharacter], decompositions)
// 	} else {
// 		buildFolder := "docs"
// 		if _, err := os.Stat(buildFolder); os.IsNotExist(err) {
// 			os.Mkdir(buildFolder, os.ModePerm)
// 		}

// 		total := len(decompositions)         // Total number of entries
// 		semaphore := make(chan struct{}, 64) // Limit to 24 goroutines

// 		for character, entries := range decompositions {
// 			wg.Add(1)
// 			semaphore <- struct{}{} // Acquire a token

// 			go func(character string, entries []CharacterEntry) {
// 				defer wg.Done()
// 				defer func() { <-semaphore }() // Release the token

// 				filename := filepath.Join(buildFolder, character+".html")
// 				file, err := os.Create(filename)
// 				if err != nil {
// 					fmt.Println("Error creating file:", err)
// 					return
// 				}
// 				defer file.Close()

// 				for _, entry := range entries {
// 					htmlContent := buildHTML(character, entry, decompositions, 0)
// 					_, err = file.WriteString(htmlContent)
// 					if err != nil {
// 						fmt.Println("Error writing to file:", err)
// 						return
// 					}
// 				}

// 				// Update progress
// 				mu.Lock()
// 				progress++
// 				fmt.Printf("\rProgress: %d/%d", progress, total)
// 				mu.Unlock()
// 			}(character, entries)
// 		}

// 		wg.Wait() // Wait for all goroutines to finish
// 		fmt.Println("\nAll entries processed.")
// 	}
// }

// var MAX_DEPTH = 24

// func buildHTML(character string, entry CharacterEntry, decompositions DecompositionMap, depth int) string {
// 	decompString, components := processDecomposition(entry.Decomposition, decompositions)
// 	tags := strings.Join(entry.Tags, ", ")

// 	var result strings.Builder
// 	borderStyle := "solid"
// 	if depth == 0 {
// 		borderStyle = "none" // No border for the top-level element
// 	}

// 	result.WriteString(fmt.Sprintf("<div style='border-left: 2px %s; padding-left: 0.5rem;'>", borderStyle))
// 	result.WriteString(fmt.Sprintf("<!DOCTYPE html>\n<html>\n<head>\n<title>%s</title>\n</head>\n<body>\n<h1>%s</h1>\n<p>Decomposition: %s</p>\n<p>Tags: %s</p>\n", character, character, decompString, tags))

// 	if len(components) > 1 {
// 		for _, comp := range components {
// 			if subEntries, ok := decompositions[comp]; ok {
// 				for _, subEntry := range subEntries {
// 					subHTML := buildHTML(comp, subEntry, decompositions, depth+1)
// 					result.WriteString(subHTML)
// 				}
// 			}
// 		}
// 	}

// 	result.WriteString("</body>\n</html>\n</div>")
// 	return result.String()
// }

// func processDecomposition(decomp string, decompositions DecompositionMap) (string, []string) {
// 	var result strings.Builder
// 	var components []string

// 	runes := []rune(decomp)
// 	for _, rune := range runes {
// 		char := string(rune)
// 		if isCJKCompositionIndicator(rune) {
// 			continue // Skip composition indicators
// 		} else {
// 			result.WriteString(fmt.Sprintf("<a href=\"%s.html\">%s</a>", char, char))
// 			components = append(components, char)
// 		}
// 	}

// 	return result.String(), components
// }

// func isCJKCompositionIndicator(char rune) bool {
// 	switch char {
// 	case '\u2FF0', '\u2FF1', '\u2FF2', '\u2FF3', '\u2FF4',
// 		'\u2FF5', '\u2FF6', '\u2FF7', '\u2FF8', '\u2FF9',
// 		'\u2FFA', '\u2FFB':
// 		return true
// 	}
// 	return false
// }

// // Function to generate a test HTML file for a specific character
// func generateTestFile(character string, entries []CharacterEntry, decompositions DecompositionMap) {
// 	filename := character + ".html"
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		fmt.Println("Error creating test file:", err)
// 		return
// 	}
// 	defer file.Close()

// 	for _, entry := range entries {
// 		htmlContent := buildHTML(character, entry, decompositions, 0)
// 		_, err = file.WriteString(htmlContent)
// 		if err != nil {
// 			fmt.Println("Error writing to test file:", err)
// 			return
// 		}
// 	}
// 	fmt.Printf("Test file generated for character: %s\n", character)
// }
