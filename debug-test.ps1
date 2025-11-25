#!/usr/bin/env pwsh

Write-Host "🔍 Debug Test - Checking System Status" -ForegroundColor Green

# Test 1: Create a job
Write-Host "`n1. Creating test job..." -ForegroundColor Yellow
try {
    $body = @{ topic = "Test Topic" } | ConvertTo-Json
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/jobs" -Method POST -ContentType "application/json" -Body $body
    Write-Host "✅ Job created: ID $($response.data.id)" -ForegroundColor Green
    $jobId = $response.data.id
}
catch {
    Write-Host "❌ Failed to create job: $($_.Exception.Message)" -ForegroundColor Red
    exit
}

# Test 2: Check job immediately
Write-Host "`n2. Checking job status immediately..." -ForegroundColor Yellow
try {
    $job = Invoke-RestMethod -Uri "http://localhost:8080/api/job/$jobId" -Method GET
    Write-Host "Status: $($job.data.status)" -ForegroundColor Cyan
    Write-Host "Output length: $($job.data.output.Length)" -ForegroundColor Cyan
}
catch {
    Write-Host "❌ Failed to get job: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 3: Wait and check again
Write-Host "`n3. Waiting 10 seconds for processing..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

try {
    $job = Invoke-RestMethod -Uri "http://localhost:8080/api/job/$jobId" -Method GET
    Write-Host "Status: $($job.data.status)" -ForegroundColor Cyan
    Write-Host "Output length: $($job.data.output.Length)" -ForegroundColor Cyan
    
    if ($job.data.output -and $job.data.output.Length -gt 0) {
        Write-Host "`n📄 Generated Content Preview:" -ForegroundColor Green
        Write-Host $job.data.output.Substring(0, [Math]::Min(200, $job.data.output.Length))
        if ($job.data.output.Length -gt 200) {
            Write-Host "..." -ForegroundColor Gray
        }
    } else {
        Write-Host "❌ No output generated" -ForegroundColor Red
    }
}
catch {
    Write-Host "❌ Failed to get job: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4: List all jobs
Write-Host "`n4. Listing all jobs..." -ForegroundColor Yellow
try {
    $jobs = Invoke-RestMethod -Uri "http://localhost:8080/api/jobs" -Method GET
    foreach ($job in $jobs.data) {
        Write-Host "ID: $($job.id) | Status: $($job.status) | Topic: $($job.topic)" -ForegroundColor White
    }
}
catch {
    Write-Host "❌ Failed to list jobs: $($_.Exception.Message)" -ForegroundColor Red
}
