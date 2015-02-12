package network

import("testing")


func TestGetHostIP(t *testing.T){
	// legg til ip for hosten
	ip_answer := "192.168.1.135"
	ip := GetHostIP()
	if(ip_answer != ip && ip != "is_offline"){
		t.Error("Got ip: "+ip)
	}
}

// func TestInitNetwork(t *testing.T){
// 	ip := "192.168.1"
// 	//Tester intialize 
// 	//-finne den lokale ip-addressen
// 	//-skal legge inn en broadcastadresse inn i dictonary
// 	//-skal leggge til
// 	//-Skal få inn en liste med alle TCP connections som er mulig og koble til.


// }	

// func TestPing(t* testing){
// 	//Testen fjerner en fra lista  
// }

