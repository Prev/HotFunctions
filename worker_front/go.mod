module github.com/Prev/HotFunctions/worker_front

go 1.13

require (
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/Prev/HotFunctions/worker_front/types v1.0.0
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.4.2-0.20190717161051-705d9623b7c1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/frankban/quicktest v1.14.4 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nwaples/rardecode v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pierrec/lz4 v2.4.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sevlyar/go-daemon v0.1.5
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/grpc v1.26.0 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)

replace github.com/Prev/HotFunctions/worker_front/types => ./types

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190717161051-705d9623b7c1
