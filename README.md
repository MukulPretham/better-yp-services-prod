# better-yp-services-prod

A production-grade distributed website monitoring and notification system.

## Overview

This project is designed to monitor the status and latency of multiple websites across various regions, and notify designated users when issues are detected. It uses a microservices architecture written in Go, leverages Redis streams for messaging, and PostgreSQL for persistent storage. The system is containerized using Docker for easy deployment.

**Main Components:**
- **Publisher**: Periodically scans all websites from the database and publishes them as tasks to a Redis stream.
- **Consumer**: Deployed per region, consumes website monitoring tasks, checks site status/latency, updates the database, and triggers notifications if failures are detected.
- **Notification Server (nf-server)**: Listens to notification events and sends email alerts to users associated with affected websites/regions.

## Architecture

```
+----------------+      +----------------+      +---------------------+      +-----------------+
|   Publisher    | ---> |    Redis       | ---> |    Consumer(s)      | ---> | nf-server       |
| (Scheduler)    |      |   (Streams)    |      |  (per Region)       |      | (Notification)  |
+----------------+      +----------------+      +---------------------+      +-----------------+
        |                      |                       |                              |
        |   Websites/tasks     |    Monitored events   |   Notification events        |
        |                     |---------------------> |--------------------------->  |
+----------------+      +----------------+      +---------------------+      +-----------------+
|  Postgres DB   | <-------------------------------------------------------------|
+----------------+
```

- **Publisher**: Fetches websites, serializes them, and pushes them to a Redis stream (`websites`). Runs on a schedule (e.g., every 3 minutes).
- **Consumer(s)**: Each region has a consumer group. They read from the `websites` stream, check website health, and update status/latency in PostgreSQL. On failure, they push notification events to the `notifications` stream.
- **nf-server**: Listens to the `notifications` stream, fetches the relevant user emails, and sends notification emails using SMTP.

## Features

- Distributed monitoring across multiple regions.
- Efficient, at-least-once delivery of monitoring tasks using Redis Streams and consumer groups.
- Automatic failure detection and notification.
- Modular, containerized design for scalability.

## Folder Structure

```
.
├── publisher/         # Publishes website monitoring tasks
│   ├── main.go
│   ├── utils/
│   └── redis_utils/
├── consumer/          # Consumes tasks, checks websites, updates DB, triggers notifications
│   ├── main.go
│   ├── helpers/
│   └── Dockerfile
├── nf-server/         # Notification server (sends emails)
│   ├── main.go
│   ├── helpers.go
│   └── Dockerfile
├── go.work
└── Dockerfiles
```

## Technology Stack

- **Go**: All services are implemented in Go.
- **Redis**: Used for inter-service messaging (streams).
- **PostgreSQL**: Main data store for websites, users, status, latency.
- **Docker**: Containerization for each service.
- **SMTP**: For sending notification emails.

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.24+
- Redis server
- PostgreSQL server

### Environment Variables

Ensure to set up the following environment variables for each service (see `.env` files):

- `REGION` (for consumers): Region identifier for the consumer.
- `FromEmail`, `AppPassword`, `SmptHost`, `SmptIP`: For email notifications.
- Database credentials for PostgreSQL.
- Redis connection information.

### Building and Running

You can build and run each service using Docker:

```bash
# Build and run publisher
cd publisher
docker build -t publisher .
docker run --env-file ../.env publisher

# Similarly, build and run consumer and nf-server
```

Orchestrate all services and dependencies using Docker Compose (recommended for production).

## How It Works

1. **Publisher** connects to the database, fetches all websites, and pushes them to a Redis stream.
2. **Consumer(s)** for each region consume these tasks, check the status and latency of each website, and update the database. If a website is down and was previously up, it pushes a notification event.
3. **nf-server** listens for notification events and sends emails to users associated with the affected website/region.

## Extending and Customizing

- Add more regions by deploying additional consumers with different `REGION` environment values.
- Customize notification logic in `nf-server` as needed.
- Extend website/user schema in the database for more features.

## License

MIT

## Author

[Mukul Pretham](https://github.com/MukulPretham)
