package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

func (d DbTableRow) getPartners() (rows []DbTableRow, err error) {
	url, err := getIdUrl(d.Id)
	if err != nil {
		return
	}
	lines, err := getHtmlPage(url)
	if err != nil {
		return
	}
	for _, gene := range scrapePartners(lines) {
		if gene != d.Id {
			rows = append(rows, DbTableRow{
				Analyses: map[string]struct{}{
					"sv": struct{}{},
				},
				Id:    gene,
				Class: "gene",
			})
		}
	}
	rows = append(rows, d)
	return
}

func getIdUrl(id string) (url string, err error) {
	if atlasIds[id] == "" {
		err = errors.New(fmt.Sprintf("Id for %s is not registered", id))
		return
	}
	url = fmt.Sprintf("%s/gene-fusions/?id=%s", session.Web.AtlasGO, atlasIds[id])
	return
}

func getHtmlPage(url string) (lines []string, err error) {
	body, err := sendHttpRequest(url)
	if err != nil {
		return
	}
	lines = strings.Split(string(body), "\n")
	return
}

func sendHttpRequest(url string) (body []byte, err error) {
	response, err := http.Get(url)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not reach url %s", url))
		return
	}
	if response.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Request to %s returned %s", url, response.Status))
		return
	}
	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not read response to request to %s", url))
	}
	return
}

func scrapePartners(lines []string) (partners []string) {
	regex := regexp.MustCompile(`</ul>`)
	for _, line := range lines {
		if regex.MatchString(line) {
			pairs := getPairs(line)
			partners = append(partners, pairs...)
			break
		}
	}
	sort.Strings(partners)
	partners = unique(partners)
	return
}

func getPairs(line string) (partners []string) {
	pairs := strings.Split(line, `<li class="border list-group-item">`)
	for _, pair := range pairs {
		partner := extractPartners(pair)
		if len(partner[0]) > 1 {
			partners = append(partners, partner...)
		}
	}
	return
}

func extractPartners(line string) (partners []string) {
	regexps := []*regexp.Regexp{
		regexp.MustCompile("<.*>"),
		regexp.MustCompile(" \\(\\d?\\d?[X,Y]?[p,q]?\\d?\\d?\\.?\\d?\\d?\\d?\\)"),
	}
	for _, regex := range regexps {
		line = regex.ReplaceAllString(line, "")
	}
	partners = strings.Split(line, " ")
	return
}

func unique(slice []string) (uniqueSlice []string) {
	keys := make(map[string]bool)
	for _, entry := range slice {
		if !keys[entry] {
			keys[entry] = true
			uniqueSlice = append(uniqueSlice, entry)
		}
	}
	return
}
