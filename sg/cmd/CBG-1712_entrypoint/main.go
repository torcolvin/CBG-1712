package main

import (
	"bytes"
	"net/http"
	"log"
	"os"
	"os/exec"
	"github.com/cenkalti/backoff/v4"
)

func main() {
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
	setupBucket()
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("sync_gateway exited non-zero %+v", err)

	}
	sgExited <- struct{}{}
}

func waitForServerUp() {
	url := "http://Administrator:password@localhost:4985"
	getOp := func() error {
		_, err := http.Get("http://Administrator:password@localhost:4985")
		return err
	}
	err := backoff.Retry(getOp, backoff.NewExponentialBackOff())
	if err != nil {
		log.Fatalf("Did not see server %s online with error %s", url, err)
	}
}
func setupBucket() {
	url := "http://Administrator:password@localhost:4985/db1/"
	resp, err:= http.Post(url, "Content-Type: application/json",
	bytes.NewBuffer([]byte(`{"num_index_replicas":0, "bucket": "b1"}'`)))
	if err != nil {
		log.Fatalf("Error POST to %s with %s", url, err)
	}
	log.Printf("response to %s %+v\n", url , resp)
}
