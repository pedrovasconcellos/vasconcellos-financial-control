# Provisionamento AWS com Terraform

## Visão geral
Este conjunto de módulos cria a infraestrutura mínima para hospedar a API no ECS Fargate de forma econômica (meta < USD 100/mês):

- ECR para armazenar a imagem do backend.
- S3 para recibos e armazenamento de objetos.
- SQS + DLQ para pipeline assíncrono.
- Cognito User Pool + App Client.
- Instância EC2 `t4g.micro` dedicada ao MongoDB (custos ~USD 8/mês + armazenamento).
- ECS Cluster + Task Definition + ALB público.
- Orçamento AWS (Budget) alertando ao atingir o limite configurado.

## Pré-requisitos
- Terraform 1.6+
- AWS CLI configurada com credenciais do ambiente (homolog/produção)
- Bucket S3 remoto para state (opcional, porém recomendado)
- Imagem Docker da API publicada no ECR (ver seção abaixo)

## Passo a passo

1. **Configurar backend remoto (opcional)**

   Edite `main.tf` ou configure no CLI:
   ```bash
   terraform init -backend-config="bucket=<bucket-state>" \
                  -backend-config="key=financial-control/terraform.tfstate" \
                  -backend-config="region=us-east-1"
   ```

2. **Criar arquivo de variáveis** (baseado em `terraform.tfvars.example`):
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   # Edite o arquivo com:
   # - environment: homolog ou production
   # - container_image: URI do ECR (ver próximo item)
   # - cpu_architecture: X86_64 ou ARM64 conforme imagem
   ```

3. **Construir e publicar a imagem**
   ```bash
   docker buildx build --platform linux/amd64,linux/arm64 -t financial-api:latest .
   aws ecr get-login-password --region us-east-1 | docker login \
       --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
   docker tag financial-api:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/financial-control-api:latest
   docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/financial-control-api:latest
   ```

4. **Aplicar a infraestrutura**
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

5. **Outputs importantes**
   - `alb_dns_name`: endpoint público da API.
   - `cognito_user_pool_id` e `cognito_user_pool_client_id`: use-os para configurar Cognito/Frontend.
   - `mongo_private_ip`: utilize para acessar o MongoDB (apenas via VPC).

## Operação
- Para atualizar a imagem, publique a nova tag no ECR e execute `terraform apply` (ou force um redeploy via ECS console).
- `desired_count` pode ser escalado alterando a variável ou via AWS Console (o serviço ignora mudanças acidentais graças ao bloco `lifecycle`).
- O orçamento (`aws_budgets_budget`) dispara alertas via console AWS; configure notificações adicionais (SNS/E-mail) manualmente conforme políticas da conta.

## Limpeza
Para remover os recursos:
```bash
terraform destroy
```
*Obs.: Remover buckets S3 vazios manualmente caso existam objetos.*
