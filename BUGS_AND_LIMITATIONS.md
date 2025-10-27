# Catálogo de Bugs e Limitações Conhecidas

Documento preparado por um arquiteto de software para orientar correções futuras realizadas por desenvolvedores humanos ou agentes de IA.

---

## Back-end (Go)

### Falta de controle transacional nas atualizações de saldo
- **Contexto**: `TransactionUseCase.RecordTransaction` ajusta saldo via `AccountRepository.AdjustBalance` e depois grava transação.  
- **Risco**: Em caso de crash entre as duas operações, saldo e transação ficam inconsistentes.  
- **Correção**: Introduzir sessão/transaction no MongoDB ou compensação outbox.

### CORS excessivamente permissivo
- **Contexto**: `internal/interfaces/http/router.go` permite `AllowOrigins: ["*"]`.  
- **Risco**: APIs expostas para qualquer origem, facilitando ataques CSRF ou uso indevido em produção.  
- **Mitigação**: Restringir origens via config.

---

## Front-end (React)

### Armazenamento de tokens no `localStorage` sem renovação
- **Contexto**: `AuthProvider` apenas persiste tokens, sem refresh flow.  
- **Problema**: Sessões expiram após `expiresIn` e o usuário só percebe ao fazer uma chamada; não há estratégia de renovação.  
- **Sugestão**: Implementar refresh automático ou revogar/cerrar sessão antes da expiração.

---

## Observabilidade e Operação

### Falta de tracing e métricas
- **Contexto**: Projeto só usa logger, sem instrumentação.  
- **Consequência**: Dificuldades na análise de performance e problemas em produção.  
- **Sugestão**: Adotar OpenTelemetry ou middleware de métricas (Prometheus).

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
1. Restringir CORS e validar origens configuráveis (evitar `*` fora de dev).  
2. Documentar/automatizar provisionamento seguro de credenciais e infraestrutura AWS real.

Este catálogo deve ser mantido atualizado a cada release para garantir clareza sobre débitos técnicos e bugs em aberto.
