package main

import "github.com/asaskevich/govalidator"

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}
