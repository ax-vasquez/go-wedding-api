# README

The API used for our wedding site, built in Go. This is a simple API used to save wedding guest preferences
and responses using a stack built from scratch.

## Endpoints

All endpoints are prefixed with the current API version route: `/api/v1`

### `DELETE`

| Method | Endpoint | Parameters | Success Response |
| ------ | -------- | ---------------- | ---------------- |
| `DELETE` | `/user/:id` | `id` | `{ "status": 202, "message": "Deleted user", "data": { "records": 1 } }` |
| `DELETE` | `/user/:id/invitee/:invitee_id` | `id`, `invitee_id` | `{ "status": 202, "message": "", "data": { "records": 1 } }` |

### `GET`

| Method | Endpoint | Query Parameters | Success Response |
| ------ | -------- | ---------------- | ---------------- |
| `GET` | `/users` | `ids=<ID_1>,<ID_2>...` | `{ "status": 200, "message": "", "data": { "users": <MATCHING_USERS> } }` |
| `GET` | `/user/:id/invitees` | none | `{ "status": 200, "message": "", "data": { "users": <INVITED_USERS_FOR_USER> } }` |

### `PATCH`

| Method | Endpoint | Body Fields | Success Response |
| ------ | -------- | ----------- | ---------------- |
| `PATCH` | `/user` | `id` (required), `first_name`, `last_name`, `email`, `is_admin`, `is_going`, `can_invite_others`, `hors_douevres_selection_id`, `entree_selection_id`  | `{ "status": 202, "message": "Updated user", "data": { "records": 1 } }` |

### `POST`

| Method | Endpoint | Parameters | Body Fields | Success Response |
| ------ | -------- | ---------- | ----------- | ---------------- |
| `POST` | `/user` | none | `first_name` (required), `last_name` (required), `email` (required), `is_admin`, `is_going`, `can_invite_others`, `hors_douevres_selection_id`, `entree_selection_id` | `{ "status": 201, "message": "Created new user", "data": { "records": 1 } }` |
| `POST` | `/user/:id/invite-user` | `id` - the inviting user's ID | `first_name` (required), `last_name` (required), `email` (required), `is_admin`, `is_going`, `can_invite_others`, `hors_douevres_selection_id`, `entree_selection_id` | `{ "status": 201, "message": "Created new user", "data": { "records": 1 } }` |

## Resetting your local database

If you get your database into a weird state, it's often simplest to just delete the database and re-create it. _To be clear, this is only a valid
option if you're running locally._ Ideally, we should be hashing all of these weird things out before we are anywhere near deploying things.

To reset the DB, you'll need to do the following steps (in this order):
1. Run `docker-compose down`
1. Delete the `./local-data` directory
1. Run `docker-compose up -d`

You need to make sure the Postgres container has stopped before deleting the `/local-data` directory because it writes to this directory while in
operation. Then, once the container is stopped (and no longer writing to the `/local-data` directory), delete the `/local-data` directory. If
you don't delete the `/local-data` directory before restarting, the Postgres container will simply re-use what's defined there and it won't be reset.

Additionally, if you delete the application database manually (currently `"gorm"`) using a `DELETE DATABASE` command, _the application
will not re-create it on the next startup and it will fail to connect_. **In cases like this, it's best to simply reset your database using the
steps above.**.
