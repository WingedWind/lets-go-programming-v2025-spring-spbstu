echo "Without tags:"
first_command="go build -o main ./"
echo $first_command

eval $first_command
echo
./main
echo

echo "With \"extended\" tag:"
second_command="go build -o main -tags extended ./"
echo $second_command

eval $second_command
echo
./main