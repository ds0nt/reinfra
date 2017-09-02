# reinfra
[![Build Status](https://travis-ci.org/ds0nt/reinfra.svg?branch=master)](https://travis-ci.org/ds0nt/reinfra)

Opinionated microservices infrastructure framework. This is an initial release that is probably still broken, but I'm implementing it in my personal set of microservices and will give it some love.

A few key features

 - easily define service infrastructure by composing components into a service struct
 - standardized port numbers and environment variables
 - service readiness

The mission here is to make creating repeated microservies within the same cluster really really easy. We should be able to define more components to model the infrastructure needed by our service, with a simple API and that works well in a running cluster, for instance with docker-compose or kubernetes.

## Creating Service Definition


This is a bit taken from my auth service as an example

```go

package authv2

import (	
	"github.com/ds0nt/reinfra"
	"github.com/ds0nt/reinfra/components"
	"github.com/ds0nt/reinfra/service"
)

type Authv2Service struct {

    // always add a service into the struct so that reinfra handles it
	
    *service.Service

    // add some reinfra components
    
	*components.GRPCServer    // a grpc server
	*components.MetricsServer // a prometheus metrics http server
	*components.Postgres      // and a postgres client

    // other fields you want in your struct
	tableName string
}

func InitNewService() *Authv2Service {
	svc := &Authv2Service{
		tableName: "credentials",
	}
    // initialize service.
    // this reflects the service, and components, and creates the objects
    // assigning them to our pointers.
	reinfra.Init(svc)
	
    // now that our postgres pointer is not nil, we can use it's methods
    // like so, let's set the AfterConnect handler
    svc.Postgres.SetAfterConnect(svc.initTables)
	
    // register our grpc service into the services grpc server
    // note: svc.Server() could also be written like svc.GRPCServer.Server()
    apiv1.RegisterAuthv2Server(svc.Server(), svc)
	
    return svc
}

func Start() error {
	svc := InitNewService()

    // reinfra.Run runs all of the services components
	for err := range reinfra.Run(context.Background(), svc) {
		logrus.Errorln(err)
	}

	return nil
}
// ...
```