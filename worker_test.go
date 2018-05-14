package main

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	"testing"
	"time"
)

// TODO assert taken path is expected path
func TestWalkDeadlock(t *testing.T) {
	WorkerQueue = make(chan chan *models.Strategy, 1)
	testCases := []struct {
		name         string
		w            Worker
		t            *models.Tree
		expectedPath []int //Tree.Id
	}{
		{
			name: "greater-than-or-equal-to",
			w:    NewWorker(0, WorkerQueue),
			t: &models.Tree{
				Id: 0,
				Conditions: []models.Condition{
					{ConditionType: "greater-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.072},
					{ConditionType: "greater-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-bid", Value: 0.066},
				},
				Action: models.Action{OrderType: "limit-buy", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 1, Value: 0.08},
				Left:   nil,
				Right: &models.Tree{
					Id: 1,
					Conditions: []models.Condition{
						{ConditionType: "greater-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-ask", Value: 0.075},
					},
					Action: models.Action{OrderType: "limit-buy", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 2, Value: 0.08},
					Left:   nil,
					Right: &models.Tree{
						Id: 3,
						Conditions: []models.Condition{
							{ConditionType: "greater-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "volume", Value: 8000},
						},
						Action: models.Action{OrderType: "limit-buy", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 3, Value: 0.08},
						Left: &models.Tree{
							Id: 2,
							Conditions: []models.Condition{
								{ConditionType: "greater-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.065},
							},
							Action: models.Action{OrderType: "limit-buy", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 4, Value: 0.08},
							Left:   nil,
							Right:  nil,
						},
						Right: nil,
					},
				},
			},
			expectedPath: []int{3, 2},
		},
		{
			name: "less-than-or-equal-to",
			w:    NewWorker(0, WorkerQueue),
			t: &models.Tree{
				Id: 0,
				Conditions: []models.Condition{
					{ConditionType: "less-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "volume", Value: 9000},
				},
				Action: models.Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 1, Value: 0.088},
				Left:   nil,
				Right: &models.Tree{
					Id: 1,
					Conditions: []models.Condition{
						{ConditionType: "less-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-bid", Value: 0.069},
						{ConditionType: "less-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-ask", Value: 0.07},
					},
					Action: models.Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 2, Value: 0.088},
					Left:   nil,
					Right: &models.Tree{
						Id: 4,
						Conditions: []models.Condition{
							{ConditionType: "less-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.071},
						},
						Action: models.Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 3, Value: 0.088},
						Left: &models.Tree{
							Id: 3,
							Conditions: []models.Condition{
								{ConditionType: "less-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-bid", Value: 0.072},
							},
							Action: models.Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 4, Value: 0.088},
							Left: &models.Tree{
								Id: 2,
								Conditions: []models.Condition{
									{ConditionType: "less-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-ask", Value: 0.075},
								},
								Action: models.Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 5, Value: 0.088},
								Left:   nil,
								Right:  nil,
							},
							Right: nil,
						},
						Right: nil,
					},
				},
			},
			expectedPath: []int{4, 3, 2},
		},
		{
			name: "percentage-increase",
			w:    NewWorker(0, WorkerQueue),
			t: &models.Tree{
				Id: 0,
				Conditions: []models.Condition{
					{ConditionType: "percentage-increase", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.1, TimeframeInMS: 1},
				},
				Action: models.Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 1, Value: 0.079},
				Left:   nil,
				Right: &models.Tree{
					Id: 2,
					Conditions: []models.Condition{
						{ConditionType: "percentage-increase", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-bid", Value: 0.07, TimeframeInMS: 1},
					},
					Action: models.Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 2, Value: 0.079},
					Left: &models.Tree{
						Id: 1,
						Conditions: []models.Condition{
							{ConditionType: "percentage-increase", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-ask", Value: 0.065, TimeframeInMS: 1},
							{ConditionType: "percentage-increase", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "volume", Value: 0.3, TimeframeInMS: 1},
						},
						Action: models.Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 3, Value: 0.079},
						Left:   nil,
						Right:  nil,
					},
					Right: &models.Tree{
						Id: 3,
						Conditions: []models.Condition{
							{ConditionType: "percentage-increase", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.07, TimeframeInMS: 1},
						},
						Action: models.Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 4, Value: 0.079},
						Left:   nil,
						Right:  nil,
					},
				},
			},
			expectedPath: []int{2, 1},
		},
		{
			name: "percentage-decrease",
			w:    NewWorker(0, WorkerQueue),
			t: &models.Tree{
				Id: 3,
				Conditions: []models.Condition{
					{ConditionType: "percentage-decrease", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.045, TimeframeInMS: 0},
				},
				Action: models.Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 1, Value: 0.079},
				Left: &models.Tree{
					Id: 2,
					Conditions: []models.Condition{
						{ConditionType: "percentage-decrease", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-bid", Value: 0.065, TimeframeInMS: 0},
					},
					Action: models.Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 2, Value: 0.079},
					Left: &models.Tree{
						Id: 1,
						Conditions: []models.Condition{
							{ConditionType: "percentage-decrease", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-ask", Value: 0.058, TimeframeInMS: 0},
						},
						Action: models.Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 3, Value: 0.079},
						Left: &models.Tree{
							Id: 0,
							Conditions: []models.Condition{
								{ConditionType: "percentage-decrease", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "volume", Value: 0.1, TimeframeInMS: 0},
							},
							Action: models.Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 4, Value: 0.079},
							Left:   nil,
							Right:  nil,
						},
						Right: nil,
					},
					Right: nil,
				},
				Right: nil,
			},
			expectedPath: []int{3, 2, 1, 0},
		},
	}

	done := make(chan bool)
	timeout := make(chan bool)
	timer := time.NewTimer(time.Second * 5)
	go func() {
		<-timer.C
		timeout <- true
	}()
	go func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.w.Walk(tc.t, tc.t)
			})
		}
		done <- true
	}(t)
	select {
	case <-timeout:
		t.Fatal("The Walk function has reached a deadlock. Timeout reached.")
	case <-done:
	}
}
