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
        body, _ := io.ReadAll(r.Body)
        
        var data map[string]interface{}
        json.Unmarshal(body, &data)
        
        // Log the tool call
        if params, ok := data["params"].(map[string]interface{}); ok {
            if args, ok := params["arguments"].(map[string]interface{}); ok {
                log.Printf("Calling tool [%v] with arguments: %v", data["tool"], args)
            }
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.Write(body)
    })

    http.HandleFunc("/after", func(w http.ResponseWriter, r *http.Request) {
        body, _ := io.ReadAll(r.Body)
        
        var data map[string]interface{}
        json.Unmarshal(body, &data)
        
        // Count results if available
        if content, ok := data["content"].([]interface{}); ok {
            log.Printf("Tool returned %d results", len(content))
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.Write(body)
    })

    log.Println("Starting interceptor on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
