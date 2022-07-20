package handlers

import (
	"encoding/json"
	"golang.org/x/net/context"
	"log"
	"strings"
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

	connection.GetSession().Select("id", "name").From("department").Where("id = ?", request.GetId()).LoadOne(&department)

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

	if request.GetDepartment() == "" {
		_, err = connection.GetSession().Select("id", "name", "department").From("students").Load(&students)
		if err != nil {
			log.Fatalf("Error in getting session is: %+v", err)
		}
	} else {
		_, err = connection.GetSession().Select("id", "name", "department").From("students").Where("department=?", request.GetDepartment()).Load(&students)
		if err != nil {
			log.Fatalf("Error in getting session is: %+v", err)
		}
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

	connection.GetSession().QueryRow("SELECT  logintime FROM attendance WHERE studentid = $1 ORDER BY logintime desc ", request.Attendance.GetStudentId()).
		Scan(&attendance.LoginTime)
	log.Printf("login time is   %v", attendance.LoginTime)
	var response um.GetLoginResponse

	var rollNo int32
	connection.GetSession().Select("rollno").From("students").Where("id=?", request.Attendance.GetStudentId()).Load(rollNo)
	go notifyWhenStudentLogsIn(request, rollNo)

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

	log.Printf("Logout Started   Recieved message from client: %v", request.Attendance)

	var attendance um.Attendance
	//_, err = connection.GetSession().Select("id").From("attendance").Where("studentid=?", request.Attendance.GetStudentId()).Load(&attendance.Id)
	//if err != nil {
	//	log.Fatalf("Error in getting session is: %+v", err)
	//}

	connection.GetSession().QueryRow("SELECT  logintime,id FROM attendance WHERE studentid = $1 ORDER BY logintime desc", request.Attendance.GetStudentId()).
		Scan(&attendance.LoginTime, &attendance.Id)

	log.Printf("login id  is   %v Logintime is %v", attendance.Id, attendance.LoginTime)
	var response um.GetLogoutResponse

	if attendance.Id == 0 {
		response = um.GetLogoutResponse{LoginMessage: "User have not logged in yet"}
		return &response, nil
	}
	go notifyWHenStudentsLogoutBefore8Hours(attendance.LoginTime, time.Now(), request)
	//_, err = connection.GetSession().UpdateBySql("update into attendance(logouttime) where id VALUES (?,?)", attendance.Id, time.Now()).Exec()
	_, err = connection.GetSession().Update("attendance").Where("id=?", attendance.GetId()).Set("logouttime", time.Now()).Exec()
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	response = um.GetLogoutResponse{LoginMessage: "Successfully entered logout time"}
	return &response, nil
}

func (u *universityManagementServer) GetStaff(ctx context.Context, request *um.GetStaffRequest) (*um.GetStaffResponse, error) {
	connection, err := u.connectionManager.GetConnection()
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}
	defer u.connectionManager.CloseConnection()

	log.Printf("GetStaffed   Recieved message from client: %v", request.RollNo)

	var department string
	_, err = connection.GetSession().Select("department").From("students").Where("rollno=?", request.GetRollNo()).Load(&department)
	if err != nil {
		log.Fatalf("Error in getting session is: %+v", err)
	}
	log.Printf("departmemnt %v", department)
	var departmentId int32
	_, err = connection.GetSession().Select("id").From("department").Where("name=?", department).Load(&departmentId)
	if err != nil {
		log.Fatalf("Error in getting session is: %+v", err)
	}
	log.Printf("department id %v ", departmentId)
	var staffId []int32
	_, err = connection.GetSession().Select("staffid").From("staff_department").Where("departmentid=?", departmentId).Load(&staffId)
	if err != nil {
		log.Fatalf("Error in getting session is: %+v", err)
	}
	log.Printf("staff id  %v", staffId)

	var staff []um.Staff

	_, err = connection.GetSession().Select("staffid", "name").From("staff").Where("staffid IN ?", staffId).Load(&staff)
	if err != nil {
		log.Fatalf("Error in getting session is: %+v", err)
	}
	log.Printf("staff %v", staff[0].Staffid)

	var staff_list *um.GetStaffResponse = &um.GetStaffResponse{}

	for _, staff := range staff {
		s := um.Staff{
			Staffid: staff.Staffid,
			Name:    staff.Name,
		}
		staff_list.Staff = append(staff_list.Staff, &s)
	}
	return staff_list, nil

}

func notifyWhenStudentLogsIn(request *um.GetLoginRequest, rollNo int32) {
	if rollNo == 0 {
		log.Printf("Student logged in without roll No")
	} else {
		log.Printf("Student with id %v joined with roll No %v", request.Attendance.StudentId, rollNo)
	}

}

func notifyWHenStudentsLogoutBefore8Hours(loginTime string, logoutTime time.Time, request *um.GetLogoutRequest) {

	log.Printf("login time is   %v", loginTime)
	loginTime = strings.Replace(loginTime, "T", " ", 1)
	log.Printf(" new login time is   %v", loginTime)

	lt, err := time.Parse("2006-01-02 15:04:05", loginTime[0:19])
	if err != nil {
		log.Fatalf("Logintime Error in notifyWHenStudentsLogoutBefore8Hours: %+v", err)
	}
	ct, err := time.Parse("2006-01-02 15:04:05", logoutTime.String()[0:19])
	if err != nil {
		log.Fatalf("currentime Error in notifyWHenStudentsLogoutBefore8Hours: %+v", err)
	}

	lt = lt.Add(8 * time.Hour)
	log.Printf("login time after 8 bourws %v", lt)
	log.Printf("current time  %v", ct)
	if lt.After(ct) {
		log.Printf("Student with id %v logged off before completing 8 hours", request.Attendance.GetStudentId())

	}

}
func NewUniversityManagementHandler(connectionmanager connection.DatabaseConnectionManager) um.UniversityManagementServiceServer {
	return &universityManagementServer{
		connectionManager: connectionmanager,
	}
}
