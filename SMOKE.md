# Smoke

Avec Compose déjà up (`docker compose up -d`).

Tests unitaires / build: `cd backend && go test ./...` · `cd frontend && npm run build`

## Checklist

- [ ] UI http://localhost:3847 → **API ok**
- [ ] Kanban charge les candidatures
- [ ] Nueva → créer une candidature → revient au kanban
- [ ] Chat (Nueva) → « stats » / « lista » répond
- [ ] Drag & drop change le status

## Script (API)

```powershell
.\scripts\smoke.ps1
```

Lit `APP_API_KEY` depuis `.env`. Base API: `http://localhost:8088` (ou `$env:SMOKE_API`).
