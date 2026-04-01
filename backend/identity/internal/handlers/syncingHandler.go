package handlers

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/Egot3/supel/backend/identity/internal/database/repositories"
	"github.com/Egot3/supel/backend/identity/internal/models"
	"github.com/Egot3/supel/backend/identity/internal/types"
	"github.com/rabbitmq/amqp091-go"
)

func HandleSyncMessage(ch <-chan amqp091.Delivery){
	for {
		msg, ok := <-ch
		if !ok {
			break
		}

		var body types.SyncData
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			log.Println("Bad rabbitMQ requets: ", err)
			continue
		}

		if err := repositories.UpsertUser(context.Background(), models.User{
			UUID: body.UUID,
			Role: body.Role,
		}); err != nil {
			log.Println("Couldn't upsert a user: ", err)
			continue
		}

		syncModel := models.Sync{
			Source: msg.Exchange,
			LastSync: time.Now(),
			LastMessageId: strconv.FormatUint(msg.DeliveryTag, 10),
		}
		if err := repositories.UpdateSyncing(context.Background(), syncModel); err != nil{
			log.Printf("Couldn't update/create update entry: %v", err)
			msg.Nack(false, false)
			continue
		}

		msg.Ack(false)
	}
}