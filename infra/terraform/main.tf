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

# --- Networking helpers ---

data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# --- Budgets guardrail ---
resource "aws_budgets_budget" "monthly_cap" {
  name                = "${local.name_prefix}-cost-cap"
  budget_type         = "COST"
  limit_amount        = var.monthly_cost_limit
  limit_unit          = "USD"
  time_unit           = "MONTHLY"
  cost_types {
    include_credit = true
  }
  time_period {
    start = "2024-01-01_00:00"
  }
  lifecycle {
    prevent_destroy = true
  }
}

# --- ECR repository ---
resource "aws_ecr_repository" "api" {
  name                 = "${local.name_prefix}-api"
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }
  tags = local.tags
}

# --- S3 bucket para recibos ---
resource "random_id" "bucket" {
  byte_length = 4
}

resource "aws_s3_bucket" "receipts" {
  bucket = "${local.name_prefix}-receipts-${random_id.bucket.hex}"
  tags   = local.tags
}

resource "aws_s3_bucket_public_access_block" "receipts" {
  bucket = aws_s3_bucket.receipts.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "receipts" {
  bucket = aws_s3_bucket.receipts.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# --- SQS queue + DLQ ---
resource "aws_sqs_queue" "dlq" {
  name                       = "${local.name_prefix}-transactions-dlq"
  message_retention_seconds  = 1209600
  visibility_timeout_seconds = 30
  tags                       = local.tags
}

resource "aws_sqs_queue" "transactions" {
  name                       = "${local.name_prefix}-transactions"
  visibility_timeout_seconds = 60
  message_retention_seconds  = 1209600
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.dlq.arn,
    maxReceiveCount      = 5
  })
  tags = local.tags
}

# --- Cognito ---
resource "aws_cognito_user_pool" "this" {
  name = "${local.name_prefix}-users"
  auto_verified_attributes = ["email"]
  password_policy {
    minimum_length    = 8
    require_lowercase = true
    require_numbers   = true
    require_symbols   = false
    require_uppercase = true
  }
  tags = local.tags
}

resource "aws_cognito_user_pool_client" "this" {
  name                   = "${local.name_prefix}-client"
  user_pool_id           = aws_cognito_user_pool.this.id
  generate_secret        = false
  explicit_auth_flows    = ["ALLOW_USER_PASSWORD_AUTH", "ALLOW_REFRESH_TOKEN_AUTH"]
  prevent_user_existence_errors = "ENABLED"
}

# --- MongoDB EC2 (instância econômica) ---
resource "aws_security_group" "mongo" {
  name        = "${local.name_prefix}-mongo-sg"
  description = "MongoDB access"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    from_port       = 27017
    to_port         = 27017
    protocol        = "tcp"
    security_groups = [aws_security_group.service.id]
    description     = "ECS tasks"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.tags
}

data "aws_ssm_parameter" "mongo_ami" {
  name = var.mongo_ami_ssm_parameter
}

resource "aws_instance" "mongo" {
  ami                         = data.aws_ssm_parameter.mongo_ami.value
  instance_type               = var.mongo_instance_type
  subnet_id                   = element(data.aws_subnets.default.ids, 0)
  vpc_security_group_ids      = [aws_security_group.mongo.id]
  associate_public_ip_address = true
  iam_instance_profile        = null

  user_data = <<-EOT
              #!/bin/bash
              set -xe
              dnf install -y docker
              systemctl enable docker
              systemctl start docker
              docker run -d --restart unless-stopped -p 27017:27017 --name mongo mongo:6
              EOT

  root_block_device {
    volume_size = 20
    volume_type = "gp3"
  }

  tags = merge(local.tags, { "Name" = "${local.name_prefix}-mongo" })
}

# --- Security groups para ALB e serviço ---
resource "aws_security_group" "alb" {
  name        = "${local.name_prefix}-alb-sg"
  description = "Allow HTTP"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.tags
}

resource "aws_security_group" "service" {
  name        = "${local.name_prefix}-svc-sg"
  description = "Allow traffic from ALB"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description     = "From ALB"
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.tags
}

# --- Load balancer ---
resource "aws_lb" "api" {
  name               = "${local.name_prefix}-alb"
  load_balancer_type = "application"
  internal           = false
  security_groups    = [aws_security_group.alb.id]
  subnets            = data.aws_subnets.default.ids
  tags               = local.tags
}

