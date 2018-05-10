package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

/*
	untuk mapping nama fungsi ke fungsinya
	contoh dibawah, nama fungsi "yield"
	fungsinya adalah fungsi anonymous

	layoutFuncs
	FuncMap allows it to call functions defined in the map
*/
var layoutFuncs = template.FuncMap{
	"yield": func() (string, error) {
		return "Testing", fmt.Errorf("yield called inappropriately")
	},
}

/*
	layout ini dibuat untuk merender "yield" dengan
	string pada variabel layoutFuncs, ini perlu dilakukan
	agar kalo proses render dengan fungsi RenderTemplate
	gagal halaman html tetap memberi pesan kegagalan dgn
	'Yield Calles Inapproprientely', dan tidak terjadi
	kondisi halaman html kosong/blank yang akan membi-
	ngungkan user pengguna.
*/
var layout = template.Must(
	template.
		New("layout.html").
		Funcs(layoutFuncs).
		ParseFiles("template/layout.html"),
)

var templates = template.Must(template.New("t").ParseGlob("template/**/*.html"))

var errorTemplate = `
<html>
	<body>
		<h1>Error rendering template %s</h1>
		<p>%s</p>
	</body>
</html>
`

func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}

	data["CurrentUser"] = RequestUser(r)
	data["Flash"] = r.URL.Query().Get("flash")

	/*
		kondisi program memapping string "yield" di template
		untuk d replace atau gantikan dengan w sebagai response-writer
	*/
	funcs := template.FuncMap{
		"yield": func() (template.HTML, error) {
			/*
				buf  : layout yang akan di mapping ke file tujuan
				name : nama file atau path file html tujuan yang akan di mapping
				data : optional, bisa di-isi atau dikosongin
			*/
			buf := bytes.NewBuffer(nil)
			err := templates.ExecuteTemplate(buf, name, data)
			return template.HTML(buf.String()), err
		},
	}

	/*
		layout dilakukan cloning, agar variabel layout aslinya
		tidak ditimpa, karena kita akan terus menerus mengganti
		nilai dari layout nantinya setiap pindah atau berubah layar.

		sehingga dengan cloning, nilai dari variabel layout aslinya
		tidak berubah, yang kita ubah adalah layout cloning-nya
		dan yang di render ke halaman html adalah layout cloningnya.
	*/
	layoutClone, _ := layout.Clone()
	layoutClone.Funcs(funcs)
	err := layoutClone.Execute(w, data)

	if err != nil {
		http.Error(
			w,
			fmt.Sprintf(errorTemplate, name, err),
			http.StatusInternalServerError,
		)
	}
}
