terraform {
  required_version = ">= 0.12"
}

provider "aws" {
  region  = "eu-west-1"
  profile = "personal"
}

resource "aws_db_instance" "default" {
  apply_immediately = true
  allocated_storage = 100
  backup_retention_period = 14
  copy_tags_to_snapshot = true
  delete_automated_backups = false
  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]
  engine = "postgres"
  engine_version = "11.5"
  final_snapshot_identifier = "rds-testing-final-snapshot"
  iam_database_authentication_enabled = false
  identifier = "rds-testing"
  instance_class = "db.t3.medium"
  multi_az = false
  name = "rdstesting"
  skip_final_snapshot = true
  storage_encrypted = true
  storage_type = "gp2"
  username = "postgres"
  password = "testing123"
}

output "identifier" {
  value = aws_db_instance.default.id
}
