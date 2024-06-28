package main

import (
	"fmt"
	"encoding/json"
	"os"
	"io/ioutil"
	"sync"
	"path/filepath"
	"github.com/jcelliott/lumber"
)

const version = "1.0.0"

type (
	Logger interface{
		Debug(string, ...interface{})
		Info(string, ...interface{})
		Warn(string, ...interface{})
		Error(string, ...interface{})
		Fatal(string, ...interface{})
    Trace(string, ...interface{})
	}
	
	Driver struct{
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

type Options struct{
	Logger
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)
	opts := Options{}
	
	if options != nil{
		opts = *options
	}
	
	if opts.Logger == nil{
		opts.Logger = lumber.NewConsoleLogger((lumber.DEBUG))
	}

	driver := Driver{
		dir: dir,
		mutexes: make(map[string]*sync.Mutex),
		log: opts.Logger,
	}

	if _, err := os.Stat(dir); err == nil{
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return &driver, nil
	}
	opts.Logger.Debug("Creating the database at '%s'....\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

func (d *Driver) Write(collection, resource string , v interface{}) error {
	
	if collection == ""{
		return fmt.Errorf("Missing Collection - No place to save records\n")
	}
	
	if resource == ""{
		return fmt.Errorf("Missing Resource - No place to save rescords(No name)\n")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()
	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource+".json")
	tmpPath := fnlPath + ".tmp"
	
  if err := os.MkdirAll(dir, 0755); err != nil{
		return err
	}

  b, err := json.MarshalIndent(v, "", "\t")
	if err != nil{
		return err
	}
	
	b = append(b, byte('\n'))
	
  if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil{
		return err
	}
	return os.Rename(tmpPath, fnlPath)
}

func (d *Driver) Read(collection, resource string, v interface{}) error {
  if collection == ""{
		return fmt.Errorf("Unable to read records\n")
	}
	
	if resource == ""{
		return fmt.Errorf("Missing Resource - No place to read records(No name)\n")
	}

	record := filepath.Join(d.dir, collection, resource)
	if _, err := Stat(record); err != nil{
		return err
	}

	b, err := ioutil.ReadFile(record + ".json")
  if err != nil{
    return err
  }
  
  return json.Unmarshal(b, &v)
}

func (d *Driver) ReadAll(collection string)([]string, error) {
	if collection == ""{
		return nil, fmt.Errorf("Unable to read records\n")
	}
	
  dir := filepath.Join(d.dir, collection)
	
	if _, err := Stat(dir); err != nil{
		return nil, err
	}
	
	files, _ := ioutil.ReadDir(dir)
  var records []string
	for _, file := range files{
		b, err :=ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil{
			return nil, err
		}
		records = append(records, string(b))
	}
	return records, nil
}

func (d *Driver)Delete(collection, resource string) error {
	path := filepath.Join(collection, resource)
	mutex := d.getOrCreateMutex(collection)
  mutex.Lock()
	mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err := Stat(dir);{
	case fi == nil, err!= nil:
		return fmt.Errorf("Unable to find file %v\n", path)
		
	case fi.Mode().IsDir():
		return os.RemoveAll(dir)

	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}
	return nil
}

func (d *Driver) getOrCreateMutex(collection string)*sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, ok := d.mutexes[collection]
	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}
	return m
}

func Stat(path string) (fi os.FileInfo,err error){
	if fi, err = os.Stat(path); os.IsNotExist(err){
		fi, err = os.Stat(path + ".json")
	}
	return
}

type Address struct{
	City    string
	State   string
	Country string
	Pincode json.Number
}

type User struct{
	Name    string
	Age     json.Number
	Contact string
  Company string
	Address Address
}

func main() {
	dir := "./"
	db, err := New(dir, nil)
	if err != nil{
	  fmt.Println("Error: ", err)
	}
	employees := []User{
		{"John", "25", "9876543210", "Google", 
Address {"New york", "New york", "USA", "12345"}},
		{"Brian", "30", "98651891093", "Microsoft", Address{"Texas", "Texas", "USA", "14325"}},
		{"Alice", "28", "9845623178", "Apple", Address{"California", "San Francisco", "USA", "94101"}},
{"Sarah", "35", "9832145678", "Amazon", Address{"Washington", "Seattle", "USA", "98101"}},
{"Michael", "40", "9812345678", "Facebook", Address{"California", "Menlo Park", "USA", "94025"}},
{"Jessica", "32", "9823456789", "Netflix", Address{"California", "Los Gatos", "USA", "95032"}},
	}

	for _, value := range employees{
		db.Write("users", value.Name, User{
			Name: value.Name,
		  Age: value.Age,
			Contact: value.Contact,
			Company: value.Company,
			Address: value.Address,
		})
	}

  records, err := db.ReadAll("users")
	if err != nil{
		fmt.Println("Error: ", err)
	}
	fmt.Println(records)

  allUsers := []User{}
	for _, f := range records{
		employeeFound := User{}
		if err := json.Unmarshal([]byte(f), &employeeFound); err != nil{
			fmt.Println("Error: ", err)
		}
		allUsers = append(allUsers, employeeFound)
	}
	fmt.Println((allUsers))

	/*if err := db.Delete("users", "John"); err != nil{
		fmt.Println("Error", err)
	}*/
	
}
 