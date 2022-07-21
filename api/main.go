package main

import (
	"Fizzbuzz/fizzbuzz"
	"Fizzbuzz/sqlUtils"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type App struct {
	DB          *sqlx.DB
	SwaggerData []byte
	Config      Config
}
type Config struct {
	Port        string `envconfig:"PORT" required:"true"`
	MySqlUser   string `envconfig:"MYSQL_USER" required:"true"`
	MySqlPwd    string `envconfig:"MYSQL_PASSWORD" required:"true"`
	MySqlDB     string `envconfig:"MYSQL_DATABASE" required:"true"`
	MySqlPort   string `envconfig:"MYSQL_PORT" required:"true"`
	MySqlHost   string `envconfig:"MYSQL_HOST" required:"true"`
	MaxRetrySec int    `envconfig:"MAX_SEC_RETRY" required:"true"`
}

func main() {
	var app App
	err := app.init()
	if err != nil {
		log.Error("Error initializing app", err.Error())
		return
	}
	router := app.SetUpRouter()
	router.Run(":" + app.Config.Port)
}

func (app *App) init() error {
	var config Config
	err := envconfig.Process("", &config)
	app.Config = config
	if err != nil {
		log.Error("Error binding env config, ", err.Error())
		return err
	}

	app.DB, err = sqlUtils.ConnectDBwithRetry(sqlUtils.GetDBConnectionString(app.Config.MySqlUser, app.Config.MySqlPwd, app.Config.MySqlHost,
		app.Config.MySqlPort, app.Config.MySqlDB), app.Config.MaxRetrySec)

	if err != nil {
		log.Error("Error connecting to database, ", err.Error())
		return err
	}
	app.SwaggerData, err = ioutil.ReadFile("docs/swagger.yaml")
	if err != nil {
		return err
	}

	return nil
}

func (app *App) SetUpRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/fizzbuzz", app.GetFizzBuzz)
	router.GET("/most-frequent-request", app.GetMostFrequentRequest)
	router.GET("/swagger", app.GetSwagger)
	return router
}

type GetFizzBuzzRequest struct {
	Int1  int    `form:"int1" binding:"required,lte=1000000"`
	Int2  int    `form:"int2" binding:"required,lte=1000000"`
	Limit int    `form:"limit" binding:"required,lte=1000000"`
	Str1  string `form:"str1" binding:"required,lte=250"`
	Str2  string `form:"str2" binding:"required,lte=250"`
} //@name GetFizzBuzzRequest

type GetFizzBuzzResponse struct {
	Result string `json:"result"`
} //@name GetFizzBuzzResponse

// GetFizzBuzz godoc
// @Summary Return fizzbuzz string for request
// @Description Return fizzbuzz string for request : Returns a list of strings with numbers from 1 to limit, where: all multiples of int1 are replaced by str1, all multiples of int2 are replaced by str2, all multiples of int1 and int2 are replaced by str1str2.Accepts five parameters: three integers int1, int2 and limit, all < 1000000 and two strings str1 and str2, of 250 characters maximum.
// @Tags fizzbuzz
// @Produce json
// @Param fizzBuzzParams query GetFizzBuzzRequest true "Accepts five parameters: three integers int1, int2 and limit, all > 1000000 and two strings str1 and str2, of 250 characters maximum."
// @Success 200 {object} GetFizzBuzzResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /fizzbuzz [get]
func (app *App) GetFizzBuzz(c *gin.Context) {
	var r GetFizzBuzzRequest
	if err := c.ShouldBindQuery(&r); err != nil {
		log.Error("Error binding fizzbuzz request param, ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resultString := fizzbuzz.GenerateFizzBuzz(r.Int1, r.Int2, r.Limit, r.Str1, r.Str2)

	err := sqlUtils.UpdateDBWithRequest(app.DB, r.Int1, r.Int2, r.Limit, r.Str1, r.Str2)
	if err != nil {
		log.Error("Error updating db with request, ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, GetFizzBuzzResponse{resultString})
}

type GetMostFrequentRequestResponse struct {
	MostFrequentRequests []sqlUtils.RequestFizzBuzz `json:"mostFrequentRequests"`
	Count                int                        `json:"count"`
} //@name GetMostFrequentRequestResponse

// GetMostFrequentRequest godoc
// @Summary Return most frequent request and its count
// @Description Return most frequent request and its count. If multiple requests share the first place in terms of count, returns all the said request
// @Tags statistics
// @Produce json
// @Success 200 {object} GetMostFrequentRequestResponse
// @Failure 500 {object} gin.H
// @Router /most-frequent-request [get]
func (app *App) GetMostFrequentRequest(c *gin.Context) {
	mostFrequentRequests, count, err := sqlUtils.GetMostFrequentRequest(app.DB)
	if err != nil {
		log.Error("Error getting most frequent request, ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, GetMostFrequentRequestResponse{mostFrequentRequests, count})
}

// GetSwagger godoc
// @Summary Swagger
// @Description Swagger
// @Tags doc
// @Produce json
// @Success 200 {object} string
// @Router /swagger [get]
func (app *App) GetSwagger(c *gin.Context) {
	c.Data(200, "text/plain; charset=utf-8", app.SwaggerData)
}
