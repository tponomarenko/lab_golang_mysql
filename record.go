package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Record struct {
	Id          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

type RecordRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

func newRecordFromRequest(request *RecordRequest) *Record {
	return &Record{
		Id:          uuid.Must(uuid.NewUUID()).String(),
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		PhoneNumber: request.PhoneNumber,
	}
}

type Endpoint struct {
	db *sql.DB
}

func (this *Endpoint) getRecordById(recordId string) (*Record, error) {
	results, err := this.db.Query("SELECT * FROM records WHERE id = ?;", recordId)
	if err != nil {
		log.Print(err.Error())
		return nil, fmt.Errorf("cannot get record from database")
	}

	if results.Next() {

		var record Record

		_ = results.Scan(&record.Id, &record.FirstName, &record.LastName, &record.PhoneNumber)

		return &record, nil
	} else {
		return nil, nil
	}
}

func (this *Endpoint) getAllRecords() (*[]Record, error) {
	results, err := this.db.Query("SELECT * FROM records;")
	if err != nil {
		log.Print(err.Error())
		return nil, fmt.Errorf("cannot get records from database")
	}

	records := []Record{}
	for results.Next() {
		var record Record

		_ = results.Scan(&record.Id, &record.FirstName, &record.LastName, &record.PhoneNumber)
		records = append(records, record)
	}

	return &records, nil
}

func (this *Endpoint) saveRecord(record *Record) error {
	_, err := this.db.Query(
		`INSERT INTO records (id, first_name, last_name, phone_number) VALUES (?, ?, ?, ?) AS new
                   ON DUPLICATE KEY UPDATE first_name=new.first_name, last_name=new.last_name, phone_number=new.phone_number;`,
		record.Id, record.FirstName, record.LastName, record.PhoneNumber,
	)

	if err != nil {
		log.Print(err.Error())
		return fmt.Errorf("cannot save record to database")
	}

	return nil
}

func (this *Endpoint) deleteRecord(recordId string) error {
	_, err := this.db.Query("DELETE FROM records WHERE id = ?;", recordId)

	if err != nil {
		log.Print(err.Error())
		return fmt.Errorf("cannot delete record from database")
	}

	return nil
}

func (this *Endpoint) GetRecords(ctx *gin.Context) {
	allRecords, err := this.getAllRecords()
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, &gin.H{"message": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, allRecords)
}

func (this *Endpoint) GetRecordById(ctx *gin.Context) {
	recordId := ctx.Param("recordId")

	record, err := this.getRecordById(recordId)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, &gin.H{"message": err.Error()})
		return
	}

	if record == nil {
		ctx.IndentedJSON(http.StatusNotFound, &gin.H{"message": "record not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, record)
}

func (this *Endpoint) AddRecord(ctx *gin.Context) {
	var request RecordRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		log.Println(err.Error())
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "could not parse submitted record"})
		return
	}

	record := newRecordFromRequest(&request)
	err = this.saveRecord(record)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, &gin.H{"message": err.Error()})
		return
	}

	savedRecord, _ := this.getRecordById(record.Id)
	ctx.IndentedJSON(http.StatusCreated, savedRecord)
}

func (this *Endpoint) DeleteRecord(ctx *gin.Context) {
	recordId := ctx.Param("recordId")

	record, err := this.getRecordById(recordId)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, &gin.H{"message": err.Error()})
		return
	}

	if record == nil {
		ctx.IndentedJSON(http.StatusNotFound, &gin.H{"message": "record not found"})
		return
	}

	if this.deleteRecord(recordId) != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, &gin.H{"message": err.Error()})
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}

func (this *Endpoint) UpdateRecord(ctx *gin.Context) {
	recordId := ctx.Param("recordId")

	oldRecord, err := this.getRecordById(recordId)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, &gin.H{"message": err.Error()})
		return
	}

	if oldRecord == nil {
		ctx.IndentedJSON(http.StatusNotFound, &gin.H{"message": "record not found"})
		return
	}

	var record Record
	err = ctx.BindJSON(&record)
	if err != nil {
		log.Println(err.Error())
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "could not parse submitted record"})
		return
	}

	if recordId != record.Id {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "cannot update record id"})
		return
	}

	if this.saveRecord(&record) != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, &gin.H{"message": err.Error()})
		return
	}

	savedRecord, _ := this.getRecordById(record.Id)
	ctx.IndentedJSON(http.StatusOK, savedRecord)
}

func NewEndpoint(settings *Settings) (*Endpoint, error) {
	var connectionString string
	var driverName string

	if settings.DbEngine == "postgresql" {
		connectionString = fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			settings.DbUsername, settings.DbPassword, settings.DbHost, settings.DbPort, settings.DbName,
		)
		driverName = "postgresql"
	} else {
		connectionString = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			settings.DbUsername, settings.DbPassword, settings.DbHost, settings.DbPort, settings.DbName,
		)
		driverName = "mysql"
	}

	db, err := sql.Open(driverName, connectionString)

	if err != nil {
		log.Printf("Cannot connect to the database at %s.\n", connectionString)
		return nil, err
	}

	return &Endpoint{db: db}, nil
}
