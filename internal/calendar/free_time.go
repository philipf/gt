package calendar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func FreeTime() {
	// Read the access token from the environment variable ACCESS_TOKEN
	var accessToken = os.Getenv("ACCESS_TOKEN")

	if accessToken == "" {
		panic("ACCESS_TOKEN environment variable not set")
	}

	//GetMe(accessToken)
	getFindMeetingTimes(accessToken)
}

func getFindMeetingTimes(accessToken string) {

	location, err := time.LoadLocation("Pacific/Auckland")
	if err != nil {
		fmt.Println("Error loading time zone:", err)
		return
	}
	currentDate := time.Now().In(location)

	// add 1 day to the current date
	// currentDate = currentDate.AddDate(0, 0, 1)

	startTime := currentDate.Format("2006-01-02") + "T09:00:00"
	endTime := currentDate.Format("2006-01-02") + "T17:00:00"

	findMeetingTimesRequest := FindMeetingTimesRequest{
		Attendees: []Attendee{
			// {
			// 	Type: "optional",
			// 	EmailAddress: EmailAddress{
			// 		Address: "test@test.com",
			// 	},
			// },
		},
		TimeConstraint: TimeConstraint{
			Timeslots: []TimeSlot{
				{
					Start: DateTimeTimeZone{
						DateTime: startTime,
						TimeZone: "New Zealand Standard Time",
					},
					End: DateTimeTimeZone{
						DateTime: endTime,
						TimeZone: "New Zealand Standard Time",
					},
				},
			},
		},

		IsOrganizerOptional: false,
		//MeetingDuration:           "PT1H",
		//Meeting duration should be at least 30 minutes
		MeetingDuration:           "PT30M",
		ReturnSuggestionReasons:   false,
		MinimumAttendeePercentage: 50,
		MaxCandidates:             10,
	}

	findMeetingTimesRequestBody, err := json.Marshal(findMeetingTimesRequest)

	fmt.Println(">>>")
	fmt.Println(string(findMeetingTimesRequestBody))
	fmt.Println("<<<")

	if err != nil {
		panic(err)
	}

	fmt.Println(string(findMeetingTimesRequestBody))

	req, err := http.NewRequest("POST", "https://graph.microsoft.com/v1.0/me/findMeetingTimes", bytes.NewBuffer(findMeetingTimesRequestBody))

	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "outlook.timezone=\"New Zealand Standard Time\"")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

	var findMeetingTimesResponse FindMeetingTimesResponse

	err = json.Unmarshal(body, &findMeetingTimesResponse)

	if err != nil {
		panic(err)
	}

	fmt.Println("\nMeeting Time Suggestions:")

	fmt.Printf("%v to %v\n", startTime, endTime)

	// loop through all the meeting time suggestions and print the start and end times
	for i, meetingTimeSuggestion := range findMeetingTimesResponse.MeetingTimeSuggestions {
		// parse the start and end times
		startTime, _ := time.Parse("2006-01-02T15:04:00.0000000", meetingTimeSuggestion.MeetingTimeSlot.Start.DateTime)
		endTime, _ := time.Parse("2006-01-02T15:04:00.0000000", meetingTimeSuggestion.MeetingTimeSlot.End.DateTime)

		fmt.Printf(" %d] %v - %v \n", i+1, startTime.Format(time.Kitchen), endTime.Format(time.Kitchen))
	}

	fmt.Println()

	for i, meetingTimeSuggestion := range findMeetingTimesResponse.MeetingTimeSuggestions {
		// parse the start and end times
		startTime, _ := time.Parse("2006-01-02T15:04:00.0000000", meetingTimeSuggestion.MeetingTimeSlot.Start.DateTime)
		endTime, _ := time.Parse("2006-01-02T15:04:00.0000000", meetingTimeSuggestion.MeetingTimeSlot.End.DateTime)

		fmt.Printf(" %d] %v - %v \n", i+1, startTime.Format("15:04"), endTime.Format("15:04"))
	}
}

type MeResponse struct {
	DisplayName       string   `json:"displayName"`
	Mail              string   `json:"mail"`
	UserPrincipalName string   `json:"userPrincipalName"`
	MobilePhone       string   `json:"mobilePhone"`
	BusinessPhones    []string `json:"businessPhones"`
}

type FindMeetingTimesRequest struct {
	Attendees      []Attendee     `json:"attendees"`
	TimeConstraint TimeConstraint `json:"timeConstraint"`
	//LocationConstraint               LocationConstraint `json:"locationConstraint"`
	IsOrganizerOptional bool `json:"isOrganizerOptional"`
	//SuggestedTimeSuggestionExpansion Expansion          `json:"suggestedTimeSuggestionExpansion"`
	MaxCandidates             int    `json:"maxCandidates"`
	MeetingDuration           string `json:"meetingDuration"`
	ReturnSuggestionReasons   bool   `json:"returnSuggestionReasons"`
	MinimumAttendeePercentage int    `json:"minimumAttendeePercentage"`
}

type Attendee struct {
	Type         string       `json:"type"`
	EmailAddress EmailAddress `json:"emailAddress"`
}

type EmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type TimeConstraint struct {
	Timeslots []TimeSlot `json:"timeslots"`
}

type TimeSlot struct {
	Start DateTimeTimeZone `json:"start"`
	End   DateTimeTimeZone `json:"end"`
}

type DateTimeTimeZone struct {
	DateTime string `json:"dateTime"`
	TimeZone string `json:"timeZone"`
}

type Expansion struct {
	NumberOfTimeSuggestions              int               `json:"numberOfTimeSuggestions"`
	MaximumNonWorkHoursSuggestionsPerDay int               `json:"maximumNonWorkHoursSuggestionsPerDay"`
	MaximumSuggestionsPerDay             int               `json:"maximumSuggestionsPerDay"`
	MinimumSuggestionQuality             SuggestionQuality `json:"minimumSuggestionQuality"`
}

type SuggestionQuality struct {
	AttendeeConflict string `json:"attendeeConflict"`
	Unknown          string `json:"unknown"`
}

type FindMeetingTimesResponse struct {
	MeetingTimeSuggestions []MeetingTimeSuggestion `json:"meetingTimeSuggestions"`
	EmptySuggestionsReason string                  `json:"emptySuggestionsReason"`
}

type MeetingTimeSuggestion struct {
	Confidence            float32                `json:"confidence"`
	OrganizerAvailability string                 `json:"organizerAvailability"`
	SuggestionReason      string                 `json:"suggestionReason"`
	AttendeeAvailability  []AttendeeAvailability `json:"attendeeAvailability"`
	MeetingTimeSlot       TimeSlot               `json:"meetingTimeSlot"`
}

type AttendeeAvailability struct {
	Attendee     Attendee `json:"attendee"`
	Availability string   `json:"availability"`
}
