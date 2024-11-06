package gcs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	video "cloud.google.com/go/videointelligence/apiv1"
	videopb "cloud.google.com/go/videointelligence/apiv1/videointelligencepb"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Attraction struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}


func TransformVideoData(url string) (Attraction, error) {
	fileBytes, err := DownloadVideoFile(url)

	var attraction Attraction

	if err != nil {
		return attraction, fmt.Errorf("could not download video: %w", err)
	}
	extractedText, err := ExtractTextInVideo(fileBytes)

	if err != nil {
		return attraction, fmt.Errorf("something went wrong when using gemini to analyze the video file %w", err)
	}

	// fmt.Println("extractedText: " + extractedText)

	attraction, err = getAttractionFromText(extractedText)

	if err != nil {
		return attraction, fmt.Errorf("something went wrong when using gemini to extract text from passage %w", err)
	}

	return attraction, nil
}

// gets the name and location of the attraction mentioned in the video
func getAttractionFromText(extractedText string) (Attraction, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))

	var attraction Attraction

	if err != nil {
		return attraction, err
	}

	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	prompt := "can you give me the name and location of the attraction from this passage, ignore everything not related to the review, format the answer to be in the form name | location exactly" + extractedText
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return attraction, err
	}

	fmt.Println("response:")
	for _, cand := range resp.Candidates {
		fmt.Println(cand.Content.Parts[0].(genai.Text))
	}

	result, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)

	if !ok {
		return attraction, errors.New("response returned was not text")
	}

	resultStr := string(result)
	fmt.Println("result str: " + resultStr)
	pieces := strings.Split(resultStr, "|")

	if !strings.Contains(resultStr, "|") || len(pieces) < 2 {
		return attraction, errors.New("missing Data from passage, could not get name/location or both")
	}

	// removing leading/trailing spaces
	attraction.Name = strings.TrimSpace(pieces[0])
	attraction.Location = strings.TrimSpace(pieces[1])

	return attraction, nil
}

func ExtractTextInVideo(fileBytes []byte) (string, error) {
	ctx := context.Background()

	// creates a client
	client, err := video.NewClient(ctx)

	if err != nil {
		return "", fmt.Errorf("video.NewClient: %w", err)
	}
	defer client.Close()

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		InputContent: fileBytes,
		Features: []videopb.Feature{
			videopb.Feature_TEXT_DETECTION,
		},
	})

	if err != nil {
		return "", fmt.Errorf("annotateVideo: %w", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return "", fmt.Errorf("wait: %w", err)
	}

	// Only one video was processed, so get the first result.
	result := resp.GetAnnotationResults()[0]

	// sort by time in the video
	sort.Slice(result.TextAnnotations, func(i, j int) bool {
		t1 := result.TextAnnotations[i].GetSegments()[0].Segment.StartTimeOffset.AsDuration().Seconds()
		t2 := result.TextAnnotations[j].GetSegments()[0].Segment.StartTimeOffset.AsDuration().Seconds()
		return t1 < t2
	})

	resultText := ""
	confidence_threshold, _ := strconv.ParseFloat(os.Getenv("GEMINI_API_VIDEO_CONFIDENCE_THRESHOLD"), 32)

	for _, annotation := range result.TextAnnotations {
		// Get the first text segment.
		segment := annotation.GetSegments()[0]

		// don't take if confidence is too low, configurable in environment variables file
		if segment.GetConfidence() < float32(confidence_threshold) {
			continue
		}

		resultText += annotation.GetText() + " "
	}

	return strings.TrimSpace(resultText), nil
}

func DownloadVideoFile(url string) ([]byte, error) {
	resp, rerr := http.Get(url)
	if rerr != nil {
		return nil, rerr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("couldn't get 200 response, statusCode: " + string(resp.StatusCode))
	}

	bytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.New("couldnt read bytes from the video\n: " + err.Error())
	}

	return bytes, nil
}

func GenerateTripName(attraction Attraction) (string, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))

	trip_name := ""
	
	if err != nil {
		return trip_name, err
	}

	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	prompt := "can you give me just the name of the country or city based on the following location: " + attraction.Location
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return trip_name, err
	}

	for _, cand := range resp.Candidates {
		fmt.Println(cand.Content.Parts[0].(genai.Text))
	}

	result, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)

	if !ok {
		return trip_name, errors.New("response returned was not text")
	}

	trip_name = string(result)

	return trip_name, nil
}