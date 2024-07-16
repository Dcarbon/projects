TAG=dcarbon/projects:prod.1

docker build -t $TAG .
if [[ "$1" == "push" ]];then
    docker push $TAG
fi
