// Copyright 2019 Axetroy. All rights reserved. MIT license.
package news

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/middleware"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type UpdateParams struct {
	Title   *string           `json:"title"`
	Content *string           `json:"content"`
	Type    *model.NewsType   `json:"type"`
	Tags    *[]string         `json:"tags"`
	Status  *model.NewsStatus `json:"status"`
}

func Update(context controller.Context, newsId string, input UpdateParams) (res schema.Response) {
	var (
		err          error
		data         schema.News
		tx           *gorm.DB
		shouldUpdate bool
	)

	defer func() {

		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		if tx != nil {
			if err != nil || !shouldUpdate {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	tx = database.Db.Begin()

	adminInfo := model.Admin{Id: context.Uid}

	// 判断管理员是否存在
	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
			return
		}
	}

	newsInfo := model.News{
		Id: newsId,
	}

	if err = tx.First(&newsInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NewsNotExist
			return
		}
		return
	}

	if input.Title != nil {
		shouldUpdate = true
		newsInfo.Title = *input.Title
	}

	if input.Content != nil {
		shouldUpdate = true
		newsInfo.Content = *input.Content
	}

	if input.Type != nil {
		shouldUpdate = true
		newsInfo.Type = *input.Type
	}

	if input.Status != nil {
		shouldUpdate = true
		newsInfo.Status = *input.Status
	}

	if input.Tags != nil {
		shouldUpdate = true
		newsInfo.Tags = *input.Tags
	}

	if err = tx.Save(&newsInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NewsNotExist
			return
		}
		return
	}

	if err = mapstructure.Decode(newsInfo, &data.NewsPure); err != nil {
		return
	}

	data.CreatedAt = newsInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = newsInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func UpdateRouter(context *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input UpdateParams
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	id := context.Param("news_id")

	if err = context.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = Update(controller.Context{
		Uid: context.GetString(middleware.ContextUidField),
	}, id, input)
}
