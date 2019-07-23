TEST_NAMESPACE=mongodb-terraform-builder

default: build

clean:
	-docker ps -f 'label=ns=$(TEST_NAMESPACE)' -q | xargs docker rm -f
	-docker images -f 'label=ns=$(TEST_NAMESPACE)' -q | xargs docker rmi -f

build-image: clean
	@docker build -t $(TEST_NAMESPACE):latest --label "ns=$(TEST_NAMESPACE)" .

build: build-image
	@$(eval SOURCE_CODE = $(shell realpath ../../))
	@docker run --mount type=bind,source=$(SOURCE_CODE),target=/code $(TEST_NAMESPACE):latest

.PHONY: default clean build-image build
