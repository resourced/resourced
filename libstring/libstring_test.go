package libstring

import (
	"encoding/json"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestReplaceTildeWithRoot(t *testing.T) {
	path := "~/resourced"
	toBeTested := strings.Replace(path, "~", "/root", 1)

	if toBeTested != "/root/resourced" {
		t.Errorf("~ is not expanded correctly. Path: %v", toBeTested)
	}
}

func TestExpandTildeAndEnv(t *testing.T) {
	toBeTested := ExpandTildeAndEnv("~/resourced")

	if runtime.GOOS == "darwin" {
		if !strings.HasPrefix(toBeTested, "/Users") {
			t.Errorf("~ is not expanded correctly. Path: %v", toBeTested)
		}
	}

	toBeTested = ExpandTildeAndEnv("$GOPATH/src/github.com/resourced/resourced/tests/data/script-reader/darwin-memory.py")
	gopath := os.Getenv("GOPATH")

	if !strings.HasPrefix(toBeTested, gopath) {
		t.Errorf("$GOPATH is not expanded correctly. Path: %v", toBeTested)
	}
}

func TestGeneratePassword(t *testing.T) {
	_, err := GeneratePassword(8)
	if err != nil {
		t.Errorf("Generating password should not fail. err: %v", err)
	}
}

func TestGetIP(t *testing.T) {
	goodAddress := "127.0.0.1:55555"
	badAddress := "tasty:cakes"

	goodIP := GetIP(goodAddress)
	if goodIP == nil {
		t.Error("Should be able to parse '%v'", goodAddress)
	}

	if goodIP.String() != strings.Split(goodAddress, ":")[0] {
		t.Error("goodIP.String() should be the same as split goodAddress")
	}

	badIP := GetIP(badAddress)
	if badIP != nil {
		t.Error("Should not be able to parse '%v'", badAddress)
	}
}

func TestCSVtoJSON(t *testing.T) {
	csv := `pxname,svname,qcur,qmax,scur,smax,slim,stot,bin,bout,dreq,dresp,ereq,econ,eresp,wretr,wredis,status,weight,act,bck,chkfail,chkdown,lastchg,downtime,qlimit,pid,iid,sid,throttle,lbtot,tracked,type,rate,rate_lim,rate_max,check_status,check_code,check_duration,hrsp_1xx,hrsp_2xx,hrsp_3xx,hrsp_4xx,hrsp_5xx,hrsp_other,hanafail,req_rate,req_rate_max,req_tot,cli_abrt,srv_abrt,comp_in,comp_out,comp_byp,comp_rsp,lastsess,last_chk,last_agt,qtime,ctime,rtime,ttime,
http-in,FRONTEND,,,1,44,100,1675024,446801740,46344342820,42517,0,15,,,,,OPEN,,,,,,,,,1,2,0,,,,0,2,0,37,,,,0,1232521,401385,64289,433,41020,,2,37,1698691,,,31985314033,16729870273,47036,656760,,,,,,,,
http-in,IPv4-direct,,,1,33,100,301589,49674092,1835207008,14230,0,4,,,,,OPEN,,,,,,,,,1,2,1,,,,3,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
http-in,IPv4-cached,,,0,32,100,908022,313412130,38738243976,142,0,0,,,,,OPEN,,,,,,,,,1,2,2,,,,3,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
http-in,IPv6-direct,,,0,31,100,398999,76560960,5631508745,27267,0,11,,,,,OPEN,,,,,,,,,1,2,3,,,,3,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
http-in,local,,,0,0,100,0,0,0,0,0,0,,,,,OPEN,,,,,,,,,1,2,4,,,,3,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
http-in,local-https,,,0,5,100,78690,7154558,139383091,878,0,0,,,,,OPEN,,,,,,,,,1,2,5,,,,3,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
www,www,0,0,0,20,20,947559,350334547,45102970884,,0,,0,124,0,0,UP,1,1,0,10,2,1848240,39,,1,3,1,,947287,,2,0,,37,L7OK,200,3,217,753577,130896,62511,5,0,0,,,,90231,1305,,,,,3,OK,,0,1,16,972,
www,bck,0,0,0,6,10,53,18178,137531,,0,,0,0,0,0,UP,1,0,1,0,0,3916052,0,,1,3,2,,41,,2,0,,8,L7OK,200,1,0,26,6,21,0,0,0,,,,0,0,,,,,1848243,OK,,239,1,1,305,
www,BACKEND,0,22,0,42,100,947832,350532539,45103149435,137,0,,2,124,0,0,UP,1,1,1,,0,3916052,0,,1,3,0,,947328,,1,0,,37,,,,0,753603,130902,62841,426,59,,,,,90244,1305,31142359176,16393468547,47036,500655,3,,,0,1,16,972,
git,www,0,0,0,2,2,8132,3939635,135252254,,0,,0,0,0,0,UP,1,1,0,10,2,1848240,39,,1,4,1,,8125,,2,0,,2,L7OK,200,4,0,3941,4168,23,0,0,0,,,,120,1,,,,,390,OK,,4,1,1846,2555,
git,bck,0,0,0,0,2,0,0,0,,0,,0,0,0,0,UP,1,0,1,0,0,3916052,0,,1,4,2,,0,,2,0,,0,L7OK,200,2,0,0,0,0,0,0,0,,,,0,0,,,,,-1,OK,,0,0,0,0,
git,BACKEND,0,2,0,4,2,8132,3942374,135252254,0,0,,0,0,0,0,UP,1,1,1,,0,3916052,0,,1,4,0,,8125,,1,0,,2,,,,0,3941,4168,23,0,0,,,,,120,1,42189251,13364582,0,1711,390,,,4,1,1846,2555,
demo,BACKEND,0,0,1,26,20,215236,30025083,1033600591,0,0,,7,0,0,0,UP,0,0,0,,0,3916052,0,,1,17,0,,0,,1,1,,17,,,,0,215018,208,0,7,2,,,,,33,0,800765606,323037144,0,154394,0,,,6,0,1,155,
`

	inJson, err := CSVtoJSON(csv)
	if err != nil {
		t.Fatalf("Failed to parse CSV to JSON. Error: %v", err)
	}

	if string(inJson) == "" {
		t.Fatalf("JSON cannot be empty. JSON: %v", inJson)
	}

	var data []map[string]interface{}

	err = json.Unmarshal(inJson, &data)
	if err != nil {
		t.Fatalf("Invalid JSON. Error: %v", err)
	}
}
