package ws

import "context"

// SyncMessage defines the structure for synchronization messages.
type SyncMessage struct {
	ClientID string
	Data     []byte
}

// SyncAdapter abstracts the pub/sub system used to sync messages across nodes.
type SyncAdapter interface {
	// Publish sends a sync message to the backend.
	Publish(ctx context.Context, msg SyncMessage) error
	// Subscribe starts listening for sync messages and calls handler on each message.
	Subscribe(ctx context.Context, handler func(msg SyncMessage))
}
