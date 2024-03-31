// COVID-19 Vaccination Appointment Application : Built by Pallavi Limaye - 06/03/2021
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func writeApptCSVFile(wg *sync.WaitGroup, apptArray []time.Time) {
	defer wg.Done()

	defer func() {
		if r := recover(); r != nil {
			println("Panic:" + r.(string))
		}
	}()
	mu.Lock()
	{
		csvfile, err := os.Create(filePath + csvPath + apptFileName)
		defer csvfile.Close()

		if err != nil {
			log.Fatalf("failed creating file: %s", err)
			panic(fmt.Sprintf("Unable to open Appointment Available CSV file for writing"))
			return
		}

		csvwriter := csv.NewWriter(csvfile)

		for _, appt := range apptArray {

			timeNow := time.Now()
			if timeNow.Before(appt) {
				d := appt.Format(layoutDateTime)
				dateArray := strings.Split(d, " ")

				if err := csvwriter.Write(dateArray); err != nil {
					log.Fatalln("error writing record to csv:", err)
					return
				}
			}
		}
		csvwriter.Flush()
		if err := csvwriter.Error(); err != nil {
			log.Fatal(err)
		}
	}
	mu.Unlock()

}

func readApptCSVFile(wg *sync.WaitGroup) []time.Time {
	var apptArray []time.Time

	csvFile, err := os.Open(filePath + csvPath + apptFileName)

	if err != nil {
		panic(fmt.Sprintf("Unable to open Appointment Available CSV file"))
		return apptArray
	}

	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()

	if err != nil {
		panic(fmt.Sprintf("Unable to read Appointment Available CSV file"))
		return apptArray
	}

	var newApptAvailable time.Time
	var newDate, newTime string
	for _, line := range csvLines {
		newDate = line[0]
		newTime = line[1]
		newApptAvailable = datePlusTime(newDate, newTime)
		apptArray = append(apptArray, newApptAvailable)
	}
	return apptArray
}

func readpersonCSVFile(wg *sync.WaitGroup) (*linkedList, *BST, *BST) {
	readList := &linkedList{head: nil, size: 0}
	bstUserName := &BST{root: nil}
	bstID := &BST{root: nil}

	csvFile, err := os.Open(filePath + csvPath + personDataFileName)
	if err != nil {
		panic(fmt.Sprintf("Unable to open Person Data CSV file"))
		return readList, bstUserName, bstID
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		panic(fmt.Sprintf("Unable to read Person Data CSV file"))
		return readList, bstUserName, bstID
	}

	for _, line := range csvLines {
		i := 0
		onePerson := person{
			Identification:     line[i],
			Username:           line[i+1],
			Password:           line[i+2],
			Firstname:          line[i+3],
			Lastname:           line[i+4],
			Dob:                line[i+5],
			Phone:              line[i+6],
			Address:            line[i+7],
			Email:              line[i+8],
			VaccinationQualify: convertToBool(line[i+9]),
			FirstVaccineDone:   convertToBool(line[i+10]),
			FirstVaccineDate:   line[i+11],
			FirstVaccineTime:   line[i+12],
			SecondVaccineDone:  convertToBool(line[i+13]),
			SecondVaccineDate:  line[i+14],
			SecondVaccineTime:  line[i+15],
		}

		// Starting goroutines for three functions.
		//Adding nodes to the linked list, and usernames and identification to BST

		wg.Add(3)
		go readList.addNode(onePerson, wg)
		go bstUserName.insert(onePerson.Username, wg)
		go bstID.insert(onePerson.Identification, wg)
		wg.Wait()

	}
	return readList, bstUserName, bstID
}

func (p *linkedList) writePersonCSVFile(wg *sync.WaitGroup) {
	defer wg.Done()

	defer func() {
		if r := recover(); r != nil {
			println("Panic:" + r.(string))
		}
	}()
	mu.Lock()
	{
		csvfile, err := os.Create(filePath + csvPath + personDataFileName)
		defer csvfile.Close()

		if err != nil {
			log.Fatalf("failed creating file: %s", err)
			panic(fmt.Sprintf("Unable to create Person Data CSV file for writing"))
			return
		}

		csvwriter := csv.NewWriter(csvfile)

		personllhead := p.head
		personllsize := p.size
		if personllsize != 0 && personllhead != nil {
			for i := 1; i < personllsize+1; i++ {
				oneperson, _ := p.get(i)

				var line []string

				line = append(line, oneperson.Identification, oneperson.Username, oneperson.Password,
					oneperson.Firstname, oneperson.Lastname, oneperson.Dob, oneperson.Phone,
					oneperson.Address, oneperson.Email, strconv.FormatBool(oneperson.VaccinationQualify),
					strconv.FormatBool(oneperson.FirstVaccineDone), oneperson.FirstVaccineDate,
					oneperson.FirstVaccineTime, strconv.FormatBool(oneperson.SecondVaccineDone),
					oneperson.SecondVaccineDate, oneperson.SecondVaccineTime)

				if err := csvwriter.Write(line); err != nil {
					log.Fatalln("error writing record to csv:", err)
					return
				}
			}
			csvwriter.Flush()
			if err := csvwriter.Error(); err != nil {
				log.Fatal(err)
				return
			}

		}
	}
	mu.Unlock()
}
