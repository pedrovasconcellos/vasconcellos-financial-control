output "alb_dns_name" {
  description = "Endpoint público da API"
  value       = aws_lb.api.dns_name
}

output "ecr_repository_url" {
  description = "URI do repositório ECR para push da imagem"
  value       = aws_ecr_repository.api.repository_url
}

output "s3_receipts_bucket" {
  description = "Bucket onde recibos são armazenados"
  value       = aws_s3_bucket.receipts.bucket
}

output "sqs_queue_url" {
  description = "URL da fila SQS de transações"
  value       = aws_sqs_queue.transactions.url
}

output "cognito_user_pool_id" {
  description = "ID do User Pool Cognito"
  value       = aws_cognito_user_pool.this.id
}

output "cognito_user_pool_client_id" {
  description = "ID do App Client Cognito"
  value       = aws_cognito_user_pool_client.this.id
}

output "mongo_private_ip" {
  description = "IP privado da instância Mongo"
  value       = aws_instance.mongo.private_ip
}
