package pdp

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/open-policy-agent/opa/rego"
)

type PolicyEngine struct {
	lock       sync.Mutex
	repository string
	rego       *rego.Rego
	query      *rego.PreparedEvalQuery
}

const (
	policyQuery = "x = data.pdp"
)

func NewPolicyEngine(path string, changeMonitor bool) (*PolicyEngine, error) {
	svc := PolicyEngine{repository: path}
	err := svc.Load(changeMonitor)
	if err != nil {
		return &PolicyEngine{}, err
	}

	return &svc, nil
}

func (svc *PolicyEngine) Evaluate(input PolicyInput) (PolicyResponse, error) {
	return PolicyResponse{Allow: true}, nil
}

func (svc *PolicyEngine) Load(changeMonitor bool) error {
	err := svc.loadPolicy()
	if err != nil {
		return err
	}

	// TODO: Switch to inotify/kqueue
	if changeMonitor {
		d, err := time.ParseDuration(policyEvalChangeMonitorInterval)
		if err != nil {
			return err
		}

		go func() {
			time.Sleep(d)

			err := svc.loadPolicy()
			if err != nil {
				log.Printf("Failed to reload policy: %s", err.Error())
			}
		}()
	}

	return nil
}

func (svc *PolicyEngine) loadPolicy() error {
	queryFn := rego.Query(policyQuery)
	policyDoc := rego.Load([]string{svc.repository}, nil)

	r := rego.New(queryFn, policyDoc)
	q, err := r.PrepareForEval(context.Background())

	if err != nil {
		return err
	}

	svc.lock.Lock()
	defer svc.lock.Unlock()

	svc.rego = r
	svc.query = &q

	return nil
}
