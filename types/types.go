/*
Copyright © 2024 PACLabs
*/
package types

/*
 *  all of the test-related structs we need for this project
 */

import (
	"fmt"
	"raygun/opa"
	"time"
)

type TestSuite struct {
	Opa         opa.OpaConfig
	Name        string `yaml:"suite"`
	Description string
	Directory   string `yaml:"-"`
	Tests       []TestRecord
}

func (suite TestSuite) String() string {

	return fmt.Sprintf("Suite: %s with %d Tests and OPA config: %v", suite.Name, len(suite.Tests), suite.Opa.String())
}

type TestRecord struct {
	Suite        TestSuite `yaml:"-"`
	Name         string
	Description  string
	Expects      []TestExpectation
	Input        TestInput
	DecisionPath string `yaml:"decision-path"`
}

func (tr TestRecord) String() string {

	return fmt.Sprintf("Test: %s (%s)", tr.Name, tr.Description)
}

type CombinedResult struct {
	ResultList []TestSuiteResult
}

type TestSuiteResult struct {
	Source  TestSuite
	Failed  []TestResult
	Passed  []TestResult
	Skipped []TestResult
}

func (tsr TestSuiteResult) String() string {

	return fmt.Sprintf("Suite Results: %s - Passed: %d, Failed: %d, Skipped: %d",
		tsr.Source.Name, len(tsr.Passed), len(tsr.Failed), len(tsr.Skipped))

}

type TestResult struct {
	Source   TestRecord
	Actual   string
	Status   string // fail, pass, skip
	Start    time.Time
	End      time.Time
	Duration time.Duration
}

func (tr TestResult) String() string {

	return fmt.Sprintf("TestResult: %s - status: %s", tr.Source.Name, tr.Status)
}

type TestExpectation struct {
	ExpectationType string `yaml:"type"` // exact, substring, jsonpath
	Target          string
}

func (te TestExpectation) String() string {

	return fmt.Sprintf("TestExpectation: Type: %s  - Target: %s", te.ExpectationType, te.Target)

}

type TestInput struct {
	InputType string `yaml:"type"` // filepath, inline
	Value     string
}

func (ti TestInput) String() string {

	if ti.InputType == "filepath" {
		return fmt.Sprintf("TestInput File: %s", ti.Value)
	}

	return fmt.Sprintf("TestInput: %s...", ti.Value[:20])
}
