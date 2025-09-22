mkdir build/
rm -r build/*

go build -o build/vlados.exe github.com/TrueHopolok/VladOS

touch build/run.sh
echo "./vlados.exe -config=configs/release.cfg" > build/run.sh
chmod 777 build/run.sh 

cp -r static/ build/static/

mkdir build/configs/
cp configs/release.cfg build/configs/release.cfg
cp configs/bot.key build/configs/bot.key

zip -r build/vlados.zip build