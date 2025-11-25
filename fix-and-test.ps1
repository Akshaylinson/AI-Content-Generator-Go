#!/usr/bin/env pwsh

Write-Host "🔧 Fix and Test System" -ForegroundColor Green

Write-Host "1. Checking model detection..." -ForegroundColor Yellow
try {
    $modelStatus = Invoke-RestMethod -Uri "http://localhost:8080/api/model-status" -Method GET
    Write-Host "✅ Model Status: $($modelStatus.data.message)" -ForegroundColor Green
}
catch {
    Write-Host "❌ Model API Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n2. Creating test job..." -ForegroundColor Yellow
try {
    $testJob = @{ 
        topic = "Test: AI and Machine Learning"
        type = "blog" 
    } | ConvertTo-Json
    
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/jobs" -Method POST -ContentType "application/json" -Body $testJob
    
    if ($response.success) {
        $jobId = $response.data.id
        Write-Host "✅ Created test job ID: $jobId" -ForegroundColor Green
        
        Write-Host "`n3. Force processing..." -ForegroundColor Yellow
        $processResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/process" -Method POST
        Write-Host "✅ $($processResponse.data.message)" -ForegroundColor Green
        
        Write-Host "`n4. Waiting for processing (15 seconds)..." -ForegroundColor Yellow
        Start-Sleep -Seconds 15
        
        Write-Host "`n5. Checking result..." -ForegroundColor Yellow
        $jobResult = Invoke-RestMethod -Uri "http://localhost:8080/api/job/$jobId" -Method GET
        
        Write-Host "Job Status: $($jobResult.data.status)" -ForegroundColor Cyan
        
        if ($jobResult.data.output) {
            Write-Host "`n📄 Generated Content:" -ForegroundColor Green
            $preview = if ($jobResult.data.output.Length -gt 300) { 
                $jobResult.data.output.Substring(0, 300) + "..." 
            } else { 
                $jobResult.data.output 
            }
            Write-Host $preview -ForegroundColor White
            Write-Host "`n✅ SUCCESS: Content generated!" -ForegroundColor Green
        } else {
            Write-Host "❌ No content generated" -ForegroundColor Red
        }
    }
}
catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n🌐 Dashboard: http://localhost:8080" -ForegroundColor Green
