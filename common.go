package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var parsedAllowBeforeTime = func()int{
	parsedTime, err := time.Parse("15:04", allowBefore)
	if err != nil {
		log.Fatal(err)
	}
	return parsedTime.Hour()+parsedTime.Minute()
}()
var parsedBlockAfterTime = func()int{
	parsedTime, err := time.Parse("15:04", blockAfter)
	if err != nil {
		log.Fatal(err)
	}
	return parsedTime.Hour()+parsedTime.Minute()
}()

func isPageBlocked(hostname string) (blocked bool){
	if allowBefore != "" {
		currentDate := time.Now()
		currentTimeInDay := currentDate.Hour()+currentDate.Minute()
		if currentTimeInDay > parsedAllowBeforeTime {
			return true
		}else{return false}
	}
	if blockAfter != "" {
		currentDate := time.Now()
		currentTimeInDay := currentDate.Hour()+currentDate.Minute()
		if currentTimeInDay > parsedBlockAfterTime {
			return true
		}
	}

	if allTasksCompleted {return false}

	var match = false

	for _,blockedDomain := range blockedDomains {
		if strings.Contains(hostname, blockedDomain){
			match = true; break
		}
	}

	if (match && blockType=="blacklist") ||  (!match && blockType=="whitelist") {
		return true
	}
	return false
}

func expandPath(path string) (string, error) {
	// Check if the path starts with a tilde
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		// Replace the tilde with the home directory
		path = strings.Replace(path, "~", homeDir, 1)
	}

	// Get the absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// readLines reads a file and returns its content as a slice of strings.
func readLines(filename string) ([]string, error) {
	var lines []string
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

//timeStr format must be a yyyy-mm-dd with an optional hh:mm
func ConvertToUnixEpoch(dateTimeStr string) int {
	dateTimeStr = extractDate(dateTimeStr)
	timeStr := extractTime(dateTimeStr)
	if timeStr == "" {
		dateTimeStr = dateTimeStr+" 00:00"
	}else{dateTimeStr = dateTimeStr+" "+timeStr}

	parsedTime, err := time.Parse("2006-01-02 15:04", dateTimeStr)
	if err != nil {
		return 0
	}

	return int(parsedTime.Unix())
}
