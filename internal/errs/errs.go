package errs

import "fmt"

var UserNotFoundError = fmt.Errorf("no such user")
var ValidationError = fmt.Errorf("validation failed")
var AccessError = fmt.Errorf("access forbiden")
var AdNotFoundError = fmt.Errorf("no such ad")
var WrongProtoBufDataError = fmt.Errorf("wrong field given")
