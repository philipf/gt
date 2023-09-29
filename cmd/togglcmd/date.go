package togglcmd

import (
	"time"

	"github.com/spf13/cobra"
)

func GetDateRange(cmd *cobra.Command) (time.Time, time.Time) {
	// First check if any of the fixed periods have been selected (today, yesterday, etc)
	today, err := cmd.Flags().GetBool("today")
	cobra.CheckErr(err)

	now := time.Now()

	if today {
		// Set the start date to today, trimming the time
		sd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		// Set the end date to tomorrow, trimming the time and subtracting a second
		ed := sd.AddDate(0, 0, 1).Add(-time.Second)
		return sd, ed
	}

	// check if yesterday has been selected
	yesterday, err := cmd.Flags().GetBool("yesterday")
	cobra.CheckErr(err)

	if yesterday {
		// Set the start date to yesterday, trimming the time
		now := now.AddDate(0, 0, -1)
		sd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		// Set the end date to today, trimming the time and subtracting a second
		ed := sd.AddDate(0, 0, 1).Add(-time.Second)
		return sd, ed
	}

	// check if eyesterday has been selected
	eyesterday, err := cmd.Flags().GetBool("eyesterday")
	cobra.CheckErr(err)

	if eyesterday {
		// Set the start date to yesterday, trimming the time
		now := now.AddDate(0, 0, -2)
		sd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		// Set the end date to yesterday, trimming the time and subtracting a second
		ed := sd.AddDate(0, 0, 1).Add(-time.Second)
		return sd, ed
	}

	// check if thisweek has been selected
	thisweek, err := cmd.Flags().GetBool("thisweek")
	cobra.CheckErr(err)

	if thisweek {
		// TODO: get logic from the gt-at package
		thisSunday := sundayOfTheWeek(time.Now())
		sd := time.Date(thisSunday.Year(), thisSunday.Month(), thisSunday.Day(), 0, 0, 0, 0, thisSunday.Location())
		ed := sd.AddDate(0, 0, 7).Add(-time.Second)

		return sd, ed
	}

	// check if lastweek has been selected
	lastweek, err := cmd.Flags().GetBool("lastweek")
	cobra.CheckErr(err)

	if lastweek {
		thisSunday := sundayOfTheWeek(time.Now())
		lastSunday := thisSunday.AddDate(0, 0, -7)
		sd := time.Date(lastSunday.Year(), lastSunday.Month(), lastSunday.Day(), 0, 0, 0, 0, lastSunday.Location())
		ed := sd.AddDate(0, 0, 7).Add(-time.Second)

		return sd, ed
	}

	sdStr, err := cmd.Flags().GetString("startDate")
	cobra.CheckErr(err)

	edStr, err := cmd.Flags().GetString("endDate")
	cobra.CheckErr(err)

	loc, err := time.LoadLocation("Local")
	cobra.CheckErr(err)

	sd, err := time.ParseInLocation("2006/01/02", sdStr, loc)
	cobra.CheckErr(err)

	ed, err := time.ParseInLocation("2006/01/02", edStr, loc)
	cobra.CheckErr(err)

	ed = ed.AddDate(0, 0, 1).Add(-time.Second)

	return sd, ed
}

// SundayOfTheWeek returns the date of the Sunday of the week based on a provided date.
func sundayOfTheWeek(t time.Time) time.Time {
	// Subtract the weekday number from the given date.
	// Since Sunday = 0 in time.Weekday, it gives the exact offset we need.
	offset := int(t.Weekday())
	return t.AddDate(0, 0, -offset)
}
