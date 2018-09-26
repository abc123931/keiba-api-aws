test: dep-ensure set-env
	@go test ./gethorsename/

systest: dep-ensure set-env
	@sam local invoke GetHorseName -e systests/gethorsename_event.json
	@sam local invoke GetHorseId -e systests/gethorseid_event.json
	@sam local invoke GetRaceName -e systests/getracename_event.json

dep-ensure: 
	@dep ensure

set-env:
	@sh dev_config.sh
