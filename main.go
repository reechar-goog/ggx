package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	directory "google.golang.org/api/admin/directory/v1"
)

func toYamlFormat(member *directory.Member) string {
	// var result string
	switch typ := member.Type; typ {
	case "GROUP":
		return "group:" + member.Email
	case "USER":
		if strings.HasSuffix(member.Email, "gserviceaccount.com") {
			return "serviceAccount:" + member.Email
		}
		return "user:" + member.Email
	}
	// if(member.Type =="GROUP"){

	// }
	// result = strings
	return "UH OH"
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// func main() {
// 	b, err := ioutil.ReadFile("credentials.json")
// 	if err != nil {
// 		log.Fatalf("Unable to read client secret file: %v", err)
// 	}

// 	// If modifying these scopes, delete your previously saved token.json.
// 	config, err := google.ConfigFromJSON(b, directory.AdminDirectoryUserReadonlyScope, directory.AdminDirectoryGroupReadonlyScope, directory.AdminDirectoryGroupMemberReadonlyScope)
// 	if err != nil {
// 		log.Fatalf("Unable to parse client secret file to config: %v", err)
// 	}
// 	client := getClient(config)

// 	srv, err := directory.New(client)
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve directory Client %v", err)
// 	}

// 	groups, err := srv.Groups.List().Domain("reechar.co").Do()
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve groups %v", err)
// 	}

// 	for _, group := range groups.Groups {
// 		fmt.Println(group.Name)
// 		g, err := srv.Members.List(group.Id).Do()
// 		if err != nil {
// 			// log.Fatalln("couldn't get")
// 			log.Fatalf("Unable to retrieve groups %v", err)
// 		}
// 		for _, gmember := range g.Members {
// 			// fmt.Println(gmember.Type + " " + gmember.Email)
// 			fmt.Println(toYamlFormat(gmember))
// 		}
// 	}

// 	r, err := srv.Users.List().Customer("my_customer").MaxResults(10).
// 		OrderBy("email").Do()
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve users in domain: %v", err)
// 	}

// 	if len(r.Users) == 0 {
// 		fmt.Print("No users found.\n")
// 	} else {
// 		fmt.Print("Users:\n")
// 		for _, u := range r.Users {
// 			fmt.Printf("%s (%s)\n", u.PrimaryEmail, u.Name.FullName)
// 		}
// 	}
// }

//Policy is for iam policy
type Policy struct {
	Bindings []*Binding `yaml:"Bindings"`
}

//Binding for binding
type Binding struct {
	Members []string `yaml:"members"`
	Role    string   `yaml:"role"`
}

func main() {
	// info, err := os.Stdin.Stat()
	// if err != nil {
	// 	panic(err)
	// }

	// if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
	// 	fmt.Println("The command is intended to work with pipes.")
	// 	fmt.Println("Usage: fortune | gocowsay")
	// 	return
	// }

	reader := bufio.NewReader(os.Stdin)
	// var output []rune
	input, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	// input, err := reader.ReadBytes()
	c := Policy{}
	err = yaml.Unmarshal(input, &c)

	// for _, x := range c.Bindings {
	// 	fmt.Println(x.Role + " " + x.Members[0])
	// }
	groupFiltered := Policy{}
	newBindings := make([]*Binding, 0)
	for _, binding := range c.Bindings {
		tempBind := new(Binding)
		tempBind.Role = binding.Role
		_ = buildMapOfSources(binding.Members)
		for i, member := range binding.Members {

			if strings.HasPrefix(member, "group") {
				binding.Members[i] = member + " #THIS IS A GROUP"
				// binding.Members = append(binding.Members, "")
				// copy(binding.Members[i+1:], binding.Members[i:])
				// binding.Members[i] = "INSERTED MEMBER"
				binding.Members = insertStringToSlice(binding.Members, "insertedMember", i+1)
				// i++
			}
		}
		tempBind.Members = binding.Members
		// tempBind.Members[0] = tempBind.Members[0] + " #from group1"
		newBindings = append(newBindings, tempBind)
	}
	// c.Bindings = newBindings
	groupFiltered.Bindings = newBindings
	// for _, binding := range groupFiltered.Bindings {
	// 	fmt.Println(binding.Role)
	// }
	yamlbob, err := yaml.Marshal(groupFiltered)
	removeDash := string(yamlbob)
	removeDash = strings.Replace(removeDash, "'", "", -1)
	fmt.Println(removeDash)
	// fmt.Println(string(yamlbob))

	// yamlbob2, err := yaml.Marshal()
	// fmt.Println(string(yamlbob2))
	// fmt.Println(string(input))

	// for {
	// 	input, _, err := reader.ReadRune()
	// 	if err != nil && err == io.EOF {
	// 		break
	// 	}
	// 	output = append(output, input)
	// }

	// for j := 0; j < len(output); j++ {
	// 	fmt.Printf("%c", output[j])
	// }
}

func insertStringToSlice(slice []string, target string, index int) []string {
	result := make([]string, index+1)
	copy(result, slice[:index])
	result[index] = target
	result = append(result, slice[index:]...)
	return result
}

func buildMapOfSources(members []string) map[string][]string {
	results := make(map[string][]string)
	workQueue := make([]string, 0)

	results, workQueue = process(members, workQueue, results)

	for len(workQueue) > 0 {
		fmt.Println("hello")
		fmt.Printf("Queue: %v", workQueue)
		pop := workQueue[0]
		newMembers := getMembers(strings.Replace(pop, "group:", "", -1))
		results, workQueue = process(newMembers, workQueue, results)
		workQueue = workQueue[1:]
	}

	return results
}

func process(members []string, workQueue []string, results map[string][]string) (map[string][]string, []string) {
	fmt.Printf("Processing: %v\n", members)
	for _, member := range members {
		switch {
		case strings.HasPrefix(member, "user"):
			results[member] = []string{"policy"}
		case strings.HasPrefix(member, "serviceAccount"):
			results[member] = []string{"policy"}
		case strings.HasPrefix(member, "group"):
			workQueue = append(workQueue, member)
		}
	}
	return results, workQueue
	// return results, workQueue
}

func getMembers(groupName string) []string {
	fmt.Printf("api getting: %v\n", groupName)
	results := make([]string, 0)
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, directory.AdminDirectoryUserReadonlyScope, directory.AdminDirectoryGroupReadonlyScope, directory.AdminDirectoryGroupMemberReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := directory.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve directory Client %v", err)
	}
	g, err := srv.Members.List(groupName).Do()
	for _, gmember := range g.Members {
		// fmt.Println(gmember.Type + " " + gmember.Email)
		fmt.Println(toYamlFormat(gmember))
		results = append(results, toYamlFormat(gmember))
	}

	return results
}
