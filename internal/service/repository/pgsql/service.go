package pgsql

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"go.uber.org/zap"

	"github.com/gradusp/crispy/internal/model"
)

type ServiceRepo struct {
	log  *zap.SugaredLogger
	pool *pgxpool.Pool
}

func NewServiceRepo(pool *pgxpool.Pool, l *zap.SugaredLogger) *ServiceRepo {
	return &ServiceRepo{
		log:  l,
		pool: pool,
	}
}

func (sr *ServiceRepo) Create(ctx context.Context, cl *model.Cluster, s *model.Service) (*model.Service, error) {
	c, err := sr.pool.Acquire(ctx)
	if err != nil {
		sr.log.Error(err)
		return nil, err
	}
	defer c.Release()

	query := `
insert into controller.services (cluster_id, routing_type, balancing_type, bandwidth, proto, addr, port)
values ($1, $2, $3, $4, $5, $6, $7)
returning id;`

	if err := c.QueryRow(ctx, query, cl.ID,
		s.RoutingType, s.BalancingType, s.Bandwidth, s.Proto, s.Addr.To4(), s.Port).Scan(&s.ID); err != nil {
		sr.log.Error(err)
		return nil, err
	}

	return s, nil
}

// TODO: implement
func (sr *ServiceRepo) Get(ctx context.Context) ([]*model.Service, error) {
	c, err := sr.pool.Acquire(ctx)
	if err != nil {
		sr.log.Error(err)
		return nil, err
	}
	defer c.Release()

	query := `select id, cluster_id, routing_type, balancing_type, bandwidth, proto, addr, port from controller.services`
	services, err := c.Query(ctx, query)
	if err != nil {
		sr.log.Error(err)
		return nil, err
	}

	var r []*model.Service
	for services.Next() {
		var s model.Service
		err = services.Scan(&s.ID, &s.ClusterID, &s.RoutingType, &s.BalancingType, &s.Bandwidth, &s.Proto, &s.Addr, &s.Port)
		if err != nil {
			sr.log.Error(err)
			return nil, err
		}
		r = append(r, &s)
	}
	err = services.Err()

	return r, err
}

func (sr *ServiceRepo) GetByID(ctx context.Context, s *model.Service) (*model.Service, error) {
	c, err := sr.pool.Acquire(ctx)
	if err != nil {
		sr.log.Error(err)
		return nil, err
	}
	defer c.Release()

	clusterRow := c.QueryRow(ctx,
		"select cluster_id, routing_type, balancing_type, bandwidth, proto, addr, port from controller.services where id=$1;",
		s.ID)
	if err = clusterRow.Scan(&s.ClusterID, &s.RoutingType, &s.BalancingType, &s.Bandwidth, &s.Proto, &s.Addr, &s.Port); err != nil {
		sr.log.Error(err)
		return nil, err
	}

	var (
		reals        []*model.Real
		healthchecks []*model.Healthcheck
	)

	rq, err := c.Query(ctx, "select id, addr, port, hc_addr, hc_port from controller.reals where service_id=$1", s.ID)
	if err != nil {
		sr.log.Error("reals can't be selected", err)
		return nil, err
	}

	for rq.Next() {
		var r model.Real
		err = rq.Scan(&r.ID, &r.Addr, &r.Port, &r.HealthcheckAddr, &r.HealthcheckPort)
		if err != nil {
			sr.log.Error(err)
			return nil, err
		}
		reals = append(reals, &r)
	}
	//err = rq.Err() // TODO

	// FIXME: two queries can race for db conn
	hcq, err := c.Query(ctx,
		"select id, hello_timer, response_timer, alive_threshold, dead_threshold, quorum, hysteresis from controller.healthchecks where service_id=$1",
		s.ID)
	if err != nil {
		sr.log.Error("healthchecks can't be selected: ", err)
		return nil, err
	}

	for hcq.Next() {
		var hc model.Healthcheck
		err = hcq.Scan(&hc.ID, &hc.HelloTimer, &hc.ResponseTimer, &hc.AliveThreshold, &hc.DeadThreshold, &hc.Quorum, &hc.Hysteresis)
		if err != nil {
			sr.log.Error(err)
			return nil, err
		}
		healthchecks = append(healthchecks, &hc)
	}
	//err = hcq.Err() // TODO

	//if err = c.QueryRow(ctx, "select id, name from controller.clusters where id=$1", s.ClusterID).Scan(&s.Cluster.ID, &s.Cluster.Name); err != nil {
	//	sr.log.Error(err)
	//	return nil, err
	//}

	s.Reals = reals
	s.Healthchecks = healthchecks

	return s, nil
}

// TODO: implement
func (sr *ServiceRepo) Update(ctx context.Context) error {
	panic("implement my repo")
}

func (sr *ServiceRepo) Delete(ctx context.Context, s *model.Service) error {
	c, err := sr.pool.Acquire(ctx)
	if err != nil {
		sr.log.Error(err)
		return err
	}
	defer c.Release()

	r, err := c.Exec(ctx, "delete from controller.services where id=$1", s.ID)
	if err != nil {
		sr.log.Error(err)
		return err
	}
	if r.RowsAffected() != 1 {
		sr.log.Debug("delete for non-existing Service ID")
		return nil
	}

	return nil
}
