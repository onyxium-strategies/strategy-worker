package main

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)
	db, err := models.InitDB("localhost")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	env.db = db
	os.Exit(m.Run())
}

func TestWalk(t *testing.T) {
	WorkerQueue = make(chan chan WorkRequest, 1)
	testCases := []struct {
		w Worker
		t *Tree
	}{
		{
			w: NewWorker(0, WorkerQueue),
			t: &Tree{
				Conditions: []Condition{
					{ConditionType: "greater-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.02},
				},
				Action: Action{OrderType: "limit-buy", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 10, Value: 0.08},
				Left:   nil,
				Right:  nil,
			},
		},
		{
			w: NewWorker(0, WorkerQueue),
			t: &Tree{
				Conditions: []Condition{
					{ConditionType: "less-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.14},
				},
				Action: Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 10, Value: 0.088},
				Left:   nil,
				Right:  nil,
			},
		},
		{
			w: NewWorker(0, WorkerQueue),
			t: &Tree{
				Conditions: []Condition{
					{ConditionType: "greater-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.02},
					{ConditionType: "percentage-increase", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.001, TimeframeInMS: 3600000},
				},
				Action: Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 10, Value: 0.079},
				Left:   nil,
				Right: &Tree{
					Conditions: []Condition{
						{ConditionType: "greater-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.02},
						{ConditionType: "percentage-decrease", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.001, TimeframeInMS: 3600000},
					},
					Action: Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 10, Value: 0.079},
				},
			},
		},
	}

	for _, c := range testCases {
		c.w.Walk(c.t, c.t)
	}
}
