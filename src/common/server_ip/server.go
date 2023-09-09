package server_ip

import (
	"errors"
	"net"

	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/src/common/http/api"
)

var lan_ip = ""
var pub_ip = ""

// GetLocalIpv4 get local IP address.
func GetLocalIpv4() (string, error) {
	if lan_ip != "" {
		return lan_ip, nil
	}
	localAddr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, v := range localAddr {
		if ipNet, ok := v.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				lan_ip = ipNet.IP.String()
				return lan_ip, nil
			}
		}
	}
	return "", errors.New("failed to obtain the local IP address")
}

func GetPublicIpv4() (string, error) {

	if pub_ip != "" {
		return pub_ip, nil
	}

	// //get
	ip, err := getPublicIpv4()
	if err != nil {
		return "", err
	}

	pub_ip = ip
	return pub_ip, nil
}

func getPublicIpv4() (string, error) {
	ip, err := getPublicIpv4_dnsService()
	if err == nil {
		return ip, nil
	}
	basic.Logger.Errorln("getPublicIpv4_dnsService error:", err)

	// back up
	ip, err = getPublicIpv4_backup_ipconfigme()
	if err == nil {
		return ip, nil
	}
	basic.Logger.Errorln("getPublicIpv4_backup_ipconfigme error:", err)

	ip, err = getPublicIpv4_backup_ipify()
	if err == nil {
		return ip, nil
	}
	basic.Logger.Errorln("getPublicIpv4_backup_ipify error:", err)

	ip, err = getPublicIpv4_backup_bigdatacloud()
	if err == nil {
		return ip, nil
	}
	basic.Logger.Errorln("getPublicIpv4_backup_bigdatacloud error:", err)

	return "", err
}

func getPublicIpv4_dnsService() (string, error) {
	type Msg_Resp_PublicIpv4 struct {
		api.API_META_STATUS
		Ip string `json:"ip"`
	}

	resp := &Msg_Resp_PublicIpv4{}
	err := api.Get_("https://api.dns.coreservice.io/api/common/public_ipv4", "", 10, resp)
	if err != nil {
		return "", err
	}
	if resp.Meta_status < 0 {
		return "", errors.New(resp.Meta_message)
	}

	return resp.Ip, nil
}

func getPublicIpv4_backup_ipify() (string, error) {
	type Msg_Resp_PublicIpv4 struct {
		Ip string `json:"ip"`
	}

	resp := &Msg_Resp_PublicIpv4{}

	// try from public api
	err := api.Get_("https://api.ipify.org?format=json", "", 10, resp)
	if err != nil {
		return "", err
	}

	netip := net.ParseIP(resp.Ip)
	if netip == nil {
		return "", errors.New("ip not find")
	}
	if netip.To4() == nil {
		return "", errors.New("ip not find")
	}

	return resp.Ip, nil
}

func getPublicIpv4_backup_bigdatacloud() (string, error) {
	type Msg_Resp_PublicIpv4 struct {
		IpString string `json:"ipString"`
	}

	resp := &Msg_Resp_PublicIpv4{}

	// try from public api
	err := api.Get_("https://api.bigdatacloud.net/data/client-ip", "", 10, resp)
	if err != nil {
		return "", err
	}

	netip := net.ParseIP(resp.IpString)
	if netip == nil {
		return "", errors.New("ip not find")
	}
	if netip.To4() == nil {
		return "", errors.New("ip not find")
	}

	return resp.IpString, nil
}

func getPublicIpv4_backup_ipconfigme() (string, error) {
	type Msg_Resp_PublicIpv4 struct {
		IpAddr string `json:"ip_addr"`
	}

	resp := &Msg_Resp_PublicIpv4{}

	// try from public api
	err := api.Get_("https://ifconfig.me/all.json", "", 10, resp)
	if err != nil {
		return "", err
	}

	netip := net.ParseIP(resp.IpAddr)
	if netip == nil {
		return "", errors.New("ip not find")
	}
	if netip.To4() == nil {
		return "", errors.New("ip not find")
	}

	return resp.IpAddr, nil
}
