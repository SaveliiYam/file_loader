# Build
docker build -t my-go-app .

# Get binaries
docker run --rm -v "$PWD":/output my-go-app cp /app/output.exe /output/output.exe
