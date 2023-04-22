up:
	docker-compose -f ./microservices/docker-compose.yml up -d

down:
	docker-compose -f ./microservices/docker-compose.yml down --remove-orphans

build:
	make service
	make docker

service:
	# cd ./src/saiEthManager && go build -o ../../microservices/saiEthManager/build/sai-eth-manager
	# cd ./src/saiGNMonitor && go build -o ../../microservices/saiGNMonitor/build/sai-gn-monitor
	cd ./src/saiStorage && go mod tidy && go build -o ../../microservices/saiStorage/build/sai-storage
	cd ./src/saiAuth && go mod tidy && go build -o ../../microservices/saiAuth/build/sai-auth
	# cd ./src/saiContractExplorer && go mod tidy && go build -o ../../microservices/saiContractExplorer/build/sai-contract-explorer
	# cd ./src/saiEthIndexer/cmd/app && go mod tidy && go build -o ../../../../microservices/saiEthIndexer/build/sai-eth-indexer
	# cd ./src/saiEthInteraction && go mod tidy && go build -o ../../microservices/saiEthInteraction/build/sai-eth-interaction
	# cp ./src/saiEthManager/config/config.json ./microservices/saiEthManager/build/config.json
	# cp ./src/saiGNMonitor/config/config.json ./microservices/saiGNMonitor/build/config.json
	cp ./src/saiStorage/config.json ./microservices/saiStorage/build/config.json
	cp ./src/saiAuth/config.json ./microservices/saiAuth/build/config.json
	# cp ./src/saiContractExplorer/config/config.json ./microservices/saiContractExplorer/build/config.json
	# cp ./src/saiEthIndexer/config/config.json ./microservices/saiEthIndexer/build/config/config.json
	# cp ./src/saiEthInteraction/config.yml ./microservices/saiEthInteraction/build/config.yml
	# cp ./src/saiEthInteraction/contracts.json ./microservices/saiEthInteraction/build/contracts.json

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


## integration tests
integration-test:
	make build-test
	make --ignore-errors integration-test-all
	@make integration-test-down-quiet
	@echo -e "\033[0;34mIntegration tests done"

build-test:
	make service
	make docker-integration-test

docker-integration-test:
	docker-compose -f ./microservices/docker-compose-integration-test.yml up -d --build

integration-test-up:
	docker-compose -f ./microservices/docker-compose-integration-test.yml up -d

integration-test-down:
	docker-compose -f ./microservices/docker-compose-integration-test.yml down --remove-orphans 

integration-test-down-quiet:
	@docker-compose -f ./microservices/docker-compose-integration-test.yml down --remove-orphans > /dev/null 2>&1
	
integration-test-all:
	cd ./src/saiAuth && go clean -testcache && go test -v ./integration-test/


## load tests
load-test:
	make build-test
	cd ./src/vegetaTool && go run .
	@make integration-test-down-quiet
	@echo "\033[0;34mLoad tests done"

## install vegeta
	wget https://github.com/tsenart/vegeta/releases/download/v12.8.4/vegeta_12.8.4_linux_386.tar.gz && tar -xvf vegeta_12.8.4_linux_386.tar.gz



## prepare requests for yandex-tank
yandex-tank-make-ammo:
	echo "POST||/register||||{\"key\":\"1\",\"password\":\"123456\"}"| ./make_ammo.py > ammo.txt

## run yandex-tank to test sai_auth
yandex-tank-run-sai_auth:
	cd src/yandex.tank/scripts && chmod +x test_sai_auth.sh && ./test_sai_auth.sh
	make integration-test-up
	cd src/ && docker run --net host --rm -v ./yandex.tank:/var/loadtest -it yandex/yandex-tank -c load.yaml ammo.txt
	make integration-test-down

## run yandex-tank to test sai_storage
yandex-tank-run-sai_storage:
	cd src/yandex.tank/scripts && chmod +x test_sai_storage.sh && ./test_sai_storage.sh
	make integration-test-up
	cd src/ && docker run --net host --rm -v ./yandex.tank:/var/loadtest -it yandex/yandex-tank -c load.yaml ammo.txt
	make integration-test-down