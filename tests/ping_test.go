package tests_test

import (
	"context"
	"time"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	pb "github.com/guidewire-oss/fern-ginkgo-client/reporter" // Import generated protobuf package
)

var _ = Describe("PingService Real Test", func() {

	var (
		client pb.PingServiceClient
		conn   *grpc.ClientConn
	)

	BeforeEach(func() {

		if os.Getenv("GRPC_EXECUTE") != "TRUE" {
			Skip("Skipping real GRPC test as GRPC_EXECUTE != TRUE")
		}

		var err error
		conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
		Expect(err).NotTo(HaveOccurred())

		client = pb.NewPingServiceClient(conn)
	})

	AfterEach(func() {
		conn.Close()
	})

	It("should send a real ping and get a response", func() {
		resp, err := client.Ping(context.Background(), &pb.PingRequest{Message: "Pong"})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).NotTo(BeNil())
		Expect(resp.Message).To(Equal("Pong"))
	})
})

