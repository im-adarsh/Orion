package orion

import (
	"github.com/carousell/Orion/orion/handlers"
)

//RegisterEncoder allows for registering an HTTP request encoder to arbitrary urls
//Note: this is normally called from protoc-gen-orion autogenerated files
func RegisterEncoder(svr Server, serviceName, method, httpMethod, path string, encoder handlers.Encoder) {
	if e, ok := svr.(handlers.Encodeable); ok {
		e.AddEncoder(serviceName, method, httpMethod, path, encoder)
	}
}
