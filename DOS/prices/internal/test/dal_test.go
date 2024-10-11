package verification

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"

	"blaucorp.com/prices/internal/dal"
)

func TestWithMongoDBContainer(t *testing.T) {
	// Request for a MongoDB container
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0", // or the version you need
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(30 * time.Second),
	}

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}
	defer mongoC.Terminate(ctx)

	// Get the host and port for the MongoDB container
	host, err := mongoC.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %s", err)
	}

	port, err := mongoC.MappedPort(ctx, "27017")
	if err != nil {
		t.Fatalf("failed to get container port: %s", err)
	}

	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port.Port())
	var client *mongo.Client
	err = dal.SetUpConnMongoDB(&client, mongoURI)
	if err != nil {
		t.Fatalf("failed to set up mongo client: %s", err)
	}
	defer client.Disconnect(ctx)

	log.Println("MongoDB Testcontainer is running and connected")
}
