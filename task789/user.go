package task789

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"os"
	"strings"
	"sync"
)

type user struct {
	login  string
	pass   string
	role   string
	access []string
}

type logins struct {
	list map[string]user
	mu   *sync.Mutex
}

var users = logins{list: make(map[string]user), mu: &sync.Mutex{}}

func loadLogins() string {
	users.mu.Lock()
	defer users.mu.Unlock()

	file, err := os.Open("task789/logins.txt")
	if err != nil {
		file, err := os.Create("task789/logins.txt")
		if err != nil {
			return "error"
		}
		defer file.Close()
		file.WriteString("su 21232f297a57a5a743894a0e4a801fc3 superuser all\n")

		return "ok"
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		splitted := strings.Split(scanner.Text(), " ")
		if len(splitted) < 4 {
			return "error"
		}
		if splitted[2] != "superuser" && splitted[2] != "admin" && splitted[2] != "editor" && splitted[2] != "user" {
			return "error"
		}
		newUser := user{login: splitted[0], pass: splitted[1], role: splitted[2], access: strings.Split(splitted[3], ";")}
		users.list[newUser.login] = newUser
	}

	return "ok"
}

func (l *logins) checkLogin(login, password string) (bool, string) {
	users.mu.Lock()
	defer users.mu.Unlock()

	user, ok := l.list[login]
	if !ok {
		return false, "not a user"
	}
	if !user.checkPass(password) {
		return false, "wrong password"
	}
	return true, "ok"
}

func (u *user) checkPass(password string) bool {
	return u.pass == password
}

func (u *user) checkRole(role string) bool {
	return u.role == role
}

func (u *user) checkPoolAccess(p string) (bool, string) {
	words := strings.Split(p, ".")
	if len(words) != 1 {
		return false, "invalid pool name"
	}
	for _, pool := range u.access {
		if pool == words[0] {
			return true, "ok"
		}
	}
	return false, "ok"
}

func (u *user) checkSchemaAccess(ps string) (bool, string) {
	words := strings.Split(ps, ".")
	if len(words) != 2 {
		return false, "invalid schema name"
	}
	_, res := u.checkPoolAccess(words[0])
	if res != "ok" {
		return false, res
	}

	for _, schema := range u.access {
		if schema == words[1] {
			return true, "ok"
		}
	}
	return false, "ok"
}

func (u *user) checkCollectionAccess(psc string) (bool, string) {
	words := strings.Split(psc, ".")
	if len(words) != 3 {
		return false, "invalid collection name"
	}
	_, res := u.checkSchemaAccess(words[0] + "." + words[1])
	if res != "ok" {
		return false, res
	}

	for _, collection := range u.access {
		if collection == words[2] {
			return true, "ok"
		}
	}
	return false, "ok"
}

func blocks(username string) string {
	u, ok := users.list[username]
	if !ok {
		return "blocks_notuser.html"
	}
	ok, _ = users.checkLogin(u.login, u.pass)
	if !ok {
		return "blocks_notuser.html"
	}
	return "blocks_notuser.html"
}

func isLoggedIn(w http.ResponseWriter, r *http.Request) (bool, string, string) {
	_, erruser := r.Cookie("currentUser")
	_, errpassword := r.Cookie("currentPassword")

	if erruser != nil {
		cookie := &http.Cookie{
			Name:   "currentUser",
			Value:  "",
			MaxAge: 604800,
		}
		http.SetCookie(w, cookie)

		if errpassword != nil {
			cookie = &http.Cookie{
				Name:   "currentPassword",
				Value:  "",
				MaxAge: 604800,
			}
			http.SetCookie(w, cookie)
		}
		return false, "", "ok"
	} else {
		if errpassword != nil {
			cookie := &http.Cookie{
				Name:   "currentPassword",
				Value:  "",
				MaxAge: 604800,
			}
			http.SetCookie(w, cookie)
			return false, "", "ok"
		}
	}

	user, _ := r.Cookie("currentUser")
	password, _ := r.Cookie("currentPassword")

	ok, res := users.checkLogin(user.Value, password.Value)

	return ok, user.Value, res
}

func login(w http.ResponseWriter, username, pass string) bool {
	users.mu.Lock()
	defer users.mu.Unlock()

	user, ok := users.list[username]
	if !ok {
		return false
	}
	password := hashPassword(pass)
	if !user.checkPass(password) {
		return false
	}

	cookie := &http.Cookie{
		Name:   "currentUser",
		Value:  username,
		MaxAge: 604800,
	}
	http.SetCookie(w, cookie)
	cookie = &http.Cookie{
		Name:   "currentPassword",
		Value:  password,
		MaxAge: 604800,
	}
	http.SetCookie(w, cookie)
	return true
}

func register(w http.ResponseWriter, username, pass string) string {
	users.mu.Lock()
	defer users.mu.Unlock()

	password := hashPassword(pass)

	if _, ok := users.list[username]; ok {
		return "username is already taken"
	}
	newUser := user{login: username, pass: password, role: "user", access: []string{"none"}}
	users.list[username] = newUser

	file, err := os.OpenFile("task789/logins.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		file, err := os.Create("task789/logins.txt")
		if err != nil {
			return "server error"
		}
		defer file.Close()

		for _, user := range users.list {
			file.WriteString(user.login + " " + user.pass + " " + user.role + " " + strings.Join(user.access, ";") + "\n")
		}
		return "ok"
	}
	defer file.Close()
	file.WriteString(username + " " + password + " " + "user" + " " + strings.Join(newUser.access, ";") + "\n")
	return "ok"
}

func logout(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "currentUser",
		Value:  "",
		MaxAge: 1,
	}
	http.SetCookie(w, cookie)
	cookie = &http.Cookie{
		Name:   "currentPassword",
		Value:  "",
		MaxAge: 1,
	}
	http.SetCookie(w, cookie)
}

func hashPassword(original string) string {
	hasher := md5.New()
	hasher.Write([]byte(original))
	hashed := hasher.Sum(nil)
	return hex.EncodeToString(hashed)
}
