package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rohit907/http-grpc/invoicer"
	"google.golang.org/grpc"
)

type CurrencyServer struct {
	invoicer.UnimplementedCurrencyServer
}

func (my CurrencyServer) Create(ctx context.Context, req *invoicer.CreateRequest) (*invoicer.CreateResponse, error) {
	return &invoicer.CreateResponse{
		Pdf:  req.From,
		Docx: req.VAt,
	}, nil

}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", ":9000")
		if err != nil {
			log.Fatalf("Error while listening the port %s ", err)
		}
		server := grpc.NewServer()
		service := &CurrencyServer{}
		invoicer.RegisterCurrencyServer(server, service)
		err = server.Serve(lis)
		if err != nil {
			log.Fatalf("server registery failed", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		router := gin.New()
		router.GET("/", func(c *gin.Context) {
			conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"data": "message",
				})
			}
			defer conn.Close()

			grpcClient := invoicer.NewCurrencyClient(conn)
			response, err := grpcClient.Create(context.Background(), &invoicer.CreateRequest{
				Amount: &invoicer.Amount{
					Amount:   0,
					Currency: "USD",
				},
				From: "test",
				To:   "adas",
				VAt:  "asdas",
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})

			}
			c.JSON(http.StatusOK, gin.H{
				"data": fmt.Sprintf("GRPC server called, %s", response.Docx),
			})
		})
		router.GET("/id", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"id": "http endpoint",
			})
		})

		router.Run()

	}()
	wg.Wait()

}
