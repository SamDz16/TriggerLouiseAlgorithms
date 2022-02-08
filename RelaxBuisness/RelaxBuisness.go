package RelaxBuisness

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"unsafe"
)

type Query struct {
	query   string
	parents []string
}

// START : SPECIAL METHODS
// ###########################################################################################################################################################################

func TpGenerateLevelTripplePatterns(triplePatterns []string, level int) []string {
	combinations := []string{}
	n := len(triplePatterns)

	indexes := []int{}
	for i := 0; i < level; i++ {
		indexes = append(indexes, i)
	}

	//  liste_combinaisons.append(tuple([e[index] for index in indices]))
	triple := ""
	for i, index := range indexes {
		if i == len(indexes)-1 {
			triple += triplePatterns[index]
		} else {
			triple += triplePatterns[index] + " . "
		}
	}
	combinations = append(combinations, triple)

	if level == 0 {
		return []string{" "}
	}

	if level == n {
		return combinations
	}

	i := level - 1

	for i != -1 {
		indexes[i] += 1

		for j := i + 1; j < level; j++ {
			indexes[j] = indexes[j-1] + 1
		}

		if indexes[i] == (n - level + i) {
			i -= 1
		} else {
			i = level - 1
		}

		temp := []string{}
		for _, index := range indexes {
			temp = append(temp, triplePatterns[index])
		}

		// fmt.Println("temp", temp)

		triple := ""
		for i, t := range temp {
			if i == len(temp)-1 {
				triple += t
			} else {
				triple += t + " . "
			}
		}

		combinations = append(combinations, triple)
	}
	return combinations
}

func TpMakeQueries(tripplePatterns []string) []string {
	var res []string

	for _, t := range tripplePatterns {
		res = append(res, "select * where {"+t+"}")
	}
	return res
}

func TpMakeLattice(q string) []interface{} {

	var initialQuery Query
	initialQuery.query = q

	// Get all the tripple patterns
	tripplePatterns := GetQueryTripplePatterns(initialQuery)

	// triplePatternsNbr will contain the length of all the individual triple patterns of the initialQuery
	triplePatternsNbr := len(tripplePatterns)

	level := triplePatternsNbr
	var allTripplePatterns []string

	for i := 0; i < triplePatternsNbr+1; i++ {
		var temp []string = GenerateLevelTripplePatterns(tripplePatterns, level)
		for _, t := range temp {
			allTripplePatterns = append(allTripplePatterns, t)
		}
		level--
	}

	res := TpMakeQueries(allTripplePatterns)

	var ret []interface{}

	for _, q := range res {
		ret = append(ret, q)
	}

	return ret
}

// END : SPECIAL METHODS
// ###########################################################################################################################################################################

// START : UTILITY FUNCTIONS
// ###########################################################################################################################################################################

func IntToByte(num int) byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr[0]
}

// END : UTILITY FUNCTIONS
// ###########################################################################################################################################################################

// START : BASE ALGORITHM FUNCTIONS
// ###########################################################################################################################################################################

func ExecuteSPARQLQuery(requestUrl string, sparqlQuery string) []byte {
	// Make the HTTP request
	// res, err := http.DefaultClient.Get(requestUrl)
	data := url.Values{}
	data.Set("query", sparqlQuery)

	u, _ := url.ParseRequestURI(requestUrl)
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(r)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return []byte{}
	}

	// Read the response body
	dataBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return dataBody
}

func GetQueryTripplePatterns(initialQuery Query) []string {
	triplePatternsStr := ""
	start := false
	for _, ch := range initialQuery.query {

		char := string(ch)

		if char == "{" || char == "}" {
			if char == "{" {
				start = true
				continue
			} else {
				start = false
				continue
			}
		}

		if start {
			triplePatternsStr += char
		}
	}

	return strings.Split(triplePatternsStr, " . ")
}

