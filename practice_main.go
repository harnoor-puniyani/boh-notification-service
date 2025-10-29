package newmain

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func homeHandler(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/" {
		http.Error(w,"not the root path",http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Hello from the go notification service")
}


type Shapes interface{
	Area() float64
}

type circle struct {
radius int
}

func (c circle) Area() float64{
return 2*3.14*float64(c.radius)
}

func findArea(s Shapes) float64{
	return s.Area()
}

func userHandler(w http.ResponseWriter, r *http.Request){
	pathParts := strings.Split(r.URL.Path,"/")

	if len(pathParts) <3 || pathParts[2] =="" {
		http.Error(w , "Bad reqeust: missing the userID",http.StatusBadRequest)
	}

	userId := pathParts[2]

	log.Printf("Received request to notify the user: %s",userId)

	fmt.Println("received userID ", userId,"okay again userId",userId)

	fmt.Fprintln(w,"received userID ", userId,"okay again userId",userId)

}

func healthHandler(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w,`{"status":"okay"}`)
	w.Header().Set("Content-Type","applicaction/json")
}

func newmain(){


	// a := circle{radius: 2}

	// fmt.Printf("%f",findArea(a))




	http.HandleFunc("/",homeHandler)
	http.HandleFunc("/health",healthHandler)
	http.HandleFunc("/notify/",userHandler)

	port := ":5000"

	fmt.Printf("Notification service starting on port %s\n", port)

	log.Fatal(http.ListenAndServe(port,nil))
}