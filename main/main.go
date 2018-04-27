package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "encoding/json"
)

//Notes:
///users/id
//Under each unique user(collection), there is only one item(document), the ToDoList struct (which only has a single array of strings)
var session *mgo.Session
var sErr error
//Struct for storing in database 
type ToDoList struct {
	Username string `json:"username"`
	Items []string `json:"items"`
}



//Handler function for basic request
func BaseEndPoint(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Home Page!")
}

func GetList(w http.ResponseWriter, r *http.Request){ //Reads list from the desired user in the database

	//The desired user id will be passed in the http request
	vars := mux.Vars(r)
	id := vars["userID"]

	//Copy new connection to the database
	newSession := session.Copy()
	defer newSession.Close()
	//Get the lists collection
	c := newSession.DB("todolist").C("lists")

	var list ToDoList

	err := c.Find(bson.M{"username": id}).One(&list)
	if err != nil {
        respondJsonError(w, "DB error", http.StatusInternalServerError)
        fmt.Println("failed to find ",id)
        log.Println("Failed find list: ", err)
        return
    }
    if list.Username == "" {
        respondJsonError(w, "List not found", http.StatusNotFound)
        return
    }

    respBody, err := json.MarshalIndent(list, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    //Respond to request (using fxn)
	respondJson(w, respBody, http.StatusOK)
	
}

func CreateList(w http.ResponseWriter, r *http.Request){ //Writes new list to the desired user in the database
	
	fmt.Println("create list called\n")
	//Copy new connection to the database
	
	newSession := session.Copy()
	defer newSession.Close()
	
	//Read in the desired list post from request
	var list ToDoList
	decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&list)
    if err != nil {
        respondJsonError(w, "Error reading body", http.StatusBadRequest)
        fmt.Println(err)
        return
    }
    //Insert the list into the collection
    c := newSession.DB("todolist").C("lists")
    err = c.Insert(list)
    if err != nil {
        if mgo.IsDup(err) {
            respondJsonError(w, "List already exists", http.StatusBadRequest)
            return
        }
        respondJsonError(w, "Database error", http.StatusInternalServerError)
        log.Println("Error inserting list: ", err)
        return
    }
    //Respond to request
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Location", r.URL.Path+"/"+list.Username)
    w.WriteHeader(http.StatusCreated)


}

func UpdateList(w http.ResponseWriter, r *http.Request){
	fmt.Println("update list called\n")

	//The desired user id will be passed in the http request
	vars := mux.Vars(r)
	id := vars["userID"]

	//Copy new connection to the database
	newSession := session.Copy()
	defer newSession.Close()

	//Read in the desired list post from request
	var list ToDoList
	decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&list)
    if err != nil {
        respondJsonError(w, "Error reading body", http.StatusBadRequest)
        fmt.Println(err)
        return
    }

    //Update the list in the collection
    c := newSession.DB("todolist").C("lists")
    err = c.Update(bson.M{"username": id}, &list)
    if err != nil {
        switch err {
            default:
                respondJsonError(w, "DB error", http.StatusInternalServerError)
                log.Println("Failed update list: ", err)
                return
            case mgo.ErrNotFound:
                respondJsonError(w, "list not found", http.StatusNotFound)
                return
        }
    }
    //Simple response to request
    w.WriteHeader(http.StatusNoContent)

}

func DeleteList(w http.ResponseWriter, r *http.Request){
	//The desired user id will be passed in the http request
	vars := mux.Vars(r)
	id := vars["userID"]

	//Copy new connection to the database
	newSession := session.Copy()
	defer newSession.Close()

	//Update the list into the collection
    c := newSession.DB("todolist").C("lists")
    err := c.Remove(bson.M{"username": id})
    if err != nil {
        switch err {
            default:
                respondJsonError(w, "DB error", http.StatusInternalServerError)
                log.Println("Failed update list: ", err)
                return
            case mgo.ErrNotFound:
                respondJsonError(w, "list not found", http.StatusNotFound)
                return
        }
    }
    //Respond to request
    w.WriteHeader(http.StatusNoContent)

}



func ensureIndex(s *mgo.Session) {  
    session := s.Copy()
    defer session.Close()

    c := session.DB("todolist").C("lists")

    index := mgo.Index{
        Key:        []string{"username"},
        Unique:     true,
        DropDups:   true,
        Background: true,
        Sparse:     true,
    }
    err := c.EnsureIndex(index)
    if err != nil {
        panic(err)
    }
}

//Respond to a request successfully using this
func respondJson(w http.ResponseWriter, json []byte, code int){
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(code)
    w.Write(json)
}
//Respond to request with an err using this
func respondJsonError(w http.ResponseWriter, message string, code int) {  
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(code)
    fmt.Fprintf(w, "{message: %q}", message)
}


func main(){
	//Show program start in terminal
	fmt.Print("running \n")

	//Get inital database session
    session, sErr = mgo.Dial("localhost")
    if sErr != nil {
        panic(sErr)
    }
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)

    //Make sure there are unique indexes
    ensureIndex(session) 

    //Routing
	r := mux.NewRouter()
	r.HandleFunc("/", BaseEndPoint).Methods("GET")
	r.HandleFunc("/user/{userID}",GetList).Methods("GET")
	r.HandleFunc("/user",CreateList).Methods("POST")
	r.HandleFunc("/user/{userID}",UpdateList).Methods("PUT")
	r.HandleFunc("/user/{userID}",DeleteList).Methods("DELETE")


	if err := http.ListenAndServe(":3000", r); err != nil{
		log.Fatal(err)
	}



}







