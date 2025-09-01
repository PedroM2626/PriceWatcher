# Check if PostgreSQL is installed and running
try {
    $pgService = Get-Service -Name postgresql* -ErrorAction Stop
    if ($pgService.Status -ne "Running") {
        Write-Host "Starting PostgreSQL service..."
        Start-Service -Name $pgService.Name
    }
    
    # Check if the database exists
    $dbExists = psql -U postgres -tAc "SELECT 1 FROM pg_database WHERE datname='pricewatcher'"
    
    if (-not $dbExists) {
        Write-Host "Creating database 'pricewatcher'..."
        psql -U postgres -c "CREATE DATABASE pricewatcher;"
        
        # Create the necessary tables
        Write-Host "Creating database tables..."
        psql -U postgres -d pricewatcher -f "$PSScriptRoot/migrations/001_initial_schema.up.sql"
        
        Write-Host "Database setup completed successfully!"
    } else {
        Write-Host "Database 'pricewatcher' already exists."
    }
} catch {
    Write-Host "Error setting up the database: $_"
    Write-Host "Please ensure PostgreSQL is installed and running, and that the 'postgres' user has the correct permissions."
}
