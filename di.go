package fall

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	constructors = make(map[string]func() (any, error))
	instances    = sync.Map{}
	mu           sync.Mutex
)

func Register(name string, constructor func() (any, error)) {
	mu.Lock()
	defer mu.Unlock()
	constructors[name] = constructor
}

func Resolve(name string) (any, error) {
	mu.Lock()
	defer mu.Unlock()
	return resolve(name)
}

func resolve(name string) (any, error) {
	if instance, ok := instances.Load(name); ok {
		return instance, nil
	}

	construtor, ok := constructors[name]
	if !ok {
		return nil, fmt.Errorf("dependency not registered: %s", name)
	}

	instance, err := construtor()
	if err != nil {
		return nil, fmt.Errorf("construction failed for %s: %w", name, err)
	}

	if err := autoInject(instance); err != nil {
		return nil, fmt.Errorf("injection failed for %s: %w", name, err)
	}

	if initializable, ok := instance.(Initializable); ok {
		if err := initializable.Init(); err != nil {
			return nil, fmt.Errorf("initialization failed for %s: %w", name, err)
		}
	}

	instances.Store(name, instance)
	return instance, nil
}

func Store(name string, instance any) {
	instances.Store(name, instance)
}

func ResolveControllers() []Controller {
	var controllers []Controller
	mu.Lock()
	defer mu.Unlock()
	for name := range constructors {
		instance, err := resolve(name)
		if err != nil {
			PanicIfError(err)
		}
		if controller, ok := instance.(Controller); ok {
			controllers = append(controllers, controller)
		}
	}
	return controllers
}

func autoInject(instance interface{}) error {
	val := reflect.ValueOf(instance)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to struct")
	}

	elem := val.Elem()
	return injectRecursive(elem)
}

func injectRecursive(val reflect.Value) error {
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Processa campos embedded recursivamente
		if field.Anonymous && fieldVal.Kind() == reflect.Struct {
			if err := injectRecursive(fieldVal.Addr().Elem()); err != nil {
				return err
			}
			continue
		}

		// Processa tags DI
		tag := field.Tag.Get("fall")
		if tag == "" {
			continue
		}

		dependency, err := resolve(tag)
		if err != nil {
			return fmt.Errorf("failed to resolve %s for field %s: %w", tag, field.Name, err)
		}

		if !fieldVal.CanSet() {
			return fmt.Errorf("cannot set field %s", field.Name)
		}

		fieldVal.Set(reflect.ValueOf(dependency))
	}

	return nil
}
