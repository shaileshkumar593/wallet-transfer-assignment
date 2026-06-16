# create-project.ps1

mkdir wallet-transfer

mkdir wallet-transfer\cmd\server
mkdir wallet-transfer\internal\handler
mkdir wallet-transfer\internal\service
mkdir wallet-transfer\internal\repository
mkdir wallet-transfer\internal\domain
mkdir wallet-transfer\internal\database
mkdir wallet-transfer\internal\dto
mkdir wallet-transfer\migrations
mkdir wallet-transfer\tests

ni wallet-transfer\cmd\server\main.go -ItemType File

ni wallet-transfer\internal\handler\transfer_handler.go -ItemType File
ni wallet-transfer\internal\service\transfer_service.go -ItemType File

ni wallet-transfer\internal\repository\wallet_repository.go -ItemType File
ni wallet-transfer\internal\repository\transfer_repository.go -ItemType File
ni wallet-transfer\internal\repository\ledger_repository.go -ItemType File
ni wallet-transfer\internal\repository\idempotency_repository.go -ItemType File

ni wallet-transfer\internal\domain\wallet.go -ItemType File
ni wallet-transfer\internal\domain\transfer.go -ItemType File
ni wallet-transfer\internal\domain\ledger.go -ItemType File
ni wallet-transfer\internal\domain\errors.go -ItemType File

ni wallet-transfer\internal\database\postgres.go -ItemType File

ni wallet-transfer\internal\dto\request.go -ItemType File
ni wallet-transfer\internal\dto\response.go -ItemType File

ni wallet-transfer\migrations\001_init.sql -ItemType File

ni wallet-transfer\tests\transfer_test.go -ItemType File
ni wallet-transfer\tests\idempotency_test.go -ItemType File
ni wallet-transfer\tests\concurrency_test.go -ItemType File

ni wallet-transfer\docker-compose.yml -ItemType File
ni wallet-transfer\Dockerfile -ItemType File
ni wallet-transfer\README.md -ItemType File:wq
