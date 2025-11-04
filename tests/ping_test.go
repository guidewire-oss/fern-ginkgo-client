package tests_test

import (
	"context"
	"os"

	pb "github.com/guidewire-oss/fern-ginkgo-client/reporter" // Import generated protobuf package
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	if os.Getenv("GRPC_EXECUTE") == "TRUE" {
		var _ = Describe("PingService Real Test", func() {

			var (
				client pb.PingServiceClient
				conn   *grpc.ClientConn
			)

			BeforeEach(func() {
				var err error
				conn, err = grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
				Expect(err).NotTo(HaveOccurred())

				client = pb.NewPingServiceClient(conn)
			})

			AfterEach(func() {
				if conn != nil {
					_ = conn.Close()
				}
			})

			It("should send a real ping and get a response", func() {
				resp, err := client.Ping(context.Background(), &pb.PingRequest{Message: "Pong"})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.Message).To(Equal("Pong"))
			})
		})
	}
}
