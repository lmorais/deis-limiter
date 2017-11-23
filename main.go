package main

import (
    "fmt"
    "log"
    "sync"
    "os"
    "os/exec"
    "bufio"
)

func main() {
    apps := []string{}

    f, err := os.Open("apps.txt")
    
    if err != nil {
        log.Fatal(err)
    }

    defer f.Close()

    scanner := bufio.NewScanner(f)

    for scanner.Scan() {
        apps = append(apps, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    results := make(chan string)

    var wg sync.WaitGroup

    wg.Add(len(apps))
    
    fmt.Printf("Setting limits for apps: %s \n", apps)

    for _, app := range apps {
        go func(app string) {
            defer wg.Done()
            cmd := exec.Command("deis", "limits:set", "cmd=1G", "-a", app)
            if err := cmd.Run(); err != nil {
                log.Fatal(err)
            } else {
                results <- app
            }
        }(app)
    }

    go func() {
        for result := range results {
            fmt.Printf("App: %s, memory limited to 1G \n", result)
        }
    }()

    wg.Wait()
}