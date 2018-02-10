.PHONY: gen_protos
gen_protos:
	protoc --proto_path=definitions --go_out=plugins=grpc:generated definitions/*.proto