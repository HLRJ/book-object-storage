package es8

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
)

// Metadata 元数据结构体
type Metadata struct {
	Name    string //对象的名字
	Version int    //对象的版本
	Size    int64  //对象该版本的大小
	Hash    string //对象该版本的散列值
	// 如果某个版本是一个删除标记，其size为0，hash为空字符串
}

// 元数据
type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total struct {
			Value    int
			Relation string
		}
		Hits []hit
	}
}

// 通过对象名字和版本id拿到对象的元数据 并将其反序列化
func getMetadata(name string, versionId int) (meta Metadata, err error) {
	// url地址为我们环境设置的ES_SERVER,  然后将其拼接成一个带文件地址的请求
	url := fmt.Sprintf("http://%s/metadata/_doc/%s_%d/_source", os.Getenv("ES_SERVER"), name, versionId)
	// 请求会返回一个结果
	result, err := http.Get(url)
	if err != nil {
		return
	}
	if result.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to get %s_%d:%d", name, versionId, result.StatusCode)
		return
	}
	// 将请求返回的结果的body读出来反序列化
	result2, _ := ioutil.ReadAll(result.Body)
	// 写入meat中
	err = json.Unmarshal(result2, &meta)
	if err != nil {
		return Metadata{}, err
	}
	return
}

//通过对象的名字拿到  最新版本的对象文件

func SearchLatestVersion(name string) (meta Metadata, err error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=name:%s&size=1&sort=version:desc", os.Getenv("ES_SERVER"), url2.PathEscape(name))
	result, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	if result.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to search latest metadata:%s", result.StatusCode)
		return
	}
	result2, err := ioutil.ReadAll(result.Body)
	if err != nil {
		fmt.Println(err)
	}
	var sr searchResult
	err = json.Unmarshal(result2, &sr)
	if err != nil {
		return Metadata{}, err
	}
	// 如果 返回的切片长度不为0，说明取值成功
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}
	// 如果 取值失败 就返回空白
	return
}

// GetMetadata 根据对象的名字和版本号，如果版本号为0返回最新的元数据信息，否则调用内部函数getMetadata
func GetMetadata(name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

// 向ES服务上传一个新的元数据

func PutMetadata(name string, version int, size int64, hash string) error {
	document := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash":"%s"}`, name, version, size, hash)
	client := http.Client{}
	//op_type=create参数   如果同时有多个客户端上传同一个元数据，结果会发生冲突，只有第一个文档被成功创建。之后的PUT请求会返回409
	url := fmt.Sprintf("http://%s/metadata/_doc/%s_%d?op_type=create", os.Getenv("ES_SERVER"), name, version)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(document))
	request.Header.Set("Content-Type", "application/json")
	result, err := client.Do(request)
	if err != nil {
		return err
	}
	// 如果响应的状态码为冲突，就递归调用自己，给版本加上一
	if result.StatusCode == http.StatusConflict {
		return PutMetadata(name, version+1, size, hash)
	}
	// 如果创建不成功，返回错误
	if result.StatusCode != http.StatusCreated {
		result2, _ := ioutil.ReadAll(result.Body)
		return fmt.Errorf("fail to put metadata:%d %s", result.StatusCode, string(result2))
	}
	return nil
}

// AddVersion  首先调用SearchLatestVersion 获取对象最新的版本，然后在版本号上加1调用PutMetadata
func AddVersion(name, hash string, size int64) error {
	version, err := SearchLatestVersion(name)
	if err != nil {
		return err
	}
	return PutMetadata(name, version.Version+1, size, hash)
}

// 用于搜索某个对象或所有对象的全部版本， 返回元数据切片
func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?sort=name,version&from=%d&size=%d", os.Getenv("ES_SERVER"), from, size)
	//如果name不为空字符串，就搜索指定对象的所有版本，否则搜索所有对象的所有版本
	if name != "" {
		url += "&q=name:" + name
	}
	result, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	metas := make([]Metadata, 0)
	result2, _ := ioutil.ReadAll(result.Body)
	var sr searchResult
	err = json.Unmarshal(result2, &sr)
	if err != nil {
		return nil, err
	}
	// 遍历取元数据 存入切片
	for i := range sr.Hits.Hits {
		metas = append(metas, sr.Hits.Hits[i].Source)
	}
	return metas, nil
}

// DelMetadata 删除指定版本的对象的元数据
func DelMetadata(name string, version int) {
	url := fmt.Sprintf("http://%s/metadata/_doc/%s_%d", os.Getenv("ES_SERVER"), name, version)
	client := http.Client{}
	request, _ := http.NewRequest("DELETE", url, nil)
	client.Do(request)
}
