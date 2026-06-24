# AlaRent Backend

Backend for AstaRent, a rental housing platform for students and young professionals in Astana.

## Stack

- Go + Gin
- PostgreSQL + GORM
- JWT access tokens + refresh tokens
- WebSocket chat

## Local Setup

Create `.env` in the project root:

```env
PORT=8080
JWT_SECRET=change_me
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=astarent
```

Start PostgreSQL, then run:

```bash
go run ./cmd
```

The API will be available at:

```text
http://localhost:8080
```

## Auth

```text
POST /auth/register
POST /auth/login
POST /auth/refresh
POST /auth/logout
```

Register:

```json
{
  "name": "Owner",
  "email": "owner@test.com",
  "phone": "+77001234567",
  "password": "123456",
  "role": "landlord"
}
```

Login returns:

```json
{
  "access_token": "...",
  "refresh_token": "...",
  "user": {}
}
```

Protected requests use:

```text
Authorization: Bearer ACCESS_TOKEN
```

## Profile

```text
GET /me
PATCH /me
```

Patch example:

```json
{
  "name": "New name",
  "phone": "+77007654321",
  "avatar_url": "https://example.com/avatar.jpg",
  "bio": "Student looking for housing"
}
```

## Apartments

```text
GET /apartaments
GET /apartaments/:id
POST /apartaments
PATCH /apartaments/:id
DELETE /apartaments/:id
```

Create/update/delete require a landlord token. A landlord can update and delete only their own apartments.

Create example:

```json
{
  "title": "2-room apartment",
  "description": "Near university",
  "type": "apartment",
  "price": 220000,
  "district": "Есиль",
  "address": "Mangilik El 10",
  "rooms": 2,
  "floor": 5,
  "has_furniture": true,
  "has_wifi": true,
  "has_washer": true,
  "photo_urls": [
    "https://example.com/photo1.jpg"
  ]
}
```

Filters:

```text
GET /apartaments?q=university&district=Есиль&type=apartment&min_price=100000&max_price=300000&rooms=2&has_wifi=true&has_furniture=true&has_washer=true&sort=price_asc
```

Sort values:

```text
price_asc
price_desc
newest
oldest
```

## Favorites

```text
POST /apartaments/:id/favorite
DELETE /apartaments/:id/favorite
GET /me/favorites
```

## Chat

Start or open a conversation by apartment:

```text
POST /apartaments/:id/conversation
```

List conversations:

```text
GET /conversations
```

Messages:

```text
GET /conversations/:id/messages
POST /conversations/:id/messages
```

Send message:

```json
{
  "text": "Hello, is this apartment available?"
}
```

WebSocket:

```text
GET /ws/conversations/:id
Authorization: Bearer ACCESS_TOKEN
```

Send WebSocket JSON:

```json
{
  "text": "Hello from WebSocket"
}
```
