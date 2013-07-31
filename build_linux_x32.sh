CGO_ENABLED=0 GOOS=linux GOARCH=386 go build main.go
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build setup.go
rm -rf output
mkdir output
chmod +x main
mv main output/godaily
mv setup output/
cp config.ini output/
cp createdb.sql output/
cp -r static output/
cp -r templates output/
