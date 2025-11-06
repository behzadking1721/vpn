Mock server for local testing

Endpoints provided:
- GET /api/config -> returns contents of config/settings.json
- PUT /api/config -> replaces in-memory config and saves to config/settings.json
- GET /api/servers -> returns servers array from config
- POST /api/servers -> add a server (returns 201)
- PUT /api/servers/{id} -> update server by id
- DELETE /api/servers/{id} -> remove server by id

Run:

```bash
# from repository root
go run ./mock_server
```

The server listens on :8080 and enables CORS for development.
