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
	resource := struct {
		Source  *preve.Source
		Version *preve.Version `validate:"omitempty"`
	}{}

	if err := json.NewDecoder(os.Stdin).Decode(&resource); err != nil {
		return errors.Wrap(err, "failed to parse json from stdin")
	}

	validator := preve.NewValidator()
	if err := validator.Struct(resource); err != nil {
		return errors.Wrap(err, "invalid source section")
	}

	versions, err := getNewVersions(resource.Source, resource.Version)
	if err != nil {
		return errors.Wrap(err, "failed to get new versions")
	}

	for i, j := 0, len(versions)-1; i < j; i, j = i+1, j-1 {
		versions[i], versions[j] = versions[j], versions[i]
	}

	if err := json.NewEncoder(os.Stdout).Encode(versions); err != nil {
		return errors.Wrap(err, "failed to encode versions")
	}

	return nil
}

func filterPullRequestEvent(events []*github.Event) []*github.Event {
	r := events[:0]
	for _, event := range events {
		if event.GetType() == "PullRequestEvent" {
			r = append(r, event)
		}
	}
	return r
}

func getNewVersions(source *preve.Source, current *preve.Version) ([]*preve.Version, error) {
	tokens := strings.Split(source.Repo, "/")
	owner, repo := tokens[0], tokens[1]

	var versions []*preve.Version

	gh := preve.MustGitHubClient(source.BaseURL)
	ctx := context.Background()
	opt := &github.ListOptions{}
	for {
		events, resp, err := gh.Activity.ListRepositoryEvents(ctx, owner, repo, opt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get repository events")
		}
		events = filterPullRequestEvent(events)
		for _, event := range events {
			payload, err := event.ParsePayload()
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse event payload")
			}
			prEvent := payload.(*github.PullRequestEvent)
			if prEvent.GetAction() != source.When {
				continue
			}
			version := &preve.Version{
				EventID:       event.GetID(),
				PullRequestID: strconv.Itoa(prEvent.GetNumber()),
			}
			if current == nil {
				versions = append(versions, version)
				return versions, nil
			}
			if preve.IsOlderVersion(current, version) {
				return versions, nil
			}
			versions = append(versions, version)
		}
		opt.Page = resp.NextPage
	}
	return versions, nil
}
