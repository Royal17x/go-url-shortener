$PROTO_DIR = "proto"
$OUT_DIR = "internal/pb"

protoc `
    --proto_path=$PROTO_DIR `
    --go_out=$OUT_DIR --go_opt=paths=source_relative `
    --go-grpc_out=$OUT_DIR --go-grpc_opt=paths=source_relative `
    $PROTO_DIR/*.proto