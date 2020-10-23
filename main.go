package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type WorkOrderEntry struct {
	ResourceID          string
	RefId               string
	URI                 string
	ContainerIndicator1 string
	ContainerIndicator2 string
	ContainerIndicator3 string
	Title               string
	ComponentID         string
}

var workOrderPtr = flag.String("workorder", "digitization_work_order_report.tsv", "the location of the work order")

func main() {
	flag.Parse()

	//open the workorder as a slice of WorkOrderEntries
	workOrder, err := openWorkOrder(*workOrderPtr)
	if err != nil {
		panic(err)
	}

	//get the name of the directory from the first line of the workorder
	directoryName := strings.Replace(workOrder[1].ResourceID, ".", "-", 1)

	//create the root directory
	err = os.Mkdir(directoryName, 0755)
	if err != nil {
		panic(err)
	}

	//create the metadata directory
	metadataDir := filepath.Join(directoryName, "metadata")
	err = os.Mkdir(metadataDir, 0755)
	if err != nil {
		panic(err)
	}

	//copy the work order to the metadata directory
	err = CopyWorkOrder(*workOrderPtr, metadataDir)
	if err != nil {
		panic(err)
	}

	//create cuid directories
	for _, entry := range workOrder {
		subdir := filepath.Join(directoryName, entry.ComponentID)
		err := os.Mkdir(subdir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func CopyWorkOrder(workorder string, mdLocation string) error {
	wo, err := ioutil.ReadFile(workorder)
	if err != nil {
		return err
	}

	wo2, err := os.Create(filepath.Join(mdLocation, "digitization_work_order_report.tsv"))
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(wo2)
	writer.Write(wo)
	writer.Flush()
	wo2.Close()

	return nil
}

func openWorkOrder(fileLoc string) ([]WorkOrderEntry, error) {
	var workOrder = []WorkOrderEntry{}
	workOrderFile, err := os.Open(fileLoc)
	if err != nil {
		return workOrder, err
	}

	scanner := bufio.NewScanner(workOrderFile)
	scanner.Scan() // skip the header

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "\t")
		workOrder = append(workOrder, WorkOrderEntry{
			line[0], line[1], line[2], line[3], line[4], line[5], line[6], line[7],
		})
	}

	if scanner.Err() != nil {
		return workOrder, err
	}

	return workOrder, nil
}
