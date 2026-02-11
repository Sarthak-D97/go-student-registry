
# Go_Student_Registry 🎓

A robust, high-performance RESTful API for student and video management built with **Go (Golang)**, **Gin Framework**, **PostgreSQL**, and **Redis**.

This project demonstrates a production-ready architecture using the **MVC Pattern** (Model-View-Controller). It features **JWT Authentication**, structured logging, a "Cache-Aside" pattern for performance, and fully interactive API documentation via **Swagger**.

## 🚀 Features

* **Gin Web Framework:** High-performance HTTP web framework for routing and middleware.
* **Swagger Documentation:** Interactive API docs available at `/docs/index.html`.
* **JWT Authentication:** Secure access to private routes using JSON Web Tokens.
* **Dual-Entity Management:** Full CRUD operations for **Students** and **Videos**.
* **PostgreSQL Database:** Reliable, relational storage for all persistent data.
* **High-Speed Caching:** Implements **Redis** to cache database queries, significantly reducing latency.
* **Containerized Environment:** Fully Dockerized setup with **Docker Compose**.
* **Graceful Shutdown:** Cleanly closes DB connections and HTTP listeners on exit.

## 🛠️ Tech Stack

* **Language:** Go (1.21+)
* **Framework:** [Gin Gonic](https://github.com/gin-gonic/gin)
* **Database:** PostgreSQL 16
* **Cache:** Redis 7 (Alpine)
* **Documentation:** Swagger (Swaggo)
* **Auth:** JWT (JSON Web Tokens)
* **Containerization:** Docker & Docker Compose

## 📂 Project Structure

```bash
go_stuAPI/
├── cmd/
│   └── stuAPI/        # Main application entry point
├── controller/        # HTTP Handlers (Gin Context)
├── service/           # Business Logic Layer
├── repository/        # Database Access Layer (GORM/SQL)
├── internal/
│   ├── config/        # Configuration loading
│   └── platform/      # DB & Redis connections
├── docs/              # Swagger generated files
├── middlewares/       # Auth & Logging middleware
├── Dockerfile         # Multi-stage build
├── docker-compose.yml # Orchestration for App, Postgres, Redis
└── main.go            # Application Entry point

```

## 🔌 API Endpoints

### 📖 Documentation

The full interactive documentation is available in your browser:
👉 **[http://localhost:5000/docs/index.html](https://www.google.com/search?q=http://localhost:5000/docs/index.html)**

### Public Routes

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/login` | Authenticate user & get **JWT Token** |
| `GET` | `/docs/*` | Swagger UI Access |

### Protected Routes (Requires `Authorization: Bearer <token>`)

**Students**

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/api/students` | Create a student |
| `GET` | `/api/students/{id}` | Get student by ID (**Cached**) |
| `GET` | `/api/students/` | List all students (**Cached**) |
| `PUT` | `/api/students/{id}` | Update student details |
| `DELETE` | `/api/students/{id}` | Remove student |

**Videos**

| Method | Endpoint | Description |
| --- | --- | --- |
| `GET` | `/api/videos` | List all videos |
| `POST` | `/api/videos` | Add a new video |
| `PUT` | `/api/videos/{id}` | Update video metadata |
| `DELETE` | `/api/videos/{id}` | Delete a video |

## ⚙️ Getting Started

### Prerequisites

* [Docker Desktop](https://www.docker.com/products/docker-desktop/) installed.

### Installation & Deployment (Docker)

The fastest way to get started is using Docker Compose. This will spin up the API on **port 5000**, PostgreSQL on **port 5432**, and Redis on **port 6379**.

1. **Clone the Repository:**

```bash
git clone [https://github.com/Sarthak-D97/go_stuAPI.git](https://github.com/Sarthak-D97/go_stuAPI.git)
cd go_stuAPI

```

2. **Run with Docker Compose:**

```bash
docker-compose up --build

```

*(Note: The first run might take a moment to initialize the Database schema).*

### Testing the API

**1. Access Swagger UI**
Open your browser and navigate to:
[http://localhost:5000/docs/index.html](https://www.google.com/search?q=http://localhost:5000/docs/index.html)

**2. Login to get a Token**
Use the `Try it out` button on the `/login` endpoint.

* **Username:** `student_user`
* **Password:** `secret123`
* *Copy the returned "token" string.*

**3. Authorize in Swagger**

* Click the **Authorize** button at the top of the Swagger page.
* Enter `Bearer <YOUR_TOKEN_HERE>`.
* Now you can test all the `/api/students` and `/api/videos` endpoints!

## 📄 License

This project is open-source and available under the [MIT License](https://opensource.org/licenses/MIT).

```

```