test:
	TEST_ZONE_NAME=lihuaio.com go test .
build:
	docker build -t tttlkkkl/cert-manage-alidns .
	docker push tttlkkkl/cert-manage-alidns