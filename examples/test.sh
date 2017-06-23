#!/usr/bin/env bash

echo '###### service ####'
go run ./service/get.go


echo '##### bucket ####'

go run ./bucket/delete.go
go run ./bucket/put.go
go run ./bucket/putACL.go
go run ./bucket/putCORS.go
go run ./bucket/putLifecycle.go
go run ./bucket/putTagging.go
go run ./bucket/get.go
go run ./bucket/getACL.go
go run ./bucket/getCORS.go
go run ./bucket/getLifecycle.go
go run ./bucket/getTagging.go
go run ./bucket/getLocation.go
go run ./bucket/head.go
go run ./bucket/listMultipartUploads.go
go run ./bucket/delete.go
go run ./bucket/deleteCORS.go
go run ./bucket/deleteLifecycle.go
go run ./bucket/deleteTagging.go


echo '##### object ####'

go run ./bucket/putCORS.go
go run ./object/put.go
go run ./object/putACL.go
go run ./object/append.go
go run ./object/get.go
go run ./object/head.go
go run ./object/getAnonymous.go
go run ./object/getACL.go
go run ./object/listParts.go
go run ./object/options.go
go run ./object/initiateMultipartUpload.go
go run ./object/uploadPart.go
go run ./object/completeMultipartUpload.go
go run ./object/abortMultipartUpload.go
go run ./object/delete.go
go run ./object/deleteMultiple.go
