# Fixtures

JSON fixtures for the bureau MCP server. Each file is a JSON array whose objects use **snake_case** keys matching database columns (GORM default naming).

- Timestamps use RFC3339 (e.g., `2024-01-01T00:00:00Z`).
- `deleted_at` is `null` for active rows.
- Foreign keys reference IDs in the corresponding fixture files.
- `person_data_sources.json` represents the join table for `persons` ↔ `data_sources`.

## Fixture runner

The fixture runner loads JSON files in dependency order and inserts them using GORM.

```bash
go run ./cmd/fixtures -dir ./fixtures
```

To truncate tables before loading:

```bash
go run ./cmd/fixtures -dir ./fixtures -truncate
```

## Cadastro

Os arquivos do bloco cadastral (`personal_informations.json`, `addresses.json`,
`person_addresses.json`) são gerados por `generate_cadastro.py`, na raiz de
`mcp-servers/`, que monta a visão específica desta fonte a partir do cadastro do
banco. Edite o gerador, não os arquivos: as divergências entre as fontes são
propositais e precisam continuar coerentes com as demais.

## Cenários

Os registros de id 1 a 30 são os 30 cenários de avaliação de crédito gerados por
`generate_cenarios.py`, na raiz de `mcp-servers/`, que declara cada cenário uma
única vez (no dicionário `CENARIOS`) e o deriva para as quatro fontes. Eles se
distribuem em faixas `bom` (1-3), `ruim` (4-6) e `complexo` (7-30), na
proporção 10% / 10% / 80% — deliberada, para imitar uma carteira real: casos
limpos e casos obviamente reprováveis são raros, e a maioria exige ler as
quatro fontes juntas, que frequentemente discordam entre si. Edite o gerador,
não os arquivos, e rode `generate_cadastro.py` em seguida.
