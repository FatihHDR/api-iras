# API IRAS - Demo & Testing

## Cara Menjalankan Aplikasi

### 1. Setup Database PostgreSQL
```sql
CREATE DATABASE api_iras;
```

### 2. Update file .env
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=api_iras
DB_SSLMODE=disable
```

### 3. Jalankan aplikasi
```bash
go run cmd/server/main.go
```

Atau build dan jalankan:
```bash
go build -o api-server.exe cmd/server/main.go
./api-server.exe
```

## Testing API Endpoints

### 1. Health Check
```bash
curl -X GET http://localhost:8080/health
```

### 2. Create Category
```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-token" \
  -d '{
    "name": "Electronics",
    "description": "Electronic devices and accessories"
  }'
```

### 3. Get All Categories
```bash
curl -X GET http://localhost:8080/api/v1/categories
```

### 4. Register User
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 5. Create Product
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-token" \
  -d '{
    "name": "Laptop Gaming",
    "description": "High performance gaming laptop",
    "price": 15000000,
    "stock": 10,
    "category_id": 1
  }'
```

### 6. Get Products with Pagination
```bash
curl -X GET "http://localhost:8080/api/v1/products?page=1&limit=10"
```

## Catatan Penting

1. **Database**: Pastikan PostgreSQL sudah running dan database sudah dibuat
2. **Authentication**: Saat ini menggunakan simple token validation (demo-token)
3. **Auto Migration**: Database schema akan otomatis dibuat saat aplikasi pertama kali dijalankan
4. **Environment**: Gunakan file .env untuk konfigurasi yang lebih aman

## Response Format

Semua endpoint menggunakan format response yang konsisten:

```json
{
  "success": true,
  "message": "Success message",
  "data": {
    // actual data
  }
}
```

Untuk error:
```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error information"
}
```
