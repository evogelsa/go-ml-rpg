all: docker-clean docker-build docker-create

docker-build:
	docker build --no-cache --pull -t go-ml-rpg .

docker-create:
	docker create \
		-p 8081:8081 \
		-v ${CURDIR}/src/:/go/src/go-ml-rpg/ \
		--name go-ml-rpg \
		go-ml-rpg

docker-start:
	docker start go-ml-rpg

docker-stop:
	docker stop go-ml-rpg

docker-clean:
	-docker rm go-ml-rpg
	-docker rmi go-ml-rpg
