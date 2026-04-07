# Go Auth App

A production-ready, modular authentication REST API for user and admin management, built in Go using Gin, Gorm, and PostgreSQL.

**Key Features:**
- Email-based registration & verification
- JWT access/refresh with rotation, blacklist, and RBAC
- Secure password reset via email (SendGrid)
- Role-based user/admin access
- API endpoints for profile, listing users (admin)
- Modular clean project structure with dotenv-based config
- Ready for local dev or containerized deploy (Docker)
- Built-in Makefile tasks for DB migration & seeding
- Mocks and isolated tests for core logic

---

## 🚀 Features At a Glance

- **Register:** Users sign up, receive verification email
- **Verify Email:** Activate account via emailed token
- **Login:** Obtain access & refresh JWTs (only after verification)
- **Token Rotation:** Each refresh token is single-use, immediately invalid after use
- **Logout:** Blacklist/invalidate refresh token
- **Forgot/Reset Password:** Email reset link, reset with verified token
- **Profile & Admin Listing:** Users see their info; admins can list all users
- **RBAC:** Endpoints restricted by role (user/admin)
- **Security:** bcrypt password hashing
- **Testing:** Extensive mocks, logic isolation

---

## 🏁 Quickstart

### Prerequisites

- Go 1.18+ ([download](https://golang.org/dl/))
- Docker (for Postgres, Redis; optional for development)

### Installation

```sh
git clone https://github.com/your-username/go-auth-app.git
cd go-auth-app
go mod tidy
cp .env.example .env        # Fill in DB, JWT, email config!
```

Edit `.env`:
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_SSLMODE`
- `JWT_SECRET`
- `SENDGRID_API_KEY`, `SENDGRID_EMAIL`
- `ADMIN_EMAIL`, `ADMIN_PASSWORD`
- `REDIS_HOST`, `REDIS_PORT`

---

### Database & Redis: Migrate & Seed

Prepare Docker (recommended for local dev):

```sh
docker compose up -d
```

Or use your own Postgres/Redis setup and update `.env` accordingly.

Use the Makefile for DB tasks:

**Migrate:**
```sh
make migrate
```
- Auto-creates required tables (User, Token, etc.)

**Seed:**
```sh
make seed
```
- Seeds initial admin user (see `.env` for `ADMIN_EMAIL`, `ADMIN_PASSWORD`)

**Quick one-step (migrate + seed):**
```sh
make db-setup
```
- Shortcut: runs both steps above

---

### Launch the Server

```sh
go run main.go
```

Server will listen on: [http://localhost:8080](http://localhost:8080)

You can also deploy using Docker:
```sh
docker build -t go-auth-app .
docker run --env-file .env -p 8080:8080 go-auth-app
```

---

## 📚 API Overview

### Registration & Email Verification

- **POST `/register`**
  ```json
  { "name": "Alice", "email": "alice@email.com", "password": "supersecure" }
  ```
  _Sends verification email—login disabled until activation._

- **GET `/verify-email?token=...`**
  - Visit the emailed link to activate account

- **POST `/resend`**
  ```json
  { "email": "alice@email.com" }
  ```
  _Request a new verification email_

---

### Auth & Tokens

- **POST `/login`**
  ```json
  { "email": "alice@email.com", "password": "supersecure" }
  ```
  _Returns:_
  ```json
  { "access_token": "JWT...", "refresh_token": "..." }
  ```
- **POST `/refresh-token`**
  ```json
  { "refresh_token": "..." }
  ```
  _Returns new tokens; previous refresh token is immediately invalidated._

- **POST `/logout`**
  _Headers:_ `Authorization: Bearer <access_token>`
  ```json
  { "refresh_token": "..." }
  ```

---

### Password Reset

- **POST `/forgot-password`**
  ```json
  { "email": "alice@email.com" }
  ```
  _Sends a password reset email._

- **POST `/reset-password`**
  ```json
  { "token": "<reset_token>", "new_password": "yourNewPassword" }
  ```

---

### User & Admin

- **GET `/profile`**  
  _(Authenticated; token required)_  
  Returns user profile.

- **GET `/users`**  
  _(Admin only; token required)_  
  Returns list of all users.

---

## 🧪 Running Tests

```sh
go test ./tests/...
```
- Covers registration, login, token logic, guards, email flow, password reset
- Uses mocks to isolate business logic

---

## 📂 Project Structure

```
.
├── controllers/     # API route handlers
├── dto/             # Request/response schemas & validation
├── models/          # GORM models (User, Token, etc.)
├── repositories/    # DB access layer
├── services/        # Business logic & integrations
├── middleware/      # JWT, role guards, auth
├── config/          # .env loader, DB, mail, JWT config
├── tests/           # Unit/integration tests, mocks
├── main.go          # Entrypoint
└── dockerfile       # Multi-stage Docker build
```

---

## 🤝 Contributing

1. Fork this repository
2. Create feature branch: `git checkout -b feat/my-feature`
3. Commit & push changes
4. [Open a Pull Request](https://github.com/your-username/go-auth-app/pulls)

---

## 📄 License

MIT License

---
