$repo = "clemilsonazevedo/look-news"
$apiUrl = "https://api.github.com/repos/$repo/releases/latest"

Write-Host "Look News - instalador"
Write-Host ""
Write-Host "Buscando a ultima versao..."

$release = Invoke-RestMethod -Uri $apiUrl -Headers @{ "User-Agent" = "look-news-installer" }

$asset = $release.assets | Where-Object { $_.name -like "*.exe" } | Select-Object -First 1

if (-not $asset) {
  Write-Host "Nao encontrei um .exe na ultima release. Confira em:"
  Write-Host "https://github.com/$repo/releases"
  exit 1
}

$setupPath = Join-Path $env:TEMP $asset.name

Write-Host "Baixando..."
Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $setupPath

Write-Host "Executando instalador..."
Start-Process -FilePath $setupPath -Wait

Remove-Item $setupPath -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "Pronto!"