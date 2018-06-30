package pagerduty

type ExtensionSchema struct {
	APIObject
}

type Extension struct {
	APIObject
	EndpointUrl      string          `json:"endpoint_url"`
	Name             string          `json:"name"`
	ExtensionSchema  ExtensionSchema `json:"extension_schema"`
	ExtensionObjects []APIObject     `json:"extension_objects"`
}

func NewExtension(opts ...ExtensionOptFunc) *Extension {
	ext := &Extension{
		APIObject: APIObject{
			Type: ExtensionResourceType,
		},
		EndpointUrl: "PLACEHOLDER",
		Name:        "PLACEHOLDER",
		ExtensionSchema: ExtensionSchema{
			APIObject: APIObject{
				Type: "PLACEHOLDER",
				ID:   "PLACEHOLDER",
			},
		},
		ExtensionObjects: []APIObject{},
	}
	for _, opt := range opts {
		opt(ext)
	}
	return ext
}

type ExtensionOptFunc func(*Extension)

func ExtensionWithService(serviceId string) ExtensionOptFunc {
	return func(extension *Extension) {
		svcRef := APIObject{
			ID:   serviceId,
			Type: ServiceResourceType,
		}
		extension.ExtensionObjects = append(extension.ExtensionObjects, svcRef)
	}
}

func ExtensionWithName(name string) ExtensionOptFunc {
	return func(extension *Extension) {
		extension.Name = name
	}
}

func ExtensionWithEndpoint(endpoint string) ExtensionOptFunc {
	return func(extension *Extension) {
		extension.EndpointUrl = endpoint
	}
}

func ExtensionWithSchema(r ExtensionSchema) ExtensionOptFunc{
	return func(extension *Extension) {
		extension.ExtensionSchema = r
	}
}

func ExtensionWithObjects(r ...APIObject) ExtensionOptFunc{
	return func(extension *Extension) {
		extension.ExtensionObjects = r
	}
}


func (c *Client) GetExtension(id string, opts ...ResourceRequestOptionFunc) (*Extension, error) {
	res, err := c.GetResource(ExtensionResourceType, id, opts...)
	if err != nil {
		return nil, err
	}
	obj := res.(Extension)
	return &obj, nil
}
