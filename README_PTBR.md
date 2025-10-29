# Plataforma de Controle Financeiro

## Visão Geral

Financial Control é uma plataforma de finanças pessoais construída sobre um backend Go com arquitetura limpa e um frontend React (Vite + Material UI). O sistema persiste dados no MongoDB, integra com Amazon S3 e SQS (opcionalmente via LocalStack para desenvolvimento local), utiliza Amazon Cognito para autenticação e transfere recálculos de orçamento para uma AWS Lambda Go que consome a fila de transações.

## Destaques da Arquitetura

- **Backend com arquitetura limpa** – entidades de domínio, repositórios e casos de uso vivem em `src/internal/`; adaptadores HTTP e preocupações de infraestrutura ficam em seus próprios pacotes.
- **Frontend com tipagem segura** – React + TypeScript com Material UI, React Router e React Query para estado e acesso a dados.
- **Pipeline assíncrono** – cada transação registrada emite uma mensagem SQS; a Lambda atualiza os totais de execução de orçamento.
- **Recibos seguros** – recibos de transação são criptografados com AES-256 e armazenados no S3; URLs pré-assinadas são retornadas para a UI.
- **Ambientes configuráveis** – a configuração é composta por padrões, um arquivo YAML referenciado via `CONFIG_FILE` e variáveis de ambiente.

## Estrutura de Diretórios

```
.
├── src/
│   ├── cmd/                         # Pontos de entrada da aplicação
│   │   ├── api/                     # Ponto de entrada da API HTTP
│   │   └── lambdas/
│   │       └── transaction_processor/   # Handler AWS Lambda
│   ├── internal/                    # Entidades de domínio, casos de uso e adaptadores
│   ├── frontend/                    # Aplicação React (Vite + Material UI)
│   └── configs/                     # Templates de configuração
├── scripts/                         # Utilitários (bootstrap LocalStack, migrações, seeds)
├── infra/terraform/                # Stack de provisionamento AWS opcional
├── docker/                          # Dockerfiles
├── docker-compose.yml
└── Makefile
```

## Estratégia de Identificadores

Todas as coleções agora usam strings UUID (`uuid.NewString()` em Go).

## Requisitos

| Ferramenta | Versão |
|------------|--------|
| Go  | 1.24 ou mais recente |
| Node.js | 20.19 ou mais recente (requerido pelo Vite 7) |
| npm | vem com Node 20 |
| Docker & Docker Compose | versão estável mais recente |
| AWS CLI | configurado para deployments de homolog/produção |
| Chave AES-256 | Chave Base64 de 32 bytes para criptografia de recibos (`security.encryptionKey`) |

> **Secrets**: Fora do desenvolvimento local, armazene valores sensíveis (chave de criptografia, credenciais de banco, chaves AWS) no AWS Secrets Manager ou Systems Manager Parameter Store e injete-os no runtime. Nunca versionize chaves de criptografia.

## Desenvolvimento Local

1. **Configure as credenciais**
   ```bash
   cp src/configs/local_credentials.example.yaml src/configs/local_credentials.yaml
   ```
   Gere uma chave de criptografia (`openssl rand -base64 32`) e preencha `security.encryptionKey`. O arquivo de exemplo já aponta para os serviços criados pelo Docker Compose.

2. **Inicie a stack**
   ```bash
   docker compose up --build
   ```
   Isso inicializa:
   - API em `http://localhost:8080` (`/api/v1/health` para verificação de prontidão)
   - Frontend em `http://localhost:5173`
   - MongoDB em `mongodb://localhost:27017`
   - LocalStack (S3, SQS, Cognito) configurado por `scripts/localstack/00-bootstrap.sh` (cria `financial-transactions-queue` com uma fila de mensagens mortas `financial-transactions-dlq` e uma política de redirecionamento de 5 tentativas)
   - Worker Lambda de transação (serviço `transaction-lambda`) fazendo polling contínuo da fila com `LAMBDA_LOCAL=true`

3. **Execute os serviços manualmente (opcional)**
   ```bash
   # API
   export CONFIG_FILE=src/configs/local_credentials.yaml
   go run ./src/cmd/api

   # Build da Lambda (para testes locais)
   GOOS=linux GOARCH=amd64 go build -o bin/transaction_processor ./src/cmd/lambdas/transaction_processor

   # Frontend com hot reload
   cd src/frontend
   npm install
   npm run dev
   ```

4. **Popule dados**
   ```bash
   # Script completo de seed (cria usuários, contas, categorias, transações, orçamentos e metas)
   mongosh financial-control scripts/seed_complete.js
   ```
   Ajuste o comando `docker exec ... mongosh` se estiver executando MongoDB dentro da stack do compose (o nome do container padrão é `financial-control-mongo-1`).

5. **Testes**
   ```bash
   go test ./...
   npm run build        # Type-checks do frontend
   make lambda-build    # Produz o artefato da Lambda
   ```

## Referência de Seed de Dados

| Script | Propósito |
|--------|-----------|
| `scripts/seed_complete.js` | Script completo de seed que cria dois usuários (vasconcellos e teste), suas contas, categorias, transações, orçamentos e metas para testes e desenvolvimento. |

