package transport

import (
	"encoding/json"
	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/security"
	"testing"
)

func TestSignMessageES256(t *testing.T) {
	keys := security.NewEcdsaKey()
	if err:=keys.Generate();err !=nil {
		t.Error("SerializedKey generation error",err)
		t.FailNow()
	}
	private , pub := keys.ExportX509EncodedKeys()
	t.Log(private)
	t.Log(pub)
	msg := fimpgo.NewFloatMessage("evt.sensor.report", "temp_sensor", 35.5, nil, nil, nil)
	got, err := SignMessageES256(msg, nil, "alex@gmail.com", keys, nil)
	if err !=nil {
		t.Error("Signing error",err)
		t.FailNow()
	}

	bmsg,err := got.SerializeToJson()
	if err = keys.Generate();err !=nil {
		t.Error("SerializedKey generation error",err)
		t.FailNow()
	}
	t.Log("X:",keys.PublicKey().X.Text(16))
	t.Log("Y:",keys.PublicKey().Y.Text(16))

	t.Log(string(bmsg))

	keys2 := security.NewEcdsaKey()
	err = keys2.ImportX509PublicKey(pub)
	if err != nil {
		t.Error("Wrong key")
		t.FailNow()
	}

	result,err := GetVerifiedMessageES256(got,keys2)
	if err != nil || result == nil {
		t.Error("Signature is not valid")
	}else {
		t.Log(result)
	}

}

func TestSignMessageES256_TestKey(t *testing.T) {
	keys := security.NewEcdsaKey()
	if err:=keys.Generate();err !=nil {
		t.Error("SerializedKey generation error",err)
		t.FailNow()
	}
	private , pub := keys.ExportJsonEncodedKeys()
	bprivate,err := json.Marshal(private)
	if err !=nil {
		t.Error("Serialize error",err)
		t.FailNow()
	}
	bpub,err := json.Marshal(pub)
	if err !=nil {
		t.Error("Serialize error",err)
		t.FailNow()
	}
	t.Log(string(bprivate))
	t.Log(string(bpub))
	msg := fimpgo.NewFloatMessage("evt.sensor.report", "temp_sensor", 35.5, nil, nil, nil)
	got, err := SignMessageES256(msg, nil, "alex@gmail.com", keys, nil)
	if err !=nil {
		t.Error("SerializedKey generation error",err)
		t.FailNow()
	}

	bmsg,err := got.SerializeToJson()
	if err !=nil {
		t.Error("Serialize error",err)
		t.FailNow()
	}

	t.Log(string(bmsg))

	keys2 := security.NewEcdsaKey()
	err = keys2.ImportJsonPublicKey(pub)
	if err != nil {
		t.Error("Wrong key")
		t.FailNow()
	}
	result,err := GetVerifiedMessageES256(got,keys2)
	if err != nil || result == nil {
		t.Error("Signature is not valid")
	}else {
		t.Log(result)
	}
}

func TestSignMessageES256_TestVerify(t *testing.T) {
	pubKey := security.JsonEcKey{
		X: "f149023bb33138b6bfc6458c50b21c4ed67004b25d8ae201a2fea3731a9da694",
		Y: "6ce19554fbb2104de755c01aeb247fc3fb99b122b3ade55bbb2332b3c9acdd66",
	}
	keys := security.NewEcdsaKey()
	err := keys.ImportJsonPublicKey(pubKey)
	if err != nil {
		t.Error("Keys import error")
		t.FailNow()
	}
	signedMsgTxt := "{\"type\":\"evt.transport.signed\",\"serv\":\"temp_sensor\",\"val_t\":\"bin\",\"val\":\"eyJjb3JpZCI6bnVsbCwiY3RpbWUiOiIyMDIwLTA1LTE1VDExOjI5OjQ4Ljc1NzY4MSIsInByb3BzIjpudWxsLCJzZXJ2IjoiZmhidXRsZXIiLCJ0YWdzIjpudWxsLCJ0eXBlIjoiY21kLmdhdGV3YXkuZ2V0X2FjdGl2ZV9hZGFwdGVycyIsInVpZCI6IjlmMjJlYzkwLTk2OGUtMTFlYS1kZTdmLWViYTk3ODNhZjQwMCIsInZhbF90Ijoic3RyaW5nIiwidmVyIjpudWxsLCJ2YWwiOiIiLCJyZXNwX3RvIjoicHQ6ajEvbXQ6cnNwL3J0OmNsb3VkL3JuOnJlbW90ZS1jbGllbnQvYWQ6c21hcnRob21lLWFwcCIsInNyYyI6ImFwcCJ9\",\"tags\":null,\"props\":{\"sig\":\"SCwiI0yRhv4vydSND-Khpi2uCkoSjLOHmdZeKnELmkMtZOnxCuVpMs1A9zNPfXBprL3xN4_n8WT__IM8kpjEhA\",\"user_id\":\"alex@gmail.com\"},\"ver\":\"1\",\"corid\":\"\",\"ctime\":\"2020-05-14T10:56:32.385+02:00\",\"uid\":\"6ad4ae68-7458-44a9-8cdc-fcc8551689e5\"}"
	signedMsg,err := fimpgo.NewMessageFromBytes([]byte(signedMsgTxt))
	t.Log("Signature:",signedMsg.Properties["sig"])
	if err != nil {
		t.Error("Wrong message")
		t.FailNow()
	}
	innnerMsg,err := GetVerifiedMessageES256(signedMsg,keys)
	if err != nil {
		t.Error("Message can't be verified. Err:",err)
		t.FailNow()
	}
	t.Logf("%+v",*innnerMsg)
}