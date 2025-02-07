# Hashitalks 2025 Demo

This is a demo for Hashitalks 2025.

## Setup

- engine `nomad job run dagger-engine.nomad.hcl`
- task `nomad job run dagger-job.nomad.hcl`
- test task
  `nomad job dispatch -meta repository="github.com/kpenfound/greetings-api" -meta commit="main" dagger-job`
