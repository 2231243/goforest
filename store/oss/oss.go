package oss

import (
	"io"
	"net/http"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/util/gconv"
)

type Client struct {
	Bucket *oss.Bucket
}

func (c *Client) Get(path string) (io.ReadCloser, error) {
	return c.Bucket.GetObject(path, nil)
}
func (c *Client) Put(path string, r io.Reader) error {
	return c.Bucket.PutObject(path, r)
}

func (c *Client) GetHeader(path string) (http.Header, error) {
	return c.Bucket.GetObjectMeta(path)
}

func (c *Client) SetObjectAcl(path string, acl string) error {
	return c.Bucket.SetObjectACL(path, oss.ACLType(acl))
}

func (c *Client) Copy(dest, src string) error {
	_, err := c.Bucket.CopyObject(src, dest)
	return err
}

func (c *Client) ContentLength(object string) (int64, error) {
	h, err := c.GetHeader(object)
	if err != nil {
		return 0, err
	}
	return gconv.Int64(h.Get("Content-Length")), nil
}
func (c *Client) PutObjectFromFile(objectKey string, filePath string) error {
	return c.Bucket.PutObjectFromFile(objectKey, filePath)
}

type OssOption struct {
	EndPoint, AccessKeyID, AccessKeySecret, Bucket string
}

func New(o OssOption) *Client {
	c, _ := oss.New(o.EndPoint, o.AccessKeyID, o.AccessKeySecret)
	bucket, _ := c.Bucket(o.Bucket)
	return &Client{Bucket: bucket}
}
