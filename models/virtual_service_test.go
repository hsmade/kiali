package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/kiali/kiali/models"
)

func TestVirtualServiceHasRequestTimeout(t *testing.T) {
	cases := map[string]struct {
		vsYAML          []byte
		expectedTimeout bool
	}{
		"Has timeout": {
			expectedTimeout: true,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: reviews
spec:
  hosts:
  - reviews
  http:
  - route:
    - destination:
        host: reviews
        subset: v2
    timeout: 0.5s
`),
		},
		"No timeout": {
			expectedTimeout: false,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: reviews
spec:
  hosts:
  - reviews
  http:
  - route:
    - destination:
        host: reviews
        subset: v2
`),
		},
		"Multiple timeouts": {
			expectedTimeout: true,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: reviews
spec:
  hosts:
  - reviews
  http:
  - route:
    - destination:
        host: reviews
        subset: v2
    timeout: 0.5s
  - route:
    - destination:
        host: reviews
        subset: v1
    timeout: 2.5s
`),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			var vs models.VirtualService
			assert.NoError(yaml.Unmarshal(tc.vsYAML, &vs))

			assert.Equal(vs.HasRequestTimeout(), tc.expectedTimeout)
		})
	}

	// Testing nil case
	var vs *models.VirtualService
	assert.False(t, vs.HasRequestTimeout())
}

func TestVirtualServiceHasFaultInjection(t *testing.T) {
	cases := map[string]struct {
		vsYAML                 []byte
		expectedFaultInjection bool
	}{
		"Has fault": {
			expectedFaultInjection: true,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: ratings
spec:
  hosts:
  - ratings
  http:
  - fault:
      delay:
        fixedDelay: 7s
        percentage:
          value: 100
    match:
    - headers:
        end-user:
          exact: jason
    route:
    - destination:
        host: ratings
        subset: v1
  - route:
    - destination:
        host: ratings
        subset: v1
`),
		},
		"No fault": {
			expectedFaultInjection: false,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: ratings
spec:
  hosts:
  - ratings
  http:
  - match:
    - headers:
        end-user:
          exact: jason
    route:
    - destination:
        host: ratings
        subset: v1
  - route:
    - destination:
        host: ratings
        subset: v1
`),
		},
		"Multiple faults": {
			expectedFaultInjection: true,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: ratings
spec:
  hosts:
  - ratings
  http:
  - fault:
      delay:
        fixedDelay: 7s
        percentage:
          value: 100
    match:
    - headers:
        end-user:
          exact: jason
    route:
    - destination:
        host: ratings
        subset: v1
  - route:
    - destination:
        host: ratings
        subset: v1
    fault:
      delay:
        fixedDelay: 5s
        percentage:
          value: 10
`),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			var vs models.VirtualService
			assert.NoError(yaml.Unmarshal(tc.vsYAML, &vs))

			assert.Equal(vs.HasFaultInjection(), tc.expectedFaultInjection)
		})
	}

	// Testing nil case
	var vs *models.VirtualService
	assert.False(t, vs.HasFaultInjection())
}

func TestVirtualServiceHasTrafficShifting(t *testing.T) {
	cases := map[string]struct {
		vsYAML                  []byte
		expectedTrafficShifting bool
	}{
		"Has traffic shifting": {
			expectedTrafficShifting: true,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: reviews-route
spec:
  hosts:
  - reviews.prod.svc.cluster.local
  http:
  - route:
    - destination:
        host: reviews.prod.svc.cluster.local
        subset: v2
      weight: 25
    - destination:
        host: reviews.prod.svc.cluster.local
        subset: v1
      weight: 75
`),
		},
		"Single destination with no weight": {
			expectedTrafficShifting: false,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: reviews-route
spec:
  hosts:
  - reviews.prod.svc.cluster.local
  http:
  - route:
    - destination:
        host: reviews.prod.svc.cluster.local
`),
		},
		"Single destination with weight": {
			expectedTrafficShifting: false,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: reviews-route
spec:
  hosts:
  - reviews.prod.svc.cluster.local
  http:
  - route:
    - destination:
        host: reviews.prod.svc.cluster.local
        subset: v1
      weight: 100
`),
		},
		"No routes": {
			expectedTrafficShifting: false,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: reviews-route
spec:
  hosts:
  - reviews.prod.svc.cluster.local
`),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			var vs models.VirtualService
			assert.NoError(yaml.Unmarshal(tc.vsYAML, &vs))

			assert.Equal(vs.HasTrafficShifting(), tc.expectedTrafficShifting)
		})
	}

	// Testing nil case
	var vs *models.VirtualService
	assert.False(t, vs.HasTrafficShifting())
}

