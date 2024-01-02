terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = "https://review2.qualilabs.net/"
  space = "Shira"
  token = ""
}

resource "torque_space_resource" "new_space" {
  name  = "newspacetest"
  color = "darkGreen"
  icon  = "flow"
  # space_members = ["admontomer@gmail.com"]
  # space_admins  = ["sgeller1980@gmail.com"]
  # associated_agents = ["tomer-test"]
  associated_kubernetes_agent = [
    {
      default_namespace       = "vido-sb"
      default_service_account = "vido-sb"
      agent_name              = "review2-aks"
    }
  ]
  ## TODO: associated_repos = []
}
