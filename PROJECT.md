# Decisões Técnicas do Projeto

## Arquitetura Geral

- O backend em Go segue Clean Architecture com camadas explícitas:
  - `internal/domain`: entidades de negócio, interfaces de repositório e portas para serviços externos.
  - `internal/usecase`: orquestração de regras de negócio sem depender de detalhes de infraestrutura.
  - `internal/infrastructure`: adaptadores para MongoDB, AWS (Cognito, S3, SQS), autenticação e logger.
  - `internal/interfaces/http`: entrega HTTP construída com Gin, incluindo middleware de autenticação.
- O serviço HTTP principal está em `cmd/api/main.go`. Toda construção de dependências acontece ali, mantendo os módulos injetáveis.
- A lambda (`cmd/lambdas/transaction_processor`) reutiliza os mesmos repositórios e configurações do backend, garantindo consistência e evitando duplicação de lógica.

## Persistência e Dados

- MongoDB foi escolhido pela flexibilidade na modelagem de transações, contas e metas.
- Cada repositório Mongo cria índices primários/compostos para garantir unicidade e consultas eficientes (`user_id + name`, `user_id + occurred_at`, etc.).
- O cálculo de relatórios agrega dados em memória a partir de coleções de transações, orçamentos e metas, evitando pipelines complexos no Mongo neste primeiro momento.

## Integração com AWS

- SDK v2 da AWS é utilizado com wrappers finos (`internal/infrastructure/aws`) para S3 (upload/presign), SQS (publicação) e Cognito (autenticação). Informações sensíveis (ex.: `security.encryptionKey`) são injetadas via Secrets Manager/Parameter Store.
- A configuração suporta LocalStack através de `aws.useLocalstack` e `aws.endpoint`, permitindo rodar tudo localmente sem credenciais reais.
- SQS recebe eventos de transações para processamento assíncrono. A lambda atualiza os gastos dos orçamentos conforme as mensagens chegam.
- O upload de recibos usa S3 com geração de URL pré-assinada para consumo pelo frontend.

## Autenticação

- Modo padrão utiliza Cognito (`auth.mode = cognito`) com fluxo `USER_PASSWORD_AUTH`.
- Para desenvolvimento offline há um provedor local (`auth.mode = local`) que carrega usuários definidos em `local.authUsers` dentro do YAML de configuração.
- Middleware HTTP valida tokens (Cognito ou local) e sincroniza o usuário na base caso ainda não exista.

## Configuração

- `CONFIG_FILE` aponta para um YAML com credenciais não sensíveis. Em homologação/produção esse arquivo é materializado via Secrets Manager e a chave AES (`security.encryptionKey`) é lida de um secret dedicado.
- `config/local_credentials.example.yaml` documenta todas as chaves e serve como ponto de partida para desenvolvimento.
- O carregamento com Viper permite mesclar defaults, arquivo e variáveis de ambiente.

## Frontend

- Aplicação React com Vite + TypeScript, Material UI, React Router e React Query.
- Armazenamento de tokens feito via `AuthProvider` com `localStorage`, ajustando automaticamente o header `Authorization` do Axios.
- Páginas fundamentais: Dashboard (relatórios), Accounts, Transactions, Budgets e Goals. Cada uma consome as rotas REST do backend.
- Organização em componentes reutilizáveis (`AppLayout`), pages e serviços (`services/api.ts`).

## Docker e Desenvolvimento Local

- `docker-compose.yml` sobe a API, frontend, MongoDB e LocalStack. Script `scripts/localstack/00-bootstrap.sh` cria fila, bucket e estrutura Cognito automaticamente.
- Makefile centraliza comandos (`api-build`, `api-test`, `lambda-build`, `frontend-build`).
- Para rodar localmente sem Docker basta informar `CONFIG_FILE=config/local_credentials.yaml`, subir MongoDB e invocar `go run ./cmd/api`.

## Validação e Qualidade

- Execução obrigatória de `go test ./...` após alterações no backend/lambda.
- `npm run build` garante type-check do frontend.
- `make lambda-build` gera o binário Linux pronto para empacotamento em artefatos de deploy.

## Padrões de Código

- Código em inglês, comentários em português e somente quando agregam contexto arquitetural.
- Dependências externas expostas via interfaces nos pacotes de domínio/porto para facilitar substituição ou testes.
- Nenhum pacote interno deve depender de outro externo na direção errada (ex.: usecases não importam adaptadores de infraestrutura).

## Próximos Passos Sugeridos

- Adicionar automação de provisionamento (Terraform/SAM) para infraestrutura real.
- Implementar testes de integração com MongoDB usando contêiner temporário.
- Expandir a lambda para notificar metas alcançadas e enviar resumos por e-mail/SNS.
