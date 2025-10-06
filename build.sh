mkdir build/
rm -r build/*

go build -o build/vlados.exe github.com/TrueHopolok/VladOS

touch build/run.sh
echo "./vlados.exe -config=configs/release.cfg" > build/run.sh
chmod 777 build/run.sh 

touch build/README.md
echo "To launch: execute **run.sh** script." >> build/README.md
echo "Required a **bot.key** file in **./configs/** directory with a bot's telegram API key to launch." >> build/README.md

cp -r static/ build/static/

mkdir build/configs/
cp configs/release.cfg build/configs/release.cfg

zip -r build/vlados.zip build