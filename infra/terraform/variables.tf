variable "project" {
  description = "Nome base do projeto"
  type        = string
  default     = "financial-control"
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

variable "docdb_master_username" {
  description = "Username master do DocumentDB"
  type        = string
  default     = "financialcontrol"
}

variable "monthly_cost_limit" {
  description = "Limite de custo mensal em USD para o AWS Budgets"
  type        = string
  default     = "60"
}

variable "extra_tags" {
  description = "Tags adicionais"
  type        = map(string)
  default     = {}
}

