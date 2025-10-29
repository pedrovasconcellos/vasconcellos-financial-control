# Provisionamento AWS com App Runner + DocumentDB Serverless

## Visão Geral
Esta infraestrutura é otimizada para custos (~USD 45-50/mês) usando:
- **App Runner** para hospedar a API (USD 7/mês)
- **DocumentDB Serverless** para MongoDB (USD 36/mês)
- Remoção de ALB (App Runner tem balancer embutido)
- Cognito, S3, SQS para recursos auxiliares

## Custo Estimado Total

| Componente | Custo/mês (USD) |
|------------|-----------------|
| App Runner | USD 7 |
| DocumentDB Serverless (min 2 DCU @ USD 0.0822/DCU-hora) | USD 36 |
| S3 (receipts bucket) | USD 0.50 |
| SQS + DLQ | USD 0 |
| Cognito | USD 0 (dentro do free tier) |
| CloudWatch Logs | USD 1-2 |
| **Total** | **USD 45-50/mês** |

**Nota**: O cálculo do DocumentDB é: 2 DCUs mínimos × 730h/mês × USD 0.0822/DCU-hora = USD 36/mês (modo STANDARD, mais econômico).


## Arquivos Principais

- `main.tf` - Recursos AWS (App Runner, DocumentDB, etc)
- `variables.tf` - Variáveis de entrada
- `outputs.tf` - Outputs da infraestrutura
- `terraform.tfvars.example` - Exemplo de configuração

## Como Usar

### 1. Configurar variáveis
```bash
cp terraform.tfvars.example terraform.tfvars
# Edite o arquivo com suas configurações
```

### 2. Build e push da imagem Docker
```bash
# Build da imagem
docker buildx build --platform linux/amd64 -t financial-api:latest .

# Login no ECR
aws ecr get-login-password --region us-east-1 | docker login \
    --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com

# Push para ECR
docker tag financial-api:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/financial-control-api:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/financial-control-api:latest
```

### 3. Aplicar Terraform
```bash
terraform init
terraform plan
terraform apply
```

## Outputs Importantes

Após `terraform apply`, você terá:
- `app_runner_url` - URL pública da sua API
- `docdb_endpoint` - Endpoint do DocumentDB (sensível)
- `cognito_user_pool_id` - Para configurar autenticação
- `cognito_client_id` - Cliente Cognito
- `s3_bucket` - Bucket para recibos
- `sqs_queue_url` - Fila de mensagens

## Diferenças vs Infra Anterior

| Aspecto | Infra Anterior (Fargate) | Nova (App Runner) |
|---------|--------------------------|-------------------|
| Serviço | ECS Fargate + ALB | App Runner (ALB embutido) |
| DB | EC2 t4g.micro | DocumentDB Serverless |
| Custo | USD 68/mês | USD 45-50/mês |
| Complexidade | Alta | Baixa |
| Escalabilidade | Manual | Automática |

## Vantagens

✅ **Custo reduzido**: USD 45-50 vs USD 68  
✅ **Sem ALB**: App Runner já inclui load balancer  
✅ **Escalabilidade automática**: DocumentDB Serverless  
✅ **Menos infra**: Menos recursos para gerenciar  
✅ **Cold start evitado**: MinCapacity = 2 DCU (~4GB RAM)  

## Limitações

⚠️ **DocumentDB é mais caro** que EC2 + MongoDB  
⚠️ **Dependência da AWS**: Migração mais difícil  
⚠️ **Menos controle**: App Runner não permite configurações avançadas

## Recomendação

Se custo é prioridade, considere voltar para EC2 + MongoDB (~USD 10/mês).  
Esta infra usa DocumentDB para escalabilidade e gerenciamento automático.

