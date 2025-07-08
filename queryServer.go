package main

import (
	"fmt"
	"net/http"
)

func queryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Set the response header to plaintext
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if allTasksCompleted {
			fmt.Fprintf(w, "tasks-completed")
		} else {
			fmt.Fprintf(w, "tasks-uncompleted")
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func QueryServer() {
	var err error
	http.HandleFunc("/", queryHandler) // Set the handler for the root path
	fmt.Println("Query server is listening on port "+queryPort)
	if certPath != "" && keyPath != "" {
		err = http.ListenAndServeTLS(":"+queryPort, certPath,keyPath, nil) // Start the server
	}else{
		err = http.ListenAndServe(":"+queryPort, nil) // Start the server
	}
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
