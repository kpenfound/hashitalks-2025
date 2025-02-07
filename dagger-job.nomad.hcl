job "dagger-job" {
  type = "batch"

  parameterized {
    payload       = "forbidden"
    meta_required = ["repository", "commit"]
    meta_optional = ["dagger_cloud_token"]
  }

  group "dagger" {
    task "dagger-job" {
      driver = "docker"

      env {
        DAGGER_CLOUD_TOKEN = "${NOMAD_META_dagger_cloud_token}"
      }

      config {
        image = "registry.dagger.io/engine:v0.15.3"
        entrypoint = ["/usr/local/bin/dagger"]
        args = ["-m", "${NOMAD_META_repository}@${NOMAD_META_commit}", "call", "check"]

        mount {
          type = "bind"
          target = "/var/run/buildkit"
          source = "/var/run/buildkit"
          readonly = false
          bind_options {
            propagation = "rshared"
          }
        }
      }
    }
    restart {
      attempts = 0
      mode = "fail"
    }
  }
  reschedule {
    attempts = 0
  }
}
