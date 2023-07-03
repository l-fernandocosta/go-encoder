package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/queue"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

type JobManager struct {
	Db               *gorm.DB
	Domain           domain.Job
	MessageChannel   chan amqp.Delivery
	JobReturnChannel chan JobWorkerResult
	RabbitMQ         *queue.RabbitMQ
}

type JobNotificationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewJobManager(db *gorm.DB, rabbitMQ *queue.RabbitMQ, jobReturnChannel chan JobWorkerResult, messageChannel chan amqp.Delivery) *JobManager {
	return &JobManager{
		Db:               db,
		Domain:           domain.Job{},
		MessageChannel:   messageChannel,
		JobReturnChannel: jobReturnChannel,
		RabbitMQ:         rabbitMQ,
	}
}

func (j *JobManager) Start(ch *amqp.Channel) {
	vs := NewVideoService()
	vs.VideoRepository = repositories.VideoRepositoryDB{Db: j.Db}

	js := JobService{
		JobRepository: repositories.JobRepositoryDB{Db: j.Db},
		VideoService:  vs,
	}
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))
	if err != nil {
		log.Fatalf("Error loading var: CONCURRENCY_WORKERS")
	}

	for qtdProcesses := 0; qtdProcesses < concurrency; qtdProcesses++ {
		go JobWorker(j.MessageChannel, j.JobReturnChannel, js, j.Domain, qtdProcesses)
	}

	for jobResult := range j.JobReturnChannel {
		if jobResult.Error != nil {
			err = j.checkParseErrors(jobResult)
		} else {
			err = j.notifySuccess(jobResult, ch)
		}

		if err != nil {
			jobResult.Message.Reject(false)
		}
	}
}

func (jm *JobManager) checkParseErrors(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Printf("MessageID #{jobResult.Message.DeliveryTag}. Error parsing JOB: #{jobResult.Job.ID}.")
	} else {
		log.Printf("MessageID #{jobResult.Message.DeliveryTag}. Error parsing MESSAGE: #{jobResult.Error}.")
	}

	errorMsg := JobNotificationError{
		Message: string(jobResult.Message.Body),
		Error:   jobResult.Error.Error(),
	}

	jobJson, err := json.Marshal(errorMsg)

	if err != nil {
		return err
	}

	err = jm.notify(jobJson)

	if err != nil {
		return err
	}

	err = jobResult.Message.Reject(false)

	if err != nil {
		return err
	}

	return nil

}

func (jm *JobManager) notify(jobJson []byte) error {
	err := jm.RabbitMQ.Notify(
		string(jobJson),                                //msg
		"application/json",                             //content type
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),          //exchange
		os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"), // routingKey
	)

	if err != nil {
		return err
	}

	return nil
}

func (jm *JobManager) notifySuccess(jobResult JobWorkerResult, ch *amqp.Channel) error {
	jobJson, err := json.Marshal(jobResult.Job)

	if err != nil {
		return err
	}

	err = jm.notify(jobJson)

	if err != nil {
		return err
	}

	err = jobResult.Message.Ack(false)

	if err != nil {
		return err
	}

	return nil
}
