package main
 
import (
   "fmt"
   "html/template"
   "net/http"
   "TestApp/templates"
   "TestApp/models"
   "strconv"
   "TestApp/utilities"
)

var resultTemplate = template.Must(template.New("result").Parse( templates.ResultsPage))
var errorTemplate = template.Must(template.New("error").Parse( templates.ErrorPage))
var addUserTemplate = template.Must(template.New("error").Parse( templates.AddUserPage))
var allResultTemplate = template.Must(template.New("result").Parse( templates.AllResultsPage))
var dashBoardTemplate = template.Must( template.New("dashboard").Parse( templates.DashBoardPage))

func notify(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, templates.NotifyPage )
}
 
func addUser(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, templates.AddUserPage )
}

func handleStationEvent(w http.ResponseWriter, r *http.Request) {
        errorTemplate.Execute(w, "Event received")
        stationId, _ := strconv.Atoi( r.FormValue("stationId") )
        eventId, _ := strconv.Atoi( r.FormValue("eventId") )
        if stationId != 0 {
                users.HandleEventStation( stationId, eventId )
        } 
}

func showAllUsers(w http.ResponseWriter, r *http.Request) {
        results :=  users.GetAllUsers()
        if results != nil {
                errn := allResultTemplate.Execute( w, results )
                if errn != nil {
                        http.Error(w, errn.Error(), http.StatusInternalServerError)
                } 
        } else {
                errorTemplate.Execute(w, "No entries in database yet")
        }
}

func dashBoard(w http.ResponseWriter, r *http.Request) {
        if r.FormValue("licensePlate") != ""{
                result :=  users.GetUser( r.FormValue("licensePlate") ) 
                if result.Email != "" {
                        users.AddToQueue( r.FormValue("name") )    
                } else {
                        errorTemplate.Execute(w, "Sorry, License Plate hasnt been added to database yet. Please register the vehicle")
                        return 
                }
        }
        dashBoard, err := users.GetDashBoard()
        if err != 0 {
                errn := dashBoardTemplate.Execute( w, dashBoard )
                if errn != nil {
                        http.Error(w, errn.Error(), http.StatusInternalServerError)
                } 
        } else {
                errorTemplate.Execute(w, "No entries in database yet")
        }
}

func result(w http.ResponseWriter, r *http.Request) {
        result :=  users.GetUser( r.FormValue("licensePlate") )
        if result.Email != "" {
                errn := resultTemplate.Execute( w, result )
                email.SendEmail( result.Email, result.LicensePlate )
                if errn != nil {
                        http.Error(w, errn.Error(), http.StatusInternalServerError)
                } 
        } else {
                errorTemplate.Execute(w, "Sorry, License Plate hasnt been added to database yet. Please register the vehicle")
        }
}

func addedUser(w http.ResponseWriter, r *http.Request) {
        result, error :=  users.InsertUser( r.FormValue("licensePlate"), 
                                     r.FormValue("name"),
                                     r.FormValue("email")   )
        if error == -1 {
                errorTemplate.Execute(w, "Duplicate entry( License Plate ) exists in database, User not added ")
                return        
        }
        if result == nil {
                errorTemplate.Execute(w, "User added successfully")
        } else {
                errorTemplate.Execute(w, "User cannon be added, Try again")
        }
}

func main() {
        http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/", dashBoard )
	http.HandleFunc("/result", result)
        http.HandleFunc("/addUser", addUser)
        http.HandleFunc("/addedUser", addedUser)
        http.HandleFunc("/notify", notify )
        http.HandleFunc("/showAllUsers", showAllUsers)
        http.HandleFunc("/handleStationEvent", handleStationEvent)
	err := http.ListenAndServe( ":9999", nil)
	if err != nil {
	        panic(err)
	}
}