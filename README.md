[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/Guidewire/fern-reporter/badge)](https://securityscorecards.dev/viewer/?uri=github.com/Guidewire/fern-reporter)
![Coverage](https://img.shields.io/badge/Coverage-23.0%25-red)

![Fern](https://github.com/guidewire/fern-reporter/raw/main/docs/images/logo-no-background.png)


## Introduction

Welcome to the Fern Project, an innovative open-source solution designed to enhance Ginkgo test reports. This project is focused on capturing, storing, and analyzing test data to provide insights into test performance and trends. The Fern Project is ideal for teams using Ginkgo, a popular BDD-style Go testing framework, offering a comprehensive overview of test executions and performance metrics.

### Integrating the Client into Ginkgo Test Suites

1. **Add the Fern dependency to your test project**:

   ```bash
   go get -u github.com/guidewire/fern-reporter
   ```
2. **Add the Fern Client to your Ginkgo test suite**:
   
   Import the fern client package
   ```go
   import fern "github.com/guidewire/fern-ginkgo-client/client"
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