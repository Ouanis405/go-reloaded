package main

import (
	"fmt"
	"os"
)

func main() {
	// Vérification du nombre d'arguments
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Erreur : Veuillez fournir un fichier d'entrée et un fichier de sortie.")
		fmt.Println("Utilisation : go run . <fichier_entree.txt> <fichier_sortie.txt>")
		os.Exit(1)
	}

	inputFile := args[0]
	outputFile := args[1]

	// Lecture du fichier d'entrée
	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Erreur lors de la lecture du fichier %s : %v\n", inputFile, err)
		os.Exit(1)
	}

	// Traitement du texte
	processedText := ProcessText(string(content))

	// Écriture du résultat dans le fichier de sortie
	err = os.WriteFile(outputFile, []byte(processedText+"\n"), 0644)
	if err != nil {
		fmt.Printf("Erreur lors de l'écriture dans le fichier %s : %v\n", outputFile, err)
		os.Exit(1)
	}

	fmt.Println("Traitement terminé avec succès ! Fichier généré :", outputFile)
}
