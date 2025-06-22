all_packages=$(go list ./modules/...)

for pkg in $all_packages; do
    if [[ "$pkg" == "github.com/TrueHopolok/VladOS/modules/vos" ]]; then
        go test "$pkg"
    else
        go test "$pkg" -config=./configs/test.cfg
    fi
done