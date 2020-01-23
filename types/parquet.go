package types

// ProtoRecord is a simple wrapper for seralized protobufs
type ProtoRecord struct {
	Data string `parquet:"name=data, type=BYTE_ARRAY"`
}
