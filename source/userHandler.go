// COVID-19 Vaccination Appointment Application : Built by Pallavi Limaye - 06/03/2021
package main

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func (a *appData) viewapptHandler(w http.ResponseWriter, req *http.Request) {
	currentPerson := a.getPerson(w, req)
	message := []string{}
	//io.WriteString(w, currentPerson.Firstname)
	if currentPerson.VaccinationQualify {
		message = currentPerson.printAllAppt()
	} else {
		message = printNotQualifiedMessage()
	}
	tpl.ExecuteTemplate(w, "viewappt.gohtml", message)
}

func (a *appData) makeapptHandler(w http.ResponseWriter, req *http.Request) {
	currentPerson := a.getPerson(w, req)

	makeapptmessage := makeapptstruct{}

	if currentPerson.VaccinationQualify {
		a.apptArray, makeapptmessage = currentPerson.makeNewAppt(a.apptArray)
	} else {
		makeapptmessage.Message = printNotQualifiedMessage()
	}

	appointmentDate := ""
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		appointmentDate = req.FormValue("appttime")

		a.apptArray = currentPerson.updateNewAppt(appointmentDate, makeapptmessage.Apptmade, a.apptArray)

		//update person
		a.updatePerson(currentPerson, w, req)

		http.Redirect(w, req, "/index", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "makeappt.gohtml", makeapptmessage)
}

func (a *appData) deleteapptHandler(w http.ResponseWriter, req *http.Request) {
	currentPerson := a.getPerson(w, req)

	deleteapptmessage := deleteapptstruct{}

	if currentPerson.VaccinationQualify {
		deleteapptmessage = currentPerson.deleteAppt()
	} else {
		deleteapptmessage.Message = printNotQualifiedMessage()
	}
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		_ = req.FormValue("submit")
		a.apptArray = currentPerson.updatedeleteAppt(deleteapptmessage.Apptdelete, a.apptArray)

		//update person
		a.updatePerson(currentPerson, w, req)

		http.Redirect(w, req, "/index", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "deleteappt.gohtml", deleteapptmessage)
}

func (a *appData) profileHandler(w http.ResponseWriter, req *http.Request) {
	currentPerson := a.getPerson(w, req)
	tpl.ExecuteTemplate(w, "profile.gohtml", currentPerson)
}

func (a *appData) getPerson(w http.ResponseWriter, req *http.Request) person {
	// get current session cookie
	myCookie, err := req.Cookie("myCookie")
	if err != nil { //cannot find cookie, create new cookie
		id, _ := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}
	}
	http.SetCookie(w, myCookie)

	// if the user exists already, get user
	var currentPerson person
	if username, ok := a.mapSessions[myCookie.Value]; ok {
		if username != "admin" {
			currentPerson, _, _ = a.personList.searchUserName(username)
		}
	}
	return currentPerson
}

func (a *appData) updatePerson(p person, w http.ResponseWriter, req *http.Request) {
	// get current session cookie
	myCookie, err := req.Cookie("myCookie")
	if err != nil { //cannot find cookie, create new cookie
		id, _ := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}
	}
	http.SetCookie(w, myCookie)

	// if the user exists already, get user
	if username, ok := a.mapSessions[myCookie.Value]; ok {
		if username != "admin" {
			currentPerson, index, _ := a.personList.searchUserName(username)
			if p.Username == currentPerson.Username {
				a.personList.writePersonData(p, index)
			}
		}
	}
}
