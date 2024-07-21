resource "aws_secretsmanager_secret" "simple-bank-db-password" {
  name = "simple-bank-db-password"
}

resource "aws_secretsmanager_secret_version" "simple-bank-db-password" {
  secret_id = aws_secretsmanager_secret.simple-bank-db-password.id
  secret_string = var.DB_PASSWORD
}