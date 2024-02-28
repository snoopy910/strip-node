package sequencer

func StartSequencer(httpPort string) {
	keepAlive := make(chan string)

	go startHTTPServer(httpPort)

	<-keepAlive
}