func GenerateLevelTripplePatterns(triplePatterns []string, level int) []string {
	combinations := []string{}
	n := len(triplePatterns)

	indexes := []int{}
	for i := 0; i < level; i++ {
		indexes = append(indexes, i)
	}

	//  liste_combinaisons.append(tuple([e[index] for index in indices]))
	triple := ""
	for i, index := range indexes {
		if i == len(indexes)-1 {
			triple += triplePatterns[index]
		} else {
			triple += triplePatterns[index] + " . "
		}
	}
	combinations = append(combinations, triple)

	if level == 0 {
		return []string{" "}
	}

	if level == n {
		return combinations
	}

	i := level - 1

	for i != -1 {
		indexes[i] += 1

		for j := i + 1; j < level; j++ {
			indexes[j] = indexes[j-1] + 1
		}

		if indexes[i] == (n - level + i) {
			i -= 1
		} else {
			i = level - 1
		}

		temp := []string{}
		for _, index := range indexes {
			temp = append(temp, triplePatterns[index])
		}

		// fmt.Println("temp", temp)

		triple := ""
		for i, t := range temp {
			if i == len(temp)-1 {
				triple += t
			} else {
				triple += t + " . "
			}
		}

		combinations = append(combinations, triple)
	}
	return combinations
}

func MakeQueries(tripplePatterns []string, queries *[]Query) {
	for _, t := range tripplePatterns {
		var q Query
		q.query = "select * where {" + t + "}"
		*queries = append(*queries, q)
	}
}

func MakeLattice(initialQuery Query, queries *[]Query) {

	// Get all the tripple patterns
	tripplePatterns := GetQueryTripplePatterns(initialQuery)

	// triplePatternsNbr will contain the length of all the individual triple patterns of the initialQuery
	triplePatternsNbr := len(tripplePatterns)

	level := triplePatternsNbr
	var allTripplePatterns []string

	for i := 0; i < triplePatternsNbr+1; i++ {
		var temp []string = GenerateLevelTripplePatterns(tripplePatterns, level)
		for _, t := range temp {
			// q.parents = []string{}

			allTripplePatterns = append(allTripplePatterns, t)
		}
		level--
	}

	MakeQueries(allTripplePatterns, queries)
}

