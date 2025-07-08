package main

import (
	"strings"
	"time"
)

var allTasksCompleted = true
func main(){
	go checkIfTasksCompleted()

	//go QueryServer()
	ProxyServer()
}

func checkIfTasksCompleted(){
	for true {
		currentDate := time.Now().Unix()
		
		uncompletedTaskFound := false

		for _,agendaItem := range GetAgendaItems() {
			itemType := agendaItem.AgendaItem[0]
			itemText := agendaItem.AgendaItem[1]

			var scheduledStr string
			if agendaItem.Properties["Scheduled"]!="" {
				scheduledStr = agendaItem.Properties["Scheduled"]
			}

			//Only Works With uncompleted tasks with blck tag
			if (itemType != "CANCELLED" && itemType != "INFO" && itemType != "DONE" && itemType != "DUE") &&
			(strings.Contains(itemText, "#blck") || strings.Contains(itemText, ":blck:")) {
				//If current time is bigger than the scheduled date (task is not completed)
				//or task does not have a scheduled property
				if scheduledStr=="" || (scheduledStr!="" && ConvertToUnixEpoch(scheduledStr) < int(currentDate)) {
					uncompletedTaskFound = true
					break
				}			
			}
		}

		if uncompletedTaskFound {allTasksCompleted=false}else{allTasksCompleted=true}

		//Wait for 5 minutes
		time.Sleep(time.Minute*3)
	}
}
