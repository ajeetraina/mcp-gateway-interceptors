package main

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    http.HandleFunc("/before", func(w http.ResponseWriter, r *http.Request) {
        body, err := io.ReadAll(r.Body)
        if err != nil {
            log.Printf("Error reading request: %v", err)
            http.Error(w, "Error", http.StatusInternalServerError)
            return
        }
        
        // Log but don't modify
        var data map[string]interface{}
        if err := json.Unmarshal(body, &data); err == nil {
            if params, ok := data["params"].(map[string]interface{}); ok {
                if args, ok := params["arguments"].(map[string]interface{}); ok {
                    log.Printf("Calling tool [%v] with arguments: %v", data["tool"], args)
                }
            }
        }
        
        // CRITICAL: Pass through unchanged
        w.Header().Set("Content-Type", "application/json")
        w.Write(body)
    })

    http.HandleFunc("/after", func(w http.ResponseWriter, r *http.Request) {
        body, err := io.ReadAll(r.Body)
        if err != nil {
            log.Printf("Error reading response: %v", err)
            http.Error(w, "Error", http.StatusInternalServerError)
            return
        }
        
        // Log but don't modify
        var data map[string]interface{}
        if err := json.Unmarshal(body, &data); err == nil {
            if content, ok := data["content"].([]interface{}); ok {
                log.Printf("Tool returned %d results", len(content))
            }
        }
        
        // CRITICAL: Pass through unchanged
        w.Header().Set("Content-Type", "application/json")
        w.Write(body)
    })

    log.Println("Starting interceptor on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
