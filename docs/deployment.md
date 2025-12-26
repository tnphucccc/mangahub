# MangaHub Deployment Guide

---

## 1. Overview

This guide covers deploying MangaHub in different environments, from local development to Docker-based production deployment.

**Deployment Options:**

- **Local Development**: Run servers directly with `go run`
- **Docker Development**: Run servers in containers with hot-reload
- **Docker Compose Production**: Multi-container deployment (recommended)
- **Manual Production**: Build binaries and deploy to servers

---

## 2. Prerequisites

### For Local Development

- **Go**: 1.19 or later
- **Node.js**: 18 or later
- **Yarn**: 4.0 or later
- **SQLite3**: Should be available on most systems
- **Git**: For cloning the repository

### For Docker Deployment

- **Docker**: 20.10 or later
- **Docker Compose**: 1.29 or later (or use `docker compose` plugin)

### System Requirements

| Component  | Minimum                      | Recommended           |
| ---------- | ---------------------------- | --------------------- |
| CPU        | 1 core                       | 2+ cores              |
| RAM        | 512MB                        | 2GB+                  |
| Disk Space | 500MB                        | 2GB+                  |
| OS         | Linux, macOS, Windows (WSL2) | Linux (Ubuntu 20.04+) |

---

## 3. Local Development Deployment

### 3.1 Initial Setup

**1. Clone Repository**:

```bash
git clone https://github.com/tnphucccc/mangahub.git
cd mangahub
```

**2. Install Go Dependencies**:

```bash
go mod download
```

**3. Install Node.js Dependencies**:

```bash
yarn install
```

**4. Setup Database**:

```bash
# Run migrations
make migrate-up

# Seed with sample data
make seed
```

You should see:

```
Applied migration: 001_create_users_table
Applied migration: 002_create_manga_table
Applied migration: 003_create_user_progress_table
All migrations applied successfully

Seeding database...
Seeding users...
  Created user: testuser (testuser@example.com)
  Created user: alice (alice@example.com)
  Created user: bob (bob@example.com)
Seeding manga from data/manga.json...
  Imported 200 manga from JSON file
Database seeding completed successfully!
```

### 3.2 Running Backend Servers

**Start all 4 Go servers** (requires 4 terminal windows):

**Terminal 1 - API Server** (HTTP + WebSocket):

```bash
make run-api
```

Expected output:

```
HTTP API Server starting on localhost:8080
WebSocket Server starting on ws://localhost:9093/ws
Endpoints available:
  - Health check: GET /health (HTTP)
  - WebSocket: GET /ws (WebSocket)
  - Register: POST /api/v1/auth/register (HTTP)
  - Login: POST /api/v1/auth/login (HTTP)
```

**Terminal 2 - TCP Server**:

```bash
make run-tcp
```

Expected output:

```
Starting TCP Progress Sync Server...
TCP Progress Sync Server started successfully
Listening on port: 9090
Waiting for client connections...
```

**Terminal 3 - UDP Server**:

```bash
make run-udp
```

Expected output:

```
Starting UDP Notification Server...
UDP Notification Server started successfully
Listening on port: 9091
Ready to broadcast chapter release notifications
```

**Terminal 4 - gRPC Server**:

```bash
make run-grpc
```

Expected output:

```
gRPC Internal Service listening on :9092
```

### 3.3 Running Frontend

**Terminal 5 - Next.js Web App**:

```bash
make js-dev
# OR
yarn workspace @mangahub/web dev
```

Expected output:

```
  ▲ Next.js 16.1.0
  - Local:        http://localhost:3000
  - Network:      http://192.168.1.100:3000

 ✓ Starting...
 ✓ Ready in 2.3s
```

**Access the application**:

- Web UI: http://localhost:3000
- API: http://localhost:8080
- API Docs: http://localhost:8080/health

---

## 4. Docker Development Deployment

### 4.1 Build Docker Images

**Build all server images**:

```bash
docker-compose build
```

This builds the multi-stage Dockerfile:

1. **Builder Stage**: Compiles Go binaries
2. **Final Stage**: Minimal Alpine Linux with binaries

