# Book Search API Documentation

## Overview

The Best Reads API provides endpoints for searching books with a smart fallback mechanism. It first searches the local database, and if no results are found, it queries the Open Library API.

## Base URL

- Development: `http://localhost:8080/api/v1`
- Production: `https://api.bestreads.com/api/v1`

## Endpoints

### 1. Search Books

Search for books by query string.

**Endpoint:** `GET /books/search`

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `q` | string | Yes | Search query (searches in title, author, and description) |
| `limit` | integer | No | Maximum results from Open Library (default: 10, max: 100) |

**Request Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/books/search?q=harry%20potter" \
  -H "Accept: application/json"
```

**Response Example (200 OK):**

```json
[
  {
    "ID": 0,
    "ISBN": "",
    "Title": "Harry Potter and the Order of the Phoenix",
    "Author": "J.K. Rowling",
    "CoverURL": "",
    "RatingAvg": 0,
    "Description": "",
    "ReleaseDate": 2003
  },
  {
    "ID": 1,
    "ISBN": "978-0-439-02340-2",
    "Title": "test0",
    "Author": "test0",
    "CoverURL": "",
    "RatingAvg": 0,
    "Description": "test0",
    "ReleaseDate": 0
  }
]
```

**Error Response (400 Bad Request):**

```json
{
  "error": "Query parameter 'q' is required"
}
```

**Error Response (500 Internal Server Error):**

```json
{
  "error": "Error searching in Open Library: timeout"
}
```

### 2. Health Check

Check if the API is running and healthy.

**Endpoint:** `GET /health`

**Request Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/health" \
  -H "Accept: application/json"
```

**Response Example (200 OK):**

```json
{
  "status": "ok"
}
```

## Search Behavior

### Local Search First

- The API searches the local PostgreSQL database first
- Searches across `title`, `author`, and `description` fields
- Case-insensitive pattern matching

### Open Library Fallback

- If no local results are found, the API queries the Open Library API
- Returns up to 10 results by default
- External results have `ID: 0` to indicate they're not from local database

### Combined Results

- Results from both sources are returned in a single array
- Local results typically appear first

## Response Schema

### Book Object

```typescript
{
  ID: number,              // Unique database ID (0 for external API results)
  ISBN: string,            // ISBN identifier
  Title: string,           // Book title
  Author: string,          // Primary author
  CoverURL: string,        // URL to cover image
  RatingAvg: number,       // Average rating
  Description: string,     // Book description
  ReleaseDate: number      // Publication year
}
```

## Error Codes

| Code | Message | Description |
|------|---------|-------------|
| 400 | Query parameter 'q' is required | The search query parameter was not provided |
| 500 | Error searching in database | Database query failed |
| 500 | Error searching in Open Library | External API request failed |

## Rate Limiting

- No rate limiting currently implemented
- Open Library API has its own rate limits

## CORS

- CORS is not configured by default
- Configure in the application settings if needed

## Examples

### Search for local books

```bash
curl -s "http://localhost:8080/api/v1/books/search?q=test"
```

### Search with Open Library fallback

```bash
curl -s "http://localhost:8080/api/v1/books/search?q=harry+potter"
```

### Search for a specific ISBN

```bash
curl -s "http://localhost:8080/api/v1/books/search?q=978-0-439"
```

## Implementation Details

### Technology Stack

- **Framework:** Fiber (Go)
- **Database:** PostgreSQL with GORM
- **External API:** Open Library API (<https://openlibrary.org/search.json>)
- **HTTP Client:** Resty

### Search Query

Local database search uses case-insensitive LIKE pattern matching:

```sql
WHERE LOWER(title) LIKE LOWER(?) 
  OR LOWER(author) LIKE LOWER(?) 
  OR LOWER(description) LIKE LOWER(?)
```

### Logging

All requests and errors are logged using Zerolog with the following levels:

- `INFO` - Successful searches
- `WARN` - Missing query parameters
- `ERROR` - Search failures

Each log entry includes:

- Request ID
- Timestamp
- Caller information
- Custom fields (query, result count, etc.)

## OpenAPI/Swagger

The API is documented using OpenAPI 3.0.0. You can find the specification in `docs/swagger.yaml`.

To view the interactive documentation:

1. Use Swagger UI with the `docs/swagger.yaml` file
2. Or visit: <https://editor.swagger.io> and import the YAML file

## Future Enhancements

- [ ] Implement caching for frequently searched books
- [ ] Add pagination support
- [ ] Add filtering by author, ISBN, or publication date
- [ ] Implement user ratings and reviews
- [ ] Add book recommendations based on search history
- [ ] Support for multiple languages
- [ ] Implement rate limiting
