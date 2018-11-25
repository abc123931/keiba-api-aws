test: dep-ensure set-env
	@go test ./gethorsename/
	@go test ./getcourseresult/
	@go test ./line_bot_test/

systest: dep-ensure set-env
	@sam local invoke GetHorseName -e systests/gethorsename_event.json
	@sam local invoke GetHorseId -e systests/gethorseid_event.json
	@sam local invoke GetRaceName -e systests/getracename_event.json
	@sam local invoke GetCourseResult -e systests/getcourseresult_event.json

dep-ensure: 
	@dep ensure

set-env:
	@sh dev_config.sh
