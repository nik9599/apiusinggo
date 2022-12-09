package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// important if we update the record we check for rowAffected
//and if insert  the recor we check for last record inserted
//and for delete recode we check row affected

// setting up conection with mysql

var db *sql.DB

func getMySQLDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@(127.0.0.1:3306)/Practice2?parseTime=true")

	if err != nil {
		log.Fatal(err)
	}
	return db
}

// creating structure for reading the data from database

type student struct {
	Sid    string `json="sid,omitempty"`
	Name   string `json="name,omitempty"`
	Course string `json="course,omitempty"`
}

func main() {

	// declaring router
	r := mux.NewRouter()

	//decalring api

	r.HandleFunc("/students", getStudents).Methods("GET")
	r.HandleFunc("/students", addStudents).Methods("POST")
	r.HandleFunc("/students/{sid}", updateStudents).Methods("PUT")
	r.HandleFunc("/students/{sid}", deleteStudents).Methods("DELETE")

	// starting an server
	http.ListenAndServe(":8080", r)
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	// creating obj of mysql  conection
	db = getMySQLDB()

	rows, err := db.Query("select * from studentinfo1")

	defer db.Close()

	ss := []student{}
	s := student{}

	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {

		for rows.Next() {
			rows.Scan(&s.Sid, &s.Name, &s.Course)
			ss = append(ss, s)
		}
		json.NewEncoder(w).Encode(ss)
	}

}
func addStudents(w http.ResponseWriter, r *http.Request) {

	db = getMySQLDB()
	defer db.Close()

	s := student{}

	json.NewDecoder(r.Body).Decode(&s)

	sid, _ := strconv.Atoi(s.Sid)

	result, err := db.Exec("insert into studentinfo1(sid,name,course) values(?,?,?)", sid, s.Name, s.Course)

	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err := result.LastInsertId()

		if err != nil {
			json.NewEncoder(w).Encode("{error:record not inserted}")
		} else {
			json.NewEncoder(w).Encode(s)
		}
	}

}
func updateStudents(w http.ResponseWriter, r *http.Request) {

	db = getMySQLDB()
	defer getMySQLDB().Close()

	vars := mux.Vars(r) // gettign sid input

	s := student{}

	//fetching information and storing in object of structure and for takin input we use NewDecoder

	json.NewDecoder(r.Body).Decode(&s)

	sid, _ := strconv.Atoi(vars["sid"])

	result, err := db.Exec("update studentinfo1 set name=?,course=? where sid=?", s.Name, s.Course, sid)

	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err := result.RowsAffected()

		if err != nil {
			json.NewEncoder(w).Encode("{error :Record Not Inserted}")
		} else {
			json.NewEncoder(w).Encode(s)
		}
	}

}
func deleteStudents(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	defer getMySQLDB().Close()

	vars := mux.Vars(r) // gettign sid input

	sid, _ := strconv.Atoi(vars["sid"])

	result, err := db.Exec("delete from studentinfo1  where sid=?", sid)

	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err := result.RowsAffected()

		if err != nil {
			json.NewEncoder(w).Encode("{error : data not deleted}")
		} else {
            json.NewEncoder(w).Encode("{ data  deleted}")
		}
	}

}
