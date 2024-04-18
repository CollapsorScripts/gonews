package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"newsaggr/pkg/aggregator"
	"newsaggr/pkg/api"
	"newsaggr/pkg/config"
	"newsaggr/pkg/database"
	"newsaggr/pkg/database/model"
	"newsaggr/pkg/logger"
	"newsaggr/pkg/rss"
	"os"
	"os/signal"
	"time"
)

// preload - Предзагрузка конфигурации и миграция бд
func preload() error {
	//Инициализация конфига
	err := config.Init()
	if err != nil {
		return err
	}

	//Инициализация базы данных
	if _, err := database.Init(); err != nil {
		return errors.New("Ошибка при инициализации бд: " + err.Error())
	}

	//Миграции
	{
		err := database.GetDB().AutoMigrate(&model.News{})
		if err != nil {
			return errors.New("Ошибка при миграции: " + err.Error())
		}
	}

	return nil
}

func main() {
	if err := preload(); err != nil {
		logger.Error("%v", err)
		return
	}

	//Получаем конфигурацию
	cfgXML, err := rss.GetData()
	if err != nil {
		logger.Error("%v", err)
		return
	}

	//Создаем агрегатор
	aggr := aggregator.New(cfgXML)
	//Запускаем агрегатор
	go aggr.Start()

	//Запуск сервера
	{
		srv := api.New()
		var wait time.Duration
		flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
		flag.Parse()

		go func() {
			if err := srv.Server.ListenAndServe(); err != nil {
				log.Println(err)
			}
		}()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		<-c

		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		err := srv.Server.Shutdown(ctx)
		if err != nil {
			logger.Error("%s", err.Error())
		}
		log.Println("shutting down")
		os.Exit(0)
	}
}
