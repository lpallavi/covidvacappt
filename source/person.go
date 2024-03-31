// COVID-19 Vaccination Appointment Application : Built by Pallavi Limaye - 06/03/2021
package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type makeapptstruct struct {
	Firstname     string
	Lastname      string
	Message       []string
	Apptmade      string
	Possibleappts []string
}

type deleteapptstruct struct {
	Firstname  string
	Lastname   string
	Message    []string
	Apptdelete string
}
type registerError struct {
	Firstname      string
	Lastname       string
	Identification string
	Username       string
	Password       string
	Dob            string
	Phone          string
	Address        string
	Email          string
}

type person struct {
	Identification     string
	Username           string
	Password           string
	Firstname          string
	Lastname           string
	Dob                string
	Phone              string
	Address            string
	Email              string
	VaccinationQualify bool
	FirstVaccineDone   bool
	FirstVaccineDate   string
	FirstVaccineTime   string
	SecondVaccineDone  bool
	SecondVaccineDate  string
	SecondVaccineTime  string
}

func (p *person) printUserCredentials() (string, string) {
	return p.Username, p.Password
}

func (p *person) checkUserData() (registerError, string) {
	regTemplate := registerError{}
	checkError := "nil"

	id := regexp.QuoteMeta(p.Identification)

	if id != p.Identification {
		checkError = "Special characters in Identification number"
		regTemplate.Identification = "No special characters allowed in Identification number"
	} else {
		// First and last character must be an alphabet
		idMatch := regexp.MustCompile(`^[a-zA-Z]+.*[a-zA-Z]+$`)
		if !idMatch.MatchString(id) {
			regTemplate.Identification = "Incorrect Singapore Identification number"
			checkError = "Error in Identification number"
		} else {
			if len(id) != 9 {
				regTemplate.Identification = "Incorrect Singapore Identification number"
				checkError = "Error in Identification number"
			} else {

			}

		}
	}

	return regTemplate, checkError
}

func (p *person) printAllInfo() {
	fmt.Println("Your Profile Information:")
	fmt.Println("First Name:\t\t\t", p.Firstname)
	fmt.Println("Last Name:\t\t\t", p.Lastname)
	fmt.Println("Username:\t\t\t", p.Username)
	fmt.Println("Identification:\t\t\t", p.Identification)
	fmt.Println("Date Of Birth:\t\t\t", p.Dob)
	fmt.Println("Phone:\t\t\t\t", p.Phone)
	fmt.Println("Email:\t\t\t\t", p.Email)
	fmt.Println("Address:\t\t\t", p.Address)
	fmt.Println("Vaccination Qualification:\t", yesOrNo(p.VaccinationQualify))
	fmt.Println("First Vaccination Complete:\t", yesOrNo(p.FirstVaccineDone))
	if p.FirstVaccineDate != "" {
		fmt.Println("First Vaccination Date:\t\t", p.FirstVaccineDate)
		fmt.Println("First Vaccination Time:\t\t", p.FirstVaccineTime)
	}
	fmt.Println("Second Vaccination Complete:\t", yesOrNo(p.SecondVaccineDone))
	if p.SecondVaccineDate != "" {
		fmt.Println("Second Vaccination Date:\t", p.SecondVaccineDate)
		fmt.Println("Second Vaccination Time:\t", p.SecondVaccineTime)
	}
}

