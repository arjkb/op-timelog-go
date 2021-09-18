package main

import (
	"testing"
)

func TestExtractData(t *testing.T) {
	inputLine := "123 4.56 Meeting with Carrie, Kathy & John"

	expectedWorkpackage := 123
	expectedDuration := "4.56"
	expectedDescription := "Meeting with Carrie, Kathy & John"

	actualWorkpackage, actualDuration, actualDescription := extractData(inputLine)

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
