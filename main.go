package main

import (
	"go-final/controlloer"

	"github.com/gin-gonic/gin"
)

// https://sql2gorm.mccode.info     model
func main() {

	// pass, _ := hashPassword("123456")
	// fmt.Println(pass)
	gin.SetMode(gin.ReleaseMode)
	controlloer.StartServer()

}

// $2a$10$wX8IEuJiq.E.RDIa3aWEVuM86hEz07.UP2g4LLvMvZb.nqkmSvyRC
