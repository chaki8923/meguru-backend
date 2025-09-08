terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# Variables
variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-northeast-1"
}

variable "project_name" {
  description = "Project name for tagging and naming resources"
  type        = string
  default     = "my-go-app"
}

variable "app_runner_cpu" {
  description = "App Runner CPU"
  type        = string
  default     = "0.5 vCPU"
  validation {
    condition     = contains(["0.25 vCPU", "0.5 vCPU", "1 vCPU", "2 vCPU"], var.app_runner_cpu)
    error_message = "App Runner CPU must be one of: 0.25 vCPU, 0.5 vCPU, 1 vCPU, 2 vCPU."
  }
}

variable "app_runner_memory" {
  description = "App Runner Memory"
  type        = string
  default     = "1 GB"
  validation {
    condition     = contains(["0.5 GB", "1 GB", "2 GB", "3 GB"], var.app_runner_memory)
    error_message = "App Runner Memory must be one of: 0.5 GB, 1 GB, 2 GB, 3 GB."
  }
}

variable "app_runner_min_size" {
  description = "App Runner minimum instances"
  type        = number
  default     = 1
}

variable "app_runner_max_size" {
  description = "App Runner maximum instances"
  type        = number
  default     = 2
}

variable "app_runner_port" {
  description = "App Runner port"
  type        = number
  default     = 8080
}

variable "image_repository_type" {
  description = "Image repository type"
  type        = string
  default     = "ECR"
  validation {
    condition     = contains(["ECR_PUBLIC", "ECR"], var.image_repository_type)
    error_message = "Image repository type must be either ECR_PUBLIC or ECR."
  }
}

variable "image_identifier" {
  description = "Image identifier"
  type        = string
}

variable "database_name" {
  description = "Database name"
  type        = string
  default     = "myappdb"
}

variable "database_username" {
  description = "Database username"
  type        = string
  default     = "meguru"
}

variable "database_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

variable "min_aurora_capacity" {
  description = "Minimum Aurora capacity"
  type        = number
  default     = 0.5
}

variable "max_aurora_capacity" {
  description = "Maximum Aurora capacity"
  type        = number
  default     = 1
}

variable "r2_endpoint" {
  description = "R2 endpoint URL"
  type        = string
}

variable "r2_access_key" {
  description = "R2 access key"
  type        = string
  sensitive   = true
}

variable "r2_secret_key" {
  description = "R2 secret key"
  type        = string
  sensitive   = true
}

variable "r2_bucket_url" {
  description = "R2 bucket URL"
  type        = string
}

variable "r2_public_bucket_domain" {
  description = "R2 public bucket domain"
  type        = string
}

variable "gemini_api_key" {
  description = "Gemini API key for AI services"
  type        = string
  sensitive   = true
}

# Lambda function for Gemini API calls
resource "aws_iam_role" "gemini_lambda_role" {
  name = "${var.project_name}-gemini-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name = "${var.project_name}-GeminiLambdaRole"
  }
}

resource "aws_iam_role_policy_attachment" "gemini_lambda_basic_execution" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.gemini_lambda_role.name
}

resource "aws_iam_role_policy_attachment" "gemini_lambda_vpc_execution" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
  role       = aws_iam_role.gemini_lambda_role.name
}

# App Runner Instance Role (for Lambda invocation)
resource "aws_iam_role" "app_runner_instance_role" {
  name = "${var.project_name}-apprunner-instance-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "tasks.apprunner.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name = "${var.project_name}-AppRunnerInstanceRole"
  }
}

resource "aws_iam_role_policy" "app_runner_lambda_invoke" {
  name = "lambda-invoke-policy"
  role = aws_iam_role.app_runner_instance_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = "lambda:InvokeFunction"
        Resource = aws_lambda_function.gemini_api.arn
      }
    ]
  })
}

resource "aws_lambda_function" "gemini_api" {
  function_name = "${var.project_name}-gemini-api"
  role         = aws_iam_role.gemini_lambda_role.arn
  handler      = "bootstrap"
  runtime      = "provided.al2023"
  timeout      = 120

  # Placeholder zip file - we'll update this later
  filename         = "gemini-lambda.zip"
  source_code_hash = filebase64sha256("gemini-lambda.zip")

  vpc_config {
    subnet_ids         = [aws_subnet.database_1.id, aws_subnet.database_2.id]
    security_group_ids = [aws_security_group.lambda.id]
  }

  environment {
    variables = {
      GEMINI_API_KEY = var.gemini_api_key
    }
  }

  tags = {
    Name = "${var.project_name}-GeminiLambda"
  }
}


