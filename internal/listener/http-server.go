package listener

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// HttpFileServer handles a simple HTTP server using Gin.
type HttpFileServer struct {
	listener    net.Listener
	port        int
	server      *http.Server // Gin engine is wrapped by http.Server
	engine      *gin.Engine
	FileToServe string
}

func (s *HttpFileServer) fileHandlerGin(c *gin.Context) {
	if s.FileToServe == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file to serve"})
		return
	}
	name := sanitizeFilename(filepath.Base(s.FileToServe))
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
	c.File(s.FileToServe)
}

// NewHttpFileServer creates and starts a new HttpFileServer using Gin.
// It listens on a random available TCP port.
func NewHttpFileServer() (*HttpFileServer, error) {
	gin.SetMode(gin.ReleaseMode) // Or gin.DebugMode
	router := gin.New()          // Use gin.Default() if you want default middleware (logger, recovery)

	listener, err := net.Listen("tcp", "0.0.0.0:0") // :0 means a random available port
	if err != nil {
		return nil, fmt.Errorf("failed to listen on a random port: %w", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port

	srv := &http.Server{
		Handler: router, // Gin engine acts as the handler
	}

	service := &HttpFileServer{
		listener: listener,
		port:     port,
		server:   srv,
		engine:   router,
		// FileToServe will be set later via SetFileToServe
	}

	// Now that service is initialized, we can set its method as a handler
	router.GET("/file", service.fileHandlerGin)

	go func() {
		if err := service.server.Serve(service.listener); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server (Gin) error: %v", err)
		}
	}()

	log.Printf("HttpFileServer (Gin) listening on port %d", port)

	return service, nil
}

// Port returns the port the service is listening on.
func (s *HttpFileServer) Port() int {
	return s.port
}

// Addresses returns a list of address strings for the server,
// preferring non-loopback IPs. Falls back to localhost.
func (s *HttpFileServer) Addresses() []string {
	var addresses []string
	localIPs, err := getLocalIPs()

	if err != nil {
		log.Printf("Warning: Error getting local IPs for Addresses(): %v. Only localhost might be accurate.", err)
	}

	if err == nil && len(localIPs) > 0 {
		for _, ip := range localIPs {
			// Skip loopback addresses
			if net.ParseIP(ip).IsLoopback() {
				continue
			}
			addresses = append(addresses, fmt.Sprintf("%s:%d", ip, s.port))
		}
	}

	// Remove duplicates if any (e.g. if a localIP was 127.0.0.1)
	// This is a simple way, more robust would be to parse and compare IPs
	uniqueAddresses := make([]string, 0, len(addresses))
	seen := make(map[string]bool)
	for _, addr := range addresses {
		if !seen[addr] {
			uniqueAddresses = append(uniqueAddresses, addr)
			seen[addr] = true
		}
	}
	if len(uniqueAddresses) == 0 {
		// If no unique addresses found, fallback to localhost
		uniqueAddresses = []string{fmt.Sprintf("localhost:%d", s.port)}
	}
	return uniqueAddresses
}

func (s *HttpFileServer) SetFileToServe(file string) {
	s.FileToServe = file
}

// Close gracefully shuts down the server.
func (s *HttpFileServer) Close() error {
	log.Printf("Shutting down HttpFileServer (Gin) on port %d", s.port)
	if s.server != nil {
		// For Gin, http.Server.Close() is still the way to close the underlying listener
		// and stop accepting new connections. For graceful shutdown of active connections,
		// http.Server.Shutdown(ctx) would be used.
		return s.server.Close()
	}
	return nil
}
