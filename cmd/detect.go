package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-enry/go-enry/v2"
	"github.com/simonwhitaker/gibo/utils"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var ignoreDirs = []string{".git", "tests", "build", "dist"}

var detectPath string

func init() {
	giboCmd.AddCommand(detectCmd)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("On getting current working directory: %v\n", err)
	}

	detectCmd.Flags().StringVarP(&detectPath, "path", "p", cwd, "Path to search. Defaults to current working directory.")
}

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect the languages used in the current working directory and its subdirectories",
	Run: func(cmd *cobra.Command, args []string) {
		langCount := make(map[string]int)
		err := filepath.WalkDir(detectPath, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() && slices.Contains(ignoreDirs, d.Name()) {
				return filepath.SkipDir
			}
			if !d.IsDir() {
				lang, safe := enry.GetLanguageByExtension(path)
				if safe && lang != "" {
					v, ok := langCount[lang]
					if ok {
						langCount[lang] = v + 1
					} else {
						langCount[lang] = 1
					}
				}
			}
			return nil
		})
		if err != nil {
			log.Fatalf("On walking current working directory: %v\n", err)
		}

		langs := make([]string, len(langCount))
		i := 0
		for lang := range langCount {
			langs[i] = lang
			i++
		}

		sort.Slice(langs, func(i, j int) bool {
			return langCount[langs[i]] > langCount[langs[j]]
		})

		for _, lang := range langs {
			if err := utils.PrintBoilerplate(lang); err != nil {
				fmt.Printf("# [gibo] %v\n", err)
			}
		}
	},
}
