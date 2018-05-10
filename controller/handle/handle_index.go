package handle

import (
	"example-go-web/controller"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Display Home Page
	images, err := controller.GlobalImageStore.FindAll(0)

	if err != nil {
		panic(err)
	}

	controller.RenderTemplate(w, r, "index/home", map[string]interface{}{
		"Images": images,
	})
}
