package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Patients struct {
	ID                      int
	Name                    string
	Age                     int
	Gender                  string
	Address                 string
	City                    string
	Phone                   string
	Disease                 string
	Selected_specialisation string
	Patient_history         string
}
type Doctor struct {
	ID                               int
	Name                             string
	Gender                           string
	Address                          string
	City                             string
	Phone                            string
	Specialisation                   string
	Opening_time                     string
	Closing_time                     string
	Availability_time                string
	Availability                     string
	Available_for_home_visit         string
	Available_for_online_consultancy string
	Fees                             int
}

type Database interface {
	AddPatient(p *Patients) error
	GetPatient(p *Patients) (*Patients, error)
	UpdatePatient(p *Patients) error
	DeletePatient(p *Patients) error

	AddDoctor(d *Doctor) error
	GetDoctorByLocation(d *Doctor) (*Doctor, error)
	UpdateDoctor(d *Doctor) error
	DeleteDoctor(d *Doctor) error
}

type HTTPHandler struct {
	function Database
}

type MySQLDatabase struct {
	db *sql.DB
}

func NewMySQLDatabase(connectionString string) (*MySQLDatabase, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &MySQLDatabase{db}, nil
}

// Get Doctor By City and Specialisation - Read Operation

func (m *MySQLDatabase) GetDoctorByLocation(d *Doctor) (*Doctor, error) {

	sql_query := fmt.Sprintf("SELECT * FROM DOCTOR WHERE City='%s' AND Specialisation='%s'", d.City, d.Specialisation)
	details, err := m.db.Query(sql_query)
	if err != nil {
		return nil, err
	}
	defer details.Close()
	var doctor Doctor
	if details.Next() {
		err = details.Scan(&doctor.ID, &doctor.Name, &doctor.Gender, &doctor.Address, &doctor.City, &doctor.Phone, &doctor.Specialisation,
			&doctor.Opening_time, &doctor.Closing_time, &doctor.Availability_time,
			&doctor.Availability, &doctor.Available_for_home_visit, &doctor.Available_for_online_consultancy, &doctor.Fees)
		if err != nil {
			return nil, err
		}
	}
	return &doctor, nil
}

func (h *HTTPHandler) GetDoctorByLocation(c *gin.Context) {

	var doctor Doctor

	err := c.BindJSON(&doctor)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	result, err := h.function.GetDoctorByLocation(&doctor)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"list of doctor as per your choice of location": result})

}

// Add Doctor to Database-  CREATE OPERATION

func (m *MySQLDatabase) AddDoctor(d *Doctor) error {
	sql_query := fmt.Sprintf(`INSERT INTO Doctor (Name,Gender,Address,City,Phone,Specialisation,Opening_time,Closing_time,Availability_time,Availability,Available_for_home_visit,Available_for_online_consultancy,Fees) VALUES ( '%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s',%d)`, d.Name, d.Gender, d.Address, d.City, d.Phone, d.Specialisation, d.Opening_time, d.Closing_time, d.Availability_time, d.Availability, d.Available_for_home_visit, d.Available_for_online_consultancy, d.Fees)
	_, err := m.db.Exec(sql_query)
	return err
}

func (h *HTTPHandler) AddDoctor(c *gin.Context) {
	var doctor Doctor
	if err := c.BindJSON(&doctor); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := h.function.AddDoctor(&doctor); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusCreated, doctor)
}

// Update

func (m *MySQLDatabase) UpdateDoctor(d *Doctor) error {
	update_query := fmt.Sprintf("UPDATE Doctor SET Address='%s',City='%s',Phone='%s', Opening_time ='%s',Closing_time='%s',Fees=%d WHERE Id=%d", d.Address, d.City, d.Phone, d.Opening_time, d.Closing_time, d.Fees, d.ID)
	fmt.Println(update_query)
	_, err := m.db.Exec(update_query)
	return err
}

func (h *HTTPHandler) UpdateDoctort(c *gin.Context) {
	var doctor Doctor
	err := c.BindJSON(&doctor)
	if err != nil {
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	err = h.function.UpdateDoctor(&doctor)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusCreated, doctor)
}

// Delete Doctor to Database-  delete OPERATION

func (m *MySQLDatabase) DeleteDoctor(d *Doctor) error {
	sql_query := fmt.Sprintf("DELETE FROM Doctor WHERE ID= %d", d.ID)
	_, err := m.db.Exec(sql_query)
	return err
}

func (h *HTTPHandler) DeleteDoctor(c *gin.Context) {
	var doctor Doctor
	err := c.BindJSON(&doctor)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.function.DeleteDoctor(&doctor)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusCreated, doctor)
}

//  ################################################### PATIENT ######################################################## //

// Add Patient to Database-  CREATE OPERATION

func (m *MySQLDatabase) AddPatient(p *Patients) error {
	sql_query := fmt.Sprintf(`INSERT INTO patient(Name,Age,Gender,Address,City,Phone,Disease,Selected_Specialisation,Patient_history) 
	VALUES('%s',%d,'%s','%s','%s','%s','%s','%s','%s')`, p.Name, p.Age, p.Gender, p.Address,
		p.City, p.Phone, p.Disease, p.Selected_specialisation, p.Patient_history)
	_, err := m.db.Exec(sql_query)
	return err
}

