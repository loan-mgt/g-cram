package service

import (
	"context"
	"encoding/json"
	"fmt"
	"loan-mgt/g-cram/internal/config"
	"loan-mgt/g-cram/internal/db"
	"loan-mgt/g-cram/internal/db/sqlc"
	"loan-mgt/g-cram/internal/server/ws"
	"math"
)

func StartUploadDoneListener(ws *ws.WebSocketManager, db *db.Store, cfg *config.Config) {

	notification, err := NewAMQPConnection(config.New(), "uploader_done")
	if err != nil {
		panic(err)
	}

	defer notification.Conn.Close()

	// Consume messages from queue
	msgs, err := notification.Consume()
	if err != nil {
		fmt.Println(err)
		return
	}

	for d := range msgs {
		var payload struct {
			MediaId   string `json:"mediaId"`
			UserId    string `json:"userId"`
			Timestamp int64  `json:"timestamp"`
		}

		if err := json.Unmarshal(d.Body, &payload); err != nil {
			fmt.Println("Error unmarshal ling message:", err)
			continue
		}

		param := sqlc.GetMediaParams{
			MediaID:   payload.MediaId,
			UserID:    payload.UserId,
			Timestamp: payload.Timestamp,
		}

		media, err := db.GetMedia(context.Background(), param)
		if err != nil {
			fmt.Println("Error getting media:", err)
			continue
		}

		// update step +1 and set done
		err = db.SetMediaStep(context.Background(), sqlc.SetMediaStepParams{
			MediaID:   media.MediaID,
			UserID:    media.UserID,
			Timestamp: media.Timestamp,
			Step:      media.Step + 1,
		})
		if err != nil {
			fmt.Println("Error updating media step:", err)
			continue
		}

		err = db.SetMediaDone(context.Background(), sqlc.SetMediaDoneParams{
			MediaID:   media.MediaID,
			UserID:    media.UserID,
			Timestamp: media.Timestamp,
			Done:      1,
		})
		if err != nil {
			fmt.Println("Error updating media done:", err)
			continue
		}

		if err := ws.SendMessageToClient(payload.UserId, "Success"); err != nil {
			fmt.Println("Error sending message to WebSocket:", err)
		}

		paramCount := sqlc.CountUserMediaInJobParams{
			UserID:    payload.UserId,
			Timestamp: payload.Timestamp,
		}

		nbMedia, err := db.CountUserMediaInJob(context.Background(), paramCount)
		if err != nil {
			fmt.Println("Error getting nb media:", err)
			nbMedia = 0
			continue
		}

		nbMediaDone, err := db.CountUserMediaInJobAtStep(context.Background(), sqlc.CountUserMediaInJobAtStepParams{
			UserID:    payload.UserId,
			Timestamp: payload.Timestamp,
			Step:      2,
		})
		if err != nil {
			fmt.Println("Error getting nb media done:", err)
			nbMediaDone = 0
			continue
		}

		err = PushNotification(db, cfg, payload.UserId, fmt.Sprintf("Job done: %d/%d", nbMediaDone, nbMedia), fmt.Sprintf("%d", payload.Timestamp))
		if err != nil {
			fmt.Println("Error sending push notification:", err)
		}

		if nbMedia == nbMediaDone {

			jobSpace, err := db.GetJobSpace(context.Background(), sqlc.GetJobSpaceParams{
				UserID:    payload.UserId,
				Timestamp: payload.Timestamp,
			})
			if err != nil {
				fmt.Println("Error getting job space:", err)
				continue
			}

			spaceSaved := math.Round((1 - (jobSpace.SumNewSize.Float64 / jobSpace.SumOldSize.Float64)) * 100)

			err = PushNotificationFull(db, cfg, payload.UserId, fmt.Sprintf("Job done: %.0f%% space saved", spaceSaved), fmt.Sprintf("%dDONE", payload.Timestamp), false)
			if err != nil {
				fmt.Println("Error sending push notification:", err)
			}
		}

	}

}
