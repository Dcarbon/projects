TAG=dcarbon/projects:prod1

docker build -t $TAG .
if [[ "$1" == "push" ]];then
    docker push $TAG
fi
