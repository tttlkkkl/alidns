tag?=latest

test:
	TEST_ZONE_NAME=lihuaio.com go test .
build:
	docker build -t tttlkkkl/cert-manage-alidns:$(tag) .
	docker push tttlkkkl/cert-manage-alidns:$(tag)