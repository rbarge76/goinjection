package goinjection

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
)

type Application struct {
	services []interface{}
	config   map[string]interface{}
}

func NewApplication() *Application {

	return &Application{config: make(map[string]interface{})}
}

func (this *Application) ReadConfig(name string) error {

	raw, err := ioutil.ReadFile(name)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(raw, &this.config)

	if err != nil {
		return err
	}

	return nil
}

func (this *Application) AddDefaultValue(name string, value interface{}) {
	this.config[name] = value
}

func (this *Application) AddService(service interface{}) error {
	return this.AddServiceWithName(service, "")

}

func (this *Application) AddServiceWithName(service interface{}, name string) error {
	// make sure service is interface

	if reflect.TypeOf(service).Kind() != reflect.Ptr {
		return errors.New("Service must be a pointer")
	}

	this.services = append(this.services, service)

	return nil
}

func (this *Application) fieldValue(field reflect.StructField) (interface{}, error) {
	for _, service := range this.services {
		serviceType := reflect.TypeOf(service)
		if field.Type.Kind() == reflect.Interface {
			if serviceType.Implements(field.Type) {
				return service, nil
			}
		} else if field.Type.Kind() == reflect.Struct {

			return nil, errors.New("Injected type must not be a struct (use pointer or interface)")
		} else {
			if serviceType.Elem() == field.Type.Elem() {
				return service, nil
			}
		}

	}
	return nil, nil
}

func (this *Application) searchValue(key string, packageName string) interface{} {

	log.Println(packageName)

	keyWithPackage := packageName + "." + key

	log.Println("Looking for ", keyWithPackage)
	if res := this.lookupValue(keyWithPackage); res != nil {
		return res
	}
	log.Println("Looking for ", key)
	return this.lookupValue(key)

}

func (this *Application) lookupValue(key string) interface{} {

	value := os.Getenv(key)

	if value != "" {
		return value
	}

	value = os.Getenv(strings.ToUpper(key))
	if value != "" {
		return value
	}

	return this.config[key]
}

const INJECTTAG = "inject"
const INJECTOPTINALTAG = "inject"

const INJECTVALUE = "injectValue"
const INJECTVALUEOPTINAL = "injectValueOptional"

func (this *Application) Wire() error {
	for _, service := range this.services {
		meta := reflect.TypeOf(service).Elem()

		packageName := meta.PkgPath()

		log.Println("Checking service ", meta.Name())

		for i := 0; i < meta.NumField(); i++ {
			field := meta.Field(i)
			tag := string(field.Tag)
			tags := strings.Split(tag, ":")

			if len(tags) == 0 {
				continue
			}

			if tags[0] == INJECTTAG || tags[0] == INJECTOPTINALTAG {

				log.Println("Found tag on " + field.Name)
				linkedService, err := this.fieldValue(field)
				if err != nil {
					return err
				}

				// infect the type
				if linkedService != nil {
					fieldValue := reflect.ValueOf(service).Elem()
					linkValue := reflect.ValueOf(linkedService)
					fieldValue.Field(i).Set(linkValue)

				} else if tags[0] == INJECTTAG {
					return errors.New("Could not find dependency for " + field.Name)
				}

			} else if tags[0] == INJECTVALUE || tags[0] == INJECTVALUEOPTINAL {

				value := this.searchValue(tags[1], packageName)

				if value != nil {

					valueType := reflect.TypeOf(value)

					fieldValue := reflect.ValueOf(service).Elem()
					linkValue := reflect.ValueOf(value)

					// make sure these are the same type

					if fieldValue.Field(i).Type() != valueType {
						return errors.New("Values do not match")
					}

					fieldValue.Field(i).Set(linkValue)

				} else if tags[0] == INJECTVALUE {
					return errors.New("Could not find dependnecy value " + tags[1])
				}

			}
		}
	}

	for _, service := range this.services {
		if m, ok := interface{}(service).(ApplicationSetup); ok {
			meta := reflect.TypeOf(service).Elem()
			log.Println("Settingup " + meta.Name())
			err := m.DoSetup()
			log.Println("Setting complete")
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (this *Application) Close() {

}