func (p *person) printAllAppt() []string {
	message := []string{}
	message = append(message, fmt.Sprintf("Appointment information for %s %s\n", p.Firstname, p.Lastname))
	timeFirst := datePlusTime(p.FirstVaccineDate, p.FirstVaccineTime)
	timeSecond := datePlusTime(p.SecondVaccineDate, p.SecondVaccineTime)
	timeFirst29 := timeFirst.Add(time.Hour * 24 * time.Duration(29))
	timeNow := time.Now()

	if timeFirst.IsZero() {
		message = append(message, fmt.Sprintf("You have no appointments yet. Please make new appointments."))
	} else {
		if timeFirst.After(time.Now()) {
			message = append(message, fmt.Sprintf("First Vaccine Appointment is on: %s", timeFirst.Format(layoutDateTime)))
		} else {
			if p.FirstVaccineDone {
				message = append(message, fmt.Sprintf("First Vaccine was given on: %s", timeFirst.Format(layoutDateTime)))

			} else {
				message = append(message, fmt.Sprintf("You missed your first Vaccine appointment dated: %s", timeFirst.Format(layoutDateTime)))
				message = append(message, fmt.Sprintf("Please delete first appointment and make new appointments"))
			}
		}

		if timeSecond.IsZero() {
			if p.FirstVaccineDone {
				//check if first vaccine is  done more than 29 days before now
				if timeNow.After(timeFirst29) {
					message = append(message, fmt.Sprintf("28 days have already elapsed after your first vaccination dose."))
					message = append(message, fmt.Sprintf("You can no longer make a second vaccination appointment. Contact MOH immediately!"))
					return message
				}
			} else {
				message = append(message, fmt.Sprintf("You have no appointment for second vaccine. Please make new appointment."))
			}
		} else {
			if timeSecond.After(timeNow) {
				message = append(message, fmt.Sprintf("Second Vaccine Appointment is on: %s", timeSecond.Format(layoutDateTime)))
			} else {
				if p.SecondVaccineDone {
					message = append(message, fmt.Sprintf("Second Vaccine was given on: %s", timeSecond.Format(layoutDateTime)))
				} else {
					message = append(message, fmt.Sprintf("You missed your second Vaccine appointment dated: %s", timeSecond.Format(layoutDateTime)))
					if p.FirstVaccineDone {
						if timeNow.After(timeFirst29) {
							message = append(message, fmt.Sprintf("28 days have already elapsed after your first vaccination dose."))
							message = append(message, fmt.Sprintf("You can no longer make a second vaccination appointment. Contact MOH immediately!"))
						} else {
							message = append(message, fmt.Sprintf("You can still make a new appointment for second vaccination."))
							message = append(message, fmt.Sprintf("Delete second appointment before making new appointment"))
							message = append(message, fmt.Sprintf("Hurry up! Slots are limited. Or contact MOH immediately!"))
						}
					}
				}
			}
		}
	}
	return message
}

