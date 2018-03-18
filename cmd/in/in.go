package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/orisano/preve"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dest := os.Args[1]
	resource := struct {
		Source  *preve.Source
		Version *preve.Version
	}{}

	if err := json.NewDecoder(os.Stdin).Decode(&resource); err != nil {
		return errors.Wrap(err, "failed to parse resource from stdin")
	}

	filename := filepath.Join(dest, "id")
	if err := ioutil.WriteFile(filename, []byte(resource.Version.PullRequestID), 0644); err != nil {
		return errors.Wrapf(err, "failed to write %q", filename)
	}

	output := struct {
		Version  *preve.Version        `json:"version"`
		MetaData []preve.MetadataField `json:"metadata"`
	}{
		Version:  resource.Version,
		MetaData: []preve.MetadataField{},
	}
	if err := json.NewEncoder(os.Stdout).Encode(&output); err != nil {
		return errors.Wrap(err, "failed to encode results")
	}

	return nil
}
