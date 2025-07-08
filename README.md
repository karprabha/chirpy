# Chirpy 🐦

A Twitter-like social media API built with Go that allows users to create accounts, authenticate, and post short messages called "chirps".

## Features

- 🔐 **User Authentication** - JWT-based authentication with refresh tokens
- 📝 **Chirp Management** - Create, read, and delete chirps (max 140 characters)
- 🔍 **Content Filtering** - Automatic profanity filtering
- 👑 **Premium Features** - Chirpy Red subscription via webhooks
- 📊 **Admin Dashboard** - Metrics and system management
- 🗄️ **PostgreSQL Database** - Robust data persistence with migrations
- 🚀 **RESTful API** - Clean HTTP API design

## Tech Stack

- **Language**: Go 1.24.4
- **Database**: PostgreSQL
- **ORM**: sqlc for type-safe SQL queries
- **Authentication**: JWT tokens with bcrypt password hashing
- **Migrations**: Goose for database schema management
- **Environment**: dotenv for configuration

## Quick Start

### Prerequisites

- Go 1.24.4 or higher
- PostgreSQL database
- Environment variables configured

### Installation

1. Clone the repository:

```bash
git clone https://github.com/karprabha/chirpy.git
cd chirpy
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment variables:

```bash
# Create a .env file in the root directory
DB_URL=postgres://username:password@localhost:5432/chirpy_db?sslmode=disable
JWT_SECRET=your-secret-key-here
POLKA_KEY=your-polka-webhook-key
PLATFORM=dev
```

4. Run database migrations:

```bash
goose -dir sql/schema postgres "postgres://username:password@localhost/gator?sslmode=disable" up
```

5. Start the server:

```bash
go run cmd/server/main.go
```

The server will start on port 8080.

## API Documentation

Comprehensive API documentation is available in the `/docs` folder:

- [Authentication API](docs/auth.md) - User login, registration, and token management
- [Users API](docs/users.md) - User profile management
- [Chirps API](docs/chirps.md) - Chirp creation, retrieval, and management
- [Admin API](docs/admin.md) - Administrative endpoints and metrics
- [Webhooks API](docs/webhooks.md) - External integrations and premium features
- [Health Check API](docs/health.md) - Server health monitoring endpoint

## Project Structure

```
chirpy/
├── cmd/server/          # Application entry point
├── internal/
│   ├── auth/           # Authentication utilities
│   ├── config/         # Configuration management
│   ├── database/       # Database models and queries
│   ├── handler/        # HTTP handlers
│   └── middleware/     # HTTP middleware
├── sql/
│   ├── queries/        # SQL queries for sqlc
│   └── schema/         # Database migrations
├── docs/               # API documentation
└── assets/             # Static assets
```

## Database Schema

The application uses PostgreSQL with the following main tables:

- **users** - User accounts with authentication
- **chirps** - User posts/messages
- **refresh_tokens** - JWT refresh token management

## Authentication

Chirpy uses JWT tokens for authentication:

- **Access tokens** - Short-lived (1 hour) for API access
- **Refresh tokens** - Long-lived (60 days) for token renewal
- **Password hashing** - bcrypt for secure password storage

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## API Base URL

When running locally: `http://localhost:8080`

## Health Check

Test if the server is running:

```bash
curl http://localhost:8080/api/healthz
```

## Support

For questions or issues, please open an issue on GitHub.
