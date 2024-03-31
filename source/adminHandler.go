// COVID-19 Vaccination Appointment Application : Built by Pallavi Limaye - 06/03/2021
package main

import (
	"fmt"
	"net/http"
)

func (a *appData) adminHandler(w http.ResponseWriter, req *http.Request) {

	if ok, path := a.alreadyLoggedIn(req); !ok {
		http.Redirect(w, req, path, http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "admin.gohtml", nil)
}

func (a *appData) listallusersHandler(w http.ResponseWriter, req *http.Request) {

	if ok, path := a.alreadyLoggedIn(req); !ok {
		http.Redirect(w, req, path, http.StatusSeeOther)
		return
	}
	adminToTemplate := a.personList.printAllUsers()
	adminToTemplate.Deleteuser = "no"
	tpl.ExecuteTemplate(w, "admin.gohtml", adminToTemplate)
}

func (a *appData) deleteuserHandler(w http.ResponseWriter, req *http.Request) {

	if ok, path := a.alreadyLoggedIn(req); !ok {
		http.Redirect(w, req, path, http.StatusSeeOther)
		return
	}
	adminToTemplate := a.personList.printAllUsers()
	adminToTemplate.Deleteuser = "yes"

	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		userToDeleteIndex := req.FormValue("user")
		// add one to userToDeleteIndex, as index in form is a range starting from 0
		personRemoved, _ := a.personList.remove(convertToInt(userToDeleteIndex) + 1)
		wg.Add(2)
		go a.bstUserName.delete(personRemoved.Username, &wg)
		go a.bstID.delete(personRemoved.Identification, &wg)
		wg.Wait()

		http.Redirect(w, req, "/admin", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "admin.gohtml", adminToTemplate)
}

func (a *appData) viewapptsbydateHandler(w http.ResponseWriter, req *http.Request) {

	if ok, path := a.alreadyLoggedIn(req); !ok {
		http.Redirect(w, req, path, http.StatusSeeOther)
		return
	}
	adminToTemplate := listApptByDate(a.apptArray)
	adminToTemplate.ApptAdd = "no"

	tpl.ExecuteTemplate(w, "admin.gohtml", adminToTemplate)
}

func (a *appData) addapptsfordateHandler(w http.ResponseWriter, req *http.Request) {

	if ok, path := a.alreadyLoggedIn(req); !ok {
		http.Redirect(w, req, path, http.StatusSeeOther)
		return
	}
	adminToTemplate := adminStruct{}
	adminToTemplate.Message = append(adminToTemplate.Message, fmt.Sprintf("Choose a date to add appointments from 9:00am upto 5:45pm"))
	adminToTemplate.ApptAdd = "yes"

	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		addtoThisDate := req.FormValue("appt")
		addThisDate := convertToTime(addtoThisDate)
		a.apptArray = addApptArray(a.apptArray, addThisDate)
		http.Redirect(w, req, "/admin", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "admin.gohtml", adminToTemplate)
}
