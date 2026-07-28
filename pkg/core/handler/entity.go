package handle_swagger

// Nguồn chi tiết
//https://swagger.io/specification/v2/

type SwaggerDoc struct {
	Swagger             string                    `json:"swagger"`
	Info                Info                      `json:"info"`
	Host                string                    `json:"host"`
	BasePath            string                    `json:"basePath"`
	Schemes             []string                  `json:"schemes"`
	Paths               map[string]PathItem       `json:"paths"`
	Definitions         map[string]Schema         `json:"definitions"`
	SecurityDefinitions map[string]SecurityScheme `json:"securityDefinitions,omitempty"`
	Tags                []Tag                     `json:"tags,omitempty"`
}
type SecurityRequirement map[string][]string

type SecurityScheme struct {
	Type             string            `json:"type"`
	Description      string            `json:"description,omitempty"`
	Name             string            `json:"name,omitempty"`
	In               string            `json:"in,omitempty"`
	Flow             string            `json:"flow,omitempty"`
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

type Info struct {
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Version     string   `json:"version"`
	Contact     *Contact `json:"contact,omitempty"`
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type PathItem map[string]Operation

type Operation struct {
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	Tags        []string              `json:"tags,omitempty"`
	Consumes    []string              `json:"consumes,omitempty"`
	Produces    []string              `json:"produces,omitempty"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	Responses   map[string]Response   `json:"responses"`
	Security    []SecurityRequirement `json:"security,omitempty"`
}

type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"` // path, query, header, body
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required"`
	Type        string      `json:"type,omitempty"`
	Schema      *SchemaRef  `json:"schema,omitempty"`
	Example     interface{} `json:"example,omitempty"`
}

type SchemaRef struct {
	Ref string `json:"$ref,omitempty"`
}

type Response struct {
	Description string     `json:"description"`
	Schema      *SchemaRef `json:"schema,omitempty"`
}

type Schema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
	Items      *SchemaRef          `json:"items,omitempty"`
	Ref        string              `json:"$ref,omitempty"`
}

type Property struct {
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
	Example     interface{} `json:"example,omitempty"`
	Format      string      `json:"format,omitempty"`
	Minimum     *int        `json:"minimum,omitempty"`
	Maximum     *int        `json:"maximum,omitempty"`
	MinLength   *int        `json:"minLength,omitempty"`
	MaxLength   *int        `json:"maxLength,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
	Ref         string      `json:"$ref,omitempty"`
	Items       *SchemaRef  `json:"items,omitempty"` // <-- thêm dòng này

}
