package report

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/ethanefung/bubble-datepicker"
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
		// Use bubble-datepicker for date selection
		startDate, err := getDateWithPicker("Select Start Date")
		cobra.CheckErr(err)
		
		endDate, err := getDateWithPicker("Select End Date")
		cobra.CheckErr(err)
		
		// Format dates for command flags
		startDateStr = startDate.Format("2006/01/02")
		endDateStr = endDate.Format("2006/01/02")
		
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

// Bubble Tea model for datepicker
type datePickerModel struct {
	title      string
	datePicker datepicker.Model
	quitting   bool
	err        error
	selectedDate time.Time
}

func initialModel(title string) datePickerModel {
	dp := datepicker.New(time.Now())
	dp.SelectDate() // Activate date selection
	
	return datePickerModel{
		title:      title,
		datePicker: dp,
		selectedDate: time.Time{},
	}
}

func (m datePickerModel) Init() tea.Cmd {
	return nil
}

func (m datePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			// Store selected date and quit
			m.selectedDate = m.datePicker.Time
			return m, tea.Quit
		}
	}
	
	// Handle datepicker updates
	dp, cmd := m.datePicker.Update(msg)
	m.datePicker = dp
	
	return m, cmd
}

func (m datePickerModel) View() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("%s\n\n", m.title))
	s.WriteString(m.datePicker.View())
	s.WriteString("\n\nPress Enter to confirm selection, Esc to cancel\n")
	
	return s.String()
}

// getDateWithPicker launches a BubbleTea program with a datepicker
// and returns the selected date
func getDateWithPicker(title string) (time.Time, error) {
	m := initialModel(title)
	p := tea.NewProgram(m)
	
	finalModel, err := p.Run()
	if err != nil {
		return time.Time{}, err
	}
	
	// Get the final state
	finalState := finalModel.(datePickerModel)
	
	if finalState.quitting && finalState.selectedDate.IsZero() {
		return time.Now(), nil // Default to today if user quit
	}
	
	return finalState.selectedDate, nil
}
