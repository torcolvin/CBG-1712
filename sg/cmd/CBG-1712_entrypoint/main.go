package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/cenkalti/backoff/v4"
)

const (
	adminUsername = "Administrator"
	adminPassword = "password"
	adminURL      = "http://localhost:4985"
)

func main() {
	log.Default().SetPrefix("SYNC_GATEWAY_ENTRYPOINT ")
	sgExited := make(chan struct{})
	runSG(sgExited)
	<-sgExited
}

func runSG(sgExited chan struct{}) {
	cmd := exec.Command("/CBG-1712/sync_gateway/sync_gateway", "/CBG-1712/sync_gateway_config.json")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatalf("sync_gateway failed to launch +%v", err)
	}
	waitForServerUp()
	setupDB()
	configureDB()
	setupUser()
	primeCache()
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("sync_gateway exited non-zero %+v", err)

	}
	sgExited <- struct{}{}
}

func waitForServerUp() {
	log.Printf("Make sure admin interface is online")
	waitForGetURL(adminURL)
	log.Printf("admin interface is online")
}

func waitForGetURL(url string) {
	getOp := func() error {
		client := http.Client{Timeout: 5 * time.Second}
		req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
		if err != nil {
			return err
		}
		req.SetBasicAuth(adminUsername, adminPassword)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("%s GET return %d: %s", url, resp.StatusCode, resp.Body)
		}
		return nil
	}
	err := backoff.Retry(getOp, backoff.NewExponentialBackOff())
	if err != nil {
		log.Fatalf("URL%s GET with error %s", url, err)
	}
	err = getOp()
	if err != nil {
		log.Fatalf("URL%s GET with error %s", url, err)
	}

}

func putAdmin(url string, payload *bytes.Buffer) {
	doHttp(http.MethodPut, url, payload)
}

//func postAdmin(url string, payload *bytes.Buffer) {
//	doHttp(http.MethodPost, url, payload)
//}

func doHttp(method, url string, payload *bytes.Buffer) {
	client := http.Client{Timeout: 50 * time.Second}
	var err error
	var req *http.Request
	if payload != nil {
		req, err = http.NewRequest(method, url, payload)
	} else {
		req, err = http.NewRequest(method, url, http.NoBody)

	}
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(adminUsername, adminPassword)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close() //nolint
	if err != nil {
		log.Fatalf("Error %s to %s with %s", method, url, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading body")
	}
	if resp.StatusCode >= 400 {
		log.Fatalf("%s %s return %d: %s", url, method, resp.StatusCode, body)
	}
	log.Printf("%s %s return %d", url, method, resp.StatusCode)
	waitForGetURL(url)
}

func setupDB() {
	log.Print("Setting up a SG db")
	url := adminURL + "/db/"
	putAdmin(url, bytes.NewBuffer([]byte(`{"bucket": "travel-sample","num_index_replicas": 0}`)))
	waitForGetURL(url)
	log.Print("Finished setting up a SG db")
}

func configureDB() {
	log.Print("Re-confgiuring a SG db")
	//url := "http://localhost:4985/db/_config"
	//postAdmin(url, bytes.NewBuffer([]byte(`{"import_docs": true, "cache" : { "channel_cache": {"max_length": 100000}}}`)))
	//waitForGetURL(url)
	log.Print("Finished configuring up a SG db")
}

func setupUser() {
	log.Print("Setting up a user")
	url := adminURL + "/db/_user/guest"
	putAdmin(url, bytes.NewBuffer([]byte(`{"password": "guest", "admin_channels": ["*"] }`)))
	waitForGetURL(url)
	log.Print("Finished setting up a user")
}

func get(url string) {
	doHttp(http.MethodGet, url, nil)
}

func primeCache() {
	log.Print("Put all data into cache")
	url := adminURL + "/db"
	get(url)
	log.Printf("Primed cache")

	//log.Print("Put prime channels")
	//url = adminURL + "/db/_changes?filter=sync_gateway%2Fbychannel&channels=ABC"
	//get(url)
	//log.Printf("Primed cache")
}
