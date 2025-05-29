# YukLomba API Documentation (v2.0)

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
  "username": "string (required, min 1)",
  "email": "string (required, email format)",
  "password": "string (required, min 8)",
  "university": "string (required)",
  "interests": "string (required)"
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
- `400 Bad Request`: Invalid request body
- `400 Bad Request`: Validation errors (returns specific error messages)
- `409 Conflict`: Email already registered

#### Login
```
POST /auth/login
```

**Request Body:**
```json
{
  "email": "string (required)",
  "password": "string (required)"
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
- `400 Bad Request`: Missing email/password
- `401 Unauthorized`: Invalid credentials
- `403 Forbidden`: Account not activated

#### Google OAuth2 Flow

1. Initiate Google Login
```
GET /auth/google
```

**Response:**  
302 Redirect to Google OAuth consent screen

2. OAuth Callback
```
GET /auth/google/callback
```

**Query Parameters:**
- `code`: Authorization code from Google
- `state`: CSRF protection token

**Success Response (200 OK):**
```json
{
  "access_token": "string",
  "token_type": "Bearer",
  "expires_in": 86400
}
```

**Error Responses:**
- `400 Bad Request`: Missing code/state
- `401 Unauthorized`: Invalid OAuth code
- `500 Internal Server Error`: Google API communication failure

#### Complete Registration (Set Role)
```
POST /auth/complete-registration
```

**Headers:**
```
Authorization: Bearer <your_temp_token>
```

**Request Body:**
```json
{
  "role": "string (required, enum: student|organizer)"
}
```

**Response (200 OK):**
```json
{
  "message": "User registration completed successfully",
  "user": {
    "id": "uuid",
    "username": "string",
    "email": "string",
    "role": "string"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid role
- `401 Unauthorized`: Invalid/missing token
- `403 Forbidden`: Role already set

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
  "role": "string",
  "university": "string",
  "interests": "string"
}
```

**Error Responses:**
- `401 Unauthorized`: Missing/invalid token

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

### Competition Endpoints

#### Get All Competitions
```
GET /competitions
```

**Query Parameters:**
- `title`: Filter by title (min 3 chars)
- `type`: Filter by competition type
- `category`: Filter by category
- `before`: Filter competitions with deadline before date (ISO format)
- `after`: Filter competitions with deadline after date (ISO format)

**Response (200 OK):**
```json
{
  "competitions": [
    {
      "id": "uuid",
      "title": "string",
      "type": "string",
      "description": "string",
      "image": ["url1", "url2"],
      "organizer": {
        "id": "uuid",
        "name": "string"
      },
      "deadline": "ISO8601",
      "category": "string",
      "rules": "string",
      "eventLink": "url",
      "results": "string",
      "createdAt": "ISO8601",
      "updatedAt": "ISO8601"
    }
  ],
  "total": 0
}
```

**Error Responses:**
- `400 Bad Request`: Invalid filter parameters
- `500 Internal Server Error`: Failed to fetch competitions

#### Get Competition by ID
```
GET /competitions/:id
```

**Parameters:**
- `id`: Competition UUID

**Response (200 OK):**
```json
{
  "id": "uuid",
  "title": "string (required)",
  "type": "string (required)",
  "description": "string (required)",
  "image": ["url1", "url2"],
  "organizer": {
    "id": "uuid",
    "name": "string"
  },
  "deadline": "ISO8601 (required, future date)",
  "category": "string (required)",
  "rules": "string",
  "eventLink": "url",
  "results": "string",
  "createdAt": "ISO8601",
  "updatedAt": "ISO8601"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid UUID format
- `404 Not Found`: Competition not found

#### Create Competition
```
POST /competitions
```

**Headers:**
```
Authorization: Bearer <organizer_token>
```

**Request Body:**
```json
{
  "title": "string (required)",
  "type": "string (required)",
  "description": "string (required)",
  "image": ["url1", "url2 (required, min 1 url)"],
  "deadline": "ISO8601 (required, future date)",
  "category": "string (required)",
  "rules": "string",
  "eventLink": "url (valid URL)"
}
```

**Response (201 Created):**
```json
{
  "id": "uuid",
  "message": "Competition created successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Validation errors
- `401 Unauthorized`: Missing/invalid token
- `403 Forbidden`: User not an organizer
- `500 Internal Server Error`: Creation failed


#### Update Competition
```
PUT /competitions/:id
```

**Headers:**
```
Authorization: Bearer <organizer_token>
```

**Parameters:**
- `id`: Competition UUID

**Request Body:**
```json
{
  "title": "string (required)",
  "type": "string (required)",
  "description": "string (required)",
  "image": ["url1", "url2 (required, min 1 url)"],
  "deadline": "ISO8601 (required, future date)",
  "category": "string (required)",
  "rules": "string",
  "eventLink": "url (valid URL)",
  "results": "string"
}
```

**Response (200 OK):**
```json
{
  "message": "Competition updated successfully",
  "competition": {
    "id": "uuid",
    "title": "string",
    "deadline": "ISO8601",
    "updatedAt": "ISO8601"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Validation errors
- `401 Unauthorized`: Missing/invalid token  
- `403 Forbidden`: User not authorized
- `404 Not Found`: Competition not found
- `500 Internal Server Error`: Update failed

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
