package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/orisano/preve"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	param := struct {
		Source  preve.Source
		Version *preve.Version `validate:"omitempty"`
	}{}

	if err := json.NewDecoder(os.Stdin).Decode(&param); err != nil {
		return errors.Wrap(err, "failed to parse json from stdin")
	}

	validator := preve.NewValidator()
	if err := validator.Struct(param); err != nil {
		return errors.Wrap(err, "invalid source section")
	}

	tokens := strings.Split(param.Source.Repo, "/")
	owner, repo := tokens[0], tokens[1]

	gh := preve.MustGithubClient(param.Source.BaseURL)

	ctx := context.Background()
	events, _, err := gh.Activity.ListRepositoryEvents(ctx, owner, repo, nil)
	if err != nil {
		return errors.Wrap(err, "failed to get repository events")
	}

	current := param.Version

	var versions []*preve.Version
	for _, event := range events {
		if event.GetType() != "PullRequestEvent" {
			continue
		}
		payload, err := event.ParsePayload()
		if err != nil {
			return errors.Wrap(err, "broken events")
		}
		prEvent := payload.(*github.PullRequestEvent)
		if prEvent.GetAction() != param.Source.When {
			continue
		}

		version := &preve.Version{
			EventID:       event.GetID(),
			PullRequestID: strconv.Itoa(prEvent.GetNumber()),
		}
		if current == nil {
			versions = append(versions, version)
			break
		}
		if preve.IsOlderVersion(current, version) {
			break
		}
		versions = append(versions, version)
	}

	for i, j := 0, len(versions)-1; i < j; i, j = i+1, j-1 {
		versions[i], versions[j] = versions[j], versions[i]
	}

	if err := json.NewEncoder(os.Stdout).Encode(versions); err != nil {
		return errors.Wrap(err, "failed to encode versions")
	}

	return nil
}
