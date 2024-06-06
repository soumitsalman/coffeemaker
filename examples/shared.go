package examples

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	ds "github.com/soumitsalman/coffeemaker/sdk/beansack"
)

func localFileStore(contents []ds.Bean) {
	data, _ := json.MarshalIndent(contents, "", "\t")
	filename := fmt.Sprintf("test_REDDIT_%s", time.Now().Format("2006-01-02-15-04-05.json"))
	os.WriteFile(filename, data, 0644)
}
