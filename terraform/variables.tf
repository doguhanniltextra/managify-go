variable "aws_region" {
  description = "AWS Region"
  type        = string
  default     = "eu-central-1"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "managify"
}

variable "vpc_cidr" {
  description = "VPC CIDR"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  description = "Public Subnet CIDRs"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnet_cidrs" {
  description = "Private Subnet CIDRs"
  type        = list(string)
  default     = ["10.0.10.0/24", "10.0.11.0/24"]
}

variable "availability_zones" {
  description = "Availability Zones"
  type        = list(string)
  default     = ["eu-central-1a", "eu-central-1b"]
}

# Database
# Database
variable "mongo_uri" {
  description = "MongoDB Atlas Connection URI"
  type        = string
  sensitive   = true
}

# ECS & App Configuration
variable "container_image" {
  description = "Docker image URL"
  type        = string
  # Default not needed if passed from pipeline, but good to have empty if not strictly validated
}

variable "ecs_cpu" {
  description = "ECS Task CPU"
  type        = number
  default     = 256
}

variable "ecs_memory" {
  description = "ECS Task Memory"
  type        = number
  default     = 512
}

variable "ecs_desired_count" {
  description = "Desired number of ECS tasks"
  type        = number
  default     = 1
}

# Environment Variables mapping
variable "secret_key" {
  description = "App Secret Key"
  type        = string
  sensitive   = true
}

variable "api_max_limiter" {
  description = "API Max Limiter"
  type        = string
  default     = "100"
}

variable "rate_limit_expiration" {
  description = "Rate Limit Expiration"
  type        = string
  default     = "1h"
}

variable "smtp_from" {
  description = "SMTP From Address"
  type        = string
}

variable "smtp_password" {
  description = "SMTP Password"
  type        = string
  sensitive   = true
}

variable "smtp_host" {
  description = "SMTP Host"
  type        = string
}

variable "smtp_port" {
  description = "SMTP Port"
  type        = string
}

variable "mongo_db_name" {
  description = "Name of the MongoDB database"
  type        = string
  default     = "managify"
}

variable "swagger_enabled" {
  description = "Enable Swagger"
  type        = string
  default     = "true"
}

variable "metrics_enabled" {
  description = "Enable Metrics"
  type        = string
  default     = "true"
}

variable "app_port" {
  description = "Application Port"
  type        = string
  default     = "8080"
}

variable "google_client_id" {
  description = "Google OAuth Client ID"
  type        = string
}

variable "google_client_secret" {
  description = "Google OAuth Client Secret"
  type        = string
  sensitive   = true
}

variable "google_redirect_uri" {
  description = "Google OAuth Redirect URI"
  type        = string
}
