package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestRDSSnapshot(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/rds",
	}

	// At the end of the test, run `terraform destroy`
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`
	terraform.InitAndApply(t, terraformOptions)

	// Get RDS DB identifier
	identifier := terraform.Output(t, terraformOptions, "identifier")

	// Get AWS client
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})
	if err != nil {
		fmt.Println("Error creating session ", err)
	}
	client := rds.New(sess)

	// Create DB snapshot, emulating manual action
	create_input := rds.CreateDBSnapshotInput{
		DBInstanceIdentifier: aws.String(identifier),
		DBSnapshotIdentifier: aws.String("rds-testing-manual-snapshot"),
	}
	snapshot, err := client.CreateDBSnapshot(&create_input)
	if err != nil {
		fmt.Println("Error creating snapshot ", err)
	}
	fmt.Println(snapshot)

	// Wait for snapshot to be created
	snapshot_request_time := time.Now()
	describe_input := rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: aws.String("rds-testing"),
		DBSnapshotIdentifier: aws.String("rds-testing-manual-snapshot"),
	}
	var status string
	for status != "available" {
		time.Sleep(time.Second * 3)
		snapshots, err := client.DescribeDBSnapshots(&describe_input)
		if err != nil {
			fmt.Println("Error describing snapshot ", err)
		}
		status = *snapshots.DBSnapshots[0].Status
	}
	fmt.Println("Snapshot available")

	// Assert on snapshot creation time
	elapsed := time.Since(snapshot_request_time)
	assert.LessOrEqual(t, int64(elapsed), int64(time.Minute*3))
}
