package quiz

// Naming conventions:
//
// file names
// naming the file xxx_test.go is important, this allows Go to identify the file as a test file
// 'go build' will also ignore these files
// there's usually a test file for each source file, just good practice
// caveats:
// - 'export_test.go' to access unexported variables in external tests
// - 'xxx_internal_test.go' for internal tests
// - 'example_xxx_test.go' for examples in isolated files

// function names
// starts with TestXxx where Xxx is the name of the function/type being tested
// if testing method of type, TestXxx_Xxx
// finally, add _xxx for test case. TestXxx_Xxx_xxx or TestXxx_xxx
// unless you use table driven tests
// if you're using examples as tests, name them ExampleXxx

// variable names
// 'want' and 'got'
// if got != want
// t.Errorf("function, got, want")

// testing Signals
// t.Log and t.Logf only show when a test fails or you run 'go test -v'
// t.Fail test fails but continues running
// t.FailNow test fails and stops running
// Error = Log + Fail
// Fatal = Log + FailNow

// Examples
// the "Output:" comment tells Go what output to expect
// use "Unordered output:" for unordered outputs

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func areEqual(a, b [][]string) bool {
	// Check if lengths are different
	if len(a) != len(b) {
		return false
	}

	// Compare each row
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}

	return true
}

func ExampleRetrieveProblems() {
	file := "problems_test.csv"
	gotProblems, err := RetrieveProblems(file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", gotProblems)
	// Output:
	// [[5+3 8]]
}

func deepCopy(a [][]string) [][]string {
	b := make([][]string, len(a))
	for i := range a {
		copy(b[i], a[i])
	}
	return b
}

func TestDeepCopy(t *testing.T) {
	problems := [][]string{{"5+8", "13"}, {"6*3", "18"}}
	want := deepCopy(problems)
	problems[0] = []string{"6-4", "2"}
	if areEqual(want, problems) {
		t.Errorf("Did not deep copy, got %v, want %v", problems, want)
	}
}

func ExampleCreateQuiz() {
	problems := [][]string{{"5+8", "13"}, {"6*3", "18"}}
	in, err := os.CreateTemp("", "example")
	if err != nil {
		panic(err)
	}
	defer os.Remove(in.Name())

	_, err = in.WriteString("13\n" + "18\n")
	if err != nil {
		panic(err)
	}

	_, err = in.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}

	CreateQuiz(in, problems, false, 0)
	// Output:
	// Problem #0: 5+8=Problem #1: 6*3=
	// Complete! Elapsed time: 0.00 seconds
	// You scored 2 out of 2
}
