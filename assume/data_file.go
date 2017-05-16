package assume

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"time"
)

func dataFilePath() string {
	cu, _ := user.Current()
	return fmt.Sprintf("/tmp/awstools-%s.data", cu.Uid)
}

func readCredentialsExpirationTimestamps() map[string]int64 {
	ts := map[string]int64{}
	f, err := os.Open(dataFilePath())
	if os.IsNotExist(err) {
		return ts
	} else if err != nil {
		log.Fatalln("ERROR: cannot open file:", err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&ts)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	return ts
}

func saveCredentialsExpirationTimestamps(timestamps map[string]int64) {
	f, err := os.OpenFile(dataFilePath(), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln("ERROR: cannot open file:", err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(timestamps)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

}

func saveProfileExpirationTimestamp(profile string, expiration time.Time) {
	ts := readCredentialsExpirationTimestamps()
	ts[profile] = expiration.Unix()
	saveCredentialsExpirationTimestamps(ts)
}

func readProfileExpirationTimestamp(profile string) time.Time {
	ts := readCredentialsExpirationTimestamps()[profile]
	return time.Unix(ts, 0)
}
