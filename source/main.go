// COVID-19 Vaccination Appointment Application : Built by Pallavi Limaye - 06/03/2021
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	TimeFormatISO    = "2006-01-02"
	TimeFormatUS     = "January 02, 2006"
	TimeFormatTime   = "15:04:05"
	layoutDateTime   = "2006-01-02 15:04:05"
	layoutDateTimeUS = "January 02, 2006 3:4 pm"
	Kitchen          = "3:04PM"

	//Age Criteria for qualifying for vaccination
	ageQualification = 70

	// number of available appointments shown at a time
	numOfApptsShown = 15

	// admin login is admin, admin password is "goisthebest"
	adminName = "admin"
	adminPW   = "$2a$04$WfaVDeh35aAFJdEgk9wsG.QZkKB463q/tAOruw0qlx6R8bwPuZMeW"

	filePath           = "./"
	csvPath            = "./csv/"
	htmlPath           = "./html/"
	certPath           = "./cert/"
	apptFileName       = "freeappointments.csv"
	personDataFileName = "personData.csv"
)

var (
	//loc = time.FixedZone("UTC+8", +8*60*60)

	// wg is used to wait for the program to finish.
	wg  sync.WaitGroup
	tpl *template.Template

	mu sync.Mutex
)

type appData struct {
	personList  *linkedList
	bstUserName *BST
	bstID       *BST
	apptArray   []time.Time
	mapSessions map[string]string
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	// Initialize data structures and appt array
	appDS := appData{}
	appDS.mapSessions = make(map[string]string)

	// read CSV files and write to data structures concurrently
	appDS.apptArray = readApptCSVFile(&wg)
	// Read all names and their data from CSV file
	appDS.personList, appDS.bstUserName, appDS.bstID = readpersonCSVFile(&wg)

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("./Assets"))))

	http.HandleFunc("/", appDS.indexHandler)
	http.HandleFunc("/register", appDS.registrationHandler)
	http.HandleFunc("/login", appDS.loginHandler)
	http.HandleFunc("/logout", appDS.logoutHandler)
	http.HandleFunc("/viewappt", appDS.viewapptHandler)
	http.HandleFunc("/makeappt", appDS.makeapptHandler)
	http.HandleFunc("/deleteappt", appDS.deleteapptHandler)
	http.HandleFunc("/profile", appDS.profileHandler)
	http.HandleFunc("/admin", appDS.adminHandler)
	http.HandleFunc("/listallusers", appDS.listallusersHandler)
	http.HandleFunc("/deleteuser", appDS.deleteuserHandler)
	http.HandleFunc("/viewapptsbydate", appDS.viewapptsbydateHandler)
	http.HandleFunc("/addapptsfordate", appDS.addapptsfordateHandler)

	fmt.Println("Listening....")
	//log.Fatal(http.ListenAndServe(":5221", nil))
	log.Fatal(http.ListenAndServeTLS(":5221", certPath+"cert.pem", certPath+"key.pem", nil))

}

func (a *appData) indexHandler(w http.ResponseWriter, req *http.Request) {
	currentPerson := a.getPerson(w, req)
	tpl.ExecuteTemplate(w, "index.gohtml", currentPerson)
}
