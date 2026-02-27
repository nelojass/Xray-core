package main

import (
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/core"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/features/stats"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"testing"
)

func TestRunConfig(t *testing.T) {
	testCases := []struct {
		Input  string
		Output string
	}{
		{
			Input: `{
				"log": {
					"loglevel": "info"
				}
		}`,
			Output: "",
		},
		{
			Input: `{
				"log": {
					// abcd
					"loglevel": "info",
				}
		}`,
			Output: "line 5 char 5",
		},
		{
			Input: `{
				"port": 1,
				"inbounds": [{
					"protocol": "test"
				}]
		}`,
			Output: "parse json config",
		},
		{
			Input: `{
				"inbounds": [{
					"port": 1,
					"listen": 0,
					"protocol": "test"
				}]
		}`,
			Output: "line 1 char 1",
		},
	}
	for _, testCase := range testCases {
		err := RunJsonConfig(testCase.Input)
		if err != nil {
			errString := err.Error()
			if !strings.Contains(errString, testCase.Output) {
				t.Error("unexpected output from json: ", testCase.Input, ". expected ", testCase.Output, ", but actually ", errString)
			}
		}
	}
}

const Private_Key = "uKyoqS6z6qT5RwSPLoKVZlnO3csElxml1ENgYfRNuHg"
const Public_key = "Jp0PpbLrYWk34_m9oyl7cDkeml49KuDWgvbqqc6w6BU"

func TestRunInstanceConfig(t *testing.T) {
	testCases := []struct {
		Input  string
		Output string
	}{
		{
			Input: `{
				"log": {
						"loglevel": "debug"
					},
				"inbounds":[
					{
						"port":10881,
						"listen":"127.0.0.1",
						"protocol":"http"
					}
				],
				"outbounds":[
					{
						"protocol": "vless",
						"mux":{
							"enabled":false, //reality不能使用多路复用，原理来说是正常的
							"concurrent":8
						},
						"settings": {
							"vnext": [
								{
									"address": "127.0.0.1", // 服务端的域名或 IP
									"port": 1443,
									"users": [
										{
											"id": "a3edfebb-670e-42b4-907f-f58dad9e0a15", // 与服务端一致
											"flow": "xtls-rprx-vision", // 与服务端一致
											"encryption": "none"
										}
									]
								}
							]
						},
						"streamSettings": {
							"network": "raw",
							"security": "reality",
							"realitySettings": {
								"show": true, // 选填，若为 true，输出调试信息
								"fingerprint": "chrome", // 选填，使用 uTLS 库模拟客户端 TLS 指纹，默认 chrome
								"serverName": "google.com", // 服务端 serverNames 之一
								"password": "Jp0PpbLrYWk34_m9oyl7cDkeml49KuDWgvbqqc6w6BU", // 服务端私钥生成的公钥，对客户端来说就是密码
								"shortId": "0123456789abcdef", // 服务端 shortIds 之一
								"mldsa65Verify": "", // 选填，服务端 mldsa65Seed 生成的公钥，对证书进行抗量子的额外验证
								"spiderX": "" // 爬虫初始路径与参数，建议每个客户端不同
							}
						}
					}
				]
			}`,
			Output: "",
		},
		{
			Input: `{
				"log": {
						"loglevel": "debug"
				},
				"policy":{
					"levels": {
						"0": {
							"handshake": 4,
							"connIdle": 300,
							"uplinkOnly": 2,
							"downlinkOnly": 5,
							"statsUserUplink": true,
							"statsUserDownlink": true,
							"bufferSize": 10240
						}
					},
					"system": {
						"statsInboundUplink": true,
						"statsInboundDownlink": false
					}
				},
				"stats":{},
				"inbounds":[
					{
						"listen": "127.0.0.1",
						"port": 1443,
						"protocol": "vless",
						"settings": {
							"clients": [
								{
									"level": 0,
									"email":"love@v2fly.org",
									"id": "a3edfebb-670e-42b4-907f-f58dad9e0a15", // 必填，执行 ./xray uuid 生成，或 1-30 字节的字符串
									"flow": "xtls-rprx-vision" // 选填，若有，客户端必须启用 XTLS
								}
							],
							"decryption": "none"
						},
						"streamSettings": {
							"network": "raw",
							"security": "reality",
							"realitySettings": {
								"show": true, // 选填，若为 true，输出调试信息
								"target": "google.com:443", // 必填，格式同 VLESS fallbacks 的 dest
								"xver": 0, // 选填，格式同 VLESS fallbacks 的 xver
								"serverNames": [ // 必填，客户端可用的 serverName 列表，暂不支持 * 通配符
									"google.com",
									"www.google.com"
								],
								"privateKey": "uKyoqS6z6qT5RwSPLoKVZlnO3csElxml1ENgYfRNuHg", // 必填，执行 ./xray x25519 生成
								"minClientVer": "", // 选填，客户端 Xray 最低版本，格式为 x.y.z
								"maxClientVer": "", // 选填，客户端 Xray 最高版本，格式为 x.y.z
								"maxTimeDiff": 0, // 选填，允许的最大时间差，单位为毫秒
								"shortIds": [ // 必填，客户端可用的 shortId 列表，可用于区分不同的客户端
									"", // 若有此项，客户端 shortId 可为空
									"0123456789abcdef" // 0 到 f，长度为 2 的倍数，长度上限为 16
								],
								"mldsa65Seed": "", // 选填，执行 ./xray mldsa65 生成，对证书进行抗量子的额外签名
								// 下列两个 limit 为选填，可对未通过验证的回落连接限速，bytesPerSec 默认为 0 即不启用
								// 回落限速是一种特征，不建议启用，如果您是面板/一键脚本开发者，务必让这些参数随机化
								"limitFallbackUpload": {
									"afterBytes": 0, // 传输指定字节后开始限速
									"bytesPerSec": 0, // 基准速率（字节/秒）
									"burstBytesPerSec": 0 // 突发速率（字节/秒），大于 bytesPerSec 时生效
								},
								"limitFallbackDownload": {
									"afterBytes": 0, // 传输指定字节后开始限速
									"bytesPerSec": 0, // 基准速率（字节/秒）
									"burstBytesPerSec": 0 // 突发速率（字节/秒），大于 bytesPerSec 时生效
								}
							}
						}
					}
				],
				"outbounds":[
					{
						"protocol":"freedom"
					}
				]
			}`,
			Output: "",
		},
	}

	instance, err := core.StartInstance("json", []byte(testCases[1].Input))
	if err != nil {
		t.Error("start instance error:", err.Error())
		return
	}

	stat := instance.GetFeature(stats.ManagerType())
	manager := stat.(stats.Manager)

	value, err := getStats("user>>>love@v2fly.org>>>traffic>>>uplink", true, manager)
	if err != nil {
		t.Fatal("register counter error", err)
	}

	t.Logf("register counter: %d", value)

	return

	for _, testCase := range testCases {
		err := RunJsonConfig(testCase.Input)
		if err != nil {
			errString := err.Error()
			if !strings.Contains(errString, testCase.Output) {
				t.Error("unexpected output from json: ", testCase.Input, ". expected ", testCase.Output, ", but actually ", errString)
			}
		}
	}

	select {}
}

func getStats(name string, isReset bool, s stats.Manager) (int64, error) {
	c := s.GetCounter(name)
	if c == nil {
		return 0, status.Error(codes.NotFound, name+" not found.")
	}
	var value int64
	if isReset {
		value = c.Set(0)
	} else {
		value = c.Value()
	}
	return value, nil
}
