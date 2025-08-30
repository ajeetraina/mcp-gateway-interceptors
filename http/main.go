package main

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/before", func(w http.ResponseWriter, r *http.Request) {
        body, _ := io.ReadAll(r.Body)
        
        var data map[string]interface{}
        json.Unmarshal(body, &data)
        
        // Log the call
        if params, ok := data["params"].(map[string]interface{}); ok {
            if args, ok := params["arguments"].(map[string]interface{}); ok {
                log.Printf("[HTTP-BEFORE] Tool: %v, Query: %v", 
                    data["tool"], args["query"])
            }
        }
        
        // Pass through unchanged
        w.Header().Set("Content-Type", "application/json")
        w.Write(body)
    })

    http.HandleFunc("/after", func(w http.ResponseWriter, r *http.Request) {
        body, _ := io.ReadAll(r.Body)
        
        var data map[string]interface{}
        json.Unmarshal(body, &data)
        
        // Log results
        if content, ok := data["content"].([]interface{}); ok {
            log.Printf("[HTTP-AFTER] Results: %d", len(content))
        }
        
        // Pass through unchanged
        w.Header().Set("Content-Type", "application/json")
        w.Write(body)
    })

    log.Println("Interceptor server started on :8080")
    http.ListenAndServe(":8080", nil)
}
