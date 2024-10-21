provider "aws" {
  region = var.aws_region
}

module "ecs" {
  source = "./modules/ecs"

  environment      = var.environment
  ecs_cluster_name = var.ecs_cluster_name
  ecs_task_family  = var.ecs_task_family
  ecs_service_name = var.ecs_service_name
  container_port   = var.container_port
  desired_count    = var.desired_count
  cpu              = var.cpu
  memory           = var.memory
}

variable "aws_region" {}
variable "environment" {}
variable "ecs_cluster_name" {}
variable "ecs_task_family" {}
variable "ecs_service_name" {}
variable "container_port" {}
variable "desired_count" {}
variable "cpu" {}
variable "memory" {}