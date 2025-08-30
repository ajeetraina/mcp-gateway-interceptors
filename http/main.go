package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "time"
)

func main() {
    // Health check endpoint
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // Before interceptor
    http.HandleFunc("/before", func(w http.ResponseWriter, r *http.Request) {
        body, err := io.ReadAll(r.Body)
        if err != nil {
            log.Printf("❌ Error reading request: %v", err)
            http.Error(w, "Error reading request", http.StatusBadRequest)
            return
        }

        log.Printf("📥 BEFORE Interceptor received %d bytes", len(body))
        
        var data map[string]interface{}
        if err := json.Unmarshal(body, &data); err != nil {
            log.Printf("⚠️ JSON parse error: %v", err)
            // Pass through anyway
            w.Header().Set("Content-Type", "application/json")
            w.Write(body)
            return
        }

        // Extract query
        query := "unknown"
        if params, ok := data["params"].(map[string]interface{}); ok {
            if args, ok := params["arguments"].(map[string]interface{}); ok {
                if q, ok := args["query"].(string); ok {
                    query = q
                }
            }
        }

        log.Printf("🔵 BEFORE: Tool=%v, Query='%s'", data["tool"], query)
        
        // Add timestamp
        data["intercepted_before"] = time.Now().Format(time.RFC3339)
        
        response, _ := json.Marshal(data)
        w.Header().Set("Content-Type", "application/json")
        w.Write(response)
    })

    // After interceptor
    http.HandleFunc("/after", func(w http.ResponseWriter, r *http.Request) {
        body, err := io.ReadAll(r.Body)
        if err != nil {
            log.Printf("❌ Error reading response: %v", err)
            http.Error(w, "Error reading response", http.StatusBadRequest)
            return
        }

        log.Printf("📤 AFTER Interceptor received %d bytes", len(body))
        
        var data map[string]interface{}
        if err := json.Unmarshal(body, &data); err != nil {
            log.Printf("⚠️ JSON parse error in after: %v", err)
            // Pass through anyway
            w.Header().Set("Content-Type", "application/json")
            w.Write(body)
            return
        }

        // Count results
        resultCount := 0
        if content, ok := data["content"].([]interface{}); ok {
            resultCount = len(content)
        }

        log.Printf("🟢 AFTER: Processed %d results", resultCount)
        
        // Add metadata
        data["intercepted_after"] = time.Now().Format(time.RFC3339)
        data["result_count"] = resultCount
        
        response, _ := json.Marshal(data)
        w.Header().Set("Content-Type", "application/json")
        w.Write(response)
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("🚀 Interceptor server starting on port %s", port)
    log.Printf("📍 Endpoints: /health, /before, /after")
    
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
}
