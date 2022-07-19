package handlers

import (
	"encoding/json"
	"golang.org/x/net/context"
	"log"
	"time"
	"university-management-golang/db/connection"
	um "university-management-golang/protoclient/university_management"
)

type universityManagementServer struct {
	um.UniversityManagementServiceServer
	connectionManager connection.DatabaseConnectionManager
}

func (u *universityManagementServer) GetDepartment(ctx context.Context, request *um.GetDepartmentRequest) (*um.GetDepartmentResponse, error) {
	connection, err := u.connectionManager.GetConnection()
	defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	log.Printf("Recieved message from client: %v", request.Id)
	var department um.Department

	connection.GetSession().Select("id", "name").From("departments").Where("id = ?", request.GetId()).LoadOne(&department)

	_, err = json.Marshal(&department)
	if err != nil {
		log.Fatalf("Error while marshaling %+v", err)
	}
	log.Printf("Worked with id and the result is : %v", department)

	return &um.GetDepartmentResponse{Department: &um.Department{
		Id:   department.Id,
		Name: department.Name,
	}}, nil
}

func (u *universityManagementServer) GetStudent(ctx context.Context, request *um.GetStudentRequest) (*um.GetStudentResponse, error) {
	connection, err := u.connectionManager.GetConnection()
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}
	defer u.connectionManager.CloseConnection()

	log.Printf("Recieved message from client: %v", request.Department)
	var students []um.Student
	_, err = connection.GetSession().Select("id", "name", "department").From("students").Where("department=?", request.GetDepartment()).Load(&students)
	if err != nil {
		log.Fatalf("Error in getting session is: %+v", err)
	}

	var student_list *um.GetStudentResponse = &um.GetStudentResponse{}

	for _, st := range students {
		stu := um.Student{
			Id:         st.Id,
			Name:       st.Name,
			Department: st.Department,
		}
		student_list.Students = append(student_list.Students, &stu)
	}
	return student_list, nil
}

func (u *universityManagementServer) RecordStudentLoginTime(ctx context.Context, request *um.GetLoginRequest) (*um.GetLoginResponse, error) {
	connection, err := u.connectionManager.GetConnection()
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}
	defer u.connectionManager.CloseConnection()

	log.Printf("Recieved message from client: %v", request.Attendance)
	var attendance um.Attendance

	connection.GetSession().QueryRow("SELECT  logintime FROM attendance WHERE studentid = $1 ", request.Attendance.GetStudentId()).
		Scan(&attendance.LoginTime)
	log.Printf("login time is   %v", attendance.LoginTime)
	var response um.GetLoginResponse
	if attendance.LoginTime != "" {
		log.Printf("login time is   %v", attendance.LoginTime)
		loginTime, err := time.Parse("2006-01-02", attendance.LoginTime[0:10])
		if err != nil {
			log.Fatalf("Logintime Error: %+v", err)
		}
		currentTime, err := time.Parse("2006-01-02", time.Now().String()[0:10])
		if err != nil {
			log.Fatalf("currentime Error: %+v", err)
		}

		log.Printf("login time is   %v", loginTime)

		if loginTime.Equal(currentTime) {
			response = um.GetLoginResponse{LoginMessage: "Already Logged in for a day!!"}
			return &response, nil
		}
		_, err = connection.GetSession().InsertBySql("INSERT into attendance(studentid, logintime) VALUES (?,?)", request.Attendance.GetStudentId(), time.Now()).Exec()

		if err != nil {
			log.Fatalf("Error: %+v", err)
		}

		response = um.GetLoginResponse{LoginMessage: "Successfully entered login time"}
		return &response, nil
	}

	_, err = connection.GetSession().InsertBySql("INSERT into attendance(studentid, logintime) VALUES (?,?)", request.Attendance.GetStudentId(), time.Now()).Exec()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	response = um.GetLoginResponse{LoginMessage: "Successfully entered login time"}
	return &response, nil
}

func (u *universityManagementServer) RecordStudentLogoutTime(ctx context.Context, request *um.GetLogoutRequest) (*um.GetLogoutResponse, error) {
	connection, err := u.connectionManager.GetConnection()
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}
	defer u.connectionManager.CloseConnection()

	log.Printf("Recieved message from client: %v", request.Attendance)

	var attendance um.Attendance
	_, err = connection.GetSession().Select("id").From("attendance").Where("studentid=?", request.Attendance.GetStudentId()).Load(&attendance.Id)
	if err != nil {
		log.Fatalf("Error in getting session is: %+v", err)
	}
	log.Printf("login id  is   %v", attendance.Id)
	var response um.GetLogoutResponse

	if attendance.Id == 0 {
		response = um.GetLogoutResponse{LoginMessage: "User have not logged in yet"}
		return &response, nil
	}
	//_, err = connection.GetSession().UpdateBySql("update into attendance(logouttime) where id VALUES (?,?)", attendance.Id, time.Now()).Exec()
	_, err = connection.GetSession().Update("attendance").Where("id=?", attendance.GetId()).Set("logouttime", time.Now()).Exec()
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	response = um.GetLogoutResponse{LoginMessage: "Successfully entered logout time"}
	return &response, nil
}

func notify(request *um.GetLoginRequest) {
	log.Println("jhh")

}

func NewUniversityManagementHandler(connectionmanager connection.DatabaseConnectionManager) um.UniversityManagementServiceServer {
	return &universityManagementServer{
		connectionManager: connectionmanager,
	}
}
