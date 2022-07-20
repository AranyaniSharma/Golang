package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"university-management-golang/protoclient/university_management"
)

const (
	host = "localhost"
	port = "2345"
)

func main() {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error occured %+v", err)
	}
	client := university_management.NewUniversityManagementServiceClient(conn)
	var departmentID int32 = 4
	departmentResponse, err := client.GetDepartment(context.TODO(), &university_management.GetDepartmentRequest{Id: departmentID})
	if err != nil {
		log.Fatalf("Error occured while fetching department for id %d,err: %+v", departmentID, err)
	}

	log.Println("response from server", departmentResponse)

	var department = "Electronics"
	studentResponse, err := client.GetStudent(context.TODO(), &university_management.GetStudentRequest{Department: department})
	if err != nil {
		log.Fatalf("Error occured while fetching student %+v err: ", department, err)
	}
	log.Println("response from server", studentResponse)

	loginAttendace := university_management.Attendance{StudentId: 5}

	loginResponse, err := client.RecordStudentLoginTime(context.TODO(), &university_management.GetLoginRequest{Attendance: &loginAttendace})
	if err != nil {
		log.Fatalf("Error occured while entering lohin timings %+v err: ", loginResponse, err)
	}
	log.Println("response from server", loginResponse)

	logoutAttendance := university_management.Attendance{StudentId: 2}

	logoutResponse, err := client.RecordStudentLogoutTime(context.TODO(), &university_management.GetLogoutRequest{Attendance: &logoutAttendance})
	if err != nil {
		log.Fatalf("Error occured while entering logout timings %+v err: ", logoutAttendance, err)
	}
	log.Println("response from server", logoutResponse)
}
