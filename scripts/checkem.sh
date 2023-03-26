## Example script for checking multiple decks
declare -a arr=("yakujakku" "jediguest")

## now loop through the above array
for i in "${arr[@]}"
do
    echo "$i"
go run main.go --deck ../mtg-sealed-league/one/pools/${i}/deck.txt --pool ../mtg-sealed-league/one/pools/${i}/final.txt
    echo ""
done
