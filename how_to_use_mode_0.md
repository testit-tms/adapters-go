1. export TMS_TOKEN=<YOUR_TOKEN>

2. testit testrun create --url https://team-s01g.testit.software --project-id 01963990-832c-72ce-b082-eec991e243af  --testrun-name "New test run"  --output output.txt

3. export TMS_TEST_RUN_ID=$(cat output.txt)  

# you can use pwd to get absolute path in bash
4. pwd
# Set TMS_CONFIG_FILE as absulute path to config
5. export TMS_CONFIG_FILE=/c/..../adapters-go/tms.config.json

6. setup config (testRunId will be overrided with env's TMS_TEST_RUN_ID)
{
    "url": "https://team-s01g.testit.software",
    "privateToken": "bURoMUpKNk4wYjJVeWZ2NTl0",
    "projectId": "01963990-832c-72ce-b082-eec991e243af",
    "configurationId": "01963990-835b-763a-9ae5-f159a73c6dc2",
    "testRunId": "7694ebc7-c34b-4c44-a804-24015cfdd4e5",
    "adapterMode": "0"
}

7. cd examples
8. go test ./...

9. Check your test run in Test IT.

