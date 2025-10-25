# Catálogo de Bugs e Limitações Conhecidas

Documento preparado por um arquiteto de software para orientar correções futuras realizadas por desenvolvedores humanos ou agentes de IA.

---

## Plataforma e Ferramentas

### Node.js desatualizado em relação ao Vite
- **Contexto**: `frontend/package.json` exige `vite@7.1.12`, que por sua vez demanda Node.js `>= 20.19.0`.  
- **Sintoma**: Builds locais exibem warning: `Vite requires Node.js version 20.19+ or 22.12+`.  
- **Risco**: Ambiente com Node anterior pode falhar em features futuras do Vite e causa inconsistências em CI/CD.  
- **Ação sugerida**: Padronizar toolchain para Node 20.19+ (ou 22 LTS) e registrar no README/CI.

### Script do LocalStack assume `jq`
- **Contexto**: `scripts/localstack/00-bootstrap.sh` usa `jq` para extrair IDs.  
- **Risco**: Se a imagem `localstack/localstack` for trocada por uma sem `jq`, o bootstrap quebra silenciosamente.  
- **Ação sugerida**: Garantir instalação explícita ou substituir por comandos AWS nativos (`--query`).

---

## Back-end (Go)

### Tokens de desenvolvimento sem expiração
- **Contexto**: `internal/infrastructure/auth/local_auth.go` mantém tokens em `map[string]LocalAuthUser`.  
- **Problema**: Nenhum TTL/cleanup → mapa cresce indefinidamente em sessões repetidas.  
- **Impacto**: Vazamento de memória em ambientes de QA e comportamento irrealista vs. Cognito.  
- **Correção**: Implementar expiração simples (timestamp + GC) ou usar cache com TTL.

### Mensagens de erro do Cognito expostas ao cliente
- **Contexto**: `internal/interfaces/http/handler/auth_handler.go` propaga `err.Error()` direto.  
- **Risco**: Vazamento de detalhes internos (ex.: stack ou mensagens AWS) e mensagens pouco amigáveis.  
- **Correção**: Mapear erros para respostas sanitizadas (401 com texto padrão).

### Importação de recibos totalmente em memória
- **Contexto**: `transaction_handler.go` lê o upload com `io.ReadAll`.  
- **Problema**: Falta limite de tamanho/streaming → uploads grandes travam API e consomem muita RAM.  
- **Mitigação**: Validar `file.Size`, impor limites (ex.: 5 MB) e fazer streaming direto para S3.

### Falta de controle transacional nas atualizações de saldo
- **Contexto**: `TransactionUseCase.RecordTransaction` ajusta saldo via `AccountRepository.AdjustBalance` e depois grava transação.  
- **Risco**: Em caso de crash entre as duas operações, saldo e transação ficam inconsistentes.  
- **Correção**: Introduzir sessão/transaction no MongoDB ou compensação outbox.

### Lambda atualiza orçamento somando `payload.Amount` absoluto
- **Contexto**: `src/lambdas/transaction_processor/main.go` trata toda mensagem como despesa.  
- **Problema**: Não diferencia receitas vs. despesas → metas de orçamento infladas para transações de crédito.  
- **Correção**: Incluir campo de tipo (income/expense) e ajustar lógica para receitas não subirem o gasto.

### CORS excessivamente permissivo
- **Contexto**: `internal/interfaces/http/router.go` permite `AllowOrigins: ["*"]`.  
- **Risco**: APIs expostas para qualquer origem, facilitando ataques CSRF ou uso indevido em produção.  
- **Mitigação**: Restringir origens via config.

### Falta de paginação nos endpoints de listagem
- **Contexto**: `ListAccounts`, `ListTransactions`, `ListBudgets`, etc. retornam todos registros.  
- **Problema**: Para bases grandes, requests ficam lentos e pesados.  
- **Correção**: Introduzir parâmetros de paginação/filtragem e índices apropriados.

---

## Front-end (React)

### Falta de validação de formulário nas telas principais
- **Contexto**: Páginas como `AccountsPage` e `TransactionsPage` só confiam no backend para validação.  
- **Impacto**: UX fraca (usuário não vê erros imediatamente) e risco de submissões inválidas.  
- **Correção**: Adicionar validações lado cliente (React Hook Form/Yup ou validações manuais).

### Armazenamento de tokens no `localStorage` sem renovação
- **Contexto**: `AuthProvider` apenas persiste tokens, sem refresh flow.  
- **Problema**: Sessões expiram após `expiresIn` e o usuário só percebe ao fazer uma chamada; não há estratégia de renovação.  
- **Sugestão**: Implementar refresh automático ou revogar/cerrar sessão antes da expiração.

### Tabelas sem tratamento de carregamento/erro granulado
- **Contexto**: telas mostram `Loading...` como string na tabela ou alerta genérico.  
- **Problema**: Para UX profissional, é necessário skeleton/feedback adequado.  
- **Correção**: Usar `Skeleton` do MUI e mensagens específicas por recurso.

---

## Pipeline Assíncrono

### Ausência de DLQ e reprocessamento
- **Contexto**: Lambda que processa SQS não trata falhas além de logar.  
- **Problema**: Mensagem que falhar permanece na fila e é reprocessada indefinidamente sem observabilidade.  
- **Correção**: Configurar DLQ + alarmes e adicionar idempotência na lambda.

### Falta de idempotência na lambda
- **Contexto**: Update de orçamento soma `payload.Amount` sem verificar se a transação já foi processada.  
- **Risco**: Reexecuções (SQS redelivery) duplicam gastos.  
- **Correção**: Registrar `transactionId` processados (ex.: coleção `processed_events`) antes de aplicar o delta.

---

## Observabilidade e Operação

### Falta de tracing e métricas
- **Contexto**: Projeto só usa logger, sem instrumentação.  
- **Consequência**: Dificuldades na análise de performance e problemas em produção.  
- **Sugestão**: Adotar OpenTelemetry ou middleware de métricas (Prometheus).

### Logger wrappers não utilizados
- **Contexto**: `structuredLogger` em `cmd/api/main.go` foi declarado, mas não utilizado.  
- **Problema**: Código morto que confunde manutenção.  
- **Correção**: Remover tipo ou utilizá-lo de fato.

---

## Segurança e Configuração

### Configuração sensível em arquivo local sem criptografia
- **Contexto**: `config/local_credentials.example.yaml` incentiva uso de arquivo plano com credenciais.  
- **Risco**: Em ambientes reais, arquivo pode vazar.  
- **Correção**: Integrar com AWS Secrets Manager/Parameter Store ou Vault, documentando processo.

### Modo local com credenciais padrão expostas
- **Contexto**: Usuário demo (`demo@local.dev`) com senha conhecida.  
- **Problema**: Facilita invasões em ambientes que replicam defaults.  
- **Correção**: Exigir substituição via variável de ambiente ou script pós-instalação.

---

## Backlog Técnico Prioritário
1. Atualizar toolchain Node e CI para Node 20.19+.  
2. Introduzir limites/streaming no upload de recibos.  
3. Implementar expiração de tokens e sanitização de respostas de autenticação.  
4. Revisar lambda para suportar idempotência e diferenciar tipos de transação.  
5. Restringir CORS e validar origens configuráveis.  
6. Adicionar paginação nos endpoints críticos.  
7. Documentar/automatizar provisionamento seguro de credenciais e infraestrutura AWS real.

Este catálogo deve ser mantido atualizado a cada release para garantir clareza sobre débitos técnicos e bugs em aberto.
