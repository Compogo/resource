package helper

import (
	"github.com/Compogo/resource/domain"
	"github.com/Compogo/resource/protobuf/resource"
)

// NewFromResource преобразует domain.Resource в protobuf.Resource.
func NewFromResource(r *domain.Resource) (*resource.Resource, error) {
	t, err := TypeLinkProtoType.Get(r.Type)
	if err != nil {
		return nil, err
	}

	return &resource.Resource{
		Type:  t,
		Limit: r.Limit,
		Usage: r.Usage,
	}, nil
}

// NewFromProtoResource преобразует protobuf.Resource в domain.Resource.
func NewFromProtoResource(r *resource.Resource) (*domain.Resource, error) {
	t, err := ProtoTypeLinkType.Get(r.GetType())
	if err != nil {
		return nil, err
	}

	return &domain.Resource{
		Type:  t,
		Limit: r.GetLimit(),
		Usage: r.GetUsage(),
	}, nil
}
