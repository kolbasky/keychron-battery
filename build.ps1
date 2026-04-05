param(
    [Parameter(Mandatory=$true)]
    [string]$Version
)

$ErrorActionPreference = "Stop"

if ($Version -notmatch '^\d+\.\d+\.\d+$') {
    Write-Error "Version must be in format x.y.z (e.g., 1.0.0)"
    exit 1
}

$parts = $Version -split '\.'
$major = $parts[0]; $minor = $parts[1]; $patch = $parts[2]; $build = "0"
$fullVersion = "$Version.$build"

Write-Host "Building Keychron Battery Monitor v$Version" -ForegroundColor Cyan

# Ensure we're in the project root
Set-Location $PSScriptRoot
$configPath = "internal\config\config.go"

# ── 1. Toggle DEV_MODE to false ──
if (!(Test-Path $configPath)) {
    Write-Error "Config file not found: $configPath"
    exit 1
}

$configOriginal = Get-Content $configPath -Raw
$configProd = $configOriginal -replace 'DEV_MODE\s*=\s*true', 'DEV_MODE = false'
[System.IO.File]::WriteAllText($configPath, $configProd, (New-Object System.Text.UTF8Encoding($false)))
Write-Host "DEV_MODE set to false" -ForegroundColor Yellow

# ── 2. Build (wrapped in try/finally to guarantee restore) ──
try {
    # Helper: write UTF-8 without BOM
    function Set-UTF8NoBOM($Path, $Content) {
        [System.IO.File]::WriteAllText($Path, $Content, (New-Object System.Text.UTF8Encoding($false)))
    }

    # Update versioninfo.json
    $viPath = "versioninfo.json"
    $viContent = Get-Content $viPath -Raw
    $viContent = $viContent -replace '"FileVersion": "\d+\.\d+\.\d+\.0"', "`"FileVersion`": `"$fullVersion`""
    $viContent = $viContent -replace '"ProductVersion": "v\d+\.\d+\.\d+"', "`"ProductVersion`": `"v$Version`""
    $viContent = $viContent -replace '"Major": \d+', "`"Major`": $major"
    $viContent = $viContent -replace '"Minor": \d+', "`"Minor`": $minor"
    $viContent = $viContent -replace '"Patch": \d+', "`"Patch`": $patch"
    $viContent = $viContent -replace '"Build": \d+', "`"Build`": $build"
    Set-UTF8NoBOM $viPath $viContent
    Write-Host "Updated versioninfo.json" -ForegroundColor Green

    # Ensure goversioninfo is available
    if (!(Get-Command goversioninfo -ErrorAction SilentlyContinue)) {
        Write-Host "Installing goversioninfo..." -ForegroundColor Yellow
        go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest
    }
    $gobin = go env GOBIN
    if ([string]::IsNullOrEmpty($gobin)) { $gobin = "$(go env GOPATH)\bin" }
    if ($env:PATH -notlike "*$gobin*") { $env:PATH = "$gobin;$env:PATH" }

    # Generate resource.syso
    Write-Host "Generating resource.syso..." -ForegroundColor Yellow
    goversioninfo -64 -icon "icon.ico" -manifest "manifest.xml" -o "resource.syso"
    if ($LASTEXITCODE -ne 0) { throw "goversioninfo failed" }

    # Build executable
    Write-Host "Building keybat.exe..." -ForegroundColor Yellow
    go build -ldflags="-H=windowsgui -s -w -buildid=" -trimpath -buildvcs=false -o keybat.exe .
    if ($LASTEXITCODE -ne 0) { throw "go build failed" }

    # Optional: UPX compression (separate file)
    if (Get-Command upx -ErrorAction SilentlyContinue) {
        Write-Host "Compressing with UPX..." -ForegroundColor Yellow
        $upxOut = "keybat-upx.exe"
        $before = (Get-Item keybat.exe).Length
        upx -f --best --lzma -o "$upxOut" keybat.exe 2>$null | Out-Null
        if ($LASTEXITCODE -eq 0) {
            $after = (Get-Item $upxOut).Length
            $saved = [math]::Round(($before - $after) / 1MB, 2)
            $pct = [math]::Round((1 - ($after / $before)) * 100, 1)
            Write-Host "Saved ~$saved MB ($pct% smaller)" -ForegroundColor Green
        }
    }

    # Summary
    $size = [math]::Round((Get-Item keybat.exe).Length / 1MB, 2)
    Write-Host "`nv$Version built successfully!" -ForegroundColor Green
    Write-Host "Output: keybat.exe ($size MB)" -ForegroundColor Cyan
    if (Test-Path "keybat-upx.exe") {
        $upxSize = [math]::Round((Get-Item "keybat-upx.exe").Length / 1MB, 2)
        Write-Host "Compressed: keybat-upx.exe ($upxSize MB)" -ForegroundColor Cyan
    }
}
finally {
    # ── 3. ALWAYS restore DEV_MODE to true ──
    [System.IO.File]::WriteAllText($configPath, $configOriginal, (New-Object System.Text.UTF8Encoding($false)))
    Write-Host "`nRestored DEV_MODE = true in config.go" -ForegroundColor Green
}
