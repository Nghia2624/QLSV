package kafka

import (
	"context"
	"log"
	"qlsvgo/internal/domain/event"
	"qlsvgo/internal/mapping"
	"reflect"
	"time"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// KafkaService provides unified Kafka operations
type KafkaService struct {
	brokers []string
}

// NewKafkaService creates a new Kafka service
func NewKafkaService(brokers []string) *KafkaService {
	return &KafkaService{
		brokers: brokers,
	}
}

// EventPublisher handles event publishing
type EventPublisher struct {
	writer *kafka.Writer
	entity string
}

// NewEventPublisher creates a new event publisher
func (ks *KafkaService) NewEventPublisher(topic, entity string) *EventPublisher {
	return &EventPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(ks.brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
		entity: entity,
	}
}

// Publish publishes an event
func (ep *EventPublisher) Publish(evt *event.Event) error {
	data, err := mapping.SerializeEvent(evt)
	if err != nil {
		log.Printf("Failed to serialize %s event: %v", ep.entity, err)
		return err
	}

	err = ep.writer.WriteMessages(context.Background(), kafka.Message{
		Value: data,
	})
	if err != nil {
		log.Printf("Failed to publish %s event: %v", ep.entity, err)
		return err
	}

	log.Printf("✓ %s event published: %s", ep.entity, evt.Type)
	return nil
}

// Close closes the publisher
func (ep *EventPublisher) Close() error {
	return ep.writer.Close()
}

// EventConsumer handles event consumption
type EventConsumer struct {
	reader     *kafka.Reader
	collection *mongo.Collection
	entity     string
	mapper     func(*event.Event) interface{}
}

// NewEventConsumer creates a new event consumer
func (ks *KafkaService) NewEventConsumer(topic, groupID, entity string, collection *mongo.Collection, mapper func(*event.Event) interface{}) *EventConsumer {
	return &EventConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  ks.brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
		collection: collection,
		entity:     entity,
		mapper:     mapper,
	}
}

// Start starts consuming events
func (ec *EventConsumer) Start() {
	log.Printf("Starting %s Kafka consumer...", ec.entity)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		msg, err := ec.reader.ReadMessage(ctx)
		cancel()

		if err != nil {
			if err.Error() != "context deadline exceeded" {
				log.Printf("%s Kafka read error: %v", ec.entity, err)
			}
			continue
		}

		log.Printf("%s event received, data length: %d bytes", ec.entity, len(msg.Value))

		evt, err := mapping.DeserializeEvent(msg.Value)
		if err != nil {
			log.Printf("%s event deserialize error: %v", ec.entity, err)
			continue
		}

		log.Printf("%s event type: %s", ec.entity, evt.Type)

		entity := ec.mapper(evt)
		if entity != nil {
			log.Printf("%s mapped successfully", ec.entity)
			entityID := getEntityID(entity)

			_, err = ec.collection.UpdateOne(
				context.Background(),
				map[string]interface{}{"_id": entityID},
				map[string]interface{}{"$set": entity},
				options.Update().SetUpsert(true),
			)

			if err != nil {
				log.Printf("Mongo upsert %s error: %v", ec.entity, err)
			} else {
				log.Printf("✓ %s synced to MongoDB: %s", ec.entity, entityID)
			}
		} else {
			log.Printf("Failed to map %s from event", ec.entity)
		}
	}
}

// Close closes the consumer
func (ec *EventConsumer) Close() error {
	return ec.reader.Close()
}

// getEntityID extracts ID from entity using reflection
func getEntityID(entity interface{}) string {
	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Try to get ID field
	idField := v.FieldByName("ID")
	if idField.IsValid() && idField.Kind() == reflect.String {
		return idField.String()
	}

	return ""
}
