package runner

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"raygun/config"
	"raygun/log"
	"raygun/types"
	"raygun/util"
	"strings"
)

type TestRunner struct {
	Source types.TestRecord
}

func NewTestRunner(test types.TestRecord) TestRunner {
	testRunner := TestRunner{Source: test}

	return testRunner
}

func (tr TestRunner) Post() (string, error) {

	postUrl := fmt.Sprintf("http://localhost:%d%s", config.OpaPort, tr.Source.DecisionPath)

	bodyString := ""

	switch tr.Source.Input.InputType {
	case "inline":
		// we're readying directly from the .raygun file, and it may or may not already have
		// the input key
		bodyString = optionally_add_input_key(tr.Source.Input.Value)
	case "json-file":

		// we're reading the data from a json file, which may or may not have an input tag already
		log.Debug("Suite Directory: %s , filename: %s", tr.Source.Suite.Directory, tr.Source.Input.Value)
		tmp, err := util.ReadFile(tr.Source.Suite.Directory, tr.Source.Input.Value)
		if err != nil {
			return "", err
		}

		bodyString = optionally_add_input_key(tmp)

	default:
		return "", fmt.Errorf("unsupported input type: %s", tr.Source.Input.InputType)
	}

	//	log.Debug("BodyString: %s", bodyString)

	bodyBytes := []byte(bodyString)

	response, err := http.Post(postUrl, "application/json", bytes.NewReader(bodyBytes))

	if err != nil {
		log.Error("Attempted to complete POST to %s with payload %s -> %s", postUrl, bodyString, err.Error())
		return "", err
	}

	defer response.Body.Close()

	builderBuffer := new(strings.Builder)

	_, err = io.Copy(builderBuffer, response.Body)

	if err != nil {
		log.Error("Error reading body of response: %s", err.Error())
		return "", err
	}

	log.Debug("Response Content: %s", builderBuffer.String())

	return builderBuffer.String(), nil
}

func (tr TestRunner) Evaluate(response string) (types.TestResult, error) {

	result := types.TestResult{}

	result.Source = tr.Source

	switch tr.Source.Expects.ExpectationType {
	case "substring":
		compressed_actual := util.RemoveAllWhitespace(response)
		result.Actual = response

		expected := util.RemoveAllWhitespace(tr.Source.Expects.Target)
		if strings.Contains(compressed_actual, expected) {
			result.Status = config.PASS
		} else {
			result.Status = config.FAIL
		}

	default:
		log.Fatal("Unsupported ExpectationType for %s -> %s", tr.Source, tr.Source.Expects.ExpectationType)
	}

	return result, nil
}

/*
 *  Keeping it really simple until we know we need something more sophisticated
 */
func optionally_add_input_key(json string) string {

	no_whitespace := util.RemoveAllWhitespace(json)

	if strings.HasPrefix(no_whitespace, "{\"input\"") {
		return json
	}

	return fmt.Sprintf("{\"input\":%s}", json)

}
