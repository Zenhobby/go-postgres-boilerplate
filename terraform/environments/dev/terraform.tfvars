environment         = "dev"
ecs_cluster_name    = "go-postgres-boilerplate-cluster-dev"
ecs_task_family     = "go-postgres-boilerplate-task-dev"
ecs_service_name    = "go-postgres-boilerplate-service-dev"
container_port      = 8080
desired_count       = 1
cpu                 = "256"
memory              = "512"