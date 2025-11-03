# Financas - Sistema de Gestão Financeira Pessoal

# Sobre o Projeto
O Financas é uma aplicação web completa para gestão financeira pessoal, desenvolvida em Go. O sistema permite que usuários controlem suas finanças através de categorização de transações, estabelecimento de metas financeiras, acompanhamento de progresso e geração de relatórios detalhados. A aplicação oferece uma solução robusta e segura para ajudar usuários a alcançarem seus objetivos financeiros.


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

- **Relatórios Financeiros**
  - Relatórios completos
  - Acesso restrito a usuários ativados

- **Objetivos Financeiros**
  - CRUD completo de objetivos financeiros
  - Ajuda no planejamento de objetivos estabelecidos
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
