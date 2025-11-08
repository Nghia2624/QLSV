# 🎓 QLSVGO - Hệ Thống Quản Lý Sinh Viên

[![Go Version](https://img.shields.io/badge/Go-1.24.1-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com/)

> Hệ thống quản lý sinh viên hiện đại được xây dựng theo kiến trúc **CQRS** và **Event-Driven Architecture** sử dụng Golang, PostgreSQL, MongoDB, và Apache Kafka.

---

## 📋 Mục Lục

- [Tổng Quan](#-tổng-quan)
- [Kiến Trúc Hệ Thống](#-kiến-trúc-hệ-thống)
- [Tính Năng](#-tính-năng)
- [Công Nghệ Sử Dụng](#-công-nghệ-sử-dụng)
- [Cấu Trúc Dự Án](#-cấu-trúc-dự-án)
- [Cài Đặt & Chạy](#-cài-đặt--chạy)
- [API Documentation](#-api-documentation)
- [Phân Quyền](#-phân-quyền)
- [Cấu Hình](#-cấu-hình)
- [Monitoring](#-monitoring)
- [Bảo Mật](#-bảo-mật)
- [Roadmap](#-roadmap)

---

## 🎯 Tổng Quan

**QLSVGO** là hệ thống quản lý sinh viên được thiết kế theo mô hình **CQRS (Command Query Responsibility Segregation)** và **Event-Driven Architecture**, cho phép tách biệt hoàn toàn các thao tác ghi (Write) và đọc (Read) để đạt được hiệu suất và khả năng mở rộng tối ưu.

### ✨ Đặc Điểm Nổi Bật

- 🔄 **CQRS Pattern**: Tách biệt Command và Query services
- 📡 **Event-Driven**: Sử dụng Kafka để đồng bộ dữ liệu giữa các services
- 🗄️ **Dual Database**: PostgreSQL cho write operations, MongoDB cho read operations
- 🔐 **JWT Authentication**: Xác thực và phân quyền dựa trên JWT
- 🐳 **Docker Ready**: Dễ dàng triển khai với Docker Compose
- 🏗️ **Clean Architecture**: Cấu trúc code rõ ràng, dễ bảo trì

---

## 🏗️ Kiến Trúc Hệ Thống

### Sơ Đồ Tổng Quan

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client Apps   │    │   Admin Panel   │    │   Mobile App    │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │      API Gateway          │
                    └─────────────┬─────────────┘
                                  │
          ┌───────────────────────┼───────────────────────┐
          │                       │                       │
    ┌─────▼─────┐           ┌─────▼─────┐           ┌─────▼─────┐
    │ Command   │           │  Query    │           │  Auth     │
    │ Service   │           │ Service   │           │ Service   │
    │ (Port 8080)│          │ (Port 8081)│          │ (Port 8080)│
    └─────┬─────┘           └─────┬─────┘           └─────┬─────┘
          │                       │                       │
          │                       │                       │
    ┌─────▼─────┐           ┌─────▼─────┐           ┌─────▼─────┐
    │PostgreSQL │           │ MongoDB   │           │  JWT      │
    │(Write DB) │           │(Read DB)  │           │  Auth     │
    └─────┬─────┘           └─────┬─────┘           └───────────┘
          │                       │
          │                       │
          └───────────┬───────────┘
                      │
              ┌───────▼───────┐
              │    Kafka      │
              │ Message Broker│
              └───────────────┘
```

### Luồng Dữ Liệu

#### 1. Write Operations (Commands)
```
Client → Command Service → PostgreSQL → Kafka Event → Query Service → MongoDB
```

**Chi tiết:**
1. Client gửi request tạo/cập nhật/xóa dữ liệu đến **Command Service**
2. Command Service xác thực JWT và phân quyền
3. Command Service lưu dữ liệu vào **PostgreSQL** (Write Database)
4. Command Service publish **event** vào Kafka topic tương ứng
5. Query Service consume event từ Kafka
6. Query Service đồng bộ dữ liệu sang **MongoDB** (Read Database)

#### 2. Read Operations (Queries)
```
Client → Query Service → MongoDB → Response
```

**Chi tiết:**
1. Client gửi request đọc dữ liệu đến **Query Service**
2. Query Service xác thực JWT và phân quyền
3. Query Service đọc dữ liệu từ **MongoDB** (tối ưu cho read)
4. Trả về response cho client

#### 3. Authentication Flow
```
Client → Login Request → Auth Service → JWT Token → Client
Client → API Request + JWT → Service → Validate JWT → Process Request
```

### Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │ HTTP Handler│  │ Middleware  │  │   Router    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │ Use Cases   │  │ Event Bus   │  │ Validators  │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                     Domain Layer                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Models    │  │  Events     │  │ Interfaces  │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                  Infrastructure Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │ PostgreSQL  │  │  MongoDB    │  │   Kafka     │        │
│  │ Repository  │  │ Repository  │  │  Service    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
```

---

## 🚀 Tính Năng

### Core Features
- ✅ **Quản lý Sinh viên**: CRUD operations cho sinh viên
- ✅ **Quản lý Giảng viên**: CRUD operations cho giảng viên
- ✅ **Quản lý Môn học**: CRUD operations cho môn học
- ✅ **Quản lý Lớp học**: CRUD operations cho lớp học
- ✅ **Đăng ký Môn học**: Sinh viên đăng ký các lớp học
- ✅ **Xác thực & Phân quyền**: JWT-based authentication với 3 roles (Admin, Teacher, Student)

### Technical Features
- ✅ **Event-Driven Synchronization**: Tự động đồng bộ dữ liệu qua Kafka
- ✅ **Health Checks**: Monitoring endpoints cho tất cả services
- ✅ **Database Migrations**: Tự động migrate khi khởi động
- ✅ **Structured Logging**: Logging với timestamps và levels

---

## 🛠️ Công Nghệ Sử Dụng

### Backend
- **Golang 1.24.1**: Ngôn ngữ lập trình chính
- **Standard Library**: HTTP server, context, time

### Databases
- **PostgreSQL 16**: Write database (ACID compliance)
- **MongoDB 7**: Read database (optimized for queries)

### Message Broker
- **Apache Kafka 7.5.0**: Event streaming platform
- **Zookeeper**: Kafka cluster coordination

### Authentication
- **JWT (golang-jwt/jwt/v5)**: JSON Web Token authentication

### Infrastructure
- **Docker & Docker Compose**: Containerization
- **Kafka UI**: Web interface for Kafka management

### Dependencies
- `github.com/lib/pq`: PostgreSQL driver
- `go.mongodb.org/mongo-driver`: MongoDB driver
- `github.com/segmentio/kafka-go`: Kafka client

---

## 📁 Cấu Trúc Dự Án

```
QLSVGO/
├── cmd/
│   ├── command-service/      # Command Service (Write operations)
│   │   ├── main.go
│   │   └── Dockerfile
│   └── query-service/         # Query Service (Read operations)
│       ├── main.go
│       └── Dockerfile
│
├── internal/
│   ├── domain/               # Domain Layer
│   │   ├── model/            # Domain models (Student, Teacher, Course, etc.)
│   │   └── event/            # Event definitions
│   │
│   ├── usecase/              # Application Layer - Business logic
│   │   ├── student.go
│   │   ├── teacher.go
│   │   ├── course.go
│   │   ├── class.go
│   │   ├── registration.go
│   │   └── user.go
│   │
│   ├── repository/           # Repository interfaces
│   │   ├── student_repository.go
│   │   ├── teacher_repository.go
│   │   └── ...
│   │
│   ├── infrastructure/       # Infrastructure Layer
│   │   ├── postgres/         # PostgreSQL implementations
│   │   ├── mongo/            # MongoDB implementations
│   │   ├── kafka/            # Kafka service
│   │   ├── jwt/              # JWT utilities
│   │   └── config/           # Configuration management
│   │
│   ├── handler/              # Presentation Layer
│   │   ├── http/             # HTTP handlers
│   │   ├── middleware/       # JWT, Role middleware
│   │   └── router/           # Router configuration
│   │
│   └── mapping/              # Data mapping utilities
│       └── unified_mapping.go
│
├── pkg/                      # Shared packages
│   ├── errors/               # Error definitions
│   ├── logger/               # Logging utilities
│   └── utils/                # Utility functions
│
├── migrations/               # Database migrations
│   └── 001_init.sql
│
├── docker-compose.yaml       # Docker Compose configuration
├── go.mod                    # Go dependencies
├── go.sum                    # Go checksums
├── README.md                 # This file
└── tongquan.md               # Detailed documentation (Vietnamese)
```

---

## 🚀 Cài Đặt & Chạy

### Prerequisites

Đảm bảo bạn đã cài đặt:
- **Docker** (version 20.10+)
- **Docker Compose** (version 2.0+)
- **Git**

### Quick Start

1. **Clone repository**
```bash
git clone <repository-url>
cd QLSVGO
```

2. **Build và khởi động tất cả services**
```bash
docker-compose up -d --build
```

3. **Kiểm tra trạng thái services**
```bash
docker-compose ps
```

4. **Xem logs**
```bash
# Xem tất cả logs
docker-compose logs -f

# Xem logs của service cụ thể
docker-compose logs -f command-service
docker-compose logs -f query-service
```

### Services & Ports

Sau khi khởi động, các services sẽ chạy trên các ports sau:

| Service | Port | URL | Description |
|---------|------|-----|-------------|
| Command Service | 8080 | http://localhost:8080 | Write operations API |
| Query Service | 8081 | http://localhost:8081 | Read operations API |
| Kafka UI | 8082 | http://localhost:8082 | Kafka management UI |
| PostgreSQL | 5433 | localhost:5433 | Write database |
| MongoDB | 27018 | localhost:27018 | Read database |
| Kafka | 9093 | localhost:9093 | Message broker |

### Health Checks

Kiểm tra trạng thái services:
```bash
# Command Service
curl http://localhost:8080/health

# Query Service
curl http://localhost:8081/health
```

### Dừng Services

```bash
# Dừng tất cả services
docker-compose down

# Dừng và xóa volumes (xóa dữ liệu)
docker-compose down -v
```

---

## 📚 API Documentation

### Base URLs
- **Command Service**: `http://localhost:8080`
- **Query Service**: `http://localhost:8081`

### Authentication APIs

#### Đăng ký User
```http
POST /api/register
Content-Type: application/json

{
  "username": "admin",
  "email": "admin@example.com",
  "password": "password123",
  "role": "admin"
}
```

#### Đăng nhập
```http
POST /api/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin"
  }
}
```

### Command APIs (Write Operations)

> **Lưu ý**: Tất cả endpoints dưới đây yêu cầu JWT token trong header `Authorization: Bearer <token>` và role `admin`.

#### Student Commands

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/students` | Tạo sinh viên mới |
| PUT | `/api/students/{id}` | Cập nhật sinh viên |
| DELETE | `/api/students/{id}` | Xóa sinh viên |

**Example - Create Student:**
```http
POST /api/students
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Nguyễn Văn A",
  "email": "nguyenvana@example.com",
  "phone": "0123456789",
  "gender": "Nam",
  "dob": "2000-01-01",
  "address": "Hà Nội"
}
```

#### Teacher Commands

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/teachers` | Tạo giảng viên mới |
| PUT | `/api/teachers/{id}` | Cập nhật giảng viên |
| DELETE | `/api/teachers/{id}` | Xóa giảng viên |

#### Course Commands

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/courses` | Tạo môn học mới |
| PUT | `/api/courses/{id}` | Cập nhật môn học |
| DELETE | `/api/courses/{id}` | Xóa môn học |

#### Class Commands

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/classes` | Tạo lớp học mới |
| PUT | `/api/classes/{id}` | Cập nhật lớp học |
| DELETE | `/api/classes/{id}` | Xóa lớp học |

#### Registration Commands

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/registrations` | Đăng ký môn học |
| PUT | `/api/registrations/{id}` | Cập nhật đăng ký |
| DELETE | `/api/registrations/{id}` | Hủy đăng ký |

### Query APIs (Read Operations)

> **Lưu ý**: Tất cả endpoints dưới đây yêu cầu JWT token trong header `Authorization: Bearer <token>`.

#### User Queries

| Method | Endpoint | Description | Role Required |
|--------|----------|-------------|---------------|
| GET | `/api/users` | Lấy danh sách tất cả users | admin |
| GET | `/api/users/{id}` | Lấy thông tin user theo ID | admin |

#### Student Queries

| Method | Endpoint | Description | Role Required |
|--------|----------|-------------|---------------|
| GET | `/api/students` | Lấy danh sách tất cả sinh viên | admin, teacher |
| GET | `/api/students/{id}` | Lấy thông tin sinh viên theo ID | admin, teacher |

#### Teacher Queries

| Method | Endpoint | Description | Role Required |
|--------|----------|-------------|---------------|
| GET | `/api/teachers` | Lấy danh sách tất cả giảng viên | admin |
| GET | `/api/teachers/{id}` | Lấy thông tin giảng viên theo ID | admin |

#### Course Queries

| Method | Endpoint | Description | Role Required |
|--------|----------|-------------|---------------|
| GET | `/api/courses` | Lấy danh sách tất cả môn học | admin, teacher, student |
| GET | `/api/courses/{id}` | Lấy thông tin môn học theo ID | admin, teacher, student |

#### Class Queries

| Method | Endpoint | Description | Role Required |
|--------|----------|-------------|---------------|
| GET | `/api/classes` | Lấy danh sách tất cả lớp học | admin, teacher, student |
| GET | `/api/classes/{id}` | Lấy thông tin lớp học theo ID | admin, teacher, student |

#### Registration Queries

| Method | Endpoint | Description | Role Required |
|--------|----------|-------------|---------------|
| GET | `/api/registrations` | Lấy danh sách tất cả đăng ký | admin, teacher |
| GET | `/api/registrations/{id}` | Lấy thông tin đăng ký theo ID | admin, teacher |

---

## 🔐 Phân Quyền

### Roles

1. **admin**: Quyền quản trị toàn bộ hệ thống
   - Full CRUD cho tất cả entities
   - Quản lý users
   - Xem tất cả dữ liệu

2. **teacher**: Quyền giảng viên
   - Read: Sinh viên, Môn học, Lớp học, Đăng ký
   - Không có quyền write

3. **student**: Quyền sinh viên
   - Read: Môn học, Lớp học
   - Không có quyền write

### Permission Matrix

| Resource | admin | teacher | student |
|----------|-------|---------|---------|
| User Management | ✅ Full | ❌ | ❌ |
| Student Management | ✅ Full | ✅ Read | ❌ |
| Teacher Management | ✅ Full | ❌ | ❌ |
| Course Management | ✅ Full | ✅ Read | ✅ Read |
| Class Management | ✅ Full | ✅ Read | ✅ Read |
| Registration Management | ✅ Full | ✅ Read | ❌ |

**Legend:**
- ✅ Full: Create, Read, Update, Delete
- ✅ Read: Read only
- ❌: No access

---

## ⚙️ Cấu Hình

### Environment Variables

Các biến môi trường được cấu hình trong `docker-compose.yaml`:

```bash
# Database Configuration
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=qlsv
POSTGRES_PASSWORD=qlsv123
POSTGRES_DB=qlsvdb

MONGO_URI=mongodb://qlsv:qlsv123@mongo:27017/qlsvdb?authSource=admin
MONGO_DB=qlsvdb

# Kafka Configuration
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC=qlsv-events
KAFKA_CONSUMER_GROUP=qlsv-query-group

# Kafka Topics
KAFKA_STUDENT_TOPIC=student-events
KAFKA_TEACHER_TOPIC=teacher-events
KAFKA_COURSE_TOPIC=course-events
KAFKA_CLASS_TOPIC=class-events
KAFKA_REGISTRATION_TOPIC=registration-events
KAFKA_USER_TOPIC=user-events

# JWT Configuration
JWT_SECRET_KEY=qlsv-super-secret-key-change-in-production
JWT_EXPIRATION=24h

# Application
APP_ENV=production
LOG_LEVEL=INFO
```

> **⚠️ Lưu ý**: Trong môi trường production, hãy thay đổi các giá trị mặc định, đặc biệt là `JWT_SECRET_KEY` và database passwords.

### Kafka Topics

Hệ thống sử dụng các Kafka topics sau:
- `student-events`: Events liên quan đến Student
- `teacher-events`: Events liên quan đến Teacher
- `course-events`: Events liên quan đến Course
- `class-events`: Events liên quan đến Class
- `registration-events`: Events liên quan đến Registration
- `user-events`: Events liên quan đến User

Mỗi entity có 3 loại event:
- `{Entity}Created`: Khi tạo mới
- `{Entity}Updated`: Khi cập nhật
- `{Entity}Deleted`: Khi xóa

---

## 📊 Monitoring

### Health Checks

Tất cả services có health check endpoint:
```bash
# Command Service
curl http://localhost:8080/health

# Query Service
curl http://localhost:8081/health
```

### Kafka UI

Truy cập Kafka UI để quản lý và monitor Kafka:
- **URL**: http://localhost:8082
- **Features**:
  - Xem danh sách topics
  - Xem messages trong topics
  - Monitor consumer groups
  - Xem broker metrics

### Logging

Hệ thống sử dụng structured logging với:
- Timestamps
- Log levels: INFO, WARN, ERROR
- Kafka consumer logs với event processing status

Xem logs:
```bash
# Tất cả services
docker-compose logs -f

# Service cụ thể
docker-compose logs -f command-service
docker-compose logs -f query-service
docker-compose logs -f kafka
```

---

## 🛡️ Bảo Mật

### Authentication
- **JWT-based authentication**: Sử dụng JSON Web Tokens
- **Token expiration**: 24 hours (có thể cấu hình)
- **Secure token storage**: Tokens được lưu ở client side

### Authorization
- **Role-based access control (RBAC)**: 3 roles (admin, teacher, student)
- **Middleware validation**: Mỗi endpoint được bảo vệ bởi middleware
- **Principle of least privilege**: Users chỉ có quyền tối thiểu cần thiết

### Data Protection
- **Input validation**: Validation cho tất cả inputs
- **SQL injection prevention**: Sử dụng parameterized queries
- **Password security**: ⚠️ Hiện tại password lưu plain text (nên hash trong production)

> **⚠️ Security Recommendations**:
> - Implement bcrypt cho password hashing
> - Sử dụng HTTPS trong production
> - Implement rate limiting
> - Add CORS configuration
> - Implement audit logging

---

## 🔮 Roadmap

### Planned Features

1. **Security Enhancements**
   - [ ] Password hashing với bcrypt
   - [ ] API rate limiting
   - [ ] CORS configuration
   - [ ] Audit logging

2. **Performance**
   - [ ] Redis caching layer
   - [ ] Database connection pooling optimization
   - [ ] Query optimization

3. **Monitoring & Observability**
   - [ ] Prometheus metrics
   - [ ] Grafana dashboards
   - [ ] Distributed tracing (Jaeger/Zipkin)

4. **Documentation**
   - [ ] Swagger/OpenAPI documentation
   - [ ] API examples và tutorials
   - [ ] Architecture decision records (ADRs)

5. **Testing**
   - [ ] Unit tests
   - [ ] Integration tests
   - [ ] End-to-end tests
   - [ ] Load testing

6. **CI/CD**
   - [ ] GitHub Actions workflow
   - [ ] Automated testing pipeline
   - [ ] Automated deployment

7. **Scalability**
   - [ ] Horizontal scaling support
   - [ ] Database sharding
   - [ ] Load balancing
   - [ ] Message queue optimization

---

## 📝 Ghi Chú

- Đảm bảo Docker và Docker Compose đã được cài đặt trước khi chạy
- Các services sẽ tự động migrate database khi khởi động
- Kafka topics sẽ được tự động tạo khi có events đầu tiên
- Trong môi trường production, hãy thay đổi tất cả credentials mặc định

---

## 📞 Liên Hệ & Hỗ Trợ

- **Tác giả**: Đỗ Ngọc Nghĩa
- **Email**: dnghia9119@gmail.com
- **Website**: [dnnghia.vercel.app](https://dnnghia.vercel.app)

---

## 📄 License

MIT License - Xem file LICENSE để biết thêm chi tiết.

---

**Made with ❤️ using Go, PostgreSQL, MongoDB, and Kafka**
