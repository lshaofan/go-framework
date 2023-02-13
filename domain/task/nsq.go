package task

import "github.com/nsqio/go-nsq"

type Nsq struct {
	Topic   string
	Body    []byte
	config  *nsq.Config
	address string
}

type Options func(*Nsq)

func NewNsq(addr string, opt ...Options) *Nsq {
	n := &Nsq{
		address: addr,
	}
	for _, o := range opt {
		o(n)
	}
	if n.config == nil {
		n.config = nsq.NewConfig()
	}

	return n
}

func WithConfig(config *nsq.Config) Options {
	return func(option *Nsq) {
		option.config = config
	}
}

// GetProducer  get producer
func (p *Nsq) GetProducer() (*nsq.Producer, error) {
	producer, err := nsq.NewProducer(p.address, p.config)
	if err != nil {
		return nil, err
	}
	// ping
	err = producer.Ping()
	if err != nil {
		return nil, err
	}
	return producer, nil
}

// Publish publish message
func (p *Nsq) Publish(topic string, msg []byte) error {
	producer, err := p.GetProducer()
	if err != nil {
		return err
	}
	err = producer.Publish(topic, msg)
	if err != nil {
		return err
	}
	return nil
}
