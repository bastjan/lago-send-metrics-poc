package prommock

import (
	"context"
	"fmt"
	"time"

	"github.com/appuio/appuio-cloud-reporting/pkg/sourcekey"
	apiv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type FakeQuerier struct {
	Queries map[string]map[string]model.SampleValue
}

func (q FakeQuerier) Query(ctx context.Context, query string, ts time.Time, _ ...apiv1.Option) (model.Value, apiv1.Warnings, error) {
	var res model.Vector
	for k, s := range q.Queries[query] {
		sk, err := sourcekey.Parse(k)
		if err != nil {
			return nil, nil, err
		}
		res = append(res, &model.Sample{
			Metric: map[model.LabelName]model.LabelValue{
				"product":  model.LabelValue(k),
				"category": model.LabelValue(fmt.Sprintf("%s:%s", sk.Zone, sk.Namespace)),
				"tenant":   model.LabelValue(sk.Tenant),
			},
			Value:     s,
			Timestamp: model.TimeFromUnixNano(ts.UnixNano()),
		})
	}
	return res, nil, nil
}
