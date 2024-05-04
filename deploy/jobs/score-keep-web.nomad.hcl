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
        "traefik.http.routers.score-keep-web.rule=Host(`34.168.2.208`)",
        "traefik.http.routers.score-keep-web.entrypoints=web,websecure",
        "traefik.http.routers.score-keep-web.tls=true",
        "traefik.http.routers.score-keep-web.tls.certresolver=myresolver",
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