func (p *person) makeNewAppt(apptArray []time.Time) ([]time.Time, makeapptstruct) {

	makeapptmessage := makeapptstruct{}
	makeapptmessage.Firstname = p.Firstname
	makeapptmessage.Lastname = p.Lastname

	timeFirst := datePlusTime(p.FirstVaccineDate, p.FirstVaccineTime)
	timeSecond := datePlusTime(p.SecondVaccineDate, p.SecondVaccineTime)
	timeFirst20 := timeFirst.Add(time.Hour * 24 * time.Duration(20))
	timeFirst29 := timeFirst.Add(time.Hour * 24 * time.Duration(29))
	timeNow := time.Now()
	// Apptmade : none = no appt being made, first = first appt being made, second = second appt being made
	makeapptmessage.Apptmade = "none"

	if p.FirstVaccineDone {
		if p.SecondVaccineDone {
			makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("Your COVID-19 Vaccination is complete. Take good rest and be safe."))
			return apptArray, makeapptmessage
		} else {
			// new appt for second vaccine. Must be 21-28 days within first vaccine
			// else contact MOH hotline
			if timeSecond.IsZero() {
				makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("Your first vaccination was done on %s", timeFirst.Format(layoutDateTime)))
				makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("Available dates and time for next appointment are:"))
				// Apptmade : none = no appt being made, first = first appt being made, second = second appt being made
				makeapptmessage.Apptmade = "second"

				counter := 0
				for _, apptItem := range apptArray {
					if timeNow.Before(apptItem) {
						if timeFirst29.After(apptItem) && timeFirst20.Before(apptItem) {
							counter++
							makeapptmessage.Possibleappts = append(makeapptmessage.Possibleappts, fmt.Sprintf("%s", apptItem.Format(layoutDateTime)))
						}
						if counter >= numOfApptsShown {
							break
						}
					}
				}

				if counter == 0 {
					makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("\nNo available appointments found! Please contact MOH Hotline\n "))
					return apptArray, makeapptmessage
				} else {
					makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("Please choose your appointment date/time:"))
					return apptArray, makeapptmessage
				}
			} else {
				makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("Your first vaccination was done on %s", timeFirst.Format(layoutDateTime)))
				if timeNow.After(timeSecond) {
					makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("You missed your second appointment on %s", timeSecond.Format(layoutDateTime)))
				} else {
					makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("You already have an appointment on %s", timeSecond.Format(layoutDateTime)))
				}
				makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("Delete this appointment before making new appointment"))
				return apptArray, makeapptmessage
			}
		}
	} else {
		// make both appts together - difference of 21-28 days

		//no first appt made
		if timeFirst.IsZero() {
			makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("Available dates and time for first appointment are:"))
			// Apptmade : none = no appt being made, first = first appt being made, second = second appt being made
			makeapptmessage.Apptmade = "first"
			counter := 0
			for _, apptItem := range apptArray {
				if timeNow.Before(apptItem) {
					counter++
					makeapptmessage.Possibleappts = append(makeapptmessage.Possibleappts, fmt.Sprintf("%s", apptItem.Format(layoutDateTime)))
				}
				if counter >= numOfApptsShown {
					break
				}
			}
			if counter == 0 {
				makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("\nNo available appointments found! Please contact MOH Hotline\n "))
				return apptArray, makeapptmessage
			} else {
				makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("Please choose your appointment date/time:"))
				return apptArray, makeapptmessage
			}
		} else { // first appt done. check if second appt is done
			if timeSecond.IsZero() {
				// For second appointment after choosing first appointment date
				if timeNow.After(timeFirst) {
					makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("You missed your first appointment on %s", timeFirst.Format(layoutDateTime)))
					makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("Delete this appointment, and make new appointments"))
				} else {
					makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("You first vaccination appointment is on %s", timeFirst.Format(layoutDateTime)))
					makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("Available dates and time for second appointment are:"))
					// Apptmade : none = no appt being made, first = first appt being made, second = second appt being made
					makeapptmessage.Apptmade = "second"

					counter := 0
					for _, apptItem := range apptArray {
						if timeNow.Before(apptItem) {
							if timeFirst29.After(apptItem) && timeFirst20.Before(apptItem) {
								counter++
								makeapptmessage.Possibleappts = append(makeapptmessage.Possibleappts, fmt.Sprintf("%s", apptItem.Format(layoutDateTime)))
							}
							if counter >= numOfApptsShown {
								break
							}
						}
					}
					if counter == 0 {
						makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("\nNo available appointments found! Please contact MOH Hotline\n "))
						return apptArray, makeapptmessage
					} else {
						makeapptmessage.Message = append(makeapptmessage.Message, fmt.Sprintf("Please choose your appointment date/time:"))
						return apptArray, makeapptmessage
					}
				}
			}
		}
	}
	return apptArray, makeapptmessage
}

func (p *person) updateNewAppt(apptTime string, apptFor string, apptArray []time.Time) []time.Time {

	timeOfAppt := dateTimeFormat(apptTime) // convert string to time.Time format

	apptArray = deleteFromApptArray(apptArray, timeOfAppt)
	d := timeOfAppt.Format(layoutDateTime)
	dateArray := strings.Split(d, " ")

	if apptFor == "first" {
		p.FirstVaccineDate = dateArray[0]
		p.FirstVaccineTime = dateArray[1]
		p.FirstVaccineDone = false
	} else if apptFor == "second" {
		p.SecondVaccineDate = dateArray[0]
		p.SecondVaccineTime = dateArray[1]
		p.SecondVaccineDone = false
	}
	return apptArray
}

