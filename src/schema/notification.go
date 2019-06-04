// Copyright 2019 Axetroy. All rights reserved. MIT license.
package schema

type NotificationPure struct {
	Id      string  `json:"id"`
	Author  string  `json:"author"`
	Title   string  `json:"title"`
	Content string  `json:"content"`
	Read    bool    `json:"read"`    // 用户是否已读
	ReadAt  string  `json:"read_at"` // 用户读取的时间
	Note    *string `json:"note"`
}

type Notification struct {
	NotificationPure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
