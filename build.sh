buildver=$1
docker build -t rate-limiting .
docker tag rate-limiting laurencechau/rate-limiting:$buildver
docker tag rate-limiting laurencechau/rate-limiting:latest
docker push laurencechau/rate-limiting:$buildver
docker push laurencechau/rate-limiting:latest