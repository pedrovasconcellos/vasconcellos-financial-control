# Plano de Melhorias Estratégicas

Documento elaborado para orientar evolução contínua do projeto.

## Toolchain e Builds

- **Imagem de produção**: considerar estágio final `scratch` ou `distroless/static` após auditar dependências, reduzindo superfície de ataque. Avaliar uso de `nonroot` para execução.

## Segurança e Configuração

- **Arquivos de credenciais**: substituir materialização em disco (`src/configs/*.yaml`) por leitura direta de AWS Secrets Manager/Parameter Store em homolog/produção. Incluir mecanismo de fallback seguro.
- **CORS**: restringir `security.allowedOrigins` no ambiente produtivo; adicionar verificação automática (ex.: erro se `*` for usado fora de dev).

## API e Domínio

- **Consistência transacional**: mover ajustes de saldo e gravação de transações para transações MongoDB (sessões) ou adotar padrão outbox/event sourcing.

## Lambda e Pipeline Assíncrono

- **Observabilidade**: instrumentar lambda com métricas (tempo de processamento, mensagens processadas) e logs estruturados.
- **Infra as Code**: documentar templates Terraform/CloudFormation para fila, lambda e integrações, garantindo reprodutibilidade.

## Front-end

- **Validação de formulários**: implementar validação cliente (React Hook Form/Yup) para contas, transações, orçamentos e metas, com feedback visual consistente.
- **Renovação de tokens**: implementar fluxo de refresh (quando Cognito habilitar) ou logout automático próximo ao `expiresIn`.
- **Code-splitting**: responder ao warning de chunk > 500 kB (usar `React.lazy`, `manualChunks` no Vite) para melhorar performance.
- **Testes de UI**: adicionar suite de testes (Cypress/Playwright) para fluxos críticos.

## Observabilidade e DevOps

- **Logs e métricas**: integrar OpenTelemetry ou middleware de métricas (Prometheus) na API, registrar traços por request e contadores de chamadas por use case.
- **Pipelines CI/CD**: montar pipelines com estágios de build/test/lint para backend, frontend e lambda, incluindo scans de segurança (Snyk/Dependabot).
- **Kubernetes readiness**: preparar manifests Helm ou Kustomize (deployments, ConfigMaps, Secrets, HPA) para futura orquestração em EKS.

## Documentação

- **Runbooks**: criar runbook operacional (procedimentos de deploy, rollback, rotação de credenciais).
- **SLIs/SLOs**: definir indicadores de disponibilidade e latência; documentar objetivos para produção.

Este roadmap deve ser revisado trimestralmente, priorizando itens conforme impacto/dor identificados em produção ou homologação.
