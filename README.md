# Argon - CI/CD Runner

A lightweight, extensible, and modular CI/CD Runner platform inspired by Jenkins/GitHub Actions but built from scratch using **microservices**, **RabbitMQ**, **Docker**, and **multi-language components**. This system listens to GitHub webhooks, fetches and builds code, runs tests in sandboxed environments, streams logs in real-time, and sends completion notifications.

---

## Tech Stack Overview

| Component            | Tech Used          | Description                                            |
| -------------------- | ------------------ | ------------------------------------------------------ |
| Webhook Listener     | Spring Boot (Java) | Receives GitHub webhooks, validates them, pushes jobs  |
| Orchestrator         | Go                 | Clones repo, parses `.runnerci.yml`, dispatches tasks  |
| Sandbox Executor     | Go                 | Runs job in Docker, applies limits, returns output     |
| Logger Service       | Python + MongoDB   | Logs everything from all services, persists to MongoDB |
| Notification Service | Python + AWS SES   | Sends email notifications on success/failure           |
| Webhook Queue        | RabbitMQ           | Decouples Webhook listener & Orchestrator              |
| Sandbox Queue        | RabbitMQ           | Decouples Orchestrator & Sanbox Executor               |
| Log Queue            | RabbitMQ           | All services emit logs to this queue                   |
| Notification Queue   | RabbitMQ           | Final notification events are passed here              |
| Dashboard (Planned)  | React/Node         | View job status, logs, history in real-time            |

---

## System Intuition & Design Philosophy

This project was born from the idea of **demystifying how CI/CD systems work** and building a production-mimicking runner system using:

* **Polyglot microservices** — right tool for the right job (Java for REST, Go for performance, Python for scripting)
* **Event-driven design** — using RabbitMQ as the backbone for inter-service communication
* **Separation of concerns** — listener, orchestrator, executor, logger, notifier all work independently
* **Security-first thinking** — validating incoming GitHub webhooks using HMAC
* **Sandboxed task execution** — leveraging Docker for safe, isolated execution
* **Observability** — streaming and storing logs in MongoDB, planning for real-time dashboards

---

## Workflow Overview

```mermaid
graph TD
    A[GitHub Push / PR] --> B[Webhook Listener (Spring Boot)]
    B -->|Validates & Pushes| C[webhook.queue]
    C -->|Consumes| D[Orchestrator (Go)]
    D -->|Clones Repo & Parses .runnerci.yml| E[sandbox.queue]
    E -->|Consumes| F[Sandbox Executor (Go)]
    F -->|Executes Jobs in Docker| G[notification.queue]

    %% Logging
    B -->|Logs| L[Logger Service]
    D -->|Logs| L
    F -->|Logs| L

    %% Notification
    G -->|Consumes| H[Notification Service (Python)]
    H -->|Sends Email via AWS SES| I[Developer]

    %% Persistent Logging
    L -->|Persists Logs| J[MongoDB]

```

---

## .runnerci.yml Format

The orchestrator reads a file named `.runnerci.yml` in the repo root like:

```yml
version: 1.0

jobs:
  build:
    image: python:3.10
    steps:
      - name: Install dependencies
        run: pip install -r requirements.txt

      - name: Run main script
        run: python main.py
```

This allows future support for:

* Parallel/Sequential job execution
* Caching
* Conditional steps
* Artifacts

---


## Job Execution

* Jobs are executed in **Docker containers** (isolated sandbox)
* Each step’s `stdout` and `stderr` are captured
* Execution is timeout-limited
* Logs are streamed to `log.queue` in real-time

---

## Centralized Logging

* All services push logs to a common RabbitMQ `log.queue`
* Python-based logger consumes this and writes structured logs into **MongoDB**
* Each log entry contains:

  * `timestamp`
  * `service`
  * `job_id`
  * `message`
  * `severity`

---

## Notifications

* When all jobs complete, a final message is pushed to `notification.queue`
* The Notification Service:

  * Consumes the message
  * Generates a summary
  * Sends an email via AWS SES (Amazon Simple Email)

---

## Real-Time Dashboard (Planned)

The dashboard will:

* Show job status (running, failed, success)
* Display real-time logs using WebSocket
* Allow re-running jobs
* Show historical builds from MongoDB

Tech stack:

* React + Tailwind UI
* WebSocket for live updates
* Express/Flask backend to expose logs

---

## Parallelism & Future Scaling

* Orchestrator can dispatch **multiple jobs concurrently** (Go goroutines)
* Executor can be scaled horizontally (multiple consumers for `sandbox.queue`)
* RabbitMQ ensures fair job distribution
* Can extend to **Kubernetes runners** for elastic scaling

---

## Unit Tests & Integration Tests

* **Webhook Listener**: JUnit tests for signature validation and request handling
* **Orchestrator**: Go unit tests for YAML parsing, Git handling, queue push
* **Sandbox**: Go integration tests with Docker
* **Logger**: Pytest + Mock RabbitMQ + MongoDB for log ingestion
* **Notification**: SES mocking using `moto` (for local testing)

Also planned:

* CI workflow to run tests before pushing
* Local test suite using `docker-compose -f docker-compose.test.yml`

---

## Future Improvements

* Dashboard UI with historical builds
* Auth system (GitHub OAuth)
* Artifacts upload/download
* Slack integration for team notifications
* Caching for faster re-builds
* Cron job triggers (scheduled jobs)
* ML model to detect flaky tests
* SSO and secrets vault (e.g., HashiCorp Vault)

---

## Inspiration

* Jenkins, GitHub Actions
* Drone CI
* Temporal workflows
* Google Cloud Build
