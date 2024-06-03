# Ticket Trasaction Service

## Name

transaction-service

## Description

Transaction Service is service that used to create transaction order ticket and view order list

## Installation

1. Ensure, already install golang 1.21 or up
2. Create file .env

```bash
    cp env_example .env
```

3. Fill out the env configuration

```bash
# DATABASE
DATABASE_USER=
DATABASE_PASSWORD=
DATABASE_HOST=
DATABASE_PORT=
DATABASE_SCHEMA=
DATABASE_CONN_MAX_LIFETIME=
DATABASE_MAX_OPEN_CONN=
DATABASE_MAX_IDLE_CONN=

# AWS
AWS_REGION=
AWS_COGNITO_USER_POOL_ID=

# MIDTRANS
MIDTRANS_SERVER_KEY=

# REDIS
CACHER_SERVICE=
CACHER_HOST=
CACHER_PORT=
CACHER_PASSWORD=
CACHER_DEFAULT_EXP=

AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=

SQS_TRANSACTION_URL=
SQS_MAIL_URL=
SQS_TICKET_URL=
SQS_TICKET_FAILED_URL=
SQS_TICKET_SUCCESS_URL=
```

4. Install dependencies:

```bash
make install
```

5. Run in development:

```bash
make run
```

## Test

1. Run unit test

```bash
make unit-test
```

2. Show local coverage (in html)

```bash
make coverage
```

## High Level Design Architecture

![picture](/assets/high-level-architecture.png)

## Low Level Design Architecture

![picture](assets/low-level.jpg)

## ERD

![picture](assets/erd.png)

## Authors

- **Syamsul Bachri** - [Github](https://github.com/SyamSolution)

## Development Tools

- [Fiber](https://gofiber.io/) Rest Framework
- [Zap](https://github.com/uber-go/zap) Log Management
- [Go mod](https://go.dev/ref/mod) Depedency Management
- [Docker](https://www.docker.com/) Container Management
- [Amazon SQS](https://aws.amazon.com/sqs/) Event Management
