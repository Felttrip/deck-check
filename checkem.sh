declare -a arr=("lambardi" "lonny" "mc" "mettledrum" "redfort" "shiffelburger" "warlockguest" "yakujakku")

## now loop through the above array
for i in "${arr[@]}"
do
    echo "$i"
go run main.go --deck ../mtg-sealed-league/snc/pools/${i}/deck.txt --pool ../mtg-sealed-league/snc/pools/${i}/week6.txt
    echo ""
done