# Data sources
data "aws_availability_zones" "available" {
  state = "available"
}

# VPC
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "${var.project_name}-VPC"
  }
}

# Subnets
resource "aws_subnet" "public_1" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = data.aws_availability_zones.available.names[0]
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.project_name}-PublicSubnet1"
  }
}

resource "aws_subnet" "public_2" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.2.0/24"
  availability_zone       = data.aws_availability_zones.available.names[1]
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.project_name}-PublicSubnet2"
  }
}

resource "aws_subnet" "database_1" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.3.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = {
    Name = "${var.project_name}-DatabaseSubnet1"
  }
}

resource "aws_subnet" "database_2" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.4.0/24"
  availability_zone = data.aws_availability_zones.available.names[1]

  tags = {
    Name = "${var.project_name}-DatabaseSubnet2"
  }
}

# Internet Gateway
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "${var.project_name}-IGW"
  }
}

# Route Table
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "${var.project_name}-PublicRouteTable"
  }
}

resource "aws_route_table_association" "public_1" {
  subnet_id      = aws_subnet.public_1.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "public_2" {
  subnet_id      = aws_subnet.public_2.id
  route_table_id = aws_route_table.public.id
}

# Security Groups
resource "aws_security_group" "lambda" {
  name_prefix = "${var.project_name}-lambda-"
  vpc_id      = aws_vpc.main.id
  description = "Security group for Lambda functions"

  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTPS outbound for Gemini API"
  }

  tags = {
    Name = "${var.project_name}-LambdaSecurityGroup"
  }
}

resource "aws_security_group" "app_runner" {
  name_prefix = "${var.project_name}-apprunner-"
  vpc_id      = aws_vpc.main.id
  description = "Security Group for App Runner VPC Connector"

  tags = {
    Name = "${var.project_name}-AppRunnerSG"
  }
}

resource "aws_security_group" "database" {
  name_prefix = "${var.project_name}-database-"
  vpc_id      = aws_vpc.main.id
  description = "Security Group for Aurora Database"

  tags = {
    Name = "${var.project_name}-DatabaseSG"
  }
}

# Security Group Rules (separate to avoid circular dependency)
resource "aws_security_group_rule" "app_runner_egress_to_db" {
  type                     = "egress"
  from_port                = 5432
  to_port                  = 5432
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.database.id # これがDBのSGを指す（宛先として機能）
  security_group_id        = aws_security_group.app_runner.id  # このルールが適用されるSG
}

resource "aws_security_group_rule" "database_ingress_from_app_runner" {
  type                     = "ingress"
  from_port                = 5432
  to_port                  = 5432
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.app_runner.id # こちらは正しい
  security_group_id        = aws_security_group.database.id
}

# Aurora Resources
resource "aws_db_subnet_group" "main" {
  name       = "${var.project_name}-db-subnet-group"
  subnet_ids = [aws_subnet.database_1.id, aws_subnet.database_2.id]

  tags = {
    Name = "${var.project_name}-DBSubnetGroup"
  }
}

resource "aws_rds_cluster_parameter_group" "main" {
  family = "aurora-postgresql15"
  name   = "${var.project_name}-cluster-pg"
  
  parameter {
    name  = "rds.force_ssl"
    value = "0"
  }

  tags = {
    Name = "${var.project_name}-ClusterParameterGroup"
  }
}

resource "aws_rds_cluster" "main" {
  cluster_identifier                = "${var.project_name}-aurora-cluster"
  engine                           = "aurora-postgresql"
  engine_version                   = "15.3"
  database_name                    = var.database_name
  master_username                  = var.database_username
  master_password                  = var.database_password
  db_subnet_group_name             = aws_db_subnet_group.main.name
  vpc_security_group_ids           = [aws_security_group.database.id]
  db_cluster_parameter_group_name  = aws_rds_cluster_parameter_group.main.name
  backup_retention_period          = 1
  deletion_protection              = false
  skip_final_snapshot             = true

  serverlessv2_scaling_configuration {
    max_capacity = var.max_aurora_capacity
    min_capacity = var.min_aurora_capacity
  }

  tags = {
    Name = "${var.project_name}-AuroraCluster"
  }
}

resource "aws_rds_cluster_instance" "main" {
  identifier           = "${var.project_name}-aurora-instance-1"
  cluster_identifier   = aws_rds_cluster.main.id
  instance_class       = "db.serverless"
  engine               = aws_rds_cluster.main.engine
  engine_version       = aws_rds_cluster.main.engine_version
  publicly_accessible = false

  tags = {
    Name = "${var.project_name}-AuroraInstance1"
  }
}

