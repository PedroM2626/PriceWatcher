# Check if psql is available
try {
    $psqlPath = Get-Command psql -ErrorAction Stop | Select-Object -ExpandProperty Source
    Write-Host "Found psql at: $psqlPath"
} catch {
    Write-Host "Error: psql command not found. Please ensure PostgreSQL is installed and added to your PATH."
    exit 1
}

# Database connection parameters
$env:PGPASSWORD = "postgres"
$dbName = "pricewatcher"

# Check if database exists
$dbExists = psql -U postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$dbName'"

if (-not $dbExists) {
    Write-Host "Creating database '$dbName'..."
    psql -U postgres -c "CREATE DATABASE $dbName;"
    
    # Run migrations
    $migrationFile = "$PSScriptRoot\migrations\001_initial_schema.up.sql"
    if (Test-Path $migrationFile) {
        Write-Host "Running migrations..."
        psql -U postgres -d $dbName -f $migrationFile
        Write-Host "Database setup completed successfully!"
    } else {
        Write-Host "Error: Migration file not found at $migrationFile"
        exit 1
    }
} else {
    Write-Host "Database '$dbName' already exists."
}

# Clear the password from environment
$env:PGPASSWORD = ""
