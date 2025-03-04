package wsapiclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"

	vmx "github.com/johlandabee/govmx"
)

// GetVM Auxiliar function to get the data of the VM and don't repeat code
// Input: c: pointer at the client of the API server, i: string with the ID yo VM
func GetVM(c *Client, i string) (*MyVm, error) {
	var vms []MyVm
	var vm MyVm
	// If you want see the path of the VM it's necessary getting all VMs
	// because the API of VmWare Workstation doesn't permit see this the another way
	response, err := c.httpRequest("vms", "GET", bytes.Buffer{})
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapitools.go Fu: GetVM Message: The request at the server API failed %s", err)
		return nil, err
	}
	err = json.NewDecoder(response).Decode(&vms)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: GetVM Message: I can't read the json structure %s", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: GetVM Obj: List of VMs %#v\n", vms)
	for tempvm, value := range vms {
		if value.IdVM == i {
			vm = vms[tempvm]
		}
	}

	vm.Denomination, err = GetDisplayName(vm.Path)
	if err != nil {
		return nil, err
	}
	vm.Description, err = GetAnnotation(vm.Path)
	if err != nil {
		return nil, err
	}
	response, err = c.httpRequest("vms/"+i, "GET", bytes.Buffer{})
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapitools.go Fu: GetVM Message: The request at the server API failed %s", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: GetVM Obj:Body of VM %#v\n", response)
	err = json.NewDecoder(response).Decode(&vm)
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapitools.go Fu: GetVM Message: I can't read the json structure %s", err)
		return nil, err
	}
	response, err = c.httpRequest("vms/"+i+"/power", "GET", bytes.Buffer{})
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapitools.go Fu: GetVM Message: The request at the server API failed %s", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: GetVM Obj:Body of power %#v\n", response)
	err = json.NewDecoder(response).Decode(&vm)
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapitools.go Fu: GetVM Message: I can't read the json structure %s", err)
		return nil, err
	}
	return &vm, nil
}

// GetVMFromFile - With this function we can obtain a vmx.VirtualMachine structure
// with all the possible values that we have in the file.
// Input: p: string, the complete path of the vxm file that we want to read
// Output: string, vmx.VirtualMachine structure, and error if you obtain some error in the function
func GetVMFromFile(p string) (vmx.VirtualMachine, error) {
	vm := new(vmx.VirtualMachine)
	data, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: GetVMFromFile Message: Failed %s, please make sure the config file exists", err)
		return *vm, err
	}

	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: GetVMFromFile Obj: Data File %#v\n", string(data))
	err = vmx.Unmarshal(data, vm)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: GetVMFromFile Obj: %#v", err)
		return *vm, err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: GetVMFromFile Obj: VM  %#v\n", vm)
	return *vm, nil
}

// SetVMToFile - With this function we can save a vmx.VirtualMachine structure
// with all the possible values that we have in the file.
// Input: p: string, with the parameter we want to change
// Output: error if you obtain some error in the function
func SetVMToFile(vm vmx.VirtualMachine, p string) error {
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetVMToFile Message: parameters %#v, %#v", vm, p)
	data, err := vmx.Marshal(vm)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: SetVMToFile Message: Failed to save the VMX structure in memory %s", err)
		return err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetVMToFile Obj: Data after read vm %#v\n", string(data))
	err = ioutil.WriteFile(p, data, 0644)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: SetVMToFile Message: Failed writing in file %s, please make sure the config file exists", err)
		return err
	}
	return err
}

// GetAnnotation - With this function we can obtain the value of the description of VM
// Input: p: string, the complete path of the vxm file that we want to read
// Output: string, Value of the Annotation field of the VM, error if you obtain some error in the fuction
func GetAnnotation(p string) (string, error) {
	vm, err := GetVMFromFile(p)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: GetAnnotation Message: Failure to obtain the value of the Description %s", err)
		return "", err
	}
	return vm.Annotation, nil
}

