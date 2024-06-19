package solver

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ConstructResponse struct {
	DataToSign string `json:"dataToSign"`
}

// Construct the solution and return the data to sign
func Construct(
	solver string,
	intent *([]byte),
	operationIndex int,
) (string, error) {
	operationIndexStr := strconv.FormatUint(uint64(operationIndex), 10)
	req, err := http.NewRequest("POST", solver+"/construct?operationIndex="+operationIndexStr, bytes.NewBuffer(*intent))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return "", errors.New(string(body))
	}

	var response ConstructResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.DataToSign, nil
}

type SolveResponse struct {
	Result string `json:"result"`
}

// Solve the intent with the signature and return the result
func Solve(
	solver string,
	intent *([]byte),
	operationIndex int,
	signature string,
) (string, error) {
	operationIndexStr := strconv.FormatUint(uint64(operationIndex), 10)
	req, err := http.NewRequest("POST", solver+"/solve?operationIndex="+operationIndexStr+"&signature="+signature, bytes.NewBuffer(*intent))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return "", errors.New(string(body))
	}

	var response SolveResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Result, nil
}

type StatusResponse struct {
	Status string `json:"status"` // "pending", "success", "failure"
}

const (
	SOLVER_OPERATION_STATUS_PENDING = "pending"
	SOLVER_OPERATION_STATUS_SUCCESS = "success"
	SOLVER_OPERATION_STATUS_FAILURE = "failure"
)

// Check the status of the operation
func CheckStatus(
	solver string,
	intent *([]byte),
	operationIndex int,
) (string, error) {
	operationIndexStr := strconv.FormatUint(uint64(operationIndex), 10)
	req, err := http.NewRequest("POST", solver+"/status?operationIndex="+operationIndexStr, bytes.NewBuffer(*intent))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return "", errors.New(string(body))
	}

	var response StatusResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Status, nil
}

type OutputResponse struct {
	Output string `json:"output"`
}

func GetOutput(
	solver string,
	intent *([]byte),
	operationIndex int,
) (string, error) {
	operationIndexStr := strconv.FormatUint(uint64(operationIndex), 10)
	req, err := http.NewRequest("POST", solver+"/output?operationIndex="+operationIndexStr, bytes.NewBuffer(*intent))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return "", errors.New(string(body))
	}

	var response OutputResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Output, nil
}
