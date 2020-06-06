# terratest-example

## Description

Example of using [terratest](https://terratest.gruntwork.io/) to test
infrastructure code.

One example using the AWS Hello World from their quick start, which starts an
HTTP server on an EC2 instance and waits for it to respond 200.

Another example of my own using RDS, testing a manual snapshot is completed
within 3 minutes.

## Running

Setup a `personal` profile in local AWS config, then run `cd test && go test -v
-timeout 30m ./rds_test.go`. Set env vars `SKIP_x=true` to disable any of the
RDS test stages:

- `terraform_apply`
- `create_manual_snapshot`
- `terraform_destroy`
- `cleanup_snapshots`
