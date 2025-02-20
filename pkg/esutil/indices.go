package esutil

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/loungeup/go-loungeup"
)

const (
	globalIndexSuffix   = "global"
	wildcardIndexSuffix = "alias"
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

type IndicesMaker struct {
	platform loungeup.Platform
}

func MakeIndices(platform loungeup.Platform) *IndicesMaker { return &IndicesMaker{platform} }

func (m *IndicesMaker) At(t time.Time) *Indices {
	prefix := makeIndexPrefix(m.platform)

	if m.platform == loungeup.PlatformDevelopment || m.platform == loungeup.PlatformStudio {
		return makeIndices(prefix, globalIndexSuffix)
	}

	return &Indices{
		Bookings: makeBookingIndex(prefix, formatBookingIndexTime(t)),
		Guests:   makeGuestIndex(prefix, formatGuestIndexTime(t)),
	}
}

func (m *IndicesMaker) Wildcard() *Indices {
	return makeIndices(makeIndexPrefix(m.platform), wildcardIndexSuffix)
}

func ParseResponseBody(response *esapi.Response) (string, error) {
	if response.IsError() {
		return "", fmt.Errorf("invalid status code: %d", response.StatusCode)
	}

	defer response.Body.Close()

	bodyBuilder := &strings.Builder{}
	if _, err := io.Copy(bodyBuilder, response.Body); err != nil {
		return "", fmt.Errorf("could not read response body: %w", err)
	}

	return bodyBuilder.String(), nil
}

func formatBookingIndexTime(t time.Time) string {
	year := t.Year()

	if year < 2020 || year > 2100 {
		return globalIndexSuffix
	}

	if year < 2022 { //nolint:mnd
		return strconv.Itoa(year)
	}

	return t.Format("2006-01")
}

func formatGuestIndexTime(t time.Time) string {
	if t.Year() < 2023 { //nolint:mnd
		return globalIndexSuffix
	}

	return t.Format("2006-01")
}

func makeIndices(prefix, suffix string) *Indices {
	return &Indices{
		Bookings: makeBookingIndex(prefix, suffix),
		Guests:   makeGuestIndex(prefix, suffix),
	}
}

func makeBookingIndex(prefix, suffix string) string {
	return strings.Join([]string{prefix, "guestbookings", suffix}, "-")
}

func makeGuestIndex(prefix, suffix string) string {
	return strings.Join([]string{prefix, "guestcards", suffix}, "-")
}

func makeIndexPrefix(p loungeup.Platform) string {
	switch p {
	case loungeup.PlatformDevelopment:
		return "development"
	case loungeup.PlatformStudio:
		return "studio"
	default:
		return "production"
	}
}
