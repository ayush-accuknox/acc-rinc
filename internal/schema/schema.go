package schema

import (
	"fmt"

	"github.com/accuknox/rinc/internal/db"
	"github.com/accuknox/rinc/types/ceph"
	"github.com/accuknox/rinc/types/connectivity"
	"github.com/accuknox/rinc/types/dass"
	"github.com/accuknox/rinc/types/imagetag"
	"github.com/accuknox/rinc/types/longjobs"
	"github.com/accuknox/rinc/types/pod"
	"github.com/accuknox/rinc/types/pv"
	"github.com/accuknox/rinc/types/rabbitmq"
	"github.com/accuknox/rinc/types/resource"

	"github.com/invopop/jsonschema"
)

// Generate generates json schema from go structs.
func Generate(target string) ([]byte, error) {
	r := new(jsonschema.Reflector)
	r.FieldNameTag = "-"

	var schema *jsonschema.Schema

	switch target {
	case db.CollectionRabbitmq:
		schema = r.Reflect(rabbitmq.Metrics{})
	case db.CollectionCeph:
		schema = r.Reflect(ceph.Metrics{})
	case db.CollectionImageTag:
		schema = r.Reflect(imagetag.Metrics{})
	case db.CollectionDass:
		schema = r.Reflect(dass.Metrics{})
	case db.CollectionLongJobs:
		schema = r.Reflect(longjobs.Metrics{})
	case db.CollectionPVUtilizaton:
		schema = r.Reflect(pv.Metrics{})
	case db.CollectionResourceUtilization:
		schema = r.Reflect(resource.Metrics{})
	case db.CollectionConnectivity:
		schema = r.Reflect(connectivity.Metrics{})
	case db.CollectionPodStatus:
		schema = r.Reflect(pod.Metrics{})
	default:
		return nil, fmt.Errorf("invalid target: %q", target)
	}

	out, err := schema.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("marshalling schema to json: %w", err)
	}
	return out, nil
}
