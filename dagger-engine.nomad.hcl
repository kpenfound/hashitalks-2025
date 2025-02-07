job "dagger-engine" {
  type = "system"

  group "dagger" {
    task "dagger-engine" {
      driver = "docker"

      config {
        image = "registry.dagger.io/engine:v0.15.3"
        privileged = true
        cap_add = ["all"]
        mount {
          type = "bind"
          target = "/var/run/buildkit"
          source = "/var/run/buildkit"
          readonly = false
          bind_options {
            propagation = "rshared"
          }
        }
        mount {
          type = "bind"
          target = "/var/lib/dagger"
          source = "/var/lib/dagger"
          readonly = false
          bind_options {
            propagation = "rshared"
          }
        }
      }
      resources {
        memory = 3096
      }
    }
  }
}
