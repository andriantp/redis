# Redis Learning Series

This repository contains the source code and examples accompanying my Redis articles on Medium.

## Articles

### Chapter #1: Securing Redis with ACL and TLS

Learn how to secure Redis using:

* TLS encryption
* Access Control Lists (ACL)
* User authentication
* Certificate generation and configuration
* Verifying secure connections using `redis-cli`

Medium [Article](https://andriantriputra.medium.com/redis-securing-redis-with-acl-and-tls-1a2854fe725f)

### Chapter #2: Exploring Redis Features with Go and Rust

Hands-on examples of Redis features implemented in both Go and Rust, including:

* Strings
* Hashes
* Pub/Sub
* Sorted Sets
* Expiration and TTL

Medium [Article](https://andriantriputra.medium.com)

## Project Structure

```text
.
├── docker
│   ├── certs
│   ├── redis.conf
│   ├── redis-data
│   ├── redis.docker-compose.yml
│   └── users.acl
├── example
│   ├── go-redis
│   └── rust-redis
├── Makefile
└── README.md
```

## Requirements

* Docker
* Docker Compose
* Make
* Go (for Go examples)
* Rust (for Rust examples)

## Environment

Generated Environment:

```bash
make setup
```

## Running Redis

Start Redis with TLS and ACL enabled:

```bash
make up
```

Check that the container is running:

```bash
docker ps
```

Stop Redis:

```bash
make down
```

## Verifying the Connection

Use `redis-cli` with TLS enabled:

```bash
redis-cli \
  --tls \
  --cacert docker/certs/ca.crt \
  --cert docker/certs/client.crt \
  --key docker/certs/client.key \
  -u "rediss://<username>:<password>@localhost:6379"
```

## Examples

### Go

```bash
cd example/go-redis
go run .
```

### Rust

```bash
cd example/rust-redis
cargo run
```

## Learning Goals

This repository focuses on understanding Redis from both operational and application perspectives:

* How to secure Redis in production environments.
* How ACL and TLS work together.
* How Redis authentication is configured.
* How Redis data structures behave in practice.
* How to interact with Redis using Go and Rust.
* Building confidence through hands-on experimentation.

## 👤 Author

**Andrian Tri Putra**

* Medium: https://andriantriputra.medium.com
* GitHub: https://github.com/andriantp
* GitHub (alternative): https://github.com/AndrianTriPutra

## 📄 License

Licensed under the Apache License 2.0.