Build time: ~2-5 minutes (first time, cached afterwards)

### 4.2 Run with Docker Compose

**Start all services**:

```bash
make run-all
# OR
docker-compose up -d
```

Expected output:

```
Creating network "mangahub_default" with the default driver
Creating mangahub_tcp-server_1  ... done
Creating mangahub_udp-server_1  ... done
Creating mangahub_grpc-server_1 ... done
Creating mangahub_api-server_1  ... done

✓ All servers started!
  API Server:   http://localhost:8080
  WebSocket:    ws://localhost:9093
  TCP Server:   localhost:9090
  UDP Server:   localhost:9091
  gRPC Server:  localhost:9092
```

**View logs**:

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api-server
docker-compose logs -f tcp-server
```

**Stop all services**:

```bash
make stop-all
# OR
docker-compose down
```

### 4.3 Docker Compose Service Details

```yaml
services:
  api-server:
    ports:
      - "8080:8080" # HTTP API
      - "9093:9093" # WebSocket
    environment:
      - TCP_HOST=tcp-server # Internal Docker network
      - UDP_HOST=udp-server
      - GRPC_HOST=grpc-server
    volumes:
      - ./data:/app/data # Database persistence
      - ./configs:/app/configs # Configuration files
    depends_on:
      - tcp-server
      - udp-server
      - grpc-server

  tcp-server:
    command: /app/tcp-server
    ports:
      - "9090:9090"

  udp-server:
    command: /app/udp-server
    ports:
      - "9091:9091/udp" # Note: UDP protocol

  grpc-server:
    command: /app/grpc-server
    ports:
      - "9092:9092"
    volumes:
      - ./data:/app/data # Shares database with API
      - ./configs:/app/configs
```

---

## 5. Production Deployment

### 5.1 Docker Compose Production Setup

**1. Create Production Configuration**:

Create `configs/prod.yaml`:

```yaml
server:
  host: "0.0.0.0" # Listen on all interfaces
  http_port: "8080"
  tcp_port: "9090"
  udp_port: "9091"
  grpc_port: "9092"
  websocket_port: "9093"

database:
  path: "/app/data/mangahub.db"

jwt:
  secret: "CHANGE-THIS-TO-SECURE-RANDOM-STRING" # ⚠️ CHANGE THIS!
  expiry_days: 7
```

**2. Update docker-compose.yml for Production**:

Create `docker-compose.prod.yml`:

```yaml
version: "3.8"

services:
  api-server:
    build:
      context: .
      dockerfile: docker/Dockerfile
    ports:
      - "8080:8080"
      - "9093:9093"
    environment:
      - CONFIG_PATH=/app/configs/prod.yaml
      - TCP_HOST=tcp-server
      - UDP_HOST=udp-server
      - GRPC_HOST=grpc-server
    volumes:
      - ./data:/app/data
      - ./configs:/app/configs
    restart: unless-stopped
    depends_on:
      - tcp-server
      - udp-server
      - grpc-server

  tcp-server:
    build:
      context: .
      dockerfile: docker/Dockerfile
    command: /app/tcp-server
    ports:
      - "9090:9090"
    environment:
      - CONFIG_PATH=/app/configs/prod.yaml
    volumes:
      - ./configs:/app/configs
    restart: unless-stopped

  udp-server:
    build:
      context: .
      dockerfile: docker/Dockerfile
    command: /app/udp-server
    ports:
      - "9091:9091/udp"
    environment:
      - CONFIG_PATH=/app/configs/prod.yaml
    volumes:
      - ./configs:/app/configs
    restart: unless-stopped

  grpc-server:
    build:
      context: .
      dockerfile: docker/Dockerfile
    command: /app/grpc-server
    ports:
      - "9092:9092"
    environment:
      - CONFIG_PATH=/app/configs/prod.yaml
    volumes:
      - ./data:/app/data
      - ./configs:/app/configs
    restart: unless-stopped

  # Nginx reverse proxy (optional but recommended)
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - api-server
    restart: unless-stopped
