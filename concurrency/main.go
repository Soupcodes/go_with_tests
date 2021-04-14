package main

type WebsiteChecker func(string) bool

type result struct {
	string
	bool
}

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)
	resultsChan := make(chan result)

	// _ is the index of the slice, url is the element
	for _, url := range urls {
		go func(url string) {
			resultsChan <- result{url, wc(url)}
		}(url)
	}

	// for r := range would loop forever here because there's no break condition, or done channel to inform it to finish
	for i := 0; i < len(urls); i++ {
		r := <-resultsChan
		results[r.string] = r.bool
	}

	return results
}
