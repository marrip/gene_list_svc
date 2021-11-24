package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

func sendHttpRequest(url string) (body []byte, err error) {
	response, err := http.Get(url)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not reach url %s", url))
		return
	}
	if response.Status != "200 OK" {
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

func getHtmlPage(url string) (lines []string, err error) {
	body, err := sendHttpRequest(url)
	if err != nil {
		return
	}
	lines = strings.Split(string(body), "\n")
	return
}

func (s Session) getAtlasIndexUrl(firstLetter string) string {
	return fmt.Sprintf("%s/Indexbyalpha/idxa_%s.html", s.AtlasRootUrl, firstLetter)
}

func (s Session) findHref(gene string, lines []string) (url string, err error) {
	geneRegex := regexp.MustCompile(fmt.Sprintf(">%s<", gene))
	for _, line := range lines {
		if geneRegex.MatchString(line) {
			start := strings.Index(line, "HREF")
			end := strings.Index(line, ".html")
			url = fmt.Sprintf("%s%s", s.AtlasRootUrl, line[(start+7):(end+5)])
			log.Printf("Found %s for gene %s", url, gene)
			return
		}
	}
	err = errors.New(fmt.Sprintf("Could not find link for gene %s", gene))
	return
}

func (s Session) findGeneUrl(gene string) (url string, err error) {
	lines, err := getHtmlPage(s.getAtlasIndexUrl(gene[0:1]))
	if err != nil {
		return
	}
	url, err = s.findHref(gene, lines)
	return
}

func extractFusions(line string) (fusions []string, err error) {
	regexps := []*regexp.Regexp{
		regexp.MustCompile("<.*>"),
		regexp.MustCompile(" \\(\\d?\\d?[X,Y]?[p,q]?\\d?\\d?\\.?\\d?\\d?\\d?\\)"),
	}
	for _, regex := range regexps {
		line = regex.ReplaceAllString(line, "")
	}
	fusions = strings.Split(line, "::")
	return
}

func unique(slice []string) (uniqueSlice []string) {
	keys := make(map[string]bool)
	for _, entry := range slice {
		if keys[entry] {
			keys[entry] = true
			uniqueSlice = append(uniqueSlice, entry)
		}
	}
	return
}

func getFusionPairs(line string) (fusions []string, err error) {
	pairs := strings.Split(line, "</TD>")
	for _, pair := range pairs {
		var fusion []string
		fusion, err = extractFusions(pair)
		if err != nil {
			return
		}
		if len(fusion[0]) > 1 {
			fusions = append(fusions, fusion...)
		}
	}
	return
}

func getFusionGenePairs(lines []string) (fusions []string, err error) {
	startRegex := regexp.MustCompile(">Fusion genes<")
	endRegex := regexp.MustCompile(">DNA/RNA<")
	fusionRegex := regexp.MustCompile("::")
	var enable bool
	for _, line := range lines {
		if startRegex.MatchString(line) {
			enable = true
		} else if endRegex.MatchString(line) {
			break
		}
		if enable && fusionRegex.MatchString(line) {
			var pairs []string
			pairs, err = getFusionPairs(line)
			if err != nil {
				return
			}
			fusions = append(fusions, pairs...)
		}
	}
	sort.Strings(fusions)
	fusions = unique(fusions)
	return
}

func (s Session) getFusionGenes(entity Entity) (genes []string, err error) {
	url, err := s.findGeneUrl(entity.Id)
	if err != nil {
		return
	}
	lines, err := getHtmlPage(url)
	if err != nil {
		return
	}
	genes, err = getFusionGenePairs(lines)
	return
}
