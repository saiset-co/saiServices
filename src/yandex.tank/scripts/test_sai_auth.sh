#!/usr/bin/env bash

set -e
#set -x

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

function register() {
  local USER=$1

  res=$(curl --request GET 'http://localhost:8800/register' \
  --header 'Token: 12345' \
  --header 'Content-Type: application/json' \
  --data $USER  2>/dev/null)

  echo $res
}

function login() {
    local USER=$1

    res=$(curl --request GET 'http://localhost:8800/login' \
      --header 'Token: 12345' \
      --header 'Content-Type: application/json' \
      --data $USER 2> /dev/null | jq -r ".at.name" 2> /dev/null || echo -n "" )

    echo $res
}

#for i in {1..100}
#do
#  register "{\"key\":$i,\"password\":\"12345\"}"
#done

tokens=()
for i in {1..10}
do
  echo -e "$RED Login {\"key\":$i,\"password\":\"12345\"} $NC"
  token=$(login "{\"key\":$i,\"password\":\"12345\"}")
  echo -e "$GREEN Token: $token $NC"

  tokens[i]=$token
done

rm -rf ammo.txt

for i in {1..100}
do
  userId=$((1 + $RANDOM % 10))
  loginId=$((20 + $RANDOM % 30))
  newUserId=$((1000 + $RANDOM % 2000))

  case $((1 + $RANDOM % 3)) in
    1)
      echo -e "$BLUE Generated ammo $YELLOW $i: $NC register user"
      echo "GET||/register||||||{\"key\":\"$newUserId\",\"password\":\"12345\"}"| ./make_ammo_v2.py >> ammo.txt
      ;;

    2)
      echo -e "$BLUE Generated ammo $YELLOW $i: $NC login user"
      echo "GET||/login||||||{\"key\":\"$loginId\",\"password\":\"12345\"}"| ./make_ammo_v2.py >> ammo.txt
      ;;

    3)
      echo -e "$BLUE Generated ammo $YELLOW $i: $NC check access"
      echo "GET||/access||Token: ${tokens[$userId]}||||{\"collection\":\"test\", \"method\": \"get\"}"| ./make_ammo_v2.py >> ammo.txt
      ;;
  esac
done
