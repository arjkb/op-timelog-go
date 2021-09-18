package main

import (
	"testing"
)

// TestExtractData calls extractData with a correct input,
// checking for valid return values.
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

// TestExtractDataHourDuration calls extractData with a duration that ends in
// .00, checking for valid return values.
func TestExtractDataHourDuration(t *testing.T) {
	inputLine := "123 4.00 Meeting with Carrie, Kathy & John"
	expectedWorkpackage := 123
	expectedDuration := "4.00"
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

// TestExtractDataHalfHourDuration calls extractData with a duration that ends
// in .50, checking for valid return values.
func TestExtractDataHalfHourDuration(t *testing.T) {
	inputLine := "123 2.50 Meeting with Carrie, Kathy & John"
	expectedWorkpackage := 123
	expectedDuration := "2.50"
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

// TestExtractDataWithoutWP calls extractData without a work package
// checking for an error.
func TestExtractDataWithoutWP(t *testing.T) {
	inputLine := "1.25 Meeting with Carrie, Kathy & John"
	wp, dur, desc, err := extractData(inputLine)
	if wp != 0 || dur != "" || desc != "" || err == nil {
		t.Fatalf(`extractData(%q), got (%v, %v, %v, %v), expected err == nil`, inputLine, wp, dur, desc, err)
	}
}

// TestExtractDataWithoutDuration calls extractData without a duration
// checking for an error.
func TestExtractDataWithoutDuration(t *testing.T) {
	inputLine := "123 Meeting with Carrie, Kathy & John"
	wp, dur, desc, err := extractData(inputLine)
	if wp != 0 || dur != "" || desc != "" || err == nil {
		t.Fatalf(`extractData(%q), got (%v, %v, %v, %v), expected err == nil`, inputLine, wp, dur, desc, err)
	}
}

// TestExtractDataWithoutDescription calls extractData without a description
// checking for an error.
func TestExtractDataWithoutDescription(t *testing.T) {
	inputLine := "123 2.50 Meeting with Carrie, Kathy & John"
	wp, dur, desc, err := extractData(inputLine)
	if wp != 0 || dur != "" || desc != "" || err == nil {
		t.Fatalf(`extractData(%q), got (%v, %v, %v, %v), expected err == nil`, inputLine, wp, dur, desc, err)
	}
}

// TestExtractDataSingleWord calls extractData with just a single word,
// checking for an error.
func TestExtractDataSingleWord(t *testing.T) {
	inputLine := "foo"
	wp, dur, desc, err := extractData(inputLine)
	if wp != 0 || dur != "" || desc != "" || err == nil {
		t.Fatalf(`extractData(%q), got (%v, %v, %v, %v), expected err == nil`, inputLine, wp, dur, desc, err)
	}
}

// TestExtractDataBlank calls extractData with an empty string,
// checking for an error.
func TestExtractDataBlank(t *testing.T) {
	inputLine := ""
	wp, dur, desc, err := extractData(inputLine)
	if wp != 0 || dur != "" || desc != "" || err == nil {
		t.Fatalf(`extractData(%q), got (%v, %v, %v, %v), expected err == nil`, inputLine, wp, dur, desc, err)
	}
}
