// Turorial: https://www.zenrows.com/blog/web-scraping-golang#set-up-go-project
// Scraping from https://www.atozflowers.com/flower/

package main 
 
import ( 
	"encoding/csv" 
	"github.com/gocolly/colly" 
	"log" 
	"os" 
	"fmt"
) 
 
// defining a data structure to store the scraped data 
type PokemonProduct struct { 
	url, image, name, price string 
} 

type FlowerList struct {
	url string
	name string
}

type FlowerPage struct {
	name string
	imageURL string
}
 
// it verifies if a string is present in a slice 
func contains(s []string, str string) bool { 
	for _, v := range s { 
		if v == str { 
			return true 
		} 
	} 
 
	return false 
} 
 
func main() { 
	scrapeAllFlowers()
}

func scrapeAllFlowers(){
	// initializing the slice of structs that will contain the scraped data 	
	// var pokemonProducts []PokemonProduct 
	var flowers []FlowerList 
 
	// initializing the list of pages to scrape with an empty slice 
	var pagesToScrape []string 
 
	// the first pagination URL to scrape 
	pageToScrape := "https://www.atozflowers.com/flower/" 
 
	// initializing the list of pages discovered with a pageToScrape 
	pagesDiscovered := []string{ pageToScrape } 
 
	// current iteration 
	i := 1 
	// max pages to scrape 
	// total number of flowers / flowers on a page
	limit := 2
 
	// initializing a Colly instance 
	c := colly.NewCollector() 
	// setting a valid User-Agent header 
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36" 
 
	// iterating over the list of pagination links to implement the crawling logic 
	c.OnHTML("a.page", func(e *colly.HTMLElement) { 
		// discovering a new page 
		newPaginationLink := e.Attr("href") 
 
		// if the page discovered is new 
		if !contains(pagesToScrape, newPaginationLink) { 
			// if the page discovered should be scraped 
			if !contains(pagesDiscovered, newPaginationLink) { 
				pagesToScrape = append(pagesToScrape, newPaginationLink) 
			} 
			pagesDiscovered = append(pagesDiscovered, newPaginationLink) 
		} 
	}) 
 
	// scraping the product data 
	c.OnHTML("li.col-smx-12", func(e *colly.HTMLElement) { 
		flowerList := FlowerList{}

		flowerList.url =  e.ChildAttr("a", "href")
		flowerList.name = e.ChildText("h4")

		scrapeFlowerImage(e.ChildText("h4"), e.ChildAttr("a", "href"))

		flowers = append(flowers, flowerList)
	}) 
 
	c.OnScraped(func(response *colly.Response) { 
		// until there is still a page to scrape 
		if len(pagesToScrape) != 0 && i < limit { 
			// getting the current page to scrape and removing it from the list 
			pageToScrape = pagesToScrape[0] 
			pagesToScrape = pagesToScrape[1:] 
 
			// incrementing the iteration counter 
			i++ 
 
			// visiting a new page 
			c.Visit(pageToScrape) 
		} 
	}) 
 
	// visiting the first page 
	c.Visit(pageToScrape) 
}

func scrapeFlowerImage(flowerName string, url string){
	var flowerImageURLS []FlowerPage

 
	// initializing a Colly instance 
	c := colly.NewCollector() 
	// setting a valid User-Agent header 
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36" 
 

	// scraping the product data 
	c.OnHTML("li.col-3", func(e *colly.HTMLElement) { 
		flowerPage := FlowerPage{}

		flowerPage.name = flowerName
		flowerPage.imageURL =  e.ChildAttr("a", "href")


		fmt.Println("img url: ", e.ChildAttr("a", "href"), "\n")
		flowerImageURLS = append(flowerImageURLS, flowerPage)
	}) 
  
	// visiting the first page 
	c.Visit(url) 

	// opening the CSV file 
	file, err := os.Create(("flower_image_urls/" + flowerName + ".csv")) 
	if err != nil { 
		log.Fatalln("Failed to create output CSV file", err) 
	} 
	defer file.Close() 

	// initializing a file writer 
	writer := csv.NewWriter(file) 

	// defining the CSV headers 
	headers := []string{ 
		"name", 
		"ImageURL", 
	} 
	// writing the column headers 
	writer.Write(headers) 

	// adding each Pokemon product to the CSV output file 
	for _, image := range flowerImageURLS { 
		// converting a PokemonProduct to an array of strings 
		record := []string{ 
			image.name,
			image.imageURL, 
		} 

		// writing a new CSV record 
		writer.Write(record) 
	} 
	defer writer.Flush() 
}
