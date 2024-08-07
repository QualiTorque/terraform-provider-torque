output "read_only" {
  value = data.torque_environment.env.read_only
}

output "workflow" {
  value = data.torque_environment.env.is_workflow
}
