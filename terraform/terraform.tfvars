project_name = "managify"
aws_region   = "eu-central-1"

# App Secrets
secret_key    = "your-secret-key-here"
smtp_password = "your-smtp-password"

# App Config
smtp_from   = "noreply@example.com"
smtp_host   = "smtp.example.com"
smtp_port   = "587"

# ECS
container_image = "your-account-id.dkr.ecr.eu-central-1.amazonaws.com/managify:latest"

# Optional Overrides
# api_max_limiter = "100"
# rate_limit_expiration = "1h"
# mongo_db_name = "managify"
