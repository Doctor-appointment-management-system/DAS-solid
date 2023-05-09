package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

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
	update_query := fmt.Sprintf("UPDATE Doctor SET Address='%s',City='%s',Phone='%s', Opening_time ='%s',Closing_time='%s',Fees %d, WHERE Id=%d", d.Address, d.City, d.Phone, d.Opening_time, d.Closing_time, d.Fees, d.ID)
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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS PATIENT(ID INT,Name VARCHAR(20),Gender VARCHAR(8),Address VARCHAR(255),City VARCHAR(20),State VARCHAR(20),Mobile_no VARCHAR(15))")
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

	handler := &HTTPHandler{}

	router := gin.Default()

	router.GET("doctor/getdoctorbylocation", handler.GetDoctorByLocation)

	router.POST("doctor/add_doctor", handler.AddDoctor)

	router.PUT("doctor/update_doctor", handler.UpdateDoctort)
	router.DELETE("doctor/delete_doctor", handler.DeleteDoctor)
	router.Run("localhost:8080")
}
