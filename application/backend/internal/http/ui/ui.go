package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
	v8 "rogchap.com/v8go"
)

var config = viper.New()

const javascriptSourceFileConfig = "javascript_source_file"

func init() {
	config.AutomaticEnv()
	config.SetEnvPrefix("ui")
	config.SetDefault(javascriptSourceFileConfig, "../frontend/build/bundle.js")
}

func PrintJSError(e error) {
	err, ok := e.(*v8.JSError) // JavaScript errors will be returned as the JSError struct
	if !ok {
		fmt.Println("is not JSError:", e.Error()) // the message of the exception thrown
		return
	}
	fmt.Println(err.Message)    // the message of the exception thrown
	fmt.Println(err.Location)   // the filename, line number and the column where the error occured
	fmt.Println(err.StackTrace) // the full stack trace of the error, if available

	fmt.Printf("javascript error: %v", err)        // will format the standard error message
	fmt.Printf("javascript stack trace: %+v", err) // will format the full error stack trace
}

func NewUI(frontendBuildDir string) (http.Handler, error) {
	iso := v8.NewIsolate()
	defer iso.Dispose()

	// Get value from configuration
	distFile := config.GetString(javascriptSourceFileConfig)

	bundle, err := os.ReadFile(distFile)
	if err != nil {
		return nil, fmt.Errorf("ui read bundle dist file: %w", err)
	}

	// var process = { env: ` + jsonStringify(os.Environ()) + ` };
	javascript := `
	var self = this;
	var process = { env: ` + jsonStringify(os.Environ()) + ` };
	` + string(bundle)

	script, err := iso.CompileUnboundScript(javascript, "", v8.CompileOptions{})
	if err != nil {
		PrintJSError(err)
		return nil, fmt.Errorf("ui render compile unbound script: %w", err)
	}

	cachedData := script.CreateCodeCache()

	return httpHandlerFunc(javascript, frontendBuildDir, cachedData)
}

func httpHandlerFunc(javascript, frontendBuildDir string, cachedData *v8.CompilerCachedData) (http.Handler, error) {
	// react-script build generated main.jst
	jsFiles, err := os.ReadDir(frontendBuildDir + "/static/js/")
	if err != nil {
		return nil, fmt.Errorf("read static/js directory: %w", err)
	}

	var scripts string
	for _, file := range jsFiles {
		if !(strings.Contains(file.Name(), "main")) {
			// if !(strings.Contains(file.Name(), "main") || strings.Contains(file.Name(), "chunk")) {
			continue
		}
		if !strings.Contains(file.Name(), ".js") ||
			strings.HasSuffix(file.Name(), ".map") ||
			strings.HasSuffix(file.Name(), ".LICENSE.txt") {
			continue
		}
		scripts += `<script defer="defer" src="/static/js/` + file.Name() + `"></script>`
	}
	fmt.Println("Serving JS file:", scripts)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		iso := v8.NewIsolate() // create a new JavaScript VM
		defer iso.Dispose()

		v8Ctx := v8.NewContext(iso) // new context within the VM
		defer v8Ctx.Close()         // Dispose context after using

		var wg sync.WaitGroup
		defer wg.Wait()

		done := make(chan struct{})
		defer close(done)

		wg.Add(1)
		go func() {
			defer wg.Done()

			select {
			case <-done:
				// success
			case <-ctx.Done():
				// request cancelled
				isolate := v8Ctx.Isolate()   // get the Isolate from the context
				isolate.TerminateExecution() // terminate the execution
			}
		}()

		script, err := iso.CompileUnboundScript(string(javascript), "", v8.CompileOptions{CachedData: cachedData})
		if err != nil {
			PrintJSError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		_, err = script.Run(v8Ctx)
		if err != nil {
			PrintJSError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// Get the global object from context
		global := v8Ctx.Global()

		// Then get the function you want to call on that object
		funcVal, err := global.Get("renderApp")
		if err != nil {
			PrintJSError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("cannot get global renderApp function"))
		}

		if funcVal.IsFunction() {
			// Call the function
			renderAppFunction, err := funcVal.AsFunction()
			if err != nil {
				PrintJSError(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			reqUrl, err := v8.NewValue(iso, r.URL.String())
			if err != nil {
				fmt.Println("setting request url value to pass to renderApp:", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			val, err := renderAppFunction.Call(global, reqUrl)
			if err != nil {
				PrintJSError(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			// Get HTML and styles from return object
			obj, err := val.AsObject()
			if err != nil {
				fmt.Println("renderApp return an object:", err.Error(), "obj:", val)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			htmlVal, err := obj.Get("html")
			if err != nil {
				fmt.Println("renderApp returned object contains html:", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			stylesVal, err := obj.Get("styles")
			if err != nil {
				fmt.Println("renderApp returned object contains styles:", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			html := htmlVal.String()
			styles := stylesVal.String()

			// Write the HTML and styles to our HTTP response
			responseBody := fmt.Sprintf(`
				<!doctype html>
				<html>
				  <head>
					%s
					%s
				  </head>
				  <body>
					<div id="root">%s</div>
				  </body>
				</html>`,
				styles,
				scripts,
				html,
			)
			if _, err := w.Write([]byte(responseBody)); err != nil {
				fmt.Println("writing request response:", err.Error())
				return
			}
			// Write result to HTTP response
			// fmt.Fprintf(w, "<html><body>%s</body></html>", val)
		} else {
			PrintJSError(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("renderApp is not a function"))
			return
		}
	}), nil
}

// Convert os.Environ() output to a JSON string
func jsonStringify(environ []string) string {
	envMap := make(map[string]string)
	for _, env := range environ {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envMap[pair[0]] = pair[1]
		}
	}
	jsonData, err := json.Marshal(envMap)
	if err != nil {
		return "{}"
	}
	return string(jsonData)
}
