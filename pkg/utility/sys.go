package utility

import (
	"fmt"
	"runtime"
	"strings"
)

func GetFrameData(skip int) (scope, caller, callee, location string) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", "unknown", "unknown", "unknown"
	}

	funcName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(funcName, ".")

	if len(parts) < 2 {
		return "unknown", "unknown", "unknown", "unknown"
	}

	// Extract scope (struct name if method, package name if function)
	scopeParts := strings.Split(parts[0], "/")
	scope = scopeParts[len(scopeParts)-1]

	// Extract method/function name
	caller = parts[len(parts)-1]

	// For the callee, we'll use the function name
	if len(parts) >= 3 {
		callee = parts[len(parts)-2]
	} else {
		callee = caller
	}

	location = fmt.Sprintf("%s:%d", file, line)

	fmt.Println(parts)
	fmt.Println(location)

	return
}
