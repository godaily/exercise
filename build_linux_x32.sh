CGO_ENABLED=0 GOOS=linux GOARCH=386 gopm build
rm -rf output
mkdir output
chmod +x godaily
mv godaily output/
cp config.ini output/
cp createdb.sql output/
cp -r static output/
cp -r templates output/
