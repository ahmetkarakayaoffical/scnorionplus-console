package servers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func GetLatestServerReleaseFromAPI(tmpDir string) (*scnorionplus_nats.scnorionplusRelease, error) {
	latestServerReleasePath := filepath.Join(tmpDir, "latest.json")

	if _, err := os.Stat(latestServerReleasePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("latest server releases json file doesn't exist, reason: %v", err)
	}

	data, err := os.ReadFile(latestServerReleasePath)
	if err != nil {
		return nil, err
	}

	r := scnorionplus_nats.scnorionplusRelease{}
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	return &r, nil
}
