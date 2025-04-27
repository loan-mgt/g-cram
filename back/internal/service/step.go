package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"loan-mgt/g-cram/internal/config"
	"loan-mgt/g-cram/internal/db"
	"loan-mgt/g-cram/internal/db/sqlc"
	"loan-mgt/g-cram/internal/server/ws"
)

func StartDownloadDoneListener(ws *ws.WebSocketManager, db *db.Store, cfg *config.Config) {
	ctx := context.Background()
	notification, err := NewAMQPConnection(config.New(), "downloader_done")
	if err != nil {
		panic(err)
	}
	defer notification.Conn.Close()

	next, err := NewAMQPConnection(config.New(), "cramer")
	if err != nil {
		panic(err)
	}
	defer next.Conn.Close()

	msgs, err := notification.Consume()
	if err != nil {
		return
	}

	for d := range msgs {
		var payload struct {
			MediaID   string `json:"media_id"`
			UserID    string `json:"user_id"`
			Timestamp int64  `json:"timestamp"`
			FileSize  int64  `json:"file_size"`
		}

		if err := json.Unmarshal(d.Body, &payload); err != nil {
			continue
		}

		params := sqlc.SetMediaNewSizeParams{
			MediaID:   payload.MediaID,
			UserID:    payload.UserID,
			Timestamp: payload.Timestamp,
			NewSize: sql.NullInt64{
				Int64: payload.FileSize,
				Valid: true,
			},
		}

		if err := db.SetMediaNewSize(ctx, params); err != nil {
			fmt.Println(err)
			continue
		}

		media, err := db.GetMedia(ctx, sqlc.GetMediaParams{
			MediaID:   payload.MediaID,
			UserID:    payload.UserID,
			Timestamp: payload.Timestamp,
		})
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err := ws.SendMessageToClient(payload.UserID, fmt.Sprintf("{\"event\":\"DownloadDone\",\"step\":\"%d\"}", media.Step)); err != nil {
			fmt.Println(err)
			continue
		}

		if err := PushNotification(db, cfg, payload.UserID, fmt.Sprintf("Notification: %s", media.Filename), "download-done:"+media.Filename); err != nil {
			fmt.Println(err)
			continue
		}

		nextPayload := struct {
			MediaID   string `json:"media_id"`
			UserID    string `json:"user_id"`
			Timestamp int64  `json:"timestamp"`
		}{
			MediaID:   payload.MediaID,
			UserID:    payload.UserID,
			Timestamp: payload.Timestamp,
		}

		jsonData, err := json.Marshal(nextPayload)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err := next.SendRequest(jsonData); err != nil {
			fmt.Println(err)
			continue
		}

	}

}
