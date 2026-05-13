package helper

import (
	"github.com/Compogo/resource/domain"
	"github.com/Compogo/resource/protobuf/resource"
	"github.com/Compogo/types/linker"
)

var (
	TypeLinkProtoType = linker.NewLinker[domain.Type, resource.Type](
		linker.Link(domain.Memory, resource.Type_Memory),
		linker.Link(domain.CPU, resource.Type_CPU),
	)

	ProtoTypeLinkType = linker.NewLinker[resource.Type, domain.Type](
		linker.Link(resource.Type_Memory, domain.Memory),
		linker.Link(resource.Type_CPU, domain.CPU),
	)

	StateLinkProtoState = linker.NewLinker[domain.State, resource.State](
		linker.Link(domain.Normal, resource.State_Normal),
		linker.Link(domain.Warning, resource.State_Warning),
		linker.Link(domain.Alarm, resource.State_Alarm),
	)

	ProtoStateLinkState = linker.NewLinker[resource.State, domain.State](
		linker.Link(resource.State_Normal, domain.Normal),
		linker.Link(resource.State_Warning, domain.Warning),
		linker.Link(resource.State_Alarm, domain.Alarm),
	)
)
