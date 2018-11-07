package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

//Collection models a postman collection
type Collection struct {
	Info      CollectionInfo `json:"info"`
	Items     []Item         `json:"item"`
	Auth      Auth           `json:"auth"`
	Event     []Event        `json:"event"`
	Variables []Variable     `json:"variable"`
}

//Auth is the collection Auth Values
type Auth struct {
	Type   string     `json:"type"`
	Basic  []Variable `json:"basic,omitempty"`
	Bearer []Variable `json:"bearer,omitempty"`
	Digest []Variable `json:"digest,omitempty"`
	OAuth1 []Variable `json:"oauth1,omitempty"`
	OAuth2 []Variable `json:"oauth2,omitempty"`
	Hawk   []Variable `json:"hawk,omitempty"`
	AWSV4  []Variable `json:"awsv4,omitempty"`
	NTLM   []Variable `json:"ntlm,omitempty"`
}

//CollectionInfo contains the name and schema version
type CollectionInfo struct {
	Name        string `json:"name"`
	Description string `json:"description,,omitempty"`
	Schema      string `json:"schema"`
}

//Item represents a postman request item
type Item struct {
	Name    string  `json:"name"`
	Event   []Event `json:"event"`
	Request Request `json:"request"`
}

//Request is the main
type Request struct {
	Auth    Auth     `json:"auth"`
	Method  string   `json:"method"`
	Headers []KeyVal `json:"header"`
	Body    Body     `json:"body"`
	URL     URL      `json:"url"`
}

//KeyVal is a key Value struct
type KeyVal struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//Body is for HTTP request body
type Body struct {
	Mode       string     `json:"mode"`
	Raw        string     `json:"raw"`
	FormData   []Variable `json:"formdata"`
	URLEncoded []Variable `json:"urlencoded"`
}

//Variable is a collection variable
type Variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

//Event is used to model collection events like tests and pre-request scripts
type Event struct {
	Listen string `json:"listen"`
	Script `json:"script"`
}

//Script is found within events
type Script struct {
	Type string   `json:"type"`
	Exec []string `json:"exec"`
}

//URL represents the URL of the request
type URL struct {
	Raw   string   `json:"raw"`
	Host  []string `json:"host"`
	Path  []string `json:"path"`
	Query []KeyVal `json:"query"`
}

func splitCmd(cmd *cobra.Command, args []string) {
	status := func() int {
		//read the file and unmarshal to struct
		fmt.Printf("Reading Collection File : %v \n", args[0])
		collection, err := readCollectionFile(args[0])
		//create directory
		fmt.Printf("Creating Collection Directory : %v \n", collection.Info.Name)
		err = createCollectionDirectory(collection.Info.Name)
		if err != nil {
			return 1
		}
		//put down the collection variables
		fmt.Println("Saving Collection Details")
		err = saveCollectionDetails(collection, collection.Info.Name)
		if err != nil {
			return 1
		}
		//put the requests down
		fmt.Println("Saving Requests")
		err = saveRequestFiles(collection.Items, collection.Info.Name)
		if err != nil {
			return 1
		}
		return 0
	}()
	os.Exit(status)
}

//readCollectionFile reads the collection file and returns a collection or error
func readCollectionFile(filename string) (*Collection, error) {
	plan, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Unable to read File at : %v \n", filename)
		return nil, err
	}
	var in Collection
	err = json.Unmarshal(plan, &in)
	return &in, err
}

func createCollectionDirectory(directoryName string) error {
	fmt.Println("Changing Directory Permissins to 0777")
	return os.Mkdir(directoryName, 0777)
}

func saveCollectionDetails(coll *Collection, directoryName string) error {
	//naming convention is type:collection.json
	//info
	info, err := json.Marshal(coll.Info)
	if err != nil {
		fmt.Printf("Error generating Info JSON : %v \n", err)
		return err
	}
	err = writeFile(info, coll.Info.Name, "info:collection")
	if err != nil {
		fmt.Printf("Error Writing Info File : %v \n", err)
		return err
	}
	//event
	event, err := json.Marshal(coll.Event)
	if err != nil {
		fmt.Printf("Error generating Event JSON : %v \n", err)
		return err
	}
	err = writeFile(event, directoryName, "event:collection")
	if err != nil {
		fmt.Printf("Error Writing Event File : %v \n", err)
		return err
	}
	//auth
	auth, err := json.Marshal(coll.Auth)
	if err != nil {
		fmt.Printf("Error generating Auth JSON : %v \n", err)
		return err
	}
	err = writeFile(auth, coll.Info.Name, "auth:collection")
	if err != nil {
		fmt.Printf("Error Writing Auth File : %v \n", err)
		return err
	}

	//variable
	variables, err := json.Marshal(coll.Variables)
	if err != nil {
		fmt.Printf("Error generating Variables JSON : %v \n", err)
		return err
	}
	err = writeFile(variables, coll.Info.Name, "variables:collection")
	if err != nil {
		fmt.Printf("Error Writing Varaibles File : %v \n", err)
		return err
	}
	return nil
}

func saveRequestFiles(items []Item, directoryName string) error {
	for _, item := range items {
		itemJSON, err := json.Marshal(item)
		if err != nil {
			fmt.Printf("Error generating item JSON : %v \n", err)
			return err
		}
		err = writeFile(itemJSON, directoryName, fmt.Sprintf("Request:%v:%v", item.Name, item.Request.Method))
		if err != nil {
			fmt.Printf("Error Writing item File : %v \n", err)
			return err
		}
	}
	return nil
}

func writeFile(contents []byte, directoryName, fileName string) error {
	if !strings.HasSuffix(directoryName, "/") {
		directoryName = fmt.Sprintf("%v/", directoryName)
	}
	path := fmt.Sprintf("%v%v", directoryName, fileName)
	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("Unable to create file : %v \n", path)
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = w.Write(contents)
	if err != nil {
		fmt.Printf("Unable to create file : %v \n", path)
		return err
	}
	w.Flush()
	return err
}
