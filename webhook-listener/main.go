// A generated module for WebhookListener functions

package main

import (
	"dagger/webhook-listener/internal/dagger"
)

type WebhookListener struct{}

// Run the GitHub Webhook Listener to dispatch nomad jobs
func (m *WebhookListener) Listen(
	// GitHub Webhook Secret
	ghWebhookSecret *dagger.Secret,
	// Nomad Agent Address
	nomadAddr string,
	// Nomad ACL Token
	nomadToken *dagger.Secret,
	// Dagger Cloud Token
	cloudToken *dagger.Secret,
) *dagger.Service {
	return dag.Container().
		From("golang:alpine").
		WithSecretVariable("GH_WEBHOOK_SECRET_KEY", ghWebhookSecret).
		WithEnvVariable("NOMAD_ADDR", nomadAddr).
		WithSecretVariable("NOMAD_TOKEN", nomadToken).
		WithSecretVariable("DAGGER_CLOUD_TOKEN", cloudToken).
		WithWorkdir("/src").
		WithDirectory("/src", dag.CurrentModule().Source().Directory("server")).
		WithExposedPort(8080).
		WithDefaultArgs([]string{"go", "run", "."}).
		AsService()
}
