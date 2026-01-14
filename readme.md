# greenlight

A Go-based backend service / web API starter kit including CLI, database migrations, and Docker support.

This project provides a modular Go application skeleton useful for building REST APIs or backend microservices with migrations, logging, and Docker setup.

---

## 🚀 Features

- ✅ Go (Golang) application with modular architecture  
- ✅ CLI entrypoints under `cmd/`  
- ✅ Database migrations support  
- ✅ Docker & Docker Compose setup  
- ✅ Structured project layout (`internal/`, `cmd/`, `migrations/`)  
- ✅ Ready for API server or worker service

---

## 📁 Project Structure

```

├── cmd/                   # Command binaries / entrypoints
├── internal/              # Core application logic
├── migrations/            # DB migrations
├── docker/                # Docker related configs
├── .dockerignore
├── docker-compose.yaml
├── go.mod                 # Go module definition
├── go.sum
├── Makefile
└── README.md

````

---

## 🧠 Requirements

Make sure you have installed:

- Go 1.20+  
- Docker & Docker Compose (optional but recommended)  
- A database (PostgreSQL/MySQL depending on migration setup)

---

## 🛠️ Local Setup

1. **Clone the Project**
   ```bash
   git clone https://github.com/PriyanshuSharma23/greenlight.git
   cd greenlight```

2. **Build the app**

   ```bash
   go build ./cmd/...
   ```

3. **Run database migrations**

   > Update the database config/environment variables as needed

   ```bash
   go run ./migrations
   ```

4. **Start the app**

   ```bash
   go run ./cmd/server
   ```

   Replace `server` with the appropriate command if different in `cmd/`

---

## 🐳 Using Docker

1. **Start services**

   ```bash
   docker compose up --build
   ```

2. Access your app at:

   ```
   http://localhost:8080
   ```

   (Adjust port based on your configuration)

---

## ✅ Projects & Usage

This repository is intended to be used as:

✅ Backend service starter kit
✅ API server basis
✅ Microservice or CLI tool foundation

Feel free to extend it with your own handlers, database models, middleware, and configuration!
