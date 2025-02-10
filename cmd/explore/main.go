package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"muzz-project/storage/mysql"
	"net"

	"muzz-project/service"
	"muzz-project/service/protos"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

var (
	port        string
	host        string
	database    string
	password    string
	user        string
	maxPageSize int
)

func init() {
	flag.StringVar(&port, "port", "8080", "port to listen on")
	flag.StringVar(&host, "host", "0.0.0.0", "host to listen on")
	flag.StringVar(&database, "db", "testdb", "database name")
	flag.StringVar(&password, "password", "rootpassword", "database password")
	flag.StringVar(&user, "user", "root", "database user")
	flag.IntVar(&maxPageSize, "maxPageSize", 1000, "maximum number of db rows to be returned in one query")
}

func main() {
	flag.Parse()

	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", user, password, "db", database)
	db, err := sql.Open("mysql", dsn)
	log.Printf("connected to db")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	s := mysql.NewMysqlStorage(db, maxPageSize)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	protos.RegisterExploreServiceServer(grpcServer, service.NewExploreService(s, maxPageSize))
	log.Printf("server listening at %s", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
