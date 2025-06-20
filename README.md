# Sistema de Aluguel de Carros

Este repositório contém um sistema web de aluguel de carros desenvolvido em Go puro, sem o uso de frameworks, como parte da disciplina GAC116 – Programação Web da Universidade Federal de Lavras (UFLA).

## Objetivo

Desenvolver uma aplicação web funcional utilizando a linguagem Go, implementando manualmente o servidor HTTP, roteamento, templates, autenticação, e comunicação com banco de dados para gerenciar clientes, carros e locações.

## Funcionalidades

- Estrutura básica do servidor HTTP em Go
- Renderização de templates HTML com `html/template`
- Organização manual das rotas e handlers
- Gerenciamento de sessões e autenticação (planejado)
- Modelos para clientes, carros e locações (planejado)
- Interface web para cadastro e visualização de dados (planejado)

## Tecnologias utilizadas

- Go 1.24+
- Banco de dados SQLite ou PostgreSQL (a definir)

## Requisitos do sistema

- Go 1.24 ou superior instalado e configurado no PATH

## Regras de commit

Para manter a organização e facilitar a leitura do histórico de alterações, utilize mensagens de commit padronizadas no seguinte formato:

```
<tipo>: <descrição breve da mudança>
```

### Tipos de commit

- **feat**: Adição de nova funcionalidade ao sistema.
  - Exemplo: `feat: adicionar modelo de Cliente`

- **fix**: Correção de bug ou comportamento inesperado.
  - Exemplo: `fix: corrigir erro na criação de Aluguel`

- **docs**: Alterações na documentação (README, comentários, etc.).
  - Exemplo: `docs: atualizar instruções de execução`

- **style**: Alterações de formatação que não afetam a lógica do código (espaços, quebras de linha, etc.).
  - Exemplo: `style: ajustar identação do views.py`

- **refactor**: Refatoração de código (reestruturação interna sem alterar comportamento).
  - Exemplo: `refactor: separar lógica de validação em função auxiliar`

- **test**: Criação ou modificação de testes.
  - Exemplo: `test: adicionar testes para modelo de Carro`

- **chore**: Tarefas de manutenção que não afetam a lógica do sistema (atualização de dependências, configs, etc.).
  - Exemplo: `chore: atualizar .gitignore`

## Autores

- Thallys Henrique Martins
- Gabriel Marcos Lopes
- Bruno de Almeida de Paula

Curso de Ciência da Computação - Universidade Federal de Lavras (UFLA)
