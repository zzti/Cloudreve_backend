package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// 1. 清理旧文件
	os.RemoveAll("assets/build")
	os.Remove("assets.zip")

	// 2. 执行 Yarn 构建
	// 我们把 "assets" 作为第一个参数传入，指定为命令执行的工作目录
	runCmdInDir("assets", "yarn", "install", "--network-timeout", "1000000")
	runCmdInDir("assets", "yarn", "run", "build")

	// 3. 压缩 build 目录到 assets.zip
	if err := zipWithStructure("assets/build", "assets.zip"); err != nil {
		fmt.Printf("压缩失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("资产构建完成: assets.zip")
}

// 增加 dir 参数，并设置了 cmd.Dir
func runCmdInDir(dir string, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir // 强制子进程在这个目录下运行，解决 Husky 找不到路径的问题
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	// 如果是 Windows 还需要处理某些环境下找不到 yarn 的问题（可选）
	// cmd.Env = os.Environ() 

	if err := cmd.Run(); err != nil {
		fmt.Printf("在目录 %s 执行 %s 失败: %v\n", dir, name, err)
		os.Exit(1)
	}
}

// zipWithStructure 按照 assets/build/文件名 的结构进行打包
func zipWithStructure(srcDir, destZip string) error {
	out, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer out.Close()

	w := zip.NewWriter(out)
	defer w.Close()

	// 遍历 assets/build 目录
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 忽略目录条目，只处理文件
		if info.IsDir() {
			return nil
		}

		// 核心逻辑：
		// path 现在是 "assets\build\version.json" (Windows)
		// 我们需要把它在 ZIP 里的路径保持为 "assets/build/version.json"
		
		// 1. 将 Windows 的 \ 转换为 /
		standardPath := filepath.ToSlash(path)

		// 2. 创建 ZIP 中的文件条目
		f, err := w.Create(standardPath)
		if err != nil {
			return err
		}

		// 3. 写入内容
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		_, err = io.Copy(f, in)
		return err
	})
}

func zipDirNew(src, dest string) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	w := zip.NewWriter(out)
	defer w.Close()

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// 如果是目录，不需要单独创建条目（zip 写入文件时会自动创建父目录）
		if info.IsDir() {
			return nil
		}

		// 1. 获取相对路径（例如从 assets/build/locales... 变成 locales\en-US\...）
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// 2. 重要修复：ZIP 规范规定目录分隔符必须是正斜杠 "/"，将 Windows 路径分隔符" \" 统一转换为 ZIP 标准的"/"
		// 这样无论在什么系统上压缩，ZIP 内部的路径都是 locales/en-US/...
		zipPath := filepath.ToSlash(relPath)

		// 3. 写入 ZIP
		f, err := w.Create(zipPath)
		if err != nil {
			return err
		}
		
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		
		_, err = io.Copy(f, in)
		return err
	})
}


func zipDirOld(src, dest string) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	w := zip.NewWriter(out)
	defer w.Close()

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		relPath, _ := filepath.Rel(src, path)
		f, err := w.Create(relPath)
		if err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		_, err = io.Copy(f, in)
		return err
	})
}
