package main

import (
	"database/sql"
	"log"
	"net/http"
	httpHandler "qlsvgo/internal/handler/http"
	"qlsvgo/internal/handler/middleware"
	"qlsvgo/internal/infrastructure/config"
	"qlsvgo/internal/infrastructure/kafka"
	infraPostgres "qlsvgo/internal/infrastructure/postgres"
	"qlsvgo/internal/usecase"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	// Kết nối Postgres
	db, err := sql.Open("postgres", cfg.BuildPostgresDSN())
	if err != nil {
		log.Fatal("Postgres connect error:", err)
	}
	defer db.Close()

	// Khởi tạo unified Kafka service
	kafkaService := kafka.NewKafkaService([]string{cfg.KafkaBrokers})

	// Khởi tạo Kafka EventPublishers
	studentPublisher := kafkaService.NewEventPublisher("student-events", "Student")
	teacherPublisher := kafkaService.NewEventPublisher("teacher-events", "Teacher")
	coursePublisher := kafkaService.NewEventPublisher("course-events", "Course")
	classPublisher := kafkaService.NewEventPublisher("class-events", "Class")
	registrationPublisher := kafkaService.NewEventPublisher("registration-events", "Registration")
	userPublisher := kafkaService.NewEventPublisher("user-events", "User")
	defer studentPublisher.Close()
	defer teacherPublisher.Close()
	defer coursePublisher.Close()
	defer classPublisher.Close()
	defer registrationPublisher.Close()
	defer userPublisher.Close()

	// Đợi Kafka sẵn sàng
	log.Println("Waiting for Kafka to be ready...")
	time.Sleep(10 * time.Second)
	log.Println("Kafka publishers initialized")

	// Repository thực tế
	studentRepo := &infraPostgres.StudentPostgresRepository{DB: db}
	teacherRepo := &infraPostgres.TeacherPostgresRepository{DB: db}
	courseRepo := &infraPostgres.CoursePostgresRepository{DB: db}
	classRepo := &infraPostgres.ClassPostgresRepository{DB: db}
	registrationRepo := &infraPostgres.RegistrationPostgresRepository{DB: db}
	userRepo := &infraPostgres.UserPostgresRepository{DB: db}

	// Usecase với EventBus
	studentUC := &usecase.StudentUsecase{CommandRepo: studentRepo, EventBus: studentPublisher.Publish}
	teacherUC := &usecase.TeacherUsecase{CommandRepo: teacherRepo, EventBus: teacherPublisher.Publish}
	courseUC := &usecase.CourseUsecase{CommandRepo: courseRepo, EventBus: coursePublisher.Publish}
	classUC := &usecase.ClassUsecase{CommandRepo: classRepo, EventBus: classPublisher.Publish}
	registrationUC := &usecase.RegistrationUsecase{CommandRepo: registrationRepo, EventBus: registrationPublisher.Publish}
	userUC := &usecase.UserUsecase{Repo: userRepo, EventBus: userPublisher.Publish}

	// Inject vào handlers
	httpHandler.StudentUC = studentUC
	httpHandler.TeacherUC = teacherUC
	httpHandler.CourseUC = courseUC
	httpHandler.ClassUC = classUC
	httpHandler.RegistrationUC = registrationUC
	httpHandler.UserUC = userUC
	httpHandler.JWTSecret = cfg.JWTSecret

	jwtMw := middleware.JWTMiddleware(cfg.JWTSecret)
	adminMw := middleware.RoleMiddleware("admin")

	mux := http.NewServeMux()

	// Healthcheck
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Auth/User (không cần JWT cho register/login)
	mux.HandleFunc("/api/register", httpHandler.RegisterHandler)
	mux.HandleFunc("/api/login", httpHandler.LoginHandler)

	// Command operations only (Create, Update, Delete)
	// Student Commands
	mux.HandleFunc("/api/students", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			jwtMw(adminMw(http.HandlerFunc(httpHandler.CreateStudentHandler))).ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/students/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			jwtMw(adminMw(http.HandlerFunc(httpHandler.UpdateStudentHandler))).ServeHTTP(w, r)
		case http.MethodDelete:
			jwtMw(adminMw(http.HandlerFunc(httpHandler.DeleteStudentHandler))).ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Teacher Commands
	mux.HandleFunc("/api/teachers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			jwtMw(adminMw(http.HandlerFunc(httpHandler.CreateTeacherHandler))).ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/teachers/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			jwtMw(adminMw(http.HandlerFunc(httpHandler.UpdateTeacherHandler))).ServeHTTP(w, r)
		case http.MethodDelete:
			jwtMw(adminMw(http.HandlerFunc(httpHandler.DeleteTeacherHandler))).ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Course Commands
	mux.HandleFunc("/api/courses", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			jwtMw(adminMw(http.HandlerFunc(httpHandler.CreateCourseHandler))).ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/courses/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			jwtMw(adminMw(http.HandlerFunc(httpHandler.UpdateCourseHandler))).ServeHTTP(w, r)
		case http.MethodDelete:
			jwtMw(adminMw(http.HandlerFunc(httpHandler.DeleteCourseHandler))).ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Class Commands
	mux.HandleFunc("/api/classes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			jwtMw(adminMw(http.HandlerFunc(httpHandler.CreateClassHandler))).ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/classes/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			jwtMw(adminMw(http.HandlerFunc(httpHandler.UpdateClassHandler))).ServeHTTP(w, r)
		case http.MethodDelete:
			jwtMw(adminMw(http.HandlerFunc(httpHandler.DeleteClassHandler))).ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Registration Commands
	mux.HandleFunc("/api/registrations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			jwtMw(adminMw(http.HandlerFunc(httpHandler.CreateRegistrationHandler))).ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/registrations/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			jwtMw(adminMw(http.HandlerFunc(httpHandler.UpdateRegistrationHandler))).ServeHTTP(w, r)
		case http.MethodDelete:
			jwtMw(adminMw(http.HandlerFunc(httpHandler.DeleteRegistrationHandler))).ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Println("Command Service started on :8080")
	log.Println("Command Service: Write operations only (Create, Update, Delete)")
	http.ListenAndServe(":8080", mux)
}
