run:
	docker-compose up --build

test:
	docker-compose -f docker-compose.yml -f docker-compose.override.yml up && \
	cd api_test && \
	go test

unit_test:
	cd models && go test

lint:
	gometalinter.v2 ./...
