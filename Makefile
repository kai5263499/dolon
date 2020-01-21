gen-protos:
	protoc --proto_path=protos --go_out=types protos/*.proto