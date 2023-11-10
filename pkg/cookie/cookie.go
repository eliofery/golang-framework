package cookie

import "net/http"

func New(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
}

func Set(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, New(name, value))
}

func Get(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func Delete(w http.ResponseWriter, name string) {
	cookie := New(name, "")
	cookie.MaxAge = -1

	http.SetCookie(w, cookie)
}
