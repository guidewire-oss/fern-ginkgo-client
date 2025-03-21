
## Introduction

Welcome to the Fern Project, an innovative open-source solution designed to enhance Ginkgo test reports. This project is focused on capturing, storing, and analyzing test data to provide insights into test performance and trends. The Fern Project is ideal for teams using Ginkgo, a popular BDD-style Go testing framework, offering a comprehensive overview of test executions and performance metrics.

### Integrating the Client into Ginkgo Test Suites

1. **Add the Fern dependency to your test project**:

   ```bash
   go get -u github.com/guidewire-oss/fern-ginkgo-client
   ```
2. **Add the Fern Client to your Ginkgo test suite**:
   
   In the example below, from this project, the constant pkg.ProjectName is set to the name of this project:
   ```go
   const ProjectName = "Fern Ginkgo Client"
   ```
   Import the fern client package into the Ginkgo test suite file:
   ```go
   import fern "github.com/guidewire-oss/fern-ginkgo-client/pkg/client"
   ```
   Add ReportAfterSuite to call the Fern ReportTestResult.
   ```go
   var _ = ReportAfterSuite("", func(report Report) {
      fern.ReportTestResult(pkg.ProjectName, report)
   })
   ```
   ReportTestResult takes two parameters, the name of the project and the report. It uses an environment variable
   `FERN_REPORTER_BASE_URL` to determine the location of the Fern Reporter service. If it is not set, then it will 
   default to `http://localhost:8080/`.

3. **Run Your Tests**: After adding the client, run your Ginkgo tests normally.

   How to execute the tests  :
   ```
   make test
   ```

   To add flags to test suites, add a label to the test suite file. For an example, look in 
   [tests/adder_suite_test.go](tests/adder_suite_test.go)
   ```go
   func TestAdder(t *testing.T) {
       RegisterFailHandler(Fail)
       RunSpecs(t, "Adder Suite", Label("this-is-a-suite-level-label"))
   }
   
   var _ = ReportAfterSuite("", func(report Report) {
       fern.ReportTestResult(pkg.ProjectName, report)
   })
   ```
   The fern report will have three test reports with the name `Fern Ginkgo Client`, three suite runs, and a spec run for 
   each test. If there is a label on the suite, then that gets stored and associated with each spec run that was within
   the suite. 

### See Also
1. [Fern UI](https://github.com/Guidewire/fern-ui)
2. [Fern Ginkgo Reporter](https://github.com/Guidewire/fern-reporter)