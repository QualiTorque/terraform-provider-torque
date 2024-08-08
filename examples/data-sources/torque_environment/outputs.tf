output "id" {
  value = data.torque_environment.env.id
}

output "last_used" {
  value = data.torque_environment.env.last_used
}

output "owner" {
  value = data.torque_environment.env.owner_email
}

output "is_eac" {
  value = data.torque_environment.env.is_eac
}

output "blueprint_name" {
  value = data.torque_environment.env.blueprint_name
}

output "initiator" {
  value = data.torque_environment.env.initiator_email
}

output "inputs" {
  value = data.torque_environment.env.inputs
}

output "tags" {
  value = data.torque_environment.env.tags
}

output "outputs" {
  value = data.torque_environment.env.outputs
}
output "blueprint_commit" {
  value = data.torque_environment.env.blueprint_commit
}
output "blueprint_repository_name" {
  value = data.torque_environment.env.blueprint_repository_name
}
output "name" {
  value = data.torque_environment.env.name
}
output "status" {
  value = data.torque_environment.env.status
}
output "start_time" {
  value = data.torque_environment.env.start_time
}
output "end_time" {
  value = data.torque_environment.env.end_time
}

output "errors" {
  value = data.torque_environment.env.errors
}
