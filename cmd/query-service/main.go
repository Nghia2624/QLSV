package main

import (
	"context"
	"log"
	"net/http"
	"qlsvgo/internal/domain/event"
	httpHandler "qlsvgo/internal/handler/http"
	"qlsvgo/internal/infrastructure/config"
	"qlsvgo/internal/infrastructure/kafka"
	infraMongo "qlsvgo/internal/infrastructure/mongo"
	"qlsvgo/internal/mapping"
	"qlsvgo/internal/usecase"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.LoadConfig()

	// Kết nối MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURL))
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}
	db := client.Database("qlsvdb")

	// MongoDB Collections
	studentCol := db.Collection("students")
	teacherCol := db.Collection("teachers")
	courseCol := db.Collection("courses")
	classCol := db.Collection("classes")
	registrationCol := db.Collection("registrations")
	userCol := db.Collection("users")

	// MongoDB Repositories
	studentRepo := &infraMongo.StudentMongoRepository{Collection: studentCol}
	teacherRepo := &infraMongo.TeacherMongoRepository{Collection: teacherCol}
	courseRepo := &infraMongo.CourseMongoRepository{Collection: courseCol}
	classRepo := &infraMongo.ClassMongoRepository{Collection: classCol}
	registrationRepo := &infraMongo.RegistrationMongoRepository{Collection: registrationCol}
	userRepo := &infraMongo.UserMongoRepository{Collection: userCol}

	// Usecase với QueryRepo
	studentUC := &usecase.StudentUsecase{QueryRepo: studentRepo}
	teacherUC := &usecase.TeacherUsecase{QueryRepo: teacherRepo}
	courseUC := &usecase.CourseUsecase{QueryRepo: courseRepo}
	classUC := &usecase.ClassUsecase{QueryRepo: classRepo}
	registrationUC := &usecase.RegistrationUsecase{QueryRepo: registrationRepo}

	// Inject vào handlers
	httpHandler.StudentUC = studentUC
	httpHandler.TeacherUC = teacherUC
	httpHandler.CourseUC = courseUC
	httpHandler.ClassUC = classUC
	httpHandler.RegistrationUC = registrationUC
	httpHandler.UserRepo = userRepo

	// Khởi tạo unified Kafka service
	kafkaService := kafka.NewKafkaService([]string{cfg.KafkaBrokers})

	// Khởi tạo Kafka EventConsumers với wrapper functions
	studentConsumer := kafkaService.NewEventConsumer("student-events", "query-service-group", "Student", studentCol, func(evt *event.Event) interface{} {
		return mapping.EventToStudent(evt)
	})
	teacherConsumer := kafkaService.NewEventConsumer("teacher-events", "query-service-group", "Teacher", teacherCol, func(evt *event.Event) interface{} {
		return mapping.EventToTeacher(evt)
	})
	courseConsumer := kafkaService.NewEventConsumer("course-events", "query-service-group", "Course", courseCol, func(evt *event.Event) interface{} {
		return mapping.EventToCourse(evt)
	})
	classConsumer := kafkaService.NewEventConsumer("class-events", "query-service-group", "Class", classCol, func(evt *event.Event) interface{} {
		return mapping.EventToClass(evt)
	})
	registrationConsumer := kafkaService.NewEventConsumer("registration-events", "query-service-group", "Registration", registrationCol, func(evt *event.Event) interface{} {
		return mapping.EventToRegistration(evt)
	})
	userConsumer := kafkaService.NewEventConsumer("user-events", "query-service-group", "User", userCol, func(evt *event.Event) interface{} {
		return mapping.EventToUser(evt)
	})

	defer studentConsumer.Close()
	defer teacherConsumer.Close()
	defer courseConsumer.Close()
	defer classConsumer.Close()
	defer registrationConsumer.Close()
	defer userConsumer.Close()

	log.Println("Query Service: Kafka consumers initialized")

	// Start Kafka consumers
	go studentConsumer.Start()
	go teacherConsumer.Start()
	go courseConsumer.Start()
	go classConsumer.Start()
	go registrationConsumer.Start()
	go userConsumer.Start()

	mux := http.NewServeMux()

	// Healthcheck
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Query operations only (Read)
	// User Queries
	mux.HandleFunc("/api/users", httpHandler.GetAllUsersHandler)
	mux.HandleFunc("/api/users/", httpHandler.GetUserHandler)

	// Student Queries
	mux.HandleFunc("/api/students", httpHandler.GetAllStudentsHandler)
	mux.HandleFunc("/api/students/", httpHandler.GetStudentHandler)

	// Teacher Queries
	mux.HandleFunc("/api/teachers", httpHandler.GetAllTeachersHandler)
	mux.HandleFunc("/api/teachers/", httpHandler.GetTeacherHandler)

	// Course Queries
	mux.HandleFunc("/api/courses", httpHandler.GetAllCoursesHandler)
	mux.HandleFunc("/api/courses/", httpHandler.GetCourseHandler)

	// Class Queries
	mux.HandleFunc("/api/classes", httpHandler.GetAllClassesHandler)
	mux.HandleFunc("/api/classes/", httpHandler.GetClassHandler)

	// Registration Queries
	mux.HandleFunc("/api/registrations", httpHandler.GetAllRegistrationsHandler)
	mux.HandleFunc("/api/registrations/", httpHandler.GetRegistrationHandler)

	log.Println("Query Service started on :8081")
	log.Println("Query Service: Read operations only (Get, GetAll)")
	http.ListenAndServe(":8081", mux)
}
