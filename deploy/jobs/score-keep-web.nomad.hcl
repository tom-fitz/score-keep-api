job "score-keep-web" {
  datacenters = ["dc1"]
  type        = "service"

  group "web" {
    count = 2

    network {
      port "http" {
        to = 3000
      }
    }

    service {
      name = "score-keep-web"
      port = "http"

      tags = [
        "traefik.enable=true",
        "traefik.http.routers.score-keep-web.rule=Host(`scorekeepapp.com`)",
      ]
    }

    task "web" {
      driver = "docker"

      config {
        image = "tomfitzgerald406/score-keep-web:latest"
        ports = ["http"]
      }

      env {
        API_URL = "http://score-keep-api.service.consul:4000"
      }

      resources {
        cpu    = 500
        memory = 256
      }
    }
  }
}