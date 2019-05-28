#! /bin/bash

repo=$1
ss=$(echo $repo | tr "/" " ")
IFS=' ' read -a arr <<< "${ss}"

len=${#arr[@]}
end=${arr[len-1]}
name=${end%????} # remove '.git'

lower_name=$(echo "$name" | awk '{print tolower($0)}')

port=$2
if [ -z "$port" ]
then
    port=8080 # if no port provided default to 8080
fi

# find available port to map service to
# https://unix.stackexchange.com/questions/55913/whats-the-easiest-way-to-find-an-unused-local-port
eport=$RANDOM
quit=0

while [ "$quit" -ne 1 ]; do
  netstat -a | grep $eport >> /dev/null
  if [ $? -gt 0 ]; then
    quit=1
  else
    eport=`expr $eport + 1`
  fi
done

echo image name - $lower_name
echo repo - $1
echo internal port - $port
echo external port - $eport

docker image build -t $lower_name \
    --build-arg repo=$repo \
    --build-arg port=$port \
    --build-arg name=$name .

docker run -p $eport:$port --name $lower_name $lower_name