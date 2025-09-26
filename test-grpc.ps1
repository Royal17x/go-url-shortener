param (
    [string]$action = ""
)

$server = "host.docker.internal:50051"

# временная папка для JSON файлов
$tmpDir = "$env:TEMP/grpc"
if (-not (Test-Path $tmpDir)) { New-Item -ItemType Directory -Path $tmpDir | Out-Null }

function RunGrpc($jsonFile, $method) {
    $hostPath = $tmpDir -replace '\\','/'
    docker run --rm -i -v "${hostPath}:/tmp" fullstorydev/grpcurl -plaintext $server -d "/tmp/$jsonFile" $method
}

switch ($action) {

    "list" {
        Write-Host "Available gRPC services and methods:"
        docker run --rm fullstorydev/grpcurl -plaintext $server list
    }

    "register" {
        Write-Host "Registering new user..."
        $jsonFile = "user.json"
        $json = @{
            username = "testuser"
            email    = "test@example.com"
            password = "123456"
        } | ConvertTo-Json -Compress
        $json | Set-Content -Path "$tmpDir/$jsonFile" -Encoding UTF8
        RunGrpc $jsonFile "pb.AuthService/Register"
    }

    "login" {
        Write-Host "Logging in..."
        $jsonFile = "login.json"
        $json = @{
            email    = "test@example.com"
            password = "123456"
        } | ConvertTo-Json -Compress
        $json | Set-Content -Path "$tmpDir/$jsonFile" -Encoding UTF8
        RunGrpc $jsonFile "pb.AuthService/Login"
    }

    "shorten" {
        Write-Host "Shortening URL..."
        $jsonFile = "shorten.json"
        $json = @{
            userId      = 1
            originalUrl = "https://google.com"
        } | ConvertTo-Json -Compress
        $json | Set-Content -Path "$tmpDir/$jsonFile" -Encoding UTF8
        RunGrpc $jsonFile "pb.URLService/ShortenURL"
    }

    "resolve" {
        Write-Host "Resolving short code..."
        $jsonFile = "resolve.json"
        $json = @{
            shortCode = "abc123"
        } | ConvertTo-Json -Compress
        $json | Set-Content -Path "$tmpDir/$jsonFile" -Encoding UTF8
        RunGrpc $jsonFile "pb.URLService/ResolveURL"
    }

    Default {
        Write-Host "Unknown action. Use: list | register | login | shorten | resolve"
    }
}
