// COVID-19 Vaccination Appointment Application : Built by Pallavi Limaye - 06/03/2021
package main

import (
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"
	bcrypt "golang.org/x/crypto/bcrypt"
)

func (a *appData) registrationHandler(w http.ResponseWriter, req *http.Request) {
	if ok, _ := a.alreadyLoggedIn(req); ok {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	regTemplate := registerError{}
	checkError := ""

	var newPerson person
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		Firstname := req.FormValue("firstname")
		Lastname := req.FormValue("lastname")
		Identification := req.FormValue("identification")
		Username := req.FormValue("username")
		Password := req.FormValue("password")
		Dob := req.FormValue("dob")
		//Dob := "1976-01-03"
		Phone := req.FormValue("phone")
		Address := req.FormValue("address")
		Email := req.FormValue("email")

		newPerson.Identification = Identification
		newPerson.Username = Username
		newPerson.Firstname = Firstname
		newPerson.Lastname = Lastname
		newPerson.Phone = Phone
		newPerson.Dob = Dob
		newPerson.Email = Email
		newPerson.Address = Address

		// Check input data for security
		regTemplate, checkError = newPerson.checkUserData()

		fmt.Println("Error in form : ", checkError)
		if checkError == "nil" { //no error

			idNode := a.bstID.search(Identification)

			if idNode == nil { // ID is not in BST
				newPerson.Identification = Identification
				usernameNode := a.bstUserName.search(Username)

				if usernameNode == nil { // username is not in BST

					newPerson.Username = Username
					pwBytes, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.MinCost)
					if err != nil {
						fmt.Println(err)
					} else {

						newPerson.Password = string(pwBytes)
						// Check if person qualifies for vaccine
						if calculateAge(newPerson.Dob) >= ageQualification {
							newPerson.VaccinationQualify = true
						} else {
							newPerson.VaccinationQualify = false
						}
						// Everything  is OK, add person to linked list and username and ID to the two BSTs
						wg.Add(3)
						go a.personList.addNode(newPerson, &wg)
						go a.bstUserName.insert(Username, &wg)
						go a.bstID.insert(Identification, &wg)
						wg.Wait()
					}
				} else {
					http.Error(w, "Username already taken", http.StatusForbidden)
					return
				}
			} else {
				http.Error(w, "Identification already taken", http.StatusForbidden)
				return
			}
			// redirect to main index
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}
	}
	tpl.ExecuteTemplate(w, "register.gohtml", regTemplate)
}

func (a *appData) loginHandler(w http.ResponseWriter, req *http.Request) {
	if ok, _ := a.alreadyLoggedIn(req); ok {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	loginmessage := ""
	var currentPerson person

	// process form submission
	if req.Method == http.MethodPost {
		Username := req.FormValue("username")
		Password := req.FormValue("password")

		// Check if admin login
		if Username == "admin" {
			err := bcrypt.CompareHashAndPassword([]byte(adminPW), []byte(Password))

			if err != nil {
				//http.Error(w, "Username and/or password do not match", http.StatusUnauthorized)
				//return
				loginmessage = "Username and/or password do not match"
			} else {
				// create session
				id, _ := uuid.NewV4()
				myCookie := &http.Cookie{
					Name:  "myCookie",
					Value: id.String(),
				}
				http.SetCookie(w, myCookie)
				a.mapSessions[myCookie.Value] = Username

				http.Redirect(w, req, "/admin", http.StatusSeeOther)
				return
			}

		} else {

			// check if user exist with username
			usernameNode := a.bstUserName.search(Username)
			if usernameNode == nil { // username is not in BST
				//http.Error(w, "Username and/or password do not match", http.StatusUnauthorized)
				//return
				loginmessage = "Username and/or password do not match"
			} else {
				var err error
				currentPerson, _, err = a.personList.searchUserName(Username)

				if err != nil {
					//http.Error(w, "Username and/or password do not match", http.StatusUnauthorized)
					//return
					loginmessage = "Username and/or password do not match"
				} else {
					err = bcrypt.CompareHashAndPassword([]byte(currentPerson.Password), []byte(Password))

					if err != nil {
						//http.Error(w, "Username and/or password do not match", http.StatusForbidden)
						//return
						loginmessage = "Username and/or password do not match"
					} else {
						// create session

						id, _ := uuid.NewV4()
						myCookie := &http.Cookie{
							Name:  "myCookie",
							Value: id.String(),
						}

						http.SetCookie(w, myCookie)
						a.mapSessions[myCookie.Value] = Username
						http.Redirect(w, req, "/", http.StatusSeeOther)
						return
					}
				}
			}
		}
	}
	tpl.ExecuteTemplate(w, "login.gohtml", loginmessage)
}

func (a *appData) logoutHandler(w http.ResponseWriter, req *http.Request) {
	if ok, path := a.alreadyLoggedIn(req); !ok {
		http.Redirect(w, req, path, http.StatusSeeOther)
		return
	}

	myCookie, _ := req.Cookie("myCookie")

	//var currentPerson person
	if _, ok := a.mapSessions[myCookie.Value]; ok {
		wg.Add(2)
		a.personList.writePersonCSVFile(&wg)
		writeApptCSVFile(&wg, a.apptArray)
		wg.Wait()
	}

	// delete the session
	delete(a.mapSessions, myCookie.Value)

	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, myCookie)

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (a *appData) alreadyLoggedIn(req *http.Request) (bool, string) {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false, "/login"
	}
	username := a.mapSessions[myCookie.Value]
	if username == "admin" {
		return true, "/admin"
	}
	usernameNode := a.bstUserName.search(username)
	if usernameNode == nil { // username is not in BST
		return false, "/login"
	} else {
		return true, "/"
	}
}
