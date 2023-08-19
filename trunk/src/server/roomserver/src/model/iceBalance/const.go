package iceBalance

const ICE_OK = 0
const ICE_NO_FOUND = 1
const ICE_FULL_LOAD = 2
const ICE_JSON_PARSE_FALIED = 3
const ICE_HAVE_REGISTERED = 4
var ICE_ERROR_MSG = []string{
	"successful",	// 0
	"ice_no_found",	// 1
	"ice_full_loading", // 2
	"ice_json_parse_falied",	// 3
	"ice_have_registered"}