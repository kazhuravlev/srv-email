package sender

import (
	"github.com/kazhuravlev/srv-email/contracts"
	"github.com/kazhuravlev/srv-email/services/kafka"
)

type taskContainer struct {
	Message      contracts.Msg
	KafkaMessage kafka.Message
}
