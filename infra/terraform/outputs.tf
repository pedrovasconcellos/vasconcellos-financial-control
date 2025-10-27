# Outputs para App Runner + DocumentDB Serverless

output "app_runner_url" {
  description = "App Runner service URL pública"
  value       = aws_apprunner_service.api.service_url
}

output "docdb_endpoint" {
  description = "DocumentDB cluster endpoint (conexão privada)"
  value       = aws_docdb_cluster.mongo.endpoint
  sensitive   = true
}

output "cognito_user_pool_id" {
  description = "ID do User Pool Cognito"
  value       = aws_cognito_user_pool.this.id
}

output "cognito_client_id" {
  description = "ID do App Client Cognito"
  value       = aws_cognito_user_pool_client.this.id
}

output "s3_bucket" {
  description = "Bucket S3 para recibos"
  value       = aws_s3_bucket.receipts.bucket
}

output "sqs_queue_url" {
  description = "URL da fila SQS de transações"
  value       = aws_sqs_queue.transactions.url
}

output "sqs_dlq_url" {
  description = "URL da Dead-Letter Queue"
  value       = aws_sqs_queue.dlq.url
}

output "app_runner_service_arn" {
  description = "ARN do serviço App Runner"
  value       = aws_apprunner_service.api.arn
}

output "docdb_cluster_id" {
  description = "ID do cluster DocumentDB"
  value       = aws_docdb_cluster.mongo.cluster_identifier
}

