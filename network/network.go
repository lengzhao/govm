package network

import "context"

// MessageHandler Used to handle messages received in broadcast mode
type MessageHandler = func(from string, topic string, data []byte) error

// RequestHandler Used to handle requests in point-to-point mode
type RequestHandler = func(from string, reqType string, data []byte) ([]byte, error)

// MessageFilter Used to filter broadcast messages
type MessageFilter = func(from string, topic string, data []byte) bool

// NetworkInterface Network interface, defines all externally provided functions
type NetworkInterface interface {
	// Run Start the network module and run
	Run(ctx context.Context) error

	// BroadcastMessage Broadcast message to specified topic
	BroadcastMessage(topic string, data []byte) error

	// RegisterMessageHandler Register broadcast message handler
	RegisterMessageHandler(topic string, handler MessageHandler)

	// RegisterRequestHandler Register point-to-point request handler
	RegisterRequestHandler(requestType string, handler RequestHandler)

	// SendRequest Send point-to-point request
	SendRequest(peerID string, requestType string, data []byte) ([]byte, error)

	// ConnectToPeer Connect to specified peer
	ConnectToPeer(addr string) error

	// GetPeers Get list of connected peers
	GetPeers() []string

	// GetLocalAddresses Get local peer address list
	GetLocalAddresses() []string

	// GetLocalPeerID Get local peer ID
	GetLocalPeerID() string

	// RegisterMessageFilter Register broadcast message filter
	RegisterMessageFilter(topic string, filter MessageFilter)
}
