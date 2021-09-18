package main

import (
	"testing"
)

func TestExtractData(t *testing.T) {
	inputLine := "123 4.56 Meeting with Carrie, Kathy & John"

	expectedWorkpackage := 123
	expectedDuration := "4.56"
	expectedDescription := "Meeting with Carrie, Kathy & John"

	actualWorkpackage, actualDuration, actualDescription, err := extractData(inputLine)

	if err != nil {
		t.Fatalf(`extractData(%q) want (%v, %q, %q), got err = %v`, inputLine, expectedWorkpackage, expectedDuration, expectedDescription, err)
	}

	if actualWorkpackage != expectedWorkpackage {
		t.Fatalf(`extractData(%q) workpackage, expected = %v, got = %v`, inputLine, expectedWorkpackage, actualWorkpackage)
	}

	if actualDuration != expectedDuration {
		t.Fatalf(`extractData(%q) duration, expected = %q, got = %q`, inputLine, expectedDuration, actualDuration)
	}

	if actualDescription != expectedDescription {
		t.Fatalf(`extractData(%q) description, expected = %q, got = %q`, inputLine, expectedDescription, actualDescription)
	}
}

func TestExtractDataWithIncorrectFormat(t *testing.T) {
	inputLine := "123 Meeting with Carrie, Kathy & John"

	wp, dur, desc, err := extractData(inputLine)

	if wp != 0 || dur != "" || desc != "" || err == nil {
		t.Fatalf(`extractData(%q), got (%v, %v, %v, %v), expected err == nil`, inputLine, wp, dur, desc, err)
	}
}
