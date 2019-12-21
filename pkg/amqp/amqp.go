package amqp

import (
	"github.com/streadway/amqp"
)

// Queues is an entity witholding connections to queues and references to channels/queues
type Queue struct {
	AmqpConnection *amqp.Connection
	Channel        *amqp.Channel
	Entities       map[string]amqp.Queue
}

// NewAmqpConnection takes a set of queue names, creates
// those channels, and establishes connections to such queues
func NewAmqpConnection(connstring string, queueNames []string) (Queue, error) {
	conn, err := amqp.Dial(connstring)

	if err != nil {
		return Queue{}, err
	}

	ch, err := conn.Channel()

	if err != nil {
		return Queue{}, err
	}

	var (
		queue  amqp.Queue
		queues = make(map[string]amqp.Queue)
	)
	for _, channelName := range queueNames {
		queue, err = ch.QueueDeclare(
			channelName,
			true,  /*durable connection*/
			false, /*auto delete*/
			false, /*exclusive*/
			false, /*nowai*t*/
			nil /*arg amqp.table*/)

		if err != nil {
			return Queue{}, err
		}
		queues[channelName] = queue
	}

	return Queue{AmqpConnection: conn, Channel: ch, Entities: queues}, nil
}

// SendMessageToQueue sends messages to queues
func (q Queue) SendMessageToQueue(message string, queueName string) error {
	publishedMsg := amqp.Publishing{
		DeliveryMode: 2, // persistent msg delivery
		Body:         []byte(message),
		Priority:     4,
		Type:         "SendWelcomeEmail",
	}

	err := q.Channel.Publish(
		"",                         // exchange string
		q.Entities[queueName].Name, // key string
		false,                      // mandatory
		false,                      // immediate
		publishedMsg)

	if err != nil {
		return err
	}

	return nil
}

// ConsumerMessageFromQueue consumes messages from queues
func (q Queue) ConsumerMessageFromQueue(message string, queueName string) ([]interface{}, error) {
	var response []interface{}

	msgs, err := q.Channel.Consume(
		q.Entities[queueName].Name, // queue string
		"",                         // consumer string
		true,                       // auto ack
		false,                      // exclusive bool
		false,                      // no local
		false,                      // no wait
		nil)

	if err != nil {
		return nil, err
	}

	for m := range msgs {
		response = append(response, m.Body)
	}

	return response, nil
}
