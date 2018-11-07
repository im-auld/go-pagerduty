package pagerduty

import (
	"net/http"
)

type ResourceRequestOptionFunc func(*http.Request) error

func setQueryParam(key, value string, request *http.Request) error {
	params := request.URL.Query()
	params.Add(key, value)
	request.URL.RawQuery = params.Encode()
	return nil
}

func WithParams(p map[string]string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		params := request.URL.Query()
		for key, value := range p {
			params.Add(key, value)
		}
		request.URL.RawQuery = params.Encode()
		return nil
	}
}

func WithDateRange(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("date_range", value, request)
	}
}

func WithEditable(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("editable", value, request)
	}
}

func WithEscalationPolicyIDs(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("escalation_policy_ids", value, request)
	}
}

func WithFilter(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("filter", value, request)
	}
}

func WithIncludes(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("includes", value, request)
	}
}

func WithIsOverview(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("is_overview", value, request)
	}
}

func WithOverflow(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("overflow", value, request)
	}
}

func WithQuery(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("query", value, request)
	}
}

func WithServiceIDs(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("service_ids", value, request)
	}
}

func WithSince(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("since", value, request)
	}
}

func WithSortBy(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("sort_by", value, request)
	}
}

func WithStatuses(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("statuses", value, request)
	}
}

func WithTeamIDs(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("team_ids", value, request)
	}
}

func WithTimeZone(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("time_zone", value, request)
	}
}

func WithUntil(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("until", value, request)
	}
}

func WithUserIDs(value string) ResourceRequestOptionFunc {
	return func(request *http.Request) error {
		return setQueryParam("user_ids", value, request)
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
