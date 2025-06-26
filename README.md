# API IRAS - REST API with Golang, Gin & GORM

REST API yang dibangun menggunakan Golang dengan framework Gin dan GORM sebagai ORM. Project ini mengikuti struktur folder yang biasa digunakan oleh senior developer.

## 🚀 Fitur

- **Framework**: Gin (HTTP web framework)
- **ORM**: GORM (Object Relational Mapping)
- **Database**: PostgreSQL (mudah diganti ke database lain)
- **Validation**: Built-in request validation
- **Middleware**: CORS, Logger, Error Handler, Authentication
- **Architecture**: Clean Architecture dengan separation of concerns
- **Environment**: Konfigurasi menggunakan .env file

## 📁 Struktur Project

```
api-iras/
├── cmd/
│   └── server/           # Entry point aplikasi
│       └── main.go
├── internal/             # Code yang tidak bisa diakses dari luar
│   ├── config/          # Konfigurasi aplikasi
│   ├── controllers/     # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models dan structs
│   ├── routes/          # Route definitions
│   └── services/        # Business logic
├── pkg/                 # Code yang bisa diakses dari luar
│   └── utils/           # Utility functions
├── .env                 # Environment variables
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## 🛠️ Setup dan Instalasi

### Prerequisites

- Go 1.21+
- PostgreSQL
- Git

### 1. Clone Repository

```bash
git clone <repository-url>
cd api-iras
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Setup Database

Buat database PostgreSQL dan update file `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=api_iras
DB_SSLMODE=disable
```

### 4. Setup Environment Variables

Copy `.env` file dan sesuaikan dengan konfigurasi Anda:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=api_iras
DB_SSLMODE=disable

# Server Configuration
PORT=8080
ENV=development

# JWT Configuration (untuk implementasi di masa depan)
JWT_SECRET=your-secret-key-here

# API Configuration
API_VERSION=v1
```

### 5. Run Application

```bash
go run cmd/server/main.go
```

Aplikasi akan berjalan pada `http://localhost:8080`

## 📚 API Endpoints

### Health Check
- `GET /health` - Status aplikasi

### Categories
- `GET /api/v1/categories` - Get all categories
- `GET /api/v1/categories/:id` - Get category by ID
- `POST /api/v1/categories` - Create category (auth required)
- `PUT /api/v1/categories/:id` - Update category (auth required)
- `DELETE /api/v1/categories/:id` - Delete category (auth required)

### Products
- `GET /api/v1/products` - Get all products (with pagination)
- `GET /api/v1/products/:id` - Get product by ID
- `POST /api/v1/products` - Create product (auth required)
- `PUT /api/v1/products/:id` - Update product (auth required)
- `DELETE /api/v1/products/:id` - Delete product (auth required)

### Users
- `POST /api/v1/users/register` - Register new user
- `GET /api/v1/users` - Get all users (auth required)
- `GET /api/v1/users/:id` - Get user by ID (auth required)
- `PUT /api/v1/users/:id` - Update user (auth required)
- `DELETE /api/v1/users/:id` - Delete user (auth required)

### Authentication
- `POST /api/v1/auth/login` - Login (JWT implementation needed)

## 🔧 Development

### Build Application

```bash
go build -o bin/api-server cmd/server/main.go
```

### Run Tests

```bash
go test ./...
```

### Format Code

```bash
go fmt ./...
```

### Vet Code

```bash
go vet ./...
```

## 📝 API Usage Examples

### Create Category

```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{
    "name": "Electronics",
    "description": "Electronic devices and accessories"
  }'
```

### Get Products with Pagination

```bash
curl "http://localhost:8080/api/v1/products?page=1&limit=10"
```

### Create Product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{
    "name": "Laptop Gaming",
    "description": "High performance gaming laptop",
    "price": 15000000,
    "stock": 10,
    "category_id": 1
  }'
```

## 🏗️ Architecture

Project ini menggunakan Clean Architecture dengan pembagian yang jelas:

- **cmd/**: Entry point aplikasi
- **internal/config/**: Konfigurasi dan koneksi database
- **internal/models/**: Data models dan business entities
- **internal/services/**: Business logic layer
- **internal/controllers/**: HTTP handlers (presentation layer)
- **internal/middleware/**: HTTP middleware
- **internal/routes/**: Route definitions
- **pkg/utils/**: Utility functions yang dapat digunakan kembali

## 🔒 Security

- CORS middleware untuk cross-origin requests
- Input validation menggunakan validator package
- Error handling yang aman
- Environment variables untuk sensitive data
- Placeholder untuk JWT authentication

## 🚀 Production Deployment

### Build untuk Production

```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-server cmd/server/main.go
```

### Environment Variables untuk Production

```env
ENV=production
PORT=8080
DB_HOST=your-production-db-host
DB_NAME=your-production-db-name
JWT_SECRET=your-very-secure-jwt-secret
```

## 📈 Future Enhancements

- [ ] JWT Authentication implementation
- [ ] Role-based access control
- [ ] API documentation dengan Swagger
- [ ] Unit tests dan integration tests
- [ ] Docker containerization
- [ ] Caching dengan Redis
- [ ] Rate limiting
- [ ] Logging dengan structured logging
- [ ] Monitoring dan metrics

## 🤝 Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License.
