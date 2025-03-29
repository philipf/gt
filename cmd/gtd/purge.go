package gtd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/philipf/gt/internal/console"
	"github.com/philipf/gt/internal/settings"
	"github.com/spf13/cobra"
)

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge archived actions",
	Long:  `All files in the Kanban directory with the status set to Archive in the Front Matter will be deleted.`,
	Run: func(cmd *cobra.Command, _ []string) {
		purge(cmd)
	},
}

func init() {
	gtdCmd.AddCommand(purgeCmd)

	purgeCmd.Flags().BoolP("dry-run", "d", false, "Dry run")
	purgeCmd.Flags().BoolP("clean-all", "a", false, "Clean all notes, not only the Kanban board")
}

func purge(cmd *cobra.Command) {

	dryRun, err := cmd.Flags().GetBool("dry-run")
	cobra.CheckErr(err)

	cleanAll, err := cmd.Flags().GetBool("clean-all")
	cobra.CheckErr(err)

	searchDir := settings.GetKanbanGtdPath()

	if cleanAll {
		searchDir = settings.GetKanbanBasePath()
	}

	if searchDir == "" {
		fmt.Println("The GTD directory is not set in the config file")
		return
	}

	// Check if the directory exists
	if _, err := os.Stat(searchDir); os.IsNotExist(err) {
		fmt.Printf("The GTD directory [%s] does not exist\n", searchDir)
		return
	}

	// Define the search pattern for files
	filePattern := "*.md"

	// Define the content pattern to search for using regex
	contentPatternRegex := `^---\n(?:[^\n]*\n)*?status:\s*(Archive|Done)`

	re, err := regexp.Compile(contentPatternRegex)
	if err != nil {
		log.Fatalf("Failed to compile regex: %v", err)
	}

	// Declare a slice to store the results
	var filesToBeDeleted []string

	// Show the the directory to be searched
	fmt.Printf("Searching in [%s] for files matching [%s] with content matching [%s]\n", searchDir, filePattern, contentPatternRegex)

	// Walk through the directory and its subdirectories
	err = filepath.Walk(searchDir, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file matches the pattern
		if matched, _ := filepath.Match(filePattern, filepath.Base(path)); matched {
			// Read the file content
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Check if the content matches the pattern
			if re.Match(fileContent) {
				filesToBeDeleted = append(filesToBeDeleted, path)
				fmt.Println(path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fileCount := len(filesToBeDeleted)

	if fileCount == 0 {
		fmt.Println("No files found")
		return
	}

	// Prompt the user to confirm the deletion of the files
	fmt.Printf("Do you want to delete the %v file(s) in [%s]? (y/N): ", len(filesToBeDeleted), searchDir)
	confirmation, err := console.ReadLine()

	if err != nil {
		log.Fatal(err)
	}

	confirmation = strings.ToLower(strings.TrimSpace(confirmation))

	if confirmation != "y" {
		fmt.Println("Files not deleted, user cancelled")
		return
	}

	for _, file := range filesToBeDeleted {
		if dryRun {
			continue
		}

		err := os.Remove(file)
		if err != nil {
			fmt.Printf("Error deleting file %s: %s\n", file, err)
		}
	}

	if dryRun {
		fmt.Println("Dry run, files not deleted")
		return
	}
	fmt.Println("Files deleted")
}
