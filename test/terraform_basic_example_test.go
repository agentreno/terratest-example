package test

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
)

func TestWebServer(t *testing.T) {
  terraformOptions := &terraform.Options {
    TerraformDir: "../examples",
  }

  // At the end of the test, run `terraform destroy`
  defer terraform.Destroy(t, terraformOptions)

  // Run `terraform init` and `terraform apply`
  terraform.InitAndApply(t, terraformOptions)

  // Run `terraform output` to get the value of an output variable
  url := terraform.Output(t, terraformOptions, "url")

  // Verify that we get back a 200 OK with the expected text. It
  // takes ~1 min for the Instance to boot, so retry a few times.
  status := 200
  text := "Hello, World"
  retries := 30
  sleep := 10 * time.Second
  http_helper.HttpGetWithRetry(t, url, nil, status, text, retries, sleep)
}
