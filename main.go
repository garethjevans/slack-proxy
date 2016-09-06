package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "bytes"
	"io/ioutil"
    "github.com/newrelic/go-agent"
)

func send(r *http.Request, client *http.Client) error {
	url := "https://hooks.slack.com" + r.URL.Path
	log.Println(r.Method + " - " + url)

	body, err := ioutil.ReadAll(r.Body)
	log.Println(string(body))

	if err != nil {
        log.Fatal(err)
        return err
    }

	req, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
    if err != nil {
        log.Fatal(err)
        return err
    }

	resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
        return err
    }

	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Println(string(responseBody))
    if err != nil {
        log.Fatal(err)
        return err
    }

	defer resp.Body.Close()
	return nil
}

func main() {
    nrLicense := os.Getenv("NEW_RELIC_LICENSE_KEY")
    nrName := os.Getenv("NEW_RELIC_APP_NAME")
    port := os.Getenv("PORT")

    config := newrelic.NewConfig(nrName, nrLicense)
    app, err := newrelic.NewApplication(config)

    if err != nil {
        log.Fatal("Unable to create new relic app: ", err)
        os.Exit(1)
    }

    log.Println("Starting application server on port " + port)

    tr := &http.Transport {
        DisableCompression: true,
        DisableKeepAlives: false,
    }

    client := &http.Client{Transport: tr}

    http.HandleFunc(newrelic.WrapHandleFunc(app, "/", func(w http.ResponseWriter, r *http.Request) {
        send(r, client)
        fmt.Fprintf(w,"OK")
    }))

    err2 := http.ListenAndServe(":"+port, nil)
    if err2 != nil {
        log.Fatal("ListenAndServe: ", err2)
        os.Exit(1)
    }
}
