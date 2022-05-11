package pdp

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
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

func (svc *PolicyEngine) Evaluate(ctx context.Context, input PolicyInput) (PolicyResponse, error) {
	svc.lock.Lock()
	defer svc.lock.Unlock()

	// log.Printf("PolicyInput: %s", utils.Introspect(input))

	rs, err := svc.query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return PolicyResponse{}, err
	}

	if len(rs) == 0 || rs[0].Bindings["x"] == nil {
		return PolicyResponse{}, errors.New("Policy evaluation returned unexpected result")
	}

	x := rs[0].Bindings["x"]
	var p PolicyResponse
	err = utils.MapStruct(x, &p)
	if err != nil {
		return PolicyResponse{}, err
	}

	// log.Printf("Policy response: %s", utils.Introspect(p))
	return p, nil
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
			log.Printf("Failed to parse ticker duration for policy reload")
			return err
		}

		ticker := time.NewTicker(d)
		tickerStop := make(chan os.Signal)

		signal.Notify(tickerStop, os.Interrupt)
		go func() {
			for {
				select {
				case <-ticker.C:
					log.Printf("Re-loading policy from path: %s", svc.repository)
					err := svc.loadPolicy()
					if err != nil {
						log.Printf("Failed to reload policy: %s", err.Error())
					}
				case <-tickerStop:
					ticker.Stop()
					return
				}
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
