package data

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Cluster struct {
	decPwd   []byte
	decUser  []byte
	Host     string
	Name     string
	encPwd   []byte
	Protocol string
	encUser  []byte
}

func NewWorkingCluster(key string) *Cluster {
	name := strings.Split(key, ".")[0]
	if name == "" {
		fmt.Println("invalid working cluster name, use config to set the working cluster")
		os.Exit(1)
	}
	tmp := name + "."
	return &Cluster{
		Host:     viper.GetString(tmp + HOST),
		encPwd:   toByteArray(viper.Get(tmp + PASSWORD).([]interface{})),
		Name: name,
		Protocol: viper.GetString(tmp + PROTOCOL),
		encUser:     toByteArray(viper.Get(tmp + USER).([]interface{})),
	}
}

func NewCluster(name string) *Cluster {
	return &Cluster{
		Name: name,
	}
}

func toByteArray(ia []interface{}) []byte {
	b := make([]byte, len(ia))
	for i, v := range ia {
		b[i] = byte(v.(int))
	}
	return b
}

func (c *Cluster) String() string {
	return c.Name + "\n* " + HOST + ": " + c.Host + "\n* " + PROTOCOL + ": " + c.Protocol
}

func (c *Cluster) Password() []byte {
	if c.decPwd == nil {
		p, err := Decrypt(c.encPwd)
		cobra.CheckErr(err)
		c.decPwd = p
	}
	return c.decPwd
}

func (c *Cluster) SetPassword(b []byte) {
	c.decPwd = b
	key := CreateKey()
	e, err := Encrypt(key, []byte(b))
	cobra.CheckErr(err)
	c.encPwd = e
}

func (c *Cluster) UserName() []byte {
	if c.decUser == nil {
		u, err := Decrypt(c.encUser)
		cobra.CheckErr(err)
		c.decUser = u
	}
	return c.decUser
}

func (c *Cluster) SetUserName(b []byte) {
	c.decUser = b
	key := CreateKey()
	e, err := Encrypt(key, []byte(b))
	cobra.CheckErr(err)
	c.encUser = e
}

func (c *Cluster) WriteToViper() {
	viper.Set(c.Name, map[string]interface{}{
		HOST:     c.Host,
		PASSWORD: c.encPwd,
		PROTOCOL: c.Protocol,
		USER:     c.encUser,
	})
}
