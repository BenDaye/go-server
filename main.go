// Copyright 2019 Axetroy. All rights reserved. MIT license.
package main

import (
	"github.com/axetroy/go-server/src"
	"github.com/axetroy/go-server/src/message_queue"
)

func main() {
	go message_queue.RunMessageQueueConsumer()
	go src.ServerUserClient()
	src.ServerAdminClient()
}
