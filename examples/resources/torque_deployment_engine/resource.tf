resource "torque_deployment_engine" "engine" {
  name                     = "name"
  description              = "description"
  agent_name               = "agent"
  auth_token               = "token"
  server_url               = "https://argocd.com"
  polling_interval_seconds = 30
  all_spaces               = true
}
