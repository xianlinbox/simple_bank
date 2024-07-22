data "aws_secretsmanager_secret" "simple-bank-db-password" {
  name = "simple-bank-db-password"

}

data "aws_secretsmanager_secret_version" "simple-bank-db-password" {
  secret_id = data.aws_secretsmanager_secret.simple-bank-db-password
}

resource "aws_db_instance" "simple_bank_db" {
  allocated_storage   = 10
  db_name             = "simple-bank"
  engine              = "postgres"
  engine_version      = "16.0"
  instance_class      = "db.t3.micro"
  username            = "root"
  password            = data.aws_secretsmanager_secret_version.simple-bank-db-password
  skip_final_snapshot = true

  depends_on = [
    aws_secretsmanager_secret.simple-bank-db-password,
    aws_secretsmanager_secret_version.simple-bank-db-password
  ]
}
