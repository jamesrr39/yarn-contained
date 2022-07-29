
.PHONY: quicktest
quicktest:
	mkdir -p data
	rm -rf data/*
	cd data && YARN_CONTAINED_FORCE_DOCKER_BUILD=1 go run ../cmd/yarn-contained.go init

.PHONY: remove_image
remove_image:
	docker rmi -f jamesrr39/yarncontained