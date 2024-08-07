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
