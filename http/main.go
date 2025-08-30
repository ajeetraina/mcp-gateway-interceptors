package main

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    "os"
    "time"
)

func main() {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    http.HandleFunc("/before", func(w http.ResponseWriter, r *http.Request) {
        body, _ := io.ReadAll(r.Body)
        log.Printf("Calling tool with arguments: %s", string(body))
        
        var data map[string]interface{}
        json.Unmarshal(body, &data)
        
        // Extract and log the query
        if params, ok := data["params"].(map[string]interface{}); ok {
            if args, ok := params["arguments"].(map[string]interface{}); ok {
                log.Printf("Calling tool [%v] with arguments: %v", data["tool"], args)
            }
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.Write(body) // Pass through
    })

    http.HandleFunc("/after", func(w http.ResponseWriter, r *http.Request) {
        body, _ := io.ReadAll(r.Body)
        log.Printf("Tool response: %d bytes", len(body))
        
        w.Header().Set("Content-Type", "application/json")
        w.Write(body) // Pass through
    })

    port := "8080"
    log.Printf("Starting interceptor on port %s", port)
    http.ListenAndServe(":"+port, nil)
}
