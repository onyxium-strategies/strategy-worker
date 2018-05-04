package main

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	log "github.com/sirupsen/logrus"
	"time"
	// "gopkg.in/mgo.v2/bson"
)

// To get the current value of a metric
func getMetricValue(baseMetric string, market models.Market) float64 {
	var currentValue float64

	switch baseMetric {
	case "price-last":
		currentValue = market.Last
	case "price-ask":
		currentValue = market.Ask
	case "price-bid":
		currentValue = market.Bid
	case "volume":
		currentValue = market.Volume
	default:
		log.Errorf("Condition BaseMetric %s does not exist", baseMetric)
	}

	return currentValue
}

// Walk: https://en.m.wikipedia.org/wiki/Left-child_right-sibling_binary_tree
// Should make this a method of a Worker because now we can't send information about the worker in an error message
func walk(tree *Tree, root *Tree) {
	i := 0
	for tree != nil {

		// get latest market update
		latestMarkets, err := models.GetLatestMarket()
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("CURRENT NODE: %+v", tree)

		// If all conditions true do action then Tree.Left else go to next sibling Tree.Right
		// If order to check if all conditions are true we check if not any is false
		// All true == not any False
		doAction := true

		for _, condition := range tree.Conditions {
			if condition.ConditionType == "geq" || condition.ConditionType == "leq" {
				log.Infof("If the %s on the market %s/%s is %s than %.8f.", condition.BaseMetric, condition.BaseCurrency, condition.QuoteCurrency, condition.ConditionType, condition.Value)
			} else {
				log.Infof("If the %s on the market %s/%s has %s with %.3f percentage within %d minutes.", condition.BaseMetric, condition.BaseCurrency, condition.QuoteCurrency, condition.ConditionType, condition.Value, condition.TimeframeInMS/60000)
			}

			latestMarket := latestMarkets.Market[condition.BaseCurrency+"-"+condition.QuoteCurrency]

			switch condition.ConditionType {
			case "geq":
				currentValue := getMetricValue(condition.BaseMetric, latestMarket)
				log.Debugf("MARKET %s: %.8f", condition.BaseMetric, currentValue)
				if currentValue < condition.Value {
					doAction = false
					log.Debugf("doAction FALSE: Market %s with value %.8f is < than condition value %.8f", condition.BaseMetric, currentValue, condition.Value)
				}
			case "leq":
				currentValue := getMetricValue(condition.BaseMetric, latestMarket)
				log.Debugf("MARKET %s: %.8f", condition.BaseMetric, currentValue)
				if currentValue > condition.Value {
					doAction = false
					log.Debugf("doAction FALSE: Market %s with value %.8f is > than condition value %.8f", condition.BaseMetric, currentValue, condition.Value)
				}
			case "percentage-increase":
				historyMarkets, err := models.GetHistoryMarket(condition.TimeframeInMS)
				if err != nil {
					log.Fatal(err)
				}
				historyMarket := historyMarkets.Market[condition.BaseCurrency+"-"+condition.QuoteCurrency]

				currentValue := getMetricValue(condition.BaseMetric, latestMarket)
				pastValue := getMetricValue(condition.BaseMetric, historyMarket)

				percentage := (currentValue - pastValue) / pastValue
				log.Debugf("MARKET %s changed with %.3f", condition.BaseMetric, percentage)
				if percentage < condition.Value {
					doAction = false
					log.Debugf("doAction FALSE: Market %s with percentage difference of %.3f is < than condition value %.3f", condition.BaseMetric, percentage, condition.Value)
				}

			case "percentage-decrease":
				historyMarkets, err := models.GetHistoryMarket(condition.TimeframeInMS)
				if err != nil {
					log.Fatal(err)
				}
				historyMarket := historyMarkets.Market[condition.BaseCurrency+"-"+condition.QuoteCurrency]

				currentValue := getMetricValue(condition.BaseMetric, latestMarket)
				pastValue := getMetricValue(condition.BaseMetric, historyMarket)

				percentage := (currentValue - pastValue) / pastValue
				log.Debugf("MARKET %s changed with %.3f", condition.BaseMetric, percentage)
				if percentage > -condition.Value {
					doAction = false
					log.Debugf("COMPARISON Market %s with percentage difference of %.3f is > than condition value -%.3f", condition.BaseMetric, percentage, condition.Value)
				}
			default:
				doAction = false
				log.Warningf("Unknown ConditionType %s", condition.ConditionType)
			}

		}

		if doAction {
			switch tree.Action.ValueType {
			case "absolute":
				log.Infof("Set a %s order for %.8f %s at %.8f %s/%s for a total of %.8f %s.", tree.Action.OrderType, tree.Action.Quantity, tree.Action.QuoteCurrency, tree.Action.Value, tree.Action.BaseCurrency, tree.Action.QuoteCurrency, tree.Action.Value*tree.Action.Quantity, tree.Action.BaseCurrency)
			case "relative-above", "relative-below":
				log.Infof("Set a %s order for %.8f %s at the future rate of %s +/- %.8f %s/%s per unit.", tree.Action.OrderType, tree.Action.Quantity, tree.Action.QuoteCurrency, tree.Action.ValueQuoteMetric, tree.Action.Value, tree.Action.BaseCurrency, tree.Action.QuoteCurrency)
			case "percentage-above", "percentage-below":
				log.Infof("Set a %s order for %.8f %s at the future rate of %s * (1 +/- %.8f %s/%s per unit.", tree.Action.OrderType, tree.Action.Quantity, tree.Action.QuoteCurrency, tree.Action.ValueQuoteMetric, tree.Action.Value, tree.Action.BaseCurrency, tree.Action.QuoteCurrency)
			}

			if tree.Left == nil {
				log.Info("NO MORE STATEMENT AFTER THIS ACTION STATEMENT, I'M DONE")
				tree = nil
			} else {
				tree = tree.Left
				root = root.Left
				log.Info("JUMPING to left")
			}
			i = 0
		} else {
			if tree.Right == nil {
				log.Infof("None of the conditions of each child are true. JUMPING to root: %+v", root)
				time.Sleep(3 * time.Second) // Change in production to check each time new market data is available
				tree = root
				i += 1
				log.Infof("%d condition iterations after last action", i)
			} else {
				log.Info("None of the conditions are true. JUMPING to right sibling.")
				tree = tree.Right
			}
		}

	}

}
