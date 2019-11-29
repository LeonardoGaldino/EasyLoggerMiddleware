package namingservice

// Operation Represents the type of Operation NamingService is receiving
type Operation int

var operationValues = []string{"REGISTER", "UNREGISTER", "QUERY", "UNKNOWN"}

const (
	// REGISTER is a NamingService type of Operation
	REGISTER Operation = iota
	// UNREGISTER is a NamingService type of Operation
	UNREGISTER
	// QUERY is a NamingService type of Operation
	QUERY
	// UNKNOWN is an NamingService unknown type of Operation
	UNKNOWN
)

func (op Operation) String() string {
	return operationValues[op]
}

// StringToOperation transforms a string to its corresponding Operation
func StringToOperation(op string) Operation {
	isValid := false
	for _, opValue := range operationValues {
		if op == opValue {
			isValid = true
			break
		}
	}
	if !isValid {
		return UNKNOWN
	}
	return map[string]Operation{
		"REGISTER":   REGISTER,
		"UNREGISTER": UNREGISTER,
		"QUERY":      QUERY,
		"UNKNOWN":    UNKNOWN,
	}[op]
}

// Result Represents the result of a NamingService operation
type Result int

var resultValues = []string{"OK", "NOT_FOUND", "ERROR"}

const (
	// OK is an NamingService Operation Result
	OK Result = iota
	// NOTFOUND is an NamingService Operation Result
	NOTFOUND
	// ERROR is an NamingService Operation Result
	ERROR
)

func (result Result) String() string {
	return resultValues[result]
}

// StringToResult transforms a string to its corresponding Result
func StringToResult(result string) Result {
	isValid := false
	for _, resValue := range resultValues {
		if result == resValue {
			isValid = true
			break
		}
	}
	if !isValid {
		return ERROR
	}
	return map[string]Result{
		"OK":        OK,
		"NOT_FOUND": NOTFOUND,
		"ERROR":     ERROR,
	}[result]
}

// RequestMessage is a struct for send data to NamingService
type RequestMessage struct {
	Op   Operation
	Data string
}

// ResponseMessage is a struct for receiving data from NamingService
type ResponseMessage struct {
	Res  Result
	Data string
}
