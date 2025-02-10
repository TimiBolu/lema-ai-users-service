# Lema AI Users Service API Documentation

## Overview
The Lema AI Users Service provides a RESTful API for managing users and their posts.

**Base URL:**
```
http://localhost:8000/api
```

---

## **Endpoints**

### 1️⃣ Get Users
**Endpoint:**
```
GET /users?pageNumber={pageNumber}&pageSize={pageSize}
```
**Description:** Retrieves a paginated list of users.

**Query Parameters:**
| Parameter   | Type  | Required | Default | Description                     |
|------------|------|----------|---------|---------------------------------|
| `pageNumber` | int  | No       | `1`     | Page number                     |
| `pageSize`   | int  | No       | `4`    | Number of users per page        |

**Response Example (200 OK):**
```json
{
  "users": [
    {
      "id": "1",
      "firstname": "John",
      "lastname": "Doe",
      "email": "john@example.com",
      "address": {
        "id": "101",
        "street": "123 Main St",
        "city": "Lagos",
        "state": "LA",
        "zipCode": "10001"
      }
    }
  ],
  "total": 100
}
```

---

### 2️⃣ Get User by ID
**Endpoint:**
```
GET /users/{id}
```
**Description:** Retrieves details of a user by their ID.

**Response Example (200 OK):**
```json
{
  "id": "1",
  "firstname": "John",
  "lastname": "Doe",
  "email": "john@example.com",
  "address": {
    "id": "101",
    "street": "123 Main St",
    "city": "Lagos",
    "state": "LA",
    "zipCode": "10001"
  }
}
```

**Error Responses:**
| Status Code | Message       | Cause                  |
|------------|--------------|------------------------|
| `404`       | "User not found" | User ID does not exist |

---

### 3️⃣ Get Users Count
**Endpoint:**
```
GET /users/count
```
**Description:** Returns the total number of users in the database.

**Response Example (200 OK):**
```json
{
    "count": 100
}
```

---

### 4️⃣ Get Posts
**Endpoint:**
```
GET /posts?userId={userId}
```
**Description:** Returns all posts for a specific user.

**Query Parameters:**
| Parameter   | Type  | Required | Description                     |
|------------|------|----------|---------------------------------|
| `userId`   | string | No       | Filter posts by user ID        |

**Response Example (200 OK):**
```json
{
  "posts": [
    {
      "id": "1",
      "title": "My Post 1",
      "body": "This is the content",
      "userId": "123"
    },
    {
      "id": "2",
      "title": "My Post 2",
      "body": "This is the content",
      "userId": "123"
    }
  ]
}
```

---

### 5️⃣ Create a Post
**Endpoint:**
```
POST /posts
```
**Description:** Creates a new post.

**Request Body (JSON):**
```json
{
    "title": "My Post",
    "body": "This is the content",
    "userId": "123"
}
```

**Response Example (201 Created):**
```json
{
    "id": "10",
    "title": "My Post",
    "body": "This is the content",
    "userId": "123"
}
```

**Error Responses:**
| Status Code | Message                | Cause                          |
|------------|----------------------|--------------------------------|
| `400`       | "Missing required fields" | Title, body, or userId is missing |
| `500`       | "Database error"      | Server-side issue |

---

### 6️⃣ Delete a Post
**Endpoint:**
```
DELETE /posts/{id}
```
**Description:** Deletes a post by ID.

**Response:**
| Status Code | Message              |
|------------|----------------------|
| `204`       | No Content (Success) |
| `404`       | "Post not found" |

---

## **Error Handling**
This API returns errors in the following format:
```json
{
    "error": "User not found"
}
```

Common error codes:
- `400 Bad Request`: Invalid request body or parameters.
- `404 Not Found`: Resource not found.
- `500 Internal Server Error`: Server-side failure.

---

## **Authentication**
- No authentication required for now.
- Future versions may include token-based authentication.

---

## **Pagination**
For paginated responses:
- `pageNumber`: Defaults to `1`
- `pageSize`: Defaults to `10`

Example:
```
GET /users?pageNumber=2&pageSize=5
```

---

## **Example Requests**

### Get Users (`cURL`)
```sh
curl -X GET "http://localhost:8000/api/users?pageNumber=1&pageSize=5"
```

### Create a Post (`cURL`)
```sh
curl -X POST "http://localhost:8000/api/posts" \
     -H "Content-Type: application/json" \
     -d '{
           "title": "New Post",
           "body": "This is a test post",
           "userId": "123"
         }'
```
