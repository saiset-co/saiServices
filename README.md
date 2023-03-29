# saiServices
#### profiling
 <host:port/debug/pprof>

## saiEthInteraction

### Config
#### config.yml

eth_server: "" //For all contracts for now  
log_mode: "debug" //Debug mode

#### contracts.json

{  
&emsp;    "name": "", //Contract name, uses in api commands  
&emsp;    "server": "", //Feature update, geth server per contract  
&emsp;    "abi": "", //Contract ABI, escaped json string  
&emsp;    "address": "", //Contract address  
&emsp;    "private": "", //Private key to sign commands  
&emsp;    "gas_limit": 0 //Gas limit for the command transaction  
}

### API
#### Contract command
- request:

curl --location --request GET 'http://localhost:8804' \
&emsp;    --header 'Token: SomeToken' \
&emsp;    --header 'Content-Type: application/json' \
&emsp;    --data-raw '{"method": "api", "data": {"contract":"$name","method":"$contract_method_name", "value": "$value", "params":[{"type":"$(int|string|float...)","value":"$some_value"}]}}'

- response: {"tx_0123"} //transaction hash

#### Add contracts
- request:

curl --location --request GET 'http://localhost:8804' \
&emsp;    --header 'Token: SomeToken' \
&emsp;    --header 'Content-Type: application/json' \
&emsp;    --data-raw '{"method": "add", "data": {"contracts": [{"name":"$name", "server": "$server", "address":"$address","abi":"$abi", "private": "$private", "gas_limit":100}]}}'

- response: {"ok"}

#### Delete contracts
- request:

curl --location --request GET 'http://localhost:8804' \
&emsp;    --header 'Token: SomeToken' \
&emsp;    --header 'Content-Type: application/json' \
&emsp;    --data-raw '{"method": "delete", "data": {"names": ["$name"]}}'

- response: {"ok"}

## saiAuth
### Run in Docker
`make up`

### Run as standalone application
`microservices/saiAuth/build/sai-auth` 

### API
#### Register
- request:

 curl --location --request GET 'http://localhost:8800/register' \
 &emsp;    --header 'Token: SomeToken' \
 &emsp;    --header 'Content-Type: application/json' \
 &emsp;    --data-raw '{"key":"user","password":"12345"}'`

- response: '{\"Status\":\"Ok\"}'

#### Login
- request:

'curl --location --request GET 'http://localhost:8800/login' \
&emsp;    --header 'Token: SomeToken' \
&emsp;    --header 'Content-Type: application/json' \
&emsp;    --data-raw '{"key":"user","password":"12345"}''

- response:  '{"token":"3rwef2wef2ff23g2g","User":{"_id":"df22f23r435d","key":"user","roles":["User"]}}'

#### Access 
- request:

'curl --location --request GET 'http://localhost:8800/access' \
&emsp;    --header 'Token: 7ead9e6a0977a3bd33ffec382de1558c1ec139bf704ae19cc853094391afd145' \
&emsp;    --header 'Content-Type: application/json' \
&emsp;    --data-raw '{"collection":"users", "method": "get" }''
&emsp;    - response 

- response: 'true'
