package pagerduty

import (
	"net/http"
)

type ResourceRequestOptionFunc func(*http.Request) error

func setQueryParam(key, value string, request *http.Request) error {
	params := request.URL.Query()
	params.Add(key, value)
	qString := params.Encode()
	request.URL.RawQuery = qString
	return nil
}

func WithParams(p map[string]string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		for key, value := range p {
			params.Add(key, value)
		}
		return nil
	}
}

func WithDateRange(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		params.Add("date_range", value)
		return nil
	}
}

func WithEditable(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		params.Add("editable", value)
		return nil
	}
}

func WithEscalationPolicyIDs(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		params.Add("escalation_policy_ids", value)
		return nil
	}
}

func WithFilter(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		params.Add("filter", value)
		return nil
	}
}

func WithIncludes(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		params.Add("includes", value)
		return nil
	}
}

func WithIsOverview(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		params.Add("is_overview", value)
		return nil
	}
}

func WithOverflow(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		params.Add("overflow", value)
		return nil
	}
}

func WithQuery(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("query", value, request)
	}
}

func WithServiceIDs(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		params.Add("service_ids", value)
		return nil
	}
}

func WithSince(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		params.Add("since", value)
		return nil
	}
}

func WithSortBy(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		params.Add("sort_by", value)
		return nil
	}
}

func WithStatuses(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		params.Add("statuses", value)
		return nil
	}
}

func WithTeamIDs(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		params.Add("team_ids", value)
		return nil
	}
}

func WithTimeZone(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		params.Add("time_zone", value)
		return nil
	}
}

func WithUntil(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		params.Add("until", value)
		return nil
	}
}

func WithUserIDs(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		params.Add("user_ids", value)
		return nil
	}
}

func WithHeader(key, value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		request.Header.Set(key, value)
		return nil
	}
}

func WithHeaders(headers map[string]string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		if headers != nil {
			for k, v := range headers {
				request.Header.Set(k, v)
			}
		}
		return nil
	}
}
