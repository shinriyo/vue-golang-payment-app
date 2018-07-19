package main

import (
	"context"
	"log"
	"net"

	gpay "vue-golang-payment-app/payment-service/proto"

	payjp "github.com/payjp/payjp-go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement sa
type server struct{}

func (s *server) Charge(ctx context.Context, req *gpay.PayRequest) (*gpay.PayResponse, error) {
	// PAI の初期化
	pay := payjp.New("sk_test_bfb67d3a18aca4e01e20c7ac", nil)

	// 支払いをします。第一引数に支払い金額、第二引数に支払いの方法や設定を入れます。
	charge, err := pay.Charge.Create(int(req.Amount), payjp.Charge{
		// 現在はjpyのみサポート
		Currency: "jpy",
		// カード情報、顧客ID、カードトークンのいずれかを指定
		Card: payjp.Card{
			Number:   req.Num,
			CVC:      req.Cvc,
			ExpMonth: req.Expm,
			ExpYear:  req.Expy,
		},
		Capture: true,
		// 概要のテキストを設定できます
		Description: "Book: 'The Art of Community'",
		// 追加のメタデータを20件まで設定できます
		Metadata: map[string]string{
			"ISBN": "1449312063",
		},
	})
	if err != nil {
		return nil, err
	}

	// 支払った結果から、Response生成
	res := &gpay.PayResponse{
		Paid:     charge.Paid,
		Captured: charge.Captured,
		Amount:   int64(charge.Amount),
	}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	gpay.RegisterPayManagerServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Printf("gRPC Server started: localhost%s\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