func (h *HTTPHandler) AddPatient(c *gin.Context) {
	var patient Patients
	if err := c.BindJSON(&patient); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := h.function.AddPatient(&patient); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusCreated, patient)
}

// Get the Patient details from the Database - READ OPERATION

func (m *MySQLDatabase) GetPatient(p *Patients) (*Patients, error) {
	sql_query := fmt.Sprintf(`SELECT * FROM Patient WHERE Phone='%s'`, p.Phone)
	detail, err := m.db.Query(sql_query)
	if err != nil {
		return nil, err
	}
	defer detail.Close()
	var patient Patients
	if detail.Next() {
		err = detail.Scan(&patient.ID, &patient.Name, &patient.Age, &patient.Gender, &patient.Address,
			&patient.City, &patient.Phone, &patient.Disease, &patient.Selected_specialisation,
			&patient.Patient_history)
		if err != nil {
			return nil, err
		}
	}
	return &patient, nil
}

func (h *HTTPHandler) GetPatient(c *gin.Context) {
	var patient Patients
	err := c.BindJSON(&patient)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	result, err := h.function.GetPatient(&patient)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"patient": result})
}

// Update the patient details in the Database - UPDATE OPERATION

func (m *MySQLDatabase) UpdatePatient(p *Patients) error {
	update_query := fmt.Sprintf("UPDATE Patient SET Name='%s',Age=%d,Gender='%s',Address='%s',City='%s',Phone='%s',Diseases='%s',Selected_specialisation='%s',Patient_history='%s' WHERE Id=%d",
		p.Name, p.Age, p.Gender, p.Address, p.City, p.Phone, p.Disease, p.Selected_specialisation, p.Patient_history, p.ID)
	fmt.Println(update_query)
	_, err := m.db.Exec(update_query)
	return err
}

func (h *HTTPHandler) UpdatePatient(c *gin.Context) {
	var patient Patients
	err := c.BindJSON(&patient)
	if err != nil {
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	err = h.function.UpdatePatient(&patient)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Patient is removed from database successfully"})
}

// Delete the patient from Database - DELETE OPERATION

func (m *MySQLDatabase) DeletePatient(p *Patients) error {
	sql_query := fmt.Sprintf("DELETE FROM Patient WHERE Phone='%s'", p.Phone)
	_, err := m.db.Exec(sql_query)
	return err
}

func (h *HTTPHandler) DeletePatient(c *gin.Context) {
	var patient Patients
	err := c.BindJSON(&patient)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.function.DeletePatient(&patient)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusCreated, patient)
}
func Err(err error) {
	if err != nil {
		log.Panic(err.Error())
	}

}

func dbCreation() {

	//connecting to mysql

	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/")
	Err(err)
	defer db.Close()

	// database creation

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS das_new")
	Err(err)
}
func db_connection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/das_new")

	if err != nil {
		return nil, err
	}
	return db, nil
}

func sql_Doctor_tabel_creation() {
	db, err := db_connection()
	Err(err)
	// sql table creation

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Doctor(ID INT NOT NULL AUTO_INCREMENT, Name VARCHAR(30),Gender VARCHAR(10),Address VARCHAR(50), City VARCHAR(20),Phone VARCHAR(15),Specialisation VARCHAR(20),Opening_time VARCHAR(10),Closing_time VARCHAR(10),Availability_time VARCHAR(30),Availability VARCHAR(10),Available_for_home_visit VARCHAR(4),Available_for_online_consultancy VARCHAR(4),Fees INT ,PRIMARY KEY (ID) );")
	Err(err)
	fmt.Println("Docter Table Created")
}

func sql_Patient_tabel_creation() {
	db, err := db_connection()
	Err(err)
	// sql table creation

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS PATIENT(ID INT,Name VARCHAR(20),Age INT, Gender VARCHAR(8),Address VARCHAR(255),City VARCHAR(20,Phone VARCHAR(15),Disease VARCHAR(255),Selected_specialisation VARCHAR(50),Patient_history VARCHAR(255))")
	Err(err)
}

func main() {
	dbCreation()
	sql_Doctor_tabel_creation()
	sql_Patient_tabel_creation()

	db, err := NewMySQLDatabase("root@tcp(localhost:3306)/das_new")
	if err != nil {
		log.Fatal(err)
	}
	defer db.db.Close()

	handler := &HTTPHandler{db}

	router := gin.Default()

	router.GET("doctor/getdoctorbylocation", handler.GetDoctorByLocation)
	router.POST("doctor/add_doctor", handler.AddDoctor)
	router.PUT("doctor/update_doctor", handler.UpdateDoctort)
	router.DELETE("doctor/delete_doctor", handler.DeleteDoctor)

	router.POST("patient/add_patients", handler.AddPatient)
	router.GET("patient/get_patient", handler.GetPatient)
	router.PUT("patient/update_patient", handler.UpdatePatient)
	router.DELETE("patient/delete_patient", handler.DeletePatient)

	router.Run("localhost:8080")
}
