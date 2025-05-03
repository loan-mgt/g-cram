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
	fmt.Println("StartDownloadDoneListener")
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
			MediaId   string `json:"mediaId"`
			UserId    string `json:"userId"`
			Timestamp int64  `json:"timestamp"`
			FileSize  int64  `json:"fileSize"`
		}

		if err := json.Unmarshal(d.Body, &payload); err != nil {
			continue
		}

		fmt.Println(payload)

		params := sqlc.SetMediaOldSizeParams{
			MediaID:   payload.MediaId,
			UserID:    payload.UserId,
			Timestamp: payload.Timestamp,
			OldSize: sql.NullInt64{
				Int64: payload.FileSize,
				Valid: true,
			},
		}

		fmt.Println("media_id", payload.MediaId, "file_size", payload.FileSize, "timestamp", payload.Timestamp)

		if err := db.SetMediaOldSize(ctx, params); err != nil {
			fmt.Println(err)
			continue
		}

		media, err := db.GetMedia(ctx, sqlc.GetMediaParams{
			MediaID:   payload.MediaId,
			UserID:    payload.UserId,
			Timestamp: payload.Timestamp,
		})
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err := ws.SendMessageToClient(payload.UserId, fmt.Sprintf("{\"event\":\"DownloadDone\",\"step\":\"%d\"}", media.Step)); err != nil {
			fmt.Println(err)
			continue
		}

		nextPayload := struct {
			MediaID      string `json:"mediaId"`
			UserID       string `json:"userId"`
			Timestamp    int64  `json:"timestamp"`
			CreationDate string `json:"creationDate"`
		}{
			MediaID:      payload.MediaId,
			UserID:       payload.UserId,
			Timestamp:    payload.Timestamp,
			CreationDate: fmt.Sprintf("%d", media.CreationDate),
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

func StartCramerDoneListener(ws *ws.WebSocketManager, db *db.Store, cfg *config.Config) {
	fmt.Println("StartCramerDoneListener")
	ctx := context.Background()
	notification, err := NewAMQPConnection(config.New(), "cramer_done")
	if err != nil {
		panic(err)
	}
	defer notification.Conn.Close()

	next, err := NewAMQPConnection(config.New(), "uploader")
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
			MediaId   string `json:"mediaId"`
			UserId    string `json:"userId"`
			Timestamp int64  `json:"timestamp"`
			FileSize  int64  `json:"fileSize"`
		}

		if err := json.Unmarshal(d.Body, &payload); err != nil {
			continue
		}

		fmt.Println(payload)

		params := sqlc.SetMediaNewSizeParams{
			MediaID:   payload.MediaId,
			UserID:    payload.UserId,
			Timestamp: payload.Timestamp,
			NewSize: sql.NullInt64{
				Int64: payload.FileSize,
				Valid: true,
			},
		}

		fmt.Println("media_id", payload.MediaId, "file_size", payload.FileSize, "timestamp", payload.Timestamp)

		if err := db.SetMediaNewSize(ctx, params); err != nil {
			fmt.Println(err)
			continue
		}

		media, err := db.GetMedia(ctx, sqlc.GetMediaParams{
			MediaID:   payload.MediaId,
			UserID:    payload.UserId,
			Timestamp: payload.Timestamp,
		})
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err := ws.SendMessageToClient(payload.UserId, fmt.Sprintf("{\"event\":\"CramerDone\",\"step\":\"%d\"}", media.Step)); err != nil {
			fmt.Println(err)
			continue
		}

		accessToken, err := GetAccessToken(context.Background(), cfg, db, payload.UserId)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if accessToken == "" {
			fmt.Println("Unauthorized")
			continue
		}

		nextPayload := struct {
			MediaID   string `json:"mediaId"`
			UserID    string `json:"userId"`
			Timestamp int64  `json:"timestamp"`
			Token     string `json:"token"`
			FileName  string `json:"fileName"`
		}{
			MediaID:   payload.MediaId,
			UserID:    payload.UserId,
			Timestamp: payload.Timestamp,
			Token:     accessToken,
			FileName:  media.Filename,
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
