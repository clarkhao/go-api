package ping

import (
	"log"
	"net/http"
	"restapi/service/ping"
	"time"

	"github.com/gin-gonic/gin"

	re "restapi/utils/error"
	"restapi/utils/request"
)

type Params struct {
	Name string `form:"name" binding:"required"`
	Age  int    `form:"age"`
}
type Form struct {
	FirstName string `form:"firstname" binding:"required"`
	LastName  string `form:"lastname" binding:"required"`
}

func GetHandler(c *gin.Context) {
	var params Params
	//query string
	if err := c.ShouldBind(&params); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "missing query string",
		})
		return
	}
	//response add headers
	c.Header("hello", "world")
	c.JSON(http.StatusOK, gin.H{
		//read request from c.request
		"path": c.Request.URL.Path,
		"name": params.Name,
		"age":  params.Age,
	})
}

func PostHandler(c *gin.Context) {
	var params Params
	var form Form
	// query params
	c.ShouldBind(&params)
	ids := c.QueryMap("ids")
	// content-type: application/x-www-form-urlencoded
	// content-type: multipart/form-data
	if err := c.ShouldBind(&form); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "missing body",
		})
		return
	}
	// response add headers
	c.Header("hello", "world")
	c.JSON(http.StatusOK, gin.H{
		//read request from c.request
		"path":      c.Request.URL.Path,
		"name":      params.Name,
		"age":       params.Age,
		"ids":       ids,
		"firstname": form.FirstName,
		"lastname":  form.LastName,
	})
}

func ErrorHandler(c *gin.Context) {
	err := ping.SomeServiceErr()
	e, ok := err.(re.RequestError)
	if ok {
		c.JSON(e.Code, gin.H{
			"msg": e.Cause.Error(),
		})
	} else {
		if err == nil {
			c.JSON(200, gin.H{
				"msg": "Hello World",
			})
		} else {
			c.JSON(500, gin.H{
				"msg": "inner server mistake",
			})
		}
	}

}

type ClientUser struct {
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	ID        string    `json:"id"`
}

func ClientGetHandler(c *gin.Context) {
	items, err := request.Client.NewClient("https://641b10fb9b82ded29d494d1c.mockapi.io/api/posts").SetRequest("GET", map[string]string{}, nil, "/user").GetRequest()
	if err != nil {
		e, _ := err.(re.RequestError)
		c.JSON(e.Code, gin.H{
			"msg": e.Msg,
		})
	} else {
		c.JSON(200, gin.H{
			"data": items,
		})
	}
}

func ClientPostHandler(c *gin.Context) {
	items, err := request.Client.NewClient("https://641b10fb9b82ded29d494d1c.mockapi.io/api/posts").SetRequest("POST", map[string]string{"name": "Clark"}, nil, "/user").PostRequest()
	if err != nil {
		e, _ := err.(re.RequestError)
		c.JSON(e.Code, gin.H{
			"msg": e.Msg,
		})
	} else {
		c.JSON(200, gin.H{
			"data": items,
		})
	}
}
