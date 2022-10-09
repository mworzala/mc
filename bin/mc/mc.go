package main

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mworzala/mc-cli/internal/pkg/app"
	"github.com/mworzala/mc-cli/internal/pkg/model"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

/**
mc run latest

- Store profiles
- One data directory per profile
- Download libraries
- Launch game


mc run latest
- Check for updates
  - Read latest version manifest
  - Read current version manifest
  - If "latest" is newer, write new manifest
- Install latest version if necessary
  - blah blah
- Launch latest version

mc login
- Begin devicecode auth flow
- Store access/refresh tokens + user info
- If person is already signed in, replace the old values

mc account list
- List all account usernames and uuids

mc config account.default notmattw

mc java list
mc java discover /path/to/java/home
mc java default liberica-17


*/

var basePath = "/Users/matt/dev/projects/mmo/mc-cli/temp"

func main() {
	launcher := app.NewApp(basePath)

	cliApp := &cli.App{
		Name:     "mc",
		Version:  "0.0.1",
		Usage:    "install and launch Minecraft: Java Edition",
		HideHelp: true,
		Commands: []*cli.Command{
			{
				Name:      "run",
				Usage:     "run a Minecraft instance",
				ArgsUsage: "[version]",
				Subcommands: []*cli.Command{
					{
						Name: "latest",
						Action: func(ctx *cli.Context) error {
							return launcher.RunLatest()
						},
					},
				},
				// Default handler, interprets first argument as instance name
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() == 0 {
						return errors.New("no version specified")
					}
					fmt.Println("run", ctx.Args())
					return nil
				},
			},
			{
				Name:  "login",
				Usage: "log in to a Microsoft account",
				Action: func(ctx *cli.Context) error {
					return launcher.LoginMicrosoft()
				},
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	//resp, err := http.Get("https://launchermeta.mojang.com/mc/game/version_manifest.json")
	//if err != nil {
	//	panic(err)
	//}
	//
	//defer resp.Body.Close()
	//var manifest model.Manifest
	//err = json.NewDecoder(resp.Body).Decode(&manifest)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = writeManifest(&manifest)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = os.MkdirAll(path.Join(basePath, "versions", manifest.Latest.Release), 0755)
	//if err != nil {
	//	panic(err)
	//}
	//
	//version := manifest.Latest.Release
	//fmt.Println("Latest release:", version)
	//
	//latest := getVersion(&manifest, version)
	//if latest == nil {
	//	panic("Version not found")
	//}
	//
	//// Install
	//err = install.Vanilla(app.NewContext(basePath), latest)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// Read back the version manifest
	//f, err := os.OpenFile(path.Join(basePath, "versions", version, fmt.Sprintf("%s.json", version)), os.O_RDONLY, 0644)
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
	//var v model.Version
	//err = json.NewDecoder(f).Decode(&v)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// Microsoft auth
	////data, err := auth.BeginDeviceCodeAuth()
	////if err != nil {
	////	panic(err)
	////}
	////
	////fmt.Println(data.VerificationURL + " " + data.UserCode)
	////
	////msoToken, err := auth.PollDeviceCodeAuth(data)
	////if err != nil {
	////	panic(err)
	////}
	////
	//////msoToken := auth.DeviceCodeToken{
	//////	AccessToken:  "EwAQA+pvBAAUKods63Ys1fGlwiccIFJ+qE1hANsAAW5PTtM8BYX8T4KiUarxZExgXy5uKFT7HJt5dkPdNpr2av3BjeeXQEBlxt14AsfoyJy05EY33OqnwRuQtlcMkVaBtlvmw/BKG6a5DBlnqxl0zk5kcBr6HsxQTnfPOK59Q4hSAUoPQx7emJse/zPXALKXjSpkTrsTC8Rr/SyVtWL+1EYSHX/yn4+ZzAyGgyhoGo0fu5H86pN6blrSmEJVT3S5IOGULcVmM5+/bs2XJoVF+XjzqmJBtKA69jupr1gLJgwDL43SvjkmyXC4szcTjvcgCR6aun9TFdsBG+SZCCCJ33yMOm3bSp6mmRsgq2LVjqTqLz/5hFrUOPV9bTduHp4DZgAACMFtRX6fHM0H4AGeX1O8m4GdqdA9hVEVIPBr6i8xRxI3h114T6o2zL9FNJwtIIl/j1BKVkh650DqwbeMTyncSiyVYnAWZ5umDQoCf4y0CwwG47+2bEnBQPfaGg1KkeJLyh12U2ylPlxW5svR9HxOW2SboCp5zvyBVgftHtENMoaimr28ouQMV0K5t3Q/e79r8OmugmOXt5FUtu7qGgaobxUMLZSIhSmUK9fpDSecG8O0GtTcT4vK2xX9BaV5c0c+fe6wUZPPdNwsfmzYd0MBaixiSkRI6Giq9O1T/0HREub43Oqn8RWYsRk1KBeaiEABOHQTqkEOiS94yotcdTk+NXDUEQOckmdea8RozfzDVTz8udghdp9OTCIrzQqbXEMfRO9e49giIGHW+CzoSKsUFezueb3KyZQvZE7zrEg9Ga7SYBFEHIcGngHDHAwAfHWixKHW9m7/6+BVs4+B0c86QJwp2Kg/ElIW1D3TABgLZKom0Kq+7LGP4Uoe6IEOGj0lKocts1fOruiPnQAMURGD8g7X/9SPlncHmw7U7PdSZyQbjc+y30HCiCeBAxAep2v+oS1neSwnaV9YBE3hcTgK3Z/5eAzP+gYbUzXltMpWOZMdW2VEIeU/j8wvNUDq7hzTCB7i2OQJ6UwBo2cXAg==",
	//////	RefreshToken: "",
	//////	ExpiresIn:    0,
	//////	TokenType:    "",
	//////}
	////xblToken, err := auth.XboxLiveAuth(msoToken.AccessToken)
	////if err != nil {
	////	panic(err)
	////}
	////
	//xblToken := auth.XboxLiveToken{
	//	AccessToken: "eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhDQkMtSFMyNTYiLCJ6aXAiOiJERUYiLCJ4NXQiOiJxcEQtU2ZoOUg3NFFHOXlMN1hSZXJ5RmZvbk0iLCJjdHkiOiJKV1QifQ.tL5K2zW-Q8NArG6xGPu2bKiYuVfeWBZ9YH_bTDhKpp5GgOWB9mRliM1Kw0sE9kK8vx_cczDHJP6m-CQF8pSNtEJkw0tfqN3H1avBPhRyLu32hQCSgWv9kLv3-uMiUgZF3GrZ0CUoZnavB7Z63wzvCQYT4L-wa3E7vlKSNPtlCXToppo1zAMnhXtACI1XWeVBipmfnmNCe-Ryn5jgvC9IZL7dzJF-Z64SCQQTv-9PEOb9yp3UUjdu3f31Hw4ZNM7Sq4RWFRqaAEHsgIu2sOtKyQF0jf6u60MrSIEB7XcmdnDilWk9WL6bv_HDlboamL5LHyk75aXx5b86ruHgaxvenA.yWYOs39uOM5--nYn7MWOqA.0_o0HJzNoojaEW3jEoRXk77A0wUuoAI4u3Rf_WI94QEHlmR4KLgllZZpkxivAQ4g9cOqn7aU_ppWVThv7VZh330XIWWrw1x0c_Ah6b-lLL9pfTsajK44p33nw4wAyVxB6fUiVsGKJRl8CG5XbHEKizcz3hx6W0-yVdxTNyhU96Jb5XXtH6fJb6miP1uQn5yOEY1KkU4FkQG49W06lhTkQc2zG8J_MuTQDmgkgALczJMP88kps_yNYjfKNQUIOaSXhB9J2W7AMTBuYsDjsZMlYVGtqwFYPsdgFfSkC8nYFw-dPu1646rqWzbGW60mtdGxX3idIKw6ONGwWq4PNVxgB8xWS_i3SK1N5LjHmd2whzX-ocK5JFgZ1_0k1NuoEC_BkjQAx-8fkLivml2pgrhHwS_TDt6f7hE4m97TjQPRxkwJF-JFZ2Pr1nGC0Ie9gzdQ84M2nehjNzWKnqI0ED-mC9-0sW3QB25kmwCUhha87G6QTDJrxI2UbeOaNS970ggEPLfhnwamFMAynoIhtp4fbgW3Ahbs_-p9399ECVXZLjGU5zajMBppMaUyNsdK_9WMz9sT9a9vxKFOD7P-xD1X_uXSi9QqQyh7b6E1nT8l9Li10uiGzGenBQSGCyosccuLz_ECUQsHbsbWHBL6SKNWAlzg84uZC64hyinlu9unzgqaxAsT_xU6IvbdJywIMGIwHuUIJr0rDvRwkBGoyvYG8lgDHx2XY8xmjHSUBWthrH4NhjrE2UYUiJviTljXsbWox0z0o4JsdoZeOaZuVmHVUmZ3mTJEOp1cJvL_YC97k8kryISyd6Y-CYdXJyV8YesYNu4WqVkjiRJOfacjXUTX00TroPNWWoBFQAmqBj1m3D8L5lOOGDSCo49ssPrIvh5CF8iXbVu0Xx9e15WYEmwFg-YeEn9yPz9z-l5sP4sPXFVzJU1czsvdsa43xnCRVXYS392j5SM2_OEn5Cc5S9OUOC05v4w6aAfyLLZBd9lGkQc.01gCpktaSdwZQq6WIZKOdw",
	//	UserHash:    "5580347117603921787",
	//}
	////
	////xstsToken, err := auth.XSTSAuth(xblToken.AccessToken)
	////if err != nil {
	////	panic(err)
	////}
	//
	////xstsToken := auth.XSTSToken{
	////	AccessToken: "eyJlbmMiOiJBMTI4Q0JDLUhTMjU2IiwiYWxnIjoiUlNBLU9BRVAiLCJjdHkiOiJKV1QiLCJ6aXAiOiJERUYiLCJ4NXQiOiJzYVkzV1ZoQzdnMmsxRW9FU0Jncm9Ob2l3MVEifQ.I90KPU5oagwc7ClRGZcO8KL1sQVGAjapLLQU6C4sncIx6vkgB66bdK6bc8-hFyi1rNXW8sl5-mxapZTxPnsC1J-3bCNblJOvZ-8h5O1v1aiszQmeNYUmH2QsVnhlhaSuvggy_3vPjddHAYOwhGDlP3cKr8fqb6Do0C9hAl0fO24M8hYk-33bRlX31pmNkIKD59bYb-Z1BvRmsgovEbF-CCnigJ3VbA95d3i3AAVkW1NAG7UuIxP3-uhNA2Ag2bgiaiXLEe1HTBbGXu2i-5pJ7sXplRrLJMqEEsw1kJpU6_zcr8ya5E03J5lMi5liBGpdqRvlVC38vLBvF885gOi7Yw.qPQc6juxBNPIpnWX1BDGUw.fmpPzKtFkkvR-xYgGoWGyKZ_mja0tEniuBwqt8qIubF3MLwOnqYCasA5MaoJW8huCcar89Aq5SujhfZuJkDbMJc25a_KWIrjq1jo2QLhKNbI5V1EM7gv2iEzPyGMZvPQSo63vsC-MrgFsq2790jLLfTK8DvOzRkgqcChfH4awZH_vTfu-wMLb3zjum7X0mChRs637iwAX6Hpv_SiHhXuBaiYbt4MST5EIKPaRDjMK3inHrMvsc78YMlDs50_aMJje0rhwh4wyvzdUPUZXmOUfaDcJg14nthwSUZpITR2EOzEx_B7FcDm4tbZbwrAp9ZzDjGNJPsyFzahvONwce7hptneVjyEge3rdbn5m2jwSVe0gbDyFm7Q2Vw6sEhfsrRYjr5HC-hiAdwsgunYkJqwYns9HeeZoIuMoUh7-7PrqD_MDbaT3DLUE4-nBmEBSqqIQ_naYBhEMYZVVYl2gMAhjJMX2y3SReMKtIX99LdMabIwuxHan43b5i7IJZUEBIJ6hCozPPv5I_LMcMq_d3x5d-KTVEgTNnvyjVNVg1z1dUA3oLbCIvGMocHDmi4uPOb3LrYAhpiEM_6YI4jf47izbaDgPHfB1ygauiKvdIUPo2D5Phd9FZMi6HgwqfnpnwMoHrnV8ueaoOXZZGghRT4QDXXWR6wgA9sl_BSSPfRX9lfQHZ0Uv_3OZSAZRg97XoaHzrhZc2Gp4RbW9GNLZyo4KKfZn1732OcmNcn2QCYSzPvrCVKAt9aAywCxI0KGHrzOYeMvlXJ_dg-Q37_Qu53d0plMCIcSnN0ecAlIqiuVUh7crCkVjrMmJpE8VD6IXTt7ceOY4p66HnJ65uAoMJO7YRM101LFr9TGYFI9I_smoo_Ax0ZO58hxGt-64vUw-E6-cut3-twuheqJkVXsSObip_LIgOZL1d2PEWag_tgePn_i3a1tXHGbVBpn7_fnvPL_MRW25dihtjs0ulzVG3wR_RkEdf4LEUC7fWQKoRWTsG8LFHKKYEhrfMbfAIZ-NkffyhpKhlrdIMdPV_XNUCs3EpcaDth7201X5LzVTMcKXGcNUqKlprrocw1iAZCCnx9y3BiOPDG4ajNw4sOXbDwiEq2VklSaqInCtnKfgnIpfZGFq05xYa_Xe9cxxHATefmT9ZLllOXrZnCLUvcKKwJQAk4nA1N4FDjWsynV4cArYwhyVlFTPvgO715f1gzV2VeQwO9xwdltql-JX7fmxF_jIZ7FRhDWiKpHNPnTkDEWIGGe3-wvCWde_Cw1wUtmGxnJiqYJwAz_hQ-5523zvtcrSbUL05G_cM5X6GCghxdVC7fiRvcdQ64KgAaP-vPPK-t1-xCTZ1tTloDGMRJhkld-53M-rq43YkhjnI5dt41ND1UCe-pSLYAjGZSH_Y_sivVmquKslI4DIN3qbvQtgikkP-kF6bIR68hVUXSDxIzs9uHaGiT0VBz5rxML-vmvyha-hKqkq0TGKAEWMF0da-Mjl2l6dGFKAjNJGJMvwEbrvM9Y3DsboNxWbfMYNSGyFiOQp4Ibuj48vhzR13jVdc2sIKGFZF1R-VDmPndfK1gjEaBDBGYaP3Yir75IDZpwQggp.RpupmQAFZGsIorec4qukhA",
	////}
	////
	////mcToken, err := auth.MinecraftAuthMSO(xstsToken.AccessToken, xblToken.UserHash)
	////if err != nil {
	////	panic(err)
	////}
	//
	//mcToken := auth.MinecraftToken{
	//	AccessToken: "eyJhbGciOiJIUzI1NiJ9.eyJ4dWlkIjoiMjUzNTQ3MDkyODI4NDk2NiIsImFnZyI6IkFkdWx0Iiwic3ViIjoiYmNlZTc5ODItY2ViNC00OGVmLWE2MTctODA1ZjI4MTZhNDg4IiwibmJmIjoxNjY1MjE4MzM5LCJhdXRoIjoiWEJPWCIsInJvbGVzIjpbXSwiaXNzIjoiYXV0aGVudGljYXRpb24iLCJleHAiOjE2NjUzMDQ3MzksImlhdCI6MTY2NTIxODMzOSwicGxhdGZvcm0iOiJVTktOT1dOIiwieXVpZCI6IjhmNWFlNjA3ZmU4MGExMWM4NjUyNzg2OTZlN2EyYzA0In0.P3H6bAV5ZZM5KsqHzB3ccq909kXymp6tPv_cvidzco0",
	//	ExpiresIn:   86400,
	//}
	////
	////fmt.Println(mcToken)
	//
	//// Launch
	//var args []string
	//
	//classpath := strings.Builder{}
	//librariesPath := path.Join(basePath, "libraries")
	//for _, lib := range v.Libraries {
	//	libPath := path.Join(librariesPath, lib.Downloads.Artifact.Path)
	//	classpath.WriteString(libPath)
	//	classpath.WriteString(":")
	//}
	//classpath.WriteString(path.Join(basePath, "versions", v.Id, v.Id+".jar"))
	////fmt.Println(classpath.String())
	//
	//vars := map[string]string{
	//	// jvm
	//	"natives_directory": ".",
	//	"launcher_name":     "mc-cli",
	//	"launcher_version":  "0.0.1",
	//	"classpath":         classpath.String(),
	//	// game
	//	"version_name":      latest.Id,
	//	"game_directory":    path.Join(basePath, "tmp-instance"),
	//	"assets_root":       path.Join(basePath, "assets"),
	//	"assets_index_name": v.Assets,
	//	"auth_player_name":  "notmattw",
	//	"auth_uuid":         "aceb326fda1545bcbf2f11940c21780c",
	//	"auth_access_token": mcToken.AccessToken,
	//	"clientid":          "MTMwQUU2ODYwQUE1NDUwNkIyNUZCMzZBNjFCNjc3M0Q=",
	//	"auth_xuid":         xblToken.UserHash,
	//	"user_type":         "msa",
	//	"version_type":      latest.Type,
	//	"resolution_width":  "1920",
	//	"resolution_height": "1080",
	//}
	//_ = vars
	//
	/////Users/matt/Library/Application Support/minecraft/runtime/java-runtime-gamma/mac-os-arm64/java-runtime-gamma/jre.bundle/Contents/Home/bin/java
	////-XstartOnFirstThread
	////-Djava.library.path=.
	////-Dminecraft.launcher.brand=minecraft-launcher
	////-Dminecraft.launcher.version=2.3.443
	////-cp {CLASSPATH}
	////-Xmx16G
	////-XX:+UnlockExperimentalVMOptions
	////-XX:+UseG1GC
	////-XX:G1NewSizePercent=20
	////-XX:G1ReservePercent=20
	////-XX:MaxGCPauseMillis=50
	////-XX:G1HeapRegionSize=32M
	////-Dlog4j.configurationFile=/Users/matt/Library/Application Support/minecraft/assets/log_configs/client-1.12.xml
	////net.fabricmc.loader.impl.launch.knot.KnotClient
	////--username notmattw
	////--version 1.19.2
	////--gameDir /Users/matt/Library/Application Support/minecraft/fabric-1.19
	////--assetsDir /Users/matt/Library/Application Support/minecraft/assets
	////--assetIndex 1.19
	////--uuid aceb326fda1545bcbf2f11940c21780c
	////--accessToken eyJhbGciOiJIUzI1NiJ9.eyJ4dWlkIjoiMjUzNTQ3MDkyODI4NDk2NiIsImFnZyI6IkFkdWx0Iiwic3ViIjoiYmNlZTc5ODItY2ViNC00OGVmLWE2MTctODA1ZjI4MTZhNDg4IiwibmJmIjoxNjY1MjA3Mzg5LCJhdXRoIjoiWEJPWCIsInJvbGVzIjpbXSwiaXNzIjoiYXV0aGVudGljYXRpb24iLCJleHAiOjE2NjUyOTM3ODksImlhdCI6MTY2NTIwNzM4OSwicGxhdGZvcm0iOiJQQ19MQVVOQ0hFUiIsInl1aWQiOiI4ZjVhZTYwN2ZlODBhMTFjODY1Mjc4Njk2ZTdhMmMwNCJ9.Xe1sdcZ5w2ZAdNilwPv7KDT9vq1iG3Ytb_RB8l2uWc0
	////--clientId MTMwQUU2ODYwQUE1NDUwNkIyNUZCMzZBNjFCNjc3M0Q=
	////--xuid 2535470928284966
	////--userType msa
	////--versionType release
	//
	//replaceVars := func(s string) string {
	//	for k, v := range vars {
	//		s = strings.ReplaceAll(s, fmt.Sprintf("${%s}", k), v)
	//	}
	//	return s
	//}
	//
	//args = append(args, "-XstartOnFirstThread")
	//
	//for _, arg := range v.Arguments.JVM {
	//	if s, ok := arg.(string); ok {
	//		args = append(args, replaceVars(s))
	//	} else if m, ok := arg.(map[string]interface{}); ok {
	//		_ = m
	//		//value := m["value"]
	//		//if s, ok := value.(string); ok {
	//		//	args = append(args, replaceVars(s))
	//		//} else if a, ok := value.([]interface{}); ok {
	//		//	for _, v := range a {
	//		//		if s, ok := v.(string); ok {
	//		//			args = append(args, replaceVars(s))
	//		//		}
	//		//	}
	//		//} else {
	//		//	panic(fmt.Sprintf("unknown type: %T", value))
	//		//}
	//	} else {
	//		panic("unknown arg type")
	//	}
	//}
	//
	//args = append(args, v.MainClass)
	//
	//for _, arg := range v.Arguments.Game {
	//	if s, ok := arg.(string); ok {
	//		args = append(args, replaceVars(s))
	//	} else if m, ok := arg.(map[string]interface{}); ok {
	//		_ = m
	//		//value := m["value"]
	//		//if s, ok := value.(string); ok {
	//		//	args = append(args, replaceVars(s))
	//		//} else if a, ok := value.([]interface{}); ok {
	//		//	for _, v := range a {
	//		//		if s, ok := v.(string); ok {
	//		//			args = append(args, replaceVars(s))
	//		//		}
	//		//	}
	//		//} else {
	//		//	panic(fmt.Sprintf("unknown type: %T", value))
	//		//}
	//	} else {
	//		panic("unknown arg type")
	//	}
	//}
	////
	////for _, arg := range args {
	////	fmt.Println(arg)
	////}
	////
	//javaBin := "/Users/matt/Library/Java/JavaVirtualMachines/liberica-17.0.1/bin/java"
	//cmd := exec.Command(javaBin, args...)
	//cmd.Dir = path.Join(basePath, "tmp-instance")
	//
	//tail := len(os.Args) > 1 && os.Args[1] == "-t"
	//if tail {
	//	cmd.Stdout = os.Stdout
	//} else {
	//	cmd.Stdout = io.Discard
	//}
	//
	//if err != nil {
	//	panic(err)
	//}
	//err = cmd.Start()
	//if err != nil {
	//	panic(err)
	//}
	//
	//if tail {
	//	err = cmd.Wait()
	//	if err != nil {
	//		panic(err)
	//	}
	//}

}

func downloadFileIfNotExists(id, file string, dl model.Download) error {
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		resp, err := http.Get(dl.Url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		//todo breaks for asset index for some reason
		//if resp.ContentLength != dl.Size {
		//	return fmt.Errorf("download size mismatch for %s: %d != %d", id, resp.ContentLength, dl.Size)
		//}

		f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		bar := progressbar.DefaultBytes(resp.ContentLength, id)

		h := sha1.New()
		_, err = io.Copy(io.MultiWriter(f, h, bar), resp.Body)
		if err != nil {
			return err
		}

		hash := fmt.Sprintf("%x", h.Sum(nil))
		if hash != dl.Sha1 {
			return fmt.Errorf("download hash mismatch for %s: %s != %s", id, hash, dl.Sha1)
		}
	}

	return nil
}

func getVersion(manifest *model.Manifest, version string) *model.ManifestVersion {
	for _, v := range manifest.Versions {
		if v.Id == version {
			return v
		}
	}
	return nil
}

func writeManifest(manifest *model.Manifest) error {
	file, err := os.OpenFile(path.Join(basePath, "manifest.json"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(manifest)
	if err != nil {
		return err
	}

	return nil
}

func writeJson(file string, data interface{}) error {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(data)
	if err != nil {
		return err
	}

	return nil
}
