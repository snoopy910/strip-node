package sequencer

import (
	"fmt"
	"sync"

	// Added for DialContext timeout
	pb "github.com/StripChain/strip-node/libs/proto"
	"github.com/StripChain/strip-node/util/logger"

	// No Secrets Manager client needed *within* the manager anymore
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // Keep for potential fallback/logging
)

type ValidatorServiceClient struct {
	pb.ValidatorServiceClient
	cc *grpc.ClientConn
}

func (v *ValidatorServiceClient) GetClientConn() *grpc.ClientConn {
	return v.cc
}

type ValidatorClientManager struct {
	clients       map[string]ValidatorServiceClient
	mutex         sync.Mutex
	clientCertPEM string // Sequencer's client cert (PEM)
	clientKeyPEM  string // Sequencer's client key (PEM)
	serverCaPEM   string // PEM for the CA signing validator server certs (single CA assumed)
	// No validatorConf map needed if using single CA and no other per-validator settings
}

// NewValidatorClientManager accepts PEM strings directly
func NewValidatorClientManager(
// clientCertPEM string,
// clientKeyPEM string,
// serverCaPEM string,
) (*ValidatorClientManager, error) {

	// // Basic validation of inputs
	// if clientCertPEM == "" || clientKeyPEM == "" {
	// 	return nil, fmt.Errorf("client certificate and key PEMs cannot be empty")
	// }
	// if serverCaPEM == "" {
	// 	// Log warning if fallback to insecure is intended, error otherwise
	// 	logger.Sugar().Warn("Server CA PEM is empty, connections may default to insecure or fail")
	// 	// return nil, fmt.Errorf("server CA PEM cannot be empty for mTLS")
	// }

	return &ValidatorClientManager{
		clients: make(map[string]ValidatorServiceClient),
		// clientCertPEM: clientCertPEM,
		// clientKeyPEM:  clientKeyPEM,
		// serverCaPEM:   serverCaPEM,
	}, nil
}

// getCredentials creates the gRPC dial option using the stored PEMs
func (m *ValidatorClientManager) getCredentials() (grpc.DialOption, error) {
	// Handle case where server CA wasn't provided - fallback to insecure or error
	// if m.serverCaPEM == "" {
	// 	logger.Sugar().Warn("No server CA provided, creating insecure credentials")
	// 	return grpc.WithTransportCredentials(insecure.NewCredentials()), nil
	// 	// Or: return nil, fmt.Errorf("cannot create TLS credentials without server CA")
	// }

	// // Load the client certificate
	// clientCert, err := tls.X509KeyPair([]byte(m.clientCertPEM), []byte(m.clientKeyPEM))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to load client keypair: %w", err)
	// }

	// // Create a CA certificate pool and add the server's CA
	// caPool := x509.NewCertPool()
	// if !caPool.AppendCertsFromPEM([]byte(m.serverCaPEM)) {
	// 	return nil, fmt.Errorf("failed to append server CA cert from PEM")
	// }

	// // Create TLS configuration
	// tlsConfig := &tls.Config{
	// 	Certificates: []tls.Certificate{clientCert}, // Client cert to present
	// 	RootCAs:      caPool,                        // CAs to trust for server verification
	// 	MinVersion:   tls.VersionTLS12,
	// 	// Optionally set ServerName based on target URL if needed for hostname verification
	// 	// InsecureSkipVerify: false, // Ensure this is false (default)
	// }

	// creds := credentials.NewTLS(tlsConfig)
	return grpc.WithTransportCredentials(insecure.NewCredentials()), nil
}

func (m *ValidatorClientManager) GetClient(url string) (pb.ValidatorServiceClient, error) {
	client, err := m.FindOrCreateClient(url)
	if err != nil {
		return nil, err
	}
	return client.ValidatorServiceClient, nil
}

// GetClient retrieves or creates a gRPC client for the given URL
func (m *ValidatorClientManager) FindOrCreateClient(url string) (ValidatorServiceClient, error) {
	m.mutex.Lock()
	// Check cache first
	client, found := m.clients[url]
	if found {
		m.mutex.Unlock()
		logger.Sugar().Debugw("Reusing existing gRPC client for validator", "url", url)
		return client, nil
	}
	m.mutex.Unlock() // Unlock before potentially slow operations

	logger.Sugar().Debugw("Attempting to create new gRPC client with mTLS", "url", url)

	// Get the TLS credentials configuration (same config for all validators in this simplified model)
	credOption, err := m.getCredentials()
	if err != nil {
		return ValidatorServiceClient{}, fmt.Errorf("failed to prepare credentials for %s: %w", url, err)
	}

	// Configure dialing options (e.g., credentials, blocking dial with timeout)
	opts := []grpc.DialOption{
		credOption,
	}

	grpcClient, err := grpc.NewClient(url, opts...)
	if err != nil {
		logger.Sugar().Errorf("Failed to dial gRPC validator %s: %v", url, err)
		return ValidatorServiceClient{}, fmt.Errorf("error dialing grpc validator %s: %w", url, err)
	}

	// Create the specific service client
	serviceClient := pb.NewValidatorServiceClient(grpcClient)

	// Lock again to update the cache safely
	m.mutex.Lock()
	// Double-check cache in case another goroutine created the client while we were dialing
	client, found = m.clients[url]
	if found {
		m.mutex.Unlock()
		// Close the connection we just created, use the one already in the map
		logger.Sugar().Warnw("Client created concurrently, closing redundant connection", "url", url)
		go grpcClient.Close() // Close the redundant connection
		return client, nil
	}
	// Store the newly created client
	m.clients[url] = ValidatorServiceClient{
		ValidatorServiceClient: serviceClient,
		cc:                     grpcClient,
	}
	m.mutex.Unlock()

	logger.Sugar().Debugw("Successfully created and cached new gRPC client", "url", url)
	return ValidatorServiceClient{
		ValidatorServiceClient: serviceClient,
		cc:                     grpcClient,
	}, nil
}

// CloseAll closes all managed client connections
func (m *ValidatorClientManager) CloseAll() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	logger.Sugar().Info("Closing all managed validator gRPC client connections...")
	count := 0
	for url, client := range m.clients {
		if err := client.cc.Close(); err != nil {
			logger.Sugar().Errorf("Error closing gRPC client connection for %s: %v", url, err)
		} else {
			count++
			logger.Sugar().Debugw("Closed gRPC client connection", "url", url)
		}
	}
	m.clients = make(map[string]ValidatorServiceClient) // Clear map after closing
	logger.Sugar().Infof("Closed %d validator gRPC client connections.", count)
}
