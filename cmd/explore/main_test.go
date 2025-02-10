package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"muzz-project/service"
	"muzz-project/service/protos"
	"muzz-project/storage/mysql"
	"net"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

var (
	connectionStringVar string
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	container, connectionString := setupMySQL(ctx)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			log.Fatalf("Failed to terminate MySQL container: %v", err)
		}
	}()

	time.Sleep(5 * time.Second)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	defer db.Close()

	//wait for db to start
	time.Sleep(5 * time.Second)

	initialiseDB(db)

	code := m.Run()
	defer os.Exit(code)

}

func setupMySQL(ctx context.Context) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      "testdb",
			"MYSQL_USER":          "testuser",
			"MYSQL_PASSWORD":      "testpassword",
		},
		WaitingFor: wait.ForLog("ready for connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start MySQL container: %v", err)
	}

	portObj, err := container.MappedPort(ctx, "3306")
	if err != nil {
		log.Fatalf("Failed to get MySQL container port: %v", err)
	}

	host, _ := container.Host(ctx)
	dbConnStr := fmt.Sprintf("testuser:testpassword@tcp(%s:%s)/testdb?parseTime=true", host, portObj.Port())
	connectionStringVar = dbConnStr

	return container, dbConnStr
}

func initialiseDB(db *sql.DB) {
	_, err := db.Exec(CreateUserTable)
	if err != nil {
		log.Fatalf("Failed to create user table: %v", err)
	}

	_, err = db.Exec(CreateDecisionsTable)
	if err != nil {
		log.Fatalf("Failed to create decisions table: %v", err)
	}

	_, err = db.Exec(AddDummyUserData)
	if err != nil {
		log.Fatalf("Failed to add user data: %v", err)
	}

	_, err = db.Exec(AddDummyDecisionData)
	if err != nil {
		log.Fatalf("Failed to add decisions data: %v", err)
	}

}

