.PHONY: server
server:
	( cd ../lago && source .env && docker-compose -f docker-compose.arm64.yml up )
