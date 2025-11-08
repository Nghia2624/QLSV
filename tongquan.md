# QLSVGO - Hệ Thống Quản Lý Sinh Viên

## 📋 Tổng Quan

QLSVGO là hệ thống quản lý sinh viên được xây dựng theo kiến trúc **CQRS (Command Query Responsibility Segregation)** và **Event-Driven Architecture** sử dụng:

- **Backend**: Golang
- **Database**: PostgreSQL (Write-side) + MongoDB (Read-side)
- **Message Broker**: Apache Kafka
- **Containerization**: Docker & Docker Compose
- **Authentication**: JWT (JSON Web Token)
- **Authorization**: Role-based Access Control

## 🏗️ Kiến Trúc Tổng Thể

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

## 🔄 Luồng Đi Dữ Liệu

### 1. Write Operations (Commands)
```
Client → Command Service → PostgreSQL → Kafka Event → Query Service → MongoDB
```

**Chi tiết:**
1. Client gửi request tạo/cập nhật/xóa dữ liệu đến Command Service
2. Command Service xác thực JWT và phân quyền
3. Command Service lưu dữ liệu vào PostgreSQL
4. Command Service publish event vào Kafka topic tương ứng
5. Query Service consume event từ Kafka
6. Query Service đồng bộ dữ liệu sang MongoDB

### 2. Read Operations (Queries)
```
Client → Query Service → MongoDB → Response
```

**Chi tiết:**
1. Client gửi request đọc dữ liệu đến Query Service
2. Query Service xác thực JWT và phân quyền
3. Query Service đọc dữ liệu từ MongoDB
4. Trả về response cho client

### 3. Authentication Flow
```
Client → Login Request → Auth Service → JWT Token → Client
Client → API Request + JWT → Service → Validate JWT → Process Request
```

## 📊 Cấu Trúc Dữ Liệu

### Entities

