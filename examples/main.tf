# ---------------------------------------------------------------------------------------------------------------------
# PIN TERRAFORM VERSION TO >= 0.12
# The examples have been upgraded to 0.12 syntax
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  required_version = ">= 0.12"
}

# ---------------------------------------------------------------------------------------------------------------------
# BASIC TERRAFORM EXAMPLE
# See test/terraform_aws_example.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

data "template_file" "example" {
  template = var.example
}

data "template_file" "example2" {
  template = var.example2
}

resource "local_file" "example" {
  content  = "${data.template_file.example.rendered} + ${data.template_file.example2.rendered}"
  filename = "example.txt"
}

resource "local_file" "example2" {
  content  = data.template_file.example2.rendered
  filename = "example2.txt"
}

