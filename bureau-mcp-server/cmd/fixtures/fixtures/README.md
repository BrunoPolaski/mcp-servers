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
