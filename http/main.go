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
        
        // Extract query
        query := "unknown"
        tool := "unknown"
        if t, ok := data["tool"].(string); ok {
            tool = t
        }
        if params, ok := data["params"].(map[string]interface{}); ok {
            if args, ok := params["arguments"].(map[string]interface{}); ok {
                if q, ok := args["query"].(string); ok {
                    query = q
                }
            }
        }
        
        log.Printf("============================================")
        log.Printf("[INTERCEPTOR CHAIN - BEFORE]")
        log.Printf("Tool: %s", tool)
        log.Printf("Query: %s", query)
        log.Printf("============================================")
        
        // Pass through unchanged
        w.Header().Set("Content-Type", "application/json")
        w.Write(body)
    })

    http.HandleFunc("/after", func(w http.ResponseWriter, r *http.Request) {
        body, _ := io.ReadAll(r.Body)
        
        var data map[string]interface{}
        json.Unmarshal(body, &data)
        
        resultCount := 0
        if content, ok := data["content"].([]interface{}); ok {
            resultCount = len(content)
            
            // Show first result title if available
            if resultCount > 0 {
                if firstResult, ok := content[0].(map[string]interface{}); ok {
                    if title, ok := firstResult["title"].(string); ok {
                        log.Printf("[AFTER] First result: %s", title)
                    }
                }
            }
        }
        
        log.Printf("============================================")
        log.Printf("[INTERCEPTOR CHAIN - AFTER]")
        log.Printf("Results returned: %d", resultCount)
        log.Printf("============================================")
        
        // Pass through unchanged
        w.Header().Set("Content-Type", "application/json")
        w.Write(body)
    })

    log.Println("🚀 Docker MCP Gateway Interceptor Started")
    log.Println("📍 Ready to intercept tool calls")
    http.ListenAndServe(":8080", nil)
}
