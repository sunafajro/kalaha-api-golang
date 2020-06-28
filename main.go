package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type File struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
type Text struct {
	Cv []string `json:"cv"`
	Ru []string `json:"ru"`
}

type Audio struct {
	File File `json:"file"`
	Text Text `json:"text"`
}

type Part struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Chapters []int  `json:"chapters"`
}

type Chapter struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Pages []File `json:"pages"`
}

type Book struct {
	Id       int       `json:"id"`
	Title    string    `json:"title"`
	Cover    File      `json:"cover"`
	Back     File      `json:"back"`
	Parts    []Part    `json:"parts"`
	Chapters []Chapter `json:"chapters"`
}

type Books struct {
	Books []Book `json:"books"`
}

func handleGetRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Simple GoLang API")
}

func handleGetBooks(w http.ResponseWriter, r *http.Request) {
	if len(books.Books) > 0 {
		var emptyChapters []Chapter
		var emptyParts []Part
		var booksSimple []Book
		for i := 0; i < len(books.Books); i++ {
			book := books.Books[i]
			book.Chapters = emptyChapters
			book.Parts = emptyParts
			booksSimple = append(booksSimple, book)
		}
		jsonData, encodeError := json.Marshal(booksSimple)
		if encodeError != nil {
			log.Println(encodeError)
		}
		fmt.Fprint(w, string(jsonData))
		return
	}
	fmt.Fprint(w, "No books found.")
}

func handleGetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId, bookIdParseError := strconv.Atoi(vars["bookId"])
	if bookIdParseError != nil {
		log.Println(bookIdParseError)
	}
	if bookId > 0 {
		var book Book
		for i := 0; i < len(books.Books); i++ {
			if bookId == books.Books[i].Id {
				book = books.Books[i]
			}
		}
		if book.Id > 0 {
			jsonData, encodeError := json.Marshal(book)
			if encodeError != nil {
				log.Println(encodeError)
			}
			fmt.Fprint(w, string(jsonData))
			return
		}
	}
	fmt.Fprint(w, "The requested book not found.")
}

func handleGetBookChapter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId, bookIdParseError := strconv.Atoi(vars["bookId"])
	if bookIdParseError != nil {
		log.Println(bookIdParseError)
	}
	if bookId > 0 {
		chapterId, chapterIdParseError := strconv.Atoi(vars["chapterId"])
		if chapterIdParseError != nil {
			log.Println(chapterIdParseError)
		}
		if chapterId > 0 {
			var book Book
			for i := 0; i < len(books.Books); i++ {
				if bookId == books.Books[i].Id {
					book = books.Books[i]
				}
			}
			if book.Id > 0 {
				var chapter Chapter
				for i := 0; i < len(book.Chapters); i++ {
					if chapterId == book.Chapters[i].Id {
						chapter = book.Chapters[i]
					}
				}
				if chapter.Id > 0 {
					jsonData, encodeError := json.Marshal(chapter)
					if encodeError != nil {
						log.Println(encodeError)
					}
					fmt.Fprint(w, string(jsonData))
					return
				}
			}
		}
	}
	fmt.Fprint(w, "The requested chapter not found.")
}

func handleGetBookFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId, bookIdParseError := strconv.Atoi(vars["bookId"])
	if bookIdParseError != nil {
		log.Println(bookIdParseError)
	}
	if bookId > 0 {
		if vars["fileName"] != "" {
			if vars["fileType"] != "" {
				filePath := fileBasePath + vars["bookId"] + "/" + vars["fileName"] + "." + vars["fileType"]
				http.ServeFile(w, r, filePath)
				return
			}
		}
	}
	fmt.Fprint(w, "The requested file not found.")
}

var books Books
var fileBasePath string

func main() {
	fileBasePath = "./data/content/kala-ha_0"

	data, readError := ioutil.ReadFile("./data/books.json")
	if readError != nil {
		panic(readError)
	}
	parseError := json.Unmarshal(data, &books)
	if parseError != nil {
		panic(parseError)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handleGetRoot)
	r.HandleFunc("/api/books", handleGetBooks)
	r.HandleFunc("/api/books/{bookId}", handleGetBook)
	r.HandleFunc("/api/books/{bookId}/{chapterId}", handleGetBookChapter)
	r.HandleFunc("/api/books/{bookId}/file/{fileName}/{fileType}", handleGetBookFile)
	http.Handle("/", r)
	fmt.Println(http.ListenAndServe(":9999", nil))
}
