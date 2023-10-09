package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"github.com/srabraham/strava-oauth-helper/stravaauth"

	strava "github.com/srabraham/swagger-strava-go"
)

var (
	athleteOutFile = flag.String("athlete-out-file", "", "File in which to spew athlete details, or blank to not output such a file")
	// activitiesOutFile = flag.String("activities-out-file", "", "File in which to spew out all activities, or blank to not output such a file")
	timeout = flag.Duration("timeout", 30*time.Minute, "an overall timeout on the program")
	// workoutType       = map[int32]string{
	// 	0:  "Run",
	// 	1:  "Foot race",
	// 	2:  "Long run",
	// 	3:  "Run workout",
	// 	10: "Bike",
	// 	11: "Bike race",
	// 	12: "Bike workout",
	// }
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// authenticate with Strava API
	stravaScopes := []string{"read_all", "activity:read_all", "profile:read_all"}
	oauthCtx, err := stravaauth.GetOAuth2Ctx(ctx, strava.ContextOAuth2, stravaScopes)
	if err != nil {
		log.Fatal(err)
	}
	sClient := strava.NewAPIClient(strava.NewConfiguration())

	// get desired outputs from Strava for selected athleteID
	athlete := *getLoggedInAthleteProfile(sClient, &oauthCtx)
	// activities := *getLoggedInAthleteActivities(sClient, &oauthCtx)

	// when the athleteOutFile is not empty, output as a JSON file
	if *athleteOutFile != "" {
		// marshal athlete to JSON
		athlete_JSON, err := json.MarshalIndent(athlete, "", "    ")
		if err != nil {
			log.Fatal(err)
		}

		// write json to file
		err = os.WriteFile(*athleteOutFile, athlete_JSON, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	// read that JSON file as bytes, json.Unmarshal then type as a strava.DetailedAthlete
	fileData, err := os.ReadFile(*athleteOutFile)
	if err != nil {
		log.Fatal(err)
	}
	var newAthlete strava.DetailedAthlete
	err = json.Unmarshal(fileData, &newAthlete)
	if err != nil {
		log.Fatal(err)
	}

	// check that the data directly from the Strava API matches the unmarshalled JSON file
	log.Printf("diff is: %v", cmp.Diff(newAthlete, athlete))
	// fmt.Println("Original API variable and unmarshaled JSON file match")
	// }

}

func getLoggedInAthleteProfile(sClient *strava.APIClient, oauthCtx *context.Context) *strava.DetailedAthlete {
	athlete, _, err := sClient.AthletesApi.GetLoggedInAthlete(*oauthCtx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Got athlete:")
	spew.Dump(athlete)
	return &athlete
}
