# 🌱 Raksana Backend

A comprehensive sustainability and gamification platform backend system built with microservices architecture. This project is dedicated to BeeFest SDLC, featuring habit tracking, environmental challenges, and eco-friendly activities.

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Database Schema](#database-schema)
- [Setup Instructions](#setup-instructions)
- [API Services](#api-services)
- [Environment Variables](#environment-variables)
- [Development](#development)
- [Deployment](#deployment)

---

## 🎯 Overview

Raksana Backend is a microservices-based platform that gamifies sustainable living. The system enables users to:

- **Track Habits**: Create and manage daily sustainability habits
- **Participate in Challenges**: Complete daily environmental challenges
- **Attend Events**: Join eco-friendly community events
- **Scan Items**: Use AI to identify recyclable items and get sustainability tips
- **Earn Rewards**: Gain points, level up, and claim treasures
- **View Analytics**: Track progress through recaps and statistics
- **Explore Regions**: Discover environmental regions and tree-planting initiatives

---

## 🏗️ Architecture

The project consists of three main services:

### 1. **Backend Service** (Go + Fiber)
Main API service handling all business logic, user management, and gamification features.

### 2. **Admin Service** (Laravel + Filament)
Administrative dashboard for managing content, users, and system configurations.

### 3. **Lambda Script** (Node.js)
Serverless function for automated challenge generation using AI.

```
┌─────────────────┐
│   Mobile App    │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────────┐
│      Backend Service (Go)           │
│  ┌──────────────────────────────┐  │
│  │  REST API (Fiber)            │  │
│  │  - Auth, Users, Profiles     │  │
│  │  - Habits, Tasks, Packets    │  │
│  │  - Challenges, Events        │  │
│  │  - Quests, Treasures         │  │
│  │  - Scans, Greenprints        │  │
│  └──────────────────────────────┘  │
└─────────────┬───────────────────────┘
              │
     ┌────────┴────────┐
     ▼                 ▼
┌─────────┐      ┌──────────┐
│PostgreSQL│      │  Redis   │
└─────────┘      └──────────┘
     ▲
     │
┌────┴──────────────────┐
│  Admin Service        │
│  (Laravel + Filament) │
└───────────────────────┘

┌────────────────────────┐
│  Lambda Script         │
│  (AI Challenge Gen)    │
└────────────────────────┘
```

---

## 🛠️ Tech Stack

### Backend Service (Go)

#### Core Framework & Libraries
- **Go** `1.25.0` - Primary programming language
- **Fiber** `v2.52.9` - Fast HTTP web framework
- **PostgreSQL** with `pgx/v5` `5.7.5` - Database driver
- **Redis** `v9.12.1` - Caching and session management

#### Key Dependencies
```go
// Authentication & Security
github.com/golang-jwt/jwt/v5          v5.3.0
golang.org/x/crypto                   v0.41.0

// Validation
github.com/go-playground/validator/v10 v10.27.0

// AWS Services
github.com/aws/aws-sdk-go-v2          v1.39.0
  - config                            v1.31.2
  - credentials                       v1.18.6
  - service/s3                        v1.87.1
  - service/rekognition               v1.50.2
  - feature/s3/manager                v1.19.0

// AI & ML
github.com/google/generative-ai-go    v0.20.1
google.golang.org/api                 v0.231.0

// Firebase
firebase.google.com/go                v3.13.0+incompatible

// Configuration
github.com/spf13/viper                v1.20.1

// UUID Generation
github.com/google/uuid                v1.6.0
```

### Admin Service (Laravel + PHP)

#### Core Framework
- **PHP** `^8.2`
- **Laravel** `^12.0` - Web application framework
- **Filament** `^3.3` - Admin panel framework

#### Key Dependencies
```json
"require": {
  "php": "^8.2",
  "laravel/framework": "^12.0",
  "filament/filament": "^3.3",
  "laravel/tinker": "^2.10.1",
  
  // QR Code Generation
  "endroid/qr-code": "^5.0",
  
  // PDF Generation
  "barryvdh/laravel-dompdf": "^3.1",
  
  // Maps Integration
  "cheesegrits/filament-google-maps": "^3.0",
  
  // JWT Authentication
  "firebase/php-jwt": "^6.11",
  
  // Cloud Storage
  "league/flysystem-aws-s3-v3": "^3.0"
}

"require-dev": {
  "fakerphp/faker": "^1.23",
  "laravel/pail": "^1.2.2",
  "laravel/pint": "^1.13",
  "laravel/sail": "^1.41",
  "mockery/mockery": "^1.6",
  "nunomaduro/collision": "^8.6",
  "phpunit/phpunit": "^11.5.3"
}
```

### Lambda Script (Node.js)

#### Runtime & Dependencies
```json
"dependencies": {
  "@google/genai": "^1.16.0",  // Google Generative AI
  "pg": "^8.16.3"               // PostgreSQL client
}
```

---

## 📁 Project Structure

```
raksana-backend/
│
├── backend-service/          # Main Go API service
│   ├── app/
│   │   ├── db.go            # Database connection
│   │   ├── redis.go         # Redis client
│   │   └── router.go        # Route registration
│   ├── configs/             # Configuration files
│   │   ├── ai.go            # AI client setup
│   │   ├── aws.go           # AWS services setup
│   │   └── fcm.go           # Firebase Cloud Messaging
│   ├── handlers/            # HTTP handlers (18 handlers)
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── packet_handler.go
│   │   ├── task_handler.go
│   │   ├── challenge_handler.go
│   │   ├── event_handler.go
│   │   ├── quest_handler.go
│   │   ├── treasure_handler.go
│   │   ├── scan_handler.go
│   │   ├── memory_handler.go
│   │   ├── journal_handler.go
│   │   ├── recap_handler.go
│   │   ├── leaderboard_handler.go
│   │   ├── streak_handler.go
│   │   ├── point_handler.go
│   │   ├── history_handler.go
│   │   ├── activity_handlers.go
│   │   └── region_handler.go
│   ├── helpers/             # Utility functions
│   │   ├── config.go
│   │   ├── exp.go
│   │   ├── hash.go
│   │   ├── jwt.go
│   │   ├── point.go
│   │   ├── task.go
│   │   └── time.go
│   ├── models/              # Data transfer objects
│   ├── repositories/        # Database queries (SQLC)
│   ├── services/            # Business logic
│   ├── exceptions/          # Error handling
│   ├── main.go             # Application entry point
│   ├── go.mod              # Go dependencies
│   ├── Dockerfile          # Container configuration
│   └── docker-compose.yml  # Docker compose setup
│
├── admin-service/           # Laravel admin panel
│   ├── app/
│   │   ├── Console/
│   │   ├── Filament/       # Filament admin resources
│   │   ├── Http/
│   │   ├── Models/         # Eloquent models
│   │   └── Providers/
│   ├── config/             # Laravel configuration
│   ├── database/
│   │   ├── migrations/     # Database migrations
│   │   ├── seeders/        # Database seeders
│   │   └── database.sqlite # SQLite database
│   ├── routes/             # API routes
│   ├── resources/          # Views and assets
│   ├── composer.json       # PHP dependencies
│   └── package.json        # NPM dependencies
│
├── lambda-script/           # AWS Lambda function
│   ├── index.js            # Lambda handler
│   ├── package.json        # Node dependencies
│   └── node_modules/       # Installed packages
│
├── DB.md                   # Database schema (DBML)
├── Structured_ERD.md       # ERD documentation
└── README.md              # This file
```

---

## 💾 Database Schema

The database is organized into **5 functional domains**:

### 1. 👤 User Domain
- **users**: User accounts and authentication
- **profiles**: User profiles with levels, XP, and points
- **statistics**: User activity statistics
- **histories**: Point and XP transaction history
- **logs**: User activity logs
- **memories**: User-uploaded media files

### 2. 🎯 Habit Domain
- **packets**: Habit tracking containers
- **habits**: Habit templates
- **tasks**: Daily habit tasks
- **recaps**: Weekly/monthly summaries
- **recap_details**: Detailed recap statistics

### 3. 🏃‍♂️ Activity Domain
- **details**: Base information for activities
- **codes**: QR codes for activities
- **challenges**: Daily environmental challenges
- **participations**: Challenge completions
- **quests**: Location-based collaborative quests
- **contributions**: Quest contributions
- **events**: Community events
- **attendances**: Event attendance records
- **treasures**: Claimable rewards
- **claimed**: Treasure claims

### 4. 🌱 Sustainability Domain
- **scans**: Item scanning records
- **items**: Scanned items
- **greenprints**: Sustainability guides
- **steps**: Greenprint instructions
- **materials**: Required materials
- **tools**: Required tools
- **regions**: Environmental regions

### 5. 🔧 System Domain
- **sessions**: User session management

**Total Tables**: 30

For detailed schema information, see:
- [DB.md](./DB.md) - DBML format
- [Structured_ERD.md](./Structured_ERD.md) - Organized by domain

---

## 🚀 Setup Instructions

### Prerequisites

- **Go** 1.25.0 or higher
- **PHP** 8.2 or higher
- **PostgreSQL** 14 or higher
- **Redis** 6 or higher
- **Node.js** 18 or higher (for Lambda)
- **Docker** & **Docker Compose** (optional)

### 1. Clone Repository

```bash
git clone <repository-url>
cd raksana-backend
```

### 2. Backend Service Setup

```bash
cd backend-service

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Configure environment variables (see Environment Variables section)
nano .env

# Run migrations (ensure PostgreSQL is running)
# Migrations are handled by admin service

# Run the service
go run main.go

# Or with Docker
docker-compose up --build
```

The backend service will run on `http://localhost:3000`

### 3. Admin Service Setup

```bash
cd admin-service

# Install PHP dependencies
composer install

# Install NPM dependencies
npm install

# Copy environment file
cp .env.example .env

# Generate application key
php artisan key:generate

# Run migrations
php artisan migrate

# Seed database (optional)
php artisan db:seed

# Run the service
php artisan serve

# In another terminal, run Vite
npm run dev
```

The admin panel will run on `http://localhost:8000`

### 4. Lambda Script Setup

```bash
cd lambda-script

# Install dependencies
npm install

# Configure environment variables (AWS Lambda console)
# See Environment Variables section
```

---

## 🔌 API Services

### Backend Service Endpoints

The main API runs on port **3000** with base path `/api`

#### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/auth/profile` - Get user profile

#### User Management
- `GET /api/users/profile` - Get current user profile
- `PUT /api/users/profile` - Update profile
- `POST /api/users/profile/picture` - Upload profile picture

#### Habit Tracking
- `GET /api/packets` - List habit packets
- `POST /api/packets` - Create habit packet
- `GET /api/packets/:id` - Get packet details
- `PUT /api/packets/:id` - Update packet
- `DELETE /api/packets/:id` - Delete packet
- `GET /api/tasks` - List tasks
- `POST /api/tasks/:id/complete` - Complete task

#### Challenges
- `GET /api/challenges` - List challenges
- `GET /api/challenges/:id` - Get challenge details
- `POST /api/challenges/:id/participate` - Participate in challenge

#### Events
- `GET /api/events` - List events
- `GET /api/events/:id` - Get event details
- `POST /api/events/:id/attend` - Register attendance

#### Quests
- `GET /api/quests` - List quests
- `GET /api/quests/:id` - Get quest details
- `POST /api/quests/:id/contribute` - Contribute to quest

#### Treasures
- `GET /api/treasures` - List treasures
- `POST /api/treasures/:id/claim` - Claim treasure

#### Scanning & Greenprints
- `POST /api/scans` - Scan item (AI-powered)
- `GET /api/scans/:id/greenprints` - Get greenprints for scanned item

#### Analytics
- `GET /api/journal` - Get activity journal
- `GET /api/leaderboard` - Get leaderboard
- `GET /api/streak` - Get current streak
- `GET /api/recaps` - Get weekly/monthly recaps
- `GET /api/history` - Get point history
- `GET /api/statistics` - Get user statistics

#### Regions
- `GET /api/regions` - List environmental regions
- `GET /api/regions/:id` - Get region details

---

## 🔐 Environment Variables

### Backend Service (.env)

```env
# Server Configuration
PORT=3000

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=raksana_db

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRATION=24h

# AWS Configuration
AWS_REGION=ap-southeast-1
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_S3_BUCKET=raksana-bucket

# Google AI Configuration
GEMINI_API_KEY=your_gemini_api_key

# Firebase Configuration
FIREBASE_CREDENTIALS_PATH=./serviceAccountKey.json
```

### Admin Service (.env)

```env
APP_NAME=Raksana
APP_ENV=local
APP_KEY=base64:generated_key
APP_DEBUG=true
APP_URL=http://localhost:8000

DB_CONNECTION=pgsql
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=raksana_db
DB_USERNAME=postgres
DB_PASSWORD=your_password

CACHE_DRIVER=redis
SESSION_DRIVER=redis
QUEUE_CONNECTION=redis

REDIS_HOST=localhost
REDIS_PORT=6379

AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_DEFAULT_REGION=ap-southeast-1
AWS_BUCKET=raksana-bucket
```

### Lambda Script (Environment Variables)

```env
DB_HOST=your_rds_host
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=raksana_db

GEMINI_API_KEY=your_gemini_api_key
SYSTEM_INSTRUCTION=your_challenge_generation_prompt
```

---

## 💻 Development

### Backend Service Development

```bash
cd backend-service

# Run with hot reload (using air)
air

# Run tests
go test ./...

# Format code
go fmt ./...

# Lint code
golangci-lint run

# Build binary
go build -o bin/raksana-backend
```

### Admin Service Development

```bash
cd admin-service

# Run development server with queue and logs
composer dev

# Run tests
composer test

# Format code (Laravel Pint)
./vendor/bin/pint

# Clear cache
php artisan config:clear
php artisan cache:clear
php artisan view:clear
```

### Code Generation (SQLC)

The backend service uses SQLC for type-safe database queries:

```bash
cd backend-service

# Generate repository code from SQL queries
sqlc generate
```

---

## 🐳 Deployment

### Docker Deployment

#### Backend Service

```bash
cd backend-service

# Build and run with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

#### Dockerfile
The backend service uses multi-stage builds for optimized image size:

```dockerfile
# Stage 1: Build
FROM golang:alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./out .

# Stage 2: Run
FROM alpine:latest
WORKDIR /main
COPY --from=builder /build/out ./out
RUN chmod +x ./out
EXPOSE 3000
CMD ["./out"]
```

### AWS Lambda Deployment

```bash
cd lambda-script

# Create deployment package
zip -r function.zip index.js node_modules/

# Upload to AWS Lambda
aws lambda update-function-code \
  --function-name raksana-challenge-generator \
  --zip-file fileb://function.zip
```

### Production Considerations

- Use environment-specific configuration files
- Enable HTTPS/TLS
- Configure CORS properly
- Set up monitoring and logging
- Implement rate limiting
- Use connection pooling for database
- Enable Redis persistence
- Set up automated backups for PostgreSQL
- Use CDN for static assets
- Implement health check endpoints

---

## 📊 API Documentation

For detailed API documentation with request/response examples, see:
- Postman Collection (coming soon)
- OpenAPI/Swagger Docs (coming soon)

---

## 🤝 Contributing

1. Create a feature branch
2. Commit your changes
3. Push to the branch
4. Create a Pull Request

---

## 📝 License

This project is licensed under the MIT License.

---

## 👥 Team

Developed for BeeFest SDLC

---

## 📞 Support

For issues and questions, please create an issue in the repository.

---

**Built with ❤️ for a sustainable future 🌱**
