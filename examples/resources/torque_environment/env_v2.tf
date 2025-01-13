######################################
###### Option 1 ######################
######################################

resource "torque_environment_v2" "example" {
  name    = ""
  inputs  = {}
  outputs = {}
  grains = [{
    "grain" : data.torque_grain.grain,
    "depends-on" : data.torque_grain.grain2
    },
    {
      "grain" : data.torque_grain.grain2,
    }
  ]
}

data "torque_grain" "grain" {
  kind    = "terraform"
  source  = "repo"
  path    = "path/to/terraform/module"
  inputs  = {}
  outputs = {}
}

######################################
###### Option 2 ######################
######################################
resource "torque_environment_yaml" "example" {
  yaml = <<-EOF
    name: some_env_name
    inputs:
      agent: agent
    outputs:
      - output1
      - output2
    grains:
      - kind: terraform
      - source:
  EOF
}

resource "torque_environment_v2" "example" {
  name    = ""
  inputs  = {}
  outputs = {}
  grains = [{
    "grain" : data.torque_grain.grain,
    "depends-on" : data.torque_grain.grain2
    },
    {
      "grain" : data.torque_grain.grain2,
    }
  ]
}

resource "torque_environment_v2" "example" {
  name    = ""
  inputs  = {}
  outputs = {}
  grains = [{
    "grain" : data.torque_grain.grain,
    "depends-on" : data.torque_grain.grain2
    },
    {
      "grain" : data.torque_grain.grain2,
    }
  ]
  ######### OR optional attribute instead of name, inputs,outputs, grains...
  yaml = <<-EOF
    name: some_env_name
    inputs:
      agent: agent
    outputs:
      - output1
      - output2
    grains:
      - kind: terraform
      - source:
  EOF
}
