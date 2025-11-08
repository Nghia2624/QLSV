# QLSVGO - Student Management CQRS Microservices

## Kiến trúc tổng thể
- Golang + CQRS + Kafka + Postgres + MongoDB + Docker
- Command Service: Ghi, validate, sinh event, lưu PostgreSQL, publish Kafka
- Query Service: Lắng nghe Kafka, mapping, lưu MongoDB, phục vụ API đọc
- JWT xác thực, phân quyền (Admin, Student, Teacher)

## Cấu trúc thư mục
- `cmd/command-service`: Service ghi
- `cmd/query-service`: Service đọc
- `internal/domain/model`: Entity, struct
- `internal/domain/event`: Định nghĩa event
- `internal/usecase`: Logic nghiệp vụ
- `internal/repository`: Interface repository
- `internal/infrastructure`: Kết nối DB, Kafka, JWT
- `internal/handler`: HTTP handler, middleware
- `pkg/`: Thư viện dùng chung
- `migrations/`: File SQL khởi tạo DB

## Chạy hệ thống
```sh
docker-compose up --build
```

## Các API chính
- CRUD Sinh viên, Giảng viên, Lớp học, Môn học, Đăng ký
- Đăng nhập, đăng ký, phân quyền

## Ghi chú
- Đảm bảo đã cài Docker, Docker Compose
- Các service sẽ tự động migrate DB khi khởi động 