package utils

import "reflect"

func InjectFieldValue( target interface{}, fieldName string, toBeInjected interface{}){

	typeOfT:= reflect.ValueOf( target).Elem().FieldByName( fieldName )

	if typeOfT.CanSet() {
		typeOfT.Set( reflect.ValueOf( toBeInjected ) )
	}
}
