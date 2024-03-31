// COVID-19 Vaccination Appointment Application : Built by Pallavi Limaye - 06/03/2021
package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	age "github.com/bearbin/go-age"
)

func datePlusTime(date, timeOfDay string) time.Time {
	if date == "" {
		var t time.Time

		return t
	}
	dateWithTime := date + " " + timeOfDay
	return dateTimeFormat(dateWithTime)
}

func dateTimeFormat(dateWithTime string) time.Time {
	dateTemp, err := time.Parse(layoutDateTime, dateWithTime)
	if err != nil {
		fmt.Println(err)
	}
	return (dateTemp)
}

func convertToInt(stringInput string) int {
	number, _ := strconv.ParseInt(stringInput, 10, 0)
	return int(number)
}

func convertToBool(stringInput string) bool {
	done, _ := strconv.ParseBool(stringInput)
	return done
}

func convertToTime(stringInput string) time.Time {
	t, err := time.Parse(TimeFormatISO, stringInput)

	if err != nil {
		fmt.Println(err)
	}
	return t
}

func convertToTimeFormat(stringInput string) string {
	t, err := time.Parse(TimeFormatISO, stringInput)

	if err != nil {
		fmt.Println(err)
	}
	stringt := t.String()

	return stringt
}

func calculateAge(dob string) int {
	dobTime := convertToTime(dob)
	age := age.Age(dobTime)
	return age
}

/*
func getDateInFormat(year, month, day int) time.Time {
	dob := time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
	return dob
}*/

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func deleteFromApptArray(apptArray []time.Time, deleteThis time.Time) []time.Time {
	for i := 0; i < len(apptArray); i++ {
		arrayTime := apptArray[i]
		if arrayTime.Sub(deleteThis) == 0 {
			return append(apptArray[:i], apptArray[i+1:]...)
		}
	}
	return apptArray
}

func insertApptArray(apptArray []time.Time, addDate time.Time) []time.Time {

	for i := 0; i < len(apptArray); i++ {
		arrayTime := apptArray[i]
		if i+1 < len(apptArray) { // i+1 is not more than array size
			arrayTimeNext := apptArray[i+1]
			if arrayTime.Before(addDate) && arrayTimeNext.After(addDate) {
				tempArray := append(apptArray[:i+1], addDate)
				tempArray = append(tempArray, apptArray[i+1:]...)
				return tempArray
			} else if arrayTime.Equal(addDate) {
				// Duplicate item, no need to add
				return apptArray
			}
		} else { // last item reached
			return append(apptArray, addDate)
		}
	}
	return apptArray
}

func addApptArray(apptArray []time.Time, addDate time.Time) []time.Time {

	for i := 9; i < 18; i++ {
		for j := 0; j <= 45; j = j + 15 {
			addThis := addDate
			addThis = addThis.Add(time.Hour * time.Duration(i))
			addThis = addThis.Add(time.Minute * time.Duration(j))
			apptArray = insertApptArray(apptArray, addThis)
		}
	}
	return apptArray
}

func listApptByDate(apptArray []time.Time) adminStruct {
	adminToTemplate := adminStruct{}

	adminToTemplate.Message = append(adminToTemplate.Message, fmt.Sprintf("Available Appointments by date are: "))

	dateFound := []string{}
	timeNow := time.Now()
	i := 0
	for _, apptItem := range apptArray {
		if timeNow.Before(apptItem) {
			currentDate := apptItem.Format(TimeFormatISO)
			kitchenTime := apptItem.Format(Kitchen)
			if !findInArray(dateFound, currentDate) {
				newTime := currentDate + " => " + kitchenTime
				dateFound = append(dateFound, newTime)
				i++
			} else {
				dateFound[i-1] = dateFound[i-1] + " | " + kitchenTime
			}
		}
	}
	for _, item := range dateFound {
		adminToTemplate.Users = append(adminToTemplate.Users, fmt.Sprintf("%s", item))
	}
	if len(dateFound) == 0 {
		adminToTemplate.Message = append(adminToTemplate.Message, fmt.Sprintf("\nNo available appointments found!\n "))
	}
	return adminToTemplate
}

func findInArray(arrayName []string, arrayValue string) bool {
	for _, item := range arrayName {
		var splititem = strings.Split(item, "=>")
		itemInArray := strings.TrimSpace(splititem[0])
		if itemInArray == arrayValue {
			return true
		}
	}
	return false
}

func printNotQualifiedMessage() []string {
	message := []string{}
	message = append(message, fmt.Sprintf("You do not qualiy for Vaccination at this point. "))
	message = append(message, fmt.Sprintf("Vaccinations are rolled out only for person(s) aged %d and above.\n", ageQualification))
	message = append(message, fmt.Sprintf("Please login after you meet the age requirements,"))
	message = append(message, fmt.Sprintf("or when MOH contacts you on your phone with updated vaccination information."))
	message = append(message, fmt.Sprintf("Please contact our toll-free numbers for any assistance."))
	return message
}

func yesOrNo(status bool) string {
	if status {
		return "Yes"
	} else {
		return "No"
	}
}
