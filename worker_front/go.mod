module github.com/Prev/HotFunctions/worker_front

go 1.13

require (
	github.com/Prev/HotFunctions/worker_front/types v1.0.0
	github.com/containerd/containerd v1.3.2 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.4.2-0.20190717161051-705d9623b7c1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/nwaples/rardecode v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pierrec/lz4 v2.4.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sevlyar/go-daemon v0.1.5
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
	golang.org/x/sys v0.0.0-20200113162924-86b910548bc1 // indirect
	google.golang.org/grpc v1.26.0 // indirect
)

replace github.com/Prev/HotFunctions/worker_front/types => ./types

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190717161051-705d9623b7c1
