package goss

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestGettObject(t *testing.T) {

	client := NewS3ClientWithStaticCredentials(
		"http://192.168.10.180:9000",
		"ABCDEFGHIJKLMN",
		"abcdefghijklmn",
		"")

	testListBucket(t, client)

	testGetObject(t, client)

	testPutObject(t, client)

	testCopyObject(t, client)
}

func testListBucket(t *testing.T, client *S3Client) {
	// 获取桶信息
	outputbucket, err := client.ListBucket()

	if err != nil {
		t.Error(err)
		return
	}

	t.Log("ListBucket返回的结果：", outputbucket.Buckets)
}

func testGetObject(t *testing.T, client *S3Client) {
	// 获取桶里面的对象信息
	output, err := client.GetObject("test", "test")
	if err != nil {
		t.Error(err)
		return
	}

	body, err := ioutil.ReadAll(output.Body)

	t.Log("GetObject返回的结果：", string(body))

	metadata, _ := client.GetObjectMetadata("test", "test")
	t.Log("GetObject返回的metadata结果：", metadata)
}

func testPutObject(t *testing.T, client *S3Client) {

	sample := &struct {
		TestString string
		TestInt    int
	}{
		TestString: "xxxx",
		TestInt:    12,
	}

	s, _ := json.Marshal(sample)
	input := bytes.NewReader(s)

	metadata := make(map[string]string)
	metadata["version"] = "xxxxxxxxx"

	// 获取桶里面的对象信息
	err := client.PutObject("test", "test", input, metadata)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("PutObject的值：%+v", sample)
}

func testCopyObject(t *testing.T, client *S3Client) {
	source := "test"
	target := "test_copy"
	// 获取桶里面的对象信息
	err := client.CopyObject("test", "test"+"/"+source, target)
	if err != nil {
		t.Error(err)
		return
	}

	target = "test_same_copy"
	err = client.CopyObjectInSameBucket("test", source, target)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("CopyObject结果成功, 源文件：%s, 目标文件: %s", source, target)
}
