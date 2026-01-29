# BESQ Backend

A Golang backend API server using Gin web framework, GORM ORM, and SQLite database.

## Features

- RESTful API with Gin framework
- SQLite database for development
- User authentication system
- Role-based access (admin, cutting, pressing)
- Password hashing with bcrypt
- Auto-migration and database seeding

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

### Default Admin User

A default admin user is automatically created on first run:
- **NIK**: `admin`
- **Password**: `admin123`
- **Role**: `admin`

## API Endpoints

### Public Routes (No Authentication)
- `GET /health` - Health check endpoint
- `GET /api/` - Welcome message
- `POST /login` - User login

### Protected User Routes (Authentication Required) - `/api/users`
**All routes in this group require valid JWT token via Authorization header**

#### Any Authenticated User Can Access:
- `GET /api/users/profile` - Get current user profile
- `GET /api/users/` - Get all users
- `GET /api/users/:id` - Get specific user by ID

#### Admin-Only Routes (Admin Role Required):
- `POST /api/users/` - Create new user
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

## User Model

The User model includes the following fields:

```go
type User struct {
    ID        uint      `json:"id"`
    NIK       string    `json:"nik"`       // Unique identifier for login
    Name      string    `json:"name"`
    Password  string    `json:"-"`         // Hidden in JSON responses
    Role      string    `json:"role"`      // admin, cutting, pressing
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Valid Roles
- `admin` - Administrator access
- `cutting` - Cutting department access
- `pressing` - Pressing department access

## Middleware

The application includes JWT authentication middleware:

### AuthMiddleware()
- Validates Authorization header (Bearer token)
- Parses JWT token and extracts user information
- Sets `role` and `userId` in Gin context
- Returns 401 Unauthorized for invalid tokens

### AdminOnly()
- Requires AuthMiddleware() to run first
- Checks if user role is 'admin'
- Returns 403 Forbidden for non-admin users
- Used to protect CRUD operations

### RoleOnly(roles...)
- Allows access only to users with specific roles
- Example: `RoleOnly("admin", "cutting")` allows admin and cutting roles

## API Usage Examples

### 1. Login (Get JWT Token) - Public
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"nik":"admin","password":"admin123"}'
```
**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 2. Get Profile (Protected Route)
```bash
curl -X GET http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3. Get All Users (Any Authenticated User)
```bash
curl -X GET http://localhost:8080/api/users/ \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. Get Specific User (Any Authenticated User)
```bash
curl -X GET http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 5. Create User (Admin Only)
```bash
curl -X POST http://localhost:8080/api/users/ \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"nik":"user001","name":"John Doe","password":"password123","role":"cutting"}'
```

### 6. Update User (Admin Only)
```bash
curl -X PUT http://localhost:8080/api/users/2 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated Name","role":"pressing","password":"newpassword"}'
```

### 7. Delete User (Admin Only)
```bash
curl -X DELETE http://localhost:8080/api/users/2 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Database

The application uses SQLite for development with the database file `besq.db` created automatically in the project root.

## Security Features

- **JWT Authentication**: Stateless authentication using JSON Web Tokens
- **Password Hashing**: bcrypt with salt for secure password storage
- **Role-Based Access Control**: Admin, cutting, and pressing roles with middleware protection
- **Token Expiration**: JWT tokens expire in 24 hours
- **CORS Support**: Cross-origin requests enabled for frontend integration
- **Input Validation**: Request validation using Gin's binding features

## Error Responses

The API returns consistent error responses:

### 400 Bad Request
```json
{
  "error": "Invalid request format"
}
```

### 401 Unauthorized
```json
{
  "error": "Authorization header is required"
}
```

### 403 Forbidden
```json
{
  "error": "Access denied. Admin role required."
}
```

### 404 Not Found
```json
{
  "error": "User not found"
}
```

## Project Structure

```
backend/
├── main.go              # Main application entry point
├── go.mod               # Go module dependencies
├── controllers/
│   ├── auth.go          # JWT authentication controller
│   └── user.go          # User CRUD operations controller
├── middleware/
│   └── auth.go          # Authentication and authorization middleware
├── models/
│   └── user.go          # User model and related structs
├── .gitignore           # Git ignore patterns
└── README.md            # This documentation
```