# 🚗 SpotSync — Parking Spot Reservation API

A RESTful backend service for managing parking zones and reservations, built with **Go**, **Echo v5**, **GORM**, and **PostgreSQL (Neon)**. The API supports role-based access control with JWT authentication — **admins** manage zones and view all reservations, while **drivers** reserve spots and manage their own bookings.

> **Live URL:** [https://spotsync.onrender.com](https://spotsync.onrender.com) *(update with your actual deployment URL)*

---

## ✨ Features

- **User Authentication** — Register and login with hashed passwords (bcrypt) and JWT tokens.
- **Role-Based Access Control** — Two roles: `admin` and `driver`.
  - Admins can create, update, and delete parking zones.
  - Admins can view all reservations in the system.
  - Drivers can only create, view, and cancel their own reservations.
- **Parking Zone Management** — Full CRUD for parking zones (admin-only for write operations).
- **Reservation System** — Drivers reserve spots with license plate tracking and capacity enforcement.
- **Concurrency-Safe Capacity Check** — Uses database-level row locking (`SELECT ... FOR UPDATE`) to prevent overbooking.
- **Automatic Status Management** — Reservations have `active`, `cancelled`, and `completed` statuses.
- **Structured JSON Responses** — Consistent success/error response format across all endpoints.

---

## 🛠️ Tech Stack

| Layer          | Technology                                                      |
| -------------- | --------------------------------------------------------------- |
| Language       | Go 1.25                                                         |
| Web Framework  | [Echo v5](https://echo.labstack.com/)                           |
| ORM            | [GORM](https://gorm.io/)                                       |
| Database       | PostgreSQL ([Neon](https://neon.tech/) serverless)              |
| Authentication | JWT ([golang-jwt/jwt v5](https://github.com/golang-jwt/jwt))   |
| Validation     | [go-playground/validator v10](https://github.com/go-playground/validator) |
| Password Hash  | bcrypt (`golang.org/x/crypto`)                                  |

---

## 🏗️ Architecture

SpotSync follows a **clean layered architecture** organized by domain. Each domain is self-contained with its own handler, service, repository, entity, and DTOs.

```
cmd/
└── main.go                    ← Application entry point

internal/
├── auth/                      ← JWT token generation & validation
├── config/                    ← Environment variable loading
├── database/                  ← Database connection & auto-migration
├── httpresponse/              ← Shared success/error response structs
├── middleware/                 ← Auth middleware (JWT verification)
├── server/                    ← HTTP server setup & route registration
└── domain/
    ├── user/                  ← User registration, login, password hashing
    ├── parkingzones/          ← CRUD for parking zones
    └── reservations/          ← Reservation create, list, cancel
```

### How the Layers Interact

```mermaid
flowchart TD
    Client["Client (Postman / Frontend)"]
    MW["Auth Middleware\n(JWT Validation)"]
    H["Handler Layer\n(HTTP request/response)"]
    S["Service Layer\n(Business logic & DTO mapping)"]
    R["Repository Layer\n(Database queries via GORM)"]
    DB["PostgreSQL (Neon)"]

    Client -->|HTTP Request| MW
    MW -->|Set user claims in context| H
    H -->|Call service methods| S
    S -->|Call repository methods| R
    R -->|SQL via GORM| DB
    DB -->|Result| R
    R -->|Entity| S
    S -->|DTO Response| H
    H -->|JSON Response| Client
```

| Layer        | Responsibility                                                              |
| ------------ | --------------------------------------------------------------------------- |
| **Handler**  | Parses HTTP input, validates requests, checks authorization, returns JSON.  |
| **Service**  | Contains business logic, maps entities to response DTOs.                    |
| **Repository** | Executes database operations (CRUD, transactions, row locking).           |
| **Middleware** | Extracts and validates JWT, injects user claims into the request context.  |

---

## 🚀 Setup — Running Locally

### Prerequisites

- **Go** 1.25+ installed → [https://go.dev/dl/](https://go.dev/dl/)
- **PostgreSQL** database (local or cloud, e.g. [Neon](https://neon.tech/))

### 1. Clone the repository

```bash
git clone https://github.com/hisuvo/spotsync.git
cd spotsync
```

### 2. Create the `.env` file

Create a `.env` file in the project root with the following variables:

```env
APP_ENV=development
DB_URL=postgresql://user:password@host:5432/dbname?sslmode=require
PORT=8080
JWT_SECRET=your-secret-key-here
TOKEN_DURATION=240h
```

| Variable         | Required | Default        | Description                                           |
| ---------------- | -------- | -------------- | ----------------------------------------------------- |
| `APP_ENV`        | No       | —              | Set to `development` to allow fallback JWT secret.    |
| `DB_URL`         | **Yes**  | —              | PostgreSQL connection string.                         |
| `PORT`           | No       | `8080`         | Port the server listens on.                           |
| `JWT_SECRET`     | **Yes*** | `dev-secret-key` | Secret key for signing JWT tokens. *Required in production; optional in development mode.* |
| `TOKEN_DURATION` | No       | `24h`          | How long JWT tokens remain valid (Go duration format). |

### 3. Install dependencies & run

```bash
# Download Go modules
go mod download

# Run the server
go run ./cmd/main.go
```

Or use the Makefile:

```bash
make run      # Run the dev server
make build    # Build the binary
```

The server starts at `http://localhost:8080`. Database tables are auto-migrated on startup.

---

## 📡 API Endpoints

Base URL: `/api/v1`

### Authentication

| Method | Endpoint              | Access  | Description               |
| ------ | --------------------- | ------- | ------------------------- |
| POST   | `/api/v1/auth/register` | Public  | Register a new user       |
| POST   | `/api/v1/auth/login`    | Public  | Login and receive a JWT   |

### Parking Zones

| Method | Endpoint              | Access       | Description                    |
| ------ | --------------------- | ------------ | ------------------------------ |
| GET    | `/api/v1/zones`       | Public       | List all parking zones         |
| GET    | `/api/v1/zones/:id`   | Public       | Get a parking zone by ID       |
| POST   | `/api/v1/zones`       | 🔒 Admin     | Create a new parking zone      |
| PUT    | `/api/v1/zones/:id`   | 🔒 Admin     | Update a parking zone          |
| DELETE | `/api/v1/zones/:id`   | 🔒 Admin     | Delete a parking zone          |

### Reservations

| Method | Endpoint                                | Access              | Description                            |
| ------ | --------------------------------------- | ------------------- | -------------------------------------- |
| POST   | `/api/v1/reservations`                  | 🔒 Authenticated    | Create a new reservation               |
| GET    | `/api/v1/reservations/my-reservations`  | 🔒 Authenticated    | Get current user's reservations        |
| DELETE | `/api/v1/reservations/:id`              | 🔒 Authenticated    | Cancel a reservation (own or admin)    |
| GET    | `/api/v1/reservations`                  | 🔒 Admin            | Get all reservations in the system     |

### Access Legend

| Icon             | Meaning                                                    |
| ---------------- | ---------------------------------------------------------- |
| Public           | No authentication required.                                |
| 🔒 Authenticated | Requires a valid JWT token (`Bearer <token>` header).      |
| 🔒 Admin         | Requires JWT token **and** the user must have `admin` role. |

---

## 📝 Example Requests

### Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secret123",
    "role": "driver"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secret123"
  }'
```

### Create a Reservation (Authenticated)

```bash
curl -X POST http://localhost:8080/api/v1/reservations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "zone_id": 1,
    "license_plate": "DHK-1234"
  }'
```

---

## 📄 License

This project is for educational / assignment purposes.
