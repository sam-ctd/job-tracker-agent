# Smoke API — job-tracker-agent
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$api = if ($env:SMOKE_API) { $env:SMOKE_API.TrimEnd("/") } else { "http://localhost:8088" }

$keyLine = Get-Content (Join-Path $root ".env") -ErrorAction Stop |
  Where-Object { $_ -match "^APP_API_KEY=" } |
  Select-Object -First 1
if (-not $keyLine) { throw "APP_API_KEY manquant dans .env" }
$key = ($keyLine -replace "^APP_API_KEY=", "").Trim()
$auth = @{ Authorization = "Bearer $key"; "Content-Type" = "application/json" }

function Ok($name, $code) {
  if ($code -ge 200 -and $code -lt 300) { Write-Host "OK  $name ($code)" -ForegroundColor Green }
  else { throw "FAIL $name ($code)" }
}

Write-Host "API = $api"
$r = Invoke-WebRequest "$api/health" -UseBasicParsing
Ok "/health" $r.StatusCode

$r = Invoke-WebRequest "$api/ready" -UseBasicParsing
Ok "/ready" $r.StatusCode

$r = Invoke-WebRequest "$api/api/stats" -Headers $auth -UseBasicParsing
Ok "/api/stats" $r.StatusCode

$body = '{"message":"muéstrame las stats"}'
$r = Invoke-WebRequest "$api/api/chat" -Method POST -Headers $auth -Body $body -UseBasicParsing
Ok "/api/chat" $r.StatusCode

$create = '{"company":"SmokeCo","title":"Smoke Test","channel":"linkedin","status":"idea"}'
$r = Invoke-WebRequest "$api/api/applications" -Method POST -Headers $auth -Body $create -UseBasicParsing
Ok "POST /api/applications" $r.StatusCode
$id = ($r.Content | ConvertFrom-Json).id

$r = Invoke-WebRequest "$api/api/applications/$id" -Method PATCH -Headers $auth -Body '{"status":"enviada"}' -UseBasicParsing
Ok "PATCH status" $r.StatusCode

$r = Invoke-WebRequest "$api/api/applications/$id" -Method DELETE -Headers $auth -UseBasicParsing
Ok "DELETE" $r.StatusCode

Write-Host "Smoke API OK" -ForegroundColor Green
