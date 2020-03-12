package main

import (
	pbindex "GrpcCommon/index_service"
	pbface "GrpcCommon/xgface_service"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/net/context"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
)

var (
	stage int
)

func init() {
	flag.IntVar(&stage, "stage", 0, "execute stage")
}

type logInfo struct {
	ID        uint32
	ImagePath string
	ClusterID int32
	Time      int64
	Table     string
}

func logTableNames() (tables []string, err error) {
	rows, err := GormDb.Raw("select table_name from information_schema.tables where table_schema='xgface' and table_name like 'log_20';").Rows()
	if err != nil {
		glog.Errorln(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		rows.Scan(&name)
		tables = append(tables, name)
	}
	return
}

func getAllImageWithTime() (err error) {
	r := RedisClient.Get()
	tableNames, err := logTableNames()
	if err != nil {
		glog.Errorln("内部错误")
		return
	}
	glog.Info("日志表", tableNames)
	SQLFormat := "select %s.id, %s.name as image_path, %s.cluster_id, %s.time from %s left join %s on %s.id = %s.face_img_id order by time limit %d, %d;"
	var querySQL string
	for _, tableName := range tableNames {
		countSQL := "select count(*) from " + tableName + ";"
		count := 0
		row := GormDb.Raw(countSQL).Row()
		row.Scan(&count)
		glog.Infof("table %s has %d record", tableName, count)
		imgTableName := strings.Replace(tableName, "log_", "img_", -1)
		for i := 0; i < count; i += 1000 {
			querySQL = fmt.Sprintf(SQLFormat, tableName, imgTableName, tableName, tableName, tableName, imgTableName,
				imgTableName, tableName, i, 1000)
			rows, err := GormDb.Raw(querySQL).Rows()
			if err != nil {
				glog.Errorln(err)
				continue
			}
			num := 0
			for rows.Next() {
				rec := logInfo{
					Table: tableName,
				}
				rows.Scan(&rec.ID, &rec.ImagePath, &rec.ClusterID, &rec.Time)
				logString, err := json.Marshal(rec)
				if err != nil {
					glog.Errorln("编码失败", rec)
				}
				r.Send("LPUSH", "imageData", logString)
				num++
			}
			rows.Close()
			r.Flush()
			glog.Infof("table %s send %d message", tableName, num)
		}
	}
	return
}
func getBase64ImageFromFile(path string) (imageData string, err error) {
	realPath := Config.BasePath + strings.TrimPrefix(path, Config.ImgServer)
	data, err := ioutil.ReadFile(realPath)
	if err != nil {
		glog.Errorf("did not open file: %s\n", realPath)
		return
	}
	imageData = base64.StdEncoding.EncodeToString(data)
	return
}

func updateLogClusterID(table string, id uint32, newCID int32) error {
	err := GormDb.Table(table).Where("id = ?", id).Update("cluster_id", newCID).Error
	if err != nil {
		glog.Errorln(err)
		return err
	}
	return nil
}

//UpdateClusterID 更新index中的ID
func UpdateClusterID() (err error) {
	r := RedisClient.Get()
	xgfaceConn, err := GetXgFaceConn()
	if err != nil {
		glog.Errorln(err)
		return
	}
	indexConn, err := GetXgIndexConn()
	if err != nil {
		glog.Errorln(err)
		return
	}
	xgfaceClient := pbface.NewXgfaceServiceClient(xgfaceConn)
	indexClient := pbindex.NewIndexServiceClient(indexConn)
	count := 0
	for {
		value, err := redis.Bytes(r.Do("RPOP", "imageData"))
		if err != nil {
			glog.Errorln(err, "处理数量", count)
			break
		}
		rec := logInfo{}
		err = json.Unmarshal(value, &rec)
		if err != nil {
			glog.Errorln(err, "处理数量", count)
			break
		}
		Base64Content, err := getBase64ImageFromFile(rec.ImagePath)
		if err != nil {
			r.Do("LPUSH", "imageData", value)
			glog.Errorf("image %s failed and resend to redis", rec.ImagePath)
			continue
		}
		detectInfo, err := xgfaceClient.GetDetectInfo(context.Background(), &pbface.Request{Images: []string{Base64Content}})
		if err != nil {
			r.Do("LPUSH", "imageData", value)
			glog.Errorf("image %s failed and resend to redis", rec.ImagePath)
			continue
		}
		if len(detectInfo.Faces) == 0 {
			glog.Errorf("image %s failed detect no face", rec.ImagePath)
			err1 := updateLogClusterID(rec.Table, rec.ID, 0)
			if err1 != nil {
				glog.Errorln("更新失败", rec)
				glog.Errorln("newCID------------", 0)
			}
			continue
		}
		cluster, err := indexClient.InsertPoint(context.Background(), &pbindex.Feature{Value: detectInfo.Faces[0].Feature.Values})
		if err != nil {
			glog.Errorf("image %s insert index failed", rec.ImagePath)
			continue
		}
		if err = updateLogClusterID(rec.Table, rec.ID, cluster.Id); err != nil {
			glog.Errorln("更新失败", rec)
			glog.Errorln("newCID------------", cluster.Id)
		}
		count++
	}
	return
}

func main() {
	flag.Parse()
	defer glog.Flush()
	if stage == 0 {
		if err := getAllImageWithTime(); err != nil {
			glog.Errorln(err)
			return
		}
		if err := UpdateClusterID(); err != nil {
			glog.Errorln(err)
			return
		}
	} else if stage == 1 {
		if err := getAllImageWithTime(); err != nil {
			glog.Errorln(err)
			return
		}
	} else if stage == 2 {
		if err := UpdateClusterID(); err != nil {
			glog.Errorln(err)
			return
		}
	}
	return
}
