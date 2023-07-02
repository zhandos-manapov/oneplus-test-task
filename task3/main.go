package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
)

var data [50][7]string
var categories [50]string
var influencers *os.File
var csvwriter *csv.Writer

func main() {
	url := "https://hypeauditor.com/top-instagram-all-russia"
	err := parsePage(&url)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func parsePage(url *string) error {
	// Get request to url
	resp, err := http.Get(*url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("http request error")
	}

	// Parse html response
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return errors.New("parsing error")
	}

	// Create empty csv file
	influencers, err = os.Create("influencers.csv")
	if err != nil {
		return err
	}
	defer influencers.Close()

	// Create csv writer
	csvwriter = csv.NewWriter(influencers)

	// Recursively traverse html tree
	rowIndex := 0
	traverseHtmlTree(doc, &rowIndex)

	// Write found data from 2D array to csv file
	err = csvwriter.Write([]string{"Rank", "Instagram handle", "Influencer", "Followers", "Country", "Authentic", "Engagement", "Category"})
	if err != nil {
		return err
	}
	for i := 0; i < 50; i++ {
		category := categories[i]
		if len(category) > 2 {
			category = category[:len(category)-2]
		}
		temp := append(data[i][:], category)
		err = csvwriter.Write(temp)
		if err != nil {
			return err
		}
	}
	csvwriter.Flush()
	return nil
}

func traverseHtmlTree(node *html.Node, rowIndex *int) {
	if node == nil {
		return
	}

	// Check if the current node is the table
	if node.Type == html.ElementNode && node.Data == "div" {
		for _, attr := range node.Attr {
			if attr.Key == "class" && attr.Val == "table" {
				child := node.FirstChild
				child = child.NextSibling

				// Handle all the rows of the table
				for ; child != nil; child = child.NextSibling {
					handleRow(child, rowIndex)
				}
				return
			}
		}
	}

	// Continue searching for a table
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		traverseHtmlTree(c, rowIndex)
	}
}

func handleRow(node *html.Node, rowIndex *int) {
	if node.Type != html.ElementNode {
		return
	}
	node = node.FirstChild
	colIndex := 0

	// Search for text fields in current row
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		traverseRow(child, rowIndex, &colIndex)
	}
	*rowIndex += 1
}

func traverseRow(node *html.Node, rowIndex *int, colIndex *int) {
	// Check if current node is TextNode
	if node.Type == html.TextNode {
		firstChar := node.Data[0]
		if (firstChar >= 97 && firstChar <= 122) || (firstChar >= 65 && firstChar <= 90) || (firstChar >= 48 && firstChar <= 57) {
			for _, attr := range node.Parent.Attr {
				if attr.Key == "class" {
					// Ignore delta field
					if attr.Val == "ml-2" {
						return   
					} else if attr.Val == "tag__content ellipsis" {
						// Write multiple categories as one string 
						categories[*rowIndex] += node.Data + ", "
						return
					}
				}
			}
			// Ignore button cells
			if *colIndex > 6 {
				return
			}
			// Write found data to 2D array
			data[*rowIndex][*colIndex] = node.Data
			*colIndex += 1
		}
	}

	// Continue searching for text fields
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		traverseRow(child, rowIndex, colIndex)
	}
}
