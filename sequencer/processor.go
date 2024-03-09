package sequencer

import "log"

func ProcessIntent(intentId int64) {
	intent, err := GetIntent(intentId)
	if err != nil {
		log.Println(err)
		return
	}

	if intent.Status != INTENT_STATUS_PROCESSING {
		log.Println("intent already processed")
		return
	}
}
