# Fixtures

JSON fixtures for the Registration Validation MCP server. Each file is a JSON array whose objects use **snake_case** keys matching database columns (GORM default naming).

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

Os registros de id 1 a 10 são o conjunto preliminar, curado à mão. Os de id 11 a
40 são os 30 cenários de avaliação (10 simples, 10 intermediários, 10 complexos)
gerados por `generate_cenarios.py`, na raiz de `mcp-servers/`, que declara cada
cenário uma única vez e o deriva para as quatro fontes. Edite o gerador, não os
arquivos, e rode `generate_cadastro.py` em seguida.
