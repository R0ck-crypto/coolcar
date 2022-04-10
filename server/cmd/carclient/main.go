package main

import (
	"context"
	carpb "coolcar/car/api/gen/v1"
	"fmt"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8084", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	carClient := carpb.NewCarServiceClient(conn)
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		car, err := carClient.CreateCar(ctx, &carpb.CreateCarRequest{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("created car :%s\n", car.Id)
	}

}
