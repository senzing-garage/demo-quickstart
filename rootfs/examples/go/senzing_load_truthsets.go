package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"

	"github.com/senzing-garage/sz-sdk-go-grpc/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DataSourceKey struct {
	Data_Source string
}

type Record struct {
	Data_Source string
	Record_ID   string
}

var (
	ctx               = context.TODO()
	grpcAddress       = "localhost:8261"
	jsonDataSource    DataSourceKey
	homePath          = "./"
	jsonRecord        Record
	truthSetURLPrefix = "https://raw.githubusercontent.com/Senzing/truth-sets/refs/heads/main/truthsets/demo/"
	truthSetFileNames = []string{"customers.json", "reference.json", "watchlist.json"}
)

func testErr(err error) {
	if err != nil {
		panic(err)
	}
}

func downloadFile(url string, filepath string) error {
	outputFile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(outputFile, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func getDataSources() []string {
	result := []string{}
	for i := 0; i < len(truthSetFileNames); i++ {
		filepath := fmt.Sprintf("%s%s", homePath, truthSetFileNames[i])
		file, err := os.Open(filepath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Bytes()
			err := json.Unmarshal(line, &jsonDataSource)
			testErr(err)
			if !slices.Contains(result, jsonDataSource.Data_Source) {
				result = append(result, jsonDataSource.Data_Source)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	return result
}

func asPrettyJSON(str string) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return str
	}
	return prettyJSON.String()
}

func main() {

	// Download truth-sets.

	for i := 0; i < len(truthSetFileNames); i++ {
		url := fmt.Sprintf("%s/%s", truthSetURLPrefix, truthSetFileNames[i])
		filepath := fmt.Sprintf("%s%s", homePath, truthSetFileNames[i])
		err := downloadFile(url, filepath)
		testErr(err)
	}

	// Identify datasources.

	var dataSources = getDataSources()
	fmt.Printf("Found the following DATA_SOURCE values in the data: %v\n", dataSources)

	grpcConnection, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	testErr(err)
	szAbstractFactory := &szabstractfactory.Szabstractfactory{
		GrpcConnection: grpcConnection,
	}

	szConfig, err := szAbstractFactory.CreateConfig(ctx)
	testErr(err)

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	testErr(err)

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	testErr(err)

	oldConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	testErr(err)

	oldJsonConfig, err := szConfigManager.GetConfig(ctx, oldConfigID)
	testErr(err)

	configHandle, err := szConfig.ImportConfig(ctx, oldJsonConfig)
	testErr(err)

	for _, value := range dataSources {
		_, err := szConfig.AddDataSource(ctx, configHandle, value)
		if err != nil {
			fmt.Println(err)
		}
	}

	newJsonConfig, err := szConfig.ExportConfig(ctx, configHandle)
	testErr(err)

	newConfigID, err := szConfigManager.AddConfig(ctx, newJsonConfig, "Add TruthSet datasources")
	testErr(err)

	err = szConfigManager.ReplaceDefaultConfigID(ctx, oldConfigID, newConfigID)
	testErr(err)

	err = szAbstractFactory.Reinitialize(ctx, newConfigID)
	testErr(err)

	// Add records.

	for _, value := range truthSetFileNames {
		filepath := fmt.Sprintf("%s%s", homePath, value)
		file, err := os.Open(filepath)
		testErr(err)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Bytes()
			err := json.Unmarshal(line, &jsonRecord)
			testErr(err)
			result, err := szEngine.AddRecord(ctx, jsonRecord.Data_Source, jsonRecord.Record_ID, string(line), senzing.SzWithInfo)
			testErr(err)
			fmt.Println(result)
		}
	}

	// View results.

	customer1070Entity, err := szEngine.GetEntityByRecordID(ctx, "CUSTOMERS", "1070", senzing.SzEntityIncludeRecordSummary)
	testErr(err)
	fmt.Println(asPrettyJSON(customer1070Entity))

	searchProfile := ""
	searchQuery := `{
        "name_full": "robert smith",
        "date_of_birth": "11/12/1978"
    }`
	searchResult, err := szEngine.SearchByAttributes(ctx, searchQuery, searchProfile, senzing.SzSearchByAttributesDefaultFlags)
	testErr(err)
	fmt.Println(asPrettyJSON(searchResult))
}
