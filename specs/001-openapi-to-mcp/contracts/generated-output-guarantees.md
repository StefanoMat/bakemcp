# Garantias do projeto Node gerado

**Feature**: 001-openapi-to-mcp | **Date**: 2025-02-06

## Objetivo

Garantir que o projeto Node gerado pelo CLI tenha **estrutura sólida**, **dependências corretas** e **funcione** (instalação e execução) sem depender de “testar na mão”. Tudo verificado por automação.

---

## 1. Estrutura sólida

| Garantia | Como |
|----------|------|
| **Template versionado** | O esqueleto do projeto Node (package.json base, entry script, convenções) fica em `templates/node-fastmcp/` no repositório do CLI. Qualquer mudança é versionada e revisada. |
| **Estrutura mínima obrigatória** | O gerador sempre produz: `package.json` válido (JSON parseável), script de entrada (ex.: `index.ts` ou `index.js`), e nenhum arquivo obrigatório faltando. Contract tests validam presença e formato. |
| **Convenções fixas** | Nomes de arquivos, `scripts.start`, `type: "module"` (ou equivalente) são definidos no contrato e nos testes; o gerador não inventa variações não testadas. |

**Testes**: Em `tests/contract/` (ou equivalente), após gerar para um OpenAPI fixture: (1) listar arquivos gerados e checar que todos os obrigatórios existem; (2) parse de `package.json` e checagem de campos obrigatórios; (3) entry script existe e é legível.

---

## 2. Libs corretas e instaláveis

| Garantia | Como |
|----------|------|
| **Dependências fixas no template** | No template, `package.json` usa versões fixas (ou ranges mínimos documentados) para `fastmcp` e qualquer outra dependência. Sem versões “floating” não testadas. |
| **Apenas dependências necessárias** | O gerador só inclui o que o contrato exige (ex.: fastmcp); não adiciona libs extras por padrão. Lista de deps documentada no contrato do projeto gerado. |
| **npm install bem-sucedido** | Em ambiente de teste (CI/local), após gerar o projeto, rodar `npm install` no diretório gerado. Teste falha se houver erro (rede pode falhar; nesse caso retry ou skip com motivo claro). |

**Testes**: (1) Contract test: `package.json` contém `fastmcp` (e outras deps do contrato) com versão válida. (2) Integration test: em dir gerado, `npm install` termina com exit 0.

---

## 3. Projeto funcionando (run + MCP)

| Garantia | Como |
|----------|------|
| **npm start sobe o processo** | No mesmo dir gerado, após `npm install`, rodar `npm start` (ou o comando definido em `scripts.start`). O processo deve iniciar sem crash imediato; pode usar timeout curto (ex.: 5s) e depois encerrar o processo. Teste falha se o processo terminar com erro antes do timeout. |
| **MCP expõe as tools** | Onde for viável (ex.: com um cliente MCP em processo ou script), integration test: conectar ao MCP gerado, listar tools e invocar pelo menos uma com argumentos válidos. Confirma que o servidor responde ao protocolo MCP e que as tools correspondem às operações do OpenAPI usado. |
| **Determinismo** | Mesmo OpenAPI de entrada → mesmo conjunto de arquivos e de tools. Permite golden/snapshot tests (opcional): comparar `package.json` e lista de tools com um baseline aprovado. |

**Testes**: (1) Integration test: `CLI(fixture OpenAPI) → dir` → `npm install` → `npm start` (com timeout) → sucesso. (2) Integration test (quando houver cliente MCP em teste): list tools + invoke one → sucesso. (3) Opcional: snapshot do `package.json` e dos nomes das tools para um OpenAPI fixture fixo.

---

## 4. Requisitos de ambiente documentados

| Garantia | Como |
|----------|------|
| **Node/npm mínimos** | Documentar no quickstart e, se possível, no README gerado: versão mínima de Node (ex.: 18+) e npm (ex.: 9+). O template não usa recursos de versões mais novas sem documentar. |
| **Compatibilidade** | Se o template usar ESM, `"type": "module"` (ou equivalente) deve estar explícito e testado; o mesmo para TypeScript se o output for `.ts`. |

---

## 5. Resumo: o que falhar = build/CI falha

- **Contract tests**: estrutura e `package.json` válidos; dependências obrigatórias presentes.
- **Integration tests**: gerar → `npm install` → `npm start` (e, quando possível, list/invoke tools) com sucesso.

Assim garantimos que o projeto Node gerado tenha estrutura sólida, libs corretas e esteja funcionando, com verificação automática em cada mudança relevante.
