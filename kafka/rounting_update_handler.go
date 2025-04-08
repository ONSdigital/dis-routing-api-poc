package kafka

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/ONSdigital/dis-routing-api-poc/config"
	"github.com/ONSdigital/dis-routing-api-poc/schema"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/log.go/v2/log"
)

type RoutingUpdateHandler struct {
	Cfg           *config.Config
	Producer      kafka.IProducer
	mu            sync.Mutex
	batch         []map[string]interface{}
	flushInterval time.Duration
}

// NewRoutingUpdateHandler initializes a batch processor for routing updates
func NewRoutingUpdateHandler(kp kafka.IProducer, flushInterval time.Duration) *RoutingUpdateHandler {
	handler := &RoutingUpdateHandler{
		Producer:      kp,
		flushInterval: flushInterval,
	}
	go handler.startFlusher()
	return handler
}

// HandleEvent queues an event for batching
func (h *RoutingUpdateHandler) HandleEvent(ctx context.Context, event map[string]interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.batch = append(h.batch, event)
	log.Info(ctx, "routing update event queued", log.Data{"event": event})
}

// Handle accumulates events into a batch instead of sending immediately
func (h *RoutingUpdateHandler) Handle(ctx context.Context, _ int, msg kafka.Message) error {
	var event map[string]interface{}
	if err := json.Unmarshal(msg.GetData(), &event); err != nil {
		log.Error(ctx, "failed to decode Kafka event", err)
		return err
	}

	h.mu.Lock()
	h.batch = append(h.batch, event)
	h.mu.Unlock()

	log.Info(ctx, "routing update event queued", log.Data{"event": event})
	return nil
}

// startFlusher sends accumulated events in batch at intervals
func (h *RoutingUpdateHandler) startFlusher() {
	for {
		time.Sleep(h.flushInterval)
		h.flush()
	}
}

// flush sends the batched events to Kafka
func (h *RoutingUpdateHandler) flush() {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.batch) == 0 {
		return
	}

	batchMessage, err := json.Marshal(h.batch)
	if err != nil {
		log.Error(context.Background(), "failed to marshal Kafka batch event", err)
		return
	}

	if err := h.Producer.Send(schema.RoutingUpdatedEvent, batchMessage); err != nil {
		log.Error(context.Background(), "failed to publish batch to Kafka", err)
	}

	log.Info(context.Background(), "Kafka batch published", log.Data{"batch_size": len(h.batch)})
	h.batch = nil // Reset batch after sending
}
