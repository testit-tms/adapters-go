# Test IT TMS Adapter for Golang


## Compatibility

| Test IT | Adapters-Go         |
|---------|---------------------|
| 5.2.5   | v0.3.5              |
| 5.3     | v0.3.5-tms-5.3      |
| 5.4     | v0.4.2-tms-5.4      |
| Cloud   | v0.5.1              |

For other versions compatibility check api-client compatibility - 
https://github.com/testit-tms/api-client-golang  
and previous version of adapter


## Getting Started

### Installation

```bash
go get github.com/testit-tms/adapters-go@<necessary package version>
```

## Usage

### Configuration

| Description                                                                                                                                                                                                                                                                                                                                                                            | File property                     | Environment variable                       |
|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------|--------------------------------------------|
| Location of the TMS instance                                                                                                                                                                                                                                                                                                                                                           | url                               | TMS_URL                                    |
| API secret key [How to getting API secret key?](https://github.com/testit-tms/.github/tree/main/configuration#privatetoken)                                                                                                                                                                                                                                                            | privateToken                      | TMS_PRIVATE_TOKEN                          |
| ID of project in TMS instance [How to getting project ID?](https://github.com/testit-tms/.github/tree/main/configuration#projectid)                                                                                                                                                                                                                                                    | projectId                         | TMS_PROJECT_ID                             |
| ID of configuration in TMS instance [How to getting configuration ID?](https://github.com/testit-tms/.github/tree/main/configuration#configurationid)                                                                                                                                                                                                                                  | configurationId                   | TMS_CONFIGURATION_ID                       |
| ID of the created test run in TMS instance. <br/>It's necessary for **adapterMode** 1                                                                                                                                                                                                                                                                                                                                           | testRunId                         | TMS_TEST_RUN_ID                            |
| Adapter mode. Default value - 1. The adapter supports following modes:<br>1 - in this mode, the adapter sends all results to the test run without filtering or [with filtering CLI](#run-with-filter)<br/>2 - in this mode, the adapter creates a new test run and sends results to the new test run | adapterMode                       | TMS_ADAPTER_MODE                           |
| It enables/disables certificate validation (**It's optional**). Default value - true                                                                                                                                                                                                                                                                                                   | certValidation                    | TMS_CERT_VALIDATION                        |
| Mode of automatic creation test cases (**It's optional**). Default value - false. The adapter supports following modes:<br/>true - in this mode, the adapter will create a test case linked to the created autotest (not to the updated autotest)<br/>false - in this mode, the adapter will not create a test case                                                                    | automaticCreationTestCases        | TMS_AUTOMATIC_CREATION_TEST_CASES          |
| Mode of automatic updation links to test cases (**It's optional**). Default value - false. The adapter supports following modes:<br/>true - in this mode, the adapter will update links to test cases<br/>false - in this mode, the adapter will not update link to test cases                                                                                                         | automaticUpdationLinksToTestCases | TMS_AUTOMATIC_UPDATION_LINKS_TO_TEST_CASES |
| Enable debug logs (**It's optional**). Default value - false                                                                                                                                                                                                                                                                                                                           | isDebug                           | TMS_IS_DEBUG                               |

#### File

Create **tms.config.json** file in the project directory:

```json
{
  "url": "URL",
  "privateToken": "USER_PRIVATE_TOKEN",
  "projectId": "PROJECT_ID",
  "configurationId": "CONFIGURATION_ID",
  "testRunId": "TEST_RUN_ID",
  "automaticCreationTestCases": false,
  "automaticUpdationLinksToTestCases": false,
  "certValidation": true,
  "adapterMode": "1",
  "isDebug": true
}
```

Alternatively to set TMS_CONFIG_FILE you can place your `tms.config.json` file 

to the folder with `_test.go` files you are want to work with, 

but for multifolder structure 
`TMS_CONFIG_FILE` is prefered.


### How to run

If you specified TestRunId, then just run the command:

```bash
export TMS_CONFIG_FILE=<ABSOLUTE_PATH_TO_CONFIG_FILE>
cd examples
go test
```

To create and complete TestRun you can use
the [Test IT CLI](https://docs.testit.software/user-guide/integrations/cli.html):

```bash
export TMS_TOKEN=<YOUR_TOKEN>
testit \
  testrun create
  --url https://tms.testit.software \
  --project-id 5236eb3f-7c05-46f9-a609-dc0278896464 \
  --testrun-name "New test run" \
  --output tmp/output.txt

export TMS_TEST_RUN_ID=$(cat output.txt)  

export TMS_CONFIG_FILE=<ABSOLUTE_PATH_TO_CONFIG_FILE>
cd examples
go test

testit \
  --mode finish
  --url https://tms.testit.software \
  --testrun-id $(cat tmp/output.txt) 
```

### Run with filter
To create filter by autotests you can use the Test IT CLI (use adapterMode "1" for run with filter):

```
$ export TMS_TOKEN=<YOUR_TOKEN>
$ testit autotests_filter 
  --url https://tms.testit.software \
  --configuration-id 5236eb3f-7c05-46f9-a609-dc0278896464 \
  --testrun-id 6d4ac4b7-dd67-4805-b879-18da0b89d4a8 \
  --framework golang \
  --output tmp/filter.txt

$ export TMS_TEST_RUN_ID=6d4ac4b7-dd67-4805-b879-18da0b89d4a8
$ export TMS_ADAPTER_MODE=1

$ export TMS_CONFIG_FILE=<ABSOLUTE_PATH_TO_CONFIG_FILE>
$ go test -run "$(cat tmp/filter.txt)"
```

### Asserting usage notes

* You should use `tms.True` as asserts for correct TestResult generation, e.g.:

```
expectedValue := 4
actualValue := 5 // function call you are testing
					
tms.True(t, actualValue == expectedValue) // its be an error and failed test, cause assert with false result.
```

* You can use `tms.True` both in tms.Test and tms.Step functions 


### Metadata of autotest

Use metadata to specify information about autotest.

Description of metadata:

* `WorkItemIds` - a method that links autotests with manual tests. Receives the array of manual tests' IDs
* `DisplayName` - internal autotest name (used in Test IT)
* `ExternalId` - unique internal autotest ID (used in Test IT)
* `Title` - autotest name specified in the autotest card. If not specified, the name from the displayName method is used
* `Description` - autotest description specified in the autotest card
* `Labels` - tags listed in the autotest card
* `Links` - links listed in the autotest card ( not in the TestResult card. Additionally, there is URL validation on Link.Url and it's must be a correct URL. )
* `Step` - the designation of the step

Description of methods:

* `tms.AddLinks` - add links to the autotest TestResult (not the autotest card, for card see `TestMetadata.Links`. Additionally, there is URL validation on Link.Url and it's must be a correct URL. ).
* `tms.AddAttachments` - add attachments to the autotest result.
* `tms.AddAtachmentsFromString` - add attachments from string to the autotest result.
* `tms.AddMessage` - add message to the autotest result.

### Examples

More examples and project there: https://github.com/testit-tms/go-examples 

#### Simple test

```go
package examples

import (
  "testing"

  "github.com/testit-tms/adapters-go"
)


func TestSteps_Success(t *testing.T) {

	// links at the autotest card
	links := []tms.Link{{
		Url:         "http://google.com",
		Title:       "Link title",
		Description: "Link description",
		LinkType:    "Requirement",
	}}

	labels := []string{"Test labels"}
	parameters := map[string]interface{}{
		"param1": "value1",
	}

	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "steps success",
			// Links for autotest card
			Links:      links,
			Labels:     labels,
			Parameters: parameters,
			// other properties...
		},
		func() {

			// for TestResult attachment
			// tms.AddAtachments("tms.config.json")

			// add links to TestResult
			tms.AddLinks(tms.Link{
				Url:         "https://testit.software",
				Title:       "Link title",
				Description: "Link description",
				LinkType:    tms.LINKTYPE_RELATED,
			})
			// add message to TestResult
			tms.AddMessage("Test Message")

			// step declaration
			tms.Step(
				tms.StepMetadata{
					Name:        "step 1",
					Description: "step 1 description",
				},
				func() {
					tms.Step(tms.StepMetadata{
						Name:        "step 1.1",
						Description: "step 1.1 description",
					}, func() {
						tms.Step(tms.StepMetadata{}, func() {
							tms.True(t, true)
						})
						tms.True(t, true)
					})
					tms.True(t, true)
				},
			)
			tms.Step(
				tms.StepMetadata{
					Name:        "step 2",
					Description: "step 2 description",
				},
				func() {
					tms.Step(tms.StepMetadata{
						Name:        "step 2.1",
						Description: "step 2.1 description",
					}, func() {
						tms.Step(tms.StepMetadata{}, func() {
							tms.True(t, true)
						})
						tms.True(t, true)
					})
					tms.True(t, true)
				},
			)
		})
}
```

#### Parameterized test

```go
package examples

import (
  "testing"

  "github.com/testit-tms/adapters-go"
)

func TestParameters_success(t *testing.T) {
	tests := []struct {
		name           string
		parameters     map[string]interface{}
		stepName       string
		stepParameters map[string]interface{}
		expValue       bool
	}{
		{
			name: "add parameters success",
			parameters: map[string]interface{}{
				"param1": "value1",
				"param2": 15,
			},
			stepName: "step1",
			stepParameters: map[string]interface{}{
				"param1": "value1",
				"param2": 15,
			},
			expValue: true,
		},
		{
			name: "add parameters failed",
			parameters: map[string]interface{}{
				"param1": "value1",
				"param2": 15,
			},
			stepName: "step1",
			stepParameters: map[string]interface{}{
				"param1": "value1",
				"param2": 15,
			},
			expValue: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tms.Test(t, tms.TestMetadata{
				DisplayName: tt.name,
				Parameters:  tt.parameters,
			}, func() {
				tms.Step(
					tms.StepMetadata{
						Name:       tt.stepName,
						Parameters: tt.stepParameters,
					}, func() {
						tms.True(t, tt.expValue)
					})
			})
		})
	}
}
```

## Contributing

You can help to develop the project. Any contributions are **greatly appreciated**.

* If you have suggestions for adding or removing projects, feel free
  to [open an issue](https://github.com/testit-tms/adapters-go/issues/new) to discuss it, or create a direct pull
  request after you edit the *README.md* file with necessary changes.
* Make sure to check your spelling and grammar.
* Create individual PR for each suggestion.
* Read the [Code Of Conduct](https://github.com/testit-tms/adapters-go/blob/main/CODE_OF_CONDUCT.md) before posting
  your first idea as well.

## License

Distributed under the Apache-2.0 License.
See [LICENSE](https://github.com/testit-tms/adapters-go/blob/main/LICENSE.md) for more information.
