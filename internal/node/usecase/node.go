package usecase

import (
	"context"
	"fmt"
	"net"

	"github.com/gradusp/crispy/internal/model"
	"github.com/gradusp/crispy/internal/node"
)

type Usecase struct {
	r node.Repository
}

func NewUsecase(r node.Repository) *Usecase {
	return &Usecase{
		r: r,
	}
}

func (u Usecase) Create(ctx context.Context, cid, h string, a net.IP) (*model.Node, error) {
	n := &model.Node{
		ClusterID: cid,
		Addr:      a,
		Hostname:  h,
	}

	return u.r.Create(ctx, n)
}

// TODO: implement checks for SQL INJECTIONS (regexp?)
// TODO: check 'cid' for UUID
// TODO: check 'a' for net.IP
func (u Usecase) Get(ctx context.Context, cid, a string) ([]*model.Node, error) {
	switch {
	case cid != "" && a != "":
		return nil, node.ErrWrongQuery
	case cid != "":
		q := fmt.Sprintf("where cluster_id='%s'", cid)
		return u.r.GetByField(ctx, q)
	case a != "":
		q := fmt.Sprintf("where addr='%s'", net.ParseIP(a))
		return u.r.GetByField(ctx, q)
	default:
		return u.r.GetByField(ctx, "")
	}
}

func (u Usecase) GetByID(ctx context.Context, nid int) (*model.Node, error) {
	n := &model.Node{ID: nid}
	return u.r.GetByID(ctx, n)
}

func (u Usecase) Delete(ctx context.Context, nid int) error {
	n := &model.Node{ID: nid}
	return u.r.Delete(ctx, n)
}
