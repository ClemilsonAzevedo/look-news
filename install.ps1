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
Write-Host "  Look News — Instalador"
Write-Host "══════════════════════════════════════"
Write-Host ""

Write-Host "[1/4] Buscando a última versão..."
try {
  $release = Invoke-RestMethod -Uri $apiUrl -Headers @{ "User-Agent" = "look-news-installer" }
} catch {
  Write-Host "  ✗ Erro ao buscar a release."
  exit 1
}

$asset = $release.assets | Where-Object { $_.name -like "*.exe" } | Select-Object -First 1

if (-not $asset) {
  Write-Host "  ✗ Não encontrei um .exe na última release."
  Write-Host "    Confira em: https://github.com/$repo/releases"
  exit 1
}
Write-Host "  ✓ Versão encontrada: $($asset.name)"

$setupPath = Join-Path $env:TEMP $asset.name
Write-Host "[2/4] Baixando..."
try {
  Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $setupPath
  Write-Host "  ✓ Download concluído"
} catch {
  Write-Host "  ✗ Falha no download."
  exit 1
}

Write-Host "[3/4] Executando instalador..."
$proc = Start-Process -FilePath $setupPath -PassThru
Show-Spinner -Process $proc -Message "Instalando..."
$proc.WaitForExit()

if ($proc.ExitCode -ne 0 -and $null -ne $proc.ExitCode) {
  Write-Host "  ⚠ Instalador finalizou com código $($proc.ExitCode)"
}

Write-Host "[4/4] Limpando arquivos temporários..."
Remove-Item $setupPath -ErrorAction SilentlyContinue
Write-Host "  ✓ Limpeza concluída"

Write-Host ""
Write-Host "══════════════════════════════════════"
Write-Host "  Pronto!"
Write-Host "══════════════════════════════════════"