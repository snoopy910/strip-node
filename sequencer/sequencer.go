package sequencer

import "log"

func StartSequencer(httpPort string) {
	keepAlive := make(chan string)

	intents, err := GetIntents(INTENT_STATUS_PROCESSING)
	if err != nil {
		log.Fatal(err)
	}

	for _, intent := range intents {
		go ProcessIntent(intent.ID)
	}

	go startHTTPServer(httpPort)

	<-keepAlive
}
