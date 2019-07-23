TEST_NAMESPACE=mongodb-terraform
TEST_CONTAINER=qa_ubuntu1804

default: list

list:
	@docker ps --filter 'label=ns=$(TEST_NAMESPACE)'

build-image:
	@echo "Building images..."
	@docker build -t $(TEST_NAMESPACE)/qa_ubuntu1804 --label "ns=$(TEST_NAMESPACE)" --build-arg ssh_pub_key="$(shell cat ~/.ssh/id_rsa.pub)" qa_ubuntu1804

start-container: stop-container
	docker run -d -P --name $(TEST_CONTAINER) $(TEST_NAMESPACE)/$(TEST_CONTAINER)

stop-container:
	@$(eval CONTAINER_ID = $(shell docker ps -aqf "name=$(TEST_CONTAINER)"))
	-docker rm -f "$(CONTAINER_ID)"

remove-all-containers:
	-docker ps -f 'label=ns=$(TEST_NAMESPACE)' -q | xargs docker rm -f

remove-all-images:
	-docker images -f 'label=ns=$(TEST_NAMESPACE)' -q | xargs docker rmi -f

ssh-container:
	@docker exec -it $(TEST_CONTAINER) bash

find-port:
ifndef PORT
	$(error PORT is undefined; call this target with PORT=int)
endif
	@$(eval CONTAINER_ID = $(shell docker ps -aqf "name=$(TEST_CONTAINER)"))
	@echo "$(shell docker inspect --format='{{(index (index .NetworkSettings.Ports "$(PORT)/tcp") 0).HostPort}}' $(CONTAINER_ID))"
