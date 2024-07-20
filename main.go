package main

import (
	"AifadianCrawler/afdian"
	"AifadianCrawler/afdian/album"
	"AifadianCrawler/afdian/motion"
	"AifadianCrawler/utils"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

var (
	afdianHost    string
	authorName    string
	albumListPath string
)

func main() {
	successColor := color.New(color.Bold, color.FgGreen).FprintlnFunc()
	failColor := color.New(color.Bold, color.FgRed).FprintlnFunc()
	//记录开始时间
	startTime := time.Now()
	app := &cli.App{
		Name:            "AifadianCrawler",
		Usage:           "爱发电下载器，支持按作者或按作品集爬取数据\nGithub Link: https://github.com/PhiFever/AifadianCrawler",
		UsageText:       "eg:\n	./AifadianCrawler -au Jay motions \neg:\n\t./AifadianCrawler.exe -au Jay albums \neg:\n	./AifadianCrawler.exe -l album_list.txt",
		Version:         "0.2.1",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Destination: &afdianHost, Value: "afdian.net", Usage: "主站域名，默认为afdian.net，被封可自行更改"},
			&cli.StringFlag{Name: "author", Aliases: []string{"au"}, Destination: &authorName, Value: "", Usage: "待下载的作者id"},
			&cli.StringFlag{Name: "list", Aliases: []string{"l"}, Destination: &albumListPath, Value: "", Usage: "待下载的作品集id列表文件，每行一个id。(不能与参数-au同时使用)"},
		},
		Before: func(c *cli.Context) error {
			// 在这里可以根据需要做全局参数的预处理
			if authorName != "" && albumListPath != "" {
				return fmt.Errorf("不能同时使用 --author 和 --list 参数")
			}
			afdian.Host = fmt.Sprintf("https://%s", afdianHost)
			// 其他全局预处理任务...
			return nil
		},
		After: func(c *cli.Context) error {
			// 在这里可以根据需要做全局参数的后处理
			// 其他全局后处理任务...
			//记录结束时间
			endTime := time.Now()
			//计算执行时间，单位为秒
			successColor(os.Stdout, "处理完毕，共耗时:", utils.GetExecutionTime(startTime, endTime))
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "motions",
				Usage: "下载指定作者的所有动态",
				Action: func(c *cli.Context) error {
					return motion.GetMotions(authorName)
				},
			},
			{
				Name:  "albums",
				Usage: "下载指定作者的所有作品集",
				Action: func(c *cli.Context) error {
					return album.GetAlbums(authorName)
				},
			},
			{
				Name:  "update",
				Usage: "更新所有已经下载的作者的动态和作品集",
				Action: func(c *cli.Context) error {
					authors, err := utils.CheckAndListAuthors()
					if err != nil {
						return err
					}
					for _, author := range authors {
						log.Println("find exist author: ", author)
						if err := motion.GetMotions(author); err != nil {
							return err
						}
						if err := album.GetAlbums(author); err != nil {
							return err
						}
					}
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			//TODO:处理全局参数albumListPath
			if albumListPath != "" {
				fmt.Println("albumListPath:", albumListPath)
			} else {
				return fmt.Errorf("albumListPath=None")
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		failColor(os.Stderr, err)
	}
}
