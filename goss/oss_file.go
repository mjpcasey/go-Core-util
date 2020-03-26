package goss

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gcore/glog"
)

type OssFile struct {
	Client *S3Client
	Config OssFileConfig
}

type OssFileConfig struct {
	// 强制优先读取
	ForceOriginFile bool

	// 默认就进行压缩
	IsNotCompress bool

	// 文件路径,
	FilePath string

	Bucket string
	Key    string
}

func NewOssFile(client *S3Client, conf OssFileConfig) *OssFile {
	if !conf.IsNotCompress {
		conf.Key = conf.Key + ".zip"
	}

	return &OssFile{
		Client: client,
		Config: conf,
	}
}

func (this *OssFile) ReadLine(hookfn func(string)) (err error) {
	err = this.CheckConfig()
	if err != nil {
		return err
	}

	// 如果强制需要读取远端文件，则先获取文件更新本地文件
	// 或者本地文件不存在，需要从远端读取文件
	if this.Config.ForceOriginFile || !this.FileExistFromLocal() {
		err = this.GetObjectAndSaveFile()
		if err != nil {
			return err
		}
	}

	fin, err := os.Open(this.Config.FilePath)
	if err != nil {
		glog.Errorf("无法开启文件[%s]，错误: %v", this.Config.FilePath, err)
		return err
	}

	defer func() {
		err = fin.Close()
	}()

	scanner := bufio.NewScanner(fin)
	for scanner.Scan() {
		hookfn(scanner.Text())
	}

	return
}

func (this *OssFile) DeleteFile() (err error) {
	// 删除本地
	if this.FileExistFromLocal() {
		err = os.Remove(this.Config.FilePath)
		if err != nil {
			return err
		}
	}

	// 删除远程
	if this.FileExistOss() {
		err = this.Client.DeleteObject(this.Config.Bucket, this.Config.Key)
		return err
	}

	return nil
}

func (this *OssFile) GetObjectMetaData() (metaData map[string]string, err error) {
	metaData, err = this.Client.GetObjectMetadata(this.Config.Bucket, this.Config.Key)
	return metaData, err
}

func (this *OssFile) getObject() (data []byte, err error) {
	output, err := this.Client.GetObject(this.Config.Bucket, this.Config.Key)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}

	if compressType, ok := output.Metadata["Compress-Type"]; !(ok && *compressType == "gzip") {
		data = body
		return data, nil
	} else {
		dataBuffer := new(bytes.Buffer)

		// 压缩文件的话，先解压
		zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
		if err != nil {
			return nil, err
		}

		for _, f := range zr.File {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer func() {
				err = rc.Close()
			}()

			_, err = io.Copy(dataBuffer, rc)
			if err != nil {
				return nil, err
			}
		}
		return dataBuffer.Bytes(), nil
	}
}

func (this *OssFile) GetObjectData() (data []byte, err error) {
	if this.Config.ForceOriginFile {
		return this.getObject()
	} else {
		// 优先从本地文件获取，没有则从OSS获取
		data, err = ioutil.ReadFile(this.Config.FilePath)
		if err == nil {
			return data, nil
		} else {
			return this.getObject()
		}
	}
}

func (this *OssFile) GetObjectAndSaveFile() (err error) {
	err = this.CheckConfig()
	if err != nil {
		return err
	}

	// 尝试打开文件
	fout, err := os.Create(this.Config.FilePath)
	if err != nil {
		log.Panic("fail to create file", err)
	}

	defer fout.Close()

	data, err := this.getObject()
	if err != nil {
		return err
	}
	_, err = fout.Write(data)

	return err
}

func (this *OssFile) PutObjectFromLocal() (err error) {
	err = this.CheckConfig()
	if err != nil {
		return err
	}

	metadata := make(map[string]string)
	var data []byte
	if this.Config.IsNotCompress {
		data, err = ioutil.ReadFile(this.Config.FilePath)
		if err != nil {
			return err
		}

	} else {
		zfile := this.Config.FilePath + ".gzip"
		err = zipfile(zfile, this.Config.FilePath)
		if err != nil {
			return err
		}

		data, err = ioutil.ReadFile(zfile)
		if err != nil {
			return err
		}

		defer func() {
			err = os.Remove(zfile)
		}()
		metadata["Compress-Type"] = "gzip"
	}
	return this.Client.PutObject(
		this.Config.Bucket,
		this.Config.Key,
		bytes.NewReader(data),
		metadata)
}

func (this *OssFile) CreateLocalFile() (*os.File, error) {
	err := this.CheckConfig()
	if err != nil {
		return nil, err
	}

	dirPath, _ := filepath.Abs(this.Config.FilePath)
	dirPath = filepath.Dir(dirPath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		err = os.MkdirAll(dirPath, 0755)

		if err != nil {
			glog.Errorf("文件夹=%s, 创建本地文件夹失败: %s", dirPath, err.Error())
			return nil, err
		}
	}

	// 尝试打开文件
	fout, err := os.Create(this.Config.FilePath)
	if err != nil {
		glog.Errorf("文件=%s 创建本地文件失败: %s", this.Config.FilePath, err.Error())
		return nil, err
	}
	return fout, err
}

func (this *OssFile) CheckConfig() error {
	if this.Config.FilePath == "" ||
		this.Config.Bucket == "" ||
		this.Config.Key == "" {
		return errors.New("oss file 配置错误")
	}

	return nil
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func (this *OssFile) FileExistFromLocal() bool {
	_, err := os.Stat(this.Config.FilePath)
	return err == nil || os.IsExist(err)
}

func (this *OssFile) FileExistOss() bool {
	_, err := this.Client.GetObjectMetadata(this.Config.Bucket, this.Config.Key)
	return err == nil
}

// 判断文件是否存在
func (this *OssFile) FileExist() (exist bool) {
	if this.Config.ForceOriginFile {
		exist = this.FileExistOss()
	} else {
		exist = this.FileExistFromLocal()
		if !exist {
			exist = this.FileExistOss()
		}
	}
	return exist
}

func zipfile(zipFilePath, filePath string) (err error) {

	newZipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zw := zip.NewWriter(newZipFile)
	defer zw.Close()

	fileToZip, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = info.Name()
	header.Method = zip.Deflate
	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)

	return nil
}
