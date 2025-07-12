# HelloMix Backend Debugging Guide

## Overview
This guide explains how to debug your HelloMix Go backend using different approaches.

## Debugging Options

### 1. **Local Development with Delve (Recommended)**

**Prerequisites:**
- Go installed locally
- PostgreSQL and Redis running locally
- Delve debugger installed: `go install github.com/go-delve/delve/cmd/dlv@latest`

**Steps:**
1. Start your local database services:
   ```bash
   # PostgreSQL (port 5432)
   # Redis (port 6379)
   ```

2. Run the local debug script:
   ```bash
   ./debug-backend-local.sh
   ```

3. Connect VS Code debugger:
   - Open the project in VS Code
   - Go to Run and Debug (⌘+Shift+D)
   - Select "Attach to Docker" configuration
   - Press F5 to start debugging

**Configuration in VS Code:**
Your `.vscode/launch.json` is already configured:
```json
{
    "name": "Attach to Docker",
    "type": "go",
    "request": "attach",
    "mode": "remote",
    "remotePath": "/app",
    "port": 2345,
    "host": "localhost"
}
```

### 2. **Docker Development with Hot Reload**

For development with automatic reloading but without step-by-step debugging:

```bash
docker-compose -f docker-compose.hotreload.yml up --build
```

This uses Air for hot reload - your code changes will automatically restart the backend.

### 3. **Docker with Delve Debugging**

If you prefer Docker-based debugging:

```bash
docker-compose -f docker-compose.debug.yml up --build
```

**Note:** The current Docker setup might have issues. Use the local debugging approach instead.

### 4. **Enhanced Logging**

Your backend already has enhanced logging capabilities:

- **Debug Mode:** Set `SERVER_MODE=debug` to enable detailed logging
- **Request/Response Logging:** Automatic in debug mode
- **Structured Logging:** Uses logrus for structured log output

### 5. **API Testing**

Use the debug API testing script:
```bash
./debug-api.sh
```

This tests all your API endpoints and shows request/response details.

## Setting Breakpoints

### In VS Code:
1. Open any Go file in your backend
2. Click on the line number where you want to set a breakpoint
3. The red dot indicates a breakpoint is set
4. Start debugging (F5)
5. Make API calls to trigger your breakpoints

### Common Breakpoint Locations:
- `cmd/server/main.go` - Application startup
- `internal/api/handlers/*.go` - API request handlers
- `internal/services/*.go` - Business logic
- `internal/database/database.go` - Database operations

## Debug Commands

Once connected to the debugger, you can use these commands:

- **Continue (F5):** Resume execution
- **Step Over (F10):** Execute next line
- **Step Into (F11):** Step into function calls
- **Step Out (Shift+F11):** Step out of current function
- **Restart (Ctrl+Shift+F5):** Restart debugging session

## Troubleshooting

### "Can't connect to debugger"
1. Ensure the debug server is running (`./debug-backend-local.sh`)
2. Check that port 2345 is not in use: `lsof -i :2345`
3. Verify Delve is installed: `dlv version`

### "Backend not starting"
1. Check database connections (PostgreSQL and Redis)
2. Verify environment variables are set
3. Check the console output for error messages

### "Breakpoints not working"
1. Ensure you're running in debug mode
2. Check that source paths match between local and remote
3. Rebuild the application if needed

## Environment Variables

For local debugging, set these environment variables:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=hellomix
export DB_USER=hellomix
export DB_PASSWORD=hellomix_password
export DB_SSLMODE=disable
export REDIS_HOST=localhost
export REDIS_PORT=6379
export SERVER_MODE=debug
export SERVER_PORT=8080
```

## Example Debug Session

1. Start the debugger:
   ```bash
   ./debug-backend-local.sh
   ```

2. Set breakpoints in VS Code

3. Connect VS Code debugger (F5)

4. Make API calls:
   ```bash
   curl http://localhost:8080/api/health
   ```

5. Step through your code using VS Code debugger controls

## Log Analysis

Your backend logs include:
- Request/response details
- Database query logs
- Error stack traces
- Performance metrics

Monitor logs in real-time:
```bash
tail -f backend/logs/app.log
```

## Need Help?

If you encounter issues:
1. Check the console output for error messages
2. Verify all prerequisites are installed
3. Ensure database services are running
4. Check port availability (2345 for debugger, 8080 for API)
