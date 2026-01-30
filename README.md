# Go_Student_Registry🎓

A robust, high-performance RESTful API for student management built with **Go (Golang)**, **SQLite**, and **Redis**.

This project demonstrates a production-ready architecture using the **Go Standard Library** (Go 1.22+) enhanced with a **Redis caching layer** for ultra-fast data retrieval. It features structured logging, a "Cache-Aside" pattern, and a fully containerized environment.

## 🚀 Features

* **RESTful CRUD Operations:** Full lifecycle management for student records.
* **High-Speed Caching:** Implements **Redis** to cache database queries, significantly reducing latency for `GET` requests.
* **Standard Library Routing:** Utilizes Go 1.22+ `http.ServeMux` for native, framework-free routing.
* **Persistent & Portable Storage:** Uses **SQLite** with Write-Ahead Logging (WAL) enabled for safe concurrent access.
* **Containerized Environment:** Fully Dockerized setup with **Docker Compose** for one-command deployment.
* **Graceful Shutdown:** Cleanly closes SQLite connections, Redis clients, and HTTP listeners on exit.

## 🛠️ Tech Stack

* **Language:** Go (1.22+)
* **Database:** SQLite 3
* **Cache:** Redis 7 (Alpine)
* **Containerization:** Docker & Docker Compose
* **Router:** `net/http` (Standard Lib)
* **Logging:** `log/slog` (Standard Lib)

## 📂 Project Structure

```bash
go_stuAPI/
├── cmd/
│   └── stuAPI/        # Main application entry point
├── internal/
│   ├── config/        # Configuration loading logic
│   ├── http/
│   │   └── handlers/  # Handlers with Redis-logic integration
│   └── storage/
│       └── sqlite/    # Optimized SQLite interaction layer
├── storage/           # Local folder for persistent .db files
├── Dockerfile         # Multi-stage build for CGO/SQLite
├── docker-compose.yml # Orchestration for App and Redis
└── main.go            # Application Entry point

```

## 🔌 API Endpoints

| Method | Endpoint | Description | Cache Logic |
| --- | --- | --- | --- |
| `POST` | `/api/students` | Create a student | Invalidates List Cache |
| `GET` | `/api/students/{id}` | Get student by ID | **Cache Hit/Miss** |
| `GET` | `/api/students/` | List all students | **Cache Hit/Miss** |
| `PUT` | `/api/students/{id}` | Update student | Updates Hash & Invalidates List |
| `DELETE` | `/api/students/{id}` | Remove student | Deletes Cache Keys |

## ⚙️ Getting Started

### Prerequisites

* [Docker](https://www.docker.com/products/docker-desktop/) installed.
* *Alternatively:* Go 1.22+ and a local Redis instance.

### Installation & Deployment (Docker)

The fastest way to get started is using Docker Compose. This will spin up the API on **port 8082** and a Redis instance on **port 6379**.

1. **Clone and Enter:**
```bash
git clone https://github.com/Sarthak-D97/go_stuAPI.git
cd go_stuAPI

```


2. **Run with Docker Compose:**
```bash
docker-compose up --build

```



### Testing the API

Once the containers are running, the API is available at `http://localhost:8082`.

**Example: Create a Student**

```bash
curl -X POST http://localhost:8082/api/students \
  -H "Content-Type: application/json" \
  -d '{"name": "Sarthak", "email": "sarthak@example.com", "age": 25}'

```

**Example: Get Student (Check logs to see Redis cache hit)**

```bash
curl http://localhost:8082/api/students/1

```

## 📄 License

This project is open-source and available under the [MIT License](https://opensource.org/licenses/MIT).