func IsDirectParent(q1, q2 Query) bool {
	tp1 := GetQueryTripplePatterns(q1)
	tp2 := GetQueryTripplePatterns(q2)

	if len(tp1) == len(tp2)+1 {
		matches := 0

		for _, q2 := range tp2 {
			for _, q1 := range tp1 {
				if q2 == q1 {
					matches++
					break
				}
			}
		}
		if matches == len(tp2) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func SetSuperQueries(queries *[]Query) {
	qs := *queries

	for i := 0; i < len(*queries); i++ {
		if i == 0 {
			// The first element has no parents - root
			qs[i].parents = []string{}
		} else if i == len(qs)-1 {
			// The last element is the empty query
			for _, qTemp := range GetQueryTripplePatterns(qs[0]) {
				qs[i].parents = append(qs[i].parents, qTemp)
			}
		} else {
			//  The rest of the elements : {1, 2, ..., len(queries)-2}
			for j := 0; j < i; j++ {
				res := IsDirectParent(qs[j], qs[i])
				// fmt.Printf("j = %v - %v - %v - %v\n", j, queries[j], q, res)
				if res {
					// fmt.Println("parents", parents)
					qs[i].parents = append(qs[i].parents, strings.Join(GetQueryTripplePatterns(qs[j])[:], " . "))

					// fmt.Println(true)
					// q.parents = append(q.parents, "true")
					// q.parents = GetQueryTripplePatterns(qs[j])
				}
			}
		}
	}
}

func ContainsKey(queries *map[*Query]bool, q Query) bool {

	for k := range *queries {
		qTemp := *k
		if qTemp.query == q.query {
			return true
		}
	}
	return false
}

func FindQuery(queries []Query, query Query) (int, bool) {
	for i, q := range queries {
		if q.query == query.query {
			return i, true
		}
	}
	return -1, false
}

func RemoveQuery(listQueries []Query, index int) []Query {
	var newListQueries []Query
	for j, q := range listQueries {
		if j != index {
			newListQueries = append(newListQueries, q)
		}
	}

	return newListQueries
}

type Sparql struct {
	Results Results `json:"results"`
}

type Results struct {
	Bindings []Bindings `json:"bindings"`
}

type Bindings struct {
	Fp TP `json:"fp"`
	A  TP `json:"a"`
	N  TP `json:"n"`
	C  TP `json:"c"`
}

type TP struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func TpExecuteSPARQLQuery(requestUrl string, sparqlQuery string) int {
	// Make the HTTP request
	// res, err := http.DefaultClient.Get(requestUrl)
	data := url.Values{}
	data.Set("query", sparqlQuery)

	u, _ := url.ParseRequestURI(requestUrl)
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(r)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return -1
	}

	// Read the response body
	dataBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var results Sparql
	err = json.Unmarshal([]byte(dataBody), &results)

	if err != nil {
		return -1
	}

	return len(results.Results.Bindings)
}

func Base(q string, K byte, NBs []byte) ([]string, []string, int) {
	// Initialisations
	// ##################################################################################################################################################################### //
	// ##########################################################           INITIALIZE ALGO         ######################################################################## //
	// ##################################################################################################################################################################### //

	var initialQuery Query
	initialQuery.query = q

	// List Queries
	var listQueries []Query

	// Executed Queries : contains for each qury, the number of the results
	var executedQueries map[*Query]byte = make(map[*Query]byte)

	// List FIS : all rthe queries that fail
	var listFIS map[*Query]bool = make(map[*Query]bool)

	listXSS := &[]Query{}
	listMFIS := &[]Query{}

	// ##################################################################################################################################################################### //
	// ##########################################################              RUN ALGO             ######################################################################## //
	// ##################################################################################################################################################################### //

	MakeLattice(initialQuery, &listQueries)

	SetSuperQueries(&listQueries)
	var s int = 0

	for len(listQueries) != 0 {

		// First element of the list
		qTemp := listQueries[0]

		// Remove the first element from the list
		listQueries = listQueries[1:]

		var Nb byte
		fmt.Println(qTemp.query)

		if qTemp.query != "select * where { }" {
			Nb = NBs[s]
			s++

			// add qTemp to executedQueries list with the number Nb of the results
			executedQueries[&qTemp] = Nb
		} else {
			fmt.Println("empty query")
			Nb = 1
		}

		// Get Direct Super Queries Of 'qTemp'
		var superQueries []Query
		MakeQueries(qTemp.parents, &superQueries)

		parentsFIS := true

		i := 0

		for parentsFIS && i < len(superQueries) {
			superQuery := superQueries[i]
			if !ContainsKey(&listFIS, superQuery) {
				parentsFIS = false
			}
			i++
		}

		if Nb > K {
			// Query qTemp fails
			if parentsFIS {
				// We remove all the superqueries of qTemp from listMFIS list
				for _, qSQ := range superQueries {
					index, found := FindQuery(*listMFIS, qSQ)
					if found {
						*listMFIS = RemoveQuery(*listMFIS, index)
					}
				}

				// Since the request qTemp has failed, we add it to the list of FIS
				listFIS[&qTemp] = true

				// qTemps is the new MFIS
				*listMFIS = append(*listMFIS, qTemp)
			}
		} else {
			// qTemp has succeded
			if parentsFIS && qTemp.query != " " {
				*listXSS = append(*listXSS, qTemp)
			}
		}

	}

	var xss []string
	for _, x := range *listXSS {
		xss = append(xss, x.query)
	}

	var mfis []string
	for _, m := range *listMFIS {
		mfis = append(mfis, m.query)
	}

	return xss, mfis, len(executedQueries)
}

// END : BASE ALGORITHM FUNCTIONS
// ###########################################################################################################################################################################
