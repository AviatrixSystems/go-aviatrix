package goaviatrix

// Gateway simple struct to hold gateway details

import (
	"fmt"
	"encoding/json"
	"errors"
	//"io/ioutil"

	"github.com/davecgh/go-spew/spew"
)

type GatewayResp struct {
	CID                     string `form:"CID"`
	IntraVMRoute            string `json:"intra_vm_route"`
	VpnCidr                 string `json:"vpn_cidr"`
	HaEnabled               string `json:"ha_enabled"`
	DirectInternet          string `json:"direct_internet"`
	VpcSplunkIPPort         string `json:"vpc_splunk_ip_port"`
	EnableNat               string `json:"enable_nat"`
	InstState               string `json:"inst_state"`
	BkupGatewayZone         string `json:"bkup_gateway_zone"`
	PrivateIP               string `json:"private_ip"`
	AccountName             string `form:"account_name" json:"account_name"`
	VpcState                string `json:"vpc_state"`
	IsHagw                  string `json:"is_hagw"`
	DockerNtwkCidr          string `json:"docker_ntwk_cidr"`
	CloudnGatewayInstID     string `json:"cloudn_gateway_inst_id"`
	CloudType               int    `form:"cloud_type" json:"cloud_type"`
	VpcRegion               string `json:"vpc_region"`
	DockerNtwkName          string `json:"docker_ntwk_name"`
	SamlEnabled             string `json:"saml_enabled"`
	VpcCidr                 string `json:"vpc_cidr"`
	LicenseID               string `json:"license_id"`
	MaxConnections          string `json:"max_connections"`
	PbrEnabled              string `json:"pbr_enabled"`
	GatewayZone             string `json:"gateway_zone"`
	SandboxIP               string `json:"sandbox_ip"`
	ElbDNSName              string `json:"elb_dns_name"`
	TunnelType              string `json:"tunnel_type"`
	GwName                  string `form:"gw_name" json:"vpc_name"`
	ClientCertAuth          string `json:"client_cert_auth"`
	DockerConsulIP          string `json:"docker_consul_ip"`
	GwSubnetID              string `json:"gw_subnet_id"`
	TunnelName              string `json:"tunnel_name"`
	SplitTunnel             string `json:"split_tunnel"`
	ElbState                string `json:"elb_state"`
	AuthMethod              string `json:"auth_method"`
	BkupPrivateIP           string `json:"bkup_private_ip"`
	VpcType                 string `json:"vpc_type"`
	CloudnBkupGatewayInstID string `json:"cloudn_bkup_gateway_inst_id"`
	ClientCertSharing       string `json:"client_cert_sharing"`
	GwSecurityGroupID       string `json:"gw_security_group_id"`
	VpcSize                 string `json:"vpc_size"`
	Expiration              string `json:"expiration"`
	PublicIP                string `json:"public_ip"`
	VpcID                   string `form:"vpc_id" json:"vpc_id"`
	VendorName              string `json:"vendor_name"`
	VpnStatus               string `json:"vpn_status"`
}

type Gateway struct {
	CID         string `form:"CID"`
	CloudType   int    `form:"cloud_type" json:"cloud_type"`
	AccountName string `form:"account_name" json:"account_name"`
	GwName      string `form:"gw_name" json:"vpc_name"`
	VpcID       string `form:"vpc_id" json:"vpc_id"`
	VpcReg      string `form:"vpc_reg" json:"vpc_region"`
	VpcSize     string `form:"vpc_size" json:"vpc_size"`
	VpcNet      string `form:"vpc_net" form:"vpc_cidr"`
}

type GatewayListResp struct {
	Return  bool   `json:"return"`
	Results []Gateway `json:"results"`
	Reason  string `json:"reason"`
}

func (c *Client) CreateGateway(gateway *Gateway) (error) {
	gateway.CID=c.CID
	resp,err := c.Post(c.baseUrl, gateway)
		if err != nil {
		return err
	}
	var data ApiResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) GetGateway(gateway *Gateway) (*Gateway, error) {
	//gateway.CID=c.CID
	path := c.baseUrl + fmt.Sprintf("?CID=%s&action=list_vpcs_summary&account_name=%s", c.CID, gateway.AccountName)
	fmt.Println("PaTh: ", path)
	resp,err := c.Get(path, nil)

	if err != nil {
		return nil, err
	}
	var data GatewayListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	spew.Dump(data)
	if(!data.Return){
		return nil, errors.New(data.Reason)
	}
	gwlist:= data.Results
	for i := range gwlist {
    	if gwlist[i].GwName == gateway.GwName {
        	return &gwlist[i], nil
    	}
	}


	return nil, errors.New(fmt.Sprintf("Gateway %s not found", gateway.GwName))	
}

func (c *Client) UpdateGateway(gateway *Gateway) (error) {
	return nil
}

func (c *Client) DeleteGateway(gateway *Gateway) (error) {
	//gateway.CID=c.CID
	path := c.baseUrl + fmt.Sprintf("?action=delete_container&CID=%s&cloud_type=%d&gw_name=%s", c.CID, gateway.CloudType, gateway.GwName)
	fmt.Println("PaTh: ", path)
	resp,err := c.Delete(path, nil)

	if err != nil {
		return err
	}
	var data ApiResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	fmt.Println(data.Results)
	return nil
}
