package handlers

//
//import (
//	"encoding/json"
//	"golang.org/x/net/context"
//	"log"
//	"university-management-golang/db/connection"
//	sm "university-management-golang/protoclient/student_management"
//	um "university-management-golang/protoclient/university_management"
//)
//
//type studentManagementServer struct {
//	sm.StudentManagementServiceServer
//	connectionManager connection.DatabaseConnectionManager
//}
//
//func (s *studentManagementServer) GetStudent(ctx context.Context, request *sm.GetStudentRequest) (*sm.GetStudentResponse, error) {
//	connection, err := s.connectionManager.GetConnection()
//
//	if err != nil {
//		log.Fatalf("Error: %+v", err)
//	}
//	defer s.connectionManager.CloseConnection()
//	log.Printf("Recieved message from client: %v", request.Department)
//
//	var students []sm.Student
//	var departments []um.Department
//	result, err := connection.GetSession().Select("id", "name", "department").From("students").Where("department=?", request.GetDepartment()).Load(&students)
//	log.Printf("result is %v", result)
//	log.Printf("error is %v", err)
//
//	connection.GetSession().Select("id", "name").From("department").Load(&departments)
//	log.Printf("Worked with department and the result is : %v", departments)
//
//	connection.GetSession().Select("id", "name", "department").From("students").Load(&students)
//	log.Printf("Worked with students and the result is : %v", students)
//
//	var department1 um.Department
//	connection.GetSession().Select("id", "name").From("department").Where("id = ?", 3).LoadOne(&department1)
//	log.Printf("Worked with department with id and the result is : %v", department1)
//
//	//var department um.Department
//	//connection.GetSession().Select("id", "name").From("department").Where("id = ?", request.GetId()).LoadOne(&department)
//	//
//	log.Printf("Worked with students and the result is : %v", students)
//	_, err = json.Marshal(&students)
//	if err != nil {
//		log.Fatalf("Error while marshaling %+v", err)
//	}
//	//_, err = json.Marshal(&department)
//
//	var student_list *sm.GetStudentResponse = &sm.GetStudentResponse{}
//
//	for _, st := range students {
//		stu := sm.Student{
//			Id:         st.Id,
//			Name:       st.Name,
//			Department: st.Department,
//		}
//		student_list.Students = append(student_list.Students, &stu)
//	}
//	return student_list, nil
//
//	//studentslice := make([]sm.Student, 0, 2)
//	//for _, st := range (students) {
//	//	studentslice = append(studentslice, sm.Student{
//	//		Id:         st.Id,
//	//		Name:       st.Name,
//	//		Department: st.Department,
//	//	})
//	//}
//
//	//var student_list *sm.GetStudentResponse =&sm.GetStudentResponse{}
//	//student_list.Students=append(student_list.Students,&studentslice)
//	//return &um.GetStudentResponse{Student: &um.Department{
//	//	Id:   department.Id,
//	//	Name: department.Name,
//	//}}, nil
//	//return nil, nil
//
//}
//
//func NewStudentManagementHandler(connectionmanager connection.DatabaseConnectionManager) sm.StudentManagementServiceServer {
//
//	return &studentManagementServer{
//
//		connectionManager: connectionmanager,
//	}
//
//}
