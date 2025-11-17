# Akoko: time-tracker

1. Start services:

```bash
docker-compose up --build -d
```

## API (examples)

- Signup: `POST /signup` {"email","password"}
- Login: `POST /login` -> returns JWT
- Create Project: `POST /projects` (auth)
- List Projects: `GET /projects` (auth)
- Create Time Entry: `POST /time-entries` (auth)

Use the Authorization header: `Authorization: Bearer <token>`
