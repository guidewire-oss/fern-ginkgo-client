
## Introduction

Welcome to the Fern Project, an innovative open-source solution designed to enhance Ginkgo test reports. This project is focused on capturing, storing, and analyzing test data to provide insights into test performance and trends. The Fern Project is ideal for teams using Ginkgo, a popular BDD-style Go testing framework, offering a comprehensive overview of test executions and performance metrics.

### Integrating the Client into Ginkgo Test Suites

1. **Add the Fern dependency to your test project**:

   ```bash
   go get -u github.com/guidewire/fern-ginkgo-client
   ```
2. **Add the Fern Client to your Ginkgo test suite**:
   
   Import the fern client package
   ```go
   import fern "github.com/guidewire/fern-ginkgo-client/pkg/client"
   ```

   ```go
   var _ = ReportAfterSuite("", func(report Report) {
       f := fern.New("Example Test",
           fern.WithBaseURL("http://localhost:8080/"),
       )

       err := f.Report("example test", report)

       Expect(err).To(BeNil(), "Unable to create reporter file")
   })
   ```
   Replace `http://localhost:8080/` with your API server's URL and specify the project name in `f.Report`.

2. **Run Your Tests**: After adding the client, run your Ginkgo tests normally.

    How to execute the tests  :
    ```
    make test
    ```

    To add flags to test suits : The flag is added to file : 
    tests/adder_suite_test.go
  ```
    Example : 
    RunSpecs(t, "Adder Suite", Label("this-is-a-suite-level-label","anotherFlag"))
  ```
    The fern report will have the test report along with the labels.

**See Also**
fern-ui: `https://github.com/Guidewire/fern-ui`
fern-reporter: `https://github.com/Guidewire/fern-reporter`