# Finance Control Platform

Finance Control is a clean-architecture Go backend paired with a React + Material UI frontend for managing personal finances. It integrates with MongoDB for persistence, Amazon S3 and SQS (via LocalStack in local environments), and Cognito for authentication. Asynchronous workloads are handled by an AWS Lambda written in Go that consumes the transactional queue.

## Features

- **Authentication** powered by Amazon Cognito, with a local fallback profile for offline development.
- **Account, category, transaction, budget, goal and reporting APIs** organised strictly around clean architecture boundaries.
- **Async pipeline** where each transaction publishes an SQS message processed by a Lambda that updates budget execution.
- **Asset storage** through S3 for transaction receipts (presigned URLs returned to clients).
- **TypeScript React UI** with Material UI, React Query, and modern UX for dashboards, accounts, transactions, budgets, and goals.
- **Docker-ready** setup with MongoDB and LocalStack, plus bootstrap scripts that provision queues, topics, and Cognito assets.

## Repository layout

```
.
├── cmd/
│   └── api/               # Go HTTP API entry point
├── internal/              # Clean architecture core (entities, use cases, infra)
├── src/
│   └── lambdas/
│       └── transaction_processor/ # AWS Lambda handler
├── frontend/              # React application (Vite + MUI)
├── config/                # Environment configuration templates
├── scripts/               # LocalStack bootstrap scripts
├── Dockerfile             # API build
├── docker-compose.yml     # API + MongoDB + LocalStack + Frontend
└── Makefile               # Handy build/test targets
```

## Operating environments

| Ambiente | Localização | Observações gerais |
|----------|-------------|--------------------|
| **dev**  | Máquina do desenvolvedor | Usa Docker Compose com LocalStack (S3/SQS/Cognito) e MongoDB containerizado. Autenticação em modo `local`. |
| **homolog** | AWS (conta de staging) | Executado em infraestrutura gerenciada (ECS Fargate ou EC2). Requer Cognito, S3, SQS e DocumentDB/Mongo Atlas reais. |
| **production** | AWS (conta de produção) | Igual à homolog, porém com escalabilidade, autoscaling, logging e observabilidade reforçados. |

O mecanismo de configuração é único para todos os ambientes: a aplicação carrega defaults, e em seguida mescla um arquivo YAML indicado pela variável `CONFIG_FILE` + variáveis de ambiente. Assim conseguimos segregar comportamentos por ambiente apenas definindo um arquivo específico (ou secret) e exportando `CONFIG_FILE`.

### 1. Requisitos mínimos

- Go **1.24+**
- Node.js **20.19+** (necessário para `vite@7`)
- Docker + Docker Compose
- AWS CLI configurada com credenciais apropriadas (para homolog/produção)
- Chave de criptografia AES-256 (32 bytes base64) para proteção de recibos (`security.encryptionKey`)

### 2. Configuração do ambiente **dev** (local)

1. Copie o template:
   ```bash
   cp config/local_credentials.example.yaml config/local_credentials.yaml
   ```
   Ajuste campos se necessário; os padrões já apontam para serviços locais via Docker.
   Gere uma chave de criptografia (ex.: `openssl rand -base64 32`) e preencha `security.encryptionKey`.

2. Suba os serviços de apoio com Docker Compose:
   ```bash
   docker compose up --build
   ```
   Componentes levantados:
   - **API** (Go) em `http://localhost:8080`
   - **Frontend** (build estático servido por Nginx) em `http://localhost:5173`
   - **MongoDB** containerizado (`mongodb://localhost:27017`)
   - **LocalStack** com S3, SQS e Cognito simulados (script `scripts/localstack/00-bootstrap.sh` cria bucket/fila/pool automaticamente)

3. Executar manualmente (sem Docker) quando necessário:
   ```bash
   # Backend API
   export CONFIG_FILE=config/local_credentials.yaml
   go run ./cmd/api

   # Lambda (build para testes locais)
   GOOS=linux GOARCH=amd64 go build -o bin/transaction_processor ./cmd/lambdas/transaction_processor

   # Frontend (com hot reload)
   cd frontend
   npm install
   npm run dev
   ```
4. Testes:
   ```bash
   go test ./...
   npm run build
   ```

### 3. Configuração do ambiente **homolog** (AWS)

