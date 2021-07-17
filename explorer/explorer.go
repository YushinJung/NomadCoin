package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/YushinJung/NomadCoin/blockchain"
)

var templates *template.Template

const (
	port        string = ":4000"
	templateDir string = "explorer/templates/"
)

//parsing 파일을 하기 전에 function 이 호출되기 전에 template(pages, partials)을 일단 loading 해 놓자
//이를 하기 위해 template에 대한 정보를 갖고 있는 variable을 선언하자.

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", data)
	//templ := template.Must(template.ParseFiles("templates/pages/home.gohtml"))
	//templ.Execute(rw, data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.GetBlockchain().AddBlock(data)
		//redirection
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start(port int) {
	handler := http.NewServeMux()
	//함수로 부르지 않고 저장되어 있는 template을 가져올 것이다
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	// loaded all file in "pages folder"
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	// templates 는 template Object이고, 다른 template를 load 할 수 있음.
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	fmt.Printf("listening on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler)) // log.Fatal will stop if there is error from input
}
