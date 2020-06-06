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
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

func TestRDSSnapshot(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/rds",
	}

	snapshotIdentifier := "rds-testing-manual-snapshot"

	// Create shared AWS client
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})
	if err != nil {
		fmt.Println("Error creating session ", err)
	}
	client := rds.New(sess)

	defer test_structure.RunTestStage(t, "cleanup_snapshots", func() {
		// Cleanup manual snapshot
		delete_input := rds.DeleteDBSnapshotInput{
			DBSnapshotIdentifier: aws.String(snapshotIdentifier),
		}
		deleted_snapshot, err := client.DeleteDBSnapshot(&delete_input)
		if err != nil {
			fmt.Println("Error cleaning up snapshot ", err)
		}
		fmt.Println(deleted_snapshot)
	})

	defer test_structure.RunTestStage(t, "terraform_destroy", func() {
		terraform.Destroy(t, terraformOptions)
	})

	test_structure.RunTestStage(t, "terraform_apply", func() {
		terraform.InitAndApply(t, terraformOptions)
	})

	test_structure.RunTestStage(t, "create_manual_snapshot", func() {
		// Get RDS DB identifier
		identifier := terraform.Output(t, terraformOptions, "identifier")
		createManualSnapshot(t, identifier, snapshotIdentifier, client)
	})
}

func createManualSnapshot(t *testing.T, dbIdentifier string, snapshotIdentifier string, client *rds.RDS) {
	// Create DB snapshot, emulating manual action
	create_input := rds.CreateDBSnapshotInput{
		DBInstanceIdentifier: aws.String(dbIdentifier),
		DBSnapshotIdentifier: aws.String(snapshotIdentifier),
	}
	snapshot, err := client.CreateDBSnapshot(&create_input)
	if err != nil {
		fmt.Println("Error creating snapshot ", err)
	}
	fmt.Println(snapshot)

	// Wait for snapshot to be created
	snapshot_request_time := time.Now()
	describe_input := rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: aws.String(dbIdentifier),
		DBSnapshotIdentifier: aws.String(snapshotIdentifier),
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
