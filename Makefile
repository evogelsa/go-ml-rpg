docker-all: docker-build docker-start
	@echo "Done"

docker-build:
	@echo "Building the image from docker file..."
	docker build --no-cache --pull -t go-ml-rpg .
	@echo "Image build complete"

docker-start:
	@echo "Starting service in container..."
	docker run -p 8080:8080 -v ${CURDIR}/src/:/go/src/go-ml-rpg/ -it go-ml-rpg

docker-stop:
	@echo "Stopping the service..."
	docker stop $$(docker ps -alq)
	@echo "Service stopped"

docker-remove:
	@echo "Removing image..."
	docker rmi -f go-ml-rpg
	@echo "Image removed"

docker-clean: docker-stop docker-remove
	@echo "Clean complete"

clean:
	@echo "Removing service files created"
	rm -rf $(CREATED)
