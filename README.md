# "ToDo API"

## Technical stack

- Backend building blocks
  - [grpc-ecosystem/grpc-gateway/v2](https://github.com/grpc-ecosystem/grpc-gateway)
  - [golang-migrate/migrate/v4](https://github.com/golang-migrate/migrate)
  - [kyleconroy/sqlc](https://github.com/kyleconroy/sqlc)
    - [pq](github.com/lib/pq)
  
  - Infrastructure
    - Postgres, RabbitMQ
    - Hashicorp Nomad, Consul (Connect), Vault, Terraform
    - docker and docker-compose
    - devcontainer for reproducible development environment
