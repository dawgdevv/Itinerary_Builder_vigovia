# Vigovia Itinerary API - Complete Documentation

## Table of Contents

1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [API Base URL](#api-base-url)
4. [Authentication](#authentication)
5. [Data Models](#data-models)
6. [API Endpoints](#api-endpoints)
7. [Error Handling](#error-handling)
8. [Quick Testing Guide](#quick-testing-guide)

---

## Introduction

**Vigovia Itinerary API** is a RESTful API for creating, managing, and exporting travel itineraries. It allows users to organize trips with multiple days and activities, with the ability to export itineraries as PDF documents.

### Key Features

- Create and manage full-fidelity itineraries with user ownership
- Capture hotels, flights, transfers, daily activities, and payment plans
- Maintain inclusions and exclusions for clear client communication
- Update or remove itineraries and append activities to specific days
- Strictly validate payloads to keep data consistent
- Export itineraries to PDF and persist the file in the `/output` folder
- List all itineraries or fetch specific ones; delete when no longer needed

---

## Getting Started

### Prerequisites

- Go 1.18+ installed
- Postman or any REST client
- API running on `http://localhost:8080`

### Start the Server

```bash
cd vigovia-task
go run main.go
```

Expected output:

```
Starting Itinerary Builder API on http://localhost:8080
```

---

## API Base URL

```
http://localhost:8080
```

### Available Endpoints

| Method   | Endpoint                          | Description            | Auth Required |
| -------- | --------------------------------- | ---------------------- | ------------- |
| `GET`    | `/`                               | Health check           | No            |
| `POST`   | `/api/auth/signup`                | User registration      | No            |
| `POST`   | `/api/auth/login`                 | User login             | No            |
| `POST`   | `/api/auth/logout`                | User logout            | Yes           |
| `GET`    | `/api/auth/profile`               | Get user profile       | Yes           |
| `POST`   | `/api/itineraries`                | Create itinerary       | Yes           |
| `GET`    | `/api/itineraries`                | List all itineraries   | Yes           |
| `GET`    | `/api/itineraries/:id`            | Get specific itinerary | Yes           |
| `PUT`    | `/api/itineraries/:id`            | Update itinerary       | Yes           |
| `DELETE` | `/api/itineraries/:id`            | Delete itinerary       | Yes           |
| `POST`   | `/api/itineraries/:id/activities` | Add activity           | Yes           |
| `GET`    | `/api/itineraries/:id/export-pdf` | Export as PDF          | Yes           |

---

## Authentication

The API uses token-based authentication with support for both **Bearer tokens** in headers and **cookies** for easy Postman testing.

### Authentication Flow

1. **Sign Up** or **Login** to receive an authentication token
2. Token is automatically saved in a cookie (`auth_token`) and returned in the response
3. Use the token in subsequent requests either via:
   - **Cookie** (automatically sent by Postman/browsers)
   - **Authorization header** with format: `Bearer {token}`

### Cookie Support

After successful signup or login, the API sets an HTTP cookie named `auth_token` with the following properties:

- **Max Age**: 7 days
- **Path**: `/` (available for all endpoints)
- **HttpOnly**: Yes (prevents JavaScript access for security)
- **SameSite**: Not specified (works with Postman)

**Postman Users**: Cookies are automatically saved and sent with subsequent requests. You don't need to manually set the Authorization header if cookies are enabled.

### Token Extraction Priority

The middleware checks for authentication tokens in this order:

1. **Authorization header** (`Bearer {token}`)
2. **Cookie** (`auth_token`)

This dual approach ensures compatibility with:

- API clients that prefer header-based auth
- Browsers and Postman that support cookies

### User ID Handling

When creating or updating itineraries, the `user_id` is **automatically extracted** from the authenticated user's token. You don't need to manually provide it in the request body for authenticated endpoints.

---

### Itinerary

Complete itinerary object returned by the API:

```json
{
  "id": "20241019150405",
  "user_id": "user-123",
  "title": "Paris City Tour",
  "description": "A 3-day tour of the beautiful city of Paris",
  "start_date": "2024-11-15T00:00:00Z",
  "end_date": "2024-11-17T00:00:00Z",
  "location": "Paris, France",
  "hotels": [
    {
      "name": "Hotel Lumiere",
      "city": "Paris, France",
      "check_in": "2024-11-15T15:00:00Z",
      "check_out": "2024-11-18T11:00:00Z",
      "nights": 3
    }
  ],
  "flights": [
    {
      "airline": "Air France",
      "flight_number": "AF123",
      "departure_city": "New York, USA",
      "departure_airport": "JFK",
      "departure_time": "2024-11-14T21:30:00Z",
      "arrival_city": "Paris, France",
      "arrival_airport": "CDG",
      "arrival_time": "2024-11-15T10:45:00Z"
    }
  ],
  "transfers": [
    {
      "mode": "private car",
      "pickup": "Charles de Gaulle Airport",
      "dropoff": "Hotel Lumiere",
      "pickup_time": "11:15",
      "notes": "Driver will wait at exit gate 2"
    }
  ],
  "payment_plan": [
    {
      "installment_number": 1,
      "amount": 1200.0,
      "currency": "EUR",
      "due_date": "2024-09-01T00:00:00Z",
      "status": "Paid"
    }
  ],
  "inclusions": ["Daily breakfast", "Seine river cruise"],
  "exclusions": ["International airfare", "Travel insurance"],
  "days": [
    {
      "day_number": 1,
      "date": "2024-11-15T00:00:00Z",
      "title": "Arrival and Eiffel Tower",
      "activities": [
        {
          "period": "evening",
          "time": "19:00",
          "title": "Eiffel Tower Dinner",
          "description": "Dinner at 58 Tour Eiffel with panoramic city views",
          "location": "Eiffel Tower",
          "duration": "2 hours"
        }
      ]
    }
  ],
  "created_at": "2024-10-19T15:04:05Z",
  "updated_at": "2024-10-19T15:04:05Z"
}
```

### Request Models

#### CreateItineraryRequest

```json
{
  "user_id": "string (required)",
  "title": "string (required)",
  "description": "string (optional)",
  "start_date": "2024-11-15T00:00:00Z (required, ISO 8601)",
  "end_date": "2024-11-17T00:00:00Z (required, ISO 8601)",
  "location": "string (required)",
  "hotels": [
    {
      "name": "string (required)",
      "city": "string (required)",
      "check_in": "ISO datetime (required)",
      "check_out": "ISO datetime (required)",
      "nights": 3
    }
  ],
  "flights": [
    {
      "airline": "string (required)",
      "flight_number": "string (required)",
      "departure_city": "string (required)",
      "departure_airport": "string (required)",
      "departure_time": "ISO datetime (required)",
      "arrival_city": "string (required)",
      "arrival_airport": "string (required)",
      "arrival_time": "ISO datetime (required)"
    }
  ],
  "transfers": [
    {
      "mode": "string (required)",
      "pickup": "string (required)",
      "dropoff": "string (required)",
      "pickup_time": "string (required)",
      "notes": "string (optional)"
    }
  ],
  "payment_plan": [
    {
      "installment_number": 1,
      "amount": 1200.0,
      "currency": "EUR",
      "due_date": "ISO datetime (required)",
      "status": "string (optional)"
    }
  ],
  "inclusions": ["string"],
  "exclusions": ["string"],
  "days": [
    {
      "day_number": 1,
      "date": "ISO datetime (required)",
      "title": "Day title",
      "activities": [
        {
          "period": "morning | afternoon | evening",
          "time": "HH:MM",
          "title": "string (required)",
          "description": "string (required)",
          "location": "string (required)",
          "duration": "string (optional)"
        }
      ]
    }
  ]
}
```

#### UpdateItineraryRequest

```json
{
  "user_id": "string (optional)",
  "title": "string (optional)",
  "description": "string (optional)",
  "start_date": "ISO datetime (optional)",
  "end_date": "ISO datetime (optional)",
  "location": "string (optional)",
  "hotels": [ ... ],
  "flights": [ ... ],
  "transfers": [ ... ],
  "payment_plan": [ ... ],
  "inclusions": [ ... ],
  "exclusions": [ ... ],
  "days": [ ... ]
}
```

#### AddActivityRequest

```json
{
  "day_number": 1,
  "activity": {
    "time": "14:00",
    "title": "Activity name",
    "description": "Activity description",
    "location": "Activity location",
    "duration": "2 hours"
  }
}
```

---

## API Endpoints

### Authentication Endpoints

#### 1. User Signup

**Endpoint:** `POST /api/auth/signup`

**Content-Type:** `application/json`

**Authentication Required:** No

**Request Body:**

```json
{
  "email": "john.doe@example.com",
  "username": "johndoe",
  "password": "SecurePass123",
  "full_name": "John Doe"
}
```

**Response (201 Created):**

```json
{
  "message": "User created successfully",
  "user": {
    "id": "user-20241019150405-a1b2c3d4",
    "email": "john.doe@example.com",
    "username": "johndoe",
    "full_name": "John Doe",
    "created_at": "2024-10-19T15:04:05Z",
    "updated_at": "2024-10-19T15:04:05Z"
  },
  "token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6"
}
```

**Note:** The token is also automatically set in a cookie named `auth_token`.

---

#### 2. User Login

**Endpoint:** `POST /api/auth/login`

**Content-Type:** `application/json`

**Authentication Required:** No

**Request Body:**

```json
{
  "email": "john.doe@example.com",
  "password": "SecurePass123"
}
```

**Response (200 OK):**

```json
{
  "message": "Login successful",
  "user": {
    "id": "user-20241019150405-a1b2c3d4",
    "email": "john.doe@example.com",
    "username": "johndoe",
    "full_name": "John Doe",
    "created_at": "2024-10-19T15:04:05Z",
    "updated_at": "2024-10-19T15:04:05Z"
  },
  "token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6"
}
```

**Note:** The token is also automatically set in a cookie named `auth_token`.

---

#### 3. User Logout

**Endpoint:** `POST /api/auth/logout`

**Authentication Required:** Yes

**Headers:**

```
Authorization: Bearer {token}
```

_Or use cookie authentication_

**Response (200 OK):**

```json
{
  "message": "Logout successful"
}
```

**Note:** The `auth_token` cookie is automatically cleared.

---

#### 4. Get User Profile

**Endpoint:** `GET /api/auth/profile`

**Authentication Required:** Yes

**Headers:**

```
Authorization: Bearer {token}
```

_Or use cookie authentication_

**Response (200 OK):**

```json
{
  "user": {
    "id": "user-20241019150405-a1b2c3d4",
    "email": "john.doe@example.com",
    "username": "johndoe",
    "full_name": "John Doe",
    "created_at": "2024-10-19T15:04:05Z",
    "updated_at": "2024-10-19T15:04:05Z"
  }
}
```

---

### Itinerary Endpoints

#### 5. Health Check

**Endpoint:** `GET /`

**Description:** Verify API is running

**Response (200 OK):**

```json
{
  "message": "Itinerary API is running successfully!",
  "version": "1.0.0"
}
```

---

#### 6. Create Itinerary

**Endpoint:** `POST /api/itineraries`

**Content-Type:** `application/json`

**Authentication Required:** Yes

**Headers:**

```
Authorization: Bearer {token}
```

_Or use cookie authentication_

**Request Body:**

```json
{
  "user_id": "user-123",
  "title": "Paris City Tour",
  "description": "A 3-day tour of the beautiful city of Paris",
  "start_date": "2024-11-15T00:00:00Z",
  "end_date": "2024-11-17T00:00:00Z",
  "location": "Paris, France",
  "hotels": [
    {
      "name": "Hotel Lumiere",
      "city": "Paris, France",
      "check_in": "2024-11-15T15:00:00Z",
      "check_out": "2024-11-18T11:00:00Z",
      "nights": 3
    }
  ],
  "flights": [
    {
      "airline": "Air France",
      "flight_number": "AF123",
      "departure_city": "New York, USA",
      "departure_airport": "JFK",
      "departure_time": "2024-11-14T21:30:00Z",
      "arrival_city": "Paris, France",
      "arrival_airport": "CDG",
      "arrival_time": "2024-11-15T10:45:00Z"
    }
  ],
  "transfers": [
    {
      "mode": "private car",
      "pickup": "Charles de Gaulle Airport",
      "dropoff": "Hotel Lumiere",
      "pickup_time": "11:15",
      "notes": "Driver will wait at exit gate 2"
    }
  ],
  "payment_plan": [
    {
      "installment_number": 1,
      "amount": 1200.0,
      "currency": "EUR",
      "due_date": "2024-09-01T00:00:00Z",
      "status": "Paid"
    }
  ],
  "inclusions": ["Daily breakfast", "Seine river cruise"],
  "exclusions": ["International airfare", "Travel insurance"],
  "days": [
    {
      "day_number": 1,
      "date": "2024-11-15T00:00:00Z",
      "title": "Arrival and Eiffel Tower",
      "activities": [
        {
          "period": "evening",
          "time": "19:00",
          "title": "Eiffel Tower Dinner",
          "description": "Dinner at 58 Tour Eiffel with panoramic city views",
          "location": "Eiffel Tower",
          "duration": "2 hours"
        }
      ]
    }
  ]
}
```

**Response (201 Created):**

```json
{
  "id": "20241019150405",
  "user_id": "user-123",
  "title": "Paris City Tour",
  "description": "A 3-day tour of the beautiful city of Paris",
  "start_date": "2024-11-15T00:00:00Z",
  "end_date": "2024-11-17T00:00:00Z",
  "location": "Paris, France",
  "hotels": [
    {
      "name": "Hotel Lumiere",
      "city": "Paris, France",
      "check_in": "2024-11-15T15:00:00Z",
      "check_out": "2024-11-18T11:00:00Z",
      "nights": 3
    }
  ],
  "flights": [
    {
      "airline": "Air France",
      "flight_number": "AF123",
      "departure_city": "New York, USA",
      "departure_airport": "JFK",
      "departure_time": "2024-11-14T21:30:00Z",
      "arrival_city": "Paris, France",
      "arrival_airport": "CDG",
      "arrival_time": "2024-11-15T10:45:00Z"
    }
  ],
  "transfers": [
    {
      "mode": "private car",
      "pickup": "Charles de Gaulle Airport",
      "dropoff": "Hotel Lumiere",
      "pickup_time": "11:15",
      "notes": "Driver will wait at exit gate 2"
    }
  ],
  "payment_plan": [
    {
      "installment_number": 1,
      "amount": 1200.0,
      "currency": "EUR",
      "due_date": "2024-09-01T00:00:00Z",
      "status": "Paid"
    }
  ],
  "inclusions": ["Daily breakfast", "Seine river cruise"],
  "exclusions": ["International airfare", "Travel insurance"],
  "days": [
    {
      "day_number": 1,
      "date": "2024-11-15T00:00:00Z",
      "title": "Arrival and Eiffel Tower",
      "activities": [
        {
          "period": "evening",
          "time": "19:00",
          "title": "Eiffel Tower Dinner",
          "description": "Dinner at 58 Tour Eiffel with panoramic city views",
          "location": "Eiffel Tower",
          "duration": "2 hours"
        }
      ]
    }
  ],
  "created_at": "2024-10-19T15:04:05Z",
  "updated_at": "2024-10-19T15:04:05Z"
}
```

**Error Response (400 Bad Request):**

```json
{
  "error": "Field validation error message"
}
```

---

#### 7. List All Itineraries

**Endpoint:** `GET /api/itineraries`

**Authentication Required:** Yes

**Headers:**

```
Authorization: Bearer {token}
```

_Or use cookie authentication_

**Response (200 OK):**

```json
{
  "itineraries": [
    {
      "id": "20241019150405",
      "title": "Paris City Tour",
      "description": "A 3-day tour of the beautiful city of Paris",
      "start_date": "2024-11-15T00:00:00Z",
      "end_date": "2024-11-17T00:00:00Z",
      "location": "Paris, France",
      "days": [],
      "created_at": "2024-10-19T15:04:05Z",
      "updated_at": "2024-10-19T15:04:05Z"
    }
  ]
}
```

---

#### 8. Get Specific Itinerary

**Endpoint:** `GET /api/itineraries/:id`

**Authentication Required:** Yes

**Headers:**

```
Authorization: Bearer {token}
```

_Or use cookie authentication_

**Path Parameters:**

- `id` (string, required): Itinerary ID

**Example:** `GET /api/itineraries/20241019150405`

**Response (200 OK):**

```json
{
  "id": "20241019150405",
  "title": "Paris City Tour",
  "description": "A 3-day tour of the beautiful city of Paris",
  "start_date": "2024-11-15T00:00:00Z",
  "end_date": "2024-11-17T00:00:00Z",
  "location": "Paris, France",
  "days": [
    {
      "day_number": 1,
      "date": "2024-11-15T00:00:00Z",
      "title": "Arrival and Eiffel Tower",
      "activities": []
    }
  ],
  "created_at": "2024-10-19T15:04:05Z",
  "updated_at": "2024-10-19T15:04:05Z"
}
```

**Error Response (404 Not Found):**

```json
{
  "error": "itinerary not found"
}
```

---

#### 9. Update Itinerary

**Endpoint:** `PUT /api/itineraries/:id`

**Content-Type:** `application/json`

**Authentication Required:** Yes

**Headers:**

```
Authorization: Bearer {token}
```

_Or use cookie authentication_

**Path Parameters:**

- `id` (string, required): Itinerary ID

**Request Body (all fields optional):**

**Note:** `user_id` is automatically set from the authenticated user's token and should not be included in the request body.

```json
{
  "title": "Paris City Tour - Extended",
  "description": "A 4-day extended tour",
  "start_date": "2024-11-15T00:00:00Z",
  "end_date": "2024-11-18T00:00:00Z",
  "location": "Paris and Versailles, France"
}
```

**Response (200 OK):**

```json
{
  "id": "20241019150405",
  "title": "Paris City Tour - Extended",
  "description": "A 4-day extended tour",
  "start_date": "2024-11-15T00:00:00Z",
  "end_date": "2024-11-18T00:00:00Z",
  "location": "Paris and Versailles, France",
  "days": [],
  "created_at": "2024-10-19T15:04:05Z",
  "updated_at": "2024-10-19T16:30:20Z"
}
```

**Error Response (404 Not Found):**

```json
{
  "error": "itinerary not found"
}
```

---

#### 10. Add Activity to Itinerary

**Endpoint:** `POST /api/itineraries/:id/activities`

**Content-Type:** `application/json`

**Authentication Required:** Yes

**Headers:**

```
Authorization: Bearer {token}
```

_Or use cookie authentication_

**Path Parameters:**

- `id` (string, required): Itinerary ID

**Request Body:**

```json
{
  "day_number": 1,
  "activity": {
    "period": "evening",
    "time": "20:00",
    "title": "Dinner",
    "description": "Enjoy authentic French cuisine",
    "location": "Latin Quarter",
    "duration": "1.5 hours"
  }
}
```

**Response (200 OK):**

```json
{
  "id": "20241019150405",
  "title": "Paris City Tour",
  "description": "A 3-day tour of the beautiful city of Paris",
  "start_date": "2024-11-15T00:00:00Z",
  "end_date": "2024-11-17T00:00:00Z",
  "location": "Paris, France",
  "days": [
    {
      "day_number": 1,
      "date": "2024-11-15T00:00:00Z",
      "title": "Arrival and Eiffel Tower",
      "activities": [
        {
          "time": "20:00",
          "title": "Dinner",
          "description": "Enjoy authentic French cuisine",
          "location": "Latin Quarter",
          "duration": "1.5 hours"
        }
      ]
    }
  ],
  "created_at": "2024-10-19T15:04:05Z",
  "updated_at": "2024-10-19T15:04:05Z"
}
```

---

#### 11. Export Itinerary to PDF

**Endpoint:** `GET /api/itineraries/:id/export-pdf`

**Authentication Required:** Yes

**Headers:**

```
Authorization: Bearer {token}
```

_Or use cookie authentication_

**Path Parameters:**

- `id` (string, required): Itinerary ID

**Response (200 OK):**

- Content-Type: `application/pdf`
- Returns binary PDF file with filename `itinerary.pdf`

**Error Response (404 Not Found):**

```json
{
  "error": "itinerary not found"
}
```

---

#### 12. Delete Itinerary

**Endpoint:** `DELETE /api/itineraries/:id`

**Authentication Required:** Yes

**Headers:**

```
Authorization: Bearer {token}
```

_Or use cookie authentication_

**Path Parameters:**

- `id` (string, required): Itinerary ID

**Response (200 OK):**

```json
{
  "message": "Itinerary deleted successfully"
}
```

**Error Response (404 Not Found):**

```json
{
  "error": "itinerary not found"
}
```

---

## Error Handling

### HTTP Status Codes

| Code | Status                | Meaning                            |
| ---- | --------------------- | ---------------------------------- |
| 200  | OK                    | Request successful                 |
| 201  | Created               | Resource created successfully      |
| 400  | Bad Request           | Invalid request body or parameters |
| 404  | Not Found             | Resource not found                 |
| 500  | Internal Server Error | Server error                       |

### Error Response Format

```json
{
  "error": "Error description message"
}
```

---

## Quick Testing Guide

### Using Postman

#### Step 1: Import Collection

1. Open Postman
2. Click **Import** → Select `Vigovia_API_Postman_Collection.json`
3. All endpoints will be pre-configured

#### Step 2: Test Endpoints in Order

**Step 1: User Signup (Authentication)**

```
POST http://localhost:8080/api/auth/signup
```

Creates a new user account. The auth token will be saved in both:

- Environment variable `auth_token` (for header-based auth)
- Cookie `auth_token` (automatically sent with future requests)

**Step 2: User Login (Optional - if already signed up)**

```
POST http://localhost:8080/api/auth/login
```

Login with existing credentials. Token is saved in cookie and environment variable.

**Step 3: Get User Profile**

```
GET http://localhost:8080/api/auth/profile
```

Verify authentication is working. No need to manually set Authorization header if using cookies.

**Step 4: Health Check**

```
GET http://localhost:8080/
```

Expected: 200 OK with version info

**Step 5: Create Itinerary**

```
POST http://localhost:8080/api/itineraries
```

Copy the returned `id` for next requests. The `user_id` is automatically set from your authenticated token.

**Step 6: List All Itineraries**

```
GET http://localhost:8080/api/itineraries
```

**Step 7: Get Specific Itinerary**

```
GET http://localhost:8080/api/itineraries/{id}
```

Replace `{id}` with ID from Step 5

**Step 8: Add Activity**

```
POST http://localhost:8080/api/itineraries/{id}/activities
```

**Step 9: Update Itinerary**

```
PUT http://localhost:8080/api/itineraries/{id}
```

**Note:** Don't include `user_id` in the request body - it's automatically set from your auth token.

**Step 10: Export to PDF**

```
GET http://localhost:8080/api/itineraries/{id}/export-pdf
```

**Step 11: Delete Itinerary**

```
DELETE http://localhost:8080/api/itineraries/{id}
```

**Step 12: Logout**

```
POST http://localhost:8080/api/auth/logout
```

Invalidates your token and clears the auth cookie.

### Authentication Tips for Postman

1. **Cookie-based (Recommended for Postman)**:

   - After signup/login, cookies are automatically saved
   - No need to manually set Authorization headers
   - Just run the requests as-is

2. **Header-based (Alternative)**:

   - Use `{{auth_token}}` variable in Authorization header
   - Format: `Bearer {{auth_token}}`
   - Token is automatically saved by test scripts

3. **Checking Cookies in Postman**:
   - Click "Cookies" button (top right, below Send button)
   - View all saved cookies for localhost:8080
   - `auth_token` should appear after signup/login

## Common Issues & Solutions

### Issue: API not running

**Solution:** Run `go run main.go` from the `vigovia-task` directory

### Issue: Port 8080 already in use

**Solution:** Edit `main.go` and change the port number, then restart

### Issue: Unauthorized (401) error

**Solution:**

- Make sure you've signed up or logged in first
- Check that cookies are enabled in Postman
- Verify the `auth_token` cookie exists (click Cookies button)
- Alternatively, manually set Authorization header: `Bearer {{auth_token}}`

### Issue: User ID not appearing in itinerary

**Solution:**

- The `user_id` is now automatically extracted from your auth token
- **Don't** include `user_id` in the request body for Create/Update Itinerary
- The backend automatically sets it from the authenticated user

### Issue: Token not being saved

**Solution:**

- Enable cookies in Postman (Settings > General > Enable cookies)
- The token is saved in both cookie and environment variable
- For header-based auth, use `Bearer {{auth_token}}` format

### Issue: Validation errors on create

**Solution:** Ensure all required fields are provided:

- `title` (required)
- `start_date` (required, ISO 8601 format)
- `end_date` (required, ISO 8601 format)
- `location` (required)

### Issue: Itinerary not found (404)

**Solution:** Verify the itinerary ID exists. Use `GET /api/itineraries` to list all

---

## File Structure

```
vigovia-task/
├── main.go                              # Entry point
├── go.mod                               # Go dependencies
├── Vigovia_API_Postman_Collection.json  # Postman import file (THIS FILE)
├── handlers/
│   └── itinerary_handler.go            # HTTP handlers
├── models/
│   └── itinerary.go                    # Data models
├── services/
│   ├── itinerary_service.go            # Business logic
│   └── pdf_service.go                  # PDF generation
├── routes/
│   └── itinerary_routes.go             # Route definitions
├── storage/
│   └── memory_store.go                 # In-memory storage
├── utils/
│   └── validator.go                    # Validation utilities
└── examples/
    ├── sample_request.json             # Example request
    └── sample_response.json            # Example response
```

---

## Support & Notes

- **Version:** 1.0.0
- **Framework:** Gin (Go web framework)
- **Data Storage:** In-memory (can be upgraded to database)
- **Thread Safety:** All operations are thread-safe with RWMutex

For more information or issues, please refer to the project repository.

---
