package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// ProcessText applique toutes les transformations requises sur le texte.
func ProcessText(input string) string {
	// 1. Traitement des commandes entre parenthèses : (hex), (bin), (up), (low), (cap)
	text := applyCommands(input)

	// 2. Transformation de "a" en "an" devant une voyelle ou un 'h'
	text = fixArticles(text)

	// 3. Correction de la ponctuation standard (. , ! ? : ;)
	text = fixPunctuation(text)

	// 4. Correction des espaces autour des guillemets simples (')
	text = fixQuotes(text)

	return text
}

// applyCommands gère les modifications de mots basées sur les balises entre parenthèses
func applyCommands(input string) string {
	words := strings.Fields(input)
	var result []string

	i := 0
	for i < len(words) {
		word := words[i]

		// Si le mot commence par une parenthèse, c'est une commande
		if strings.HasPrefix(word, "(") {
			cmd := strings.Trim(word, "(),")

			switch cmd {
			case "hex":
				if len(result) > 0 {
					lastIdx := len(result) - 1
					val, err := strconv.ParseInt(result[lastIdx], 16, 64)
					if err == nil {
						result[lastIdx] = strconv.FormatInt(val, 10)
					}
				}
				i++
				continue

			case "bin":
				if len(result) > 0 {
					lastIdx := len(result) - 1
					val, err := strconv.ParseInt(result[lastIdx], 2, 64)
					if err == nil {
						result[lastIdx] = strconv.FormatInt(val, 10)
					}
				}
				i++
				continue

			case "up", "low", "cap":
				count := 1
				// Vérifie s'il y a un argument numérique (ex: (up, 2))
				if i+1 < len(words) && strings.HasSuffix(words[i+1], ")") {
					numStr := strings.TrimRight(words[i+1], ")")
					if n, err := strconv.Atoi(numStr); err == nil {
						count = n
						i++ // On saute le paramètre numérique
					}
				}

				// Applique la transformation sur les 'count' mots précédents
				for j := 1; j <= count && len(result)-j >= 0; j++ {
					targetIdx := len(result) - j
					result[targetIdx] = transformCase(result[targetIdx], cmd)
				}
				i++
				continue
			}
		}

		result = append(result, word)
		i++
	}

	return strings.Join(result, " ")
}

// transformCase modifie la casse d'un mot selon le mode demandé
func transformCase(word, mode string) string {
	switch mode {
	case "up":
		return strings.ToUpper(word)
	case "low":
		return strings.ToLower(word)
	case "cap":
		if len(word) == 0 {
			return word
		}
		runes := []rune(word)
		runes[0] = unicode.ToUpper(runes[0])
		for k := 1; k < len(runes); k++ {
			runes[k] = unicode.ToLower(runes[k])
		}
		return string(runes)
	default:
		return word
	}
}

// fixArticles remplace "a" par "an" si le mot suivant commence par une voyelle ou un 'h'
func fixArticles(input string) string {
	words := strings.Fields(input)
	vowels := "aeiouhAEIOUH"

	for i := 0; i < len(words)-1; i++ {
		word := words[i]
		nextWord := words[i+1]

		if len(nextWord) > 0 && strings.ContainsRune(vowels, rune(nextWord[0])) {
			if word == "a" {
				words[i] = "an"
			} else if word == "A" {
				words[i] = "An"
			}
		}
	}

	return strings.Join(words, " ")
}

// fixPunctuation colle la ponctuation au mot précédent et s'assure qu'un espace suit
func fixPunctuation(input string) string {
	// Supprime les espaces avant les signes de ponctuation
	reBefore := regexp.MustCompile(`\s+([.,!?:;])`)
	input = reBefore.ReplaceAllString(input, "$1")

	// Ajoute un espace après la ponctuation si elle est suivie d'une lettre ou d'un chiffre
	reAfter := regexp.MustCompile(`([.,!?:;]+)([a-zA-Z0-9])`)
	input = reAfter.ReplaceAllString(input, "$1 $2")

	return input
}

// fixQuotes supprime les espaces superflus à l'intérieur des guillemets simples
func fixQuotes(input string) string {
	re := regexp.MustCompile(`'\s*(.*?)\s*'`)
	return re.ReplaceAllString(input, "'$1'")
}
