package main

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "sync"
)

type ArtifactResponse struct {
    Items            []ArtifactItem `json:"items"`
    ContinuationToken string        `json:"continuationToken,omitempty"`
}

type ArtifactItem struct {
    Path string `json:"path"`
}

var (
    httpClient = &http.Client{
        Timeout: 30 * time.Second,
    }
    sourceNexus    = "http://source-nexus-url"
    sourceRepo     = "source-repo"
    sourceUser     = "source-username"
    sourcePassword = "source-password"
    targetNexus    = "http://target-nexus-url"
    targetRepo     = "target-repo"
    targetUser     = "target-username"
    targetPassword = "target-password"
    artifactsFile  = "artifacts_list.txt"
)

func main() {
    fetchArtifacts()
    migrateArtifacts()
    fmt.Println("Migration completed.")
}

func fetchArtifacts() {
    var continuationToken string

    file, err := os.Create(artifactsFile)
    if err != nil {
        fmt.Printf("Error creating file: %v\n", err)
        return
    }
    defer file.Close()

    for {
        url := fmt.Sprintf("%s/service/rest/v1/search/assets?repository=%s&continuationToken=%s", sourceNexus, sourceRepo, continuationToken)
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            fmt.Printf("Error creating request: %v\n", err)
            return
        }
        req.SetBasicAuth(sourceUser, sourcePassword)

        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            fmt.Printf("Error on request: %v\n", err)
            return
        }
        if resp.Body != nil {
            defer resp.Body.Close()
        }

        var response ArtifactResponse
        if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
            fmt.Printf("Error decoding response: %v\n", err)
            return
        }

        for _, item := range response.Items {
            _, err := file.WriteString(item.Path + "\n")
            if err != nil {
                fmt.Printf("Error writing to file: %v\n", err)
                return
            }
        }

        continuationToken = response.ContinuationToken
        if continuationToken == "" {
            break
        }
    }
}

func migrateArtifacts() {
    file, err := os.Open(artifactsFile)
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        return
    }
    defer file.Close()
    defer os.Remove(artifactsFile)

    var wg sync.WaitGroup
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        path := scanner.Text()
        wg.Add(1)
        go func(p string) {
            defer wg.Done()
            if err := migrateArtifact(p); err != nil {
                fmt.Printf("Failed to migrate %s: %v\n", p, err)
            } else {
                fmt.Printf("Successfully migrated %s\n", p)
            }
        }(path)
    }
    wg.Wait()

    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading file: %v\n", err)
    }
}


func migrateArtifact(path string) error {
    sourceURL := fmt.Sprintf("%s/repository/%s/%s", sourceNexus, sourceRepo, path)
    req, err := http.NewRequest("GET", sourceURL, nil)
    if err != nil {
        return fmt.Errorf("error creating the GET request: %w", err)
    }
    req.SetBasicAuth(sourceUser, sourcePassword)

    resp, err := doRequest(req)
    if err != nil {
        return fmt.Errorf("error executing the GET request: %w", err)
    }
    defer resp.Body.Close()

    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("error reading the response body: %w", err)
    }


    targetURL := fmt.Sprintf("%s/repository/%s/%s", targetNexus, targetRepo, path)
    req, err = http.NewRequest("PUT", targetURL, bytes.NewReader(data))
    if err != nil {
        return fmt.Errorf("error creating the PUT request: %w", err)
    }
    req.SetBasicAuth(targetUser, targetPassword)

    resp, err = doRequest(req)
    if err != nil {
        return fmt.Errorf("error executing the PUT request: %w", err)
    }
    defer resp.Body.Close()

    return nil
}


func doRequest(req *http.Request) (*http.Response, error) {
    var resp *http.Response
    var err error
    for i := 0; i < 3; i++ {
        resp, err = httpClient.Do(req)
        if err == nil && resp.StatusCode < 500 {
            return resp, nil
        }
        if resp != nil {
            resp.Body.Close()
        }
        time.Sleep(time.Duration(2^(i+1)) * time.Second)
    }
    return nil, err
}
