#!/bin/bash
SET="set"
for ((i=1;i<=100;i++));  
do
if [ "$1" == "$SET" ]
then
./13_sync --case 4 --set
else
./13_sync --case 4 --time
fi
done
