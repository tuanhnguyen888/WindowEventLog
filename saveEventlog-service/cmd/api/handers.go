package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"log"
	"net/http"
	entity "saveEventLog/entity"
	"strconv"
)

type SaveLogHandler struct {
	redisClient *redis.Client
	pgDB        *gorm.DB
}

func NewSaveLogHandler(redisClient *redis.Client, pgDb *gorm.DB) SaveLogHandler {
	return SaveLogHandler{
		redisClient: redisClient,
		pgDB:        pgDb,
	}
}

func (s *SaveLogHandler) getLogsFromRedis() ([]entity.EventLog, error) {
	keys, err := s.redisClient.Keys("*").Result()
	if err != nil {
		return nil, err
	}
	var eventLogs []entity.EventLog

	for _, key := range keys {
		eventLog, err := s.getLog(key)
		if err != nil {
			log.Println(err)
		}
		eventLogs = append(eventLogs, eventLog)
	}
	log.Println(keys, "2")

	return eventLogs, nil
}

func (s *SaveLogHandler) getLog(key string) (entity.EventLog, error) {
	logE := entity.EventLog{}

	value, err := s.redisClient.Get(key).Result()
	if err != nil {
		return entity.EventLog{}, err
	}
	log.Println(value, "1")

	err = json.Unmarshal([]byte(value), &logE)
	if err != nil {
		return entity.EventLog{}, err
	}

	log.Println(logE, "1")
	return logE, nil
}

func (s *SaveLogHandler) saveLogByPg(ctx *gin.Context) {
	eventLogs, err := s.getLogsFromRedis()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("%#v", eventLogs)

	for _, eventLog := range eventLogs {
		err = s.pgDB.Create(&eventLog).Error
		if err != nil {
			log.Println(err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Add Success!",
	})
}

func (s *SaveLogHandler) showLogFromPG(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page == 0 {
		page = 1
	}
	perPage := 10
	offset := (page - 1) * perPage

	// ----search
	var eventLog []entity.EventLog
	s.pgDB.Order("created").Offset(offset).Limit(perPage).Find(&eventLog)
	if len(eventLog) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": eventLog,
	})
}
