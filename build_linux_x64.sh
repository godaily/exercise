CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
rm main_linux_x64
mv main main_linux_x64
chmod +x main_linux_x64
