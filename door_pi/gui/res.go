package gui

import (
	"fmt"
	"html/template"
)

//go:generate go-bindata -o assets.go -pkg gui res/...

func (gui *GUI) loadRes() {
	gui.w.InjectCSS(string(MustAsset("res/js/vendor/bootstrap/css/bootstrap.min.css")))
	gui.w.InjectCSS(string(MustAsset("res/style.css")))
	gui.w.Eval(string(MustAsset("res/js/vendor/babel.min.js")))
	gui.w.Eval(string(MustAsset("res/js/vendor/preact.min.js")))
	gui.w.Eval(string(MustAsset("res/js/vendor/jquery-3.2.1.min.js")))
	gui.w.Eval(string(MustAsset("res/js/vendor/popper.min.js")))
	gui.w.Eval(string(MustAsset("res/js/vendor/bootstrap/js/bootstrap.bundle.min.js")))

	gui.w.Eval(fmt.Sprintf(`(function(){
		var n=document.createElement('script');
		n.setAttribute('type', 'text/babel');
		n.appendChild(document.createTextNode("%s"));
		document.body.appendChild(n);
	})()`, template.JSEscapeString(string(MustAsset("res/js/app.jsx")))))

	gui.w.Eval(`Babel.transformScriptTags()`)
}