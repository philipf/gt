package togglservices

import (
	"errors"
	"testing"
)

func TestParseProjectTitle(t *testing.T) {
	tests := []struct {
		input  string
		output ProjectTitle
		err    error
	}{
		{
			input: "",
			err:   errors.New("project cannot be empty"),
		},
		{
			input: "[ClientName|P|12345|T12345678.1] Description here",
			output: ProjectTitle{
				Client:   "ClientName",
				IsTask:   false,
				TaskID:   12345,
				TicketID: "T12345678.1",
				Project:  "Description here",
			},
		},
		{
			input: "[ClientName|S|12345] Description here",
			output: ProjectTitle{
				Client:  "ClientName",
				IsTask:  true,
				TaskID:  12345,
				Project: "Description here",
			},
		},
		{
			input: "InvalidProjectString",
			err:   errors.New("project does not match the naming convention"),
		},
		{
			input: "[ClientName|S|invalidID] Description here",
			err:   errors.New("project does not match the naming convention"),
		},
	}

	for _, tt := range tests {
		result, err := ParseProjectTitle(tt.input)
		if err != nil && tt.err == nil {
			t.Errorf("expected no error but got %v", err)
		} else if err == nil && tt.err != nil {
			t.Errorf("expected error %v but got none", tt.err)
		} else if err != nil && err.Error() != tt.err.Error() {
			t.Errorf("expected error %v but got %v", tt.err, err)
		}

		if result != tt.output {
			t.Errorf("expected output %v but got %v", tt.output, result)
		}
	}
}
