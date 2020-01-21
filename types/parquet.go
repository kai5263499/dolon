package types

type ProtoRecord struct {
	Data string `parquet:"name=data, type=BYTE_ARRAY"`
}
