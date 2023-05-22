#!/usr/bin/env bash

set -e
#set -x

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

function save() {
  local USER=$1

  res=$(curl --request GET 'http://localhost:8801/save' \
  --header 'Token: 12345' \
  --header 'Content-Type: application/json' \
  --data $USER  2>/dev/null)

  echo $res
}

function get() {
    local USER=$1

    res=$(curl --request GET 'http://localhost:8801/get' \
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
  echo -e "$RED get {\"key\":$i,\"password\":\"12345\"} $NC"
  token=$(get "{\"key\":$i,\"password\":\"12345\"}")
  echo -e "$GREEN Token: $token $NC"

  tokens[i]=$token
done

rm -rf ammo.txt

for i in {1..100}
do
  loginId=$((1+$i))
 # newUserId=$((1+$i))

  # case $((1 + $RANDOM % 2)) in
  #   1)
      echo -e "$BLUE Generated ammo $YELLOW $i: $NC save value"
      echo "POST||/save||||||{\"key\":\"$loginId\",\"value\":\"12345\"}"| ../make_ammo_v2.py >> ../ammo.txt
      # ;;

    # 2)
      echo -e "$BLUE Generated ammo $YELLOW $i: $NC get"
      echo "GET||/get||||||{\"key\":\"$loginId\",\"password\":\"12345\"}"| ../make_ammo_v2.py >> ../ammo.txt
      # ;;
  # esac
done