1. **Recursos AWS necessários** (crie via CloudFormation/Terraform ou manualmente):
   - Amazon Cognito User Pool e App Client (`USER_PASSWORD_AUTH` habilitado).
   - Amazon SQS queue (ex.: `finance-transactions-queue`) + DLQ opcional.
   - Amazon S3 bucket para recibos (`finance-control-receipts-hml`).
   - Banco de dados compatível com MongoDB (DocumentDB ou MongoDB Atlas).
   - VPC/Subnets/SG conforme política da empresa.

2. **Credenciais e configuração**:
   - Armazene um YAML com as chaves reais (ex.: `config/hml_credentials.yaml`).
   - Suba o arquivo para **AWS Secrets Manager** ou **Systems Manager Parameter Store**, e faça com que a esteira de deploy materialize o arquivo no container (ou injete variáveis via `APP_PORT`, `MONGO_URI`, `AWS_REGION` etc.).
   - Garanta `auth.mode=cognito` e endpoints reais (`aws.endpoint` vazio).

3. **Deploy da API**:
   - Construa a imagem:
     ```bash
     docker build -t finance-api:latest .
     ```
   - Publique em um registry (ECR):
     ```bash
     aws ecr create-repository --repository-name finance-api --region us-east-1
     docker tag finance-api:latest <ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/finance-api:latest
     docker push <ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/finance-api:latest
     ```
   - Execute em ECS Fargate (recomendado) ou EC2 com Systemd. Configure as variáveis:
     ```
     CONFIG_FILE=/app/config/config.yaml   # se arquivo for montado
     APP_ENVIRONMENT=homolog
     AWS_REGION=us-east-1
     ```
     Monte o arquivo de config via volume secreto, ou converta os campos do YAML em variáveis de ambiente.

4. **Deploy da Lambda (`cmd/lambdas/transaction_processor`)**:
   ```bash
   GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/lambdas/transaction_processor
   zip lambda.zip bootstrap
   aws lambda create-function \
     --function-name finance-transaction-processor-hml \
     --zip-file fileb://lambda.zip \
     --handler bootstrap \
     --runtime provided.al2 \
     --role arn:aws:iam::<ACCOUNT_ID>:role/<LAMBDA_ROLE>
   ```
   Configure variáveis da lambda:
   ```
   CONFIG_FILE=/opt/config/config.yaml   # se montado via Lambda layer/S3
   AWS_REGION=us-east-1
   AUTH_MODE=cognito
   ```
   Conecte a fila SQS como trigger e habilite DLQ.

5. **Frontend**:
   - Build estático:
     ```bash
     cd frontend
     npm install
     npm run build
     ```
   - Faça upload do conteúdo de `frontend/dist` para um bucket S3 público (com CloudFront) ou sirva via container Nginx.
   - Configure `VITE_API_URL` (ou `REACT_APP_API_URL`) apontando para o endpoint da API em homolog (ex.: `https://api-hml.suaempresa.com/api/v1`).

### 3.1 Provisionamento automatizado com Terraform

Uma stack Terraform está disponível em `infra/terraform` para criar todos os recursos AWS (ECR, ECS Fargate, ALB, S3, SQS, Cognito e instância EC2 com MongoDB). Os passos resumidos são:

1. Entre no diretório `infra/terraform` e copie o arquivo de variáveis:
   ```bash
   cd infra/terraform
   cp terraform.tfvars.example terraform.tfvars
   ```
   Ajuste `environment`, `container_image` e demais variáveis conforme o alvo (homolog ou produção).

2. Construa e publique a imagem multi-arquitetura no ECR indicado:
   ```bash
   docker buildx build --platform linux/amd64,linux/arm64 -t finance-api:latest ../../
   aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com
   docker tag finance-api:latest <ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/finance-control-api:latest
   docker push <ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/finance-control-api:latest
   ```

3. Execute:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

4. Após o apply, os outputs `alb_dns_name`, `cognito_user_pool_id` e `mongo_private_ip` serão exibidos. Utilize-os para configurar o frontend e para troubleshooting. O recurso `aws_budgets_budget` já fica configurado com limite padrão de USD 100/mês (ajustável via `monthly_cost_limit`).

> **Custos estimados:**
> - ECS Fargate (0.25 vCPU / 0.5 GB) ~ USD 15/mês.
> - ALB ~ USD 18/mês.
> - EC2 t4g.micro (MongoDB) + EBS 20 GB ~ USD 12/mês.
> - Serviços auxiliares (S3, SQS, Cognito) cobrados sob demanda (estimativa < USD 10/mês para workloads moderados).
> - Total aproximado: ~ USD 55/mês, mantendo folga em relação ao budget de USD 100.

### 4. Configuração do ambiente **production** (AWS)

Mesma topologia de homolog, com os seguintes cuidados adicionais:

