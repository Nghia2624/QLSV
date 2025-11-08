package config

import (
	"fmt"
	"os"
)

type Config struct {
	PostgresHost       string
	PostgresPort       string
	PostgresUser       string
	PostgresPassword   string
	PostgresDB         string
	MongoURL           string
	KafkaBrokers       string
	KafkaStudentTopic  string
	KafkaTeacherTopic  string
	KafkaCourseTopic   string
	KafkaClassTopic    string
	KafkaRegistrationTopic string
	KafkaUserTopic     string
	KafkaConsumerGroup string
	JWTSecret          string
	JWTExpire          string
}

func LoadConfig() *Config {
	return &Config{
		PostgresHost:       getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:       getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:       getEnv("POSTGRES_USER", "qlsv"),
		PostgresPassword:   getEnv("POSTGRES_PASSWORD", "qlsv123"),
		PostgresDB:         getEnv("POSTGRES_DB", "qlsvdb"),
		MongoURL:           getEnv("MONGO_URI", "mongodb://qlsv:qlsv123@localhost:27018/qlsvdb?authSource=admin&authMechanism=SCRAM-SHA-256"),
		KafkaBrokers:       getEnv("KAFKA_BROKERS", "localhost:9093"),
		KafkaStudentTopic:  getEnv("KAFKA_STUDENT_TOPIC", "student-events"),
		KafkaTeacherTopic:  getEnv("KAFKA_TEACHER_TOPIC", "teacher-events"),
		KafkaCourseTopic:   getEnv("KAFKA_COURSE_TOPIC", "course-events"),
		KafkaClassTopic:    getEnv("KAFKA_CLASS_TOPIC", "class-events"),
		KafkaRegistrationTopic: getEnv("KAFKA_REGISTRATION_TOPIC", "registration-events"),
		KafkaUserTopic:     getEnv("KAFKA_USER_TOPIC", "user-events"),
		KafkaConsumerGroup: getEnv("KAFKA_CONSUMER_GROUP", "qlsv-default-group"),
		JWTSecret:          getEnv("JWT_SECRET_KEY", "qlsv-super-secret-key-change-in-production"),
		JWTExpire:          getEnv("JWT_EXPIRATION", "24h"),
	}
}

func (c *Config) BuildPostgresDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
