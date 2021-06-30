package api

import (
	Err "github.com/42milez/NexusModsUpdateChecker/src/error"
	"github.com/42milez/NexusModsUpdateChecker/src/log"
	"github.com/42milez/NexusModsUpdateChecker/src/util"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"time"
)

const clientTimeout = 30 * time.Second
const secretYaml = "secret.yml"

type secret map[string][]struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func closeConn(resp *http.Response) {
	if resp == nil {
		return
	}
	if err := resp.Body.Close(); err != nil {
		util.Exit(Err.CloseConnectionFailed)
	}
}

func closeFile(f *os.File) {
	if f == nil {
		return
	}
	if err := f.Close(); err != nil {
		util.Exit(Err.CloseFileFailed)
	}
}

func getSecret(kind string, name string) (string, error) {
	if _, err := os.Stat(secretYaml); err != nil {
		if v := os.Getenv(name); v == "" {
			return "", Err.SecretNotFound
		} else {
			return v, nil
		}
	}

	var in *os.File
	var err error

	if in, err = os.Open(secretYaml); err != nil {
		return "", Err.OpenFileFailed
	}
	defer closeFile(in)

	secrets := &secret{}
	if err = yaml.NewDecoder(in).Decode(secrets); err != nil {
		log.D(err.Error())
		return "", Err.DecodeYamlFailed
	}

	for k, sec := range *secrets {
		if k != kind {
			continue
		}
		for _, v := range sec {
			if v.Name == name {
				return v.Value, nil
			}
		}
	}

	return "", Err.SecretNotFound
}
