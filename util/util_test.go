/*
 @Desc

 @Date 2020-10-09 20:12
 @Author inori
*/
package util

import (
	"fmt"
	"path"
	"testing"
)

func TestGetClientIp(t *testing.T) {
	fmt.Println(path.Base("https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1602307015829&di=e0debe701e4fa59f5ae8023d71691948&imgtype=0&src=http%3A%2F%2Ft9.baidu.com%2Fit%2Fu%3D1307125826%2C3433407105%26fm%3D79%26app%3D86%26f%3DJPEG%3Fw%3D5760%26h%3D3240"))
}
