# new-api Auto Sync and Deploy Script
# Usage: .\sync_newapi.ps1

$ErrorActionPreference = "Stop"

# Config
$LOCAL_PATH = "C:\Users\admin\Desktop\start\code\new-api"

Write-Host "========== new-api Auto Deploy ==========" -ForegroundColor Cyan

# Switch to project directory
Set-Location $LOCAL_PATH

# Get current commit hash before pull
$OLD_COMMIT = git rev-parse HEAD

# 1. Pull latest code
Write-Host "[1/3] Pulling latest code..." -ForegroundColor Yellow
git pull origin main

# Get new commit hash after pull
$NEW_COMMIT = git rev-parse HEAD

# Check if code has changed
if ($OLD_COMMIT -eq $NEW_COMMIT) {
    Write-Host "========== No Changes Detected ==========" -ForegroundColor Green
    Write-Host "Current commit: $NEW_COMMIT" -ForegroundColor Cyan
    Write-Host "Skipping container rebuild." -ForegroundColor Cyan
    Write-Host "Container Status:" -ForegroundColor Cyan
    docker-compose ps
    exit 0
}

Write-Host "Code updated: $OLD_COMMIT -> $NEW_COMMIT" -ForegroundColor Cyan

# 2. Stop old containers
Write-Host "[2/3] Stopping old containers..." -ForegroundColor Yellow
docker-compose down

# 3. Build and start containers
Write-Host "[3/3] Building and starting containers..." -ForegroundColor Yellow
docker-compose up -d --build --pull never

Write-Host "========== Deploy Complete! ==========" -ForegroundColor Green
Write-Host "Container Status:" -ForegroundColor Cyan
docker-compose ps
