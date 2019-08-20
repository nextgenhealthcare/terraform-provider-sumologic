package sumologic

import (
	"github.com/nextgenhealthcare/sumologic-sdk-go"
)

func appendFilters(sumologicFilters []interface{}) []sumologic.Filter {
	filters := make([]sumologic.Filter, 0)
	for _, v := range sumologicFilters {
		v := v.(map[string]interface{})

		f := sumologic.Filter{
			FilterType: v["filter_type"].(string),
			Name:       v["name"].(string),
			Regexp:     v["regexp"].(string),
		}

		filters = append(filters, f)
	}

	return filters
}
