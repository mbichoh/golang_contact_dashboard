package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mbichoh/contactDash/pkg/forms"
	"github.com/mbichoh/contactDash/pkg/models"
)

func (app *application) signupForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "mobile", "password")
	form.MobileNumberCheck("mobile", forms.NumberCheck)
	form.MobileCountryCheckCode("mobile", forms.NumberValid)
	form.MobileCheckPref("mobile")

	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	// CHECK : what if i give my phone number as "1234567891"?
	// CHECK: what if i give my phone number as "abcdefghij"?
	// CHECK: what if my phone number is in Congo?

	rand.Seed(time.Now().UnixNano())
	min := 10000
	max := 99999
	tokenCode := rand.Intn(max-min+1) + min

	err = app.users.Insert(form.Get("name"), form.Get("mobile"), form.Get("password"), tokenCode, false)

	if err == models.ErrDuplicateNumber {

		// CHECK: what if mysql change their error number and error messages? Your application will break

		form.Errors.Add("mobile", "Phone number already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// endpoint
	var sendMessageURL string = "https://api.amisend.com/v1/sms/send"

	// authentication

	var username string = ""
	var apikey string = ""

	// data

	messageData := map[string]string{
		"phoneNumbers": form.Get("mobile"),
		"message":      "Welcome " + form.Get("name") + " to Nathan's Amisend integretion. Your activation code is " + strconv.Itoa(tokenCode),
		"senderId":     "", // leave blank if you do not have a custom sender Id
	}

	params, _ := json.Marshal(messageData)

	request, err := http.NewRequest("POST", sendMessageURL, bytes.NewBuffer(params))

	if err != nil {
		panic(err.Error())
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Set("x-api-user", username)
	request.Header.Set("x-api-key", apikey)
	request.Header.Set("Content-Length", strconv.Itoa(len(params)))

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err.Error())
	}

	defer response.Body.Close()

	fmt.Println(string(body))
	app.session.Put(r, "flash", "Sign up successful. Please check message and verify...")

	http.Redirect(w, r, "/user/verification", http.StatusSeeOther)
}

func (app *application) loginForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)

	// CHECK: what if i give my phone number as "1234567891"?
	// CHECK: what if i give my phone number as "abcdefghij"? done
	// CHECK: what if my phone number is in Congo?

	id, err := app.users.Authenticate(form.Get("mobile"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Phone number or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err == models.ErrNotVerified {
		form.Errors.Add("activation", "Please check your message to verify...")
		app.render(w, r, "verification.page.tmpl", &templateData{
			Form: form,
		})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// CHECK : How do you confirm the phone number I gave is actually mine?
	// CHECK : You should send me a unique code and ask me to give it to you, if it matches let me login else deny me a chance

	app.session.Put(r, "userID", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) verification(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "verification.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) verified(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("activation")

	actCode, _ := strconv.Atoi(form.Get("activation"))
	id, err := app.users.Verify(actCode)
	if err == models.ErrNoRecord {
		form.Errors.Add("activation", "Invalid verification code, please check again.")
		app.render(w, r, "verification.page.tmpl", &templateData{Form: form})
		return
	} else if !form.Valid() {
		app.render(w, r, "verification.page.tmpl", &templateData{Form: form})
		return
	}
	fmt.Printf("%d", id.Token)
	idNo, err := app.users.IsVerified(actCode)
	if err != nil {
		app.serverError(w, err)
		return
	}
	fmt.Printf("Rows affected := %d\n", idNo)
	app.session.Put(r, "flash", "Account verified... Please login.")
	http.Redirect(w, r, "/user/login", 303)

}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {

	// CHECK: Its always important to check if a session actually exists before attempting to remove it -done

	if app.session.Exists == nil {
		app.session.Put(r, "flash", "No session exists.")
	} else {
		app.session.Remove(r, "userID")
		app.session.Put(r, "flash", "Logged out successfully")
		http.Redirect(w, r, "/user/login", 303)
	}

}

func (app *application) ContactHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	uid := app.session.GetInt(r, "userID")

	s, err := app.contacts.Latest(uid)
	if err != nil {
		app.serverError(w, err)
		return
	}

	g, err := app.groups.GroupFetchNames(uid)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Contacts: s,
		Groups:   g,
	})

}

func (app *application) ShowContact(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.contacts.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Contact: s,
	})

}

func (app *application) SendMessageToContact(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusSeeOther)
		return
	}
	// form := forms.New(r.PostForm)
	// form.Required("msgBody")
	// phoneNo := form.Get("phoneNo")
	// msg := form.Get("msgBody")

	//api sms
	app.session.Put(r, "flash", "Message sent successful")
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.contacts.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Contact: s,
	})
}

