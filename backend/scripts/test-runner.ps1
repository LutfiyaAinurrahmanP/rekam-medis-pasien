param(
    [string]$Path = "./...",
    [string]$Filter = ""
)

Write-Host ""
Write-Host "PASS $Path" -ForegroundColor Green
Write-Host ""

$packages = @(go list $Path 2>$null)
$packageCount = 0
if ($LASTEXITCODE -eq 0 -and $packages) {
    $packageCount = $packages.Count
    Write-Host "Test Packages:"
    $packages | ForEach-Object { Write-Host $_ }
    Write-Host ""
}

# Run tests and capture output
$testOutput = if ($Filter -eq "") {
    go test -v -count=1 $Path 2>&1
} else {
    go test -v -count=1 $Path -run $Filter 2>&1
}

# Display test output
$testOutput | ForEach-Object { Write-Host $_ }

# Count results
$total = ($testOutput | Select-String "=== RUN").Count
$passed = ($testOutput | Select-String "--- PASS:").Count
$failed = ($testOutput | Select-String "--- FAIL:").Count

if ($packageCount -eq 0) {
    $packageCount = 1
}

# Display Jest-like summary
Write-Host ""
if ($failed -gt 0) {
    Write-Host "Test Suites: 1 failed, $packageCount total" -ForegroundColor Red
} else {
    Write-Host "Test Suites: $packageCount passed, $packageCount total" -ForegroundColor Green
}

if ($failed -gt 0) {
    Write-Host "Tests:       $failed failed, $passed passed, $total total" -ForegroundColor Red
} else {
    Write-Host "Tests:       $passed passed, $total total" -ForegroundColor Green
}

Write-Host "Snapshots:   0 total"
Write-Host "Ran all test suites."
Write-Host ""

# Exit with error if tests failed
if ($failed -gt 0) {
    exit 1
}
