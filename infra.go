package reinfra

import (
	"context"
	"fmt"
	"reflect"

	"github.com/ds0nt/reinfra/components"
	"github.com/ds0nt/reinfra/service"
)

var (
	dialerType    = reflect.TypeOf((*components.GRPCDialer)(nil))
	initerType    = reflect.TypeOf((*service.Initer)(nil))
	componentType = reflect.TypeOf((*service.ServiceComponent)(nil))
)

func isHandledPkg(pkgPath string) bool {
	var (
		pkg          = dialerType.PkgPath()
		componentPkg = pkg + "/components"
		// clientPkg    = pkg + "/clients"
	)
	switch pkgPath {
	case pkg, componentPkg:
		return true
	}
	return false
}

// Init sets the values for all infra field pointers in a service
func Init(obj interface{}) {
	var (
		objT = reflect.TypeOf(obj)
		val  = reflect.ValueOf(obj).Elem()
	)

	// instantiate all pointers
	for i := 0; i < val.NumField(); i++ {
		tf := val.Type().Field(i)
		if !isHandledPkg(tf.PkgPath) {
			continue
		}
		f := val.Field(i)
		if !f.IsNil() {
			continue
		}

		fmt.Printf("Initializing %s.%s\n", objT.String(), f.Type().String())
		f.Set(reflect.New(f.Type().Elem()))

	}

	// get instantiated service
	svc := reflectService(obj)
	if svc == nil {
		panic("cannot init a non service.Service")
	}

	// run service init
	svc.Init()

	// run components init if they have it
	cmps := reflectServiceComponents(obj)
	for _, c := range cmps {
		if v, ok := c.(service.Initer); ok {
			fmt.Printf("Running Init(svc) for %s\n", reflect.TypeOf(c).String())
			v.Init(svc)
		}
	}
}

// Run is one run method to rule them all
// ready channel is closed when service is ready
func Run(ctx context.Context, obj interface{}) chan error {
	svc := reflectService(obj)
	if svc == nil {
		panic("cannot run a non service.Service")
	}

	cmps := reflectServiceComponents(obj)
	for _, c := range cmps {
		fmt.Printf("Registering %s\n", reflect.TypeOf(c).String())
		svc.RegisterComponent(c)
	}

	return svc.Run(ctx)
}

func reflectService(obj interface{}) *service.Service {
	val := reflect.ValueOf(obj).Elem()

	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		tf := val.Type().Field(i)

		if !isHandledPkg(tf.PkgPath) {
			continue
		}
		if x, ok := f.Interface().(*service.Service); ok {
			return x
		}
	}

	return nil
}

func reflectServiceComponents(obj interface{}) []service.ServiceComponent {
	cmps := []service.ServiceComponent{}

	val := reflect.ValueOf(obj).Elem()
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		tf := val.Type().Field(i)

		if !isHandledPkg(tf.PkgPath) {
			continue
		}

		valueField := f.Interface()

		fmt.Printf("Scanning %s.%s\n", val.Type().String(), tf.Type.String())
		if tf.Type.Implements(dialerType.Elem()) {
			cmps = append(cmps, components.WrapDialer(valueField.(components.GRPCDialer)))
			continue
		}

		if tf.Type.Implements(componentType.Elem()) {
			cmps = append(cmps, valueField.(service.ServiceComponent))
			continue
		}
	}

	return cmps
}
