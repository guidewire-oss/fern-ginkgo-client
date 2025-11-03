
## Introduction

Welcome to the Fern Project, an innovative open-source solution designed to enhance Ginkgo test reports. This project is focused on capturing, storing, and analyzing test data to provide insights into test performance and trends. The Fern Project is ideal for teams using Ginkgo, a popular BDD-style Go testing framework, offering a comprehensive overview of test executions and performance metrics.

### Integrating the Client into Ginkgo Test Suites

1. **Add the Fern dependency to your test project**:

   ```bash
   go get -u github.com/guidewire-oss/fern-ginkgo-client
   ```
2. **Generate Project ID by sending the below payload to `fern-reporter` (hosted in your environment)** 
```bash
curl -L -X POST http://localhost:8080/api/project \
  -H "Content-Type: application/json" \
  -d '{
    "name": "First Projects",
    "team_name": "my team",
    "comment": "This is the test project"
  }' 
```
Sample Response:
```json
{
  "uuid": "996ad860-2a9a-504f-8861-aeafd0b2ae29",
  "name": "First Projects",
  "team_name": "my team",
  "comment": "This is the test project"
}
```
3. **Add the Fern Client to your Ginkgo test suite**:
   
   In the example below, from this project, the constant pkg.ProjectName is set to the name of this project:
   ```go
   const PROJECT_ID = "996ad860-2a9a-504f-8861-aeafd0b2ae29"
   ```
   Import the fern client package into the Ginkgo test suite file:
   ```go
   import fern "github.com/guidewire-oss/fern-ginkgo-client/pkg/client"
   ```
   Add ReportAfterSuite to call the Fern ReportTestResult.    Initialize the fernClient by passing the Project ID and ClientOption. Invoke the `Report` function by passing the report Object.

   ```go
   var _ = ReportAfterSuite("", func(report Report) {
      fernReporterBaseUrl := "http://localhost:8080/"
       if os.Getenv("FERN_REPORTER_BASE_URL") != "" {
           fernReporterBaseUrl = os.Getenv("FERN_REPORTER_BASE_URL")
       }
       fernApiClient := fern.New(pkg.PROJECT_ID, fern.WithBaseURL(fernReporterBaseUrl))
       fernApiClient.Report(report)
   })

   ```
3. **Run Your Tests**: After adding the client, run your Ginkgo tests normally.

   ```bash
   ginkgo -r -p --label-filter=unit --randomize-all
   ```
4. **Authentication**: when using the Fern Platform (which has authentication enabled), set the following environment variables:
   ```shell
    export FERN_AUTH_CLIENT_ID=<Your Service Application Client ID>
    export FERN_AUTH_CLIENT_SECRET=<Your Service Application Client Secret>
    export AUTH_URL=<Base URL of your authentication server>
    export FERN_GINKGO_CLIENT_SCOPE=<Fern Platform scope for Testrun write>
   ```
5. **Test Level Tags**: if you want to set Test Run level tags, you can do this using an environment variable
TEST_RUN_TAGS which can be set to one or more values, separated by commas. For example:
    ```shell
    export TEST_RUN_TAGS="Tag1,Tag2"
    ```
6. **Environment**: if you want to set the Environment, which is stored in the Test Run in Fern, you can do this using
an environment variable TEST_ENVIRONMENT. For example:
    ```shell
    export TEST_ENVIRONMENT="local"
    ```

### See Also
1. [Fern UI](https://github.com/Guidewire/fern-ui)
2. [Fern Ginkgo Reporter](https://github.com/Guidewire/fern-reporter)