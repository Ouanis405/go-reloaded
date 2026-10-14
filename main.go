package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Kullanım: go run main.go <girdi_dosyasi> <cikis_dosyasi>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Dosya okuma hatası: %v\n", err)
		return
	}

	text := string(content)
	processedText := processText(text)

	err = os.WriteFile(outputFile, []byte(processedText), 0644)
	if err != nil {
		fmt.Printf("Dosya yazma hatası: %v\n", err)
		return
	}

	fmt.Println("Tebrikler! Proje başarıyla tamamlandı!")
}

func processText(text string) string {
	words := strings.Fields(text)
	var result []string

	for i := 0; i < len(words); i++ {
		word := words[i]

		// 1. Binary Dönüşümü
		if i+1 < len(words) && words[i+1] == "(bin)" {
			val, err := strconv.ParseInt(word, 2, 64)
			if err == nil {
				result = append(result, strconv.FormatInt(val, 10))
				i++
				continue
			}
		}

		// 2. Hex Dönüşümü
		if i+1 < len(words) && words[i+1] == "(hex)" {
			val, err := strconv.ParseInt(word, 16, 64)
			if err == nil {
				result = append(result, strconv.FormatInt(val, 10))
				i++
				continue
			}
		}

		// 3. (up) Kontrolü
		if word == "(up)" && len(result) > 0 {
			result[len(result)-1] = strings.ToUpper(result[len(result)-1])
			continue
		}

		// 4. (low) Kontrolü
		if word == "(low)" && len(result) > 0 {
			result[len(result)-1] = strings.ToLower(result[len(result)-1])
			continue
		}

		// 5. (cap) Kontrolü
		if word == "(cap)" && len(result) > 0 {
			result[len(result)-1] = capitalize(result[len(result)-1])
			continue
		}

		// 6. A / An Kuralı Kontrolü
		if strings.ToLower(word) == "a" && i+1 < len(words) {
			nextLine := strings.ToLower(words[i+1])
			if len(nextLine) > 0 {
				firstChar := nextLine[0]
				if firstChar == 'a' || firstChar == 'e' || firstChar == 'i' || firstChar == 'o' || firstChar == 'u' || firstChar == 'h' {
					if word == "A" {
						result = append(result, "An")
					} else {
						result = append(result, "an")
					}
					continue
				}
			}
		}

		result = append(result, word)
	}

	joined := strings.Join(result, " ")
	return fixPunctuation(joined)
}

// Noktalama işaretlerini düzenleyen fonksiyon
func fixPunctuation(text string) string {
	// Boşluktan sonra gelen noktalama işaretlerini birleştir
	rePunct := regexp.MustCompile(`\s+([.,\/#!$%\^&\*;:{}=\-_` + "`" + `~()]+)`)
	text = rePunct.ReplaceAllString(text, "$1")

	// Noktalama işaretinden sonra boşluk bırak
	reSpace := regexp.MustCompile(`([.,\/#!$%\^&\*;:{}=\-_` + "`" + `~()]+)([^\s])`)
	text = reSpace.ReplaceAllString(text, "$1 $2")

	// Çoklu noktaları ve grupları düzelt (örneğin ... veya !?)
	reGroups := regexp.MustCompile(`([.,?!])\s+([.,?!])`)
	for reGroups.MatchString(text) {
		text = reGroups.ReplaceAllString(text, "$1$2")
	}

	// Tek tırnak işaretlerinin etrafındaki boşlukları temizle
	reQuotes := regexp.MustCompile(`'\s*(.*?)\s*'`)
	text = reQuotes.ReplaceAllString(text, "'$1'")

	return text
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}
