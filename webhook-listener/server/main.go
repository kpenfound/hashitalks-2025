package main

import (
	"log/slog"
	"net/http"
	"os"

	github "github.com/google/go-github/v69/github"
	nomad "github.com/hashicorp/nomad/api"
)

func main() {
	listener := &Listener{
		WebhookSecretKey: os.Getenv("GH_WEBHOOK_SECRET_KEY"),
		NomadJobsClient: nomadClient(
			os.Getenv("NOMAD_ADDR"),
			os.Getenv("NOMAD_TOKEN"),
		).Jobs(),
	}
	http.Handle("/", http.HandlerFunc(listener.ServeHTTP))
	slog.Info("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}

func nomadClient(address string, token string) *nomad.Client {
	// Create a default client configuration
	config := nomad.DefaultConfig()

	config.Address = address
	config.SecretID = token

	// Create a new client
	client, err := nomad.NewClient(config)
	if err != nil {
		panic("Error creating Nomad client")
	}
	return client
}

type Listener struct {
	WebhookSecretKey string
	NomadJobsClient  *nomad.Jobs
}

func (l *Listener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte(l.WebhookSecretKey))
	if err != nil {
		slog.Error("Error validating payload: %v", err)
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		slog.Error("Error parsing webhook: %v", err)
	}
	switch event := event.(type) {
	case *github.PushEvent:
		l.ProcessPushEvent(event)
	}
}

func (l *Listener) ProcessPushEvent(ev *github.PushEvent) {
	slog.Info("Received push event for %s", ev.Repo.FullName)
	err := l.SubmitNomadJob(&NomadPayload{
		repository: ev.Repo.FullName,
		commit:     ev.HeadCommit.ID,
	})
	if err != nil {
		slog.Error("Error submitting Nomad job: %v", err)
	}
}

type NomadPayload struct {
	repository *string
	commit     *string
}

func (l *Listener) SubmitNomadJob(payload *NomadPayload) error {
	jobID := "dagger-job"

	meta := map[string]string{
		"repository": "github.com/" + *payload.repository,
		"commit":     *payload.commit,
	}
	daggerCloudToken := os.Getenv("DAGGER_CLOUD_TOKEN")
	if daggerCloudToken != "" {
		meta["dagger_cloud_token"] = daggerCloudToken
	}
	dispatch, _, err := l.NomadJobsClient.Dispatch(
		jobID,
		meta,
		nil, "", nil,
	)
	slog.Info("Dispatched Nomad Job %v", dispatch)

	return err
}
