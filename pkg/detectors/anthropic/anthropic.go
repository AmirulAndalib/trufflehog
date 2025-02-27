package anthropic

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	regexp "github.com/wasilibs/go-re2"

	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/detectors"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/detectorspb"
)

type Scanner struct {
	client *http.Client
}

// Ensure the Scanner satisfies the interface at compile time.
var _ detectors.Detector = (*Scanner)(nil)

var (
	defaultClient = common.SaneHttpClient()
	// Make sure that your group is surrounded in boundary characters such as below to reduce false positives.
	keyPat = regexp.MustCompile(`\b(sk-ant-api03-[\w\-]{93}AA)\b`)
)

// Keywords are used for efficiently pre-filtering chunks.
// Use identifiers in the secret preferably, or the provider name.
func (s Scanner) Keywords() []string {
	return []string{"sk-ant-api03"}
}

// FromData will find and optionally verify Anthropic secrets in a given set of bytes.
func (s Scanner) FromData(ctx context.Context, verify bool, data []byte) (results []detectors.Result, err error) {
	dataStr := string(data)

	matches := keyPat.FindAllStringSubmatch(dataStr, -1)

	for _, match := range matches {
		resMatch := strings.TrimSpace(match[1])

		s1 := detectors.Result{
			DetectorType: detectorspb.DetectorType_Anthropic,
			Raw:          []byte(resMatch),
		}

		if verify {
			client := s.client
			if client == nil {
				client = defaultClient
			}
			isVerified, err := verifyToken(ctx, client, resMatch)
			s1.Verified = isVerified
			s1.SetVerificationError(err, resMatch)
			if s1.Verified {
				s1.AnalysisInfo = map[string]string{
					"key": resMatch,
				}
			}
		}

		results = append(results, s1)
	}

	return results, nil
}

func verifyToken(ctx context.Context, client *http.Client, apiKey string) (bool, error) {
	// https://docs.anthropic.com/en/api/models-list
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.anthropic.com/v1/models", http.NoBody)
	if err != nil {
		return false, nil
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")

	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		return true, nil

	case http.StatusNotFound, http.StatusUnauthorized:
		// 404 is returned if api key is disabled or not found
		return false, nil

	default:
		return false, fmt.Errorf("unexpected HTTP response status %d", res.StatusCode)
	}
}

func (s Scanner) Type() detectorspb.DetectorType {
	return detectorspb.DetectorType_Anthropic
}

func (s Scanner) Description() string {
	return "Anthropic is an AI research company. The API keys can be used to access their AI models and services."
}
