output "db_password_secret" {
  value = "${aws_secretsmanager_secret_version.simple-bank-db-password.secret_string}"
}