```

**3. Deploy**:

```bash
# On production server
git clone https://github.com/tnphucccc/mangahub.git
cd mangahub

# Setup database
make migrate-up
make seed

# Build and start
docker-compose -f docker-compose.prod.yml up -d --build
```

### 5.2 Nginx Reverse Proxy (Recommended)

**nginx.conf** (example):

```nginx
events {
    worker_connections 1024;
}

http {
    upstream api_backend {
        server api-server:8080;
    }

    upstream websocket_backend {
        server api-server:9093;
    }

    server {
        listen 80;
        server_name your-domain.com;

        # HTTP API
        location /api/ {
            proxy_pass http://api_backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # WebSocket
        location /ws {
            proxy_pass http://websocket_backend;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }

        # Frontend (if serving Next.js static build)
        location / {
            root /usr/share/nginx/html;
            try_files $uri $uri/ /index.html;
        }
    }
}
```

### 5.3 Production Checklist

- [ ] Change JWT secret in `configs/prod.yaml`
- [ ] Set up HTTPS/TLS certificates (Let's Encrypt recommended)
- [ ] Configure firewall (allow ports 80, 443, block others externally)
- [ ] Set up automatic backups (database + configs)
- [ ] Configure monitoring and logging
- [ ] Set up log rotation
- [ ] Configure resource limits (Docker memory/CPU limits)
- [ ] Enable rate limiting (Nginx or application level)
- [ ] Set up health checks and auto-restart policies
- [ ] Test disaster recovery procedures

---

## 6. Manual Binary Deployment

### 6.1 Build Binaries

**For Linux Production Server**:

```bash
# Build for Linux (from any OS)
GOOS=linux GOARCH=amd64 go build -o bin/api-server ./cmd/api-server
GOOS=linux GOARCH=amd64 go build -o bin/tcp-server ./cmd/tcp-server
GOOS=linux GOARCH=amd64 go build -o bin/udp-server ./cmd/udp-server
GOOS=linux GOARCH=amd64 go build -o bin/grpc-server ./cmd/grpc-server
```

**For macOS**:

```bash
GOOS=darwin GOARCH=arm64 go build -o bin/api-server ./cmd/api-server
# ... (repeat for other servers)
```

**For Windows**:

```bash
GOOS=windows GOARCH=amd64 go build -o bin/api-server.exe ./cmd/api-server
# ... (repeat for other servers)
```

### 6.2 Deploy to Server

**1. Upload files**:

```bash
# Create deployment package
tar -czf mangahub-deploy.tar.gz bin/ configs/ data/ migrations/

# Upload to server
scp mangahub-deploy.tar.gz user@server:/opt/mangahub/
```

**2. Extract on server**:

```bash
ssh user@server
cd /opt/mangahub
tar -xzf mangahub-deploy.tar.gz
chmod +x bin/*
```

**3. Create systemd services** (Linux):

`/etc/systemd/system/mangahub-api.service`:

```ini
[Unit]
Description=MangaHub API Server
After=network.target

[Service]
Type=simple
User=mangahub
WorkingDirectory=/opt/mangahub
ExecStart=/opt/mangahub/bin/api-server
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

Create similar files for:

- `mangahub-tcp.service`
- `mangahub-udp.service`
- `mangahub-grpc.service`

**4. Start services**:

```bash
sudo systemctl daemon-reload
sudo systemctl enable mangahub-api mangahub-tcp mangahub-udp mangahub-grpc
sudo systemctl start mangahub-api mangahub-tcp mangahub-udp mangahub-grpc

# Check status
sudo systemctl status mangahub-*
```

---

## 7. Environment Variables

### Available Environment Variables

| Variable      | Default              | Description                |
| ------------- | -------------------- | -------------------------- |
| `CONFIG_PATH` | `./configs/dev.yaml` | Path to configuration file |
| `TCP_HOST`    | `localhost`          | TCP server hostname        |
| `UDP_HOST`    | `localhost`          | UDP server hostname        |
| `GRPC_HOST`   | `localhost`          | gRPC server hostname       |
| `DB_PATH`     | `./data/mangahub.db` | Database file path         |
| `JWT_SECRET`  | (from config)        | JWT signing secret         |

### Setting Environment Variables

**Docker Compose**:

```yaml
environment:
  - CONFIG_PATH=/app/configs/prod.yaml
  - TCP_HOST=tcp-server
```

**Systemd**:

```ini
[Service]
Environment="CONFIG_PATH=/opt/mangahub/configs/prod.yaml"
Environment="TCP_HOST=tcp-server"
```

**Shell**:

```bash
export CONFIG_PATH="/opt/mangahub/configs/prod.yaml"
export TCP_HOST="tcp-server"
./bin/api-server
```

---

## 8. Health Checks and Monitoring

### Health Check Endpoints

**API Server**:

```bash
curl http://localhost:8080/health
```

Response:

```json
{
  "status": "ok",
  "service": "MangaHub HTTP API Server",
  "version": "1.0.0"
}
```

### Docker Health Checks

Add to `docker-compose.prod.yml`:

```yaml
services:
  api-server:
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

### Monitoring with Logs

**View Docker logs**:

```bash
# Tail all logs
docker-compose logs -f

# Tail specific service
docker-compose logs -f api-server

# View logs since specific time
docker-compose logs --since 1h api-server

# Filter for errors
docker-compose logs | grep ERROR
```

**Systemd logs**:

```bash
# View service logs
sudo journalctl -u mangahub-api -f

# Filter by date
sudo journalctl -u mangahub-api --since "2025-12-26 10:00:00"

# View all MangaHub services
sudo journalctl -u "mangahub-*" -f
```

---

## 9. Backup and Restore

### Automated Backup Script

Create `scripts/backup.sh`:

```bash
#!/bin/bash

BACKUP_DIR="./backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup database
cp data/mangahub.db $BACKUP_DIR/mangahub_$DATE.db

# Backup configs
tar -czf $BACKUP_DIR/configs_$DATE.tar.gz configs/

# Keep only last 7 days of backups
find $BACKUP_DIR -name "mangahub_*.db" -mtime +7 -delete
find $BACKUP_DIR -name "configs_*.tar.gz" -mtime +7 -delete

echo "Backup completed: $BACKUP_DIR/mangahub_$DATE.db"
```

**Schedule with cron**:

```bash
# Add to crontab
crontab -e

# Add line (backup daily at 2 AM)
0 2 * * * /opt/mangahub/scripts/backup.sh
```

### Restore from Backup

```bash
# Stop services
docker-compose down
# OR
sudo systemctl stop mangahub-*

# Restore database
cp backups/mangahub_20251226_120000.db data/mangahub.db

# Restore configs (if needed)
tar -xzf backups/configs_20251226_120000.tar.gz

# Start services
docker-compose up -d
# OR
sudo systemctl start mangahub-*
```

---

## 10. Scaling and Performance

### Horizontal Scaling

**API Server**: Can be scaled horizontally with load balancer

```bash
docker-compose up -d --scale api-server=3
```

**Requirements for horizontal scaling**:

- Shared database (migrate to PostgreSQL)
- Redis for session storage
- Load balancer (Nginx, HAProxy)
- Shared file storage for database

### Vertical Scaling (Resource Limits)

**Docker resource limits**:

```yaml
services:
  api-server:
    deploy:
      resources:
        limits:
          cpus: "2.0"
          memory: 1G
        reservations:
          cpus: "0.5"
          memory: 512M
```

### Performance Tuning

**SQLite optimization** (for current setup):

```sql
PRAGMA journal_mode = WAL;        -- Write-Ahead Logging
PRAGMA synchronous = NORMAL;      -- Balance safety/performance
PRAGMA cache_size = -64000;       -- 64MB cache
```

**Go optimization**:

```bash
# Build with optimizations
go build -ldflags="-s -w" -o bin/api-server ./cmd/api-server
# -s: Strip debug info
# -w: Strip DWARF debugging info
```

---

## 11. Security Hardening

### Production Security Checklist

- [ ] **HTTPS/TLS**: Use Let's Encrypt for free SSL certificates
- [ ] **JWT Secret**: Use strong random secret (32+ characters)
- [ ] **Firewall**: Block all ports except 80 (HTTP), 443 (HTTPS)
- [ ] **Rate Limiting**: Implement at Nginx level
- [ ] **CORS**: Configure allowed origins
- [ ] **Security Headers**: Add via Nginx
- [ ] **Database Encryption**: Encrypt SQLite file or use encrypted filesystem
- [ ] **Secrets Management**: Use environment variables, never commit secrets
- [ ] **User Input Validation**: Sanitize all inputs
- [ ] **SQL Injection Prevention**: Already using prepared statements ✓

### Nginx Security Headers

```nginx
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
add_header Content-Security-Policy "default-src 'self';" always;
```

### Rate Limiting (Nginx)

```nginx
http {
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

    server {
        location /api/ {
            limit_req zone=api_limit burst=20 nodelay;
            proxy_pass http://api_backend;
        }
    }
}
```

---

## 12. Troubleshooting

### Common Issues

#### Port Already in Use

**Symptom**:

```
Error: listen tcp :8080: bind: address already in use
```

**Solution**:

```bash
# Find process using port
lsof -i :8080   # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Kill process
kill -9 <PID>
```

#### Docker Build Fails

**Symptom**:

```
ERROR: failed to solve: process "/bin/sh -c go build..." did not complete successfully
```

**Solution**:

```bash
# Clear Docker cache and rebuild
docker-compose build --no-cache

# Check Go version in Dockerfile (should be 1.19+)
```

#### Database Locked Error

**Symptom**:

```
Error: database is locked
```

**Solution**:

```bash
# Stop all services accessing database
docker-compose down

# Check for zombie processes
ps aux | grep mangahub

# Restart services
docker-compose up -d
```

#### WebSocket Connection Fails

**Symptom**:

```
WebSocket connection to 'ws://localhost:9093/ws' failed
```

**Check**:

1. Is WebSocket server running? `docker-compose ps`
2. Is port 9093 exposed? Check `docker-compose.yml`
3. Is firewall blocking? `sudo ufw allow 9093`
4. Check browser console for CORS errors

#### Services Can't Communicate (Docker)

**Symptom**:

```
Error: dial tcp: lookup tcp-server: no such host
```

**Solution**:

- Ensure all services in same Docker network
- Use service names (tcp-server, not localhost) for inter-service communication
- Check `depends_on` in docker-compose.yml

---

## 13. Maintenance

### Regular Maintenance Tasks

**Daily**:

- Check logs for errors: `docker-compose logs --tail=100`
- Monitor disk usage: `df -h`

**Weekly**:

- Review backup integrity
- Check for security updates: `docker-compose pull`
- Review application logs for anomalies

**Monthly**:

- Update dependencies: `go get -u ./...` and `yarn upgrade`
- Test disaster recovery procedures
- Review and rotate logs

### Updating MangaHub

**Pull latest changes**:

```bash
git pull origin main

# Rebuild and restart
docker-compose down
docker-compose build
docker-compose up -d
```

**Run new migrations** (if any):

```bash
make migrate-up
```

---

## 14. References

### Documentation

- [Architecture](./architecture.md) - System architecture
- [Database](./database.md) - Database schema and management
- [API Documentation](./api-documentation.md) - HTTP API reference
- [TCP Documentation](./tcp-documentation.md) - TCP protocol
- [UDP Documentation](./udp-documentation.md) - UDP protocol
- [WebSocket Documentation](./websocket-documentation.md) - WebSocket chat
- [gRPC Documentation](./grpc-documentation.md) - gRPC service

### External Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Nginx Documentation](https://nginx.org/en/docs/)
- [Let's Encrypt](https://letsencrypt.org/) - Free SSL certificates
- [Systemd Service Units](https://www.freedesktop.org/software/systemd/man/systemd.service.html)

---

**Last Updated**: 2025-12-26
**Version**: 1.0.0
**Deployment Status**: ✅ Production Ready
