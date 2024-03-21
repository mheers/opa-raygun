/*
Copyright © 2024 PACLabs
*/
package report

/*
 *  Generate a simple text output that looks familiar to
 */

import (
	"fmt"
	"raygun/config"
	"raygun/types"
	"strings"
)

type TextReporter struct {
	BaseReporter
}

func (tr TextReporter) Generate(results types.CombinedResult) string {

	var sb strings.Builder

	sb.WriteString("Test Results:\n")

	failureCount := 0

	for _, suite_result := range results.ResultList {

		sb.WriteString(fmt.Sprintf("   Suite: %s :\n", suite_result.Source.Name))

		if config.Verbose {
			sb.WriteString("      OPA Configuration:\n")
			sb.WriteString(fmt.Sprintf("         OPA Output Log: %s\n", suite_result.Source.Opa.LogPath))
			sb.WriteString(fmt.Sprintf("         Using OPA Bundle: %s\n", suite_result.Source.Opa.BundlePath))
		}

		for _, test_result := range suite_result.Skipped {
			sb.WriteString(fmt.Sprintf("      SKIPPED: %s\n", test_result.Source.Name))
			if config.Verbose {
				sb.WriteString(fmt.Sprintf("        - %s\n", test_result.Source.Description))
			}
		}
		for _, test_result := range suite_result.Passed {
			sb.WriteString(fmt.Sprintf("      PASSED: %s\n", test_result.Source.Name))
			if config.Verbose {
				sb.WriteString(fmt.Sprintf("        - %s\n", test_result.Source.Description))
			}
		}

		for _, test_result := range suite_result.Failed {
			failureCount++
			sb.WriteString("\n")
			sb.WriteString(fmt.Sprintf("      FAILED: %s\n", test_result.Source.Name))
			if config.Verbose {
				sb.WriteString(fmt.Sprintf("        - %s\n", test_result.Source.Description))
			}

			if config.Verbose {

				// the TrimRight at the end is to make sure we don't have a dangling ] on a single line
				sb.WriteString(fmt.Sprintf("        Comparison: %s. Expected:[%s] Actual: [%s]\n",
					test_result.Source.Expects.ExpectationType,
					test_result.Source.Expects.Target,
					strings.TrimRight(test_result.Actual, "\r\n")))

				if test_result.Source.Input.InputType == "json-file" {
					sb.WriteString(fmt.Sprintf("        Input File: %s\n", test_result.Source.Input.Value))
				}

			}
			sb.WriteString("\n")
		}

	}

	if failureCount > 0 {
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("WARNING: There are test failures: %d\n", failureCount))
	}

	return sb.String()
}
