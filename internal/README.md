This `internal` directory organizes the project according to Clean Architecture (Ports & Adapters).

Sub-packages:
- `domain`: Enterprise business rules (entities/value objects). No external deps.
- `usecase`: Application business rules. Orchestrates adapters and builds template functions.
- `adapters`: Interface adapters (GitHub, GoodReads, RSS, Literal, etc.).
- `infra`: Framework & drivers (HTTP clients, config, logging).

Code will be migrated here in small steps to avoid breaking changes.
