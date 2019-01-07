package ui

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/Josempita/ipregistry/model"
	"github.com/Josempita/ipregistry/registry"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

var messageService = &model.MessagerService{nil}
var displayMessage chan string
var clients map[string]string

//messager := make(chan string)

// RandToken generates a random @l length token.
func RandToken(l int) string {
	b := make([]byte, l)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

type Config struct {
	Assets http.FileSystem
	Host   string
	Port   string
}

func Start(cfg Config, m *model.Model) {
	config := m.GetConfig()
	router := gin.Default()
	store := sessions.NewCookieStore([]byte(RandToken(64)))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})
	clients = make(map[string]string)
	messageService.Messager = make(chan model.Details)
	go registry.Start(messageService.Messager)
	displayMessage = make(chan string)
	//keep extracting the messages from the registry
	go func() {
		for message := range messageService.Messager {
			clients[message.Name] = message.Address
			table := "<table cellpadding=\"5\">"
			data := ""
			for key, value := range clients {
				data = data + "<tr><td>" + key + "</td><td>:</td><td>" + value + "</td></tr>"
			}
			table = table + data + "</table>"
			displayMessage <- table

		}
	}()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Static("/css", config.Root+"/static/css")
	router.Static("/img", config.Root+"/static/img")
	router.LoadHTMLGlob(config.Templates)
	router.Use(sessions.Sessions("front", store))
	handler := websocket.Handler(Echo)
	router.GET("/", DisplayRadioButtons(m))
	router.GET("/clients", GetClusterClients(m))
	router.GET("/clearClients", ClearClients(m))
	router.GET("/ws", EchoPageHandler)
	router.GET("/socket", EchoHandler(&handler))
	router.Run(cfg.Host + ":" + cfg.Port)

}

const (
	cdnReact           = "https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"
	cdnReactDom        = "https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"
	cdnBabelStandalone = "https://cdnjs.cloudflare.com/ajax/libs/babel-standalone/6.24.0/babel.min.js"
	cdnAxios           = "https://cdnjs.cloudflare.com/ajax/libs/axios/0.16.1/axios.min.js"
)

func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{"link": "/import", "saveLink": "/savespread"})
}

//test stuff

type RadioButton struct {
	Name       string
	Value      string
	IsDisabled bool
	IsChecked  bool
	Text       string
}

type PageVariables struct {
	PageTitle        string
	PageRadioButtons []RadioButton

	Answer string
}

func DisplayRadioButtons(m *model.Model) gin.HandlerFunc {
	// Display some radio buttons to the user
	fn := func(c *gin.Context) {
		Title := "Alchemy Cluster Monitor"
		config := m.GetConfig()

		MyRadioButtons := []RadioButton{
			RadioButton{"fileselect", "cdsCurves.csv", false, false, config.Password},
		}

		MyPageVariables := PageVariables{
			PageTitle:        Title,
			PageRadioButtons: MyRadioButtons,
		}

		c.HTML(http.StatusOK, "import.tmpl", gin.H{"PageTitle": MyPageVariables.PageTitle, "PageRadioButtons": MyPageVariables.PageRadioButtons})
	}
	return gin.HandlerFunc(fn)
}

func GetClusterClients(m *model.Model) gin.HandlerFunc {
	// Display some radio buttons to the user
	fn := func(c *gin.Context) {
		clientDetails := make([]*model.Details, 0)

		for key, value := range clients {
			client := &model.Details{Name: key, Address: value}
			clientDetails = append(clientDetails, client)
		}
		clientsJson, _ := json.Marshal(clientDetails)
		c.Writer.Write([]byte(string(clientsJson)))
	}
	return gin.HandlerFunc(fn)
}

func ClearClients(m *model.Model) gin.HandlerFunc {
	// Display some radio buttons to the user
	fn := func(c *gin.Context) {
		clients = map[string]string{}

		c.Writer.Write([]byte("Clients cleared"))
	}
	return gin.HandlerFunc(fn)
}

func EchoPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "socket.tmpl", gin.H{"link": "/import", "saveLink": "/savespread"})
}

func EchoHandler(handler *websocket.Handler) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		handler.ServeHTTP(c.Writer, c.Request)

	}
	return gin.HandlerFunc(fn)
}

func Echo(ws *websocket.Conn) {
	// var reply string

	for {

		for message := range displayMessage {

			websocket.Message.Send(ws, message)
		}

		break

	}
}