# App Runner VPC Connector
resource "aws_apprunner_vpc_connector" "main" {
  vpc_connector_name = "${var.project_name}-vpc-connector"
  subnets           = [aws_subnet.database_1.id, aws_subnet.database_2.id]
  security_groups   = [aws_security_group.app_runner.id]

  tags = {
    Name = "${var.project_name}-AppRunnerVPCConnector"
  }
}

# IAM Role for ECR Access (conditional)
resource "aws_iam_role" "app_runner_ecr" {
  count = var.image_repository_type == "ECR" ? 1 : 0
  name  = "${var.project_name}-apprunner-ecr-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "build.apprunner.amazonaws.com"
        }
      }
    ]
  })

  inline_policy {
    name = "ECRAccessPolicy"
    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          Effect = "Allow"
          Action = [
            "ecr:GetDownloadUrlForLayer",
            "ecr:BatchGetImage",
            "ecr:BatchCheckLayerAvailability",
            "ecr:DescribeImages",
            "ecr:GetAuthorizationToken"
          ]
          Resource = "*"
        }
      ]
    })
  }

  tags = {
    Name = "${var.project_name}-AppRunnerECRRole"
  }
}

# Auto Scaling Configuration
resource "aws_apprunner_auto_scaling_configuration_version" "main" {
  auto_scaling_configuration_name = "${var.project_name}-as-config"
  min_size                       = var.app_runner_min_size
  max_size                       = var.app_runner_max_size

  tags = {
    Name      = "${var.project_name}-AppRunnerASConfig"
    ManagedBy = "Terraform"
  }
}

# App Runner Service
resource "aws_apprunner_service" "main" {
  service_name = "${var.project_name}-service"

  source_configuration {
    image_repository {
      image_identifier      = var.image_identifier
      image_repository_type = var.image_repository_type

      image_configuration {
        port = var.app_runner_port
        runtime_environment_variables = {
          PORT                    = tostring(var.app_runner_port)
          DB_HOST                 = aws_rds_cluster.main.endpoint
          DB_PORT                 = tostring(aws_rds_cluster.main.port)
          DB_USER                 = var.database_username
          DB_PASSWORD             = var.database_password
          DB_NAME                 = var.database_name
          DATABASE_URL            = "postgres://${var.database_username}:${var.database_password}@${aws_rds_cluster.main.endpoint}:${aws_rds_cluster.main.port}/${var.database_name}?sslmode=disable"
          R2_ENDPOINT             = var.r2_endpoint
          R2_ACCESS_KEY           = var.r2_access_key
          R2_SECRET_KEY           = var.r2_secret_key
          R2_BUCKET_URL           = var.r2_bucket_url
          R2_PUBLIC_BUCKET_DOMAIN = var.r2_public_bucket_domain
          GEMINI_API_KEY          = var.gemini_api_key
          GEMINI_LAMBDA_FUNCTION  = aws_lambda_function.gemini_api.function_name
          AWS_REGION              = "ap-northeast-1"
          GIN_MODE                = "release"
        }
      }
    }

    auto_deployments_enabled = true

    dynamic "authentication_configuration" {
      for_each = var.image_repository_type == "ECR" ? [1] : []
      content {
        access_role_arn = aws_iam_role.app_runner_ecr[0].arn
      }
    }
  }

  instance_configuration {
    instance_role_arn = aws_iam_role.app_runner_instance_role.arn
    cpu               = var.app_runner_cpu
    memory            = var.app_runner_memory
  }

  network_configuration {
    egress_configuration {
      egress_type       = "VPC"
      vpc_connector_arn = aws_apprunner_vpc_connector.main.arn
    }
  }

  health_check_configuration {
    healthy_threshold   = 1
    interval            = 10
    path                = "/health"
    protocol            = "HTTP"
    timeout             = 5
    unhealthy_threshold = 3
  }

  auto_scaling_configuration_arn = aws_apprunner_auto_scaling_configuration_version.main.arn

  tags = {
    Name      = "${var.project_name}-AppRunnerService"
    ManagedBy = "Terraform"
  }

  depends_on = [
    aws_rds_cluster_instance.main
  ]
}

# Outputs
output "app_runner_service_url" {
  value = aws_apprunner_service.main.service_url
}

output "aurora_cluster_endpoint" {
  value = aws_rds_cluster.main.endpoint
}

output "vpc_id" {
  value = aws_vpc.main.id
}

output "auto_scaling_config_arn" {
  value = aws_apprunner_auto_scaling_configuration_version.main.arn
}
