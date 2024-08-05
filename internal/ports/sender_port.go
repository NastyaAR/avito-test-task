package ports

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Sender struct{}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) SendEmail(ctx context.Context, recipient string, message string) error {
	duration := time.Duration(rand.Int63n(3000)) * time.Millisecond
	time.Sleep(duration)

	errorProbability := 0.1
	if rand.Float64() < errorProbability {
		return errors.New("internal error")
	}
	fmt.Printf("send message '%s' to '%s'\n", message, recipient)

	return nil
}
