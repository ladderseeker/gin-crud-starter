# Gin CRUD Starter

A simple RESTful API starter with Gin, GORM, and PostgreSQL featuring clean architecture, structured logging, consistent error handling, and Docker support.

## Features

- Clean layered architecture (controllers, services, repositories)
- Configuration management with environment variables
- Structured JSON logging with Zap
- Consistent error handling and responses
- PostgreSQL integration with GORM and connection pooling
- Request logging, CORS, and recovery middleware
- Input validation with Gin binding
- Unit tests with mocking
- Docker and Docker Compose support

## Running the Application

### Option 1: Docker Compose (Recommended)

```bash
# Start application and database
docker-compose up -d

# Seed test data
docker exec -it gin-crud-api go run cmd/seed/main.go
```

### Option 2: Local Setup

Prerequisites:
- Go 1.19+
- PostgreSQL 12+

```bash
# Setup PostgreSQL database
sudo -u postgres psql
CREATE DATABASE gin_crud;
CREATE USER postgres WITH ENCRYPTED PASSWORD 'postgres';
GRANT ALL PRIVILEGES ON DATABASE gin_crud TO postgres;
\q
```
```
# Configure and run application
# Create .env file with DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
go mod download
go run cmd/api/main.go
```

## API Endpoints

- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `GET /health` - Health check

## Test Data

run this SQL:

```sql
TRUNCATE TABLE users RESTART IDENTITY CASCADE;
INSERT INTO users 
(name, email, password, role, active, created_at, updated_at) 
VALUES 
('Admin User', 'admin@example.com', '$2a$10$pNQUxXQyxv8MwH.S8LsJx.U5XKTUJZ/a8Qs8tG2XWdnMQ396ZGgR6', 'admin', true, NOW(), NOW()),
('Regular User', 'user@example.com', '$2a$10$pNQUxXQyxv8MwH.S8LsJx.U5XKTUJZ/a8Qs8tG2XWdnMQ396ZGgR6', 'user', true, NOW(), NOW()),
('Inactive User', 'inactive@example.com', '$2a$10$pNQUxXQyxv8MwH.S8LsJx.U5XKTUJZ/a8Qs8tG2XWdnMQ396ZGgR6', 'user', false, NOW(), NOW()),
('John Smith', 'john.smith@example.com', '$2a$10$pNQUxXQyxv8MwH.S8LsJx.U5XKTUJZ/a8Qs8tG2XWdnMQ396ZGgR6', 'user', true, NOW(), NOW()),
('Jane Doe', 'jane.doe@example.com', '$2a$10$pNQUxXQyxv8MwH.S8LsJx.U5XKTUJZ/a8Qs8tG2XWdnMQ396ZGgR6', 'user', true, NOW(), NOW());
```

## Testing

Run unit tests:
```bash
go test ./...
```

