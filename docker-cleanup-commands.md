# Docker Cleanup Commands - Remove Inactive Resources (10+ days)

## Individual Commands

### Remove unused containers (stopped for 10+ days)
```bash
docker container prune --filter "until=240h"
```

### Remove unused images (unused for 10+ days)
```bash
docker image prune -a --filter "until=240h"
```

### Remove unused volumes (unused for 10+ days)
```bash
docker volume prune --filter "until=240h"
```

### Remove unused networks (unused for 10+ days)
```bash
docker network prune --filter "until=240h"
```

### Remove build cache (older than 10 days)
```bash
docker builder prune --filter "until=240h"
```

## Combined Cleanup Script

### Complete cleanup (all unused resources 10+ days old)
```bash
docker system prune -a --volumes --filter "until=240h"
```

## Alternative Time Filters

### 7 days (168 hours)
```bash
docker system prune -a --volumes --filter "until=168h"
```

### 30 days (720 hours)
```bash
docker system prune -a --volumes --filter "until=720h"
```

### Specific date (example: before January 1, 2024)
```bash
docker system prune -a --volumes --filter "until=2024-01-01T00:00:00"
```

## Safety Commands (Check before deletion)

### List containers that would be removed
```bash
docker container ls -a --filter "until=240h" --format "table {{.ID}}\t{{.Names}}\t{{.Status}}\t{{.CreatedAt}}"
```

### List images that would be removed
```bash
docker image ls --filter "until=240h" --format "table {{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{.CreatedAt}}\t{{.Size}}"
```

### List volumes that would be removed
```bash
docker volume ls --filter "dangling=true"
```

## Force Cleanup (Use with caution)

### Remove ALL stopped containers (regardless of age)
```bash
docker container prune -f
```

### Remove ALL unused images (regardless of age)
```bash
docker image prune -a -f
```

### Remove ALL unused volumes (regardless of age)
```bash
docker volume prune -f
```

### Nuclear option - remove everything unused
```bash
docker system prune -a --volumes -f
```

## Recommended Safe Approach

1. **Check what would be removed first:**
   ```bash
   docker system df  # Show disk usage
   docker container ls -a  # List all containers
   docker image ls  # List all images
   docker volume ls  # List all volumes
   ```

2. **Run with confirmation (recommended):**
   ```bash
   docker system prune -a --volumes --filter "until=240h"
   ```
   This will prompt you to confirm before deletion.

3. **Or run specific commands one by one:**
   ```bash
   docker container prune --filter "until=240h"
   docker image prune -a --filter "until=240h"
   docker volume prune --filter "until=240h"
   docker network prune --filter "until=240h"
   ```

## Notes

- `240h` = 10 days (24 hours × 10 days)
- `-a` flag includes all images, not just dangling ones
- `--volumes` includes volume cleanup
- `-f` flag forces deletion without confirmation (use carefully)
- Filters use Go duration format: `72h`, `168h`, `720h`, etc.
- Always test on non-production systems first
