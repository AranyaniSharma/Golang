package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
	migrations "university-management-golang/db"
	"university-management-golang/db/connection"
	um "university-management-golang/protoclient/university_management"

	"university-management-golang/server/internal/handlers"
)

const port = "2345"

//db
const (
	username = "postgres"
	password = "admin"
	host     = "localhost"
	dbPort   = "5436"
	dbName   = "postgres"
	schema   = "public"
)

func main() {
	err := migrations.MigrationsUp(username, password, host, dbPort, dbName, schema)

	if err != nil {
		log.Fatalf("Failed to migrate here, err: %+v\n", err)
	}

	connectionmanager := &connection.DatabaseConnectionManagerImpl{
		&connection.DBConfig{
			host, dbPort, username, password, dbName, schema,
		},
		nil,
	}

	//insertDepartmentSeedData(connectionmanager)
	//insertStudentSeedData(connectionmanager)
	//insertAttendanceSeedData(connectionmanager)

	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen to port: %s, err: %+v\n", port, err)
	}
	log.Printf("Starting to listen on port: %s\n", port)

	um.RegisterUniversityManagementServiceServer(grpcServer, handlers.NewUniversityManagementHandler(connectionmanager))

	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to start GRPC Server: %+v\n", err)
	}

}

func insertDepartmentSeedData(connectionManager connection.DatabaseConnectionManager) {
	connection, err := connectionManager.GetConnection()
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	log.Println("Cleaning up departments table")
	_, err = connection.GetSession().DeleteFrom("department").Exec()
	if err != nil {
		log.Fatalf("Could not delete from department table. Err: %+v", err)
	}

	log.Println("Inserting into departments table")
	_, err = connection.GetSession().InsertInto("department").Columns("id", "name").
		Values("1", "Computer Science").Exec()
	_, err = connection.GetSession().InsertInto("department").Columns("id", "name").
		Values("2", "Electronics").Exec()
	_, err = connection.GetSession().InsertInto("department").Columns("id", "name").
		Values("3", "Information technology").Exec()
	_, err = connection.GetSession().InsertInto("department").Columns("id", "name").
		Values("4", "Automobile").Exec()

	if err != nil {
		log.Fatalf("Could not insert into department table. Err: %+v", err)
	}

	defer connectionManager.CloseConnection()
}

func insertStudentSeedData(connectionManager connection.DatabaseConnectionManager) {
	connection, err := connectionManager.GetConnection()
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	log.Println("Cleaning up students table")
	_, err = connection.GetSession().DeleteFrom("students").Exec()
	if err != nil {
		log.Fatalf("Could not delete from students table. Err: %+v", err)
	}

	log.Println("Inserting into students table")
	_, err = connection.GetSession().InsertInto("students").Columns("id", "rollno", "name", "department").
		Values("1", "1001", "Alex", "Computer Science").
		Values("2", "1002", "Jimmy", "Electronics").
		Values("3", "1003", "Stuart", "Information technology").
		Values("4", "1004", "Andrew", "Information technology").
		Values("5", "1005", "Sara", "Computer Science").
		Values("6", "1006", "Robert", "Electronics").
		Values("7", "1007", "Will", "Electronics").Exec()

	if err != nil {
		log.Fatalf("Could not insert into students table. Err: %+v", err)
	}

	defer connectionManager.CloseConnection()
}
func insertAttendanceSeedData(connectionManager connection.DatabaseConnectionManager) {
	connection, err := connectionManager.GetConnection()
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	log.Println("Cleaning up attendance table")
	_, err = connection.GetSession().DeleteFrom("attendance").Exec()
	if err != nil {
		log.Fatalf("Could not delete from students table. Err: %+v", err)
	}

	_, err = connection.GetSession().InsertBySql("INSERT into attendance(studentid, logintime) VALUES (?,?)", 3, time.Now()).Exec()

	//log.Println("Inserting into students table")
	//_, err = connection.GetSession().InsertInto("attendance").Columns("studentid", "logintime").
	//	Values("1", "time.Now().Format(time.RFC3339)").Exec()

	if err != nil {
		log.Fatalf("Could not insert into attendance table. Err: %+v", err)
	}

	defer connectionManager.CloseConnection()
}
