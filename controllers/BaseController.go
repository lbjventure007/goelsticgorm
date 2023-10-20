package controllers

import "fmt"

type BaseController struct {
}

func (b BaseController) Log(message string) {

	fmt.Println(message)
}