func startServer(port string) {
	db, err := sql.Open("mysql", connectionStringVar)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	defer db.Close()

	s := mysql.NewMysqlStorage(db, maxPageSize)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	protos.RegisterExploreServiceServer(grpcServer, service.NewExploreService(s, maxPageSize))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func TestCountLikesForUser(t *testing.T) {
	ctx := context.Background()
	port := "50051"

	go startServer(port)

	client, conn, err := getClientAndConnection(port, time.Second*10)
	assert.NoError(t, err)

	defer conn.Close()

	out, err := client.CountLikedYou(ctx, &protos.CountLikedYouRequest{
		RecipientUserId: "1",
	})

	assert.NotNil(t, out)
	assert.Equal(t, uint64(3), out.GetCount())
	assert.NoError(t, err)
}

func TestListLikesForUser(t *testing.T) {
	ctx := context.Background()
	port := "50052"

	go startServer(port)

	time.Sleep(5 * time.Second)

	client, conn, err := getClientAndConnection(port, time.Second*10)
	assert.NoError(t, err)
	defer conn.Close()

	expectedOut := &protos.ListLikedYouResponse{
		Likers: []*protos.ListLikedYouResponse_Liker{
			{
				ActorId:       "4",
				UnixTimestamp: 0,
			},
			{
				ActorId:       "6",
				UnixTimestamp: 0,
			},
			{
				ActorId:       "10",
				UnixTimestamp: 0,
			},
		},
		NextPaginationToken: nil,
	}

	out, err := client.ListLikedYou(ctx, &protos.ListLikedYouRequest{
		RecipientUserId: "1",
	})

	assert.NotNil(t, out)

	//can't account for arbitrary timestamp
	for _, op := range out.Likers {
		op.UnixTimestamp = 0
	}
	assert.Equal(t, expectedOut.GetLikers(), out.GetLikers())
	assert.Equal(t, expectedOut.GetNextPaginationToken(), out.GetNextPaginationToken())
	assert.NoError(t, err)
}

// Should have similar output to ListLikes except, because user 1 has also liked user 4, this is not new so should be missing from output
func TestListNewLikesForUser(t *testing.T) {
	ctx := context.Background()
	port := "50053"

	go startServer(port)

	time.Sleep(5 * time.Second)

	client, conn, err := getClientAndConnection(port, time.Second*10)
	assert.NoError(t, err)
	defer conn.Close()

	expectedOut := &protos.ListLikedYouResponse{
		Likers: []*protos.ListLikedYouResponse_Liker{
			{
				ActorId:       "6",
				UnixTimestamp: 0,
			},
			{
				ActorId:       "10",
				UnixTimestamp: 0,
			},
		},
		NextPaginationToken: nil,
	}

	out, err := client.ListNewLikedYou(ctx, &protos.ListLikedYouRequest{
		RecipientUserId: "1",
	})

	assert.NotNil(t, out)

	//can't account for arbitrary timestamp
	for _, op := range out.Likers {
		op.UnixTimestamp = 0
	}
	assert.Equal(t, expectedOut.GetLikers(), out.GetLikers())
	assert.Equal(t, expectedOut.GetNextPaginationToken(), out.GetNextPaginationToken())
	assert.NoError(t, err)
}

func TestPutDecision_LikeMatch(t *testing.T) {
	ctx := context.Background()
	port := "50054"

	go startServer(port)

	time.Sleep(5 * time.Second)

	client, conn, err := getClientAndConnection(port, time.Second*10)
	assert.NoError(t, err)
	defer conn.Close()

	expectedOut := &protos.PutDecisionResponse{
		MutualLikes: true,
	}

	out, err := client.PutDecision(ctx, &protos.PutDecisionRequest{
		ActorUserId:     "2",
		RecipientUserId: "1",
		LikedRecipient:  true,
	})

	assert.NotNil(t, out)
	assert.NoError(t, err)
	assert.Equal(t, expectedOut.GetMutualLikes(), out.GetMutualLikes())

	count, err := client.CountLikedYou(ctx, &protos.CountLikedYouRequest{
		RecipientUserId: "1",
	})

	assert.NoError(t, err)
	assert.Equal(t, uint64(4), count.GetCount())
}

func TestPutDecision_LikeNoMatch(t *testing.T) {
	ctx := context.Background()
	port := "50055"

	go startServer(port)

	time.Sleep(5 * time.Second)

	client, conn, err := getClientAndConnection(port, time.Second*10)
	assert.NoError(t, err)
	defer conn.Close()

	expectedOut := &protos.PutDecisionResponse{
		MutualLikes: false,
	}

	out, err := client.PutDecision(ctx, &protos.PutDecisionRequest{
		ActorUserId:     "5",
		RecipientUserId: "1",
		LikedRecipient:  true,
	})

	assert.NotNil(t, out)
	assert.NoError(t, err)
	assert.Equal(t, expectedOut.GetMutualLikes(), out.GetMutualLikes())

	count, err := client.CountLikedYou(ctx, &protos.CountLikedYouRequest{
		RecipientUserId: "1",
	})

	assert.NoError(t, err)
	assert.Equal(t, uint64(5), count.GetCount())
}

func TestPutDecision_NoLike(t *testing.T) {
	ctx := context.Background()
	port := "50056"

	go startServer(port)

	time.Sleep(5 * time.Second)

	client, conn, err := getClientAndConnection(port, time.Second*10)
	assert.NoError(t, err)
	defer conn.Close()

	expectedOut := &protos.PutDecisionResponse{
		MutualLikes: false,
	}

	out, err := client.PutDecision(ctx, &protos.PutDecisionRequest{
		ActorUserId:     "7",
		RecipientUserId: "1",
		LikedRecipient:  false,
	})

	assert.NotNil(t, out)
	assert.NoError(t, err)
	assert.Equal(t, expectedOut.GetMutualLikes(), out.GetMutualLikes())

	count, err := client.CountLikedYou(ctx, &protos.CountLikedYouRequest{
		RecipientUserId: "1",
	})

	assert.NoError(t, err)
	assert.Equal(t, uint64(5), count.GetCount())

}

func getClientAndConnection(port string, timeout time.Duration) (protos.ExploreServiceClient, *grpc.ClientConn, error) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			client := protos.NewExploreServiceClient(conn)
			return client, conn, nil // gRPC server is ready
		}
		time.Sleep(500 * time.Millisecond) // Retry interval
	}
	return nil, nil, fmt.Errorf("gRPC server did not start in time")
}
