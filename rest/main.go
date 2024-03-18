package main

import (
	"context"
	"log"
	"os"

	"restapi/router"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

func init() {
	os.Setenv("PORT", "8080")
	// test or production
	os.Setenv("MODE", "test")
	os.Setenv("AWS_ACCESS_KEY_ID", "id")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "token")
	os.Setenv("LEMONSQUEEZY_API_KEY", "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiI5NGQ1OWNlZi1kYmI4LTRlYTUtYjE3OC1kMjU0MGZjZDY5MTkiLCJqdGkiOiI5NDFmOWI4OGIxNTI4YjBkMjcyN2RlMWI2MzVjNTMzM2ZmNDY2YzkxM2NmMTZkNzY2ZDhlZjBhNWEwYzM5ZTU0MjhjMjcxOTNjNjIzNWI1MyIsImlhdCI6MTcwODIzMDMyMy41OTk4MDcsIm5iZiI6MTcwODIzMDMyMy41OTk4MDksImV4cCI6MjAyMzg0OTUyMy41NzI3MzgsInN1YiI6IjE3OTQ5MjkiLCJzY29wZXMiOltdfQ.QLhqUL3tMqYrl5Emzm4y2A6Rlo3eNfuDi-vJTu0rwY8hCKWKjzYHVPbLK9k__aF5ci1Hv585hikCMcSLq8XQ-n7f5xvP2AjuzkoUmtUuLJqxizZUY-uYxP9chbYIofSz3KTu-lG4GbAVMpPdtRC-4OA2RvGa1mg8Nl65HmGdrzJSO1G3KcG2YHMx1C-CA0rDxo0DSF8m_UFrP8_tvebfDZr244LQSX-tkPVth5wpkSR9TMYvEW0SniZcgrY9tacOs1ysayC6HRv3vuo6NvWKoY5IPxDOw_Aj_FSrxEZMvQCv0NczTHjldG_aleNPgEFg9Eli9OufYVpZewPVOJm7ywrDFj_byQoS7Mpioz5aq5xGtikSGpVT3Rmq6YljpvVlI-M1Y5Qan5M3mH8Xc9bXSpG2nw277I4W7N_nfHnXekHUgi4t4LhZvSsLmGdKNZrx904tkOgMBFRIfmlgZ9b7Dpcnp8M8eLigId1nDUSn-ZF8QnLAf4OkN4jFYs4zpYZm")
	os.Setenv("LEMONSQUEEZY_STORE_ID", "64355")
	os.Setenv("PRODUCT_TABLE_NAME", "")
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// middleware

	// - No origin allowed by default
	// - GET,POST, PUT, HEAD methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:6006"}
	r.Use(cors.New(config))
	// router group
	v0 := r.Group("/v0")
	{
		router.Routes(v0.Group(""))
	}
	return r
}

func main() {
	log.Printf("Gin cold start")
	r := SetupRouter()
	//here the false case
	mode := os.Getenv("MODE")
	switch mode {
	//here the false case
	case "production":
		ginLambda = ginadapter.New(r)
		lambda.Start(Handler)
	default:
		log.Printf("now runnming mode is:%v", os.Getenv("MODE"))
		r.Run()
	}
}
