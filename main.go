package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error
 

type Photo struct {
	ID			int				`json:"id"`
	Title		string			`json:"title"`
	Caption		string			`json:"caption"`
	PhotoUrl	string			`json:"photourl"`
}

type Result struct {
	Code		int			`json:"code"`
	Data		interface{}	`json:"data"`
	Message		string		`json:"message"`
}

func main() {
	
    db, err = gorm.Open("mysql", "root:@/photo?charset=utf8&parseTime=True")

    if err != nil {
    	log.Println("Koneksi Gagal! error: ", err)
    } else { 
    	log.Println("Koneksi Berhasil!")
    }

		db.AutoMigrate(&Photo{})

		handleRequests()
}


func handleRequests() {
	log.Println("Start the development server http:localhost:9999")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/photos", createPhotos).Methods("POST")
	myRouter.HandleFunc("/photos", getPhotos).Methods("GET")
	myRouter.HandleFunc("/{id}", getPhoto).Methods("GET")
	myRouter.HandleFunc("/{id}", updatePhoto).Methods("PUT")
	myRouter.HandleFunc("/{id}", deletePhoto).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9999", myRouter))
}


func createPhotos(w http.ResponseWriter, r *http.Request) {

	payloads, _ := ioutil.ReadAll(r.Body)

	var photo Photo

	json.Unmarshal(payloads, &photo)

	db.Create(&photo)

	res:= Result{Code: 200, Data: photo, Message: "Berhasil menambahkan photo!"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}



func getPhotos(w http.ResponseWriter, r *http.Request) {

	photos := []Photo{}

	db.Find(&photos)

	res := Result{Code: 200, Data: photos, Message: "Berhasil menampilkan semua photo!"}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}


func getPhoto(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	photoID := vars["id"]

	var photo Photo
	db.First(&photo, photoID)

	res := Result{Code: 200, Data: photo, Message: "Berhasil menampilkan single photo!"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}


func updatePhoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	photoID := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)

	var photoUpdates Photo

	json.Unmarshal(payloads, &photoUpdates)

	var photo Photo
	db.First(&photo, photoID)

	db.Model(&photo).Update(photoUpdates)

	res:= Result{Code: 200, Data: photo, Message: "Berhasil mengubah photo!"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}



func deletePhoto(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	photoID := vars["id"]

	var photo Photo

	db.First(&photo, photoID)
	db.Delete(&photo)

	res := Result{Code: 200, Message: "Berhasil menghapus photo!"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}