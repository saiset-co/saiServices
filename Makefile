up:
	docker-compose -f ./microservices/docker-compose.yml up -d

down:
	docker-compose -f ./microservices/docker-compose.yml down --remove-orphans

build:
	make service
	make docker

service:
#	cd ./src/saiEthManager && go build -o ../../microservices/saiEthManager/build/sai-eth-manager
#	cd ./src/saiGNMonitor && go build -o ../../microservices/saiGNMonitor/build/sai-gn-monitor
	cd ./src/saiStorage && go build -o ../../microservices/saiStorage/build/sai-storage
	cd ./src/saiAuth && go build -o ../../microservices/saiAuth/build/sai-auth
#	cp ./src/saiEthManager/config/config.json ./microservices/saiEthManager/build/config.json
#	cp ./src/saiGNMonitor/config/config.json ./microservices/saiGNMonitor/build/config.json
	cp ./src/saiStorage/config/config.json ./microservices/saiStorage/build/config.json
	cp ./src/saiAuth/config/config.json ./microservices/saiAuth/build/config.json

docker:
	docker-compose -f ./microservices/docker-compose.yml up -d --build

logs:
	docker-compose -f ./microservices/docker-compose.yml logs -f

logn:
	docker-compose -f ./microservices/docker-compose.yml logs -f sai-auth

sh:
	docker-compose -f ./microservices/docker-compose.yml run --rm sai-auth sh