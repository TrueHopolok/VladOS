go build -o vlados.exe github.com/TrueHopolok/VladOS

zip vlados.zip vlados.exe static/ configs/release.cfg configs/bot.key

rm vlados.exe