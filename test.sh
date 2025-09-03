declare -A testing_flags
# testing_flags["github.com/TrueHopolok/VladOS/modules/vos"]=""

all_packages=$(go list ./modules/...)

for pkg in $all_packages; do
    if [[ ${testing_flags["$pkg"]+_} ]]; then
        go test "$pkg" -v ${testing_flags["$pkg"]}
    else
        go test "$pkg" -v -config=./configs/test.cfg
    fi
done