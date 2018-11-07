package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/patrickwalker/PostmanPat/tmpl"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

func joinCmd(cmd *cobra.Command, args []string) {
	status := func() int {
		//directory is the Name
		directoryName := args[0]
		fmt.Printf("Sourcing Collection from : %v \n", directoryName)
		info, err := loadCollectionInfo(directoryName)
		if err != nil {
			return 1
		}
		auth, err := loadCollectionAuth(directoryName)
		if err != nil {
			return 1
		}
		event, err := loadCollectionEvents(directoryName)
		if err != nil {
			return 1
		}
		variables, err := loadCollectionVariables(directoryName)
		if err != nil {
			return 1
		}
		items, err := loadCollectionItems(directoryName)
		if err != nil {
			return 1
		}
		collection := Collection{
			Auth:      auth,
			Event:     event,
			Variables: variables,
			Items:     items,
			Info:      info,
		}

		fmt.Printf("Writing collection file to %v \n", info.Name)
		err = writeCollection(collection)
		//write the collection file out using the template
		if err != nil {
			return 1
		}
		return 0
	}()
	os.Exit(status)
}

func loadCollectionInfo(directoryName string) (CollectionInfo, error) {
	val, err := readJSONFile(fmt.Sprintf("%v/info:collection", directoryName))
	if err != nil {
		fmt.Printf("Unable to read JSON File : %v \n", err)
		return CollectionInfo{}, err
	}
	var coll CollectionInfo
	err = json.Unmarshal(val, &coll)
	return coll, err
}

func loadCollectionAuth(directoryName string) (Auth, error) {
	//auth:collection
	val, err := readJSONFile(fmt.Sprintf("%v/auth:collection", directoryName))
	if err != nil {
		return Auth{}, err
	}
	var auth Auth
	err = json.Unmarshal(val, &auth)
	return auth, err
}

func loadCollectionEvents(directoryName string) ([]Event, error) {
	//event:collection
	val, err := readJSONFile(fmt.Sprintf("%v/event:collection", directoryName))
	if err != nil {
		return nil, err
	}
	var events []Event
	err = json.Unmarshal(val, &events)
	return events, err
}

func loadCollectionVariables(directoryName string) ([]Variable, error) {
	//variables:collection
	val, err := readJSONFile(fmt.Sprintf("%v/variables:collection", directoryName))
	if err != nil {
		return nil, err
	}
	var variables []Variable
	err = json.Unmarshal(val, &variables)
	return variables, err
}

func loadCollectionItems(directoryName string) ([]Item, error) {
	//regex of files
	//Request:*
	itemFiles, err := filepath.Glob(fmt.Sprintf("%v/Request:*", directoryName))
	if err != nil {
		fmt.Println("Unable to search for request files")
		return nil, err
	}
	items := make([]Item, len(itemFiles))
	for index, itemFile := range itemFiles {
		val, err := readJSONFile(itemFile)
		if err != nil {
			fmt.Printf("Unable to read : %v \n", itemFile)
		}
		var it Item
		err = json.Unmarshal(val, &it)
		if err != nil {
			fmt.Printf("Unable to parse : %v \n", val)
		}
		items[index] = it
	}

	return items, nil
}

func writeCollection(coll Collection) error {
	//generate a writer
	f, err := os.Create(fmt.Sprintf("Test-%v.postman_collection.json", coll.Info.Name))
	if err != nil {
		fmt.Printf("Unable to create destination file : %v \n", coll.Info.Name)
		return err
	}
	//get template name eventually configurable based on schema version
	return generateCollectionFile(coll, tmpl.CollectionTemplate, f)
}

func readJSONFile(fileName string) ([]byte, error) {
	jsonFile, err := os.Open(fileName)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Printf("Unable to open the file : %v \n", fileName)
		return nil, err
	}
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf("Unable to read file : %v \n", fileName)
		return nil, err
	}
	return byteValue, err
}

func generateCollectionFile(coll Collection, tmpl string, writer io.Writer) error {
	tem, err := template.New("Collection").Funcs(getTemplFunctions()).Parse(tmpl)
	if err != nil {
		fmt.Printf("Unable to parse template: %v \n", err)
		return err
	}
	err = tem.Execute(writer, coll)
	if err != nil {
		fmt.Printf("Unable to execute template: %v \n", err)
		return err
	}
	return nil
}

func getTemplFunctions() template.FuncMap {
	return template.FuncMap{
		"uuid":    uuid.NewV4,
		"endItem": endItem,
		"raw":     strconv.Quote,
	}
}

//isLastItem is used to work out if current index is last item in map
func isLastItem(index, length int) bool {
	if index == (length - 1) {
		return true
	}
	return false
}

//endItem returns a , if not the last item
func endItem(index, length int) string {
	if !isLastItem(index, length) {
		return ","
	}
	return ""
}
