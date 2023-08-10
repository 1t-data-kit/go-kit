package network

import (
	"context"
	"fmt"
	"github.com/1t-data-kit/go-kit/base"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	etcdClient "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"time"
)

type Registrar struct {
	etcd  *etcdClient.Client
	lease etcdClient.Lease
	ttl   int64 //seconds
}

var defaultTTL int64 = 5

func NewRegistrar(etcd *etcdClient.Client) *Registrar {
	return &Registrar{
		etcd:  etcd,
		lease: etcdClient.NewLease(etcd),
		ttl:   defaultTTL,
	}
}

func NewRegistrarOption(registrar *Registrar) base.Option {
	return base.NewOption(registrar)
}

func (r *Registrar) SetTTl(seconds int64) *Registrar {
	r.ttl = seconds
	return r
}

func (r *Registrar) getTTL() int64 {
	ttl := r.ttl
	if ttl < 5 {
		ttl = 5
	}

	return ttl
}

func (r *Registrar) getTTLSeconds() time.Duration {
	return time.Duration(r.getTTL()) * time.Second
}

func (r *Registrar) getHalfTTLSeconds() time.Duration {
	timeout := r.ttl / 2
	if timeout < 2 {
		timeout = 2
	}

	return time.Duration(timeout) * time.Second
}

func (r *Registrar) Register(ctx context.Context, endpoints ...Endpoint) {
	for _, endpoint := range endpoints {
		_endpoint := endpoint
		if !_endpoint.MustRegisterNetwork() {
			logrus.Infof("registry.network.Registrar.Register ignore: %s MustRegisterNetwork() has be returned false", _endpoint.Name())
			continue
		}

		go func() {
			var leaseID etcdClient.LeaseID
			var err error
			var retry int
			interval := r.getHalfTTLSeconds()
			for {
				r.waitEndpointStart(_endpoint)
				leaseID, err = r.register(ctx, _endpoint.Name(), _endpoint.Address(), leaseID)
				if err != nil {
					retry++
					logrus.Errorf("registry.network.Registrar.register error: %s[%s]%s ... retry[%d]", _endpoint.Name(), _endpoint.Address(), errors.WithStack(err), retry)
				}
				retry = 0
				time.Sleep(interval)
			}
		}()
	}

}

func (r *Registrar) waitEndpointStart(_endpoint Endpoint) {
	var retry int
	interval := r.getHalfTTLSeconds()
	for {
		retry++
		if _endpoint.Address() != "" {
			break
		}
		logrus.Debugf("registry.network.Registrar.waitEndpointStart: %s[%s] ... retry[%d]", _endpoint.Name(), _endpoint.Address(), retry)
		time.Sleep(interval)
	}
}

func (r *Registrar) register(ctx context.Context, name, address string, leaseID etcdClient.LeaseID) (etcdClient.LeaseID, error) {
	if leaseID == 0 {
		return r.add(ctx, name, address)
	}
	return r.refresh(ctx, leaseID)
}

func (r *Registrar) add(ctx context.Context, name, address string) (etcdClient.LeaseID, error) {
	ctx, cancel := context.WithTimeout(ctx, r.getHalfTTLSeconds())
	defer cancel()

	leaseResp, err := r.lease.Grant(ctx, r.getTTL())
	if err != nil {
		return 0, err
	}
	leaseID := leaseResp.ID

	manager, err := endpoints.NewManager(r.etcd, name)
	if err != nil {
		return 0, err
	}
	if err = manager.AddEndpoint(
		ctx,
		name+"/"+address,
		endpoints.Endpoint{Addr: address},
		etcdClient.WithLease(leaseID),
	); err != nil {
		return 0, err
	}
	logrus.Debugf("registry.network.Registrar.register success: %s[%s][%v]", name, address, leaseID)

	return leaseID, nil
}

func (r *Registrar) refresh(ctx context.Context, leaseID etcdClient.LeaseID) (etcdClient.LeaseID, error) {
	ctx, cancel := context.WithTimeout(ctx, r.getHalfTTLSeconds())
	defer cancel()

	resp, err := r.lease.KeepAliveOnce(ctx, leaseID)
	if err != nil {
		return leaseID, err
	}

	if resp.ID != leaseID {
		return leaseID, fmt.Errorf("registry.network.Registrar.refresh error: leaseID[%v/%v] has be changed", leaseID, resp.ID)
	}
	logrus.Debugf("registry.network.Registrar.refresh success: [%v]", leaseID)

	return leaseID, nil
}
