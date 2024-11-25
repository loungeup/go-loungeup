package esutil

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/loungeup/go-loungeup/pkg/platform"
)

const (
	globalIndexSuffix   = "global"
	wildcardIndexSuffix = "*"
)

type Indices struct {
	Bookings string
	Guests   string
}

func (i *Indices) Strings() []string {
	result := []string{}

	if i.Bookings != "" {
		result = append(result, i.Bookings)
	}

	if i.Guests != "" {
		result = append(result, i.Guests)
	}

	return result
}

type indicesMaker struct {
	platform platform.Platform
}

func MakeIndices(platform platform.Platform) *indicesMaker { return &indicesMaker{platform} }

func (m *indicesMaker) At(t time.Time) *Indices {
	switch m.platform {
	case platform.Development, platform.Studio:
		return makeIndices(makeIndexPrefix(m.platform), globalIndexSuffix)
	default:
		return makeIndices(makeIndexPrefix(m.platform), formatIndexTime(t))
	}
}

func (m *indicesMaker) Wildcard() *Indices {
	return makeIndices(makeIndexPrefix(m.platform), wildcardIndexSuffix)
}

func ParseResponseBody(response *esapi.Response) (string, error) {
	if response.IsError() {
		return "", errors.New(fmt.Sprintf("invalid status code: %d", response.StatusCode))
	}

	defer response.Body.Close()

	bodyBuilder := &strings.Builder{}
	if _, err := io.Copy(bodyBuilder, response.Body); err != nil {
		return "", fmt.Errorf("could not read response body: %w", err)
	}

	return bodyBuilder.String(), nil
}

func formatIndexTime(t time.Time) string {
	const thresholdYearForGlobal = 2023
	if t.Year() < thresholdYearForGlobal {
		return globalIndexSuffix
	}

	return t.Format("2006-01")
}

func makeIndices(prefix, suffix string) *Indices {
	return &Indices{
		Bookings: strings.Join([]string{prefix, "guestbookings", suffix}, "-"),
		Guests:   strings.Join([]string{prefix, "guestcards", suffix}, "-"),
	}
}

func makeIndexPrefix(p platform.Platform) string {
	switch p {
	case platform.Development:
		return "development"
	case platform.Studio:
		return "studio"
	default:
		return "production"
	}
}