#### 1. User
```go
type User struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email"`
    Role     string `json:"role"`
}
```

#### 2. Student
```go
type Student struct {
    ID          string `json:"id"`
    StudentCode string `json:"student_code"`
    Name        string `json:"name"`
    Email       string `json:"email"`
    Phone       string `json:"phone"`
    Gender      string `json:"gender"`
    BirthDate   string `json:"birth_date"`
    Address     string `json:"address"`
    ClassID     string `json:"class_id"`
}
```

#### 3. Teacher
```go
type Teacher struct {
    ID         string `json:"id"`
    TeacherCode string `json:"teacher_code"`
    Name       string `json:"name"`
    Email      string `json:"email"`
    Phone      string `json:"phone"`
    Gender     string `json:"gender"`
    BirthDate  string `json:"birth_date"`
    Address    string `json:"address"`
    Department string `json:"department"`
}
```

#### 4. Course
```go
type Course struct {
    ID          string `json:"id"`
    CourseCode  string `json:"course_code"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Credits     int    `json:"credits"`
    TeacherID   string `json:"teacher_id"`
}
```

#### 5. Class
```go
type Class struct {
    ID        string `json:"id"`
    ClassCode string `json:"class_code"`
    Name      string `json:"name"`
    CourseID  string `json:"course_id"`
    TeacherID string `json:"teacher_id"`
    Capacity  int    `json:"capacity"`
}
```

#### 6. Registration
```go
type Registration struct {
    ID        string `json:"id"`
    StudentID string `json:"student_id"`
    ClassID   string `json:"class_id"`
    Status    string `json:"status"`
    RegisterDate string `json:"register_date"`
}
```

## 🚀 API Endpoints

### Authentication APIs
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/register` | Đăng ký user mới | ❌ |
| POST | `/api/login` | Đăng nhập | ❌ |

### Command APIs (Write Operations)
*Tất cả endpoints dưới đây yêu cầu JWT token và role "admin"*

#### Student Commands
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/students` | Tạo sinh viên mới |
| PUT | `/api/students/{id}` | Cập nhật sinh viên |
| DELETE | `/api/students/{id}` | Xóa sinh viên |

#### Teacher Commands
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/teachers` | Tạo giáo viên mới |
| PUT | `/api/teachers/{id}` | Cập nhật giáo viên |
| DELETE | `/api/teachers/{id}` | Xóa giáo viên |

#### Course Commands
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/courses` | Tạo khóa học mới |
| PUT | `/api/courses/{id}` | Cập nhật khóa học |
| DELETE | `/api/courses/{id}` | Xóa khóa học |

#### Class Commands
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/classes` | Tạo lớp học mới |
| PUT | `/api/classes/{id}` | Cập nhật lớp học |
| DELETE | `/api/classes/{id}` | Xóa lớp học |

#### Registration Commands
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/registrations` | Đăng ký khóa học |
| PUT | `/api/registrations/{id}` | Cập nhật đăng ký |
| DELETE | `/api/registrations/{id}` | Hủy đăng ký |

### Query APIs (Read Operations)
*Tất cả endpoints dưới đây yêu cầu JWT token*

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
| GET | `/api/teachers` | Lấy danh sách tất cả giáo viên | admin |
| GET | `/api/teachers/{id}` | Lấy thông tin giáo viên theo ID | admin |

#### Course Queries
| Method | Endpoint | Description | Role Required |
|--------|----------|-------------|---------------|
| GET | `/api/courses` | Lấy danh sách tất cả khóa học | admin, teacher, student |
| GET | `/api/courses/{id}` | Lấy thông tin khóa học theo ID | admin, teacher, student |

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

### Health Check APIs
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Kiểm tra trạng thái service |

## 🔐 Bảng Phân Quyền

### Roles
1. **admin**: Quyền quản trị toàn bộ hệ thống
2. **teacher**: Quyền giáo viên (xem thông tin sinh viên, khóa học, lớp học)
3. **student**: Quyền sinh viên (xem thông tin khóa học, lớp học)

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

## 📡 Kafka Topics & Events

### Topics
1. **user-events**: Events liên quan đến User
2. **student-events**: Events liên quan đến Student
3. **teacher-events**: Events liên quan đến Teacher
4. **course-events**: Events liên quan đến Course
5. **class-events**: Events liên quan đến Class
6. **registration-events**: Events liên quan đến Registration

### Event Types
Mỗi entity có 3 loại event:
- **{Entity}Created**: Khi tạo mới
- **{Entity}Updated**: Khi cập nhật
- **{Entity}Deleted**: Khi xóa

### Event Structure
```go
type Event struct {
    ID        string    `json:"id"`
    Type      string    `json:"type"`
    Entity    string    `json:"entity"`
    Data      string    `json:"data"`
    Timestamp time.Time `json:"timestamp"`
}
```

## 🏛️ Kiến Trúc Clean Architecture

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

## 🔧 Cấu Hình Hệ Thống

### Environment Variables
```bash
# Database
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=qlsv
POSTGRES_PASSWORD=qlsv123
POSTGRES_DB=qlsvdb

MONGO_URI=mongodb://qlsv:qlsv123@mongo:27017/qlsvdb?authSource=admin
MONGO_DB=qlsvdb

# Kafka
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC=qlsv-events
KAFKA_CONSUMER_GROUP=qlsv-query-group

# JWT
JWT_SECRET_KEY=qlsv-super-secret-key-change-in-production
JWT_EXPIRATION=24h
```

### Ports
- **Command Service**: 8080
- **Query Service**: 8081
- **Kafka UI**: 8082
- **PostgreSQL**: 5433
- **MongoDB**: 27018
- **Kafka**: 9093

## 🚀 Deployment

### Prerequisites
- Docker
- Docker Compose

### Quick Start
```bash
# Clone repository
git clone <repository-url>
cd QLSVGO

# Build and start services
docker-compose up -d --build

# Check services status
docker-compose ps

# View logs
docker-compose logs -f
```

### Services
1. **PostgreSQL**: Database chính (Write-side)
2. **MongoDB**: Database đọc (Read-side)
3. **Zookeeper**: Quản lý Kafka cluster
4. **Kafka**: Message broker
5. **Command Service**: Xử lý write operations
6. **Query Service**: Xử lý read operations
7. **Kafka UI**: Giao diện quản lý Kafka

## 📈 Monitoring & Logging

### Health Checks
- Tất cả services có health check endpoint `/health`
- Docker Compose health checks cho databases và Kafka

### Logging
- Structured logging với timestamps
- Log levels: INFO, WARN, ERROR
- Kafka consumer logs với event processing status

### Monitoring
- Kafka UI: http://localhost:8082
- Service health: http://localhost:8080/health, http://localhost:8081/health

## 🔄 Event-Driven Synchronization

### Flow Diagram
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ Command     │    │   Kafka     │    │ Query       │
│ Service     │───▶│   Topic     │───▶│ Service     │
│             │    │             │    │             │
│ PostgreSQL  │    │ Event Store │    │ MongoDB     │
└─────────────┘    └─────────────┘    └─────────────┘
```

### Benefits
- **Scalability**: Read và Write operations độc lập
- **Performance**: MongoDB optimized cho read operations
- **Reliability**: Event-driven ensures eventual consistency
- **Flexibility**: Easy to add new consumers for different use cases

## 🛡️ Security

### Authentication
- JWT-based authentication
- Token expiration: 24 hours
- Secure token storage

### Authorization
- Role-based access control (RBAC)
- Middleware validation for each endpoint
- Principle of least privilege

### Data Protection
- Password hashing (plain text in demo - should be hashed in production)
- Input validation and sanitization
- SQL injection prevention through parameterized queries

## 🔮 Future Enhancements

### Planned Features
1. **Password Hashing**: Implement bcrypt for password security
2. **API Rate Limiting**: Prevent abuse
3. **Audit Logging**: Track all changes
4. **Caching**: Redis for frequently accessed data
5. **Metrics**: Prometheus + Grafana monitoring
6. **API Documentation**: Swagger/OpenAPI
7. **Testing**: Unit tests, integration tests
8. **CI/CD**: Automated deployment pipeline

### Scalability Considerations
1. **Horizontal Scaling**: Multiple instances of services
2. **Database Sharding**: Partition data across multiple databases
3. **Load Balancing**: Distribute traffic across services
4. **Caching Strategy**: Implement distributed caching
5. **Message Queue**: Handle high-volume events

---

## 📞 Support

Để biết thêm thông tin hoặc báo cáo vấn đề, vui lòng liên hệ:
- **Repository**: [GitHub Repository URL]
- **Documentation**: [Documentation URL]
- **Issues**: [GitHub Issues URL] 