Todos os scripts são escritos para `mongosh`; canalize-os através de `mongosh <database> < script.js` ou use `docker exec` ao executar MongoDB no Docker.

## Deploy (Homolog / Produção)

1. **Configuração**
   - Forneça um arquivo YAML via `CONFIG_FILE` (secret montado ou renderizado durante o deploy) e complemente com variáveis de ambiente para valores sensíveis.
   - Configure `auth.mode=cognito`, forneça a região AWS real, Client ID do Cognito, bucket S3 e URL da fila.

2. **Build e publicação do backend**
   ```bash
   docker build -t financial-api:latest .
   aws ecr create-repository --repository-name financial-api --region <region>    # uma vez
   docker tag financial-api:latest <ACCOUNT_ID>.dkr.ecr.<region>.amazonaws.com/financial-api:latest
   docker push <ACCOUNT_ID>.dkr.ecr.<region>.amazonaws.com/financial-api:latest
   ```
   Faça o deploy da imagem em ECS Fargate ou EC2 (Systemd). Forneça:
   ```
   CONFIG_FILE=/app/src/configs/config.yaml
   APP_ENVIRONMENT=homolog|production
   AWS_REGION=<region>
   ```

3. **Deploy da Lambda**
   ```bash
   GOOS=linux GOARCH=amd64 go build -o bootstrap ./src/cmd/lambdas/transaction_processor
   zip lambda.zip bootstrap
   aws lambda update-function-code \
     --function-name financial-transaction-processor \
     --zip-file fileb://lambda.zip
   ```
   Configure as mesmas variáveis de ambiente (`CONFIG_FILE`, `AWS_REGION`, `AUTH_MODE`) e conecte a fila SQS como trigger. Habilite uma DLQ para resiliência.

4. **Frontend**
   ```bash
   cd frontend
   npm install
   npm run build
   ```
Faça o deploy do conteúdo de `src/frontend/dist` para S3 + CloudFront ou outra solução de hospedagem estática. Configure `VITE_API_URL` (ex.: `https://api.company.com/api/v1`) no ambiente antes de fazer o build.

## Endpoints da API

- `POST /api/v1/auth/login`
- `GET/POST/PATCH/DELETE /api/v1/accounts`
- `GET/POST/DELETE /api/v1/categories`
- `GET/POST/PATCH /api/v1/transactions`
- `POST /api/v1/transactions/:id/receipt`
- `GET/POST /api/v1/budgets`
- `GET/POST /api/v1/goals`
- `POST /api/v1/goals/:id/progress`
- `GET /api/v1/reports/summary`

Endpoints `GET` para contas, transações, orçamentos e metas aceitam parâmetros opcionais de query `limit` e `offset` (`limit` padrão é 100, limitado a 200; `offset` padrão é 0) para suportar paginação em datasets grandes.

### Variáveis de Ambiente Comuns

| Variável | Notas |
|----------|-------|
| `CONFIG_FILE` | Caminho para o config YAML. |
| `APP_ENVIRONMENT` | `development`, `homolog`, ou `production` (usado para logging/métricas). |
| `AUTH_MODE` | `local` para dev; `cognito` em ambientes gerenciados. |
| `security.encryptionKey` | Chave AES-256 Base64 (deve vir de um secret em homolog/produção). |

Consulte `src/configs/local_credentials.example.yaml` para o schema completo.

## Automação Terraform (Opcional)

O módulo em `infra/terraform` provisiona AWS App Runner para a API, DocumentDB Serverless para MongoDB, S3, SQS (com DLQ) e Cognito. Início rápido:

```bash
cd infra/terraform
cp terraform.tfvars.example terraform.tfvars   # customize ambiente, URI da imagem, budgets, etc.

terraform init
terraform plan
terraform apply
```

Os outputs incluem a URL do App Runner, endpoint do DocumentDB, identificadores do Cognito e nome do bucket S3. Revise os custos antes de aplicar (aproximadamente USD 45-50/mês para o tamanho padrão). Consulte `infra/terraform/README_INFRA.md` para documentação detalhada da infraestrutura.

## Targets Úteis do Make

```bash
make api-build         # go build ./src/cmd/api
make api-test          # go test ./...
make lambda-build      # build binary da Lambda (linux/amd64)
make frontend-build    # npm install && npm run build
make docker-up         # inicia todos os serviços com docker-compose
make docker-down       # para todos os serviços
make docker-logs       # visualiza logs de todos os serviços
make fmt               # formata código Go
```

## Convenções e Leitura Adicional

- Mantenha a lógica de domínio dentro de `src/internal/usecase` e `src/internal/domain`; adaptadores devem permanecer finos e testáveis.
- Comentários devem estar em português quando necessário para contexto; código permanece em inglês.
- Documente mudanças arquiteturais ou de configuração em `PROJECT.md`.
- Antes de abrir um pull request, execute `go test ./...`, `make lambda-build` e `npm run build`.

Para uma análise mais profunda das decisões de design e trabalho técnico aberto, consulte:

- `PROJECT.md` – decisões arquiteturais e convenções.
- `BUGS_AND_LIMITATIONS.md` – issues conhecidos e débito técnico.
- `IMPROVEMENTS.md` – backlog de melhorias.

Happy hacking! 🚀

