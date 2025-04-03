package report

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/philipf/gt/cmd/togglcmd"
	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/cobra"
)

var reportService toggl.ReportService
var textOutput string
var jsonOutput string
var interactive bool

// reportCmd represents the action command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Report Service",
	Long:  `Report Service	`,

	Run: func(cmd *cobra.Command, args []string) {
		// Check if interactive mode is enabled
		if interactive {
			runInteractiveReport(cmd)
			return
		}
		
		// Default non-interactive mode
		runReport(cmd)
	},
}

// runReport executes the report in non-interactive mode
func runReport(cmd *cobra.Command) {
	sd, ed := togglcmd.GetDateRange(cmd)

	rpt, err := reportService.BuildReport(sd, ed)
	cobra.CheckErr(err)

	textReport, err := (cmd.Flags().GetBool("text"))
	cobra.CheckErr(err)

	jsonReport, err := (cmd.Flags().GetBool("json"))
	cobra.CheckErr(err)

	if textReport || textOutput != "" || (!textReport && !jsonReport) {
		r, err := reportService.BuildStringReport(rpt)
		cobra.CheckErr(err)

		if textReport {
			fmt.Println(*r)
		}

		if textOutput != "" {
			os.WriteFile(textOutput, []byte(*r), 0644)
		}
	}

	if jsonReport || jsonOutput != "" {
		jsonBytes, err := reportService.BuildJsonReport(rpt)
		cobra.CheckErr(err)

		if jsonReport {
			fmt.Println(string(*jsonBytes))
		}

		if jsonOutput != "" {
			os.WriteFile(jsonOutput, *jsonBytes, 0644)
		}
	}
}

// runInteractiveReport executes the report in interactive mode
func runInteractiveReport(cmd *cobra.Command) {
	fmt.Println("Interactive Report Mode")
	
	// Variables to store form values
	var selectedDateOption string
	var startDateStr, endDateStr string
	var outputText, outputJson bool
	var outputTextFile, outputJsonFile string
	
	// Date range options
	dateOptions := []huh.Option[string]{
		huh.NewOption("Today", "today"),
		huh.NewOption("Yesterday", "yesterday"),
		huh.NewOption("Day before yesterday", "eyesterday"),
		huh.NewOption("This week", "thisweek"),
		huh.NewOption("Last week", "lastweek"),
		huh.NewOption("Custom dates", "custom"),
	}
	
	// Step 1: Choose date range
	dateForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Date Range").
				Options(dateOptions...).
				Value(&selectedDateOption),
		),
	)
	
	err := dateForm.Run()
	cobra.CheckErr(err)
	
	// Step 2: If custom date range is selected, get start and end dates
	var sd, ed time.Time
	
	if selectedDateOption == "custom" {
		// Format for display
		today := time.Now().Format("2006/01/02")
		
		customDateForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Start Date (YYYY/MM/DD)").
					Placeholder(today).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("Start date is required")
						}
						_, err := time.ParseInLocation("2006/01/02", s, time.Local)
						if err != nil {
							return fmt.Errorf("Invalid date format. Please use YYYY/MM/DD")
						}
						return nil
					}).
					Value(&startDateStr),
					
				huh.NewInput().
					Title("End Date (YYYY/MM/DD)").
					Placeholder(today).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("End date is required")
						}
						_, err := time.ParseInLocation("2006/01/02", s, time.Local)
						if err != nil {
							return fmt.Errorf("Invalid date format. Please use YYYY/MM/DD")
						}
						return nil
					}).
					Value(&endDateStr),
			),
		)
		
		err = customDateForm.Run()
		cobra.CheckErr(err)
		
		// Set custom dates on the command
		cmd.Flags().Set("startDate", startDateStr)
		cmd.Flags().Set("endDate", endDateStr)
	} else {
		// Set the selected date option
		cmd.Flags().Set(selectedDateOption, "true")
	}
	
	// Get the date range based on the selected option
	sd, ed = togglcmd.GetDateRange(cmd)
	
	// Step 3: Choose output options
	var saveTextFile, saveJsonFile bool
	
	outputForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Show text report?").
				Value(&outputText),
				
			huh.NewConfirm().
				Title("Show JSON report?").
				Value(&outputJson),
				
			huh.NewConfirm().
				Title("Save text report to file?").
				Value(&saveTextFile),
				
			huh.NewConfirm().
				Title("Save JSON report to file?").
				Value(&saveJsonFile),
		),
	)
	
	err = outputForm.Run()
	cobra.CheckErr(err)
	
	// Step 4: Get file paths if needed
	if saveTextFile {
		filePathForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Text file path").
					Placeholder("/tmp/time.txt").
					Value(&outputTextFile),
			),
		)
		err = filePathForm.Run()
		cobra.CheckErr(err)
	}
	
	if saveJsonFile {
		filePathForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("JSON file path").
					Placeholder("/tmp/time.json").
					Value(&outputJsonFile),
			),
		)
		err = filePathForm.Run()
		cobra.CheckErr(err)
	}
	
	// Step 5: Generate and display the report
	rpt, err := reportService.BuildReport(sd, ed)
	cobra.CheckErr(err)
	
	// Display date range info
	fmt.Printf("\nReport for period: %s to %s\n\n", 
		sd.Format("2006-01-02"), 
		ed.Format("2006-01-02"))
	
	// Process text output
	if outputText || outputTextFile != "" {
		r, err := reportService.BuildStringReport(rpt)
		cobra.CheckErr(err)
		
		if outputText {
			fmt.Println(*r)
		}
		
		if outputTextFile != "" {
			os.WriteFile(outputTextFile, []byte(*r), 0644)
			fmt.Printf("Text report saved to: %s\n", outputTextFile)
		}
	}
	
	// Process JSON output
	if outputJson || outputJsonFile != "" {
		jsonBytes, err := reportService.BuildJsonReport(rpt)
		cobra.CheckErr(err)
		
		if outputJson {
			fmt.Println(string(*jsonBytes))
		}
		
		if outputJsonFile != "" {
			os.WriteFile(outputJsonFile, *jsonBytes, 0644)
			fmt.Printf("JSON report saved to: %s\n", outputJsonFile)
		}
	}
}

func init() {
	cobra.OnInitialize(initServices)
	togglcmd.TogglCmd.AddCommand(reportCmd)

	reportCmd.Flags().BoolP("text", "t", false, "Produces a text report")
	reportCmd.Flags().BoolP("json", "j", false, "Produces a JSON report")

	reportCmd.Flags().StringVarP(&textOutput, "ot", "", "", "Writes a text report file.")
	reportCmd.Flags().Lookup("ot").NoOptDefVal = "/tmp/time.txt"

	reportCmd.Flags().StringVarP(&jsonOutput, "oj", "", "", "Writes a JSON report file.")
	reportCmd.Flags().Lookup("oj").NoOptDefVal = "/tmp/time.json"

	reportCmd.Flags().BoolP("today", "", false, "Run a report for today")
	reportCmd.Flags().BoolP("yesterday", "", false, "Run a report for yesterday")
	reportCmd.Flags().BoolP("eyesterday", "", false, "Run a report for eyesterday")
	reportCmd.Flags().BoolP("thisweek", "", false, "Run a report for this week")
	reportCmd.Flags().BoolP("lastweek", "", false, "Run a report for last week")

	reportCmd.Flags().StringP("startDate", "s", "", "Start date")
	reportCmd.Flags().StringP("endDate", "e", "", "End date")
	
	reportCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Run report in interactive mode")
}

func initServices() {
	reportService = initialiseReportService()
}
