job "postgres" {
  datacenters = ["dc1"]
  type        = "service"

  group "db" {
    count = 1

    restart {
      attempts = 3
      delay    = "15s"
      mode     = "fail"
    }

    task "postgres" {
      driver = "docker"

      config {
        image = "postgres:13"
        port_map {
          db = 5432
        }
        volumes = [
          "local:/var/lib/postgresql/data"
        ]
      }

      env {
        POSTGRES_DB       = "myapp"
        POSTGRES_USER     = "myuser"
        POSTGRES_PASSWORD = "mysecretpassword"
      }

      logs {
        max_files     = 5
        max_file_size = 15
      }

      resources {
        cpu    = 500
        memory = 1024
      }

      service {
        name = "postgres"
        port = "db"

        check {
          type     = "tcp"
          interval = "10s"
          timeout  = "2s"
        }
      }
    }

    network {
      mode = "bridge"
      port "db" {
        static = 5432
        to     = 5432
      }
    }
  }
}