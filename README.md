# README

The API used for our wedding site, built in Go. This is a simple API used to save wedding guest preferences
and responses using a stack built from scratch.

## Endpoints

All endpoints are prefixed with the current API version route: `/api/v1`

### `DELETE`

| Method | Endpoint | Parameters | Success Response |
| ------ | -------- | ---------------- | ---------------- |
| `DELETE` | `/user/:id` | `id` | `{ "status": 202, "message": "Deleted user", "data": { "records": 1 } }` |

### `GET`

| Method | Endpoint | Query Parameters | Success Response |
| ------ | -------- | ---------------- | ---------------- |
| `GET` | `/ping` | none | `{ "status": 200, "message": "OK", "data": null }` |
| `GET` | `/users` | `ids=<ID_1>,<ID_2>...` | `{ "status": 200, "message": "", "data": { "users": <MATCHING_USERS> } }` |

### `PATCH`

| Method | Endpoint | Body Fields | Success Response |
| ------ | -------- | ----------- | ---------------- |
| `PATCH` | `/user` | `id` (required), `first_name`, `last_name`, `email`, `is_admin`, `is_going`, `can_invite_others`, `hors_douevres_selection_id`, `entree_selection_id`  | `{ "status": 202, "message": "Updated user", "data": { "records": 1 } }` |

### `POST`

| Method | Endpoint | Body Fields | Success Response |
| ------ | -------- | ----------- | ---------------- |
| `POST` | `/user` | `first_name` (required), `last_name` (required), `email` (required), `is_admin`, `is_going`, `can_invite_others`, `hors_douevres_selection_id`, `entree_selection_id` | `{ "status": 201, "message": "Created new user", "data": { "records": 1 } }` |
