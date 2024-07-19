# README

The API used for our wedding site, built in Go. This is a simple API used to save wedding guest preferences
and responses using a stack built from scratch.

## Endpoints

All endpoints are prefixed with the current API version route: `/api/v1`

| Method | Endpoint | Success Response |
| ------ | -------- | -------- |
| `GET` | `/ping` | `{ "status": 200, "message": "OK", "data": null }` |
| `GET` | `/users?ids=<ID_1>,<ID_2>...` | `{ "status": 200, "message": "", "data": { "users": <MATCHING_USERS> } }` |
| `PATCH` | `/user` | `{ "status": 202, "message": "Updated user", "data": { "records": 1 } }` |
| `POST` | `/user` | `{ "status": 201, "message": "Created new user", "data": { "records": 1 } }` |
