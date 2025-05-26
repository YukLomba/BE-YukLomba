# YukLomba API Documentation

This document provides comprehensive documentation for the YukLomba API, a platform for managing competitions and user registrations.

## Base URL

```
http://localhost:3300/api
```

## Authentication

Most endpoints require authentication using JWT tokens. Include the token in the Authorization header:

```
Authorization: Bearer <your_token>
```

### Authentication Endpoints

#### Register a new user

```
POST /auth/register
```

**Request Body:**
```json
{
  "username": "string",
  "email": "string",
  "password": "string",
  "university": "string",
  "interests": "string"
}
```

**Response (201 Created):**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": "uuid",
    "username": "string",
    "email": "string"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or missing required fields
- `400 Bad Request`: Email already in use

#### Login

```
POST /auth/login
```

**Request Body:**
```json
{
  "email": "string",
  "password": "string"
}
```

**Response (200 OK):**
```json
{
  "access_token": "string",
  "token_type": "Bearer",
  "expires_in": 86400
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or missing required fields
- `401 Unauthorized`: Invalid credentials

#### Google Login

```
POST /auth/google
```

**Request Body:**
```json
{
  "id_token": "string"
}
```

**Response (200 OK):**
```json
{
  "access_token": "string",
  "token_type": "Bearer",
  "expires_in": 86400
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or missing ID token
- `401 Unauthorized`: Invalid Google token

#### Get User Profile

```
GET /auth/profile
```

**Headers:**
```
Authorization: Bearer <your_token>
```

**Response (200 OK):**
```json
{
  "id": "uuid",
  "email": "string",
  "role": "string"
}
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid token

## User Endpoints

### Get All Users

```
GET /users
```

**Response (200 OK):**
```json
[
  {
    "id": "uuid",
    "username": "string",
    "email": "string",
    "role": "string",
    "university": "string",
    "interests": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
]
```

**Error Responses:**
- `500 Internal Server Error`: Failed to fetch users

### Get User by ID

```
GET /users/:id
```

**Parameters:**
- `id`: User UUID

**Response (200 OK):**
```json
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "role": "string",
  "university": "string",
  "interests": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID format
- `404 Not Found`: User not found

### Get User's Past Competitions

```
GET /users/:id/registration
```

**Parameters:**
- `id`: User UUID

**Response (200 OK):**
```json
[
  {
    "userId": "uuid",
    "competitionId": "uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
]
```

**Error Responses:**
- `400 Bad Request`: Invalid ID format
- `500 Internal Server Error`: Failed to fetch user's past competitions

### Create User

```
POST /users
```

**Request Body:**
```json
{
  "username": "string",
  "email": "string",
  "password": "string",
  "university": "string",
  "interests": "string"
}
```

**Response (201 Created):**
```json
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "role": "string",
  "university": "string",
  "interests": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body
- `500 Internal Server Error`: Failed to create user

### Update User

```
PUT /users/:id
```

**Parameters:**
- `id`: User UUID

**Request Body:**
```json
{
  "username": "string",
  "email": "string",
  "university": "string",
  "interests": "string"
}
```

**Response (200 OK):**
```json
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "role": "string",
  "university": "string",
  "interests": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID format or request body
- `500 Internal Server Error`: Failed to update user

## Competition Endpoints

### Get All Competitions

```
GET /competitions
```

**Response (200 OK):**
```json
[
  {
    "id": "uuid",
    "title": "string",
    "type": "string",
    "description": "string",
    "organizerId": "uuid",
    "deadline": "timestamp",
    "category": "string",
    "rules": "string",
    "eventLink": "string",
    "results": "string",
    "createdAt": "timestamp",
    "updatedAt": "timestamp"
  }
]
```

**Error Responses:**
- `500 Internal Server Error`: Failed to fetch competitions

### Get Competition by ID

```
GET /competitions/:id
```

**Parameters:**
- `id`: Competition UUID

**Response (200 OK):**
```json
{
  "id": "uuid",
  "title": "string",
  "type": "string",
  "description": "string",
  "organizerId": "uuid",
  "deadline": "timestamp",
  "category": "string",
  "rules": "string",
  "eventLink": "string",
  "results": "string",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID format
- `404 Not Found`: Competition not found

### Create Competition

```
POST /competitions
```

**Request Body:**
```json
{
  "title": "string",
  "type": "string",
  "description": "string",
  "organizerId": "uuid",
  "deadline": "timestamp",
  "category": "string",
  "rules": "string",
  "eventLink": "string"
}
```

**Response (201 Created):**
```json
{
  "id": "uuid",
  "title": "string",
  "type": "string",
  "description": "string",
  "organizerId": "uuid",
  "deadline": "timestamp",
  "category": "string",
  "rules": "string",
  "eventLink": "string",
  "results": "string",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body
- `500 Internal Server Error`: Failed to create competition

### Update Competition

```
PUT /competitions/:id
```

**Parameters:**
- `id`: Competition UUID

**Request Body:**
```json
{
  "title": "string",
  "type": "string",
  "description": "string",
  "deadline": "timestamp",
  "category": "string",
  "rules": "string",
  "eventLink": "string",
  "results": "string"
}
```

**Response (200 OK):**
```json
{
  "id": "uuid",
  "title": "string",
  "type": "string",
  "description": "string",
  "organizerId": "uuid",
  "deadline": "timestamp",
  "category": "string",
  "rules": "string",
  "eventLink": "string",
  "results": "string",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID format or request body
- `500 Internal Server Error`: Failed to update competition

### Delete Competition

```
DELETE /competitions/:id
```

**Parameters:**
- `id`: Competition UUID

**Response (200 OK):**
```json
{
  "message": "Competition deleted successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID format
- `500 Internal Server Error`: Failed to delete competition

### Get Competitions by Organizer

```
GET /competitions/organizer/:id
```

**Parameters:**
- `id`: Organizer UUID

**Response (200 OK):**
```json
[
  {
    "id": "uuid",
    "title": "string",
    "type": "string",
    "description": "string",
    "organizerId": "uuid",
    "deadline": "timestamp",
    "category": "string",
    "rules": "string",
    "eventLink": "string",
    "results": "string",
    "createdAt": "timestamp",
    "updatedAt": "timestamp"
  }
]
```

**Error Responses:**
- `400 Bad Request`: Invalid ID format
- `500 Internal Server Error`: Failed to fetch competitions

## Data Models

### User

```json
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "role": "string",
  "university": "string",
  "interests": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "competitions": [Competition],
  "joined_competitions": [Competition]
}
```

### Competition

```json
{
  "id": "uuid",
  "title": "string",
  "type": "string",
  "description": "string",
  "organizerId": "uuid",
  "deadline": "timestamp",
  "category": "string",
  "rules": "string",
  "eventLink": "string",
  "results": "string",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

### Registration

```json
{
  "userId": "uuid",
  "competitionId": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

## Environment Variables

The API requires the following environment variables to be set:

```
# Application Settings
APP_PORT=3300

# Database Settings
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=your_password
DB_NAME=your_database

# Authentication
JWT_SECRET=your_jwt_secret

# Google OAuth2
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
```

## Authentication Flow

1. **Registration**: Users register with email/password or Google OAuth
2. **Login**: Users login to receive a JWT token
3. **API Access**: Include the JWT token in the Authorization header for protected endpoints
4. **Token Validation**: The server validates the token and extracts user information

## Error Handling

All API endpoints return appropriate HTTP status codes:

- `200 OK`: Request succeeded
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request parameters or body
- `401 Unauthorized`: Authentication required or failed
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side error

Error responses follow this format:

```json
{
  "error": "Error message describing what went wrong"
}
```

## Rate Limiting

The API currently does not implement rate limiting.

## Pagination

The API currently does not implement pagination for list endpoints.
