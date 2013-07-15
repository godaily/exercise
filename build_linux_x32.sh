CGO_ENABLED=0 GOOS=linux GOARCH=386 go build main.go
rm main_linux_x32
mv main main_linux_x32
chmod +x main_linux_x32