- **Segurança**:
  - Restrinja `security.allowedOrigins` aos domínios reais (`https://app.suaempresa.com`).
  - Utilize buckets com criptografia SSE e política de acesso mínimo.
  - Habilite HTTPS em todos os endpoints (API + Frontend via CloudFront/ALB).

- **Observabilidade**:
  - Configure logs do container no CloudWatch.
  - Adicione métricas/alarme (SQS depth, erros 5xx da API, consumo da lambda).

- **Escalabilidade**:
  - Ajuste `desiredCount`/`autoScaling` no ECS ou use EKS.
  - Avalie read replicas/sharding no banco MongoDB.

- **Resiliência**:
  - Ative DLQ para SQS + alarme de mensagens mortas.
  - Considere idempotência reforçada na lambda (armazenar `transactionId` processado).

### 5. Variáveis e arquivos úteis

| Chave | Descrição |
|-------|-----------|
| `CONFIG_FILE` | Caminho para YAML com credenciais (uso obrigatório em homolog/prod via secrets). |
| `app.environment` | Pode ser `development`, `homolog`, `production` (informativo, mas útil em logs). |
| `auth.mode` | `local` (dev) ou `cognito` (homolog/prod). |
| `aws.*` | Região, chaves e endpoints; em ambiente AWS deixe `endpoint` vazio para usar o serviço real. |
| `queue.transactionQueue` | Nome lógico da fila; apontado para uma fila distinta por ambiente. |
| `storage.receiptBucket` | Bucket diferente por ambiente (ex.: `finance-control-receipts-dev|hml|prd`). |
| `security.encryptionKey` | Chave AES-256 em base64 utilizada para criptografar recibos antes do upload ao S3. |

Idealmente:
- Mantenha três arquivos de configuração:
  - `config/local_credentials.yaml` (dev)
  - `config/homolog.yaml`
  - `config/production.yaml`
- Em produção e homolog, esses arquivos devem ser gerados dinamicamente a partir de secrets; evite commitar.

### 6. Deploy automatizado sugerido

1. Esteira CI (GitHub Actions/GitLab):
   - Rodar `go test ./...`, `npm run build`, `npm audit`.
   - Construir imagem Docker (`finance-api`) e enviar ao ECR.
   - Empacotar frontend e publicar no S3/CloudFront.
   - Build e deploy da lambda (via `aws lambda update-function-code`).

2. Configurar jobs separados para homolog e produção com aprovação manual.

3. Garantir que `CONFIG_FILE` e variáveis sensíveis sejam injetadas via Secrets Manager.

### API overview

- `POST /api/v1/auth/login` – Delegates to Cognito or local credentials.
- `GET/POST/PATCH/DELETE /api/v1/accounts`
- `GET/POST/DELETE /api/v1/categories`
- `GET/POST/PATCH /api/v1/transactions`
- `POST /api/v1/transactions/:id/receipt` – Uploads receipt to S3 and returns a presigned URL.
- `GET/POST /api/v1/budgets`
- `GET/POST /api/v1/goals` and `POST /api/v1/goals/:id/progress`
- `GET /api/v1/reports/summary`

Every protected route expects a Bearer token validated via Cognito (or local session when running offline).

### Lambda pipeline

Transactions trigger an SQS message (`finance-transactions-queue`). The Go Lambda consumes the queue, enriches budget metrics, and keeps MongoDB values synchronised. The handler shares the same domain models and repositories as the API, ensuring business invariants stay consistent.

## Environment variables

Key variables honoured by the API:

| Variable | Description |
|----------|-------------|
| `CONFIG_FILE` | Path to a YAML configuration (local or Secrets Manager materialised file). |
| `APP_PORT` (via config file) | HTTP port, defaults to 8080. |
| `AWS_*` fields | Credentials and endpoints for AWS/localstack. |
| `AUTH_MODE` | `cognito` or `local`. |

Refer to `config/local_credentials.example.yaml` for the full schema.

## Frontend integration

The React app looks for `VITE_API_URL` (defaults to `http://localhost:8080/api/v1`) and handles authentication tokens automatically through the shared context provider. Material UI powers a responsive layout, while React Query keeps data fresh across the dashboard sections.

## Contributing

1. Ensure formatting with `gofmt` and `npm run lint` (if you add lint rules) before committing.
2. Keep business code inside the `internal` tree and expose adapters via constructor functions.
3. Extend the documentation in `PROJECT.md` whenever architectural decisions evolve.

Happy hacking! 🚀
