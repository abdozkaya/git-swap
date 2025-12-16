$ErrorActionPreference = 'Stop'
$Repo = "abdozkaya/git-swap"
$FileName = "git-swap-windows-amd64.exe"
$Url = "https://github.com/$Repo/releases/latest/download/$FileName"
$InstallDir = "$env:LOCALAPPDATA\Programs\git-swap"

Write-Host "⬇️  Downloading git-swap..." -ForegroundColor Cyan

# Create Directory
if (!(Test-Path -Path $InstallDir)) {
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
}

# Download
try {
    Invoke-WebRequest -Uri $Url -OutFile "$InstallDir\git-swap.exe"
} catch {
    Write-Host "Error downloading file. Please check internet connection." -ForegroundColor Red
    exit 1
}

# Add to PATH if not exists
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    Write-Host "⚙️  Adding to PATH..." -ForegroundColor Yellow
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    $env:Path += ";$InstallDir"
    Write-Host "Path updated. You might need to restart your terminal." -ForegroundColor Yellow
}

Write-Host "✅ Installed successfully!" -ForegroundColor Green
Write-Host "Run 'git-swap help' to get started."