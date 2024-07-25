*Step1:*
if he asks to create a http service
I will create it in the main.go file

router := gin.New()
router.GET("/id", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"id": "endpoint",
			})
		})


*Step2:*
If he asks for grpc service 

i.e first create a <file_name>.proto
syntax = "proto3";
option go_package = "github.com/rohit907/grpc-service/invoicer";

message Amount {
  int64 amount = 1;
  string currency = 2;
}

message CreateRequest {
  Amount amount = 1;
  string from = 2;
  string to = 3;
  string VAt=4;
}

message CreateResponse {
  string pdf = 1;
  string docx = 2;

}

service Currency {
  rpc Create(CreateRequest) returns (CreateResponse);
}


*Step3:*
create Makefile
// i wont remember this command for sure
gen:
	protoc \
    --go_out=invoicer \
    --go_opt=paths=source_relative \
    --go-grpc_out=invoicer \
    --go-grpc_opt=paths=source_relative \
    amount.proto

 make gen  


* Step 4:*
 Come back to main.go 
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


*Step 5:*
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


*Step 6:*
    wg.Add(1)
	go func() {
		router.GET("/id", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"id": "endpoint",
			})
		})
        router.Run()
    }()
	wg.Wait()


    Step 7:
    if he asks to make them tak to each other then goes inside the 
    step6 function

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