package mc

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/smithyworks/FaDO/cli"
	"github.com/smithyworks/FaDO/database"
	"github.com/smithyworks/FaDO/util"
)

var mcMutex sync.Mutex

func runCommand(name string, args ...string) (out string, err error) {
	cmd := exec.Command(name, args...)

    stdout, err := cmd.Output()
	log.Print("INFO: RUN ", name, " ", strings.Join(args, " "))
    if err != nil { return out, util.ProcessErr(err) }

    out = string(stdout)

	return
}

func SetAlias(sd database.StorageDeploymentRecord) (err error) {
	mcMutex.Lock()
	defer mcMutex.Unlock()

	proto := "http://"
	if sd.UseSSL {
		proto = "https://"
	}

	url := fmt.Sprintf("%v%v", proto, sd.Endpoint)

	_, err = runCommand("mc", "alias", "remove", sd.Alias)
	if err != nil { return util.ProcessErr(err) }

	_, err = runCommand("mc", "alias", "set", sd.Alias, url, sd.AccessKey, sd.SecretKey)

	return util.ProcessErr(err)
}

func SetNotificationTarget(sd database.StorageDeploymentRecord) (err error) {
	endpoint := fmt.Sprintf("endpoint=%v/api/notify", cli.Input.ServerURL)
	_, err = runCommand("mc", "admin", "config", "set", sd.Alias, "notify_webhook:fado", endpoint)
	return util.ProcessErr(err)
}

func Mirror(srcAlias, srcBucketName, dstAlias, dstBucketName string) (err error) {
	mcMutex.Lock()
	defer mcMutex.Unlock()

	src := fmt.Sprintf("%v/%v", srcAlias, srcBucketName)
	dst := fmt.Sprintf("%v/%v", dstAlias, dstBucketName)

	d1 := time.Now()
	_, err = runCommand("mc", "mirror", "--remove", "--overwrite", src, dst)
	d2 := time.Since(d1)
	log.Printf("INFO: Mirror operation took %v seconds.", d2.Seconds())

	return util.ProcessErr(err)
}
