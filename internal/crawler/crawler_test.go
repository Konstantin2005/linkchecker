package crawler

import (
	"errors"
	"linkchecker/config"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"
	"time"
)

func newSummary() *config.Summary {
	return &config.Summary{
		ErrorByType:  make(map[int]int),
		ProblemLinks: make(map[string]config.CheckResult),
	}
}

// после каждого теста надо чистить глобальные переменные, иначе случаи «уже посещено» будут
// влиять на другие тесты.
func resetGlobals() {
	Visited = make(map[string]bool)
	HttpClient = http.DefaultClient
}

/*** ТЕСТ AddProblem ***/
func TestAddProblem(t *testing.T) {
	defer resetGlobals()
	sum := newSummary()

	root, _ := url.Parse("https://ex.org/abc/def")
	AddProblem(sum, "https://ex.org/abc/def/page.html", 2, root, 404,
		errors.New("not found"), 123*time.Millisecond)

	if got := sum.ErrorByType[404]; got != 1 {
		t.Fatalf("ErrorByType[404]=%d, want 1", got)
	}
	cr, ok := sum.ProblemLinks["https://ex.org/abc/def/page.html"]
	if !ok {
		t.Fatalf("no record in ProblemLinks")
	}
	if cr.StatusCode != 404 || cr.Depth != 2 {
		t.Fatalf("bad CheckResult: %+v", cr)
	}
	wantPath := path.Dir("/abc/def/page.html")
	if cr.Referrer.Path != wantPath {
		t.Fatalf("Referrer.Path=%q, want %q", cr.Referrer.Path, wantPath)
	}
}

/*** ТЕСТЫ Crawl ***/
func TestCrawl(t *testing.T) {
	defer resetGlobals()

	// поднимем HTTP-сервер с тремя эндпойнтами
	mux := http.NewServeMux()

	// /
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		io := `<html><body>
				<a href="/inner">inner</a>
				<a href="/bad">bad</a>
				<a href="/file.pdf">pdf</a>
			   </body></html>`
		w.Write([]byte(io))
	})

	// /inner – успешный, html без ссылок
	mux.HandleFunc("/inner", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<p>ok</p>"))
	})

	// /bad – вернёт 404
	mux.HandleFunc("/bad", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// /file.pdf – не-html
	mux.HandleFunc("/file.pdf", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	// инициализируем необходимые глобальные переменные
	HttpClient = ts.Client()
	resetGlobals() // сброс Visited

	rootURL, _ := url.Parse(ts.URL + "/")
	sum := newSummary()

	Crawl(rootURL, rootURL, 0, 3, sum)

	/* ----------- проверки ---------- */
	t.Run("totals", func(t *testing.T) {
		if sum.TotalLinks != 4 { // /, /inner, /bad, /file.pdf
			t.Fatalf("TotalLinks=%d, want 4", sum.TotalLinks)
		}
		if sum.Successful != 3 { // 200, 200, 200 (pdf)
			t.Fatalf("Successful=%d, want 3", sum.Successful)
		}
		if len(sum.ProblemLinks) != 1 {
			t.Fatalf("Problems=%d, want 1", len(sum.ProblemLinks))
		}
		if sum.ErrorByType[404] != 1 {
			t.Fatalf("ErrorByType[404]=%d, want 1", sum.ErrorByType[404])
		}
	})

	t.Run("depthLimit", func(t *testing.T) {
		resetGlobals()
		sum2 := newSummary()
		Crawl(rootURL, rootURL, 0, 0, sum2) // глубина = 0, т.е. только корень
		if sum2.TotalLinks != 1 {
			t.Fatalf("TotalLinks=%d, want 1", sum2.TotalLinks)
		}
	})
}

/*** ТЕСТ: дубликаты (Visited) ***/
func TestCrawl_Dedup(t *testing.T) {
	defer resetGlobals()

	// страница /dup с двумя одинаковыми ссылками на one.html
	htmlWithDup := `<html><body>
		<a href="/one.html">x</a>
		<a href="/one.html">x</a>
	</body></html>`

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlWithDup))
	})
	mux.HandleFunc("/one.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<p>one</p>"))
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	HttpClient = ts.Client()
	sum := newSummary()
	root, _ := url.Parse(ts.URL + "/")

	Crawl(root, root, 0, 2, sum)

	if sum.TotalLinks != 2 { // / и /one.html, а не 3
		t.Fatalf("dedup failed: TotalLinks=%d, want 2", sum.TotalLinks)
	}
	if len(Visited) != 2 {
		t.Fatalf("Visited=%d, want 2", len(Visited))
	}
}
