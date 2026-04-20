# Managify

Managify is a Project Management Platform engineered for scalability and performance. Built with a cloud-native architecture, it leverages **Golang** for high-throughput backend services and **React** for a responsive frontend experience.

This repository contains the complete source code for the backend API, frontend application, and infrastructure definitions.

## Key Features

### Core Platform
- **Project & Task Management**: Robust system for creating projects, assigning tasks, and tracking progress with custom statuses.
- **Role-Based Access Control (RBAC)**: Granular permission system managing Users, Roles, and Team access.
- **Secure Authentication**: JWT-based stateless authentication with secure password hashing.
- **Real-time Performance**: Optimized Go backend handling concurrent requests with minimal latency.


## Technology Stack

| Category | Technologies |
|----------|--------------|
| **Backend** | Go (Golang) 1.24, Fiber Framework, JWT |
| **Frontend** | React 19, Vite, Tailwind CSS, Ant Design |
| **Database** | MongoDB (Atlas / Self-hosted) |
| **DevOps** | Docker, GitHub Actions |
| **Cloud (AWS)** | ECS Fargate, ECR, ALB, VPC, Terraform |
| **Documentation** | Swagger UI (OpenAPI 3.0) |

## Architecture Overview

The system follows a modular Monolith approach designed for easy transition to Microservices:

```
[Client (React/Mobile)] 
       │
       ▼
[Load Balancer (AWS ALB)]
       │
       ▼
[Managify API (Go/Fiber)] 
  ├── Handler Layer (HTTP)
  ├── Service Layer (Business Logic)
  └── Repository Layer (Data Access)
       │
       ▼
[MongoDB Cluster]
```

## Getting Started

### Prerequisites
- **Go**: v1.24+
- **Docker**: v20.10+
- **MongoDB**: v6.0+
- **Node.js**: v18+ (for Frontend)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/Evolve-Vitalis/managify-go.git
   cd managify-go
   ```

2. **Environment Setup**
   ```bash
   cp .env.example .env
   # Update .env with your local MongoDB URI and configs
   ```

3. **Run with Docker Compose (Recommended)**
   ```bash
   docker-compose up -d --build
   ```
   The API will be available at `http://localhost:3000` and Swagger docs at `http://localhost:3000/swagger`.

4. **Manual Run (Backend)**
   ```bash
   go mod download
   go run main.go
   ```

## Deployment

### AWS ECS (Terraform)
Infrastructure is managed via Terraform in the `/terraform` directory.

1. Configure AWS Credentials.
2. Initialize and Apply:
   ```bash
   cd terraform
   terraform init
   terraform apply -var="project_name=managify"
   ```

## API Documentation

Interactive API documentation is generated via Swagger.

- **Local**: `http://localhost:3000/swagger/index.html`
- **Specification**: `/docs/swagger.json`

## Contributing

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/access-control`).
3. Commit your changes.
4. Push to the branch and open a Pull Request.

## License

This project is licensed under the MIT License.

---
**Maintained by [Doğuhannilt](https://github.com/Evolve-Vitalis)**
