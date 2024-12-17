package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	ctx            = context.TODO()
	grpcAddress    = "localhost:8261"
	jsonDataSource DataSourceKey
	homePath       = "./"
	jsonRecord     Record
	fileName       = "senzing-example-data.json"
)

func testErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getDataSources(filePath string) []string {
	result := []string{}
	file, err := os.Open(filePath)
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
	return result
}

func updateSenzingConfig(szAbstractFactory senzing.SzAbstractFactory, dataSourceNames []string) error {

	szConfig, err := szAbstractFactory.CreateConfig(ctx)
	if err != nil {
		return err
	}

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		return err
	}

	oldConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		return err
	}

	oldJsonConfig, err := szConfigManager.GetConfig(ctx, oldConfigID)
	if err != nil {
		return err
	}

	configHandle, err := szConfig.ImportConfig(ctx, oldJsonConfig)
	if err != nil {
		return err
	}

	for _, value := range dataSourceNames {
		_, err := szConfig.AddDataSource(ctx, configHandle, value)
		if err != nil {
			fmt.Println(err)
		}
	}

	newJsonConfig, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		return err
	}

	newConfigID, err := szConfigManager.AddConfig(ctx, newJsonConfig, "Add TruthSet datasources")
	if err != nil {
		return err
	}

	err = szConfigManager.ReplaceDefaultConfigID(ctx, oldConfigID, newConfigID)
	if err != nil {
		return err
	}

	err = szAbstractFactory.Reinitialize(ctx, newConfigID)
	if err != nil {
		return err
	}

	return nil
}

func addRecords(szAbstractFactory senzing.SzAbstractFactory, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		err := json.Unmarshal(line, &jsonRecord)
		testErr(err)
		result, err := szEngine.AddRecord(ctx, jsonRecord.Data_Source, jsonRecord.Record_ID, string(line), senzing.SzWithInfo)
		testErr(err)
		fmt.Println(result)
	}
	return nil
}

func main() {

	// User input.

	inputFile := fmt.Sprintf("%s%s", homePath, fileName)

	// Create Senzing gRPC client.

	grpcConnection, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	testErr(err)
	szAbstractFactory := &szabstractfactory.Szabstractfactory{
		GrpcConnection: grpcConnection,
	}

	// Identify datasources and update Senzing configuration.

	dataSourceNames := getDataSources(inputFile)
	fmt.Printf("Found the following DATA_SOURCE values in the data: %v\n", dataSourceNames)
	err = updateSenzingConfig(szAbstractFactory, dataSourceNames)
	testErr(err)

	// Add records.

	err = addRecords(szAbstractFactory, inputFile)
	testErr(err)
}
