package ttask

import (
	"context"
	"log"
	// "fmt"

	"github.com/redis/go-redis/v9"
)

type RediStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RediStorage {
	return &RediStorage{
		client: client,
	}
}

func (s *RediStorage) storeCheckpoint(taskId string, msgId string, cpId string, msg []byte) error {
	ctx := context.Background()

	_, err := s.client.Set(ctx, "ttcp."+taskId+"."+msgId, msg, 0).Result()
	if err != nil {
		return err
	}

	_, err = s.client.SAdd(ctx, "ttcpidx."+taskId+"."+cpId, msgId).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *RediStorage) clearCheckpoint(taskId string, msgId string, cpId string) error {
	ctx := context.Background()

	_, err := s.client.Del(ctx, "ttcp."+taskId+"."+msgId).Result()
	if err != nil {
		return err
	}

	_, err = s.client.SRem(ctx, "ttcpidx."+taskId+"."+cpId, msgId).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *RediStorage) getCheckpointMessages(taskId string, cpId string) ([][]byte, error) {
	ctx := context.Background()

	msgIds, err := s.client.SMembers(ctx, "ttcpidx."+taskId+"."+cpId).Result()
	if err != nil {
		return nil, err
	}

	if len(msgIds) != 0 {
		log.Printf("recovered %v for checkpoint %s\n", len(msgIds), cpId)
	}

	msgs := [][]byte{}
	for i := range msgIds {
		k := "ttcp." + taskId + "." + msgIds[i]
		msg, err := s.client.Get(ctx, k).Result()
		if err != nil {
			return nil, err
		}

		msgs = append(msgs, []byte(msg))
	}

	return msgs, nil
}
