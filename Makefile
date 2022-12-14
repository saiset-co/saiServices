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
#	cd ./src/saiStorage && go mod tidy && go build -o ../../microservices/saiStorage/build/sai-storage
#	cd ./src/saiAuth && go mod tidy && go build -o ../../microservices/saiAuth/build/sai-auth
#	cd ./src/saiContractExplorer && go mod tidy && go build -o ../../microservices/saiContractExplorer/build/sai-contract-explorer
#	cd ./src/saiEthIndexer/cmd/app && go mod tidy && go build -o ../../../../microservices/saiEthIndexer/build/sai-eth-indexer
#	cd ./src/saiEthInteraction && go mod tidy && go build -o ../../microservices/saiEthInteraction/build/sai-eth-interaction
#	cp ./src/saiEthManager/config/config.json ./microservices/saiEthManager/build/config.json
#	cp ./src/saiGNMonitor/config/config.json ./microservices/saiGNMonitor/build/config.json
#	cp ./src/saiStorage/config/config.json ./microservices/saiStorage/build/config.json
#	cp ./src/saiAuth/config/config.json ./microservices/saiAuth/build/config.json
#	cp ./src/saiContractExplorer/config/config.json ./microservices/saiContractExplorer/build/config.json
#	cp ./src/saiEthIndexer/config/config.json ./microservices/saiEthIndexer/build/config/config.json
#	cp ./src/saiEthInteraction/config.yml ./microservices/saiEthInteraction/build/config.yml
#	cp ./src/saiEthInteraction/contracts.json ./microservices/saiEthInteraction/build/contracts.json

docker:
	docker-compose -f ./microservices/docker-compose.yml up -d --build

log:
	docker-compose -f ./microservices/docker-compose.yml logs -f

loga:
	docker-compose -f ./microservices/docker-compose.yml logs -f sai-auth

logs:
	docker-compose -f ./microservices/docker-compose.yml logs -f sai-storage

logc:
	docker-compose -f ./microservices/docker-compose.yml logs -f sai-contract-explorer

logi:
	docker-compose -f ./microservices/docker-compose.yml logs -f sai-eth-interaction

sha:
	docker-compose -f ./microservices/docker-compose.yml run --rm sai-auth sh

shs:
	docker-compose -f ./microservices/docker-compose.yml run --rm sai-storage sh

shc:
	docker-compose -f ./microservices/docker-compose.yml run --rm sai-contract-explorer sh
