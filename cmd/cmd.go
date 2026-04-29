package cmd

var Filter string

var jsonOut bool
var subJsonOut bool
var baseJsonOut bool
var baseDebug bool
var subDebug bool

var subDsnFlag string
var baseDsnFlag string
var subHost string
var baseHost string
var subPort int
var basePort int
var debug bool

const defaultRedisPort = 6379
const defaultRedisHost = "127.0.0.1"
const defaultFilter = ".*"