// SetAnnotation - With this function we can set the value of the description of VM
// Input: p: string, the complete path of the vxm file that we want to read
// v: string with the value of Annotation field
// Output: error if you obtain some error in the fuction
func SetAnnotation(p string, v string) error {
	vm, err := GetVMFromFile(p)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: SetAnnotation Message: We can't obtain the vmx object %s", err)
		return err
	}
	vm.Annotation = v
	err = SetVMToFile(vm, p)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: SetAnnotation Message: We haven't be able to save the structure in the file %s", err)
		return err
	}
	return nil
}

// GetDisplayName - With this function we can obtain the value of the name of VM
// Input: p: string, the complete path of the vxm file that we want to read
// Output: string, Value of the Denomination field of the VM, error if you obtain some error in the fuction
func GetDisplayName(p string) (string, error) {
	vm, err := GetVMFromFile(p)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: GetDisplayName Message: Failure to obtain the value of the Denomination %s", err)
		return "", err
	}
	return vm.DisplayName, nil
}

// SetAnnotation - With this function we can set the value of the denomination of VM
// Input: p: string, the complete path of the vxm file that we want to read
// v: string with the value of Denomination field, WARNING this function don't change teh PATH
// Output: error if you obtain some error in the fuction
func SetDisplayName(p string, v string) error {
	vm, err := GetVMFromFile(p)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: SetAnnotation Message: We can't obtain the vmx object %s", err)
		return err
	}
	vm.DisplayName = v
	err = SetVMToFile(vm, p)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: SetAnnotation Message: We haven't be able to save the structure in the file %s", err)
		return err
	}
	return nil
}

// SetNameDescription With this function you can setting the Denomination and Description of the VM.
// this information is in the vmx file of the machine for that you need know
// which is the file of the vm. Input: p: string with the complete path of the file,
// n: string with the denomination, d: string with the description err: variable with error if occur
func SetNameDescription(p string, n string, d string) error {
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Message: parameters %#v, %#v, %#v", p, n, d)
	data, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Message: Failed opening file %s, please make sure the config file exists", err)
		return err
	}

	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Obj: File object %#v\n", string(data))

	vm := new(vmx.VirtualMachine)
	err = vmx.Unmarshal(data, vm)
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Obj: %#v", err)
		return err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Obj: VM %#v\n", vm)

	vm.DisplayName = n
	vm.Annotation = d
	data, err = vmx.Marshal(vm)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: SetNameDescription Message: Failed to save the VMX structure in memory %s", err)
		return err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Obj: Data File %#v\n", string(data))
	err = ioutil.WriteFile(p, data, 0644)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: SetNameDescription Message: Failed writing in file %s, please make sure the config file exists", err)
		return err
	}
	// en este punto tambien tienes que cambiar el nombre del fihero cuando se cambia la denominacion
	return err
}

// SetParameter With this function you can set the value of the parameter.
// this information is in the vmx file of the machine for that you need know
// which is the file of the vm. Input: i: string with the id of the VM,
// p: string with the name or param to set, v: string with the value of param err: variable with error if occur
func (c *Client) SetParameter(i string, p string, v string) error {
	requestBody := new(bytes.Buffer)
	request, err := json.Marshal(map[string]string{
		"name":  p,
		"value": v,
	})
	if err != nil {
		return err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetParameter Obj:Request %#v\n", request)
	requestBody.Write(request)
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetParameter Obj:Request Body %#v\n", requestBody.String())
	response, err := c.httpRequest("/vms/"+i+"/configparams", "PUT", *requestBody)
	if err != nil {
		return err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetParameter Obj:response raw %#v\n", response)
	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: SetParameter Obj:Response Error %#v\n", err)
		return err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetParameter Obj:Response Body %#v\n", responseBody.String())
	// err = json.NewDecoder(responseBody).Decode(&vm)
	if err != nil {
		return err
	}
	return err
}
