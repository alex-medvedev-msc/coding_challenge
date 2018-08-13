run:
	docker-compose up --build

test:
	docker-compose -f docker-compose.yml -f docker-compose.override.yml up -d && \
	sleep 2 && \
	cd api_test && \
	go test && \
	docker-compose down -v

unit_test:
	cd models && go test

lint:
	gometalinter.v2 ./...

bench:
	docker-compose -f docker-compose.yml -f docker-compose.override.yml up -d && \
    sleep 2 && \
    cd api_test && \
    go test -bench=.&& \
    docker-compose down -v
