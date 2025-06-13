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
  default     = "ECR_PUBLIC"
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
  source_security_group_id = aws_security_group.database.id
  security_group_id        = aws_security_group.app_runner.id
}

resource "aws_security_group_rule" "database_ingress_from_app_runner" {
  type                     = "ingress"
  from_port                = 5432
  to_port                  = 5432
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.app_runner.id
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
  subnets           = [aws_subnet.public_1.id, aws_subnet.public_2.id]
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
          PORT        = tostring(var.app_runner_port)
          DB_HOST     = aws_rds_cluster.main.endpoint
          DB_PORT     = tostring(aws_rds_cluster.main.port)
          DB_USER     = var.database_username
          DB_PASSWORD = var.database_password
          DB_NAME     = var.database_name
          GIN_MODE    = "release"
        }
      }
    }

    auto_deployments_enabled = false

    dynamic "authentication_configuration" {
      for_each = var.image_repository_type == "ECR" ? [1] : []
      content {
        access_role_arn = aws_iam_role.app_runner_ecr[0].arn
      }
    }
  }

  instance_configuration {
    cpu    = var.app_runner_cpu
    memory = var.app_runner_memory
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