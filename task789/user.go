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
	Login  string
	Role   string
	pass   string
	access []string
}

type logins struct {
	List map[string]user
	mu   *sync.Mutex
}

var users = logins{List: make(map[string]user), mu: &sync.Mutex{}}

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
		newUser := user{Login: splitted[0], pass: splitted[1], Role: splitted[2], access: strings.Split(splitted[3], ";")}
		users.List[newUser.Login] = newUser
	}

	return "ok"
}

func (l *logins) checkLogin(login, password string) (bool, string) {
	users.mu.Lock()
	defer users.mu.Unlock()

	user, ok := l.List[login]
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
	return u.Role == role
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

func isLoggedIn(w http.ResponseWriter, r *http.Request) (bool, string, string) {
	_, erruser := r.Cookie("currentUser")
	_, errpassword := r.Cookie("currentPassword")

	if erruser != nil {
		updateCookie(w, "currentUser", "", 604800)

		if errpassword != nil {
			updateCookie(w, "currentPassword", "", 604800)
		}
		return false, "", "ok"
	} else {
		if errpassword != nil {
			updateCookie(w, "currentPassword", "", 604800)
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

	user, ok := users.List[username]
	if !ok {
		return false
	}
	password := hashPassword(pass)
	if !user.checkPass(password) {
		return false
	}

	updateCookie(w, "currentUser", username, 604800)
	updateCookie(w, "currentPassword", password, 604800)
	return true
}

func register(w http.ResponseWriter, username, pass string) string {
	users.mu.Lock()
	defer users.mu.Unlock()

	password := hashPassword(pass)

	if _, ok := users.List[username]; ok {
		return "username is already taken"
	}
	newUser := user{Login: username, pass: password, Role: "user", access: []string{"none"}}
	users.List[username] = newUser

	updateCookie(w, "currentUser", username, 604800)
	updateCookie(w, "currentPassword", password, 604800)

	file, err := os.OpenFile("task789/logins.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return overrideFile()
	}
	defer file.Close()
	file.WriteString(username + " " + password + " " + "user" + " " + strings.Join(newUser.access, ";") + "\n")
	return "ok"
}

func updatePassword(w http.ResponseWriter, username, pass string) string {
	users.mu.Lock()
	defer users.mu.Unlock()

	u, ok := users.List[username]
	if !ok {
		return "not a user"
	}

	password := hashPassword(pass)

	newUser := user{Login: u.Login, pass: password, Role: u.Role, access: u.access}
	users.List[u.Login] = newUser
	updateCookie(w, "currentPassword", password, 604800)

	overrideFile()

	return "ok"
}

func updateRole(w http.ResponseWriter, username, role string) string {
	users.mu.Lock()
	defer users.mu.Unlock()

	u, ok := users.List[username]
	if !ok {
		return "not a user"
	}

	newUser := user{Login: u.Login, pass: u.pass, Role: role, access: u.access}
	users.List[u.Login] = newUser

	overrideFile()

	return "ok"
}

func updateAccess(w http.ResponseWriter, username, access, read, write string) string {
	users.mu.Lock()
	defer users.mu.Unlock()

	u, ok := users.List[username]
	if !ok {
		return "not a user"
	}

	accessLine := access

	switch read + write {
	case "truetrue", "falsetrue":
		accessLine += "/w"
	case "truefalse":
		accessLine += "/r"
	case "falsefalse":
		accessLine = ""
	}

	acc := &u.access

	for i, a := range *acc {
		if len(a) >= len(access) && a[:len(access)] == access {
			if accessLine == "" {
				*acc = append((*acc)[:i], (*acc)[i+1:]...)
			} else {
				(*acc)[i] = accessLine
			}
			overrideFile()
			return "ok"
		}
	}

	if accessLine != "" {
		*acc = append(*acc, accessLine)
		overrideFile()
		return "ok"
	}

	return "ok"
}

func updateCookie(w http.ResponseWriter, cookieName, newVal string, expTime int) {
	cookie := &http.Cookie{
		Name:   cookieName,
		Value:  newVal,
		MaxAge: expTime,
	}
	http.SetCookie(w, cookie)
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

func overrideFile() string {
	file, err := os.Create("task789/logins.txt")
	if err != nil {
		return "server error"
	}
	defer file.Close()

	for _, user := range users.List {
		file.WriteString(user.Login + " " + user.pass + " " + user.Role + " " + strings.Join(user.access, ";") + "\n")
	}
	return "ok"
}

func deleteUser(username string) {
	delete(users.List, username)

	overrideFile()
}