resource "aws_lb_target_group" "api" {
  name        = "${local.name_prefix}-tg"
  port        = 8080
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = data.aws_vpc.default.id

  health_check {
    path                = "/api/v1/accounts"
    matcher             = "200-499"
    healthy_threshold   = 2
    unhealthy_threshold = 5
    timeout             = 5
    interval            = 30
  }

  tags = local.tags
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.api.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api.arn
  }
}

# --- CloudWatch Logs ---
resource "aws_cloudwatch_log_group" "api" {
  name              = "/ecs/${local.name_prefix}-api"
  retention_in_days = 30
  tags              = local.tags
}

# --- IAM roles ---
resource "aws_iam_role" "ecs_task_execution" {
  name               = "${local.name_prefix}-ecs-execution"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume.json
  tags               = local.tags
}

data "aws_iam_policy_document" "ecs_task_assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role" "ecs_task" {
  name               = "${local.name_prefix}-ecs-task"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume.json
  tags               = local.tags
}

resource "aws_iam_role_policy" "ecs_task" {
  name = "${local.name_prefix}-ecs-task-policy"
  role = aws_iam_role.ecs_task.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect   = "Allow",
        Action   = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Resource = "*"
      },
      {
        Effect   = "Allow",
        Action   = [
          "sqs:SendMessage",
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes"
        ],
        Resource = [aws_sqs_queue.transactions.arn, aws_sqs_queue.dlq.arn]
      },
      {
        Effect   = "Allow",
        Action   = [
          "s3:PutObject",
          "s3:GetObject"
        ],
        Resource = "${aws_s3_bucket.receipts.arn}/*"
      },
      {
        Effect   = "Allow",
        Action   = [
          "cognito-idp:InitiateAuth",
          "cognito-idp:DescribeUserPool",
          "cognito-idp:AdminGetUser"
        ],
        Resource = aws_cognito_user_pool.this.arn
      }
    ]
  })
}

# --- ECS Cluster ---
resource "aws_ecs_cluster" "this" {
  name = "${local.name_prefix}-cluster"
  setting {
    name  = "containerInsights"
    value = "enabled"
  }
  tags = local.tags
}

resource "aws_ecs_cluster_capacity_providers" "this" {
  cluster_name = aws_ecs_cluster.this.name
  capacity_providers = ["FARGATE", "FARGATE_SPOT"]
  default_capacity_provider_strategy {
    capacity_provider = "FARGATE"
    weight            = 1
  }
}

# --- ECS Task Definition ---
locals {
  container_env = [
    { name = "APP_ENVIRONMENT", value = var.environment },
    { name = "MONGO_URI", value = "mongodb://${aws_instance.mongo.private_ip}:27017/financial-control" },
    { name = "AWS_REGION", value = var.region },
    { name = "AWS_S3_BUCKET", value = aws_s3_bucket.receipts.bucket },
    { name = "AWS_SQS_QUEUENAME", value = aws_sqs_queue.transactions.name },
    { name = "AWS_SQS_QUEUEURL", value = aws_sqs_queue.transactions.url },
    { name = "QUEUE_TRANSACTIONQUEUE", value = aws_sqs_queue.transactions.name },
    { name = "STORAGE_RECEIPTBUCKET", value = aws_s3_bucket.receipts.bucket },
    { name = "AWS_COGNITO_USERPOOLID", value = aws_cognito_user_pool.this.id },
    { name = "AWS_COGNITO_CLIENTID", value = aws_cognito_user_pool_client.this.id },
    { name = "AUTH_MODE", value = "cognito" }
  ]
}

resource "aws_ecs_task_definition" "api" {
  family                   = "${local.name_prefix}-api"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.fargate_cpu
  memory                   = var.fargate_memory
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([
    {
      name      = "api"
      image     = var.container_image
      essential = true
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
          protocol      = "tcp"
        }
      ]
      environment = local.container_env
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.api.name
          awslogs-region        = var.region
          awslogs-stream-prefix = "api"
        }
      }
    }
  ])

  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture         = var.cpu_architecture
  }

  tags = local.tags
}

# --- ECS Service ---
resource "aws_ecs_service" "api" {
  name            = "${local.name_prefix}-svc"
  cluster         = aws_ecs_cluster.this.id
  task_definition = aws_ecs_task_definition.api.arn
  desired_count   = var.desired_count
  launch_type     = "FARGATE"
  platform_version = "1.4.0"

  network_configuration {
    subnets         = data.aws_subnets.default.ids
    security_groups = [aws_security_group.service.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.api.arn
    container_name   = "api"
    container_port   = 8080
  }

  lifecycle {
    ignore_changes = [desired_count]
  }

  depends_on = [aws_lb_listener.http]

  tags = local.tags
}