func (p *person) deleteAppt() deleteapptstruct {
	timeFirst := datePlusTime(p.FirstVaccineDate, p.FirstVaccineTime)
	timeSecond := datePlusTime(p.SecondVaccineDate, p.SecondVaccineTime)
	timeNow := time.Now()

	deleteapptmessage := deleteapptstruct{}
	deleteapptmessage.Firstname = p.Firstname
	deleteapptmessage.Lastname = p.Lastname

	if p.FirstVaccineDone { // cannot delete if done
		if p.SecondVaccineDone { // cannot delete if done
			deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("Your COVID-19 Vaccination is complete. Take good rest and be safe."))
		} else {
			//no second appt made
			if timeSecond.IsZero() {
				deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("Your first vaccination was done on %s", timeFirst.Format(layoutDateTime)))
				deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("You have no appointment for second vaccination. Please make new appointment"))
			} else {
				deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("Your first vaccination was done on %s", timeFirst.Format(layoutDateTime)))

				if timeNow.After(timeSecond) {
					deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("You missed your second appointment dated %s", timeSecond.Format(layoutDateTime)))
				} else {
					deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("Your second appointment is scheduled for %s", timeSecond.Format(layoutDateTime)))
				}

				deleteapptmessage.Apptdelete = "second"
				deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("Click confirm to delete second appointment"))

			}
		}
	} else {
		//no first appt made
		if timeFirst.IsZero() {
			//no second appt made
			if timeSecond.IsZero() {
				deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("You have no appointments made. Please make new appointments"))
			} else {
				deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("You have no appointment for first vaccination. Deleting second appointment"))
				//corner case, will not happen if everything goes well
				p.SecondVaccineDate = ""
				p.SecondVaccineTime = ""
				p.SecondVaccineDone = false
			}
		} else {
			//no second appt made
			if timeSecond.IsZero() {
				if timeNow.After(timeFirst) {
					deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("You missed your first appointment dated %s", timeFirst.Format(layoutDateTime)))
				} else {
					deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("Your first appointment is scheduled for %s", timeFirst.Format(layoutDateTime)))
				}
				deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("You have no appointments for second vaccination."))
				deleteapptmessage.Apptdelete = "first"
				deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("Click confirm to delete first appointment"))

			} else {
				if timeNow.After(timeFirst) {
					deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("You missed your first appointment dated %s", timeFirst.Format(layoutDateTime)))
				} else {
					deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("Your first appointment is scheduled for %s", timeFirst.Format(layoutDateTime)))
				}
				if timeNow.After(timeSecond) {
					deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("You missed your second appointment dated %s", timeSecond.Format(layoutDateTime)))
				} else {
					deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("Your second appointment is scheduled for %s", timeSecond.Format(layoutDateTime)))
				}
				deleteapptmessage.Apptdelete = "both"
				deleteapptmessage.Message = append(deleteapptmessage.Message, fmt.Sprintf("Click confirm to delete both appointments"))

			}
		}
	}
	return deleteapptmessage
}

func (p *person) updatedeleteAppt(appttoDelete string, apptArray []time.Time) []time.Time {

	timeFirst := datePlusTime(p.FirstVaccineDate, p.FirstVaccineTime)
	timeSecond := datePlusTime(p.SecondVaccineDate, p.SecondVaccineTime)
	timeNow := time.Now()

	if appttoDelete == "first" {
		p.FirstVaccineDate = ""
		p.FirstVaccineTime = ""
		p.FirstVaccineDone = false
		//if appt is after now, add to apptArray
		if timeNow.Before(timeFirst) {
			apptArray = insertApptArray(apptArray, timeFirst)
		}

	} else if appttoDelete == "second" {
		p.SecondVaccineDate = ""
		p.SecondVaccineTime = ""
		p.SecondVaccineDone = false
		//if appt is after now, add to apptArray
		if timeNow.Before(timeSecond) {
			apptArray = insertApptArray(apptArray, timeSecond)
		}
	} else if appttoDelete == "both" {
		p.FirstVaccineDate = ""
		p.FirstVaccineTime = ""
		p.FirstVaccineDone = false
		p.SecondVaccineDate = ""
		p.SecondVaccineTime = ""
		p.SecondVaccineDone = false
		//if appt is after now, add to apptArray
		if timeNow.Before(timeFirst) {
			apptArray = insertApptArray(apptArray, timeFirst)
		}
		//if appt is after now, add to apptArray
		if timeNow.Before(timeSecond) {
			apptArray = insertApptArray(apptArray, timeSecond)
		}

	}
	return apptArray
}
