resource "random_password" "db_password" {
  length           = 32
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "aws_secretsmanager_secret" "db_credentials" {
  name                    = "augment-fund/${var.environment}/db-credentials"
  description             = "Database credentials for augment-fund ${var.environment}"
  recovery_window_in_days = var.environment == "prod" ? 30 : 0
}

resource "aws_secretsmanager_secret_version" "db_credentials" {
  secret_id = aws_secretsmanager_secret.db_credentials.id
  secret_string = jsonencode({
    username = var.db_username
    password = random_password.db_password.result
    dbname   = var.db_name
  })
}

data "aws_iam_policy_document" "secrets_access" {
  statement {
    sid    = "AllowSecretsAccess"
    effect = "Allow"
    actions = [
      "secretsmanager:GetSecretValue",
    ]
    resources = [
      aws_secretsmanager_secret.db_credentials.arn,
    ]
  }
}

resource "aws_iam_policy" "secrets_access" {
  name        = "augment-fund-${var.environment}-secrets-access"
  description = "Allow access to augment-fund secrets"
  policy      = data.aws_iam_policy_document.secrets_access.json
}
