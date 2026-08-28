```powershell
$repo = "clemilsonazevedo/look-news"
$apiUrl = "https://api.github.com/repos/$repo/releases/latest"

function Show-Spinner {
  param(
    [System.Diagnostics.Process]$Process,
    [string]$Message
  )
  $spin = @('⠋','⠙','⠹','⠸','⠼','⠴','⠦','⠧','⠇','⠏')
  $i = 0
  while (-not $Process.HasExited) {
    Write-Host "`r  $($spin[$i]) $Message" -NoNewline
    $i = ($i + 1) % $spin.Count
    Start-Sleep -Milliseconds 100
  }
  Write-Host "`r  ✓ $Message"
}

Write-Host "══════════════════════════════════════"
Write-Host "  Look News — Installer"
Write-Host "══════════════════════════════════════"
Write-Host ""

Write-Host "[1/4] Looking for the latest version..."
try {
  $release = Invoke-RestMethod -Uri $apiUrl -Headers @{ "User-Agent" = "look-news-installer" }
} catch {
  Write-Host "  ✗ Error fetching the release."
  exit 1
}

$asset = $release.assets | Where-Object { $_.name -like "*.exe" } | Select-Object -First 1

if (-not $asset) {
  Write-Host "  ✗ Could not find a .exe in the latest release."
  Write-Host "    Check: https://github.com/$repo/releases"
  exit 1
}
Write-Host "  ✓ Version found: $($asset.name)"

$setupPath = Join-Path $env:TEMP $asset.name
Write-Host "[2/4] Downloading..."
try {
  Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $setupPath
  Write-Host "  ✓ Download completed"
} catch {
  Write-Host "  ✗ Download failed."
  exit 1
}

Write-Host "[3/4] Running installer..."
$proc = Start-Process -FilePath $setupPath -PassThru
Show-Spinner -Process $proc -Message "Installing..."
$proc.WaitForExit()

if ($proc.ExitCode -ne 0 -and $null -ne $proc.ExitCode) {
  Write-Host "  ⚠ Installer exited with code $($proc.ExitCode)"
}

Write-Host "[4/4] Cleaning up temporary files..."
Remove-Item $setupPath -ErrorAction SilentlyContinue
Write-Host "  ✓ Cleanup completed"

Write-Host ""
Write-Host "══════════════════════════════════════"
Write-Host "  Done!"
Write-Host "══════════════════════════════════════"
```
