## Example script for checking multiple decks
declare -a players=("yakujakku" "guest")
declare -a decks=("TdE0GcrTrt" "UOJVSEusmG")
declare -a pools=("Uo2i0CGxHh" "yJuMuZVGUD")

## now loop through the above array
for i in "${!players[@]}"
do
    echo "${players[$i]}"
go run main.go --deck ${decks[$i]} --pool ${pools[$i]}
    echo ""
done
