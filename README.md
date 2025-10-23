# Finanças API

Uma API RESTful para gerenciamento de finanças pessoais, desenvolvida em Go com arquitetura limpa e boas práticas.

## 📋 Funcionalidades

- **Autenticação e Autorização**
  - Login de usuários
  - Ativação de contas
  - Middleware de autenticação JWT

- **Gerenciamento de Usuários**
  - Criação de usuários
  - Ativação de contas

- **Categorias**
  - CRUD completo de categorias
  - Acesso restrito a usuários ativados

- **Transações Financeiras**
  - CRUD completo de transações
  - Filtragem por categoria
  - Acesso restrito a usuários ativados

- **Monitoramento**
  - Health check
  - Métricas expostas via expvar
  - Logs estruturados

## 🛠 Tecnologias

- **Linguagem**: Go
- **Framework Web**: Chi Router
- **Banco de Dados**: PostgreSQL
- **Autenticação**: JWT
- **Logging**: JSON estruturado
- **Monitoramento**: Expvar
- **Configuração**: Environment variables
