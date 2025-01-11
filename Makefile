
.PHONY: quicktest
quicktest:
	mkdir -p data
	rm -rf data/*
	cd data && YARN_CONTAINED_FORCE_DOCKER_BUILD=1 go run ../yarn-contained.go init

.PHONY: remove_image
remove_image:
	docker rmi -f jamesrr39/yarncontained

.PHONY: install
install:
	go build -trimpath -o ${shell go env GOBIN}/yarn-contained yarn-contained.go
	YARN_CONTAINED_FORCE_DOCKER_BUILD=1 yarn-contained --version

.PHONY: release
release:
	GITHUB_TOKEN=${GORELEASER_YARN_CONTAINED_GITHUB_TOKEN} goreleaser release --clean --release-header ${GORELEASER_RELEASE_HEADER_FILEPATH}
