// Package handler contains HTTP handlers for the Financial Control API
// @title Financial Control API
// @version 1.0
// @description API REST completa para controle financeiro pessoal
// @description Esta API permite gerenciar contas, categorias, transações, orçamentos e metas financeiras.
// @description Todos os endpoints protegidos requerem autenticação via Bearer token (JWT).
// @description
// @description ## Recursos Disponíveis
// @description - **Autenticação**: Login com Cognito ou modo local
// @description - **Contas**: Gerenciamento de contas bancárias (checking, savings, credit, cash)
// @description - **Categorias**: Organização de despesas e receitas
// @description - **Transações**: Registro e histórico de movimentações financeiras
// @description - **Orçamentos**: Controle de limites por categoria e período
// @description - **Metas**: Acompanhamento de objetivos financeiros
// @description - **Relatórios**: Sumários e análises financeiras
// @description
// @description ## Autenticação
// @description Para acessar endpoints protegidos, você precisa:
// @description 1. Fazer login em `POST /api/v1/auth/login`
// @description 2. Copiar o `accessToken` da resposta
// @description 3. Adicionar header: `Authorization: Bearer {accessToken}`
// @description
// @description ## Códigos de Status HTTP
// @description - **200 OK**: Requisição bem-sucedida
// @description - **201 Created**: Recurso criado com sucesso
// @description - **400 Bad Request**: Dados inválidos
// @description - **401 Unauthorized**: Token ausente ou inválido
// @description - **404 Not Found**: Recurso não encontrado
// @description - **500 Internal Server Error**: Erro no servidor
// @description
// @termsOfService https://github.com/pedrovasconcellos/vasconcellos-financial-control/blob/main/README.md
// @contact.name API Support
// @contact.email support@financialcontrol.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Use format: Bearer {your_token}
// @tag.name auth
// @tag.description Endpoints de autenticação
// @tag.name accounts
// @tag.description Gerenciamento de contas bancárias
// @tag.name categories
// @tag.description Categorias de transações
// @tag.name transactions
// @tag.description Registro e consulta de transações
// @tag.name budgets
// @tag.description Controle de orçamentos
// @tag.name goals
// @tag.description Metas financeiras
// @tag.name reports
// @tag.description Relatórios e análises
// @tag.name health
// @tag.description Status da API

package handler