func TestVirtualServiceHasTCPTrafficShifting(t *testing.T) {
	cases := map[string]struct {
		vsYAML                     []byte
		expectedTCPTrafficShifting bool
	}{
		"Has traffic shifting": {
			expectedTCPTrafficShifting: true,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: tcp-echo-route
spec:
  hosts:
  - tcp-echo
  tcp:
  - match:
    - port: 31400
    route:
    - destination:
        host: tcp-echo
        port:
          number: 9000
        subset: v1
      weight: 80
    - destination:
        host: tcp-echo
        port:
          number: 9000
        subset: v2
      weight: 20
`),
		},
		"Single destination with no weight": {
			expectedTCPTrafficShifting: false,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: tcp-echo-route
spec:
  hosts:
  - tcp-echo
  tcp:
  - match:
    - port: 31400
    route:
    - destination:
        host: tcp-echo
        port:
          number: 9000
        subset: v1
`),
		},
		"Single destination with weight": {
			expectedTCPTrafficShifting: false,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: tcp-echo-route
spec:
  hosts:
  - tcp-echo
  tcp:
  - match:
    - port: 31400
    route:
    - destination:
        host: tcp-echo
        port:
          number: 9000
        subset: v1
      weight: 100
`),
		},
		"No routes": {
			expectedTCPTrafficShifting: false,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: tcp-echo-route
spec:
  hosts:
  - tcp-echo
`),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			var vs models.VirtualService
			assert.NoError(yaml.Unmarshal(tc.vsYAML, &vs))

			assert.Equal(vs.HasTCPTrafficShifting(), tc.expectedTCPTrafficShifting)
		})
	}

	// Testing nil case
	var vs *models.VirtualService
	assert.False(t, vs.HasTCPTrafficShifting())
}

func TestVirtualServiceHasRequestRouting(t *testing.T) {
	cases := map[string]struct {
		vsYAML                 []byte
		expectedRequestRouting bool
	}{
		"Has http request routing": {
			expectedRequestRouting: true,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: reviews-route
spec:
  hosts:
  - reviews.prod.svc.cluster.local
  http:
  - route:
    - destination:
        host: reviews.prod.svc.cluster.local
`),
		},
		"Has tcp request routing": {
			expectedRequestRouting: true,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: reviews-route
spec:
  hosts:
  - reviews.prod.svc.cluster.local
  tcp:
  - route:
    - destination:
        host: reviews.prod.svc.cluster.local
`),
		},
		"Has tls request routing": {
			expectedRequestRouting: true,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: reviews-route
spec:
  hosts:
  - reviews.prod.svc.cluster.local
  tls:
  - route:
    - destination:
        host: reviews.prod.svc.cluster.local
`),
		},
		"Has multiple forms of request routing": {
			expectedRequestRouting: true,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: reviews-route
spec:
  hosts:
  - reviews.prod.svc.cluster.local
  http:
  - route:
    - destination:
        host: reviews.prod.svc.cluster.local
  tcp:
  - route:
    - destination:
        host: reviews.prod.svc.cluster.local
  tls:
  - route:
    - destination:
        host: reviews.prod.svc.cluster.local
`),
		},
		"Has no request routing": {
			expectedRequestRouting: false,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: reviews-route
spec:
  hosts:
  - reviews.prod.svc.cluster.local
`),
		},
		"Has no request routing but has other options": {
			expectedRequestRouting: false,
			vsYAML: []byte(`
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: reviews-route
spec:
  hosts:
  - reviews.prod.svc.cluster.local
  http:
  - timeout: 5s
`),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			var vs models.VirtualService
			assert.NoError(yaml.Unmarshal(tc.vsYAML, &vs))

			assert.Equal(vs.HasRequestRouting(), tc.expectedRequestRouting)
		})
	}

	// Testing nil case
	var vs *models.VirtualService
	assert.False(t, vs.HasRequestRouting())
}
