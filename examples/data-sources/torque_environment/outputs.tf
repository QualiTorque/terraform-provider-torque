output "read_only" {
  value = data.torque_environment.env.read_only
}

output "workflow" {
  value = data.torque_environment.env.is_workflow
}

output "id" {
  value = data.torque_environment.env.environment_id
}
