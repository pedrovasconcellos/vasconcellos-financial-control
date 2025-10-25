variable "project" {
  description = "Nome base do projeto"
  type        = string
  default     = "finance-control"
}

variable "environment" {
  description = "Identificador do ambiente (dev|homolog|production)"
  type        = string
}

variable "region" {
  description = "Região AWS alvo"
  type        = string
  default     = "us-east-1"
}

variable "container_image" {
  description = "Imagem do container (ECR URI com tag)"
  type        = string
}

variable "desired_count" {
  description = "Número de tasks ECS em execução"
  type        = number
  default     = 1
}

variable "fargate_cpu" {
  description = "CPU reservada para a task ECS (unidades Fargate)"
  type        = string
  default     = "256"
}

variable "fargate_memory" {
  description = "Memória reservada para a task ECS (MB)"
  type        = string
  default     = "512"
}

variable "cpu_architecture" {
  description = "Arquitetura da imagem (X86_64 ou ARM64)"
  type        = string
  default     = "X86_64"
}

variable "mongo_instance_type" {
  description = "Tipo da instância EC2 que hospedará o MongoDB"
  type        = string
  default     = "t4g.micro"
}

variable "mongo_ami_ssm_parameter" {
  description = "Parâmetro SSM que aponta para a AMI do Mongo host"
  type        = string
  default     = "/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-arm64"
}

variable "monthly_cost_limit" {
  description = "Limite de custo mensal em USD para o AWS Budgets"
  type        = string
  default     = "100"
}

variable "extra_tags" {
  description = "Tags adicionais"
  type        = map(string)
  default     = {}
}
