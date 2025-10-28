# Decisões Técnicas do Projeto

## Arquitetura Geral

- O backend em Go segue Clean Architecture com camadas explícitas:
  - `src/internal/domain`: entidades de negócio, interfaces de repositório e portas para serviços externos.
  - `src/internal/usecase`: orquestração de regras de negócio sem depender de detalhes de infraestrutura.
  - `src/internal/infrastructure`: adaptadores para MongoDB, AWS (Cognito, S3, SQS), autenticação e logger.
  - `src/internal/adapters/http`: entrega HTTP construída com Gin, incluindo middleware de autenticação.
- O serviço HTTP principal está em `src/cmd/api/main.go`. Toda construção de dependências acontece ali, mantendo os módulos injetáveis.
- A lambda (`src/cmd/lambdas/transaction_processor`) reutiliza os mesmos repositórios e configurações do backend, garantindo consistência e evitando duplicação de lógica.

## Persistência e Dados

- MongoDB foi escolhido pela flexibilidade na modelagem de transações, contas e metas.
- Todas as coleções utilizam **UUIDs** (`uuid.NewString()` em Go) como identificadores. Documentos legados com `ObjectId` devem ser migrados usando `scripts/convert_objectids_to_uuid.js`.
- Cada repositório Mongo cria índices compostos para garantir unicidade e consultas eficientes:
  - `accounts`: `user_id + name` (único), `user_id + created_at`
  - `transactions`: `user_id + occurred_at`, `category_id + occurred_at`
  - `budgets`: `user_id + created_at`, `user_id + category_id + period_start + period_end`
  - `goals`: `user_id + created_at`
- O cálculo de relatórios agrega dados em memória a partir de coleções de transações, orçamentos e metas, evitando pipelines complexos no Mongo neste primeiro momento.
- A coleção `processed_transactions` garante idempotência no processamento assíncrono de transações pela lambda.

## Integração com AWS

- SDK v2 da AWS é utilizado com wrappers finos (`src/internal/infrastructure/aws`) para S3 (upload/presign), SQS (publicação) e Cognito (autenticação). Informações sensíveis (ex.: `security.encryptionKey`) são injetadas via Secrets Manager/Parameter Store.
- A configuração suporta LocalStack através de `aws.useLocalstack` e `aws.endpoint`, permitindo rodar tudo localmente sem credenciais reais.
- SQS recebe eventos de transações para processamento assíncrono. A lambda processa as mensagens de forma idempotente marcando transações processadas na coleção `processed_transactions` antes de atualizar os gastos dos orçamentos, implementando rollback automático em caso de erro.
- O upload de recibos usa S3 com geração de URL pré-assinada para consumo pelo frontend.

## Autenticação

- Modo padrão utiliza Cognito (`auth.mode = cognito`) com fluxo `USER_PASSWORD_AUTH`.
- Para desenvolvimento offline há um provedor local (`auth.mode = local`) que carrega usuários definidos em `local.authUsers` dentro do YAML de configuração.
- Middleware HTTP valida tokens (Cognito ou local) e sincroniza o usuário na base caso ainda não exista.

## Configuração

- `CONFIG_FILE` aponta para um YAML com credenciais não sensíveis. Em homologação/produção esse arquivo é materializado via Secrets Manager e a chave AES (`security.encryptionKey`) é lida de um secret dedicado.
- `src/configs/local_credentials.example.yaml` documenta todas as chaves e serve como ponto de partida para desenvolvimento.
- O carregamento com Viper permite mesclar defaults, arquivo e variáveis de ambiente.

## Frontend

- Aplicação React com Vite + TypeScript, Material UI, React Router e React Query.
- Armazenamento de tokens feito via `AuthProvider` com `localStorage`, ajustando automaticamente o header `Authorization` do Axios.
- Páginas fundamentais: Dashboard (relatórios), Accounts, Transactions, Budgets e Goals. Cada uma consome as rotas REST do backend.
- Organização em componentes reutilizáveis (`AppLayout`), pages e serviços (`services/api.ts`).

## Docker e Desenvolvimento Local

- `docker-compose.yml` sobe a API, frontend, MongoDB, LocalStack e a lambda de processamento. Utiliza `depends_on` com `condition: service_healthy` para garantir inicialização ordenada dos serviços com healthchecks configurados.
- Script `scripts/localstack/00-bootstrap.sh` cria fila (`financial-transactions-queue`), dead-letter queue (`financial-transactions-dlq`), bucket S3 e estrutura Cognito automaticamente.
- A lambda rodando em modo local (`LAMBDA_LOCAL=true`) faz polling contínuo da fila SQS e processa mensagens de forma idempotente através da coleção `processed_transactions`.
- Makefile centraliza comandos (`api-build`, `api-test`, `lambda-build`, `frontend-build`).
- Para rodar localmente sem Docker basta informar `CONFIG_FILE=src/configs/local_credentials.yaml`, subir MongoDB e invocar `go run ./src/cmd/api`.

## Validação e Qualidade

- Execução obrigatória de `go test ./...` após alterações no backend/lambda.
- `npm run build` garante type-check do frontend.
- `make lambda-build` gera o binário Linux pronto para empacotamento em artefatos de deploy.

## Padrões de Código

- Código em inglês, comentários em português e somente quando agregam contexto arquitetural.
- Dependências externas expostas via interfaces nos pacotes de domínio/porto para facilitar substituição ou testes.
- Nenhum pacote interno deve depender de outro externo na direção errada (ex.: usecases não importam adaptadores de infraestrutura).

## Próximos Passos Sugeridos

- Automação de provisionamento (Terraform) está disponível em `infra/terraform` para provisionamento completo da infraestrutura AWS.
- Implementar testes de integração com MongoDB usando contêiner temporário.
- Expandir a lambda para notificar metas alcançadas e enviar resumos por e-mail/SNS.
