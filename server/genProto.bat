@echo off
set PROTO_PATH=.\rental\api
set GO_OUT_PATH=.\rental\api\gen\v1
mkdir %GO_OUT_PATH%

protoc -I=%PROTO_PATH% --go_out=%GO_OUT_PATH% --go_opt=paths=source_relative --go-grpc_out=%GO_OUT_PATH% --go-grpc_opt=paths=source_relative rental.proto
@REM protoc -I=%PROTO_PATH%  --grpc-gateway_out %GO_OUT_PATH% --grpc-gateway_opt paths=source_relative --grpc-gateway_opt grpc_api_configuration=%PROTO_PATH%\rental.yaml rental.proto

set PBTS_BIN_DIR=..\wx\miniprogram\node_modules\.bin
set PBTS_OUT_DIR=..\wx\miniprogram\service\proto_gen\rental
mkdir %PBTS_OUT_DIR%

%PBTS_BIN_DIR%\pbjs.cmd -t static -w es6 %PROTO_PATH%\rental.proto --no-create --no-encode --no-decode --no-verify --no-delimited --force-number -o %PBTS_OUT_DIR%\rental_pb_tmp.js
echo import * as $protobuf from "protobufjs"; > %PBTS_OUT_DIR%\rental_pb.js
type %PBTS_OUT_DIR%\rental_pb_tmp.js >> %PBTS_OUT_DIR%\rental_pb.js
del %PBTS_OUT_DIR%\rental_pb_tmp.js

%PBTS_BIN_DIR%\pbts.cmd -o %PBTS_OUT_DIR%\rental_pb.d.ts %PBTS_OUT_DIR%\rental_pb.js