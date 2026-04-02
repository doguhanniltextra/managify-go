terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  # Backend configuration - S3'te state sakla
  backend "s3" {
    bucket         = "managify-tfstate"
    key            = "terraform.tfstate"
    region         = "eu-central-1"
    encrypt        = true
    dynamodb_table = "managify-locks"
  }
}

provider "aws" {
  region = var.aws_region
}

module "ecr" {
  source = "./modules/ecr"

  repository_name = var.project_name
}

module "network" {
  source = "./modules/network"

  project_name         = var.project_name
  vpc_cidr             = var.vpc_cidr
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
  availability_zones   = var.availability_zones
}


module "loadbalancer" {
  source = "./modules/loadbalancer"

  project_name      = var.project_name
  vpc_id            = module.network.vpc_id
  public_subnet_ids = module.network.public_subnet_ids
}

module "ecs" {
  source = "./modules/ecs"

  project_name          = var.project_name
  vpc_id                = module.network.vpc_id
  private_subnet_ids    = module.network.private_subnet_ids
  alb_security_group_id = module.loadbalancer.alb_security_group_id
  target_group_arn      = module.loadbalancer.target_group_arn
  container_image       = var.container_image
  cpu                   = var.ecs_cpu
  memory                = var.ecs_memory
  desired_count         = var.ecs_desired_count
  aws_region            = var.aws_region
  
  container_environment = [
    { name = "MONGO_URI", value = var.mongo_uri },
    { name = "MONGO_DB", value = var.mongo_db_name },
    { name = "PORT", value = var.app_port },
    { name = "SECRET_KEY", value = var.secret_key },
    
    # SMTP
    { name = "SMTP_FROM", value = var.smtp_from },
    { name = "SMTP_PASSWORD", value = var.smtp_password },
    { name = "SMTP_HOST", value = var.smtp_host },
    { name = "SMTP_PORT", value = var.smtp_port },
    
    # App Limits
    { name = "API_MAX_LIMITER", value = var.api_max_limiter },
    { name = "RATE_LIMIT_EXPIRATION", value = var.rate_limit_expiration },

    # Features
    { name = "SWAGGER", value = var.swagger_enabled },
    { name = "METRICS", value = var.metrics_enabled },
    
    # Google OAuth
    { name = "GOOGLE_CLIENT_ID", value = var.google_client_id },
    { name = "GOOGLE_CLIENT_SECRET", value = var.google_client_secret },
    { name = "GOOGLE_REDIRECT_URI", value = var.google_redirect_uri }
  ]
}

