job "score-keep-backend" {
  datacenters = ["dc1"]
  type        = "service"

  group "api" {
    count = 2

    network {
      port "http" {
        to = 4000
      }
    }

    service {
      name = "score-keep-api"
      port = "http"
    }

    restart {
      attempts = 3
      delay    = "15s"
      mode     = "delay"
    }

    task "api" {
      driver = "docker"

      config {
        image = "tomfitzgerald406/score-keep:latest"
        ports = ["http"]
      }

      resources {
        cpu    = 500
        memory = 256
      }
    }
  }
}