func (app *application) CreateContactForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) CreateContact(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusSeeOther)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("name", "mobile")
	form.MobileNumberCheck("mobile", forms.NumberCheck)
	form.MobileCountryCheckCode("mobile", forms.NumberValid)
	form.MobileCheckPref("mobile")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	uid := app.session.GetInt(r, "userID")

	id, err := app.contacts.Insert(form.Get("name"), form.Get("mobile"), uid)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Contact created successful")
	http.Redirect(w, r, fmt.Sprintf("/contact/%d", id), http.StatusSeeOther)

}

func (app *application) FetchUpdateContact(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.contacts.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "update.page.tmpl", &templateData{
		Contact: s,
	})
}

func (app *application) UpdateContact(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("up_name", "up_contact", "up_id")
	form.MobileNumberCheck("up_contact", forms.NumberCheck)
	form.MobileCountryCheckCode("up_contact", forms.NumberValid)
	form.MobileCheckPref("up_contact")

	if !form.Valid() {
		app.render(w, r, "update.page.tmpl", &templateData{Form: form})
		return
	}
	contId, _ := strconv.Atoi(form.Get("up_id"))

	idNo, err := app.contacts.Update(form.Get("up_name"), form.Get("up_contact"), contId)
	if err != nil {
		app.serverError(w, err)
		return
	}
	fmt.Printf("Rows affected := %d\n", idNo)
	app.session.Put(r, "flash", "Contact updated successful")
	http.Redirect(w, r, fmt.Sprintf("/contact/%d", id), http.StatusSeeOther)
}

func (app *application) DelContact(w http.ResponseWriter, r *http.Request) {
	idNo, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || idNo < 1 {
		app.notFound(w)
		return
	}

	id, err := app.contacts.Delete(idNo)
	if err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Printf("%d", id)
	app.session.Put(r, "flash", "Contact deleted successful")
	http.Redirect(w, r, fmt.Sprint("/"), http.StatusSeeOther)
}

func (app *application) GroupContacts(w http.ResponseWriter, r *http.Request) {

	uid := app.session.GetInt(r, "userID")

	s, err := app.contacts.Latest(uid)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "group.page.tmpl", &templateData{
		Contacts: s,
	})
}

func (app *application) GroupedContacts(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.contacts.GetGroupedContacts(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	g, err := app.groups.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "grouped.page.tmpl", &templateData{
		Contacts: s,
		Group:    g,
	})
}

func (app *application) SendMessageToGroup(w http.ResponseWriter, r *http.Request) {

}

func (app *application) DelGroupContact(w http.ResponseWriter, r *http.Request) {
	idNo, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || idNo < 1 {
		app.notFound(w)
		return
	}

	id, err := app.groupedcontacts.DeleteContact(idNo)
	if err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Printf("%d", id)
	app.session.Put(r, "flash", "Contact from group deleted successful")
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func (app *application) DelGroup(w http.ResponseWriter, r *http.Request) {
	idNo, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || idNo < 1 {
		app.notFound(w)
		return
	}

	id, err := app.groupedcontacts.DeleteGroup(idNo)
	if err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Printf("%d", id)
	app.session.Put(r, "flash", "Group deleted successful")
	http.Redirect(w, r, fmt.Sprint("/"), http.StatusSeeOther)
}

func (app *application) DispGroupedContacts(w http.ResponseWriter, r *http.Request) {

	uid := app.session.GetInt(r, "userID")

	s, err := app.contacts.Latest(uid)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "group.page.tmpl", &templateData{
		Contacts: s,
	})
}

func (app *application) CreateGroupedContacts(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusSeeOther)
		return
	}

	contactid := r.FormValue("format")
	groupname := r.FormValue("gname")

	fmt.Printf("%s", groupname)
	fmt.Printf("%s", contactid)

	idcn, err := app.groups.GroupInsertName(groupname)
	if err != nil {
		app.serverError(w, err)
		return
	}

	abc := strings.Split(contactid, ",")
	for _, b := range abc {
		d, _ := strconv.Atoi(b)
		id, err := app.groupedcontacts.Insert(d, idcn)
		if err != nil {
			app.serverError(w, err)
			return
		}
		fmt.Printf("%d", id)
	}

	fmt.Printf("%d", idcn)
	app.session.Put(r, "flash", "Group Created successful")
	fmt.Fprintln(w, strconv.Itoa(idcn))
	// http.Redirect(w, r, "/contact/group/"+strconv.Itoa(idcn), http.StatusSeeOther)
}
