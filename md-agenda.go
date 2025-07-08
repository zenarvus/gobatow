package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func listAgendaFiles() []string {
	var aFiles []string
	for _,agendaFile := range agendaFiles {
		agendaFile,err := expandPath(agendaFile)
		if err!=nil{log.Fatal(err)}

		fileInfo,err := os.Stat(agendaFile)
		if err!=nil{log.Fatal(err)}

		if !fileInfo.IsDir() {
			aFiles = append(aFiles, agendaFile)
		//TOOD: Implement folder support
		}else{continue}
		
	}
	return aFiles
}

// parseTaskLine extracts the task type and title from a line.
func parseTaskLine(line string) (string, string) {
	re := regexp.MustCompile(`^#+ (.+): (.*)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 3 {
		return matches[1], matches[2]
	}
	return "", ""
}

type AgendaItem struct {
	Metadata    []any
	AgendaItem  []string
	Properties  map[string]string
	LogbookItems map[string][]string
}
func GetAgendaItems() []AgendaItem {
	/*{
		{
			metadata={filePath, lineNumber}
			agendaItem={type, text, fullLine}
			properties={key=value}
			logbookItems={{item type, time, progress}, ...}
		},
		{...},
		...
	}*/
	var agendaItems []AgendaItem

	for _,agendaFilePath := range listAgendaFiles(){
		fileContent, err := readLines(agendaFilePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", agendaFilePath, err)
			continue
		}

		if fileContent != nil {
			for lineNumber, line := range fileContent {
				taskType, title := parseTaskLine(line)
				if taskType != "" && title != "" {
					agendaItem := AgendaItem{
						Metadata:   []interface{}{agendaFilePath, lineNumber + 1},
						AgendaItem: []string{taskType, title, line},
					}

					agendaItem.Properties = getTaskProperties(fileContent, lineNumber+1)

					// Check for repeat indicators in Scheduled or Deadline properties
					if (agendaItem.Properties["Scheduled"] != "" && matchRepeatIndicator(agendaItem.Properties["Scheduled"])) ||
						(agendaItem.Properties["Deadline"] != "" && matchRepeatIndicator(agendaItem.Properties["Deadline"])) {
						agendaItem.LogbookItems = getLogbookEntries(fileContent, lineNumber+1)
					}

					agendaItems = append(agendaItems, agendaItem)
				}
			}
		}
	}

	return agendaItems
}

// getTaskProperties retrieves task properties from content lines.
func getTaskProperties(contentLines []string, taskLineNum int) map[string]string {
	properties := make(map[string]string)

	//It starts right after the task line
	propertyLineNum := taskLineNum + 1
	currentLine := 0

	for _, line := range contentLines {
		currentLine++

		if currentLine < propertyLineNum {
			continue
		}

		propertyPattern := regexp.MustCompile("^ *- (.+): `(.*)`")
		matches := propertyPattern.FindStringSubmatch(line)
		if len(matches) == 3 {
			key := matches[1]
			value := matches[2]

			properties[key] = value

			propertyLineNum++
		} else {break}
	}

	return properties
}

// getLogbookEntries retrieves logbook entries from content lines.
func getLogbookEntries(contentLines []string, taskLineNum int) map[string][]string {
	entries := make(map[string][]string)

	logbookStartPassed := false
	lineNumber := 0

	for _, line := range contentLines {
		lineNumber++

		// Skip task headline
		if lineNumber < taskLineNum+1 {
			continue
		}

		if strings.Contains(line, "<details logbook>") {
			logbookStartPassed = true
		}

		if logbookStartPassed {
			// Example logbook line: - DONE: `2022-12-30 18:00` `(6/10)`
			status, text := parseLogbookLine(line)
			if status != "" && text != "" {
				log := []string{status}

				time := extractLogbookTime(text)
				if time == "" {
					continue
				}
				log = append(log, time)

				progressIndicator := extractProgressIndicator(text)
				if progressIndicator != "" {
					log = append(log, progressIndicator)
				}

				date := extractDate(time)
				entries[date] = log
			}
		}

		// Stop when arriving at another header or logbook's end
		if isHeader(line) || strings.Contains(line, "</details>") {
			break
		}
	}

	return entries
}

// parseLogbookLine extracts the status and text from a logbook line.
func parseLogbookLine(line string) (string, string) {
	re := regexp.MustCompile(`^\s*-\s*([A-Z]+):\s*(.*)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 3 {
		return matches[1], matches[2]
	}
	return "", ""
}

// extractTime extracts the date and time from the text in logbook.
func extractLogbookTime(text string) string {
	re := regexp.MustCompile("`([0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+)`")
	matches := re.FindStringSubmatch(text)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}

// extractProgressIndicator extracts the progress indicator from the text.
func extractProgressIndicator(text string) string {
	re := regexp.MustCompile("`(\\([0-9]+/[0-9]+\\))`")
	matches := re.FindStringSubmatch(text)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}

// extractDate extracts the date from the time string.
func extractDate(time string) string {
	re := regexp.MustCompile("([0-9]+-[0-9]+-[0-9]+)")
	matches := re.FindStringSubmatch(time)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}
// extractDate extracts the hh:mm time from the time string.
func extractTime(time string) string {
	re := regexp.MustCompile("([0-9]+:[0-9]+)")
	matches := re.FindStringSubmatch(time)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}

// isHeader checks if the line is a header.
func isHeader(line string) bool {
	re := regexp.MustCompile(`^#+ .*`)
	return re.MatchString(line)
}

// matchRepeatIndicator checks if the string has a repeat indicator
func matchRepeatIndicator(line string) bool {
	re := regexp.MustCompile(" [\\.\\+]+[0-9]+[a-z]")
	return re.MatchString(line)
}
