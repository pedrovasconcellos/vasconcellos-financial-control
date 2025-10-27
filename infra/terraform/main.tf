# Provisionamento AWS com App Runner + DocumentDB Serverless
# Custo estimado: ~USD 60/mês

terraform {
  required_version = ">= 1.6.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.5"
    }
  }
}

provider "aws" {
  region = var.region
}

locals {
  name_prefix = "${var.project}-${var.environment}"
  tags = merge({
    "Project"     : var.project,
    "Environment" : var.environment,
    "ManagedBy"   : "terraform"
  }, var.extra_tags)
}

# --- VPC para DocumentDB Serverless ---
data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# --- Security Group para DocumentDB ---
resource "aws_security_group" "docdb" {
  name        = "${local.name_prefix}-docdb-sg"
  description = "DocumentDB Serverless access"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    from_port   = 27017
    to_port     = 27017
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "MongoDB from App Runner"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.tags
}

# --- DocumentDB Serverless Cluster ---
# Usando STANDARD mode (mais barato: USD 0.0822/DCU-hora)
resource "aws_docdb_cluster" "mongo" {
  cluster_identifier      = "${local.name_prefix}-docdb"
  engine                  = "docdb"
  engine_version          = "5.0.0"
  master_username         = var.docdb_master_username
  master_password         = random_password.docdb_password.result
  db_subnet_group_name    = aws_docdb_subnet_group.this.name
  vpc_security_group_ids  = [aws_security_group.docdb.id]
  
  # Serverless configuration - STANDARD (mais barato)
  # Cada DCU = ~2GB RAM + CPU + rede
  serverlessv2_scaling_configuration {
    max_capacity = 16.0    # 16 DCUs máximo (~32GB RAM)
    min_capacity = 2.0     # 2 DCUs mínimo (~4GB RAM, evita cold start)
  }

  skip_final_snapshot = true
  
  tags = local.tags
}

# Optional: Use I/O-Optimized mode (10% mais caro mas melhor para I/O intensivo)
# Custo: USD 0.0905/DCU-hora vs USD 0.0822/DCU-hora no Standard

resource "aws_docdb_cluster_instance" "mongo" {
  identifier         = "${local.name_prefix}-docdb-instance"
  cluster_identifier = aws_docdb_cluster.mongo.id
  instance_class     = "db.serverless"
}

resource "aws_docdb_subnet_group" "this" {
  name       = "${local.name_prefix}-docdb-subnet"
  subnet_ids = data.aws_subnets.default.ids
  
  tags = local.tags
}

# --- Password para DocumentDB ---
resource "random_password" "docdb_password" {
  length  = 32
  special = true
}

# --- App Runner Service (API) ---
resource "aws_apprunner_service" "api" {
  service_name = "${local.name_prefix}-api"

  source_configuration {
    auto_deployments_enabled = false
    
    image_repository {
      image_identifier      = var.container_image
      image_configuration {
        port = "8080"
        runtime_environment_variables = {
          APP_ENVIRONMENT              = var.environment
          MONGO_URI                    = "mongodb://${aws_docdb_cluster.mongo.master_username}:${random_password.docdb_password.result}@${aws_docdb_cluster.mongo.endpoint}:27017/financial-control?tls=true"
          AWS_REGION                   = var.region
          AWS_S3_BUCKET                = aws_s3_bucket.receipts.bucket
          AWS_SQS_QUEUENAME            = aws_sqs_queue.transactions.name
          AWS_SQS_QUEUEURL             = aws_sqs_queue.transactions.url
          QUEUE_TRANSACTIONQUEUE       = aws_sqs_queue.transactions.name
          STORAGE_RECEIPTBUCKET        = aws_s3_bucket.receipts.bucket
          AWS_COGNITO_USERPOOLID       = aws_cognito_user_pool.this.id
          AWS_COGNITO_CLIENTID         = aws_cognito_user_pool_client.this.id
          AUTH_MODE                    = "cognito"
        }
      }
      image_repository_type = "ECR"
    }
  }

  instance_configuration {
    cpu               = "0.25 vCPU"
    memory            = "0.5 GB"
    instance_role_arn = aws_iam_role.apprunner_instance.arn
  }

  tags = local.tags
}

# --- IAM Roles ---
resource "aws_iam_role" "apprunner_instance" {
  name = "${local.name_prefix}-apprunner-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Principal = {
        Service = "tasks.apprunner.amazonaws.com"
      }
      Effect = "Allow"
    }]
  })

  tags = local.tags
}

resource "aws_iam_role_policy" "apprunner_policy" {
  name = "${local.name_prefix}-apprunner-policy"
  role = aws_iam_role.apprunner_instance.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ]
        Resource = "${aws_s3_bucket.receipts.arn}/*"
      },
      {
        Effect = "Allow"
        Action = [
          "sqs:SendMessage",
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes"
        ]
        Resource = [aws_sqs_queue.transactions.arn, aws_sqs_queue.dlq.arn]
      },
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "*"
      }
    ]
  })
}

# --- S3 Bucket para Receipts ---
resource "aws_s3_bucket" "receipts" {
  bucket = "${local.name_prefix}-receipts"
  tags   = local.tags
}

resource "aws_s3_bucket_versioning" "receipts" {
  bucket = aws_s3_bucket.receipts.id
  versioning_configuration {
    status = "Enabled"
  }
}

# --- SQS Queue com DLQ ---
resource "aws_sqs_queue" "dlq" {
  name                      = "${local.name_prefix}-dlq"
  message_retention_seconds = 1209600 # 14 dias
  tags                      = local.tags
}

resource "aws_sqs_queue" "transactions" {
  name                      = "${local.name_prefix}-transactions"
  visibility_timeout_seconds = 60
  message_retention_seconds  = 345600 # 4 dias
  
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.dlq.arn
    maxReceiveCount     = 5
  })
  
  tags = local.tags
}

# --- Cognito User Pool ---
resource "aws_cognito_user_pool" "this" {
  name     = "${local.name_prefix}-users"
  tags     = local.tags
  
  username_configuration {
    case_sensitive = false
  }
  
  password_policy {
    minimum_length    = 8
    require_lowercase = true
    require_numbers   = true
    require_symbols   = true
    require_uppercase = true
  }
  
  auto_verified_attributes = ["email"]
  
  schema {
    name     = "email"
    required = true
    mutable  = true
  }
}

resource "aws_cognito_user_pool_client" "this" {
  name         = "${local.name_prefix}-client"
  user_pool_id = aws_cognito_user_pool.this.id
  
  allowed_oauth_flows                  = ["code"]
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_scopes                 = ["email", "openid", "profile"]
  
  callback_urls = ["https://example.com/callback"]
  logout_urls   = ["https://example.com/logout"]
  
  generate_secret = true
}

# --- CloudWatch Log Group ---
resource "aws_cloudwatch_log_group" "apprunner" {
  name              = "/aws/apprunner/${local.name_prefix}-api"
  retention_in_days = 30
  tags              = local.tags
}

# Outputs movidos para outputs.tf

